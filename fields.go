package form_builder

import "reflect"

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
}

func valueOf(val interface{}) reflect.Value {
	var refVal reflect.Value
	switch value := val.(type) {
	case reflect.Value:
		refVal = value
	default:
		refVal = reflect.ValueOf(val)
	}

	// With any pointers we really want to just work with their underlying
	// type.
	if refVal.Kind() == reflect.Ptr {
		// The underlying type is pretty useless if it is nil, so we need to
		// instantiate a new copy of whatever that is before using it.
		if refVal.IsNil() {
			refVal = reflect.New(refVal.Type().Elem())
		}
		refVal = refVal.Elem()
	}
	return refVal
}

func fields(strct interface{}) []field {
	refVal := valueOf(strct)
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
		refValForm := valueOf(refVal.Field(i))

		// Check unexported field
		if !refValForm.CanInterface() {
			continue
		}

		// Supports nested fields
		if refValForm.Kind() == reflect.Struct {
			nestedFields := fields(refValForm.Interface())
			for i, nestedField := range nestedFields {
				nestedFields[i].Name = typeForm.Name + "." + nestedField.Name
			}
			formFields = append(formFields, nestedFields...)
			continue
		}

		f := field{
			Label:       typeForm.Name,
			Name:        typeForm.Name,
			Type:        "text",
			Placeholder: typeForm.Name,
			Value:       refValForm.Interface(),
		}

		formFields = append(formFields, f)
	}

	return formFields
}
