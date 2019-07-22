package form_builder_test

import (
	"form_builder"
	"html/template"
	"testing"
)

var (
	tplTypeNameValue = template.Must(
		template.
			New("").
			Parse(`<input type="{{.Type}}" name="{{.Name}}" {{with .Value}}value="{{.}}"{{end}}>`),
	)
)

func TestHTML(t *testing.T) {
	tests := map[string]struct {
		tpl     *template.Template
		strct   interface{}
		want    template.HTML
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
			want: `<input type="text" name="Name" value="Alice Smith">` +
				`<input type="text" name="Email" value="alice@cc.cc">`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := form_builder.HTML(tc.tpl, tc.strct)
			if err != tc.wantErr {
				t.Fatalf("HTML() err %v; wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("HTML() got %v; want %v", got, tc.want)
			}
		})
	}
}
