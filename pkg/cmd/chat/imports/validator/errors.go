package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const maxNrOfErrors = 25

type multiError struct {
	errors []error
}

func (e multiError) Error() string {
	msgs := make([]string, 0, len(e.errors))
	for i := range e.errors {
		msgs = append(msgs, e.errors[i].Error())
	}
	return strings.Join(msgs, ",")
}

func (e *multiError) add(err error) {
	if err == nil {
		return
	}

	var errs *multiError
	if errors.As(err, &errs) {
		for _, err := range errs.errors {
			e.add(err)
		}
		return
	}

	e.errors = append(e.errors, err)
}

func (e multiError) hasErrors() bool {
	return e.len() > 0
}

func (e multiError) len() int {
	return len(e.errors)
}

type ItemError struct {
	item   rawItem
	offset int64
	err    error
}

func (e ItemError) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%s (%d)", e.Error(), e.Offset()))
}

func newItemError(item rawItem, offset int64, err error) *ItemError {
	return &ItemError{
		item:   item,
		offset: offset,
		err:    err,
	}
}

func (e *ItemError) Error() string {
	return e.err.Error()
}

func (e *ItemError) Offset() int64 {
	return e.offset
}

type WrapError struct {
	err error
}

func (e *WrapError) Error() string {
	return e.err.Error()
}

func (e WrapError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.err.Error())
}

func newValidationError(err error) error {
	return &WrapError{fmt.Errorf("validation error: %w", err)}
}

func newReferenceError(err error) error {
	return &WrapError{fmt.Errorf("reference error: %w", err)}
}

func newParseError(err error) error {
	return &WrapError{fmt.Errorf("parse error: %w", err)}
}
