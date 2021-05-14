package nested

import (
	"reflect"
	"testing"
)

func TestGetIndexValue(t *testing.T) {
	tests := []struct {
		IndexString string
		Want        IndexValue
	}{
		{
			IndexString: "4",
			Want: IndexValue{
				Type:         INDEX_NUMERIC,
				NumericValue: 4,
			},
		},
		{
			IndexString: `?(@.name=="my name")`,
			Want: IndexValue{
				Type:         INDEX_EQUALS,
				Path:         []string{"name"},
				CompareValue: `"my name"`,
			},
		},
	}

	for _, test := range tests {
		result := getIndexValue(test.IndexString)
		if !reflect.DeepEqual(result, test.Want) {
			t.Errorf("Expecting %+v, got %+v\n", test.Want, result)
		}
	}
}

func TestFromDotted(t *testing.T) {
	tests := []struct {
		Dotted string
		Want   []string
	}{
		{
			Dotted: ".a.b.c",
			Want:   []string{"a", "b", "c"},
		},
		{
			Dotted: "a.b.c",
			Want:   []string{"a", "b", "c"},
		},
		{
			Dotted: ".a[1].b[0].c",
			Want:   []string{"a[1]", "b[0]", "c"},
		},
	}

	for _, test := range tests {
		result := fromDotted(test.Dotted)
		if !stringSlicesEqual(result, test.Want) {
			t.Errorf("Expecting %s, got %s\n", test.Want, result)
		}
	}
}
