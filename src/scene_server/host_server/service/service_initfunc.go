package service

import (
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful/v3"
)

func (s *Service) initService(web *restful.WebService) {

	s.initCloudarea(web)
	s.initFavourite(web)
	s.initFindhost(web)
	s.initHost(web)
	s.initHostapplyrule(web)
	s.initHostlock(web)
	s.initModule(web)
	s.initSpecial(web)
	s.initTransfer(web)
	s.initDynamicGroup(web)
	s.initUsercustom(web)
	s.initAgent(web)
	s.initCloudHost(web)

}

func (s *Service) initCloudarea(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cloudarea", Handler: s.FindManyCloudArea})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/cloudarea", Handler: s.CreatePlatBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/cloudarea", Handler: s.CreatePlat})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/cloudarea/{bk_cloud_id}", Handler: s.UpdatePlat})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/cloudarea/{bk_cloud_id}", Handler: s.DeletePlat})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/updatemany/hosts/cloudarea_field", Handler: s.UpdateHostCloudAreaField})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cloudarea/hostcount", Handler: s.FindCloudAreaHostCount})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initFavourite(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/favorites/search", Handler: s.ListHostFavourites})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/favorites", Handler: s.AddHostFavourite})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/hosts/favorites/{id}", Handler: s.UpdateHostFavouriteByID})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/hosts/favorites/{id}", Handler: s.DeleteHostFavouriteByID})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/hosts/favorites/{id}/incr", Handler: s.IncrHostFavouritesCount})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initFindhost(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/module_relation/bk_biz_id/{bk_biz_id}",
		Handler: s.FindModuleHostRelation})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/by_service_templates/biz/{bk_biz_id}",
		Handler: s.FindHostsByServiceTemplates})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/by_set_templates/biz/{bk_biz_id}",
		Handler: s.FindHostsBySetTemplates})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/list_resource_pool_hosts",
		Handler: s.ListResourcePoolHosts})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/app/{appid}/list_hosts",
		Handler: s.ListBizHosts})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/list_hosts_without_app",
		Handler: s.ListHostsWithNoBiz})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/app/{bk_biz_id}/list_hosts_topo",
		Handler: s.ListBizHostsTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/count_by_topo_node/bk_biz_id/{bk_biz_id}",
		Handler: s.CountTopoNodeHosts})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/by_topo/biz/{bk_biz_id}",
		Handler: s.FindHostsByTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/detail_topo",
		Handler: s.ListHostDetailAndTopology})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/relation/with_topo",
		Handler: s.GetHostRelationsWithMainlineTopoInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/service_template",
		Handler: s.ListServiceTemplateIDsByHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/total_mainline_topo/biz/{bk_biz_id}",
		Handler: s.ListHostTotalMainlineTopo})
	utility.AddToRestfulWebService(web)

}

func (s *Service) initHost(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path: "/hosts/batch", Handler: s.DeleteHostBatchFromResourcePool})
	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path: "/hosts/{bk_supplier_account}/{bk_host_id}", Handler: s.GetHostInstanceProperties})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/add", Handler: s.AddHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/excel/add", Handler: s.AddHostByExcel})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/hosts/add/resource", Handler: s.AddHostToResourcePool})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/add/business_idle",
		Handler: s.AddHostToBusinessIdle})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/search", Handler: s.SearchHost})
	// search host by biz set, **only for ui**
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/hosts/biz_set/{bk_biz_set_id}",
		Handler: s.SearchHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/hosts/search/asstdetail", Handler: s.SearchHostWithAsstDetail})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/hosts/batch", Handler: s.UpdateHostBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPut,
		Path: "/hosts/property/batch", Handler: s.UpdateHostPropertyBatch})
	// TODO: Deprecated, delete this api, used in framework
	// utility.AddHandler(rest.Action{Verb: http.MethodPost,
	// Path: "/hosts/sync/new/host", Handler: s.NewHostSyncAppTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/hosts/modules/idle/set", Handler: s.MoveSetHost2IdleModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/hosts/property/clone", Handler: s.CloneHostProperty})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/hosts/update", Handler: s.UpdateImportHosts})
	// 查询业务下的主机CPU数量的特殊接口，给成本管理使用
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/count/cpu", Handler: s.CountHostCPU})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/kube/search",
		Handler: s.SearchHostWithKube})
	utility.AddToRestfulWebService(web)

}

func (s *Service) initHostapplyrule(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	// 主机属性自动应用
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/create/host_apply_rule/bk_biz_id/{bk_biz_id}", Handler: s.CreateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPut,
		Path: "/update/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}", Handler: s.UpdateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete,
		Path:    "/host/deletemany/module/host_apply_rule/bk_biz_id/{bk_biz_id}",
		Handler: s.DeleteHostApplyRule})

	utility.AddHandler(rest.Action{Verb: http.MethodGet,
		Path: "/find/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/", Handler: s.GetHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}",
		Handler: s.ListHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/createmany/host_apply_rule/bk_biz_id/{bk_biz_id}/batch_create_or_update",
		Handler: s.BatchCreateOrUpdateHostApplyRule})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/createmany/module/host_apply_plan/preview",
		Handler: s.GenerateModuleApplyPlan})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}/host_related_rules",
		Handler: s.ListHostRelatedApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/find/service_template/host_apply_status",
		Handler: s.GetTemplateHostApplyStatus})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/findmany/service_template/host_apply_rule",
		Handler: s.GetServiceTemplateHostApplyRule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/host/createmany/service_template/host_apply_plan/preview", Handler: s.GenerateTemplateApplyPlan})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/host/findmany/module/host_apply_plan/invalid_host_count", Handler: s.GetModuleInvalidHostCount})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/findmany/service_template/host_apply_plan/invalid_host_count",
		Handler: s.GetServiceTemplateInvalidHostCount})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/findmany/service_template/host_apply_rule_count",
		Handler: s.GetServiceTemplateHostApplyRuleCount})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/host/findmany/module/get_module_final_rules", Handler: s.GetModuleFinalRules})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/updatemany/module/host_apply_plan/run",
		Handler: s.UpdateModuleHostApplyRule})

	// task Execute asynchronous template host auto-apply tasks.
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/updatemany/module/host_apply_plan/task",
		Handler: s.ExecModuleHostApplyRule})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/findmany/module/host_apply_plan/status",
		Handler: s.GetHostApplyTaskStatus})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initHostlock(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/lock", Handler: s.LockHost})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/host/lock", Handler: s.UnlockHost})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/lock/search", Handler: s.QueryHostLock})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initModule(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules", Handler: s.TransferHostModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/idle", Handler: s.MoveHost2IdleModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/fault", Handler: s.MoveHost2FaultModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/recycle", Handler: s.MoveHost2RecycleModule})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/resource", Handler: s.MoveHostToResourcePool})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/resource/idle", Handler: s.AssignHostToApp})
	// get host module relation in app
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/read", Handler: s.GetHostModuleRelation})
	// transfer host to other business
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/hosts/modules/across/biz", Handler: s.TransferHostAcrossBusiness})

	// transfer resource host(multi business) to other business.
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/hosts/resource/cross/biz",
		Handler: s.TransferResourceHostsAcrossBusiness})

	// TODO: Deprecated, delete this api. delete host from business, used for framework
	// utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/hosts/module/biz/delete",
	// 	Handler: s.DeleteHostFromBusiness})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/topo/relation/read", Handler: s.GetAppHostTopoRelation})
	// 主机在资源池目录之间转移
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/transfer/resource/directory", Handler: s.TransferHostResourceDirectory})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initSpecial(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/install/bk", Handler: s.BKSystemInstall})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path: "/system/config/user_config/blueking_modify", Handler: s.FindSystemUserConfigBKSwitch})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initTransfer(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/transfer_with_auto_clear_service_instance/bk_biz_id/{bk_biz_id}/",
		Handler: s.TransferHostWithAutoClearServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodPost,
		Path:    "/host/transfer_with_auto_clear_service_instance/bk_biz_id/{bk_biz_id}/preview/",
		Handler: s.TransferHostWithAutoClearServiceInstancePreview})

	utility.AddToRestfulWebService(web)
}

// initDynamicGroup initializes dynamic grouping HTTP handlers.
func (s *Service) initDynamicGroup(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	// create new dynamic group.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodPost,
		Path:    "/dynamicgroup",
		Handler: s.CreateDynamicGroup,
	})

	// update dynamic group.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodPut,
		Path:    "/dynamicgroup/{bk_biz_id}/{id}",
		Handler: s.UpdateDynamicGroup,
	})

	// query target dynamic group.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodGet,
		Path:    "/dynamicgroup/{bk_biz_id}/{id}",
		Handler: s.GetDynamicGroup,
	})

	// delete target dynamic group.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodDelete,
		Path:    "/dynamicgroup/{bk_biz_id}/{id}",
		Handler: s.DeleteDynamicGroup,
	})

	// search(list) dynamic groups.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodPost,
		Path:    "/dynamicgroup/search/{bk_biz_id}",
		Handler: s.SearchDynamicGroup,
	})

	// execute dynamic group and get target resources.
	utility.AddHandler(rest.Action{
		Verb:    http.MethodPost,
		Path:    "/dynamicgroup/data/{bk_biz_id}/{id}",
		Handler: s.ExecuteDynamicGroup,
	})

	utility.AddToRestfulWebService(web)
}

func (s *Service) initUsercustom(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/usercustom", Handler: s.SaveUserCustom})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/usercustom/user/search", Handler: s.GetUserCustom})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/usercustom/default/model", Handler: s.GetModelDefaultCustom})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/usercustom/default/model/{obj_id}", Handler: s.SaveModelDefaultCustom})

	utility.AddToRestfulWebService(web)

}

func (s *Service) initAgent(web *restful.WebService) {

	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/bind/agent", Handler: s.BindAgent})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/host/unbind/agent", Handler: s.UnbindAgent})

	utility.AddToRestfulWebService(web)
}

// initCloudHost init cloud host related api, **dedicated for cloud host management, do not use them for other use**
func (s *Service) initCloudHost(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/cloud_hosts",
		Handler: s.AddCloudHostToBiz})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/cloud_hosts",
		Handler: s.DeleteCloudHostFromBiz})

	utility.AddToRestfulWebService(web)

}
