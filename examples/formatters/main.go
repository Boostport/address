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

	defStringFormatter := address.DefaultFormatter{
		Output: address.StringOutputter{},
	}

	defHTMLFormatter := address.DefaultFormatter{
		Output: address.HTMLOutputter{},
	}

	postalStringFormatter := address.PostalLabelFormatter{
		Output:            address.StringOutputter{},
		OriginCountryCode: "FR", // We are sending from France
	}

	postalHTMLFormatter := address.PostalLabelFormatter{
		Output:            address.HTMLOutputter{},
		OriginCountryCode: "FR", // We are sending from France
	}

	lang := "en" // Use the english names of the administrative areas, localityies and dependent localities where possible

	fmt.Println(defStringFormatter.Format(addr, lang))

	fmt.Println(defHTMLFormatter.Format(addr, lang))

	fmt.Println(postalStringFormatter.Format(addr, lang))

	fmt.Println(postalHTMLFormatter.Format(addr, lang))
}
