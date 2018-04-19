package model

// Group the interface declaration for model maintence
type Group interface {
}

// Classification the interface declaration for model classification
type Classification interface{

}

// Model the interface declaration for model maintence
type Model interface {
	CreateAttribute() Attribute
}

// Attribute the interface declaration for model attribute maintence
type Attribute interface {
	SetID(id string)
	SetName(name string)
	SetUnit(unit string)
	SetPlaceholer(placeHoler string)
	SetEditable(editable bool)
	SetRequired(required bool)
	SetKey(isKey bool)
	SetOption(option string)
	SetDescrition(des string)
}
