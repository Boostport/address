package address

type country struct {
	ID   string
	Name string

	DefaultLanguage string

	PostCodePrefix string
	PostCodeRegex  postCodeRegex

	Format          string
	LatinizedFormat string

	AdministrativeAreaNameType FieldName
	LocalityNameType           FieldName
	DependentLocalityNameType  FieldName
	PostCodeNameType           FieldName

	AllowedFields  map[Field]struct{}
	RequiredFields map[Field]struct{}
	Upper          map[Field]struct{}

	AdministrativeAreas map[string][]administrativeArea
}

type postCodeRegex struct {
	regex            string
	subdivisionRegex map[string]postCodeRegex
}

type administrativeArea struct {
	ID        string
	Name      string
	PostalKey string

	Localities []locality
}

type locality struct {
	ID   string
	Name string

	DependentLocalities []dependentLocality
}

type dependentLocality struct {
	ID   string
	Name string
}

type data map[string]country

func (d data) getCountry(countryCode string) country {
	data := generated[countryCode]
	defaults := generated["ZZ"]

	if data.Format == "" {
		data.Format = defaults.Format
	}

	if data.AdministrativeAreaNameType == 0 {
		data.AdministrativeAreaNameType = defaults.AdministrativeAreaNameType
	}

	if data.LocalityNameType == 0 {
		data.LocalityNameType = defaults.LocalityNameType
	}

	if data.DependentLocalityNameType == 0 {
		data.DependentLocalityNameType = defaults.DependentLocalityNameType
	}

	if data.PostCodeNameType == 0 {
		data.PostCodeNameType = defaults.PostCodeNameType
	}

	if len(data.AllowedFields) <= 0 {
		data.AllowedFields = defaults.AllowedFields
	}

	if len(data.RequiredFields) <= 0 {
		data.RequiredFields = defaults.RequiredFields
	}

	if len(data.Upper) <= 0 {
		data.Upper = defaults.Upper
	}

	return data
}

func (d data) hasCountry(countryCode string) bool {
	_, ok := generated[countryCode]

	return ok
}

func (d data) getAdministrativeAreaName(countryCode, administrativeAreaID, language string) string {

	data := d.getCountry(countryCode)

	lang := d.normalizeLanguage(countryCode, language)

	for _, adminArea := range data.AdministrativeAreas[lang] {
		if adminArea.ID == administrativeAreaID {
			return adminArea.Name
		}
	}

	return ""
}

func (d data) getAdministrativeAreaPostalKey(countryCode, administrativeAreaID string) string {

	data := d.getCountry(countryCode)

	lang := d.normalizeLanguage(countryCode, "")

	for _, adminArea := range data.AdministrativeAreas[lang] {
		if adminArea.ID == administrativeAreaID {
			return adminArea.PostalKey
		}
	}

	return ""
}

func (d data) getLocalityName(countryCode, administrativeAreaID, localityID, language string) string {

	data := d.getCountry(countryCode)

	lang := d.normalizeLanguage(countryCode, language)

	for _, adminArea := range data.AdministrativeAreas[lang] {
		if adminArea.ID == administrativeAreaID {
			for _, locality := range adminArea.Localities {
				if locality.ID == localityID {
					return locality.Name
				}
			}
		}
	}

	return ""
}

func (d data) getDependentLocalityName(countryCode, administrativeAreaID, localityID, dependentLocalityID, language string) string {

	data := d.getCountry(countryCode)

	lang := d.normalizeLanguage(countryCode, language)

	for _, adminArea := range data.AdministrativeAreas[lang] {
		if adminArea.ID == administrativeAreaID {
			for _, locality := range adminArea.Localities {
				if locality.ID == localityID {
					for _, dependentLocality := range locality.DependentLocalities {
						if dependentLocality.ID == dependentLocalityID {
							return dependentLocality.Name
						}
					}
				}
			}
		}
	}

	return ""
}

func (d data) normalizeLanguage(countryCode, language string) string {

	country := d.getCountry(countryCode)

	if _, ok := country.AdministrativeAreas[language]; ok {
		return language
	}

	return country.DefaultLanguage
}
