package main

import (
	"regexp"
	"sync"
)

func parallelParseLogs(levelOfParallelism int, logs []byte, format formatDescriptor) ([]map[string]varValue, error) {
	chunks := splitLogs(logs, levelOfParallelism, format.recordDelimiterPattern)
	chunks = append(chunks, len(logs))
	parsedChunks := make([][]map[string]varValue, len(chunks))

	wg := sync.WaitGroup{}
	wg.Add(len(chunks))

	chunkStart := 0
	for i, chunk := range chunks {
		chunkIndex := i
		chunkStartCopy := chunkStart
		chunkCopy := chunk
		go func() {
			defer wg.Done()
			parsedChunks[chunkIndex] = parseLogs(logs[chunkStartCopy:chunkCopy], format)
		}()

		chunkStart = chunkCopy
	}
	wg.Wait()
	return mergeSlices(parsedChunks), nil
}

func mergeSlices(parsedChunks [][]map[string]varValue) []map[string]varValue {
	var resultLen int
	for _, parsedChunk := range parsedChunks {
		resultLen += len(parsedChunk)
	}

	result := make([]map[string]varValue, resultLen)
	resultIndex := 0

	for _, parsedChunk := range parsedChunks {
		for _, item := range parsedChunk {
			result[resultIndex] = item
			resultIndex++
		}
	}

	return result
}

func splitLogs(logs []byte, chunksCount int, delimiter *regexp.Regexp) []int {
	result := []int{}
	stepSize := len(logs) / (chunksCount + 1)
	if stepSize == 0 {
		return result
	}

	currentIndex := 0

	for currentIndex < len(logs) {
		chunkStart := delimiter.FindIndex(logs[currentIndex:])
		if chunkStart == nil {
			break
		}

		result = append(result, currentIndex+chunkStart[0])
		currentIndex = currentIndex + chunkStart[0] + stepSize
	}

	return result
}
