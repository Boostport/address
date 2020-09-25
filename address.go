//go:generate go run generator/generate.go

// Package address is a library that validates and formats addresses using data generated from Google's Address
// Data Service.
package address

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/collate"
	textLanguage "golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type formatData struct {
	Country                     string
	CountryEnglish              string
	Name                        string
	Organization                string
	StreetAddress               []string
	DependentLocality           string
	Locality                    string
	AdministrativeArea          string
	AdministrativeAreaPostalKey string
	PostCode                    string
	SortingCode                 string
}

// Address represents a valid address made up of its child components.
type Address struct {
	Country            string
	Name               string
	Organization       string
	StreetAddress      []string
	DependentLocality  string
	Locality           string
	AdministrativeArea string
	PostCode           string
	SortingCode        string
}

// IsZero reports whether a represents a zero/uninitialized address
func (a Address) IsZero() bool {
	return a.Country == "" && a.Name == "" && a.Organization == "" && len(a.StreetAddress) <= 0 && a.DependentLocality == "" && a.Locality == "" && a.AdministrativeArea == "" && a.PostCode == "" && a.SortingCode == ""
}

func (a Address) toFormatData(countryData country, language string) formatData {

	f := formatData{
		Name:               a.Name,
		Organization:       a.Organization,
		StreetAddress:      []string{},
		DependentLocality:  a.DependentLocality,
		Locality:           a.Locality,
		AdministrativeArea: a.AdministrativeArea,
		PostCode:           a.PostCode,
		SortingCode:        a.SortingCode,
	}

	for _, addressLine := range a.StreetAddress {
		trimmed := strings.TrimSpace(addressLine)

		if trimmed != "" {
			f.StreetAddress = append(f.StreetAddress, trimmed)
		}
	}

	normalizedLanguage := generated.normalizeLanguage(countryData.ID, language)

	namer := display.Regions(textLanguage.MustParse(normalizedLanguage))

	country := textLanguage.MustParseRegion(countryData.ID)

	if namer != nil {
		f.Country = namer.Name(country)
	} else {
		namer = display.Regions(textLanguage.English)
		f.Country = namer.Name(country)
	}

	if a.AdministrativeArea != "" {
		if adminAreaName := generated.getAdministrativeAreaName(a.Country, a.AdministrativeArea, language); adminAreaName != "" {
			f.AdministrativeArea = adminAreaName
		}

		f.AdministrativeAreaPostalKey = generated.getAdministrativeAreaPostalKey(a.Country, a.AdministrativeArea)
	}

	if a.Locality != "" {
		if localityName := generated.getLocalityName(a.Country, a.AdministrativeArea, a.Locality, language); localityName != "" {
			f.Locality = localityName
		}
	}

	if a.DependentLocality != "" {
		if dependentLocalityName := generated.getDependentLocalityName(a.Country, a.AdministrativeArea, a.Locality, a.DependentLocality, language); dependentLocalityName != "" {
			f.DependentLocality = dependentLocalityName
		}
	}

	return f
}

// NewValid creates a new Address. If the address is invalid, an error is returned.
// In the case where an error is returned, the error is a hashicorp/go-multierror (https://github.com/hashicorp/go-multierror).
// You can use a type switch to get a list of validation errors for the address.
func NewValid(fields ...func(*Address)) (Address, error) {

	address := New(fields...)

	err := Validate(address)

	if err != nil {
		return address, fmt.Errorf("invalid address: %s", err)
	}

	return address, nil
}

// New creates a new unvalidated address. The validity of the address should be checked
// using the validator.
func New(fields ...func(*Address)) Address {

	address := Address{}

	for _, field := range fields {
		field(&address)
	}

	return address
}

// WithCountry sets the country code of an address.
// The country code must be an ISO 3166-1 country code.
func WithCountry(countryCode string) func(*Address) {
	return func(a *Address) {
		a.Country = strings.ToUpper(countryCode)
	}
}

// WithName sets the addressee's name of an address.
func WithName(name string) func(*Address) {
	return func(a *Address) {
		a.Name = name
	}
}

// WithOrganization sets the addressee's organization of an address.
func WithOrganization(organization string) func(*Address) {
	return func(a *Address) {
		a.Organization = organization
	}
}

// WithStreetAddress sets the street address of an address.
// The street address is a slice of strings, with each element representing an address line.
func WithStreetAddress(streetAddress []string) func(*Address) {
	return func(a *Address) {
		a.StreetAddress = streetAddress
	}
}

// WithDependentLocality sets the dependent locality (commonly known as the suburb) of an address.
// If the country of the address has a list of dependent localities, then the key of the dependent locality should
// be used, otherwise, the validation will fail.
func WithDependentLocality(dependentLocality string) func(*Address) {
	return func(a *Address) {
		a.DependentLocality = dependentLocality
	}
}

// WithLocality sets the locality (commonly known as the city) of an address.
// If the country of the address has a list of localities, then the key of the locality should be used, otherwise,
// the validation will fail.
func WithLocality(locality string) func(*Address) {
	return func(a *Address) {
		a.Locality = locality
	}
}

// WithAdministrativeArea sets the administrative area (commonly known as the state) of an address.
// If the country of the address has a list of administrative area, then the key of the administrative area should
// used, otherwise, the validation will fail.
func WithAdministrativeArea(administrativeArea string) func(*Address) {
	return func(a *Address) {
		a.AdministrativeArea = administrativeArea
	}
}

// WithPostCode sets the post code of an address.
func WithPostCode(postCode string) func(*Address) {
	return func(a *Address) {
		a.PostCode = postCode
	}
}

// WithSortingCode sets the sorting code of an address.
func WithSortingCode(sortingCode string) func(*Address) {
	return func(a *Address) {
		a.SortingCode = sortingCode
	}
}

// CountryData contains the address data for a country.
// The AdministrativeAreas field contains a list of nested subdivisions (administrative areas, localities and dependent
// localities) grouped by their translated languages. They are also sorted according to the sort order of the languages
// they are in.
type CountryData struct {
	Format                     string
	LatinizedFormat            string
	Required                   []Field
	Allowed                    []Field
	DefaultLanguage            string
	AdministrativeAreaNameType FieldName
	LocalityNameType           FieldName
	DependentLocalityNameType  FieldName
	PostCodeNameType           FieldName
	PostCodeRegex              PostCodeRegexData
	AdministrativeAreas        map[string][]AdministrativeAreaData
}

// PostCodeRegexData contains regular expressions for validating post codes for a given country.
// If the country has subdivisions (administrative areas, localities and dependent localities), the SubdivisionRegex
// field may contain further regular expressions to Validate the post code.
type PostCodeRegexData struct {
	Regex            string
	SubdivisionRegex map[string]PostCodeRegexData
}

// AdministrativeAreaData contains the name and ID of and administrative area. The ID must be passed to
// WithAdministrativeArea() when creating an address. The name is useful for displaying to the end user.
type AdministrativeAreaData struct {
	ID   string
	Name string

	Localities []LocalityData
}

// LocalityData contains the name and ID of and administrative area. The ID must be passed to
// WithLocalityData() when creating an address. The name is useful for displaying to the end user.
type LocalityData struct {
	ID   string
	Name string

	DependentLocalities []DependentLocalityData
}

// DependentLocalityData contains the name and ID of and administrative area. The ID must be passed to
// WithDependentLocalityData() when creating an address. The name is useful for displaying to the end user.
type DependentLocalityData struct {
	ID   string
	Name string
}

// CountryList contains a list of countries that can be used to create addresses.
type CountryList []CountryListItem

// Len returns the number of countries in the list. This is used for sorting the countries and would not generally be used
// in client code.
func (c CountryList) Len() int {
	return len(c)
}

// Swap swaps 2 countries in the list. This is used for sorting the countries and would not generally be used
// in client code.
func (c CountryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Bytes returns a country name in bytes. This is used for sorting the countries and would not generally be used
// in client code.
func (c CountryList) Bytes(i int) []byte {
	return []byte(c[i].Name)
}

// CountryListItem represents a single country, containing the ISO 3166-1 code and the name of the country.
type CountryListItem struct {
	Code string
	Name string
}

// ListCountries returns a list of countries that can be used to create addresses.
// Language must be a valid ISO 639-1 language code such as: en, jp, zh, etc.
// If the language does not have any translations or is invalid, then English is used as the fallback language.
// The returned list of countries is sorted according to the chosing language.
func ListCountries(language string) []CountryListItem {

	l, err := textLanguage.Parse(language)

	if err != nil {
		l = textLanguage.English
	}

	c := collate.New(l)
	namer := display.Regions(l)

	if namer == nil {
		namer = display.Regions(textLanguage.English)
	}

	var countries CountryList

	for countryCode := range generated {

		if countryCode == "ZZ" {
			continue
		}

		country := textLanguage.MustParseRegion(countryCode)

		countries = append(countries, CountryListItem{
			Code: countryCode,
			Name: namer.Name(country),
		})
	}

	c.Sort(countries)

	return countries
}

// GetCountry returns address information for a given country.
func GetCountry(countryCode string) CountryData {

	country := generated.getCountry(countryCode)

	return internalCountryDataToCountryData(country)
}

func internalCountryDataToCountryData(country country) CountryData {

	data := CountryData{
		Format:                     country.Format,
		LatinizedFormat:            country.LatinizedFormat,
		DefaultLanguage:            country.DefaultLanguage,
		AdministrativeAreaNameType: country.AdministrativeAreaNameType,
		LocalityNameType:           country.LocalityNameType,
		DependentLocalityNameType:  country.DependentLocalityNameType,
		PostCodeNameType:           country.PostCodeNameType,
		PostCodeRegex:              internalPostCodeRegexToPostCodeRegexData(country.PostCodeRegex),
	}

	var required []Field

	for field := range country.RequiredFields {
		required = append(required, field)
	}

	sort.Slice(required, func(i, j int) bool {
		return required[i].String() < required[j].String()
	})

	data.Required = required

	var allowed []Field

	for field := range country.AllowedFields {
		allowed = append(allowed, field)
	}

	sort.Slice(allowed, func(i, j int) bool {
		return allowed[i].String() < allowed[j].String()
	})

	data.Allowed = allowed

	administrativeAreas := map[string][]AdministrativeAreaData{}

	for lang, adminAreas := range country.AdministrativeAreas {
		administrativeAreas[lang] = internalAdministrativeAreasToAdministrativeAreaData(adminAreas)
	}

	if len(administrativeAreas) > 0 {
		data.AdministrativeAreas = administrativeAreas
	}

	return data
}

func internalPostCodeRegexToPostCodeRegexData(regex postCodeRegex) PostCodeRegexData {

	result := PostCodeRegexData{
		Regex: regex.regex,
	}

	for subID, regex := range regex.subdivisionRegex {

		if result.SubdivisionRegex == nil {
			result.SubdivisionRegex = map[string]PostCodeRegexData{}
		}

		result.SubdivisionRegex[subID] = internalPostCodeRegexToPostCodeRegexData(regex)
	}

	return result
}

func internalAdministrativeAreasToAdministrativeAreaData(areas []administrativeArea) []AdministrativeAreaData {

	var result []AdministrativeAreaData

	for _, adminArea := range areas {

		var localities []LocalityData

		for _, locality := range adminArea.Localities {

			var dependentLocalities []DependentLocalityData

			for _, dependentLocality := range locality.DependentLocalities {
				dependentLocalities = append(dependentLocalities, DependentLocalityData{
					ID:   dependentLocality.ID,
					Name: dependentLocality.Name,
				})
			}

			localityData := LocalityData{
				ID:   locality.ID,
				Name: locality.Name,
			}

			if len(dependentLocalities) > 0 {
				localityData.DependentLocalities = dependentLocalities
			}

			localities = append(localities, localityData)
		}

		adminAreaData := AdministrativeAreaData{
			ID:   adminArea.ID,
			Name: adminArea.Name,
		}

		if len(localities) > 0 {
			adminAreaData.Localities = localities
		}

		result = append(result, adminAreaData)
	}

	return result
}
