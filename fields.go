package form_builder

import "reflect"

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
}

func fields(strct interface{}) []field {
	var formFields []field

	rv := reflect.ValueOf(strct)
	t := rv.Type()

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		f := field{
			Label:       tf.Name,
			Name:        "Name",
			Type:        "Type",
			Placeholder: "Placeholder",
		}
		formFields = append(formFields, f)
	}

	return formFields
}
