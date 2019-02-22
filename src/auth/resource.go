package auth

type ResourceType string

func (r ResourceType) String() string {
	return string(r)
}

const (
	Business                  ResourceType = "business"
	Object                    ResourceType = "object"
	ObjectModule              ResourceType = "objectModule"
	ObjectSet                 ResourceType = "objectSet"
	MainlineObject            ResourceType = "mainlineObject"
	MainlineObjectTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology  ResourceType = "mainlineInstanceTopology"
	AssociationType           ResourceType = "associationType"
	ObjectAssociation         ResourceType = "objectAssociation"
	ObjectInstanceAssociation ResourceType = "objectInstanceAssociation"
	ObjectInstance            ResourceType = "objectInstance"
	ObjectInstanceTopology    ResourceType = "objectInstanceTopology"
	ObjectTopology            ResourceType = "objectTopology"
	ObjectClassification      ResourceType = "objectClassification"
	ObjectAttributeGroup      ResourceType = "objectAttributeGroup"
	ObjectAttribute           ResourceType = "objectAttribute"
	ObjectUnique              ResourceType = "objectUnique"
)
