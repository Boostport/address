package address

import (
	"errors"
	"regexp"
	"strings"
)

// Validate checks and address to determine if it is valid.
// If you want to create valid addresses, the `address.NewValid()` function does it in one call.
func Validate(address Address) error {

	var errs []error

	if !generated.hasCountry(address.Country) {
		errs = append(errs, ErrInvalidCountryCode)
		return errors.Join(errs...)
	}

	countryData := generated.getCountry(address.Country)

	if err := checkRequiredFields(address, countryData.RequiredFields); err != nil {
		errs = append(errs, err)
	}

	if err := checkAllowedFields(address, countryData.AllowedFields); err != nil {
		errs = append(errs, err)
	}

	if len(countryData.AdministrativeAreas) > 0 {

		if administrativeAreaData, ok := countryData.AdministrativeAreas[countryData.DefaultLanguage]; ok {
			err := checkSubdivisions(address, administrativeAreaData)

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if address.PostCode != "" {
		err := checkPostCode(address, countryData.PostCodeRegex)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func checkRequiredFields(address Address, requiredFields map[Field]struct{}) error {

	errors := ErrMissingRequiredFields{
		country: address.Country,
	}

	for field := range requiredFields {

		switch field {
		case Name:
			if strings.TrimSpace(address.Name) == "" {
				errors.Fields = append(errors.Fields, Name)
			}

		case Organization:
			if strings.TrimSpace(address.Organization) == "" {
				errors.Fields = append(errors.Fields, Organization)
			}

		case StreetAddress:

			isEmpty := true

			for _, addressLine := range address.StreetAddress {
				if strings.TrimSpace(addressLine) != "" {
					isEmpty = false
					break
				}
			}

			if isEmpty {
				errors.Fields = append(errors.Fields, StreetAddress)
			}

		case DependentLocality:
			if strings.TrimSpace(address.DependentLocality) == "" {
				errors.Fields = append(errors.Fields, DependentLocality)
			}

		case Locality:
			if strings.TrimSpace(address.Locality) == "" {
				errors.Fields = append(errors.Fields, Locality)
			}

		case AdministrativeArea:
			if strings.TrimSpace(address.AdministrativeArea) == "" {
				errors.Fields = append(errors.Fields, AdministrativeArea)
			}

		case PostCode:
			if strings.TrimSpace(address.PostCode) == "" {
				errors.Fields = append(errors.Fields, PostCode)
			}

		case SortingCode:
			if strings.TrimSpace(address.SortingCode) == "" {
				errors.Fields = append(errors.Fields, SortingCode)
			}
		}
	}

	if len(errors.Fields) <= 0 {
		return nil
	}

	return errors
}

func checkAllowedFields(address Address, allowedFields map[Field]struct{}) error {

	errors := ErrUnsupportedFields{
		country: address.Country,
	}

	if _, ok := allowedFields[Name]; address.Name != "" && !ok {
		errors.Fields = append(errors.Fields, Name)
	}

	if _, ok := allowedFields[Organization]; address.Organization != "" && !ok {
		errors.Fields = append(errors.Fields, Organization)
	}

	if _, ok := allowedFields[StreetAddress]; len(address.StreetAddress) > 0 && !ok {
		errors.Fields = append(errors.Fields, StreetAddress)
	}

	if _, ok := allowedFields[DependentLocality]; address.DependentLocality != "" && !ok {
		errors.Fields = append(errors.Fields, DependentLocality)
	}

	if _, ok := allowedFields[Locality]; address.Locality != "" && !ok {
		errors.Fields = append(errors.Fields, Locality)
	}

	if _, ok := allowedFields[AdministrativeArea]; address.AdministrativeArea != "" && !ok {
		errors.Fields = append(errors.Fields, AdministrativeArea)
	}

	if _, ok := allowedFields[PostCode]; address.PostCode != "" && !ok {
		errors.Fields = append(errors.Fields, PostCode)
	}

	if _, ok := allowedFields[SortingCode]; address.SortingCode != "" && !ok {
		errors.Fields = append(errors.Fields, SortingCode)
	}

	if len(errors.Fields) <= 0 {
		return nil
	}

	return errors
}

func checkSubdivisions(address Address, administrativeAreaData []administrativeArea) error {

	var errs []error

	if address.AdministrativeArea != "" {

		adminAreaIdx := -1

		for i, adminArea := range administrativeAreaData {
			if adminArea.ID == address.AdministrativeArea {
				adminAreaIdx = i
			}
		}

		if adminAreaIdx == -1 {
			errs = append(errs, ErrInvalidAdministrativeArea)
			return errors.Join(errs...)
		}

		localities := administrativeAreaData[adminAreaIdx].Localities

		localityIdx := -1

		if address.Locality == "" || len(localities) <= 0 {
			return errors.Join(errs...)
		}

		for i, locality := range localities {
			if locality.ID == address.Locality {
				localityIdx = i
			}
		}

		if localityIdx == -1 {
			errs = append(errs, ErrInvalidLocality)
			return errors.Join(errs...)
		}

		dependentLocalities := localities[localityIdx].DependentLocalities

		dependentLocalitiesIdx := -1

		if address.DependentLocality == "" || len(dependentLocalities) <= 0 {
			return errors.Join(errs...)
		}

		for i, dl := range dependentLocalities {
			if dl.ID == address.DependentLocality {
				dependentLocalitiesIdx = i
			}
		}

		if dependentLocalitiesIdx == -1 {
			errs = append(errs, ErrInvalidDependentLocality)
			return errors.Join(errs...)
		}
	}

	return errors.Join(errs...)
}

func checkPostCode(address Address, regex postCodeRegex) error {

	var errs []error

	if address.PostCode != "" && regex.regex != "" {

		country := regex

		countryRegex := regexp.MustCompile(country.regex)

		if !countryRegex.MatchString(address.PostCode) {
			errs = append(errs, ErrInvalidPostCode)
			return errors.Join(errs...)
		}

		if adminArea, ok := country.subdivisionRegex[address.AdministrativeArea]; ok {

			adminAreaRegex := regexp.MustCompile(adminArea.regex)

			if !adminAreaRegex.MatchString(address.PostCode) {
				errs = append(errs, ErrInvalidPostCode)
				return errors.Join(errs...)
			}

			if locality, ok := adminArea.subdivisionRegex[address.Locality]; ok {

				localityRegex := regexp.MustCompile(locality.regex)

				if !localityRegex.MatchString(address.PostCode) {
					errs = append(errs, ErrInvalidPostCode)
					return errors.Join(errs...)
				}

				if dependentLocality, ok := locality.subdivisionRegex[address.DependentLocality]; ok {

					dependentLocalityRegex := regexp.MustCompile(dependentLocality.regex)

					if !dependentLocalityRegex.MatchString(address.PostCode) {
						errs = append(errs, ErrInvalidPostCode)
						return errors.Join(errs...)
					}
				}
			}
		}
	}

	return errors.Join(errs...)
}
