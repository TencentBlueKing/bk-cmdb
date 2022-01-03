package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/scene_server/event_server/sync/hostIdentifier"
)

// SyncHostIdentifier sync host identifier, add hostInfo message to redis fail host list
func (s *Service) SyncHostIdentifier(ctx *rest.Contexts) {
	hostInfo := new(hostIdentifier.HostInfo)
	if err := ctx.DecodeInto(&hostInfo); err != nil {
		blog.Errorf("decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if hostInfo.HostInnerIP == "" {
		blog.Errorf("not set param bk_host_innerip, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_host_innerip"))
		return
	}
	if err := s.cache.LPush(ctx.Kit.Ctx, hostIdentifier.RedisFailHostListName, hostInfo).Err(); err != nil {
		blog.Errorf("add hostInfo to redis list error, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
