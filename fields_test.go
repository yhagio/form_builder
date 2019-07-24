package form_builder

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	tests := map[string]struct {
		arg  reflect.StructField
		want map[string]string
	}{
		"empty tag": {
			arg:  reflect.StructField{},
			want: nil,
		},
		"label tag": {
			arg: reflect.StructField{
				Tag: `form:"label=Full Name"`,
			},
			want: map[string]string{
				"label": "Full Name",
			},
		},
		"multiple tags": {
			arg: reflect.StructField{
				Tag: `form:"label=Full Name;email=Email"`,
			},
			want: map[string]string{
				"label": "Full Name",
				"email": "Email",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseTags(tc.arg)
			if len(got) != len(tc.want) {
				t.Errorf("parseTags() len = %d, want %d", len(got), len(tc.want))
			}

			for k, v := range tc.want {
				gotVal, ok := got[k]
				if !ok {
					t.Errorf("parseTags() missing key %q", k)
					continue
				}
				if gotVal != v {
					t.Errorf("parseTags()[%q] = %q; want %q", k, gotVal, v)
				}
				delete(got, k)
			}

			for gotKey, gotVal := range got {
				t.Errorf("parseTags() extra key %q, value = %q", gotKey, gotVal)
			}
		})
	}
}

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
					Type:        "text",
					Placeholder: "Name",
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
					Name:        "FullName",
					Type:        "text",
					Placeholder: "FullName",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "Alice Smith",
				},
				{
					Label:       "Email",
					Name:        "Email",
					Type:        "text",
					Placeholder: "Email",
					Value:       "alice@cc.cc",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "Alice Smith",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "Alice Smith",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "",
				},
				{
					Label:       "Age",
					Name:        "Age",
					Type:        "text",
					Placeholder: "Age",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "Alice Smith",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "123 ABC St",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
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
					Name:        "A.B.C1",
					Type:        "text",
					Placeholder: "C1",
					Value:       "C1-value",
				},
				{
					Label:       "C2",
					Name:        "A.B.C2",
					Type:        "text",
					Placeholder: "C2",
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
					Type:        "text",
					Placeholder: "Name",
					Value:       "Alice Smith",
				},
				{
					Label:       "Street",
					Name:        "Address.Street",
					Type:        "text",
					Placeholder: "Street",
					Value:       "123 ABC St",
				},
				{
					Label:       "Zip",
					Name:        "Address.Zip",
					Type:        "text",
					Placeholder: "Zip",
					Value:       12345,
				},
				{
					Label:       "Phone",
					Name:        "ContactCard.Phone",
					Type:        "text",
					Placeholder: "Phone",
					Value:       "",
				},
			},
		},
		"Struct tags": {
			strct: struct {
				LabelTest       string `form:"label=This is custom"`
				NameTest        string `form:"name=full_name"`
				TypeTest        int    `form:"type=number"`
				PlaceholderTest string `form:"placeholder=your value goes here..."`
				Nested          struct {
					MultiTest string `form:"name=NestedMulti;label=This is nested;type=email;placeholder=user@example.com"`
				}
			}{
				PlaceholderTest: "value and placeholder",
			},
			want: []field{
				{
					Label:       "This is custom",
					Name:        "LabelTest",
					Type:        "text",
					Placeholder: "LabelTest",
					Value:       "",
				},
				{
					Label:       "NameTest",
					Name:        "full_name",
					Type:        "text",
					Placeholder: "NameTest",
					Value:       "",
				},
				{
					Label:       "TypeTest",
					Name:        "TypeTest",
					Type:        "number",
					Placeholder: "TypeTest",
					Value:       0,
				},
				{
					Label:       "PlaceholderTest",
					Name:        "PlaceholderTest",
					Type:        "text",
					Placeholder: "your value goes here...",
					Value:       "value and placeholder",
				},
				{
					Label:       "This is nested",
					Name:        "NestedMulti",
					Type:        "email",
					Placeholder: "user@example.com",
					Value:       "",
				},
			},
		},
	}

	for key, tc := range tests {
		t.Run(fmt.Sprintf("%v", key), func(t *testing.T) {
			got := fields(tc.strct)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("fields():\n  got %v;\n want %v", got, tc.want)
			}

			if len(got) != len(tc.want) {
				t.Errorf("fields(): got %d; want %d", len(got), len(tc.want))
			}

			for i, gotField := range got {
				if i >= len(tc.want) {
					break
				}
				wantField := tc.want[i]
				if reflect.DeepEqual(gotField, wantField) {
					continue
				}

				if gotField.Label != wantField.Label {
					t.Errorf(" .Label = %v; want %v", gotField.Label, wantField.Label)
				}

				if gotField.Name != wantField.Name {
					t.Errorf(" .Name = %v; want %v", gotField.Name, wantField.Name)
				}

				if gotField.Type != wantField.Type {
					t.Errorf(" .Type = %v; want %v", gotField.Type, wantField.Type)
				}

				if gotField.Placeholder != wantField.Placeholder {
					t.Errorf(" .Placeholder = %v; want %v", gotField.Placeholder, wantField.Placeholder)
				}

				if gotField.Value != wantField.Value {
					t.Errorf(" .Value = %v; want %v", gotField.Value, wantField.Value)
				}
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

func TestParseTags_invalidStructTypes(t *testing.T) {
	tests := []struct {
		arg reflect.StructField
	}{
		{reflect.StructField{Tag: `form:"invalid"`}},
	}

	for _, tc := range tests {
		t.Run(string(tc.arg.Tag), func(t *testing.T) {

			defer func() {
				if err := recover(); err == nil {
					t.Errorf("parseTags() did not panic")
				}
			}()

			parseTags(tc.arg)
		})
	}
}
