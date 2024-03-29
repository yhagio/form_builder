package form_builder_test

import (
	"flag"
	"fmt"
	"form_builder"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var updateFlag bool

func init() {
	flag.BoolVar(&updateFlag, "update", false, "set the update flag to update the expected output of all golden file")
}

var (
	tplTypeNameValue = template.Must(template.New("").Parse(`<input type="{{.Type}}" name="{{.Name}}"{{with .Value}} value="{{.}}"{{end}}>`))
	tplAll           = template.Must(template.New("").Parse(`
	<label>{{.Label}}</label>
	<input
		type="{{.Type}}"
		name="{{.Name}}"
		placeholder="{{.Placeholder}}"
		{{with .Value}}value="{{.}}"{{end}}>`))
	tplErrors = template.Must(template.New("").Parse(`
		<label>{{.Label}}</label>
		<input
			class="{{with .Errors}}border-red{{end}}"
			type="{{.Type}}"
			name="{{.Name}}"
			placeholder="{{.Placeholder}}"
			{{with .Value}}value="{{.}}"{{end}}>
		{{range .Errors}}
			<p class="text-red text-xs italic">{{.}}</p>
		{{end}}`))
)

func TestHTML(t *testing.T) {
	tests := map[string]struct {
		tpl     *template.Template
		strct   interface{}
		errors  []form_builder.FieldError
		want    string
		wantErr error
	}{
		"Simple form with values": {
			tpl: tplTypeNameValue,
			strct: struct {
				Name  string
				Email string
			}{
				Name:  "Alice Smith",
				Email: "alice@cc.cc",
			},
			want: "TestHTML_basic.golden",
		},
		"A form with custom struct tags": {
			tpl: tplAll,
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
			want: "TestHTML_structTags.golden",
		},
		"A form with errors": {
			tpl: tplErrors,
			strct: struct {
				Email    string `form:"label=Email Address;placeholder=you@domain.com;type=email;name=EmailAddress"`
				Password string `form:"type=password"`
			}{
				Email:    "email@taken.com",
				Password: "badpw",
			},
			errors: []form_builder.FieldError{
				{Field: "EmailAddress", Error: "Email address is already taken"},
				{Field: "Password", Error: "Password must be between 201 and 210 characters"},
				{Field: "Password", Error: "Password must contain a greek letter"},
				{Field: "Password", Error: "Password must be a palindrome"},
				{Field: "Password", Error: "Password must contain an emoji"},
			},
			want: "TestHTML_errors.golden",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := form_builder.HTML(tc.tpl, tc.strct, tc.errors...)
			if err != tc.wantErr {
				t.Fatalf("HTML() err %v; wantErr %v", err, tc.wantErr)
			}

			gotFilename := strings.Replace(tc.want, ".golden", ".got", 1)
			os.Remove(gotFilename)

			if updateFlag {
				writeFile(t, tc.want, string(got))
				t.Logf("Updated golden file %s", tc.want)
			}

			want := template.HTML(readFile(t, tc.want))
			if got != want {
				t.Errorf("HTML() - results do not match golden file.")
				writeFile(t, gotFilename, string(got))
				t.Errorf(" To compare run: diff %s %s", gotFilename, tc.want)
			}
		})
	}
}

func writeFile(t *testing.T, filename, contents string) {
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Error creating file %v: %v", file, err)
	}
	defer file.Close()
	fmt.Fprint(file, contents)
}

func readFile(t *testing.T, filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Error opening file %v: %v", file, err)
	}
	defer file.Close()

	byte, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("Error reading file %v: %v", file, err)
	}

	return byte
}
