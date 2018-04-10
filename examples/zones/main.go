package main

import (
	"fmt"

	"github.com/Boostport/address"
)

func main() {
	addr, _ := address.New(
		address.WithCountry("AU"),
		address.WithName("John Citizen"),
		address.WithOrganization("Some Company Pty Ltd"),
		address.WithStreetAddress([]string{
			"525 Collins Street",
		}),
		address.WithLocality("Melbourne"),
		address.WithAdministrativeArea("VIC"),
		address.WithPostCode("3000"),
	)

	freeShippingToQLDAndNSW := address.Zone{
		{
			Country:            "AU",
			AdministrativeArea: "NSW",
		},
		{
			Country:            "AU",
			AdministrativeArea: "QLD",
		},
	}

	fmt.Println(freeShippingToQLDAndNSW.Contains(addr)) // false

	victorianPostCodesExceptCarltonGetDiscount := address.Zone{
		{
			Country: "AU",
			IncludedPostCodes: address.ExactMatcher{
				Ranges: []address.PostCodeRange{
					{
						Start: 3000,
						End:   3996,
					},
					{
						Start: 8000,
						End:   8873,
					},
				},
			},
			ExcludedPostCodes: address.ExactMatcher{
				Matches: []string{"3053"},
			},
		},
	}

	fmt.Println(victorianPostCodesExceptCarltonGetDiscount.Contains(addr)) // true
}
