package address

import (
	"regexp"
	"testing"
)

func TestZones(t *testing.T) {

	address1, err := New(
		WithCountry("AU"),
		WithStreetAddress(
			[]string{"525 Collins Street"},
		),
		WithLocality("Melbourne"),
		WithAdministrativeArea("VIC"),
		WithPostCode("3000"),
	)

	if err != nil {
		t.Fatalf("Error creating address 1: %s", err)
	}

	address2, err := New(
		WithCountry("CN"),
		WithStreetAddress([]string{
			"1 西河北路",
		}),
		WithDependentLocality("临翔区"),
		WithLocality("临沧市"),
		WithAdministrativeArea("53"),
		WithPostCode("677000"),
	)

	if err != nil {
		t.Fatalf("Error creating address 2: %s", err)
	}

	testCases := []struct {
		Zone     Zone
		Address  Address
		Expected bool
	}{
		{
			Zone: Zone{
				{
					Country: "AU",
				},
			},

			Address:  address1,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
				},
			},

			Address:  address1,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "NSW",
				},
			},

			Address:  address1,
			Expected: false,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
				},
			},

			Address:  address1,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					IncludedPostCodes: ExactMatcher{
						Matches: []string{
							"3200", "3171",
						},
					},
				},
			},

			Address:  address1,
			Expected: false,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					IncludedPostCodes: ExactMatcher{
						Matches: []string{
							"3100", "3200",
						},
						Ranges: []PostCodeRange{
							{
								Start: 3000,
								End:   3020,
							},
						},
					},
				},
			},

			Address:  address1,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					IncludedPostCodes: ExactMatcher{
						Matches: []string{
							"3000", "3200",
						},
						Ranges: []PostCodeRange{
							{
								Start: 3000,
								End:   3020,
							},
						},
					},
					ExcludedPostCodes: ExactMatcher{
						Matches: []string{
							"3000",
						},
					},
				},
			},

			Address:  address1,
			Expected: false,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					IncludedPostCodes: ExactMatcher{
						Matches: []string{
							"3000", "3200",
						},
						Ranges: []PostCodeRange{
							{
								Start: 3000,
								End:   3020,
							},
						},
					},
					ExcludedPostCodes: ExactMatcher{
						Matches: []string{
							"3019",
						},
						Ranges: []PostCodeRange{
							{
								Start: 3000,
								End:   3011,
							},
							{
								Start: 3000,
								End:   3020,
							},
						},
					},
				},
			},

			Address:  address1,
			Expected: false,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					ExcludedPostCodes: RegexMatcher{
						Regex: regexp.MustCompile(`3\d{3}`),
					},
				},
			},

			Address:  address1,
			Expected: false,
		},
		{
			Zone: Zone{
				{
					Country:            "AU",
					AdministrativeArea: "VIC",
					Locality:           "MELBOURNE",
					IncludedPostCodes: RegexMatcher{
						Regex: regexp.MustCompile(`3\d{3}`),
					},
				},
			},

			Address:  address1,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country: "CN",
				},
			},

			Address:  address2,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "CN",
					AdministrativeArea: "53",
				},
			},

			Address:  address2,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "CN",
					AdministrativeArea: "53",
					Locality:           "临沧市",
				},
			},

			Address:  address2,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:            "CN",
					AdministrativeArea: "53",
					Locality:           "临沧市",
					DependentLocality:  "临翔区",
				},
			},

			Address:  address2,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:  "CN",
					Locality: "临沧市",
				},
			},

			Address:  address2,
			Expected: true,
		},
		{
			Zone: Zone{
				{
					Country:           "CN",
					DependentLocality: "临翔区",
				},
			},

			Address:  address2,
			Expected: true,
		},
	}

	for i, testCase := range testCases {
		result := testCase.Zone.Contains(testCase.Address)

		if result != testCase.Expected {
			t.Errorf("Containment of address in zone does not match expected result for test %d", i)
		}
	}
}
