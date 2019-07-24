package form_builder

import (
	"reflect"
	"strings"
)

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
	Value       interface{}
}

func (f *field) apply(tags map[string]string) {
	if v, ok := tags["label"]; ok {
		f.Label = v
	}
	if v, ok := tags["name"]; ok {
		f.Name = v
	}
	if v, ok := tags["placeholder"]; ok {
		f.Placeholder = v
	}
	if v, ok := tags["type"]; ok {
		f.Type = v
	}
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

func fields(strct interface{}, parentNames ...string) []field {
	refVal := valueOf(strct)

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
			nestedParentNames := append(parentNames, typeForm.Name)
			nestedFields := fields(refValForm.Interface(), nestedParentNames...)
			formFields = append(formFields, nestedFields...)
			continue
		}

		names := append(parentNames, typeForm.Name)
		name := strings.Join(names, ".")

		f := field{
			Label:       typeForm.Name,
			Name:        name,
			Type:        "text",
			Placeholder: typeForm.Name,
			Value:       refValForm.Interface(),
		}

		f.apply(parseTags(typeForm))

		formFields = append(formFields, f)
	}

	return formFields
}

func parseTags(rsf reflect.StructField) map[string]string {
	rawTag := rsf.Tag.Get("form")
	if len(rawTag) == 0 {
		return nil
	}

	result := make(map[string]string)

	tags := strings.Split(rawTag, ";")
	for _, tag := range tags {
		kv := strings.Split(tag, "=")
		if len(kv) != 2 {
			panic("form: invalid struct tag")
		}

		k, v := kv[0], kv[1]
		result[k] = v
	}

	return result
}
