package validator

import (
	"encoding/json"
	"errors"
	"io"
)

type Decoder struct {
	source io.ReadSeeker
}

func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{source: r}
}

func (d *Decoder) Items(fn func(item Item) error) error {
	// Reset source
	if _, err := d.source.Seek(0, io.SeekStart); err != nil {
		return newParseError(err)
	}

	jd := json.NewDecoder(d.source)

	// Validate opening-token
	token, err := jd.Token()
	if err != nil {
		return newParseError(err)
	}
	if token != json.Delim('[') {
		return newParseError(errors.New("invalid format"))
	}

	// Decode items
	errs := new(multiError)
	for jd.More() {
		offset := jd.InputOffset()

		var v rawItem
		if err := jd.Decode(&v); err != nil {
			return err
		}

		item, err := newItem(&v)
		if err != nil {
			errs.add(newItemError(v, offset, newParseError(err)))
			continue
		}

		err = fn(item)
		if err != nil {
			errs.add(newItemError(v, offset, err))
		}
	}
	if errs.hasErrors() {
		return errs
	}
	return nil
}
