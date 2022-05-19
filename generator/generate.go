package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Boostport/address"
	"golang.org/x/text/language"
)

const rootURL = "https://chromium-i18n.appspot.com/ssl-address"

const numWorkers = 25

var addressFormatRegex = regexp.MustCompile(`%[NOADCSZX]`)

var urlRemoveLanguageRegex = regexp.MustCompile(`--.*`)

var postPrefixFixes = map[string]string{
	"PR": "PR ",
}

var defaultLanguageOverrides = map[string]string{
	"AQ": "en",
	"AS": "en",
	"BQ": "nl",
	"BV": "nb",
	"CW": "nl",
	"DJ": "fr",
	"GS": "en",
	"HM": "en",
	"MV": "en",
	"PG": "en",
	"PW": "en",
	"TK": "en",
	"VU": "fr",
	"WS": "en",
}

var localNameOverrides = map[string]string{
	"TV": "Tuvalu",
}

type countriesJSON struct {
	Countries string `json:"countries"`
}

type countryJSON struct {
	ID  string `json:"id"`
	Key string `json:"key"`

	Lang      string `json:"lang"`
	Languages string `json:"languages"`
	Name      string `json:"name"`

	Fmt  string `json:"fmt"`
	Lfmt string `json:"lfmt"`

	StateNameType       string `json:"state_name_type"`
	LocalityNameType    string `json:"locality_name_type"`
	SubLocalityNameType string `json:"sublocality_name_type"`
	ZipNameType         string `json:"zip_name_type"`

	Require string `json:"require"`
	Upper   string `json:"upper"`

	SubISOIDs string `json:"sub_isoids"`
	SubKeys   string `json:"sub_keys"`
	SubLNames string `json:"sub_lnames"`
	SubNames  string `json:"sub_names"`

	SubMores string `json:"sub_mores"`

	SubXRequires string `json:"sub_xrequires"`
	SubXZips     string `json:"sub_xzips"`

	SubZips   string `json:"sub_zips"`
	SubZipExs string `json:"sub_zipexs"`

	PostPrefix string `json:"post_prefix"`
	Zip        string `json:"zip"`
	Zipex      string `json:"zipex"`
}

type subdivisionJSON struct {
	ID  string `json:"id"`
	Key string `json:"key"`

	Name  string `json:"name"`
	LName string `json:"lname"`

	Lang string `json:"lang"`

	ISOID   string `json:"isoid"`
	SubKeys string `json:"sub_keys"`

	SubNames   string `json:"sub_names"`
	SubMores   string `json:"sub_mores"`
	SubLNames  string `json:"sub_lnames"`
	SubLFNames string `json:"sub_lfnames"`

	Zip       string `json:"zip"`
	ZipEx     string `json:"zipex"`
	SubZips   string `json:"sub_zips"`
	SubZipExs string `json:"sub_zipexs"`
}

type postCodeRegex struct {
	regex            string
	subdivisionRegex map[string]postCodeRegex
}

func (p postCodeRegex) toCode() string {

	// Generate postcode regex in order to avoid huge diffs when updating the data
	var ids []string

	for id := range p.subdivisionRegex {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	str := fmt.Sprintf(`{
		regex: `+"`%s`,", p.regex)

	if len(p.subdivisionRegex) > 0 {
		str += `
subdivisionRegex: map[string]postCodeRegex{
`

		for _, id := range ids {
			str += fmt.Sprintf(`"%s": %s,
`, id, p.subdivisionRegex[id].toCode())
		}

		str += `}`
	}

	str += `}`

	return str
}

type country struct {
	ID   string
	Name string

	DefaultLanguage string

	PostCodePrefix string
	PostCodeRegex  postCodeRegex

	Format          string
	LatinizedFormat string

	AdministrativeAreaNameType address.FieldName
	LocalityNameType           address.FieldName
	DependentLocalityNameType  address.FieldName
	PostCodeNameType           address.FieldName

	AllowedFields  map[address.Field]struct{}
	RequiredFields map[address.Field]struct{}
	Upper          map[address.Field]struct{}

	AdministrativeAreas map[string][]administrativeArea
}

func (c country) toCode() string {

	str := fmt.Sprintf(`{
	ID: "%s",
	Name: "%s",`, c.ID, c.Name)

	if c.DefaultLanguage != "" {
		str += fmt.Sprintf(`
DefaultLanguage: "%s",`, c.DefaultLanguage)
	} else {
		fmt.Println(c.ID)
	}

	if c.PostCodePrefix != "" {
		str += fmt.Sprintf(`
PostCodePrefix: "%s",`, c.PostCodePrefix)
	}

	if c.PostCodeRegex.regex != "" || len(c.PostCodeRegex.subdivisionRegex) > 0 {
		str += fmt.Sprintf(`
PostCodeRegex: postCodeRegex%s,`, c.PostCodeRegex.toCode())
	}

	if c.Format != "" {
		str += fmt.Sprintf(`
Format: "%s",`, c.Format)
	}

	if c.LatinizedFormat != "" {
		str += fmt.Sprintf(`
LatinizedFormat: "%s",`, c.LatinizedFormat)
	}

	if c.AdministrativeAreaNameType != 0 {
		str += fmt.Sprintf(`
AdministrativeAreaNameType: %s,`, c.AdministrativeAreaNameType.String())
	}

	if c.LocalityNameType != 0 {
		str += fmt.Sprintf(`
LocalityNameType: %s,`, c.LocalityNameType.String())
	}

	if c.DependentLocalityNameType != 0 {
		str += fmt.Sprintf(`
DependentLocalityNameType: %s,`, c.DependentLocalityNameType.String())
	}

	if c.PostCodeNameType != 0 {
		str += fmt.Sprintf(`
PostCodeNameType: %s,`, c.PostCodeNameType.String())
	}

	if len(c.AllowedFields) > 0 {

		// Generate fields in order to avoid huge diffs when updating the data
		var fields []string

		for field := range c.AllowedFields {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
AllowedFields: map[Field]struct{}{`)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `
},`
	}

	if len(c.RequiredFields) > 0 {

		// Generate fields in order to avoid huge diffs when updating the data
		var fields []string

		for field := range c.RequiredFields {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
RequiredFields: map[Field]struct{}{`)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `
},`
	}

	if len(c.Upper) > 0 {

		// Generate fields in order to avoid huge diffs when updating the data
		var fields []string

		for field := range c.Upper {
			fields = append(fields, field.String())
		}

		sort.Strings(fields)

		str += fmt.Sprintf(`
Upper: map[Field]struct{}{`)

		for _, field := range fields {
			str += fmt.Sprintf(`
%s: {},`, field)
		}

		str += `
},`
	}

	if len(c.AdministrativeAreas) > 0 {

		// Generate languages in order to avoid huge diffs when updating the address data
		var languages []string

		for language := range c.AdministrativeAreas {
			languages = append(languages, language)
		}

		sort.Strings(languages)

		str += fmt.Sprintf(`
AdministrativeAreas: map[string][]administrativeArea {`)

		for _, language := range languages {

			areas := c.AdministrativeAreas[language]

			str += fmt.Sprintf(`
"%s": {`, language)

			for _, area := range areas {
				str += fmt.Sprintf(`
%s,`, area.toCode())
			}

			str += `
},`
		}

		str += `
},`
	}

	str += `
}`

	return str
}

type administrativeArea struct {
	ID        string
	Name      string
	PostalKey string

	Localities []locality
}

func (a administrativeArea) toCode() string {

	str := fmt.Sprintf(`{
	ID: "%s",
	Name: "%s",
	PostalKey: "%s",`, a.ID, a.Name, a.PostalKey)

	if len(a.Localities) > 0 {

		str += `
Localities: []locality{
`
		for _, l := range a.Localities {
			str += l.toCode() + ",\n"
		}

		str += `
},`
	}

	str += `
}`

	return str

}

type locality struct {
	ID   string
	Name string

	DependentLocalities []dependentLocality
}

func (l locality) toCode() string {

	str := fmt.Sprintf(`{
	ID: "%s",
	Name: "%s",`, l.ID, l.Name)

	if len(l.DependentLocalities) > 0 {

		str += `
DependentLocalities: []dependentLocality{
`
		for _, dl := range l.DependentLocalities {
			str += dl.toCode() + ",\n"
		}

		str += `
},`
	}

	str += `
}`

	return str
}

type dependentLocality struct {
	ID   string
	Name string
}

func (d dependentLocality) toCode() string {

	return fmt.Sprintf(`{
	ID: "%s",
	Name: "%s",
}`, d.ID, d.Name)
}

func main() {

	fmt.Printf("Downloading address data from %s. This may take a few minutes.\n", rootURL)

	start := time.Now()

	countriesResp, err := http.Get(rootURL + "/data")

	if err != nil {
		log.Fatalf("Error getting countries from endpoint: %s", err)
	}

	if countriesResp.StatusCode != 200 {
		log.Fatalf("Error getting countries from endpoint, error code %d", countriesResp.StatusCode)
	}

	countriesUnmarshaled := &countriesJSON{}

	countriesDecoder := json.NewDecoder(countriesResp.Body)

	err = countriesDecoder.Decode(countriesUnmarshaled)

	if err != nil {
		log.Fatalf("Error unmarshaling countries JSON: %s", err)
	}

	countries := strings.Split(countriesUnmarshaled.Countries, "~")
	countries = append(countries, "ZZ") // Include the fall back ZZ (unknown) country

	countryCodeCh := make(chan string, len(countries))
	stopCh := make(chan struct{})
	resultCh := make(chan workerResult)

	for i := 0; i < numWorkers; i++ {

		w := &worker{
			countryCodes: countryCodeCh,
			stop:         stopCh,
			result:       resultCh,
		}

		w.start()
	}

	for _, country := range countries {
		countryCodeCh <- country
	}

	processedCountries := map[string]country{}

	fmt.Print("Processed: ")

	for i := 0; i < len(countries); i++ {

		result, ok := <-resultCh

		if !ok {
			break
		}

		if result.Error != nil {
			close(stopCh)
			log.Fatalf("Error processing country: %s", result.Error)
		}

		fmt.Printf("%s ", result.Country.ID)
		processedCountries[result.Country.ID] = result.Country
	}

	// Order the countries by ID, so that the order of the generated countries will be deterministic.
	// This prevents huge diffs when updating the data.
	var countriesInOrder []string

	for countryID := range processedCountries {
		countriesInOrder = append(countriesInOrder, countryID)
	}

	sort.Strings(countriesInOrder)

	fmt.Println("\nGenerating code...")
	generated := `// Code generated by address. DO NOT EDIT.
package address

var generated = data{
`

	for _, country := range countriesInOrder {

		generated += fmt.Sprintf(`
"%s":%s,`, country, processedCountries[country].toCode())
	}

	generated += `
}`

	fmt.Println("Formatting generated code...")
	formatted, err := format.Source([]byte(generated))

	if err != nil {
		log.Fatalf("Error formatting generated source: %s", err)
	}

	err = ioutil.WriteFile("data.generated.go", formatted, os.ModePerm)

	if err != nil {
		log.Fatalf("Error writing data.go: %s", err)
	}

	timeTaken := time.Since(start)

	fmt.Printf("Total time taken: %s\n", timeTaken)
}

type workerResult struct {
	Error   error
	Country country
}

type worker struct {
	countryCodes chan string
	stop         chan struct{}
	result       chan workerResult
}

func (w *worker) start() {
	go func() {
		for {
		exit:
			select {
			case countryCode := <-w.countryCodes:

				url := rootURL + "/data/" + countryCode

				countryData, err := http.Get(url)

				if err != nil {
					w.result <- workerResult{
						Error: fmt.Errorf("error getting data using url (%s): %s", url, err),
					}

					break
				}

				countryJSON, err := decodeCountryJSON(countryData.Body, url)

				if err != nil {
					w.result <- workerResult{
						Error: fmt.Errorf("error unmarhaling JSON for url (%s): %s", url, err),
					}

					break
				}

				// Sanity check latinized format
				if countryJSON.Lfmt != "" && len(getAllowedFields(countryJSON.Fmt)) != len(getAllowedFields(countryJSON.Lfmt)) {
					w.result <- workerResult{
						Error: fmt.Errorf("number of fields in the address format and latinized address format does not match for %s", countryJSON.Key),
					}

					break
				}

				// Sanity check post code regex
				if countryJSON.Zip != "" {
					err = checkPostalCodeRegex("^("+countryJSON.Zip+")$", strings.Split(countryJSON.Zipex, ","))

					if err != nil {
						w.result <- workerResult{
							Error: fmt.Errorf("error validating post code regex for %s: %s", countryJSON.Key, err),
						}

						break
					}
				}

				country := country{
					ID:   countryCode,
					Name: countryJSON.Name,

					Format:          countryJSON.Fmt,
					LatinizedFormat: countryJSON.Lfmt,

					AllowedFields:  getAllowedFields(countryJSON.Fmt),
					RequiredFields: getFields(countryJSON.Require),
					Upper:          getFields(countryJSON.Upper),
				}

				if countryJSON.Zip != "" {
					country.PostCodeRegex.regex = "^(" + countryJSON.Zip + ")$"
				}

				if countryJSON.Lang != "" {
					country.DefaultLanguage = countryJSON.Lang

				} else if lang, ok := defaultLanguageOverrides[countryCode]; ok {
					country.DefaultLanguage = lang

				} else {
					lang, _ := language.Make(fmt.Sprintf("und-%s", countryCode)).Base()
					country.DefaultLanguage = lang.String()
				}

				if countryJSON.StateNameType != "" {
					administrativeAreaNameType, err := convertFieldNameToConstant(countryJSON.StateNameType)

					if err != nil {
						w.result <- workerResult{
							Error: fmt.Errorf("error converting administrative area name type for %s: %s", countryJSON.Key, err),
						}

						break
					}

					country.AdministrativeAreaNameType = administrativeAreaNameType
				}

				if countryJSON.LocalityNameType != "" {
					localityNameType, err := convertFieldNameToConstant(countryJSON.LocalityNameType)

					if err != nil {
						w.result <- workerResult{
							Error: fmt.Errorf("error converting locality name type for %s: %s", countryJSON.Key, err),
						}

						break
					}

					country.LocalityNameType = localityNameType
				}

				if countryJSON.SubLocalityNameType != "" {
					dependentLocalityNameType, err := convertFieldNameToConstant(countryJSON.SubLocalityNameType)

					if err != nil {
						w.result <- workerResult{
							Error: fmt.Errorf("error converting dependent locality name type for %s: %s", countryJSON.Key, err),
						}

						break
					}

					country.DependentLocalityNameType = dependentLocalityNameType
				}

				if countryJSON.ZipNameType != "" {
					postCodeNameType, err := convertFieldNameToConstant(countryJSON.ZipNameType)

					if err != nil {
						w.result <- workerResult{
							Error: fmt.Errorf("error converting post code name type for %s: %s", countryJSON.Key, err),
						}

						break
					}

					country.PostCodeNameType = postCodeNameType
				}

				if prefix, ok := postPrefixFixes[countryJSON.Key]; ok {
					country.PostCodePrefix = prefix
				} else {
					country.PostCodePrefix = countryJSON.PostPrefix
				}

				// Process subdivisions
				if countryJSON.SubKeys != "" {

					// Sanity check
					if countryJSON.Languages == "" {
						w.result <- workerResult{
							Error: fmt.Errorf("%s has subkeys but does not have any languages", countryJSON.Key),
						}

						break
					}

					country.AdministrativeAreas = map[string][]administrativeArea{}

					// Get languages
					languages := strings.Split(countryJSON.Languages, "~")

					if len(languages) > 1 {
						for _, language := range languages {

							if language != countryJSON.Lang {

								countryData, err := http.Get(url + "--" + language)

								if err != nil {
									w.result <- workerResult{
										Error: fmt.Errorf("error getting language %s for country %s: %s", language, countryJSON.Key, err),
									}

									break exit
								}

								languageCountryJSON, err := decodeCountryJSON(countryData.Body, url)

								if err != nil {
									w.result <- workerResult{
										Error: fmt.Errorf("error decoding language %s for country %s: %s", language, countryJSON.Key, err),
									}

									break exit
								}

								languageAdminAreas, _, err := processAdministrativeAreas(languageCountryJSON, language)

								if err != nil {
									w.result <- workerResult{
										Error: fmt.Errorf("error processing admin areas in language %s for country %s: %s", language, countryJSON.Key, err),
									}

									break exit
								}

								for lang, adminAreas := range languageAdminAreas {
									country.AdministrativeAreas[lang] = adminAreas
								}
							} else {
								adminAreas, postCodeRegex, err := processAdministrativeAreas(countryJSON, "")

								if err != nil {
									w.result <- workerResult{
										Error: fmt.Errorf("error processing admin areas in the default language for country %s: %s", countryJSON.Key, err),
									}

									break exit
								}

								country.PostCodeRegex.subdivisionRegex = postCodeRegex

								for lang, adminAreas := range adminAreas {
									country.AdministrativeAreas[lang] = adminAreas
								}
							}
						}
					} else {
						adminAreas, postCodeRegex, err := processAdministrativeAreas(countryJSON, "")

						if err != nil {
							w.result <- workerResult{
								Error: fmt.Errorf("error processing admin areas in the default language for country %s: %s", countryJSON.Key, err),
							}

							break exit
						}

						country.PostCodeRegex.subdivisionRegex = postCodeRegex

						for lang, adminAreas := range adminAreas {
							country.AdministrativeAreas[lang] = adminAreas
						}
					}
				}

				w.result <- workerResult{
					Country: country,
				}

			case <-w.stop:
				return
			}
		}
	}()
}

func processAdministrativeAreas(countryJSON countryJSON, language string) (map[string][]administrativeArea, map[string]postCodeRegex, error) {

	result := map[string][]administrativeArea{}
	postCodeResult := map[string]postCodeRegex{}

	subISOIDs := strings.Split(countryJSON.SubISOIDs, "~")
	subNames := strings.Split(countryJSON.SubNames, "~")
	subZips := strings.Split(countryJSON.SubZips, "~")
	subMores := strings.Split(countryJSON.SubMores, "~")
	subKeys := strings.Split(countryJSON.SubKeys, "~")
	subZipExs := strings.Split(countryJSON.SubZipExs, "~")
	subLNames := strings.Split(countryJSON.SubLNames, "~")

	// Countries like China include places like Taiwan and Hong Kong in their list of administrative divisions.
	// However, these places are already in the list of countries, so we check to see if they have special post
	// code regex or required fields to filter them out
	subdivisionsToSkip := map[string]struct{}{}

	if countryJSON.SubXRequires != "" {
		for idx, requires := range strings.Split(countryJSON.SubXRequires, "~") {
			if requires != "" {
				subdivisionsToSkip[subISOIDs[idx]] = struct{}{}
			}
		}
	}

	if countryJSON.SubXZips != "" {
		for idx, xzip := range strings.Split(countryJSON.SubXZips, "~") {
			if xzip != "" {
				subdivisionsToSkip[subISOIDs[idx]] = struct{}{}
			}
		}
	}

	var processedAdministrativeAreas []administrativeArea
	var latinizedAdministrativeAreas []administrativeArea

	var ids []string

	// Deal with the case where a country has sub keys, but the list of ISO ids is blank (ex: ES)
	// Also, prefer sub-keys and treat them as authoritative when valid addresses can include
	// administrative areas that don't have an ISO code, e.g. United States addresses can include military addresses
	// with AA, AE, AP.
	useSubKeys := countryJSON.SubISOIDs == "" || countryJSON.Key == "US"
	if useSubKeys {
		ids = subKeys
	} else if countryJSON.SubKeys != "" {
		ids = subISOIDs
	}

	for i, isoID := range ids {

		if isoID == "" {
			if useSubKeys {
				isoID = subKeys[i]
			} else {
				// Skip administrative ares without iso ids due to regions being contested or
				// not recognized. (ex: Crimea and Sevastopol in Russia)
				continue
			}
		}

		if _, ok := subdivisionsToSkip[isoID]; ok {
			continue
		}

		// Sanity check
		if countryJSON.SubZips != "" && countryJSON.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			err := checkPostalCodeRegex("^"+subZips[i], strings.Split(subZipExs[i], ","))

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error checking administrative area post code regex for %s / %s against sample: %s", isoID, countryJSON.Key, err)
			}
		}

		adminArea := administrativeArea{
			ID:        isoID,
			PostalKey: subKeys[i],
		}

		if countryJSON.SubNames != "" {
			adminArea.Name = subNames[i]
		} else {
			adminArea.Name = subKeys[i]
		}

		if countryJSON.SubZips != "" && subZips[i] != "" {
			postCodeResult[isoID] = postCodeRegex{
				regex: "^" + subZips[i],
			}
		}

		var latinizedAdminArea administrativeArea

		if countryJSON.SubLNames != "" {
			latinizedAdminArea.ID = isoID
			latinizedAdminArea.Name = subLNames[i]
			latinizedAdminArea.PostalKey = subKeys[i]
		}

		if countryJSON.SubMores != "" && subMores[i] == "true" {

			url := rootURL + "/" + urlRemoveLanguageRegex.ReplaceAllString(countryJSON.ID, "") + "/" + subKeys[i]

			if language != "" {
				url += "--" + language
			}

			administrativeAreaData, err := http.Get(url)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error getting administrative area data for %s: %s", url, err)
			}

			administrativeAreaJSON, err := decodeSubdivisionJSON(administrativeAreaData.Body, url)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error decoding administrative area JSON: %s", err)
			}

			localities, subPostCodeReg, err := processLocalities(administrativeAreaJSON, language)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error processing localities for %s/%s: %s", countryJSON.Key, subKeys[i], err)
			}

			// Sanity check
			_, hasAdminPostCodeRegex := postCodeResult[isoID]

			if !hasAdminPostCodeRegex && len(subPostCodeReg) > 0 {
				return result, postCodeResult, fmt.Errorf("locality %s has postcode regexes, but the parent locality does not", administrativeAreaJSON.ID)
			}

			if len(subPostCodeReg) > 0 {
				postCodeReg := postCodeResult[isoID]
				postCodeReg.subdivisionRegex = subPostCodeReg
				postCodeResult[isoID] = postCodeReg
			}

			if len(localities[countryJSON.Lang]) > 0 {
				adminArea.Localities = localities[countryJSON.Lang]
			}

			// Consider latinized names to be english
			if administrativeAreaJSON.SubLNames != "" {

				// Sanity check
				if _, ok := localities["en"]; !ok {
					return result, postCodeResult, fmt.Errorf("%s has latinized admin areas, but does not have any latinized localities for %s", countryJSON.Key, administrativeAreaJSON.ID)
				}

				latinizedAdminArea.Localities = localities["en"]
			}
		}

		processedAdministrativeAreas = append(processedAdministrativeAreas, adminArea)

		if latinizedAdminArea.ID != "" {
			latinizedAdministrativeAreas = append(latinizedAdministrativeAreas, latinizedAdminArea)
		}
	}

	result[countryJSON.Lang] = processedAdministrativeAreas

	if len(latinizedAdministrativeAreas) > 0 {

		// Sanity check
		if len(latinizedAdministrativeAreas) != len(processedAdministrativeAreas) {
			return result, postCodeResult, fmt.Errorf("number of latinized admin areas (%d) does not match number of admin areas (%d) for %s", len(latinizedAdministrativeAreas), len(processedAdministrativeAreas), countryJSON.ID)
		}

		sort.Slice(latinizedAdministrativeAreas, func(i, j int) bool {
			return latinizedAdministrativeAreas[i].Name < latinizedAdministrativeAreas[j].Name
		})

		result["en"] = latinizedAdministrativeAreas
	}

	return result, postCodeResult, nil
}

func processLocalities(administrativeAreaJSON subdivisionJSON, language string) (map[string][]locality, map[string]postCodeRegex, error) {

	result := map[string][]locality{}
	postCodeResult := map[string]postCodeRegex{}

	subKeys := strings.Split(administrativeAreaJSON.SubKeys, "~")
	subNames := strings.Split(administrativeAreaJSON.SubNames, "~")
	subMores := strings.Split(administrativeAreaJSON.SubMores, "~")
	subZips := strings.Split(administrativeAreaJSON.SubZips, "~")
	subZipExs := strings.Split(administrativeAreaJSON.SubZipExs, "~")
	subLNames := strings.Split(administrativeAreaJSON.SubLNames, "~")

	var latinizedLocalities []locality

	var processedLocalities []locality

	for i, key := range subKeys {

		// Sanity check
		if administrativeAreaJSON.SubZips != "" && administrativeAreaJSON.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			err := checkPostalCodeRegex("^"+subZips[i], strings.Split(subZipExs[i], ","))

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error checking default locality post code regex for %s against sample: %s", administrativeAreaJSON.ID, err)
			}
		}

		defaultLocality := locality{
			ID: key, // No ISO ID at this level, so we use the key from Google's data set
		}

		if administrativeAreaJSON.SubNames != "" {
			defaultLocality.Name = subNames[i]
		} else {
			defaultLocality.Name = subKeys[i]
		}

		if administrativeAreaJSON.SubZips != "" && subZips[i] != "" {
			postCodeResult[key] = postCodeRegex{
				regex: "^" + subZips[i],
			}
		}

		var latinizedLocality locality

		if administrativeAreaJSON.SubLNames != "" {
			latinizedLocality.ID = key
			latinizedLocality.Name = subLNames[i]
		}

		if administrativeAreaJSON.SubMores != "" && subMores[i] == "true" {

			url := rootURL + "/" + urlRemoveLanguageRegex.ReplaceAllString(administrativeAreaJSON.ID, "") + "/" + subKeys[i]

			if language != "" {
				url += "--" + language
			}

			localityData, err := http.Get(url)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error getting defaultLocality area data for %s: %s", url, err)
			}

			localityAreaJSON, err := decodeSubdivisionJSON(localityData.Body, url)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error decoding defaultLocality area JSON: %s", err)
			}

			dependentLocalities, subPostCodeReg, err := processDependentLocalities(localityAreaJSON)

			if err != nil {
				return result, postCodeResult, fmt.Errorf("error processing dependent localities for %s/%s: %s", administrativeAreaJSON.ID, subKeys[i], err)
			}

			// Sanity check
			_, hasLocalityPostCodeRegex := postCodeResult[key]

			if !hasLocalityPostCodeRegex && len(subPostCodeReg) > 0 {
				return result, postCodeResult, fmt.Errorf("dependent locality %s/%s has postcode regexes, but the parent locality does not", administrativeAreaJSON.ID, subKeys[i])
			}

			if len(subPostCodeReg) > 0 {
				postCodeReg := postCodeResult[key]
				postCodeReg.subdivisionRegex = subPostCodeReg
				postCodeResult[key] = postCodeReg
			}

			if len(dependentLocalities[administrativeAreaJSON.Lang]) > 0 {
				defaultLocality.DependentLocalities = dependentLocalities[administrativeAreaJSON.Lang]
			}

			// Consider latinized names to be english
			if administrativeAreaJSON.SubLNames != "" {

				// Sanity check
				if _, ok := dependentLocalities["en"]; !ok {
					return result, postCodeResult, fmt.Errorf("%s has latinized localities, but does not have any latinized dependent localities for %s", administrativeAreaJSON.ID, subKeys[i])
				}

				latinizedLocality.DependentLocalities = dependentLocalities["en"]
			}
		}

		processedLocalities = append(processedLocalities, defaultLocality)

		if administrativeAreaJSON.SubLNames != "" {
			latinizedLocalities = append(latinizedLocalities, latinizedLocality)
		}
	}

	result[administrativeAreaJSON.Lang] = processedLocalities

	if len(latinizedLocalities) > 0 {

		// Sanity check
		if len(latinizedLocalities) != len(processedLocalities) {
			return result, postCodeResult, fmt.Errorf("number of latinized localities (%d) does not match number of localities (%d) for %s", len(latinizedLocalities), len(processedLocalities), administrativeAreaJSON.ID)
		}

		sort.Slice(latinizedLocalities, func(i, j int) bool {
			return latinizedLocalities[i].Name < latinizedLocalities[j].Name
		})

		result["en"] = latinizedLocalities
	}

	return result, postCodeResult, nil
}

func processDependentLocalities(localityJSON subdivisionJSON) (map[string][]dependentLocality, map[string]postCodeRegex, error) {

	result := map[string][]dependentLocality{}
	postCodeReg := map[string]postCodeRegex{}

	subKeys := strings.Split(localityJSON.SubKeys, "~")
	subNames := strings.Split(localityJSON.SubNames, "~")
	subZips := strings.Split(localityJSON.SubZips, "~")
	subZipExs := strings.Split(localityJSON.SubZipExs, "~")

	var processedDependentLocalities []dependentLocality

	for i, key := range subKeys {

		// Sanity check
		if localityJSON.SubZips != "" && localityJSON.SubZipExs != "" && subZips[i] != "" && subZipExs[i] != "" {
			err := checkPostalCodeRegex("^"+subZips[i], strings.Split(subZipExs[i], ","))

			if err != nil {
				return result, postCodeReg, fmt.Errorf("error checking dependent locality post code regex for %s against sample: %s", localityJSON.ID, err)
			}
		}

		dependentLocality := dependentLocality{
			ID: key, // No ISO ID at this level, so we use the key from Google's data set
		}

		if localityJSON.SubNames != "" {
			dependentLocality.Name = subNames[i]
		} else {
			dependentLocality.Name = subKeys[i]
		}

		if localityJSON.SubZips != "" && subZips[i] != "" {
			postCodeReg[key] = postCodeRegex{
				regex: "^" + subZips[i],
			}
		}

		processedDependentLocalities = append(processedDependentLocalities, dependentLocality)
	}

	result[localityJSON.Lang] = processedDependentLocalities

	// We consider latinized names to be english
	if localityJSON.SubLNames != "" {

		subLNames := strings.Split(localityJSON.SubLNames, "~")

		var latinizedDependentLocalities []dependentLocality

		for i, key := range subKeys {

			dependentLocality := dependentLocality{
				ID:   key, // No ISO ID at this level, so we use the key from Google's data set
				Name: subLNames[i],
			}

			latinizedDependentLocalities = append(latinizedDependentLocalities, dependentLocality)
		}

		sort.Slice(latinizedDependentLocalities, func(i, j int) bool {
			return latinizedDependentLocalities[i].Name < latinizedDependentLocalities[j].Name
		})

		result["en"] = latinizedDependentLocalities
	}

	return result, postCodeReg, nil
}

func decodeCountryJSON(reader io.Reader, url string) (countryJSON, error) {

	countryJSON := countryJSON{}

	countryDecoder := json.NewDecoder(reader)

	err := countryDecoder.Decode(&countryJSON)

	if err != nil {
		return countryJSON, fmt.Errorf("error unmarhaling JSON for url (%s): %s", url, err)
	}

	return countryJSON, nil
}

func decodeSubdivisionJSON(reader io.Reader, url string) (subdivisionJSON, error) {

	subdivisionJSON := subdivisionJSON{}

	subdivisionDecoder := json.NewDecoder(reader)

	err := subdivisionDecoder.Decode(&subdivisionJSON)

	if err != nil {
		return subdivisionJSON, fmt.Errorf("error unmarhaling JSON for url (%s): %s", url, err)
	}

	return subdivisionJSON, nil
}

func getFields(fields string) map[address.Field]struct{} {

	upper := map[address.Field]struct{}{}

	for _, field := range fields {

		switch string(field) {
		case "N":
			upper[address.Name] = struct{}{}
		case "O":
			upper[address.Organization] = struct{}{}
		case "A":
			upper[address.StreetAddress] = struct{}{}
		case "D":
			upper[address.DependentLocality] = struct{}{}
		case "C":
			upper[address.Locality] = struct{}{}
		case "S":
			upper[address.AdministrativeArea] = struct{}{}
		case "Z":
			upper[address.PostCode] = struct{}{}
		case "X":
			upper[address.SortingCode] = struct{}{}
		}
	}

	return upper
}

func getAllowedFields(format string) map[address.Field]struct{} {

	allowed := map[address.Field]struct{}{}

	fields := addressFormatRegex.FindAllString(format, -1)

	for _, field := range fields {

		switch field {

		case "%N":
			allowed[address.Name] = struct{}{}
		case "%O":
			allowed[address.Organization] = struct{}{}
		case "%A":
			allowed[address.StreetAddress] = struct{}{}
		case "%D":
			allowed[address.DependentLocality] = struct{}{}
		case "%C":
			allowed[address.Locality] = struct{}{}
		case "%S":
			allowed[address.AdministrativeArea] = struct{}{}
		case "%Z":
			allowed[address.PostCode] = struct{}{}
		case "%X":
			allowed[address.SortingCode] = struct{}{}
		}
	}

	return allowed
}

func convertFieldNameToConstant(fieldName string) (address.FieldName, error) {

	switch fieldName {
	case "area":
		return address.Area, nil
	case "city":
		return address.City, nil
	case "county":
		return address.County, nil
	case "department":
		return address.Department, nil
	case "district":
		return address.District, nil
	case "do_si":
		return address.DoSi, nil
	case "eircode":
		return address.Eircode, nil
	case "emirate":
		return address.Emirate, nil
	case "island":
		return address.Island, nil
	case "neighborhood":
		return address.Neighborhood, nil
	case "oblast":
		return address.Oblast, nil
	case "pin":
		return address.PINCode, nil
	case "parish":
		return address.Parish, nil
	case "post_town":
		return address.PostTown, nil
	case "postal":
		return address.PostalCode, nil
	case "prefecture":
		return address.Prefecture, nil
	case "province":
		return address.Province, nil
	case "state":
		return address.State, nil
	case "suburb":
		return address.Suburb, nil
	case "townland":
		return address.Townland, nil
	case "village_township":
		return address.VillageTownship, nil
	case "zip":
		return address.ZipCode, nil
	}

	return address.FieldName(-1), fmt.Errorf("unknown field name: %s", fieldName)
}

func checkPostalCodeRegex(regex string, postalCodes []string) error {

	postCodeRegex, err := regexp.Compile(regex)

	if err != nil {
		return errors.New("unable to compile zip regex")
	}

	for _, postCode := range postalCodes {
		if !postCodeRegex.MatchString(postCode) {
			return fmt.Errorf("sample postcode %s could not be validated by post code regex", postCode)
		}
	}

	return nil
}
