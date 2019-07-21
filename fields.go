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
	refVal := reflect.ValueOf(strct)
	// If the value is pointer, set the value of whatever the value points to
	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}
	// Make sure the value is struct
	if refVal.Kind() != reflect.Struct {
		panic("Oh oh. Only Struct is supported!")
	}

	typ := refVal.Type()

	var formFields []field
	for i := 0; i < typ.NumField(); i++ {
		typeForm := typ.Field(i)
		refValForm := refVal.Field(i)

		if !refValForm.CanInterface() {
			continue
		}

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
