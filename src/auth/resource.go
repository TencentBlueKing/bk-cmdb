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

	HostUserCustom        ResourceType = "hostUserCustom"
	HostFavorite          ResourceType = "hostFavorite"
	Host                  ResourceType = "host"
	AddHostToResourcePool ResourceType = "addHostToResourcePool"
	MoveHostToModule      ResourceType = "moveHostToModule"
	// move resource pool hosts to a business idle module
	MoveResPoolHostToBizIdleModule ResourceType = "moveResPoolHostToBizIdleModule"
	MoveHostToBizFaultModule       ResourceType = "moveHostToBizFaultModule"
	MoveHostToBizIdleModule        ResourceType = "moveHostToBizIdleModule"
	MoveHostFromModuleToResPool    ResourceType = "moveHostFromModuleToResPool"
	MoveHostToAnotherBizModule     ResourceType = "moveHostToAnotherBizModule"
	CleanHostInSetOrModule         ResourceType = "cleanHostInSetOrModule"
	MoveHostsToOrBusinessModule    ResourceType = "moveHostsToBusinessOrModule"

	Process                      ResourceType = "process"
	ProcessConfigTemplate        ResourceType = "processConfigTemplate"
	ProcessConfigTemplateVersion ResourceType = "processConfigTemplateVersion"
	ProcessBoundConfig           ResourceType = "processBoundConfig"

	NetCollector ResourceType = "netCollector"
	NetDevice    ResourceType = "netDevice"
	NetProperty  ResourceType = "netProperty"
	NetReport    ResourceType = "netReport"
)
