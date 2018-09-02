package address

import (
	"testing"
)

func TestDefaultFormatter(t *testing.T) {

	f := DefaultFormatter{
		Output: StringOutputter{},
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: "Company Pty Ltd\nJohn Smith\n525 Collins Street\nMelbourne Victoria 3000\nAustralia",
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"Suite 7, 9th Floor",
					"525 Collins Street",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
			Expected: "Company Pty Ltd\nJohn Smith\nSuite 7, 9th Floor\n525 Collins Street\nMelbourne Victoria 3000\nAustralia",
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Company Pty Ltd"),
				WithStreetAddress([]string{
					"",
					"Suite 8, 9th Floor",
					"",
					"525 Collins Street",
					"",
				}),
				WithLocality("Melbourne"),
				WithAdministrativeArea("VIC"),
				WithPostCode("3000"),
				WithCountry("AU"),
			},
			Expected: "Company Pty Ltd\nJohn Smith\nSuite 8, 9th Floor\n525 Collins Street\nMelbourne Victoria 3000\nAustralia",
		},
		{
			Address: []func(*Address){
				WithName("John Smith"),
				WithOrganization("Microsoft"),
				WithStreetAddress([]string{
					"1 Microsoft Way",
				}),
				WithLocality("Redmond"),
				WithAdministrativeArea("WA"),
				WithPostCode("98052"),
				WithCountry("US"),
			},
			Expected: "John Smith\nMicrosoft\n1 Microsoft Way\nRedmond, Washington 98052\nUnited States",
		},
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithDependentLocality("临翔区"),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "中国\n677000\n云南省临沧市临翔区\n1 西河北路\n星巴克\n司馬遷",
		},
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "中国\n677000\n云南省临沧市\n1 西河北路\n星巴克\n司馬遷",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestLatinization(t *testing.T) {

	f := DefaultFormatter{
		Output:   StringOutputter{},
		Latinize: true,
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: "Company Pty Ltd\nJohn Smith\n525 Collins Street\nMelbourne Victoria 3000\nAustralia",
		},
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithDependentLocality("临翔区"),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "司馬遷\n星巴克\n1 西河北路\n临翔区\n临沧市\n云南省, 677000\n中国",
		},
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "司馬遷\n星巴克\n1 西河北路\n临沧市\n云南省, 677000\n中国",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestLanguage(t *testing.T) {
	f := DefaultFormatter{
		Output:   StringOutputter{},
		Latinize: true,
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: "Company Pty Ltd\nJohn Smith\n525 Collins Street\nMelbourne Victoria 3000\nAustralia",
		},
		{
			Address: []func(*Address){
				WithName("Sima Qian"),
				WithOrganization("Starbucks"),
				WithStreetAddress([]string{
					"1 Xihe N. Road",
				}),
				WithDependentLocality("临翔区"),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "Sima Qian\nStarbucks\n1 Xihe N. Road\nLinxiang Qu\nLincang Shi\nYunnan Sheng, 677000\nChina",
		},
		{
			Address: []func(*Address){
				WithName("Sima Qian"),
				WithOrganization("Starbucks"),
				WithStreetAddress([]string{
					"1 Xihe N. Road",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "Sima Qian\nStarbucks\n1 Xihe N. Road\nLincang Shi\nYunnan Sheng, 677000\nChina",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "en")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestHTMLOutputter(t *testing.T) {

	f := DefaultFormatter{
		Output: HTMLOutputter{},
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: `<span class="organization">Company Pty Ltd</span><br><span class="name">John Smith</span><br><span class="address-line-1">525 Collins Street</span><br><span class="locality">Melbourne</span> <span class="administrative-area">Victoria</span> <span class="post-code">3000</span><br><span class="country">Australia</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Sima Qian"),
				WithOrganization("Starbucks"),
				WithStreetAddress([]string{
					"1 Xihe N. Road",
				}),
				WithDependentLocality("临翔区"),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: `<span class="country">中国</span><br><span class="post-code">677000</span><br><span class="administrative-area">云南省</span><span class="locality">临沧市</span><span class="dependent-locality">临翔区</span><br><span class="address-line-1">1 Xihe N. Road</span><br><span class="organization">Starbucks</span><br><span class="name">Sima Qian</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Sima Qian"),
				WithOrganization("Starbucks"),
				WithStreetAddress([]string{
					"1 Xihe N. Road",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: `<span class="country">中国</span><br><span class="post-code">677000</span><br><span class="administrative-area">云南省</span><span class="locality">临沧市</span><br><span class="address-line-1">1 Xihe N. Road</span><br><span class="organization">Starbucks</span><br><span class="name">Sima Qian</span>`,
		},
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: `<span class="country">中国</span><br><span class="post-code">677000</span><br><span class="administrative-area">云南省</span><span class="locality">临沧市</span><br><span class="address-line-1">1 西河北路</span><br><span class="organization">星巴克</span><br><span class="name">司馬遷</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Mose Boreham"),
				WithStreetAddress([]string{
					"P.O. Box 33",
				}),
				WithLocality("Vaiaku"),
				WithAdministrativeArea("FUN"),
				WithCountry("TV"),
			},
			Expected: `<span class="name">Mose Boreham</span><br><span class="address-line-1">P.O. Box 33</span><br><span class="locality">Vaiaku</span><br><span class="administrative-area">Funafuti</span><br><span class="country">Tuvalu</span>`,
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestHTMLCollapseWhitespace(t *testing.T) {
	f := DefaultFormatter{
		Output: HTMLOutputter{},
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
	}{
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: `<span class="country">中国</span><br><span class="post-code">677000</span><br><span class="administrative-area">云南省</span><span class="locality">临沧市</span><br><span class="address-line-1">1 西河北路</span><br><span class="organization">星巴克</span><br><span class="name">司馬遷</span>`,
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}

	lf := DefaultFormatter{
		Output:   HTMLOutputter{},
		Latinize: true,
	}

	latinizedTestCases := []struct {
		Address  []func(*Address)
		Expected string
	}{
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: `<span class="name">司馬遷</span><br><span class="organization">星巴克</span><br><span class="address-line-1">1 西河北路</span><br><span class="locality">临沧市</span><br><span class="administrative-area">云南省</span>, <span class="post-code">677000</span><br><span class="country">中国</span>`,
		},
	}

	for i, testCase := range latinizedTestCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using latinized test case %d: %s", i, err)
		}

		formatted := lf.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for latinized test case %d does not match the expected result", i)
		}
	}
}

func TestStringCollapseWhitespace(t *testing.T) {

	f := DefaultFormatter{
		Output: StringOutputter{},
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
	}{
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "中国\n677000\n云南省临沧市\n1 西河北路\n星巴克\n司馬遷",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}

	lf := DefaultFormatter{
		Output:   StringOutputter{},
		Latinize: true,
	}

	latinizedTestCases := []struct {
		Address  []func(*Address)
		Expected string
	}{
		{
			Address: []func(*Address){
				WithName("司馬遷"),
				WithOrganization("星巴克"),
				WithStreetAddress([]string{
					"1 西河北路",
				}),
				WithLocality("临沧市"),
				WithAdministrativeArea("53"),
				WithPostCode("677000"),
				WithCountry("CN"),
			},
			Expected: "司馬遷\n星巴克\n1 西河北路\n临沧市\n云南省, 677000\n中国",
		},
	}

	for i, testCase := range latinizedTestCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using latinized test case %d: %s", i, err)
		}

		formatted := lf.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for latinized test case %d does not match the expected result", i)
		}
	}
}

func TestPostalLabelFormatter(t *testing.T) {

	f := PostalLabelFormatter{
		Output:            StringOutputter{},
		OriginCountryCode: "FR",
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: "Company Pty Ltd\nJohn Smith\n525 Collins Street\nMELBOURNE VIC 3000\nAUSTRALIE - AUSTRALIA",
		},
		{
			Address: []func(*Address){
				WithName("Jean Dupont"),
				WithOrganization("Microsoft"),
				WithStreetAddress([]string{
					"27 Rue Pasteur",
				}),
				WithLocality("Cabourg"),
				WithPostCode("14390"),
				WithCountry("FR"),
			},
			Expected: "Microsoft\nJean Dupont\n27 Rue Pasteur\n14390 CABOURG",
		},
		{
			Address: []func(*Address){
				WithName("Walter C. Brown"),
				WithStreetAddress([]string{
					"49 Featherstone Street",
				}),
				WithLocality("London"),
				WithPostCode("EC1Y 8SY"),
				WithCountry("GB"),
			},
			Expected: "Walter C. Brown\n49 Featherstone Street\nLONDON\nEC1Y 8SY\nROYAUME-UNI - UNITED KINGDOM",
		},
		{
			Address: []func(*Address){
				WithName("Mose Boreham"),
				WithStreetAddress([]string{
					"P.O. Box 33",
				}),
				WithLocality("Vaiaku"),
				WithAdministrativeArea("FUN"),
				WithCountry("TV"),
			},
			Expected: "Mose Boreham\nP.O. BOX 33\nVAIAKU\nFUNAFUTI\nTUVALU",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestPostalLabelFormatterWithHTMLOutput(t *testing.T) {
	f := PostalLabelFormatter{
		Output:            HTMLOutputter{},
		OriginCountryCode: "FR",
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: `<span class="organization">Company Pty Ltd</span><br><span class="name">John Smith</span><br><span class="address-line-1">525 Collins Street</span><br><span class="locality">MELBOURNE</span> <span class="administrative-area">VIC</span> <span class="post-code">3000</span><br><span class="country">AUSTRALIE - AUSTRALIA</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Jean Dupont"),
				WithOrganization("Microsoft"),
				WithStreetAddress([]string{
					"27 Rue Pasteur",
				}),
				WithLocality("Cabourg"),
				WithPostCode("14390"),
				WithCountry("FR"),
			},
			Expected: `<span class="organization">Microsoft</span><br><span class="name">Jean Dupont</span><br><span class="address-line-1">27 Rue Pasteur</span><br><span class="post-code">14390</span> <span class="locality">CABOURG</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Walter C. Brown"),
				WithStreetAddress([]string{
					"49 Featherstone Street",
				}),
				WithLocality("London"),
				WithPostCode("EC1Y 8SY"),
				WithCountry("GB"),
			},
			Expected: `<span class="name">Walter C. Brown</span><br><span class="address-line-1">49 Featherstone Street</span><br><span class="locality">LONDON</span><br><span class="post-code">EC1Y 8SY</span><br><span class="country">ROYAUME-UNI - UNITED KINGDOM</span>`,
		},
		{
			Address: []func(*Address){
				WithName("Mose Boreham"),
				WithStreetAddress([]string{
					"P.O. Box 33",
				}),
				WithLocality("Vaiaku"),
				WithAdministrativeArea("FUN"),
				WithCountry("TV"),
			},
			Expected: `<span class="name">Mose Boreham</span><br><span class="address-line-1">P.O. BOX 33</span><br><span class="locality">VAIAKU</span><br><span class="administrative-area">FUNAFUTI</span><br><span class="country">TUVALU</span>`,
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}

func TestPostalLabelFormatterForInternationalMailWithoutValidOriginCountry(t *testing.T) {

	f := PostalLabelFormatter{
		Output:            StringOutputter{},
		OriginCountryCode: "01234",
	}

	testCases := []struct {
		Address  []func(*Address)
		Expected string
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
			Expected: "Company Pty Ltd\nJohn Smith\n525 Collins Street\nMELBOURNE VIC 3000",
		},
		{
			Address: []func(*Address){
				WithName("Jean Dupont"),
				WithOrganization("Microsoft"),
				WithStreetAddress([]string{
					"27 Rue Pasteur",
				}),
				WithLocality("Cabourg"),
				WithPostCode("14390"),
				WithCountry("FR"),
			},
			Expected: "Microsoft\nJean Dupont\n27 Rue Pasteur\n14390 CABOURG",
		},
		{
			Address: []func(*Address){
				WithName("Walter C. Brown"),
				WithStreetAddress([]string{
					"49 Featherstone Street",
				}),
				WithLocality("London"),
				WithPostCode("EC1Y 8SY"),
				WithCountry("GB"),
			},
			Expected: "Walter C. Brown\n49 Featherstone Street\nLONDON\nEC1Y 8SY",
		},
		{
			Address: []func(*Address){
				WithName("Mose Boreham"),
				WithStreetAddress([]string{
					"P.O. Box 33",
				}),
				WithLocality("Vaiaku"),
				WithAdministrativeArea("FUN"),
				WithCountry("TV"),
			},
			Expected: "Mose Boreham\nP.O. BOX 33\nVAIAKU\nFUNAFUTI",
		},
	}

	for i, testCase := range testCases {

		address, err := NewValid(testCase.Address...)

		if err != nil {
			t.Fatalf("Error creating valid address using test case %d: %s", i, err)
		}

		formatted := f.Format(address, "")

		if formatted != testCase.Expected {
			t.Errorf("Formatted address for test case %d does not match the expected result", i)
		}
	}
}
