package service

import (
	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/service/v2"
	"configcenter/src/common/backbone"
)

// ServiceInterface service interface
type Service interface {
	WebService() *restful.WebService
	SetEngine(*backbone.Engine)
}

func NewService() Service {
	return &v2.Service{}
}
