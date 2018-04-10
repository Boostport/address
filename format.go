package address

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	textLanguage "golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var isAlphabetRegex = regexp.MustCompile(`[A-Za-z]`)

var collapseWhitespaceRegex = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

var collapseBRRegex = regexp.MustCompile(`(?:\s|\p{Zs}|<br>){2,}`)

var funcMap = template.FuncMap{
	"toUpper": strings.ToUpper,
	"len": func(slice []string) int {
		return len(slice)
	},
	"inc": func(i int) int {
		return i + 1
	},
	"join": strings.Join,
}

// DefaultFormatter formats an address using the country's address format and includes the name of the country.
// If Latinize is set to true, in countries where a latinized address format is provided, the latinized format is used.
type DefaultFormatter struct {
	Output   Outputter
	Latinize bool
}

// Format formats an address. The language must be a valid ISO 639-1 language code. It is used to convert the keys
// in administrative areas, localities and dependent localities into their actual names. If the provided language
// does not have any translations, it falls back to the default language used by the country.
func (d DefaultFormatter) Format(address Address, language string) string {

	language = generated.normalizeLanguage(address.Country, language)

	format, isLatinized := getFormat(address.Country, d.Latinize)

	if isLatinized {
		format += "%n%country"
	} else {
		format = "%country%n" + format
	}

	t := d.Output.TransformFormat(format, map[Field]struct{}{})

	compiled := template.Must(template.New("").Funcs(funcMap).Parse(t))

	buf := bytes.NewBuffer([]byte{})

	compiled.Execute(buf, address.toFormatData(generated.getCountry(address.Country), language))

	return collapseBRRegex.ReplaceAllString(collapseWhitespaceRegex.ReplaceAllString(strings.TrimSpace(buf.String()), "\n"), "<br>")
}

// PostalLabelFormatter formats an address for postal labels. It uppercases address fields as required by the country's
// addressing standards. If the address it in the same country as the origin country, the country is omitted.
// The country name is added to the address, both in the language of the origin country as well as English, following
// recommendations of the Universal Postal Union, to avoid difficulties in transit.
// The OriginCountryCode field should be set to the ISO 3166-1 country code of the originating country.
// If Latinize is set to true, in countries where a latinized address format is provided, the latinized format is used.
type PostalLabelFormatter struct {
	Output            Outputter
	OriginCountryCode string
	Latinize          bool
}

// Format formats an address. The language must be a valid ISO 639-1 language code. It is used to convert the keys
// in administrative areas, localities and dependent localities into their actual names. If the provided language
// does not have any translations, it falls back to the default language used by the country.
func (f PostalLabelFormatter) Format(address Address, language string) string {

	language = generated.normalizeLanguage(address.Country, language)

	format, isLatinized := getFormat(address.Country, f.Latinize)

	countryData := generated.getCountry(address.Country)

	addressData := address.toFormatData(countryData, language)

	// Include the country since this is an international mail
	if generated.hasCountry(f.OriginCountryCode) && strings.ToUpper(f.OriginCountryCode) != strings.ToUpper(address.Country) {

		originLanguage, _ := textLanguage.Make(fmt.Sprintf("und-%s", f.OriginCountryCode)).Base()

		destinationCountry := textLanguage.MustParseRegion(address.Country)

		oNamer := display.Regions(textLanguage.MustParse(originLanguage.String()))
		namer := display.Regions(textLanguage.English)

		englishDestination := namer.Name(destinationCountry)

		var translatedDestination string

		if oNamer != nil {
			translatedDestination = oNamer.Name(destinationCountry)
		} else {
			translatedDestination = englishDestination
		}

		if translatedDestination != englishDestination {
			addressData.Country = fmt.Sprintf("%s - %s", strings.ToUpper(translatedDestination), strings.ToUpper(englishDestination))
		} else {
			addressData.Country = strings.ToUpper(englishDestination)
		}

		if addressData.Country != "" {
			if isLatinized {
				format += "%n%country"
			} else {
				format = "%country%n" + format
			}
		}
	}

	if isAlphabetRegex.MatchString(addressData.AdministrativeAreaPostalKey) {
		addressData.AdministrativeArea = addressData.AdministrativeAreaPostalKey
	}

	t := f.Output.TransformFormat(format, countryData.Upper)

	compiled := template.Must(template.New("").Funcs(funcMap).Parse(t))

	buf := bytes.NewBuffer([]byte{})

	compiled.Execute(buf, addressData)

	return collapseBRRegex.ReplaceAllString(collapseWhitespaceRegex.ReplaceAllString(strings.TrimSpace(buf.String()), "\n"), "<br>")
}

// Outputter defines an interface to transform an address format in Google's format into a Go template that is merged
// with the address data to produce a formatted address.
type Outputter interface {
	TransformFormat(format string, upper map[Field]struct{}) string
}

// HTMLOutputter outputs the formatted address as an HTML fragments. The address fields are annotated using the class
// attribute and `<br>`s are used for new lines.
type HTMLOutputter struct{}

// TransformFormat transforms an address format in Google's format into a HTML template. The upper map is used to
// determine which fields should be converted to UPPERCASE.
func (h HTMLOutputter) TransformFormat(format string, upper map[Field]struct{}) string {

	r := strings.NewReplacer(
		"%country", fmt.Sprintf(`{{if ne .%s "" }}<span class="country">{{.%s}}</span>{{end}}`, Country, toUpper(Country, "", upper)),
		"%N", fmt.Sprintf(`{{if ne .%s "" }}<span class="name">{{.%s}}</span>{{end}}`, Name, toUpper(Name, "", upper)),
		"%O", fmt.Sprintf(`{{if ne .%s "" }}<span class="organization">{{.%s}}</span>{{end}}`, Organization, toUpper(Organization, "", upper)),
		"%A", fmt.Sprintf(`{{$numLines:=.%s|len}}{{range $lineNo, $line := .%s}}{{$realLineNo := inc $lineNo}}<span class="address-line-{{$realLineNo}}">{{%s}}</span>{{if ne $numLines $realLineNo}}<br>{{end}}{{end}}`, StreetAddress, StreetAddress, toUpper(StreetAddress, "$line", upper)),
		"%D", fmt.Sprintf(`{{if ne .%s "" }}<span class="dependent-locality">{{.%s}}</span>{{end}}`, DependentLocality, toUpper(DependentLocality, "", upper)),
		"%C", fmt.Sprintf(`{{if ne .%s "" }}<span class="locality">{{.%s}}</span>{{end}}`, Locality, toUpper(Locality, "", upper)),
		"%S", fmt.Sprintf(`{{if ne .%s "" }}<span class="administrative-area">{{.%s}}</span>{{end}}`, AdministrativeArea, toUpper(AdministrativeArea, "", upper)),
		"%Z", fmt.Sprintf(`{{if ne .%s "" }}<span class="post-code">{{.%s}}</span>{{end}}`, PostCode, toUpper(PostCode, "", upper)),
		"%X", fmt.Sprintf(`{{if ne .%s "" }}<span class="sorting-code">{{.%s}}</span>{{end}}`, SortingCode, toUpper(SortingCode, "", upper)),
		"%n", "<br>",
	)

	return r.Replace(format)
}

// StringOutputter outputs the formatted address as a string and `\n`s are used for new lines.
type StringOutputter struct{}

// TransformFormat transforms an address format in Google's format into a string template. The upper map is used to
// determine which fields should be converted to UPPERCASE.
func (s StringOutputter) TransformFormat(format string, upper map[Field]struct{}) string {

	r := strings.NewReplacer(
		"%country", fmt.Sprintf("{{.%s}}", toUpper(Country, "", upper)),
		"%N", fmt.Sprintf("{{.%s}}", toUpper(Name, "", upper)),
		"%O", fmt.Sprintf("{{.%s}}", toUpper(Organization, "", upper)),
		"%A", fmt.Sprintf(`{{ join .%s }}`, toUpper(StreetAddress, fmt.Sprintf(`%s "\n"`, StreetAddress), upper)),
		"%D", fmt.Sprintf("{{.%s}}", toUpper(DependentLocality, "", upper)),
		"%C", fmt.Sprintf("{{.%s}}", toUpper(Locality, "", upper)),
		"%S", fmt.Sprintf("{{.%s}}", toUpper(AdministrativeArea, "", upper)),
		"%Z", fmt.Sprintf("{{.%s}}", toUpper(PostCode, "", upper)),
		"%X", fmt.Sprintf("{{.%s}}", toUpper(SortingCode, "", upper)),
		"%n", "\n",
	)

	return r.Replace(format)
}

func toUpper(field Field, fieldName string, upperCaseFields map[Field]struct{}) string {

	fieldToUse := field.String()

	if fieldName != "" {
		fieldToUse = fieldName
	}

	for upperField := range upperCaseFields {
		if field == upperField {
			return fmt.Sprintf("%s | toUpper", fieldToUse)
		}
	}

	return fmt.Sprintf("%s", fieldToUse)
}

func getFormat(countryCode string, latinized bool) (string, bool) {

	countryData := generated.getCountry(countryCode)

	if latinized && countryData.LatinizedFormat != "" {
		return countryData.LatinizedFormat, true
	}

	if countryData.Format != "" && countryData.LatinizedFormat == "" {
		return countryData.Format, true
	}

	if countryData.Format != "" {
		return countryData.Format, false
	}

	return generated.getCountry("ZZ").Format, true
}
