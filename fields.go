package form_builder

import "reflect"

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
}

func fields(strct interface{}) []field {
	var formFields []field

	refVal := reflect.ValueOf(strct)
	typ := refVal.Type()

	for i := 0; i < typ.NumField(); i++ {
		typeForm := typ.Field(i)
		refValForm := refVal.Field(i)

		f := field{
			Label:       typeForm.Name,
			Name:        "Name",
			Type:        "Type",
			Placeholder: "Placeholder",
			Value:       refValForm.Interface(),
		}

		formFields = append(formFields, f)
	}

	return formFields
}
