package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type varValue struct {
	value    interface{}
	variable variableDescriptor
}

func parseLogs(logs []byte, format formatDescriptor) []map[string]varValue {
	out := []map[string]varValue{}
	logRecords := format.recordDelimiterPattern.FindAllIndex(logs, -1)

	if logRecords == nil {
		return out
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

		varValue, err := parseLogRecord(logs[recordOffset[0]:recordEndIndex], format)

		if varValue == nil {
			if err != nil {
				logrus.Errorf("Error pasing record at %v: %s\n", recordOffset[0], err.Error())
			}

			continue
		}

		out = append(out, varValue)
	}

	return out
}

func parseLogRecord(record []byte, format formatDescriptor) (map[string]varValue, error) {
	matchedStrings := parseLogRecordRaw(record, format.recordPattern)
	if matchedStrings == nil {
		return nil, errors.New("record was not matched")
	}

	result := map[string]varValue{}

	for _, variable := range format.variables {
		rawValue := matchedStrings[variable.name]
		value, err := constructValue(rawValue, variable)
		if err != nil {
			return nil, err
		}

		result[variable.name] = value
	}

	return result, nil
}

func constructValue(rawValue string, varDescr variableDescriptor) (varValue, error) {
	switch varDescr.varType {
	case stringVarType:
		return varValue{rawValue, varDescr}, nil

	case boolVarType:
		val, err := strconv.ParseBool(rawValue)
		if err != nil {
			return varValue{}, err
		}

		return varValue{val, varDescr}, nil

	case intVarType:
		val, err := strconv.Atoi(rawValue)
		if err != nil {
			return varValue{}, err
		}

		return varValue{val, varDescr}, nil

	case timeVarType:
		val, err := time.Parse(varDescr.layout, rawValue)
		if err != nil {
			return varValue{}, err
		}

		return varValue{val, varDescr}, nil

	case durationVarType:
		val, err := time.ParseDuration(rawValue)
		if err != nil {
			return varValue{}, err
		}

		return varValue{val, varDescr}, nil
	default:
		panic(fmt.Sprintf("Unknown variable type: %v", varDescr.varType))
	}
}

func parseLogRecordRaw(record []byte, recordPattren *regexp.Regexp) (result map[string]string) {
	match := recordPattren.FindSubmatch(record)
	if match == nil {
		return nil
	}

	result = map[string]string{}

	for subIndex, subName := range recordPattren.SubexpNames() {
		result[subName] = string(match[subIndex])
	}

	return result
}
