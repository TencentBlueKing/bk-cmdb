package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"fmt"
)

type hostSnapAuditLog struct {
	audit
}

// GenerateAuditLog 生成hostsnap审计日志(属于更新操作),且curData和updateFields都不能为空
func (h *hostSnapAuditLog) GenerateAuditLog(kit *rest.Kit, hostID int64, innerIP string, preData,
	updateFields map[string]interface{}) (*metadata.AuditLog, error) {

	// 获得主机所属业务的业务ID
	input := &metadata.HostModuleRelationRequest{HostIDArr: []int64{hostID}, Fields: []string{common.BKAppIDField}}
	moduleHost, err := h.clientSet.Host().GetHostModuleRelation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("snapshot get host: %d/%s module relation failed, err:%v, rid: %s", hostID, innerIP, err, kit.Rid)
		return nil, err
	}
	if !moduleHost.Result {
		blog.Errorf("snapshot get host: %d/%s module relation failed, err: %v, rid: %s", hostID, innerIP, moduleHost.ErrMsg, kit.Rid)
		return nil, fmt.Errorf("snapshot get moduleHostConfig failed, fail to create auditLog")
	}

	var bizID int64
	if len(moduleHost.Data.Info) > 0 {
		bizID = moduleHost.Data.Info[0].AppID
	}

	// 生成审计
	auditLog := metadata.AuditLog{
		AuditType:    metadata.HostType,
		ResourceType: metadata.HostRes,
		Action:       metadata.AuditUpdate,
		OperateFrom:  metadata.FromDataCollection,
		BusinessID:   bizID,
		ResourceID:   hostID,
		ResourceName: innerIP,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				Details: &metadata.BasicContent{
					PreData:      preData,
					UpdateFields: updateFields,
				},
			},
			ModelID: common.BKInnerObjIDHost,
		},
	}
	return &auditLog, nil
}

func NewHostSnapAudit(clientSet coreservice.CoreServiceClientInterface) *hostSnapAuditLog {
	return &hostSnapAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
