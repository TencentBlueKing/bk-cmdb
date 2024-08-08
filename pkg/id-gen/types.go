package idgen

import "configcenter/src/common"

// IDGenType is the id generator type
type IDGenType string

const (
	// Biz is the business id generator type
	Biz IDGenType = "biz"
	// Set is the set id generator type
	Set IDGenType = "set"
	// Module is the module id generator type
	Module IDGenType = "module"
	// Host is the host id generator type
	Host IDGenType = "host"
	// ObjectInstance is the object instance id generator type
	ObjectInstance IDGenType = "object_instance"
	// InstAsst is the instance association id generator type
	InstAsst IDGenType = "inst_asst"
	// SetTemplate is the set template id generator type
	SetTemplate IDGenType = "set_template"
	// ServiceTemplate is the service template id generator type
	ServiceTemplate IDGenType = "service_template"
	// ProcessTemplate is the process template id generator type
	ProcessTemplate IDGenType = "process_template"
	// ServiceInstance is the service instance id generator type
	ServiceInstance IDGenType = "service_instance"
	// Process is the process id generator type
	Process IDGenType = "process"
)

// GetIDGenSequenceName get id generator sequence name by id generator type
func GetIDGenSequenceName(typ IDGenType) (string, bool) {
	sequenceName, exists := idGenTypeSeqNameMap[typ]
	return sequenceName, exists
}

var idGenTypeSeqNameMap = map[IDGenType]string{
	Biz:             common.BKTableNameBaseApp,
	Set:             common.BKTableNameBaseSet,
	Module:          common.BKTableNameBaseModule,
	Host:            common.BKTableNameBaseHost,
	ObjectInstance:  common.BKTableNameBaseInst,
	InstAsst:        common.BKTableNameInstAsst,
	SetTemplate:     common.BKTableNameSetTemplate,
	ServiceTemplate: common.BKTableNameServiceTemplate,
	ProcessTemplate: common.BKTableNameProcessTemplate,
	ServiceInstance: common.BKTableNameServiceInstance,
	Process:         common.BKTableNameBaseProcess,
}
