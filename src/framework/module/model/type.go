package model

// GroupIterator the group iterator
type GroupIterator interface {
	Next() (Group, error)
}

// Group the interface declaration for model maintence
type Group interface {
	SetID(id string)
	GetID() string
	SetName(name string)
	SetIndex(idx int)
	GetIndex() int
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDefault()
	SetNonDefault()
	Default() bool

	CreateAttribute() Attribute
	FindAttributes(attributeName string) (AttributeIterator, error)
}

// ClassificationIterator the classification iterator
type ClassificationIterator interface {
	Next() (Classification, error)
}

// Classification the interface declaration for model classification
type Classification interface {
	SetID(id string)
	SetName(name string)
	SetIcon(iconName string)
	GetID() string

	CreateModel() Model
	FindModels(modelName string) (Iterator, error)
}

// Iterator the model iterator
type Iterator interface {
	Next() (Model, error)
}

// Model the interface declaration for model maintence
type Model interface {
	SetClassification(class Classification)
	SetIcon(iconName string)
	SetID(id string)
	SetName(name string)
	SetPaused(isPaused bool)
	SetPosition(position string)
	SetSupplierAccount(ownerID string)
	SetDescription(desc string)
	SetCreator(creator string)
	SetModifier(modifier string)

	CreateAttribute() Attribute
	CreateGroup() Group
	FindAttributes(attributeName string) (AttributeIterator, error)
	FindGroups(groupName string) (GroupIterator, error)

	CreateInst() Inst
}

// AttributeIterator the attribute iterator
type AttributeIterator interface {
	Next() (Attribute, error)
}

// Attribute the interface declaration for model attribute maintence
type Attribute interface {
	SetID(id string)
	SetName(name string)
	SetUnit(unit string)
	SetPlaceholer(placeHoler string)
	SetEditable()
	SetNonEditable()
	Editable() bool
	SetRequired()
	SetNonRequired()
	Required() bool
	SetKey(isKey bool)
	SetOption(option string)
	SetDescrition(des string)
}

// Inst the instance operator interface
type Inst interface {
	SetValue(key string, value interface{})
	Save() error
}
