package address

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidCountryCode indicate that the country code used to create an address is invalid.
var ErrInvalidCountryCode = errors.New("invalid country code")

// ErrInvalidDependentLocality indicates that the dependent locality is invalid. This is usually due to the country having
// a pre-determined list of dependent localities and the value does not match any of the keys in the list of dependent localities.
var ErrInvalidDependentLocality = errors.New("invalid dependent locality")

// ErrInvalidLocality indicates that the locality is invalid. This is usually due to the country having
// a pre-determined list of localities and the value does not match any of the keys in the list of localities.
var ErrInvalidLocality = errors.New("invalid locality")

// ErrInvalidAdministrativeArea indicates that the administrative area is invalid. This is usually due to the country having
// a pre-determined list of administrative areas and the value does not match any of the keys in the list of administrative areas.
var ErrInvalidAdministrativeArea = errors.New("invalid administrative area")

// ErrInvalidPostCode indicates that the post code did not valid using the regular expressions of the country.
var ErrInvalidPostCode = errors.New("invalid post code")

// ErrMissingRequiredFields indicates the a required address field is missing. The Fields field can be used to get a list
// of missing fields.
type ErrMissingRequiredFields struct {
	country string
	Fields  []Field
}

func (e ErrMissingRequiredFields) Error() string {

	var fieldsStr []string

	for _, field := range e.Fields {
		fieldsStr = append(fieldsStr, field.String())
	}

	return fmt.Sprintf("missing required fields for %s: %s", e.country, strings.Join(fieldsStr, ","))
}

// ErrUnsupportedFields indicates that an address field as provided, but it is not supported by the address format
// of the country. The Fields field can be used to get a list of unsupported fields.
type ErrUnsupportedFields struct {
	country string
	Fields  []Field
}

func (e ErrUnsupportedFields) Error() string {

	var fieldsStr []string

	for _, field := range e.Fields {
		fieldsStr = append(fieldsStr, field.String())
	}

	return fmt.Sprintf("unsupported fields for %s: %s", e.country, strings.Join(fieldsStr, ","))
}
