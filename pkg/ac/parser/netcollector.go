package parser

import (
	meta2 "configcenter/pkg/ac/meta"
	"net/http"
	"regexp"
)

func (ps *parseStream) netCollectorRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.netCollector().
		netDevice().
		netProperty().
		netReport()

	return ps
}

const (
	findNetCollectorsPattern  = "/api/v3/collector/netcollect/collector/action/search"
	updateNetCollectorPattern = "/api/v3/collector/netcollect/collector/action/update"
	startNetCollectorPattern  = "/api/v3/collector/netcollect/collector/action/discover"
)

func (ps *parseStream) netCollector() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find all the business's net collectors
	if ps.hitPattern(findNetCollectorsPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetCollector,
					Action: meta2.FindMany,
				},
			},
		}
		return ps
	}

	// update net collector in a business.
	if ps.hitPattern(updateNetCollectorPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetCollector,
					Action: meta2.UpdateMany,
				},
			},
		}
		return ps
	}

	// start one/many net collector to collector data.
	if ps.hitPattern(startNetCollectorPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetCollector,
					Action: meta2.UpdateMany,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createNetDevicePattern         = "/api/v3/collector/netcollect/device/action/create"
	updateOrCreateNetDevicePattern = "/api/v3/collector/netcollect/device/action/batch"
	findNetDevicePattern           = "/api/v3/collector/netcollect/device/action/search"
	deleteNetDevicePattern         = "/api/v3/collector/netcollect/device/action/delete"
)

var (
	updateNetDeviceRegexp = regexp.MustCompile(`/api/v3/collector/netcollect/device/[0-9]+/action/update`)
)

func (ps *parseStream) netDevice() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create a net device
	if ps.hitPattern(createNetDevicePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetDevice,
					Action: meta2.Create,
				},
			},
		}
		return ps
	}

	// update a device
	if ps.hitRegexp(updateNetDeviceRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetDevice,
					Action: meta2.Update,
				},
			},
		}
		return ps
	}

	// update or create new net device in batch.
	if ps.hitPattern(updateOrCreateNetDevicePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetDevice,
					Action: meta2.UpdateMany,
				},
			},
		}
		return ps
	}

	// find net devices
	if ps.hitPattern(findNetDevicePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetDevice,
					Action: meta2.FindMany,
				},
			},
		}
		return ps
	}

	// delete net device patch
	if ps.hitPattern(deleteNetDevicePattern, http.MethodDelete) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetDevice,
					Action: meta2.DeleteMany,
				},
			},
		}
		return ps
	}

	// TODO: add net device import, export auth filter.
	// add import template auth filter.

	return ps
}

const (
	createNetCollectorPropertyPattern           = "/api/v3/collector/netcollect/property/action/create"
	updateOrCreateNetCollectorPropertiesPattern = "/api/v3/collector/netcollect/property/action/batch"
	findNetCollectorPropertiesPattern           = "/api/v3/collector/netcollect/property/action/search"
	deleteNetCollectorPropertiesPattern         = "/api/v3/collector/netcollect/property/action/delete"
)

var (
	updateNetCollectorPropertyRegexp = regexp.MustCompile(`/api/v3/collector/netcollect/property/[0-9]+/action/update`)
)

func (ps *parseStream) netProperty() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create property for a net collector
	if ps.hitPattern(createNetCollectorPropertyPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetProperty,
					Action: meta2.Create,
				},
			},
		}
		return ps
	}

	// update property for a net collector.
	if ps.hitRegexp(updateNetCollectorPropertyRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetProperty,
					Action: meta2.Update,
				},
			},
		}
		return ps
	}

	// update or create net collector properties.
	if ps.hitPattern(updateOrCreateNetCollectorPropertiesPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetProperty,
					Action: meta2.UpdateMany,
				},
			},
		}
		return ps
	}

	// find net collector properties
	if ps.hitPattern(findNetCollectorPropertiesPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetProperty,
					Action: meta2.Find,
				},
			},
		}
		return ps
	}

	// delete net collector properties batch
	if ps.hitPattern(deleteNetCollectorPropertiesPattern, http.MethodDelete) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetProperty,
					Action: meta2.DeleteMany,
				},
			},
		}
		return ps
	}

	// TODO: add import and export net device properties auth filter.
	// add import net collector properties template auth filter.

	return ps
}

const (
	findNetDeviceSimpleReportPattern        = "/api/v3/collector/netcollect/summary/action/search"
	findNetDeviceDetailReportPattern        = "/api/v3/collector/netcollect/report/action/search"
	findNetDeviceReportConfirmPattern       = "/api/v3/collector/netcollect/report/action/search"
	findNetDeviceReportConfirmDetailPattern = "/api/v3/collector/netcollect/report/action/confirm"
)

func (ps *parseStream) netReport() *parseStream {
	if ps.shouldReturn() {
		return ps
	}
	// find net device simple report
	if ps.hitPattern(findNetDeviceSimpleReportPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetReport,
					Action: meta2.Find,
				},
			},
		}
		return ps
	}

	// find net device detailed report
	if ps.hitPattern(findNetDeviceDetailReportPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetReport,
					Action: meta2.Find,
				},
			},
		}
		return ps
	}

	// find net device report confirm history.
	if ps.hitPattern(findNetDeviceReportConfirmPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetReport,
					Action: meta2.Find,
				},
			},
		}
		return ps
	}

	// find net device detailed report confirm history.
	if ps.hitPattern(findNetDeviceReportConfirmDetailPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta2.ResourceAttribute{
			meta2.ResourceAttribute{
				Basic: meta2.Basic{
					Type:   meta2.NetDataCollector,
					Name:   meta2.NetReport,
					Action: meta2.Find,
				},
			},
		}
		return ps
	}

	return ps
}
