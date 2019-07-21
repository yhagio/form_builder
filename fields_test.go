package form_builder

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {

	type testStructs struct {
		strct interface{}
		want  []field
	}

	tests := []testStructs{
		{
			strct: struct {
				Name string
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "",
				},
			},
		},
		{
			strct: struct {
				FullName string
			}{},
			want: []field{
				{
					Label:       "FullName",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "",
				},
			},
		},
		{
			strct: struct {
				Name  string
				Email string
				Age   int
			}{},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "",
				},
				{
					Label:       "Email",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       0,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.strct), func(t *testing.T) {
			got := fields(tc.strct)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("fields(): got %v; want %v", got, tc.want)
			}
		})
	}
}

func TestFields_labels(t *testing.T) {

	hasLabels := func(labels ...string) func(*testing.T, []field) {
		return func(t *testing.T, fields []field) {
			if len(fields) != len(labels) {
				t.Errorf("fields() len: got %d; want %d", len(fields), len(labels))
			}

			for i := 0; i < len(fields); i++ {
				if fields[i].Label != labels[i] {
					t.Errorf("fields()[%d].Label: got %s; want %s", i, fields[i].Label, labels[i])
				}
			}
		}
	}

	hasValues := func(values ...interface{}) func(*testing.T, []field) {
		return func(t *testing.T, fields []field) {
			if len(fields) != len(values) {
				t.Errorf("fields() len: got %d; want %d", len(fields), len(values))
			}

			for i := 0; i < len(fields); i++ {
				if fields[i].Value != values[i] {
					t.Errorf("fields()[%d].Value: got %s; want %s", i, fields[i].Value, values[i])
				}
			}
		}
	}

	check := func(checks ...func(*testing.T, []field)) []func(*testing.T, []field) {
		return checks
	}

	type testStructs struct {
		strct  interface{}
		checks []func(*testing.T, []field)
	}

	tests := map[string]testStructs{
		"No Value": {
			strct: struct {
				Name string
			}{},
			checks: check(hasLabels("Name")),
		},
		"Multiple values": {
			strct: struct {
				Name  string
				Email string
				Age   int
			}{
				Name:  "Alice Smith",
				Email: "alice@cc.cc",
				Age:   25,
			},
			checks: check(
				hasLabels("Name", "Email", "Age"),
				hasValues("Alice Smith", "alice@cc.cc", 25),
			),
		},
	}

	for key, tc := range tests {
		t.Run(fmt.Sprintf("%v", key), func(t *testing.T) {
			got := fields(tc.strct)

			for _, check := range tc.checks {
				check(t, got)
			}
		})
	}
}
