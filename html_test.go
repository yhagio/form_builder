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

var (
	tplTypeNameValue = template.Must(
		template.
			New("").
			Parse(`<input type="{{.Type}}" name="{{.Name}}" {{with .Value}}value="{{.}}"{{end}}>`),
	)
)

var updateFlag bool

func init() {
	flag.BoolVar(&updateFlag, "update", false, "set the update flag to update the expected output of all golden file")
}

func TestHTML(t *testing.T) {
	tests := map[string]struct {
		tpl     *template.Template
		strct   interface{}
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := form_builder.HTML(tc.tpl, tc.strct)
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
