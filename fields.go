package form_builder

import "reflect"

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
}

func fields(strct interface{}) field {
	rv := reflect.ValueOf(strct)
	t := rv.Type()

	tf := t.Field(0)

	return field{
		Label:       tf.Name,
		Name:        "Name",
		Type:        "Type",
		Placeholder: "Placeholder",
	}
}
