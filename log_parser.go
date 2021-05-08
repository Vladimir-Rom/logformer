package main

type varValue struct {
	value    interface{}
	variable variableDescriptor
}

func parseLogs(logs []byte, format formatDescriptor, out chan<- map[string]varValue) error {
	logRecords := format.recordPattern.FindAllIndex(logs, -1)

	if logRecords == nil {
		return nil
	}

	recordsCount := len(logRecords)

	for recordIndex, recordOffset := range logRecords {
		isLastRecord := recordIndex+1 == recordsCount
		var recordEndIndex int
		if isLastRecord {
			recordEndIndex = len(logs)
		} else {
			recordEndIndex = logRecords[recordIndex+1][0]
		}

		varValue := parseLogRecord(logs[recordOffset[0]:recordEndIndex], format)

		if varValue == nil {
			continue
		}

		out <- varValue
	}

	return nil
}

func parseLogRecord(record []byte, format formatDescriptor) map[string]varValue {
	return nil, nil
}
