package auth

type ResourceType string

const (
	Business                  ResourceType = "business"
	Object                    ResourceType = "object"
	MainlineObject            ResourceType = "mainlineObject"
	MainlineObjectTopology    ResourceType = "mainlineObjectTopology"
	MainlineInstanceTopology  ResourceType = "mainlineInstanceTopology"
	AssociationKind           ResourceType = "associationKind"
	ObjectAssociation         ResourceType = "objectAssociation"
	ObjectInstanceAssociation ResourceType = "objectInstanceAssociation"
	ObjectInstance            ResourceType = "objectInstance"
)
