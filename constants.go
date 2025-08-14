//go:generate go tool stringer -type=Field,FieldName -output=constant_string.go
package address

// Field is an address field type.
type Field int

const (
	Country Field = iota + 1
	Name
	Organization
	StreetAddress
	DependentLocality
	Locality
	AdministrativeArea
	PostCode
	SortingCode
)

// Key returns the corresponding one-letter abbreviation used by Google to refer to address fields.
// This is useful for parsing the address format for a country.
// See https://github.com/googlei18n/libaddressinput/wiki/AddressValidationMetadata for more information.
func (i Field) Key() string {

	switch i {
	case Country:
		return "country"
	case Name:
		return "N"
	case Organization:
		return "O"
	case StreetAddress:
		return "A"
	case DependentLocality:
		return "D"
	case Locality:
		return "C"
	case AdministrativeArea:
		return "S"
	case PostCode:
		return "Z"
	case SortingCode:
		return "X"
	}

	return ""
}

// FieldName is the name to be used when referring to a field.
// For example, in India, the post code is called PIN Code instead of Post Code.
// The field name allows you to display the appropriate form labels to the user.
type FieldName int

const (
	Area FieldName = iota + 1
	City
	County
	Department
	District
	DoSi
	Eircode
	Emirate
	Island
	Neighborhood
	Oblast
	PINCode
	Parish
	PostTown
	PostalCode
	Prefecture
	Province
	State
	Suburb
	Townland
	VillageTownship
	ZipCode
)
