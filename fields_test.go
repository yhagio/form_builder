package form

import (
	"fmt"
	"testing"
)

func TestFields(t *testing.T) {
	type form struct {
		Name string
	}

	tests := []struct {
		strct interface{}
		want  field
	}{
		{
			strct: struct {
				Name string
			}{},
			want: field{
				Label:       "Label",
				Name:        "Name",
				Type:        "Type",
				Placeholder: "Placeholder",
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.strct), func(t *testing.T) {
			got := fields(tc.strct)
			if got != tc.want {
				t.Errorf("fields(): got %v; want %v", got, tc.want)
			}
		})
	}
}
