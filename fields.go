package form_builder

type field struct {
	Label       string
	Name        string
	Type        string
	Placeholder string
}

func fields(strct interface{}) field {
	return field{
		Label:       "Label",
		Name:        "Name",
		Type:        "Type",
		Placeholder: "Placeholder",
	}
}
