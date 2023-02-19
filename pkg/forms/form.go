package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) *Form {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s cannot be blank", value))
		}
	}
	return f
}

func (f *Form) MaxLength(field string, maxLen int) *Form {
	value := f.Get(field)

	if value == "" {
		return f
	}

	if utf8.RuneCountInString(value) > maxLen {
		f.Errors.Add(field, fmt.Sprintf("%s cannot be more than maximum character length (%d)", value, maxLen))
	}

	return f
}

func (f *Form) PermittedValues(field string, opts ...string) *Form {
	value := f.Get(field)
	if value == "" {
		return f
	}
	for _, opt := range opts {
		if value == opt {
			return f
		}
	}
	f.Errors.Add(field, "This field is invalid")
	return f
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
