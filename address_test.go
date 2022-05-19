package address

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestAddressIsZero(t *testing.T) {

	tests := []struct {
		Address Address
		IsZero  bool
	}{
		{
			Address: Address{},
			IsZero:  true,
		},
		{
			Address: Address{
				Name: "John Smith",
			},
			IsZero: false,
		},
		{
			Address: Address{
				Organization: "Company Pty Ltd",
			},
			IsZero: false,
		},
		{
			Address: Address{
				StreetAddress: []string{
					"525 Collins Street",
				},
			},
			IsZero: false,
		},
		{
			Address: Address{
				DependentLocality: "test",
			},
			IsZero: false,
		},
		{
			Address: Address{
				Locality: "Melbourne",
			},
			IsZero: false,
		},
		{
			Address: Address{
				AdministrativeArea: "VIC",
			},
			IsZero: false,
		},
		{
			Address: Address{
				PostCode: "3000",
			},
			IsZero: false,
		},
		{
			Address: Address{
				Country: "AU",
			},
			IsZero: false,
		},
		{
			Address: Address{
				SortingCode: "1234",
			},
			IsZero: false,
		},
		{
			Address: Address{
				Name:         "John Smith",
				Organization: "Company Pty Ltd",
				StreetAddress: []string{
					"525 Collins Street",
				},
				Locality:           "Melbourne",
				AdministrativeArea: "VIC",
				PostCode:           "3000",
				Country:            "AU",
			},
			IsZero: false,
		},
	}

	for i, testCase := range tests {

		isZero := testCase.Address.IsZero()

		if isZero != testCase.IsZero {
			t.Errorf("Result for IsZero() in test case %d does not match expected result", i)
		}
	}
}

func TestValidAddresses(t *testing.T) {

	tests := []struct {
		Address  []func(*Address)
		Expected Address
	}{
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
			Expected: Address{
				Country:      "AU",
				Name:         "John Smith",
				Organization: "Company Pty Ltd",
				StreetAddress: []string{
					"525 Collins Street",
				},
				Locality:           "Melbourne",
				AdministrativeArea: "VIC",
				PostCode:           "3000",
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
			Expected: Address{
				Country: "AU",
				StreetAddress: []string{
					"525 Collins Street",
				},
				Locality:           "Melbourne",
				AdministrativeArea: "VIC",
				PostCode:           "3000",
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"27 Rue Pasteur",
				}),
				WithLocality("Cabourg"),
				WithPostCode("14390"),
				WithCountry("FR"),
			},
			Expected: Address{
				Country: "FR",
				StreetAddress: []string{
					"27 Rue Pasteur",
				},
				Locality: "Cabourg",
				PostCode: "14390",
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"27 Rue Pasteur",
				}),
				WithLocality("Cabourg"),
				WithPostCode("14390"),
				WithCountry("FR"),
			},
			Expected: Address{
				Country: "FR",
				StreetAddress: []string{
					"27 Rue Pasteur",
				},
				Locality: "Cabourg",
				PostCode: "14390",
			},
		},
		{
			Address: []func(*Address){
				WithName("PFC John Smith"),
				WithStreetAddress([]string{
					"PSC 1234, Box 12345",
				}),
				WithLocality("APO"),
				WithAdministrativeArea("AE"),
				WithPostCode("09204-1234"),
				WithCountry("US"),
			},
			Expected: Address{
				Country: "US",
				Name:    "PFC John Smith",
				StreetAddress: []string{
					"PSC 1234, Box 12345",
				},
				Locality:           "APO",
				AdministrativeArea: "AE",
				PostCode:           "09204-1234",
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"No.1 Jianguomenwai Avenue",
				}),
				WithDependentLocality("临翔区"),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: Address{
				Country: "CN",
				StreetAddress: []string{
					"No.1 Jianguomenwai Avenue",
				},
				DependentLocality:  "临翔区",
				Locality:           "临沧市",
				AdministrativeArea: "53",
				PostCode:           "677000",
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"Jangnyang-ro 17beon-gil",
				}),
				WithDependentLocality("북구"),
				WithLocality("포항시"),
				WithAdministrativeArea("47"),
				WithPostCode("37592"), // Invalid post code
				WithCountry("KR"),
			},
			Expected: Address{
				Country: "KR",
				StreetAddress: []string{
					"Jangnyang-ro 17beon-gil",
				},
				DependentLocality:  "북구",
				Locality:           "포항시",
				AdministrativeArea: "47",
				PostCode:           "37592",
			},
		},
	}

	for i, testCase := range tests {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		if !reflect.DeepEqual(address, testCase.Expected) {
			t.Errorf("Address in test case %d does not match expected address", i)
		}
	}
}

func TestInvalidAddresses(t *testing.T) {

	tests := []struct {
		Address []func(*Address)
	}{
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("2000"), // Invalid post code
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("VIC 3000"), // Invalid post code
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000 VIC"), // Invalid post code
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("VIC 3000 VIC"), // Invalid post code
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithLocality("Melbourne"), // Missing street address
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"525 Collins Street",
				}),
				WithDependentLocality("Toorak"), // Extraneous field
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"No.1 Jianguomenwai Avenue",
				}),
				WithAdministrativeArea("ASDF"), // Invalid administrative area
				WithPostCode("677000"),
				WithCountry("CN"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"No.1 Jianguomenwai Avenue",
				}),
				WithAdministrativeArea("53"), // Missing locality and dependent locality
				WithPostCode("677000"),
				WithCountry("CN"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"No.1 Jianguomenwai Avenue",
				}),
				WithLocality("ASDF"), // Invalid locality
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"No.1 Jianguomenwai Avenue",
				}),
				WithDependentLocality("ASDF"), // Invalid dependent locality
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
		},
		{
			Address: []func(*Address){
				WithStreetAddress([]string{
					"Jangnyang-ro 17beon-gil",
				}),
				WithDependentLocality("북구"),
				WithLocality("포항시"),
				WithAdministrativeArea("47"),
				WithPostCode("38100"), // Invalid post code
				WithCountry("KR"),
			},
		},
	}

	for i, testCase := range tests {
		_, err := NewValid(testCase.Address...)

		if err == nil {
			t.Errorf("Expected an error when creating an address using test case %d, but there was no error", i)
		}
	}
}

func TestGetCountries(t *testing.T) {

	countries := ListCountries("en")

	numCountries := len(countries)
	expected := len(generated) - 1

	if numCountries != expected {
		t.Errorf("Number of countries (%d) does not equal expected number of countries (%d)", numCountries, expected)
	}

	// Check language
	for _, country := range countries {
		if country.Code == "AU" {
			if country.Name != "Australia" {
				t.Errorf("Expected English name of Australia to be Australia, got %s", country.Name)
				break
			}
		}
	}

	countries = ListCountries("zh")

	numCountries = len(countries)
	expected = len(generated) - 1

	if numCountries != expected {
		t.Errorf("Number of countries (%d) does not equal expected number of countries (%d)", numCountries, expected)
	}

	// Check language
	for _, country := range countries {
		if country.Code == "AU" {
			if country.Name != "澳大利亚" {
				t.Errorf("Expected Chinese name of Australia to be 澳大利亚, got %s", country.Name)
				break
			}
		}
	}
}

func TestGetCountry(t *testing.T) {

	country := GetCountry("AU")

	expected := CountryData{
		Format: "%O%n%N%n%A%n%C %S %Z",
		Required: []Field{
			AdministrativeArea, Locality, PostCode, StreetAddress,
		},
		Allowed: []Field{
			AdministrativeArea, Locality, Name, Organization, PostCode, StreetAddress,
		},
		DefaultLanguage:            "en",
		AdministrativeAreaNameType: State,
		LocalityNameType:           Suburb,
		DependentLocalityNameType:  Suburb,
		PostCodeNameType:           PostalCode,
		PostCodeRegex: PostCodeRegexData{
			Regex: `^(\d{4})$`,
			SubdivisionRegex: map[string]PostCodeRegexData{
				"ACT": {
					Regex: `^29|2540|260|261[0-8]|02|2620`},
				"NSW": {
					Regex: `^1|2[0-57-8]|26[2-9]|261[189]|3500|358[56]|3644|3707`},
				"NT": {
					Regex: `^0[89]`},
				"QLD": {
					Regex: `^[49]`},
				"SA": {
					Regex: `^5|0872`},
				"TAS": {
					Regex: `^7`},
				"VIC": {
					Regex: `^[38]`},
				"WA": {
					Regex: `^6|0872`},
			},
		},
		AdministrativeAreas: map[string][]AdministrativeAreaData{
			"en": {
				{
					ID:   "ACT",
					Name: "Australian Capital Territory",
				},
				{
					ID:   "NSW",
					Name: "New South Wales",
				},
				{
					ID:   "NT",
					Name: "Northern Territory",
				},
				{
					ID:   "QLD",
					Name: "Queensland",
				},
				{
					ID:   "SA",
					Name: "South Australia",
				},
				{
					ID:   "TAS",
					Name: "Tasmania",
				},
				{
					ID:   "VIC",
					Name: "Victoria",
				},
				{
					ID:   "WA",
					Name: "Western Australia",
				},
			},
		},
	}

	if !reflect.DeepEqual(country, expected) {
		t.Errorf("Country data for AU does not match expected country data")
	}

	country = GetCountry("AC")

	expected = CountryData{
		Format: "%N%n%O%n%A%n%C%n%Z",
		Required: []Field{
			Locality, StreetAddress,
		},
		Allowed: []Field{
			Locality, Name, Organization, PostCode, StreetAddress,
		},
		DefaultLanguage:            "en",
		AdministrativeAreaNameType: Province,
		LocalityNameType:           City,
		DependentLocalityNameType:  Suburb,
		PostCodeNameType:           PostalCode,
		PostCodeRegex: PostCodeRegexData{
			Regex: `^(ASCN 1ZZ)$`,
		},
	}

	if !reflect.DeepEqual(country, expected) {
		t.Errorf("Country data for AC does not match expected country data")
	}

	country = GetCountry("KR")

	expected = CountryData{
		Format:          "%S %C%D%n%A%n%O%n%N%n%Z",
		LatinizedFormat: "%N%n%O%n%A%n%D%n%C%n%S%n%Z",
		Required: []Field{
			AdministrativeArea, Locality, PostCode, StreetAddress,
		},
		Allowed: []Field{
			AdministrativeArea, DependentLocality, Locality, Name, Organization, PostCode, StreetAddress,
		},
		DefaultLanguage:            "ko",
		AdministrativeAreaNameType: DoSi,
		LocalityNameType:           City,
		DependentLocalityNameType:  District,
		PostCodeNameType:           PostalCode,
		PostCodeRegex: PostCodeRegexData{
			Regex: `^(\d{5})$`,
			SubdivisionRegex: map[string]PostCodeRegexData{
				"11": {
					Regex: `^0[1-8]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"강남구": {
							Regex: `^06[0-3]`},
						"강동구": {
							Regex: `^05[2-4]`},
						"강북구": {
							Regex: `^01[0-2]`},
						"강서구": {
							Regex: `^07[5-8]`},
						"관악구": {
							Regex: `^08[78]`},
						"광진구": {
							Regex: `^0(?:49|5[01])`},
						"구로구": {
							Regex: `^08[23]`},
						"금천구": {
							Regex: `^08[56]`},
						"노원구": {
							Regex: `^01[6-9]`},
						"도봉구": {
							Regex: `^01[34]`},
						"동대문구": {
							Regex: `^02[4-6]`},
						"동작구": {
							Regex: `^0(?:69|70)`},
						"마포구": {
							Regex: `^0(?:39|4[0-2])`},
						"서대문구": {
							Regex: `^03[67]`},
						"서초구": {
							Regex: `^06[5-8]`},
						"성동구": {
							Regex: `^04[78]`},
						"성북구": {
							Regex: `^02[78]`},
						"송파구": {
							Regex: `^05[5-8]`},
						"양천구": {
							Regex: `^0(?:7[89]|8[01])`},
						"영등포구": {
							Regex: `^07[2-4]`},
						"용산구": {
							Regex: `^04[34]`},
						"은평구": {
							Regex: `^03[3-5]`},
						"종로구": {
							Regex: `^03[01]`},
						"중구": {
							Regex: `^04[56]|100`},
						"중랑구": {
							Regex: `^02[0-2]`},
					}},
				"26": {
					Regex: `^4[6-9]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"강서구": {
							Regex: `^467`},
						"금정구": {
							Regex: `^46[23]`},
						"기장군": {
							Regex: `^460`},
						"남구": {
							Regex: `^48[45]`},
						"동구": {
							Regex: `^48[78]`},
						"동래구": {
							Regex: `^47[789]`},
						"부산진구": {
							Regex: `^47[123]`},
						"북구": {
							Regex: `^46[56]`},
						"사상구": {
							Regex: `^4(?:69|70)`},
						"사하구": {
							Regex: `^49[345]`},
						"서구": {
							Regex: `^492`},
						"수영구": {
							Regex: `^48[23]`},
						"연제구": {
							Regex: `^47[56]`},
						"영도구": {
							Regex: `^49[01]`},
						"중구": {
							Regex: `^489`},
						"해운대구": {
							Regex: `^48[01]`},
					}},
				"27": {
					Regex: `^4[123]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"남구": {
							Regex: `^42[45]`},
						"달서구": {
							Regex: `^42[678]`},
						"달성군": {
							Regex: `^4(?:29|30)`},
						"동구": {
							Regex: `^41[0-2]`},
						"북구": {
							Regex: `^41[45]`},
						"서구": {
							Regex: `^41[78]`},
						"수성구": {
							Regex: `^42[0-2]`},
						"중구": {
							Regex: `^419`},
					}},
				"28": {
					Regex: `^2[1-3]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"강화군": {
							Regex: `^230`},
						"계양구": {
							Regex: `^21[01]`},
						"남구": {
							Regex: `^22[12]`},
						"남동구": {
							Regex: `^21[5-7]`},
						"동구": {
							Regex: `^225`},
						"부평구": {
							Regex: `^21[34]`},
						"서구": {
							Regex: `^22[6-8]`},
						"연수구": {
							Regex: `^2(?:19|20)`},
						"옹진군": {
							Regex: `^231`},
						"중구": {
							Regex: `^223`},
					}},
				"29": {
					Regex: `^6[12]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"광산구": {
							Regex: `^62[2-4]`},
						"남구": {
							Regex: `^61[67]`},
						"동구": {
							Regex: `^61[45]`},
						"북구": {
							Regex: `^61[0-2]`},
						"서구": {
							Regex: `^6(?:19|20)`},
					}},
				"30": {
					Regex: `^3[45]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"대덕구": {
							Regex: `^34[34]`},
						"동구": {
							Regex: `^34[5-7]`},
						"서구": {
							Regex: `^35[2-4]`},
						"유성구": {
							Regex: `^34[0-2]`},
						"중구": {
							Regex: `^3(?:4[89]|50)`},
					}},
				"31": {
					Regex: `^4[45]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"남구": {
							Regex: `^44[67]`},
						"동구": {
							Regex: `^44[01]`},
						"북구": {
							Regex: `^442`},
						"울주군": {
							Regex: `^4(?:49|50)`},
						"중구": {
							Regex: `^44[45]`},
					}},
				"41": {
					Regex: `^1[0-8]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"가평군": {
							Regex: `^124`},
						"고양시": {
							Regex: `^10[2-5]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"덕양구": {
									Regex: `^10[245]`},
								"일산동구": {
									Regex: `^10[2-4]`},
								"일산서구": {
									Regex: `^10[2-4]`},
							}},
						"과천시": {
							Regex: `^138`},
						"광명시": {
							Regex: `^14[23]`},
						"광주시": {
							Regex: `^12[78]`},
						"구리시": {
							Regex: `^119`},
						"군포시": {
							Regex: `^158`},
						"김포시": {
							Regex: `^10[01]`},
						"남양주시": {
							Regex: `^12[0-3]`},
						"동두천시": {
							Regex: `^113`},
						"부천시": {
							Regex: `^14[4-7]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"소사구": {
									Regex: `^14[67]`},
								"오정구": {
									Regex: `^14[45]`},
								"원미구": {
									Regex: `^14[456]`},
							}},
						"성남시": {
							Regex: `^13[1-6]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"분당구": {
									Regex: `^13[3-6]`},
								"수정구": {
									Regex: `^13[1-46]`},
								"중원구": {
									Regex: `^13[1-4]`},
							}},
						"수원시": {
							Regex: `^16[2-7]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"권선구": {
									Regex: `^16[3-6]`},
								"영통구": {
									Regex: `^16[245-7]`},
								"장안구": {
									Regex: `^16[2-4]`},
								"팔달구": {
									Regex: `^16[2-6]`},
							}},
						"시흥시": {
							Regex: `^1(?:49|5[01])`},
						"안산시": {
							Regex: `^15[2-6]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"단원구": {
									Regex: `^15[2-6]`},
								"상록구": {
									Regex: `^15[2-6]`},
							}},
						"안성시": {
							Regex: `^17[56]`},
						"안양시": {
							Regex: `^1(?:39|4[01])`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"동안구": {
									Regex: `^1(?:39|4[01])`},
								"만안구": {
									Regex: `^1(?:39|40)`},
							}},
						"양주시": {
							Regex: `^11[45]`},
						"양평군": {
							Regex: `^125`},
						"여주시": {
							Regex: `^126`},
						"연천군": {
							Regex: `^110`},
						"오산시": {
							Regex: `^181`},
						"용인시": {
							Regex: `^1(?:6[89]|7[01])`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"기흥구": {
									Regex: `^1(?:6[89]|7[01])`},
								"수지구": {
									Regex: `^16[89]`},
								"처인구": {
									Regex: `^1(?:6[89]|7[01])`},
							}},
						"의왕시": {
							Regex: `^16[01]`},
						"의정부시": {
							Regex: `^11[6-8]`},
						"이천시": {
							Regex: `^17[34]`},
						"파주시": {
							Regex: `^10[89]`},
						"평택시": {
							Regex: `^1(?:7[7-9]|80)`},
						"포천시": {
							Regex: `^111`},
						"하남시": {
							Regex: `^1(?:29|30)`},
						"화성시": {
							Regex: `^18[2-6]`},
					}},
				"42": {
					Regex: `^2[456]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"강릉시": {
							Regex: `^25[4-6]`},
						"고성군": {
							Regex: `^247`},
						"동해시": {
							Regex: `^25[78]`},
						"삼척시": {
							Regex: `^259`},
						"속초시": {
							Regex: `^24[89]`},
						"양구군": {
							Regex: `^245`},
						"양양군": {
							Regex: `^250`},
						"영월군": {
							Regex: `^262`},
						"원주시": {
							Regex: `^26[3-5]`},
						"인제군": {
							Regex: `^246`},
						"정선군": {
							Regex: `^261`},
						"철원군": {
							Regex: `^240`},
						"춘천시": {
							Regex: `^24[2-4]`},
						"태백시": {
							Regex: `^260`},
						"평창군": {
							Regex: `^253`},
						"홍천군": {
							Regex: `^251`},
						"화천군": {
							Regex: `^241`},
						"횡성군": {
							Regex: `^252`},
					}},
				"43": {
					Regex: `^2[789]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"괴산군": {
							Regex: `^280`},
						"단양군": {
							Regex: `^270`},
						"보은군": {
							Regex: `^289`},
						"영동군": {
							Regex: `^291`},
						"옥천군": {
							Regex: `^290`},
						"음성군": {
							Regex: `^27[67]`},
						"제천시": {
							Regex: `^27[12]`},
						"증평군": {
							Regex: `^279`},
						"진천군": {
							Regex: `^278`},
						"청주시": {
							Regex: `^28[0-9]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"상당구": {
									Regex: `^28[1-3578]`},
								"서원구": {
									Regex: `^28[1-35-8]`},
								"청원구": {
									Regex: `^28[13-5]`},
								"흥덕구": {
									Regex: `^28[13-6]`},
							}},
						"충주시": {
							Regex: `^27[3-5]`},
					}},
				"44": {
					Regex: `^3[1-3]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"계룡시": {
							Regex: `^328`},
						"공주시": {
							Regex: `^32[56]`},
						"금산군": {
							Regex: `^327`},
						"논산시": {
							Regex: `^3(?:29|30)`},
						"당진시": {
							Regex: `^31[78]`},
						"보령시": {
							Regex: `^33[45]`},
						"부여군": {
							Regex: `^33[12]`},
						"서산시": {
							Regex: `^3(?:19|20)`},
						"서천군": {
							Regex: `^336`},
						"아산시": {
							Regex: `^31[45]`},
						"예산군": {
							Regex: `^324`},
						"천안시": {
							Regex: `^31[0-2]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"동남구": {
									Regex: `^31[0-2]`},
								"서북구": {
									Regex: `^31[01]`},
							}},
						"청양군": {
							Regex: `^333`},
						"태안군": {
							Regex: `^321`},
						"홍성군": {
							Regex: `^322`},
					}},
				"45": {
					Regex: `^5[4-6]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"고창군": {
							Regex: `^564`},
						"군산시": {
							Regex: `^54[01]`},
						"김제시": {
							Regex: `^54[34]`},
						"남원시": {
							Regex: `^55[78]`},
						"무주군": {
							Regex: `^555`},
						"부안군": {
							Regex: `^563`},
						"순창군": {
							Regex: `^560`},
						"완주군": {
							Regex: `^553`},
						"익산시": {
							Regex: `^54[56]`},
						"임실군": {
							Regex: `^559`},
						"장수군": {
							Regex: `^556`},
						"전주시": {
							Regex: `^5(?:4[89]|5[01])`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"덕진구": {
									Regex: `^5(?:4[89]|50)`},
								"완산구": {
									Regex: `^5(?:4[89]|5[01])`},
							}},
						"정읍시": {
							Regex: `^56[12]`},
						"진안군": {
							Regex: `^554`},
					}},
				"46": {
					Regex: `^5[7-9]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"강진군": {
							Regex: `^592`},
						"고흥군": {
							Regex: `^595`},
						"곡성군": {
							Regex: `^575`},
						"광양시": {
							Regex: `^57[78]`},
						"구례군": {
							Regex: `^576`},
						"나주시": {
							Regex: `^58[23]`},
						"담양군": {
							Regex: `^573`},
						"목포시": {
							Regex: `^58[67]`},
						"무안군": {
							Regex: `^585`},
						"보성군": {
							Regex: `^594`},
						"순천시": {
							Regex: `^5(?:79|80)`},
						"신안군": {
							Regex: `^588`},
						"여수시": {
							Regex: `^59[67]`},
						"영광군": {
							Regex: `^570`},
						"영암군": {
							Regex: `^584`},
						"완도군": {
							Regex: `^591`},
						"장성군": {
							Regex: `^572`},
						"장흥군": {
							Regex: `^593`},
						"진도군": {
							Regex: `^589`},
						"함평군": {
							Regex: `^571`},
						"해남군": {
							Regex: `^590`},
						"화순군": {
							Regex: `^581`},
					}},
				"47": {
					Regex: `^(?:3[6-9]|40)\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"경산시": {
							Regex: `^38[4-6]`},
						"경주시": {
							Regex: `^38[0-2]`},
						"고령군": {
							Regex: `^401`},
						"구미시": {
							Regex: `^39[1-4]`},
						"군위군": {
							Regex: `^390`},
						"김천시": {
							Regex: `^39[56]`},
						"문경시": {
							Regex: `^3(?:69|70)`},
						"봉화군": {
							Regex: `^362`},
						"상주시": {
							Regex: `^37[12]`},
						"성주군": {
							Regex: `^400`},
						"안동시": {
							Regex: `^36[67]`},
						"영덕군": {
							Regex: `^364`},
						"영양군": {
							Regex: `^365`},
						"영주시": {
							Regex: `^36[01]`},
						"영천시": {
							Regex: `^38[89]`},
						"예천군": {
							Regex: `^368`},
						"울릉군": {
							Regex: `^402`},
						"울진군": {
							Regex: `^363`},
						"의성군": {
							Regex: `^373`},
						"청도군": {
							Regex: `^383`},
						"청송군": {
							Regex: `^374`},
						"칠곡군": {
							Regex: `^39[89]`},
						"포항시": {
							Regex: `^37[5-9]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"남구": {
									Regex: `^37[6-9]`},
								"북구": {
									Regex: `^37[5-79]`},
							}},
					}},
				"48": {
					Regex: `^5[0-3]\d{2}`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"거제시": {
							Regex: `^53[23]`},
						"거창군": {
							Regex: `^501`},
						"고성군": {
							Regex: `^529`},
						"김해시": {
							Regex: `^5(?:0[89]|10)`},
						"남해군": {
							Regex: `^524`},
						"밀양시": {
							Regex: `^504`},
						"사천시": {
							Regex: `^525`},
						"산청군": {
							Regex: `^522`},
						"양산시": {
							Regex: `^50[56]`},
						"의령군": {
							Regex: `^521`},
						"진주시": {
							Regex: `^52[6-8]`},
						"창녕군": {
							Regex: `^503`},
						"창원시": {
							Regex: `^51[2-7]`,
							SubdivisionRegex: map[string]PostCodeRegexData{
								"마산합포구": {
									Regex: `^51[237]`},
								"마산회원구": {
									Regex: `^51[23]`},
								"성산구": {
									Regex: `^51[457]`},
								"의창구": {
									Regex: `^51[1-4]`},
								"진해구": {
									Regex: `^51[5-7]`},
							}},
						"통영시": {
							Regex: `^53[01]`},
						"하동군": {
							Regex: `^523`},
						"함안군": {
							Regex: `^520`},
						"함양군": {
							Regex: `^500`},
						"합천군": {
							Regex: `^502`},
					}},
				"49": {
					Regex: `^63[0-356]\d`,
					SubdivisionRegex: map[string]PostCodeRegexData{
						"서귀포시": {
							Regex: `^63[56]`},
						"제주시": {
							Regex: `^63[0-3]`},
					}},
				"50": {
					Regex: `^30[01]\d`},
			},
		},
		AdministrativeAreas: map[string][]AdministrativeAreaData{
			"en": {
				{
					ID:   "26",
					Name: "Busan",
					Localities: []LocalityData{
						{
							ID:   "북구",
							Name: "Buk-gu",
						},
						{
							ID:   "부산진구",
							Name: "Busanjin-gu",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "동래구",
							Name: "Dongnae-gu",
						},
						{
							ID:   "강서구",
							Name: "Gangseo-gu",
						},
						{
							ID:   "금정구",
							Name: "Geumjeong-gu",
						},
						{
							ID:   "기장군",
							Name: "Gijang-gun",
						},
						{
							ID:   "해운대구",
							Name: "Haeundae-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "남구",
							Name: "Nam-gu",
						},
						{
							ID:   "사하구",
							Name: "Saha-gu",
						},
						{
							ID:   "사상구",
							Name: "Sasang-gu",
						},
						{
							ID:   "서구",
							Name: "Seo-gu",
						},
						{
							ID:   "수영구",
							Name: "Suyeong-gu",
						},
						{
							ID:   "영도구",
							Name: "Yeongdo-gu",
						},
						{
							ID:   "연제구",
							Name: "Yeonje-gu",
						},
					},
				},
				{
					ID:   "43",
					Name: "Chungcheongbuk-do",
					Localities: []LocalityData{
						{
							ID:   "보은군",
							Name: "Boeun-gun",
						},
						{
							ID:   "청주시",
							Name: "Cheongju-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "청원구",
									Name: "Cheongwon-gu",
								},
								{
									ID:   "흥덕구",
									Name: "Heungdeok-gu",
								},
								{
									ID:   "상당구",
									Name: "Sangdang-gu",
								},
								{
									ID:   "서원구",
									Name: "Seowon-gu",
								},
							},
						},
						{
							ID:   "충주시",
							Name: "Chungju-si",
						},
						{
							ID:   "단양군",
							Name: "Danyang-gun",
						},
						{
							ID:   "음성군",
							Name: "Eumseong-gun",
						},
						{
							ID:   "괴산군",
							Name: "Goesan-gun",
						},
						{
							ID:   "제천시",
							Name: "Jecheon-si",
						},
						{
							ID:   "증평군",
							Name: "Jeungpyeong-gun",
						},
						{
							ID:   "진천군",
							Name: "Jincheon-gun",
						},
						{
							ID:   "옥천군",
							Name: "Okcheon-gun",
						},
						{
							ID:   "영동군",
							Name: "Yeongdong-gun",
						},
					},
				},
				{
					ID:   "44",
					Name: "Chungcheongnam-do",
					Localities: []LocalityData{
						{
							ID:   "아산시",
							Name: "Asan-si",
						},
						{
							ID:   "보령시",
							Name: "Boryeong-si",
						},
						{
							ID:   "부여군",
							Name: "Buyeo-gun",
						},
						{
							ID:   "천안시",
							Name: "Cheonan-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "동남구",
									Name: "Dongnam-gu",
								},
								{
									ID:   "서북구",
									Name: "Seobuk-gu",
								},
							},
						},
						{
							ID:   "청양군",
							Name: "Cheongyang-gun",
						},
						{
							ID:   "당진시",
							Name: "Dangjin-si",
						},
						{
							ID:   "금산군",
							Name: "Geumsan-gun",
						},
						{
							ID:   "공주시",
							Name: "Gongju-si",
						},
						{
							ID:   "계룡시",
							Name: "Gyeryong-si",
						},
						{
							ID:   "홍성군",
							Name: "Hongseong-gun",
						},
						{
							ID:   "논산시",
							Name: "Nonsan-si",
						},
						{
							ID:   "서천군",
							Name: "Seocheon-gun",
						},
						{
							ID:   "서산시",
							Name: "Seosan-si",
						},
						{
							ID:   "태안군",
							Name: "Taean-gun",
						},
						{
							ID:   "예산군",
							Name: "Yesan-gun",
						},
					},
				},
				{
					ID:   "27",
					Name: "Daegu",
					Localities: []LocalityData{
						{
							ID:   "북구",
							Name: "Buk-gu",
						},
						{
							ID:   "달서구",
							Name: "Dalseo-gu",
						},
						{
							ID:   "달성군",
							Name: "Dalseong-gun",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "남구",
							Name: "Nam-gu",
						},
						{
							ID:   "서구",
							Name: "Seo-gu",
						},
						{
							ID:   "수성구",
							Name: "Suseong-gu",
						},
					},
				},
				{
					ID:   "30",
					Name: "Daejeon",
					Localities: []LocalityData{
						{
							ID:   "대덕구",
							Name: "Daedeok-gu",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "서구",
							Name: "Seo-gu",
						},
						{
							ID:   "유성구",
							Name: "Yuseong-gu",
						},
					},
				},
				{
					ID:   "42",
					Name: "Gangwon-do",
					Localities: []LocalityData{
						{
							ID:   "철원군",
							Name: "Cheorwon-gun",
						},
						{
							ID:   "춘천시",
							Name: "Chuncheon-si",
						},
						{
							ID:   "동해시",
							Name: "Donghae-si",
						},
						{
							ID:   "강릉시",
							Name: "Gangneung-si",
						},
						{
							ID:   "고성군",
							Name: "Goseong-gun",
						},
						{
							ID:   "횡성군",
							Name: "Hoengseong-gun",
						},
						{
							ID:   "홍천군",
							Name: "Hongcheon-gun",
						},
						{
							ID:   "화천군",
							Name: "Hwacheon-gun",
						},
						{
							ID:   "인제군",
							Name: "Inje-gun",
						},
						{
							ID:   "정선군",
							Name: "Jeongseon-gun",
						},
						{
							ID:   "평창군",
							Name: "Pyeongchang-gun",
						},
						{
							ID:   "삼척시",
							Name: "Samcheok-si",
						},
						{
							ID:   "속초시",
							Name: "Sokcho-si",
						},
						{
							ID:   "태백시",
							Name: "Taebaek-si",
						},
						{
							ID:   "원주시",
							Name: "Wonju-si",
						},
						{
							ID:   "양구군",
							Name: "Yanggu-gun",
						},
						{
							ID:   "양양군",
							Name: "Yangyang-gun",
						},
						{
							ID:   "영월군",
							Name: "Yeongwol-gun",
						},
					},
				},
				{
					ID:   "29",
					Name: "Gwangju",
					Localities: []LocalityData{
						{
							ID:   "북구",
							Name: "Buk-gu",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "광산구",
							Name: "Gwangsan-gu",
						},
						{
							ID:   "남구",
							Name: "Nam-gu",
						},
						{
							ID:   "서구",
							Name: "Seo-gu",
						},
					},
				},
				{
					ID:   "41",
					Name: "Gyeonggi-do",
					Localities: []LocalityData{
						{
							ID:   "안산시",
							Name: "Ansan-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "단원구",
									Name: "Danwon-gu",
								},
								{
									ID:   "상록구",
									Name: "Sangnok-gu",
								},
							},
						},
						{
							ID:   "안성시",
							Name: "Anseong-si",
						},
						{
							ID:   "안양시",
							Name: "Anyang-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "동안구",
									Name: "Dongan-gu",
								},
								{
									ID:   "만안구",
									Name: "Manan-gu",
								},
							},
						},
						{
							ID:   "부천시",
							Name: "Bucheon-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "오정구",
									Name: "Ojeong-gu",
								},
								{
									ID:   "소사구",
									Name: "Sosa-gu",
								},
								{
									ID:   "원미구",
									Name: "Wonmi-gu",
								},
							},
						},
						{
							ID:   "동두천시",
							Name: "Dongducheon-si",
						},
						{
							ID:   "가평군",
							Name: "Gapyeong-gun",
						},
						{
							ID:   "김포시",
							Name: "Gimpo-si",
						},
						{
							ID:   "고양시",
							Name: "Goyang-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "덕양구",
									Name: "Deogyang-gu",
								},
								{
									ID:   "일산동구",
									Name: "Ilsandong-gu",
								},
								{
									ID:   "일산서구",
									Name: "Ilsanseo-gu",
								},
							},
						},
						{
							ID:   "군포시",
							Name: "Gunpo-si",
						},
						{
							ID:   "구리시",
							Name: "Guri-si",
						},
						{
							ID:   "과천시",
							Name: "Gwacheon-si",
						},
						{
							ID:   "광주시",
							Name: "Gwangju-si",
						},
						{
							ID:   "광명시",
							Name: "Gwangmyeong-si",
						},
						{
							ID:   "하남시",
							Name: "Hanam-si",
						},
						{
							ID:   "화성시",
							Name: "Hwaseong-si",
						},
						{
							ID:   "이천시",
							Name: "Icheon-si",
						},
						{
							ID:   "남양주시",
							Name: "Namyangju-si",
						},
						{
							ID:   "오산시",
							Name: "Osan-si",
						},
						{
							ID:   "파주시",
							Name: "Paju-si",
						},
						{
							ID:   "포천시",
							Name: "Pocheon-si",
						},
						{
							ID:   "평택시",
							Name: "Pyeongtaek-si",
						},
						{
							ID:   "성남시",
							Name: "Seongnam-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "분당구",
									Name: "Bundang-gu",
								},
								{
									ID:   "중원구",
									Name: "Jungwon-gu",
								},
								{
									ID:   "수정구",
									Name: "Sujeong-gu",
								},
							},
						},
						{
							ID:   "시흥시",
							Name: "Siheung-si",
						},
						{
							ID:   "수원시",
							Name: "Suwon-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "권선구",
									Name: "Gwonseon-gu",
								},
								{
									ID:   "장안구",
									Name: "Jangan-gu",
								},
								{
									ID:   "팔달구",
									Name: "Paldal-gu",
								},
								{
									ID:   "영통구",
									Name: "Yeongtong-gu",
								},
							},
						},
						{
							ID:   "의정부시",
							Name: "Uijeongbu-si",
						},
						{
							ID:   "의왕시",
							Name: "Uiwang-si",
						},
						{
							ID:   "양주시",
							Name: "Yangju-si",
						},
						{
							ID:   "양평군",
							Name: "Yangpyeong-gun",
						},
						{
							ID:   "여주시",
							Name: "Yeoju-si",
						},
						{
							ID:   "연천군",
							Name: "Yeoncheon-gun",
						},
						{
							ID:   "용인시",
							Name: "Yongin-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "처인구",
									Name: "Cheoin-gu",
								},
								{
									ID:   "기흥구",
									Name: "Giheung-gu",
								},
								{
									ID:   "수지구",
									Name: "Suji-gu",
								},
							},
						},
					},
				},
				{
					ID:   "47",
					Name: "Gyeongsangbuk-do",
					Localities: []LocalityData{
						{
							ID:   "안동시",
							Name: "Andong-si",
						},
						{
							ID:   "봉화군",
							Name: "Bonghwa-gun",
						},
						{
							ID:   "청도군",
							Name: "Cheongdo-gun",
						},
						{
							ID:   "청송군",
							Name: "Cheongsong-gun",
						},
						{
							ID:   "칠곡군",
							Name: "Chilgok-gun",
						},
						{
							ID:   "김천시",
							Name: "Gimcheon-si",
						},
						{
							ID:   "고령군",
							Name: "Goryeong-gun",
						},
						{
							ID:   "구미시",
							Name: "Gumi-si",
						},
						{
							ID:   "군위군",
							Name: "Gunwi-gun",
						},
						{
							ID:   "경주시",
							Name: "Gyeongju-si",
						},
						{
							ID:   "경산시",
							Name: "Gyeongsan-si",
						},
						{
							ID:   "문경시",
							Name: "Mungyeong-si",
						},
						{
							ID:   "포항시",
							Name: "Pohang-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "북구",
									Name: "Buk-gu",
								},
								{
									ID:   "남구",
									Name: "Nam-gu",
								},
							},
						},
						{
							ID:   "상주시",
							Name: "Sangju-si",
						},
						{
							ID:   "성주군",
							Name: "Seongju-gun",
						},
						{
							ID:   "의성군",
							Name: "Uiseong-gun",
						},
						{
							ID:   "울진군",
							Name: "Uljin-gun",
						},
						{
							ID:   "울릉군",
							Name: "Ulleung-gun",
						},
						{
							ID:   "예천군",
							Name: "Yecheon-gun",
						},
						{
							ID:   "영천시",
							Name: "Yeongcheon-si",
						},
						{
							ID:   "영덕군",
							Name: "Yeongdeok-gun",
						},
						{
							ID:   "영주시",
							Name: "Yeongju-si",
						},
						{
							ID:   "영양군",
							Name: "Yeongyang-gun",
						},
					},
				},
				{
					ID:   "48",
					Name: "Gyeongsangnam-do",
					Localities: []LocalityData{
						{
							ID:   "창녕군",
							Name: "Changnyeong-gun",
						},
						{
							ID:   "창원시",
							Name: "Changwon-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "진해구",
									Name: "Jinhae-gu",
								},
								{
									ID:   "마산합포구",
									Name: "Masanhappo-gu",
								},
								{
									ID:   "마산회원구",
									Name: "Masanhoewon-gu",
								},
								{
									ID:   "성산구",
									Name: "Seongsan-gu",
								},
								{
									ID:   "의창구",
									Name: "Uichang-gu",
								},
							},
						},
						{
							ID:   "거창군",
							Name: "Geochang-gun",
						},
						{
							ID:   "거제시",
							Name: "Geoje-si",
						},
						{
							ID:   "김해시",
							Name: "Gimhae-si",
						},
						{
							ID:   "고성군",
							Name: "Goseong-gun",
						},
						{
							ID:   "하동군",
							Name: "Hadong-gun",
						},
						{
							ID:   "함안군",
							Name: "Haman-gun",
						},
						{
							ID:   "함양군",
							Name: "Hamyang-gun",
						},
						{
							ID:   "합천군",
							Name: "Hapcheon-gun",
						},
						{
							ID:   "진주시",
							Name: "Jinju-si",
						},
						{
							ID:   "밀양시",
							Name: "Miryang-si",
						},
						{
							ID:   "남해군",
							Name: "Namhae-gun",
						},
						{
							ID:   "사천시",
							Name: "Sacheon-si",
						},
						{
							ID:   "산청군",
							Name: "Sancheong-gun",
						},
						{
							ID:   "통영시",
							Name: "Tongyeong-si",
						},
						{
							ID:   "의령군",
							Name: "Uiryeong-gun",
						},
						{
							ID:   "양산시",
							Name: "Yangsan-si",
						},
					},
				},
				{
					ID:   "28",
					Name: "Incheon",
					Localities: []LocalityData{
						{
							ID:   "부평구",
							Name: "Bupyeong-gu",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "강화군",
							Name: "Ganghwa-gun",
						},
						{
							ID:   "계양구",
							Name: "Gyeyang-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "남구",
							Name: "Nam-gu",
						},
						{
							ID:   "남동구",
							Name: "Namdong-gu",
						},
						{
							ID:   "옹진군",
							Name: "Ongjin-gun",
						},
						{
							ID:   "서구",
							Name: "Seo-gu",
						},
						{
							ID:   "연수구",
							Name: "Yeonsu-gu",
						},
					},
				},
				{
					ID:   "49",
					Name: "Jeju-do",
					Localities: []LocalityData{
						{
							ID:   "제주시",
							Name: "Jeju-si",
						},
						{
							ID:   "서귀포시",
							Name: "Seogwipo-si",
						},
					},
				},
				{
					ID:   "45",
					Name: "Jeollabuk-do",
					Localities: []LocalityData{
						{
							ID:   "부안군",
							Name: "Buan-gun",
						},
						{
							ID:   "김제시",
							Name: "Gimje-si",
						},
						{
							ID:   "고창군",
							Name: "Gochang-gun",
						},
						{
							ID:   "군산시",
							Name: "Gunsan-si",
						},
						{
							ID:   "익산시",
							Name: "Iksan-si",
						},
						{
							ID:   "임실군",
							Name: "Imsil-gun",
						},
						{
							ID:   "장수군",
							Name: "Jangsu-gun",
						},
						{
							ID:   "정읍시",
							Name: "Jeongeup-si",
						},
						{
							ID:   "전주시",
							Name: "Jeonju-si",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "덕진구",
									Name: "Deokjin-gu",
								},
								{
									ID:   "완산구",
									Name: "Wansan-gu",
								},
							},
						},
						{
							ID:   "진안군",
							Name: "Jinan-gun",
						},
						{
							ID:   "무주군",
							Name: "Muju-gun",
						},
						{
							ID:   "남원시",
							Name: "Namwon-si",
						},
						{
							ID:   "순창군",
							Name: "Sunchang-gun",
						},
						{
							ID:   "완주군",
							Name: "Wanju-gun",
						},
					},
				},
				{
					ID:   "46",
					Name: "Jeollanam-do",
					Localities: []LocalityData{
						{
							ID:   "보성군",
							Name: "Boseong-gun",
						},
						{
							ID:   "담양군",
							Name: "Damyang-gun",
						},
						{
							ID:   "강진군",
							Name: "Gangjin-gun",
						},
						{
							ID:   "고흥군",
							Name: "Goheung-gun",
						},
						{
							ID:   "곡성군",
							Name: "Gokseong-gun",
						},
						{
							ID:   "구례군",
							Name: "Gurye-gun",
						},
						{
							ID:   "광양시",
							Name: "Gwangyang-si",
						},
						{
							ID:   "해남군",
							Name: "Haenam-gun",
						},
						{
							ID:   "함평군",
							Name: "Hampyeong-gun",
						},
						{
							ID:   "화순군",
							Name: "Hwasun-gun",
						},
						{
							ID:   "장흥군",
							Name: "Jangheung-gun",
						},
						{
							ID:   "장성군",
							Name: "Jangseong-gun",
						},
						{
							ID:   "진도군",
							Name: "Jindo-gun",
						},
						{
							ID:   "목포시",
							Name: "Mokpo-si",
						},
						{
							ID:   "무안군",
							Name: "Muan-gun",
						},
						{
							ID:   "나주시",
							Name: "Naju-si",
						},
						{
							ID:   "신안군",
							Name: "Sinan-gun",
						},
						{
							ID:   "순천시",
							Name: "Suncheon-si",
						},
						{
							ID:   "완도군",
							Name: "Wando-gun",
						},
						{
							ID:   "영암군",
							Name: "Yeongam-gun",
						},
						{
							ID:   "영광군",
							Name: "Yeonggwang-gun",
						},
						{
							ID:   "여수시",
							Name: "Yeosu-si",
						},
					},
				},
				{
					ID:   "50",
					Name: "Sejong",
					Localities: []LocalityData{
						{
							ID:   "아름동",
							Name: "Areum-dong",
						},
						{
							ID:   "보람동",
							Name: "Boram-dong",
						},
						{
							ID:   "부강면",
							Name: "Bugang-myeon",
						},
						{
							ID:   "대평동",
							Name: "Daepyeong-dong",
						},
						{
							ID:   "다정동",
							Name: "Dajeong-dong",
						},
						{
							ID:   "도담동",
							Name: "Dodam-dong",
						},
						{
							ID:   "금남면",
							Name: "Geumnam-myeon",
						},
						{
							ID:   "고운동",
							Name: "Goun-dong",
						},
						{
							ID:   "한솔동",
							Name: "Hansol-dong",
						},
						{
							ID:   "장군면",
							Name: "Janggun-myeon",
						},
						{
							ID:   "전동면",
							Name: "Jeondong-myeon",
						},
						{
							ID:   "전의면",
							Name: "Jeonui-myeon",
						},
						{
							ID:   "조치원읍",
							Name: "Jochiwon-eup",
						},
						{
							ID:   "종촌동",
							Name: "Jongchon-dong",
						},
						{
							ID:   "새롬동",
							Name: "Saerom-dong",
						},
						{
							ID:   "소담동",
							Name: "Sodam-dong",
						},
						{
							ID:   "소정면",
							Name: "Sojeong-myeon",
						},
						{
							ID:   "연동면",
							Name: "Yeondong-myeon",
						},
						{
							ID:   "연기면",
							Name: "Yeongi-myeon",
						},
						{
							ID:   "연서면",
							Name: "Yeonseo-myeon",
						},
					},
				},
				{
					ID:   "11",
					Name: "Seoul",
					Localities: []LocalityData{
						{
							ID:   "도봉구",
							Name: "Dobong-gu",
						},
						{
							ID:   "동대문구",
							Name: "Dongdaemun-gu",
						},
						{
							ID:   "동작구",
							Name: "Dongjak-gu",
						},
						{
							ID:   "은평구",
							Name: "Eunpyeong-gu",
						},
						{
							ID:   "강북구",
							Name: "Gangbuk-gu",
						},
						{
							ID:   "강동구",
							Name: "Gangdong-gu",
						},
						{
							ID:   "강남구",
							Name: "Gangnam-gu",
						},
						{
							ID:   "강서구",
							Name: "Gangseo-gu",
						},
						{
							ID:   "금천구",
							Name: "Geumcheon-gu",
						},
						{
							ID:   "구로구",
							Name: "Guro-gu",
						},
						{
							ID:   "관악구",
							Name: "Gwanak-gu",
						},
						{
							ID:   "광진구",
							Name: "Gwangjin-gu",
						},
						{
							ID:   "종로구",
							Name: "Jongno-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "중랑구",
							Name: "Jungnang-gu",
						},
						{
							ID:   "마포구",
							Name: "Mapo-gu",
						},
						{
							ID:   "노원구",
							Name: "Nowon-gu",
						},
						{
							ID:   "서초구",
							Name: "Seocho-gu",
						},
						{
							ID:   "서대문구",
							Name: "Seodaemun-gu",
						},
						{
							ID:   "성북구",
							Name: "Seongbuk-gu",
						},
						{
							ID:   "성동구",
							Name: "Seongdong-gu",
						},
						{
							ID:   "송파구",
							Name: "Songpa-gu",
						},
						{
							ID:   "양천구",
							Name: "Yangcheon-gu",
						},
						{
							ID:   "영등포구",
							Name: "Yeongdeungpo-gu",
						},
						{
							ID:   "용산구",
							Name: "Yongsan-gu",
						},
					},
				},
				{
					ID:   "31",
					Name: "Ulsan",
					Localities: []LocalityData{
						{
							ID:   "북구",
							Name: "Buk-gu",
						},
						{
							ID:   "동구",
							Name: "Dong-gu",
						},
						{
							ID:   "중구",
							Name: "Jung-gu",
						},
						{
							ID:   "남구",
							Name: "Nam-gu",
						},
						{
							ID:   "울주군",
							Name: "Ulju-gun",
						},
					},
				},
			},
			"ko": {
				{
					ID:   "42",
					Name: "강원",
					Localities: []LocalityData{
						{
							ID:   "강릉시",
							Name: "강릉시",
						},
						{
							ID:   "고성군",
							Name: "고성군",
						},
						{
							ID:   "동해시",
							Name: "동해시",
						},
						{
							ID:   "삼척시",
							Name: "삼척시",
						},
						{
							ID:   "속초시",
							Name: "속초시",
						},
						{
							ID:   "양구군",
							Name: "양구군",
						},
						{
							ID:   "양양군",
							Name: "양양군",
						},
						{
							ID:   "영월군",
							Name: "영월군",
						},
						{
							ID:   "원주시",
							Name: "원주시",
						},
						{
							ID:   "인제군",
							Name: "인제군",
						},
						{
							ID:   "정선군",
							Name: "정선군",
						},
						{
							ID:   "철원군",
							Name: "철원군",
						},
						{
							ID:   "춘천시",
							Name: "춘천시",
						},
						{
							ID:   "태백시",
							Name: "태백시",
						},
						{
							ID:   "평창군",
							Name: "평창군",
						},
						{
							ID:   "홍천군",
							Name: "홍천군",
						},
						{
							ID:   "화천군",
							Name: "화천군",
						},
						{
							ID:   "횡성군",
							Name: "횡성군",
						},
					},
				},
				{
					ID:   "41",
					Name: "경기",
					Localities: []LocalityData{
						{
							ID:   "가평군",
							Name: "가평군",
						},
						{
							ID:   "고양시",
							Name: "고양시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "덕양구",
									Name: "덕양구",
								},
								{
									ID:   "일산동구",
									Name: "일산동구",
								},
								{
									ID:   "일산서구",
									Name: "일산서구",
								},
							},
						},
						{
							ID:   "과천시",
							Name: "과천시",
						},
						{
							ID:   "광명시",
							Name: "광명시",
						},
						{
							ID:   "광주시",
							Name: "광주시",
						},
						{
							ID:   "구리시",
							Name: "구리시",
						},
						{
							ID:   "군포시",
							Name: "군포시",
						},
						{
							ID:   "김포시",
							Name: "김포시",
						},
						{
							ID:   "남양주시",
							Name: "남양주시",
						},
						{
							ID:   "동두천시",
							Name: "동두천시",
						},
						{
							ID:   "부천시",
							Name: "부천시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "소사구",
									Name: "소사구",
								},
								{
									ID:   "오정구",
									Name: "오정구",
								},
								{
									ID:   "원미구",
									Name: "원미구",
								},
							},
						},
						{
							ID:   "성남시",
							Name: "성남시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "분당구",
									Name: "분당구",
								},
								{
									ID:   "수정구",
									Name: "수정구",
								},
								{
									ID:   "중원구",
									Name: "중원구",
								},
							},
						},
						{
							ID:   "수원시",
							Name: "수원시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "권선구",
									Name: "권선구",
								},
								{
									ID:   "영통구",
									Name: "영통구",
								},
								{
									ID:   "장안구",
									Name: "장안구",
								},
								{
									ID:   "팔달구",
									Name: "팔달구",
								},
							},
						},
						{
							ID:   "시흥시",
							Name: "시흥시",
						},
						{
							ID:   "안산시",
							Name: "안산시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "단원구",
									Name: "단원구",
								},
								{
									ID:   "상록구",
									Name: "상록구",
								},
							},
						},
						{
							ID:   "안성시",
							Name: "안성시",
						},
						{
							ID:   "안양시",
							Name: "안양시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "동안구",
									Name: "동안구",
								},
								{
									ID:   "만안구",
									Name: "만안구",
								},
							},
						},
						{
							ID:   "양주시",
							Name: "양주시",
						},
						{
							ID:   "양평군",
							Name: "양평군",
						},
						{
							ID:   "여주시",
							Name: "여주시",
						},
						{
							ID:   "연천군",
							Name: "연천군",
						},
						{
							ID:   "오산시",
							Name: "오산시",
						},
						{
							ID:   "용인시",
							Name: "용인시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "기흥구",
									Name: "기흥구",
								},
								{
									ID:   "수지구",
									Name: "수지구",
								},
								{
									ID:   "처인구",
									Name: "처인구",
								},
							},
						},
						{
							ID:   "의왕시",
							Name: "의왕시",
						},
						{
							ID:   "의정부시",
							Name: "의정부시",
						},
						{
							ID:   "이천시",
							Name: "이천시",
						},
						{
							ID:   "파주시",
							Name: "파주시",
						},
						{
							ID:   "평택시",
							Name: "평택시",
						},
						{
							ID:   "포천시",
							Name: "포천시",
						},
						{
							ID:   "하남시",
							Name: "하남시",
						},
						{
							ID:   "화성시",
							Name: "화성시",
						},
					},
				},
				{
					ID:   "48",
					Name: "경남",
					Localities: []LocalityData{
						{
							ID:   "거제시",
							Name: "거제시",
						},
						{
							ID:   "거창군",
							Name: "거창군",
						},
						{
							ID:   "고성군",
							Name: "고성군",
						},
						{
							ID:   "김해시",
							Name: "김해시",
						},
						{
							ID:   "남해군",
							Name: "남해군",
						},
						{
							ID:   "밀양시",
							Name: "밀양시",
						},
						{
							ID:   "사천시",
							Name: "사천시",
						},
						{
							ID:   "산청군",
							Name: "산청군",
						},
						{
							ID:   "양산시",
							Name: "양산시",
						},
						{
							ID:   "의령군",
							Name: "의령군",
						},
						{
							ID:   "진주시",
							Name: "진주시",
						},
						{
							ID:   "창녕군",
							Name: "창녕군",
						},
						{
							ID:   "창원시",
							Name: "창원시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "마산합포구",
									Name: "마산합포구",
								},
								{
									ID:   "마산회원구",
									Name: "마산회원구",
								},
								{
									ID:   "성산구",
									Name: "성산구",
								},
								{
									ID:   "의창구",
									Name: "의창구",
								},
								{
									ID:   "진해구",
									Name: "진해구",
								},
							},
						},
						{
							ID:   "통영시",
							Name: "통영시",
						},
						{
							ID:   "하동군",
							Name: "하동군",
						},
						{
							ID:   "함안군",
							Name: "함안군",
						},
						{
							ID:   "함양군",
							Name: "함양군",
						},
						{
							ID:   "합천군",
							Name: "합천군",
						},
					},
				},
				{
					ID:   "47",
					Name: "경북",
					Localities: []LocalityData{
						{
							ID:   "경산시",
							Name: "경산시",
						},
						{
							ID:   "경주시",
							Name: "경주시",
						},
						{
							ID:   "고령군",
							Name: "고령군",
						},
						{
							ID:   "구미시",
							Name: "구미시",
						},
						{
							ID:   "군위군",
							Name: "군위군",
						},
						{
							ID:   "김천시",
							Name: "김천시",
						},
						{
							ID:   "문경시",
							Name: "문경시",
						},
						{
							ID:   "봉화군",
							Name: "봉화군",
						},
						{
							ID:   "상주시",
							Name: "상주시",
						},
						{
							ID:   "성주군",
							Name: "성주군",
						},
						{
							ID:   "안동시",
							Name: "안동시",
						},
						{
							ID:   "영덕군",
							Name: "영덕군",
						},
						{
							ID:   "영양군",
							Name: "영양군",
						},
						{
							ID:   "영주시",
							Name: "영주시",
						},
						{
							ID:   "영천시",
							Name: "영천시",
						},
						{
							ID:   "예천군",
							Name: "예천군",
						},
						{
							ID:   "울릉군",
							Name: "울릉군",
						},
						{
							ID:   "울진군",
							Name: "울진군",
						},
						{
							ID:   "의성군",
							Name: "의성군",
						},
						{
							ID:   "청도군",
							Name: "청도군",
						},
						{
							ID:   "청송군",
							Name: "청송군",
						},
						{
							ID:   "칠곡군",
							Name: "칠곡군",
						},
						{
							ID:   "포항시",
							Name: "포항시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "남구",
									Name: "남구",
								},
								{
									ID:   "북구",
									Name: "북구",
								},
							},
						},
					},
				},
				{
					ID:   "29",
					Name: "광주",
					Localities: []LocalityData{
						{
							ID:   "광산구",
							Name: "광산구",
						},
						{
							ID:   "남구",
							Name: "남구",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "북구",
							Name: "북구",
						},
						{
							ID:   "서구",
							Name: "서구",
						},
					},
				},
				{
					ID:   "27",
					Name: "대구",
					Localities: []LocalityData{
						{
							ID:   "남구",
							Name: "남구",
						},
						{
							ID:   "달서구",
							Name: "달서구",
						},
						{
							ID:   "달성군",
							Name: "달성군",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "북구",
							Name: "북구",
						},
						{
							ID:   "서구",
							Name: "서구",
						},
						{
							ID:   "수성구",
							Name: "수성구",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
					},
				},
				{
					ID:   "30",
					Name: "대전",
					Localities: []LocalityData{
						{
							ID:   "대덕구",
							Name: "대덕구",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "서구",
							Name: "서구",
						},
						{
							ID:   "유성구",
							Name: "유성구",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
					},
				},
				{
					ID:   "26",
					Name: "부산",
					Localities: []LocalityData{
						{
							ID:   "강서구",
							Name: "강서구",
						},
						{
							ID:   "금정구",
							Name: "금정구",
						},
						{
							ID:   "기장군",
							Name: "기장군",
						},
						{
							ID:   "남구",
							Name: "남구",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "동래구",
							Name: "동래구",
						},
						{
							ID:   "부산진구",
							Name: "부산진구",
						},
						{
							ID:   "북구",
							Name: "북구",
						},
						{
							ID:   "사상구",
							Name: "사상구",
						},
						{
							ID:   "사하구",
							Name: "사하구",
						},
						{
							ID:   "서구",
							Name: "서구",
						},
						{
							ID:   "수영구",
							Name: "수영구",
						},
						{
							ID:   "연제구",
							Name: "연제구",
						},
						{
							ID:   "영도구",
							Name: "영도구",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
						{
							ID:   "해운대구",
							Name: "해운대구",
						},
					},
				},
				{
					ID:   "11",
					Name: "서울",
					Localities: []LocalityData{
						{
							ID:   "강남구",
							Name: "강남구",
						},
						{
							ID:   "강동구",
							Name: "강동구",
						},
						{
							ID:   "강북구",
							Name: "강북구",
						},
						{
							ID:   "강서구",
							Name: "강서구",
						},
						{
							ID:   "관악구",
							Name: "관악구",
						},
						{
							ID:   "광진구",
							Name: "광진구",
						},
						{
							ID:   "구로구",
							Name: "구로구",
						},
						{
							ID:   "금천구",
							Name: "금천구",
						},
						{
							ID:   "노원구",
							Name: "노원구",
						},
						{
							ID:   "도봉구",
							Name: "도봉구",
						},
						{
							ID:   "동대문구",
							Name: "동대문구",
						},
						{
							ID:   "동작구",
							Name: "동작구",
						},
						{
							ID:   "마포구",
							Name: "마포구",
						},
						{
							ID:   "서대문구",
							Name: "서대문구",
						},
						{
							ID:   "서초구",
							Name: "서초구",
						},
						{
							ID:   "성동구",
							Name: "성동구",
						},
						{
							ID:   "성북구",
							Name: "성북구",
						},
						{
							ID:   "송파구",
							Name: "송파구",
						},
						{
							ID:   "양천구",
							Name: "양천구",
						},
						{
							ID:   "영등포구",
							Name: "영등포구",
						},
						{
							ID:   "용산구",
							Name: "용산구",
						},
						{
							ID:   "은평구",
							Name: "은평구",
						},
						{
							ID:   "종로구",
							Name: "종로구",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
						{
							ID:   "중랑구",
							Name: "중랑구",
						},
					},
				},
				{
					ID:   "50",
					Name: "세종",
					Localities: []LocalityData{
						{
							ID:   "고운동",
							Name: "고운동",
						},
						{
							ID:   "금남면",
							Name: "금남면",
						},
						{
							ID:   "다정동",
							Name: "다정동",
						},
						{
							ID:   "대평동",
							Name: "대평동",
						},
						{
							ID:   "도담동",
							Name: "도담동",
						},
						{
							ID:   "보람동",
							Name: "보람동",
						},
						{
							ID:   "부강면",
							Name: "부강면",
						},
						{
							ID:   "새롬동",
							Name: "새롬동",
						},
						{
							ID:   "소담동",
							Name: "소담동",
						},
						{
							ID:   "소정면",
							Name: "소정면",
						},
						{
							ID:   "아름동",
							Name: "아름동",
						},
						{
							ID:   "연기면",
							Name: "연기면",
						},
						{
							ID:   "연동면",
							Name: "연동면",
						},
						{
							ID:   "연서면",
							Name: "연서면",
						},
						{
							ID:   "장군면",
							Name: "장군면",
						},
						{
							ID:   "전동면",
							Name: "전동면",
						},
						{
							ID:   "전의면",
							Name: "전의면",
						},
						{
							ID:   "조치원읍",
							Name: "조치원읍",
						},
						{
							ID:   "종촌동",
							Name: "종촌동",
						},
						{
							ID:   "한솔동",
							Name: "한솔동",
						},
					},
				},
				{
					ID:   "31",
					Name: "울산",
					Localities: []LocalityData{
						{
							ID:   "남구",
							Name: "남구",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "북구",
							Name: "북구",
						},
						{
							ID:   "울주군",
							Name: "울주군",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
					},
				},
				{
					ID:   "28",
					Name: "인천",
					Localities: []LocalityData{
						{
							ID:   "강화군",
							Name: "강화군",
						},
						{
							ID:   "계양구",
							Name: "계양구",
						},
						{
							ID:   "남구",
							Name: "남구",
						},
						{
							ID:   "남동구",
							Name: "남동구",
						},
						{
							ID:   "동구",
							Name: "동구",
						},
						{
							ID:   "부평구",
							Name: "부평구",
						},
						{
							ID:   "서구",
							Name: "서구",
						},
						{
							ID:   "연수구",
							Name: "연수구",
						},
						{
							ID:   "옹진군",
							Name: "옹진군",
						},
						{
							ID:   "중구",
							Name: "중구",
						},
					},
				},
				{
					ID:   "46",
					Name: "전남",
					Localities: []LocalityData{
						{
							ID:   "강진군",
							Name: "강진군",
						},
						{
							ID:   "고흥군",
							Name: "고흥군",
						},
						{
							ID:   "곡성군",
							Name: "곡성군",
						},
						{
							ID:   "광양시",
							Name: "광양시",
						},
						{
							ID:   "구례군",
							Name: "구례군",
						},
						{
							ID:   "나주시",
							Name: "나주시",
						},
						{
							ID:   "담양군",
							Name: "담양군",
						},
						{
							ID:   "목포시",
							Name: "목포시",
						},
						{
							ID:   "무안군",
							Name: "무안군",
						},
						{
							ID:   "보성군",
							Name: "보성군",
						},
						{
							ID:   "순천시",
							Name: "순천시",
						},
						{
							ID:   "신안군",
							Name: "신안군",
						},
						{
							ID:   "여수시",
							Name: "여수시",
						},
						{
							ID:   "영광군",
							Name: "영광군",
						},
						{
							ID:   "영암군",
							Name: "영암군",
						},
						{
							ID:   "완도군",
							Name: "완도군",
						},
						{
							ID:   "장성군",
							Name: "장성군",
						},
						{
							ID:   "장흥군",
							Name: "장흥군",
						},
						{
							ID:   "진도군",
							Name: "진도군",
						},
						{
							ID:   "함평군",
							Name: "함평군",
						},
						{
							ID:   "해남군",
							Name: "해남군",
						},
						{
							ID:   "화순군",
							Name: "화순군",
						},
					},
				},
				{
					ID:   "45",
					Name: "전북",
					Localities: []LocalityData{
						{
							ID:   "고창군",
							Name: "고창군",
						},
						{
							ID:   "군산시",
							Name: "군산시",
						},
						{
							ID:   "김제시",
							Name: "김제시",
						},
						{
							ID:   "남원시",
							Name: "남원시",
						},
						{
							ID:   "무주군",
							Name: "무주군",
						},
						{
							ID:   "부안군",
							Name: "부안군",
						},
						{
							ID:   "순창군",
							Name: "순창군",
						},
						{
							ID:   "완주군",
							Name: "완주군",
						},
						{
							ID:   "익산시",
							Name: "익산시",
						},
						{
							ID:   "임실군",
							Name: "임실군",
						},
						{
							ID:   "장수군",
							Name: "장수군",
						},
						{
							ID:   "전주시",
							Name: "전주시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "덕진구",
									Name: "덕진구",
								},
								{
									ID:   "완산구",
									Name: "완산구",
								},
							},
						},
						{
							ID:   "정읍시",
							Name: "정읍시",
						},
						{
							ID:   "진안군",
							Name: "진안군",
						},
					},
				},
				{
					ID:   "49",
					Name: "제주",
					Localities: []LocalityData{
						{
							ID:   "서귀포시",
							Name: "서귀포시",
						},
						{
							ID:   "제주시",
							Name: "제주시",
						},
					},
				},
				{
					ID:   "44",
					Name: "충남",
					Localities: []LocalityData{
						{
							ID:   "계룡시",
							Name: "계룡시",
						},
						{
							ID:   "공주시",
							Name: "공주시",
						},
						{
							ID:   "금산군",
							Name: "금산군",
						},
						{
							ID:   "논산시",
							Name: "논산시",
						},
						{
							ID:   "당진시",
							Name: "당진시",
						},
						{
							ID:   "보령시",
							Name: "보령시",
						},
						{
							ID:   "부여군",
							Name: "부여군",
						},
						{
							ID:   "서산시",
							Name: "서산시",
						},
						{
							ID:   "서천군",
							Name: "서천군",
						},
						{
							ID:   "아산시",
							Name: "아산시",
						},
						{
							ID:   "예산군",
							Name: "예산군",
						},
						{
							ID:   "천안시",
							Name: "천안시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "동남구",
									Name: "동남구",
								},
								{
									ID:   "서북구",
									Name: "서북구",
								},
							},
						},
						{
							ID:   "청양군",
							Name: "청양군",
						},
						{
							ID:   "태안군",
							Name: "태안군",
						},
						{
							ID:   "홍성군",
							Name: "홍성군",
						},
					},
				},
				{
					ID:   "43",
					Name: "충북",
					Localities: []LocalityData{
						{
							ID:   "괴산군",
							Name: "괴산군",
						},
						{
							ID:   "단양군",
							Name: "단양군",
						},
						{
							ID:   "보은군",
							Name: "보은군",
						},
						{
							ID:   "영동군",
							Name: "영동군",
						},
						{
							ID:   "옥천군",
							Name: "옥천군",
						},
						{
							ID:   "음성군",
							Name: "음성군",
						},
						{
							ID:   "제천시",
							Name: "제천시",
						},
						{
							ID:   "증평군",
							Name: "증평군",
						},
						{
							ID:   "진천군",
							Name: "진천군",
						},
						{
							ID:   "청주시",
							Name: "청주시",
							DependentLocalities: []DependentLocalityData{
								{
									ID:   "상당구",
									Name: "상당구",
								},
								{
									ID:   "서원구",
									Name: "서원구",
								},
								{
									ID:   "청원구",
									Name: "청원구",
								},
								{
									ID:   "흥덕구",
									Name: "흥덕구",
								},
							},
						},
						{
							ID:   "충주시",
							Name: "충주시",
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(country, expected) {
		t.Errorf("Country data for KR does not match expected country data")
	}
}
func TestValidErr(t *testing.T) {
	_, err := NewValid(
		WithCountry("AasdgasdgU"), // Must be an ISO 3166-1 country code
		WithName("John Citizen"),
		WithOrganization("Some Company Pty Ltd"),
		WithStreetAddress([]string{
			"525 Collins Street",
		}),
		WithLocality("Melbourasdafsdgne"),
		WithAdministrativeArea("VIC"), // If the country has a pre-defined list of admin areas (like here), you must use the key and not the name
		WithPostCode("3000"),
	)
	if err != nil {
		// If there was an error and you want to find out which validations failed,
		// type switch it as a *multierror.Error to access the list of errors
		err = errors.Unwrap(err)
		_, ok := err.(*multierror.Error)
		if !ok {
			t.Errorf("Expected a *multierror.Error, got %T", err)
		}
	}
}
