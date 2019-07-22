package form_builder

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFields(t *testing.T) {
	var nilStructPointer *struct {
		Name string
		Age  int
	}

	type testStructs struct {
		strct interface{}
		want  []field
	}

	tests := map[string]testStructs{
		"Field name is determined from provided struct 1": {
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
		"Field name is determined from provided struct 2": {
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
		"Multiple field names are determined from struct": {
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
		"Values should be parsed as well as field names": {
			strct: struct {
				Name  string
				Email string
				Age   int
			}{
				Name:  "Alice Smith",
				Email: "alice@cc.cc",
				Age:   25,
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "Alice Smith",
				},
				{
					Label:       "Email",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "alice@cc.cc",
				},
				{
					Label:       "Age",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       25,
				},
			},
		},
		"Unexported fields should be skipped": {
			strct: struct {
				Name  string
				email string
				Age   int
			}{
				Name:  "Alice Smith",
				email: "alice@cc.cc",
				Age:   25,
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "Alice Smith",
				},
				{
					Label:       "Age",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       25,
				},
			},
		},
		"Pointers to structs should be supported": {
			strct: &struct {
				Name string
				Age  int
			}{
				Name: "Alice Smith",
				Age:  25,
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "Alice Smith",
				},
				{
					Label:       "Age",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       25,
				},
			},
		},
		"Nil pointers with a struct type should be supported": {
			strct: nilStructPointer,
			want: []field{
				{
					Label:       "Name",
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
		"Pointer fields should be supported": {
			strct: struct {
				Name *string
				Age  *int
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
					Label:       "Age",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       0,
				},
			},
		},
		"Nested structs should be supported": {
			strct: struct {
				Name    string
				Address struct {
					Street string
					Zip    int
				}
			}{
				Name: "Alice Smith",
				Address: struct {
					Street string
					Zip    int
				}{
					Street: "123 ABC St",
					Zip:    12345,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "Alice Smith",
				},
				{
					Label:       "Street",
					Name:        "Address.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "123 ABC St",
				},
				{
					Label:       "Zip",
					Name:        "Address.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       12345,
				},
			},
		},
		"Doubly nested structs should be supported": {
			strct: struct {
				A struct {
					B struct {
						C1 string
						C2 int
					}
				}
			}{
				A: struct {
					B struct {
						C1 string
						C2 int
					}
				}{
					B: struct {
						C1 string
						C2 int
					}{
						C1: "C1-value",
						C2: 123,
					},
				},
			},
			want: []field{
				{
					Label:       "C1",
					Name:        "A.B.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "C1-value",
				},
				{
					Label:       "C2",
					Name:        "A.B.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       123,
				},
			},
		},
		"Nested pointer structs should be supported": {
			strct: struct {
				Name    string
				Address *struct {
					Street string
					Zip    int
				}
				ContactCard *struct {
					Phone string
				}
			}{
				Name: "Alice Smith",
				Address: &struct {
					Street string
					Zip    int
				}{
					Street: "123 ABC St",
					Zip:    12345,
				},
			},
			want: []field{
				{
					Label:       "Name",
					Name:        "Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "Alice Smith",
				},
				{
					Label:       "Street",
					Name:        "Address.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "123 ABC St",
				},
				{
					Label:       "Zip",
					Name:        "Address.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       12345,
				},
				{
					Label:       "Phone",
					Name:        "ContactCard.Name",
					Type:        "Type",
					Placeholder: "Placeholder",
					Value:       "",
				},
			},
		},
	}

	for key, tc := range tests {
		t.Run(fmt.Sprintf("%v", key), func(t *testing.T) {
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

func TestFields_invalidTypes(t *testing.T) {
	tests := []struct {
		notAStruct interface{}
	}{
		{"string"},
		{12345},
		{nil},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.notAStruct), func(t *testing.T) {

			defer func() {
				if err := recover(); err == nil {
					t.Errorf("fields(%v) did not panic", tc.notAStruct)
				}
			}()

			fields(tc.notAStruct)
		})
	}
}
