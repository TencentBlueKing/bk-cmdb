package service

import (
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

func (s *Service) InitFunc() {
	header := make(http.Header, 0)
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}

	srvData := s.newSrvComm(header)
	go srvData.lgc.TimerDeleteHistoryTask(srvData.ctx)
}
