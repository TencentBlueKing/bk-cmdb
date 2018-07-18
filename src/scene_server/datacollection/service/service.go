package service

import (
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/rdapi"
	"configcenter/src/storage"
)

type Service struct {
	*backbone.Engine
	db    storage.DI
	cache *redis.Client
}

func (s *Service) SetDB(db storage.DI) {
	s.db = db
}

func (s *Service) SetCache(db *redis.Client) {
	s.cache = db
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/collector/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	return ws
}
