package address

import "testing"

func TestConstantKeys(t *testing.T) {

	testCases := []struct {
		Field    Field
		Expected string
	}{
		{
			Field:    Country,
			Expected: "country",
		},
		{
			Field:    Name,
			Expected: "N",
		},
		{
			Field:    Organization,
			Expected: "O",
		},
		{
			Field:    StreetAddress,
			Expected: "A",
		},
		{
			Field:    DependentLocality,
			Expected: "D",
		},
		{
			Field:    Locality,
			Expected: "C",
		},
		{
			Field:    AdministrativeArea,
			Expected: "S",
		},
		{
			Field:    PostCode,
			Expected: "Z",
		},
		{
			Field:    SortingCode,
			Expected: "X",
		},
	}

	for _, testCase := range testCases {
		if testCase.Field.Key() != testCase.Expected {
			t.Errorf("Expected key for %s to be %s, got %s", testCase.Field, testCase.Expected, testCase.Field.Key())
		}
	}
}
