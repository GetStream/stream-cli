package validator

import (
	"io"
	"regexp"

	streamchat "github.com/GetStream/stream-chat-go/v5"
)

type Results struct {
	Stats  map[string]int
	Errors []error
}

func (r *Results) HasErrors() bool {
	return len(r.Errors) > 0
}

func newResults(stats map[string]int, errs *multiError) *Results {
	return &Results{
		Stats:  stats,
		Errors: errs.errors,
	}
}

type Options func(*option)

type option struct {
	lighterValidationChannelID bool
}

func LighterValidationChannelID() Options {
	return func(o *option) {
		o.lighterValidationChannelID = true
	}
}

type Validator struct {
	decoder *Decoder
	index   *index
}

func New(r io.ReadSeeker, roles []*streamchat.Role, channelTypes channelTypeMap, options ...Options) *Validator {
	roleMap := make(roleMap, len(roles))
	for _, role := range roles {
		roleMap[role.Name] = role
	}

	opt := &option{}
	for _, o := range options {
		o(opt)
	}

	if opt.lighterValidationChannelID {
		// allows to pass ! in channel ID
		validChannelIDRe = regexp.MustCompile(`^[\w!-]*$`)
	} else {
		validChannelIDRe = regexp.MustCompile(`^[\w-]*$`)
	}

	return &Validator{
		decoder: NewDecoder(r),
		index:   newIndex(roleMap, channelTypes),
	}
}

func (v *Validator) Validate() *Results {
	errs := new(multiError)

	// first pass: validate item field and index items for reference validation
	errs.add(v.decoder.Items(func(item Item) error {
		if err := item.validateFields(); err != nil {
			return newValidationError(err)
		}
		return item.index(v.index)
	}))
	if errs.len() > maxNrOfErrors {
		return newResults(v.index.stats(), errs)
	}

	// second pass: reference validation
	errs.add(v.decoder.Items(func(item Item) error {
		if err := item.validateReferences(v.index); err != nil {
			return newReferenceError(err)
		}
		return nil
	}))
	if errs.len() > maxNrOfErrors {
		return newResults(v.index.stats(), errs)
	}

	// last pass: validate any channel members that are unaccounted for
	errs.add(v.index.validateChannelMembers())

	return newResults(v.index.stats(), errs)
}
