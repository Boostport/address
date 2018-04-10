package address

import (
	"regexp"
	"strconv"
	"strings"
)

// Zone is a list of territories. Zones are useful for determining whether an address is in a country, administrative area
// or a range of post codes to determine shipping or tax rates.
type Zone []Territory

// Contains checks to see an address is in a zone.
func (z Zone) Contains(address Address) bool {
	for _, zone := range z {
		if zone.contains(address) {
			return true
		}
	}

	return false
}

// Territory is a rule within a zone.
// It is able to match the following fields:
// - Country: An ISO 3166-1 country code
// - AdministrativeArea: The administrative area name. If the country has a list of pre-defined administrative areas,
// use the key of the administrative area.
// - Locality: The locality name. If the country has a list of pre-defined localities, use the key of the locality.
// - DependentLocality: The dependent locality name. If the country has a list of pre-defined dependent localities,
// use the key of the locality.
// - IncludedPostCodes: A PostCodeMatcher that includes the address within the territory if the post code matches.
// - ExcludedPostCodes: A PostCodeMatcher that excludes the address from the territory if the post code matches.
type Territory struct {
	Country            string
	AdministrativeArea string
	Locality           string
	DependentLocality  string
	IncludedPostCodes  PostCodeMatcher
	ExcludedPostCodes  PostCodeMatcher
}

// PostCodeMatcher returns a boolean signaling whether the post code matched or not.
type PostCodeMatcher interface {
	Match(postcode string) bool
}

func (t Territory) contains(address Address) bool {

	if t.Country != "" && address.Country != t.Country {
		return false
	}

	if t.AdministrativeArea != "" && strings.ToLower(address.AdministrativeArea) != strings.ToLower(t.AdministrativeArea) {
		return false
	}

	if t.Locality != "" && strings.ToLower(address.Locality) != strings.ToLower(t.Locality) {
		return false
	}

	if t.DependentLocality != "" && strings.ToLower(address.DependentLocality) != strings.ToLower(t.DependentLocality) {
		return false
	}

	matchIncluded := true

	if t.IncludedPostCodes != nil {

		matchIncluded = t.IncludedPostCodes.Match(address.PostCode)
	}

	matchExcluded := false

	if t.ExcludedPostCodes != nil {
		matchExcluded = t.ExcludedPostCodes.Match(address.PostCode)
	}

	return matchIncluded && !matchExcluded
}

// ExactMatcher matches post codes exactly in the list defined in the Matches field. If the post code is numeric,
// it's also possible to define a slice of ranges using the Ranges fiel.d
type ExactMatcher struct {
	Matches []string
	Ranges  []PostCodeRange
}

// Match checks to see if the post code matches a post code defined in Matches or if it is within a range defined in Ranges.
func (m ExactMatcher) Match(postCode string) bool {
	for _, match := range m.Matches {
		if postCode == match {
			return true
		}
	}

	i, err := strconv.Atoi(postCode)

	if err != nil {
		return false
	}

	for _, match := range m.Ranges {
		if i >= match.Start && i <= match.End {
			return true
		}
	}

	return false
}

// PostCodeRange defines a range of numeric post codes, inclusive of the Start and End.
type PostCodeRange struct {
	Start int
	End   int
}

// RegexMatcher defines a post code matcher that uses a regular expression.
type RegexMatcher struct {
	Regex *regexp.Regexp
}

// Match returns whether the post code is matched by the regular expression defined in the matcher.
func (m RegexMatcher) Match(postCode string) bool {
	return m.Regex.MatchString(postCode)
}
