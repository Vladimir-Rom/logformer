package main

import (
	"errors"
	"regexp"

	"gopkg.in/yaml.v2"
)

type formatDescriptor struct {
	recordDelimiterPattern regexp.Regexp
	recordPattern          regexp.Regexp
	variables              []variableDescriptor
}

type variableType int

const (
	unknownVarType variableType = iota
	stringVarType
	boolVarType
	intVarType
	timeVarType
	durationVarType
)

type variableDescriptor struct {
	name    string
	varType variableType
	layout  string
}

func getFormatDescriptor(yamlContent []byte) (result formatDescriptor, err error) {
	type formatDescriptorRawType struct {
		RecordDelimiterPattern string `yaml:"recordDelimiterPattern"`
		RecordPattern          string `yaml:"recordPattern"`
		Variables              []struct {
			Name   string `yaml:"name"`
			Type   string `yaml:"type"`
			Layout string `yaml:"layout,omitempty"`
		} `yaml:"variables"`
	}

	var formatDescriptorRaw formatDescriptorRawType

	err = yaml.Unmarshal(yamlContent, &formatDescriptorRaw)
	if err != nil {
		return result, err
	}

	recordDelimiterPattern, err := regexp.Compile(formatDescriptorRaw.RecordDelimiterPattern)
	if err != nil {
		return result, err
	}

	result.recordDelimiterPattern = *recordDelimiterPattern

	recordPattern, err := regexp.Compile(formatDescriptorRaw.RecordPattern)
	if err != nil {
		return result, err
	}

	result.recordPattern = *recordPattern

	result.variables = make([]variableDescriptor, len(formatDescriptorRaw.Variables))

	for i, v := range formatDescriptorRaw.Variables {
		result.variables[i].varType, err = parseVarType(v.Type)
		if err != nil {
			return result, err
		}

		result.variables[i].name = v.Name
		result.variables[i].layout = v.Layout
	}

	if err = validateVariables(result); err != nil {
		return result, err
	}

	return result, nil
}

func validateVariables(result formatDescriptor) error {
	regexGroups := result.recordPattern.SubexpNames()

	isVariableDefined := func(variable string) bool {
		for _, group := range regexGroups {
			if variable == group {
				return true
			}
		}

		return false
	}

	for _, v := range result.variables {
		if !isVariableDefined(v.name) {
			return errors.New("Undefined variable: " + v.name)
		}
	}

	return nil
}

func parseVarType(varType string) (variableType, error) {
	switch varType {
	case "string":
		return stringVarType, nil
	case "bool", "boolean":
		return boolVarType, nil
	case "int":
		return intVarType, nil
	case "time":
		return timeVarType, nil
	case "duration":
		return durationVarType, nil
	default:
		return unknownVarType, errors.New("Unknown variable type: " + varType)
	}
}
