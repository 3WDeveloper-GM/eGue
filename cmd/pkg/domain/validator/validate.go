package Validator

import (
	"errors"
	"regexp"
)

type Validator struct {
	errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (c *Validator) In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func (v *Validator) Errors() []error {
	errorArr := make([]error, 0)
	for key, value := range v.errors {
		errorArr = append(errorArr, errors.New(key+":"+value))
	}
	return errorArr
}
