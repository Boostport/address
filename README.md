# Address
[![GoDoc](https://godoc.org/github.com/Boostport/address?status.png)](https://godoc.org/github.com/Boostport/address)
[![wercker status](https://app.wercker.com/status/be8486d6d2d1f4baac08a701ee937863/s/master "wercker status")](https://app.wercker.com/project/byKey/be8486d6d2d1f4baac08a701ee937863)
[![Test Coverage](https://api.codeclimate.com/v1/badges/66b4554d26c85121a829/test_coverage)](https://codeclimate.com/github/Boostport/address/test_coverage)

Address is a Go library that validates and formats addresses using data generated from [Google's Address Data Service](https://chromium-i18n.appspot.com/ssl-address).

## Installation
Install the library using Go modules. For example: `go get -u github.com/Boostport/address`.

## Creating Addresses
To create a new address, use `New()`. If the address is invalid, an error will be returned.

```go
package main

import (
	"log"

	"github.com/Boostport/address"
	"github.com/hashicorp/go-multierror"
)

func main() {
	addr := address.New(
		address.WithCountry("AU"), // Must be an ISO 3166-1 country code
		address.WithName("John Citizen"),
		address.WithOrganization("Some Company Pty Ltd"),
		address.WithStreetAddress([]string{
			"525 Collins Street",
		}),
		address.WithLocality("Melbourne"),
		address.WithAdministrativeArea("VIC"), // If the country has a pre-defined list of admin areas (like here), you must use the key and not the name
		address.WithPostCode("3000"),
	)

	if err != nil {
		// If there was an error and you want to find out which validations failed,
		// type switch it as a *multierror.Error to access the list of errors
		if merr, ok := err.(*multierror.Error); ok {
			for _, subErr := range merr.Errors {
				if subErr == address.ErrInvalidCountryCode {
					log.Fatalf(subErr)
				}
			}
		}
	}

	// Use addr here
}
```

### A note about administrative areas, localities and dependent localities
An address may contain the following subdivisions:
- Administrative areas, such as a state, province, island, etc.
- Localities such as cities.
- Dependent localities such as districts, suburbs, etc.

When creating an address, certain countries have a pre-defined list of administrative areas, localities and dependent localities.
In these cases, you **MUST** use the appropriate key when calling `WithAdministrativeArea()`, `WithLocality()` and `WithDependentLocality()`,
otherwise, the address will fail validation.

In terms of the keys, for administrative areas we use the ISO 3166-2 subdivision codes from Google's data set where possible. If there is no
ISO 3166-2 code available, we use the key defined in Google's data set. In these cases, the key is a unicode string (could be in languages other than 
English).

For localities and dependent localities, there are generally no ISO codes, so we use the key defined in Google's data set. The key is a unicode string
and can be in a language other than English.

The reason for doing this is that when storing an address into a database, we need to store the values in a canonical form. Since these keys are
very stable (in general), they are safe to store. If we need to provide a visual representation of the address, we can then use the key and a language
to choose the appropriate display names.

This also allows us to do things such as rendering Canadian addresses both in French and English using a canonical address.

The library contains helpers where you can access these keys and the display names in different languages. More information available [below](#address-data-format).

## Formatting Addresses
There are 2 formatters, the `DefaultFormatter` and a `PostalLabelFormatter`.

In addition, there 2 outputters, the `StringOutputter` and the `HTMLOutputter`. The outputter takes the formatted
address from the formatters and turn them into their respective string or HTML representations. The `Outputter` is
an interface, so it's possible to implement your own version of the outputter if desired.

In some countries such as China, the address is formatted as major-to-minor (i.e. country -> administrative division -> locality ...).
It's possible to format it using a latinized format (address -> dependent locality -> locality ...) by setting the `Latinize` field in
the formatter to `true`.

This [example](examples/formatters/main.go) shows the difference between the 2 formatters and outputters (error checking omitted for brevity):

```go
package main

import (
	"fmt"

	"github.com/Boostport/address"
)

func main() {

	addr, _ := address.NewValid(
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

	lang := "en" // Use the English names of the administrative areas, localities and dependent localities where possible

	fmt.Println(defStringFormatter.Format(addr, lang))
	/* Output
	Some Company Pty Ltd
	John Citizen
	525 Collins Street
	Melbourne Victoria 3000
	Australia
	*/
	
	fmt.Println(defHTMLFormatter.Format(addr, lang))
	/* Output
	<span class="organization">Some Company Pty Ltd</span><br>
	<span class="name">John Citizen</span><br>
	<span class="address-line-1">525 Collins Street</span><br>
	<span class="locality">Melbourne</span> <span class="administrative-area">Victoria</span> <span class="post-code">3000</span><br>
	<span class="country">Australia</span>
	*/
	
	fmt.Println(postalStringFormatter.Format(addr, lang))
	/* Output
	Some Company Pty Ltd
	John Citizen
	525 Collins Street
	MELBOURNE VIC 3000
	AUSTRALIE - AUSTRALIA
	*/
	
	fmt.Println(postalHTMLFormatter.Format(addr, lang))
	/* Output
	<span class="organization">Some Company Pty Ltd</span><br>
	<span class="name">John Citizen</span><br>
	<span class="address-line-1">525 Collins Street</span><br>
	<span class="locality">MELBOURNE</span> <span class="administrative-area">VIC</span> <span class="post-code">3000</span><br>
	<span class="country">AUSTRALIE - AUSTRALIA</span>
	*/
}
```
## Zones
Zones are useful for calculating things like shipping costs or tax rates. A `Zone` consists of multiple territories, with
each `Territory` equivalent to a rule.

Territories are able to match addresses based on their `Country`, `AdministrativeArea`, `Locality`, `DependentLocality` and `PostCode`.

Note that the `Country` must be an ISO 3166-1 country code, and if there are pre-defined lists of `AdministrativeArea`s, `Locality`, and `DependentLocality`
for the country, the key must be used.

A quick [example](examples/zones/main.go):
```go
package main

import (
	"fmt"

	"github.com/Boostport/address"
)

func main() {
	addr, _ := address.NewValid(
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
```

## Address Data Format
In a lot of cases, you might need to display a form to the user to enter their address.

There is the `ListCountries()` method to get a list of available countries in your chosen language
and the `GetCountry()` method to get detailed address format information for a given country.

`GetCountry()` returns a struct like so:
```go
type CountryData struct {
	Format                     string
	LatinizedFormat            string
	Required                   []Field
	Allowed                    []Field
	DefaultLanguage            string
	AdministrativeAreaNameType FieldName
	LocalityNameType           FieldName
	DependentLocalityNameType  FieldName
	PostCodeNameType           FieldName
	PostCodeRegex              PostCodeRegexData
	AdministrativeAreas        map[string][]AdministrativeAreaData
}
```

The `Format` and `LatinizedFormat` fields are in Google's original formats (ex: `%O%n%N%n%A%n%C %S %Z` for Australia).
A description of what the tokens represent is available [here](https://github.com/googlei18n/libaddressinput/wiki/AddressValidationMetadata).

`Required` and `Allowed` represent fields that are required and allowed (not all allowed fields are required). The `Field` type
can be converted to Google's token name by calling the `Key()` method.

For administrative areas, the map contains a list of administrative areas grouped by the language they are in (the map's key).
Each list is sorted according to the language they are in. Administrative areas may contain localities, and localities
may contain dependent localities. In all cases, each element would have an ID that you should use when creating an address or a zone.

There may also be post code validation regex. There may be further structs nested inside to validate post codes for an administrative area,
locality or dependent locality. These are keyed using the appropriate ID from the list of administrative areas.


## Generating Data
### Directly in your environment
Install stringer: `go get -u golang.org/x/tools/cmd/stringer`.

To generate the data and generate the `String()` functions for the constants, simply run `go generate` from the root of the project.
This will run stringer and the generator which will download the data from Google and convert the data into Go code.

### Using docker
Run `docker-compose run generate`

## License
This library is licensed under the Apache 2 License.
