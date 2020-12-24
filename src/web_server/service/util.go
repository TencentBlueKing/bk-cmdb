package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
)

func parseModelBizID(data string) (int64, error) {
	model := &struct {
		BizID int64 `json:"bk_biz_id"`
	}{}

	if len(data) != 0 {
		if err := json.Unmarshal([]byte(data), model); nil != err {
			return 0, err
		}
	}

	return model.BizID, nil
}

func parseConfigID(data []byte) (uint64, error) {
	model := &struct {
		ConfigID uint64 `json:"config_id"`
	}{}

	if len(data) != 0 {
		if err := json.Unmarshal(data, model); nil != err {
			return 0, err
		}
	}

	return model.ConfigID, nil
}

func parseOperationChartData(chartType string, instInfo interface{}) (interface{}, error) {
	var chartData interface{}
	// chartData 的数据结构按report_type分三类
	// 1.report_type 属于 ["host_change_biz_chart"]
	// 2.report_type 属于 ["model_inst_change_chart"]
	// 3.report_type 属于 ["custom","model_inst_chart","host_cloud_chart","host_biz_chart","host_os_chart","model_and_inst_count"]
	switch chartType {
	case common.HostChangeBizChart:
		chartData = &map[string][]metadata.StringIDCount{}
	case common.ModelInstChangeChart:
		chartData = &metadata.ModelInstChange{}
	default:
		chartData = &[]metadata.StringIDCount{}
	}
	jsonInstInfo, err := json.Marshal(instInfo)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonInstInfo, chartData)
	if err != nil {
		return nil, err
	}
	return chartData, nil
}

// 依照"bk_obj_id"和"bk_property_type":"objuser"查询"cc_ObjAttDes"集合,得到"bk_property_id"的值;
// 然后以它的值为key,取得Info中的value,然后以value作为param访问ESB,得到其中文名。
func (s *Service) getUsernameMapWithPropertyList(c *gin.Context, objID string, infoList []mapstr.MapStr) (map[string]string, []string, error) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeUser,
		},
	}
	attrRsp, err := s.CoreAPI.CoreService().Model().ReadModelAttr(c, c.Request.Header, objID, &cond)
	if err != nil {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, err
	}
	if !attrRsp.Result {
		blog.Errorf("failed to search the object(%s), err: %s, rid: %s", objID, attrRsp.ErrMsg, rid)
		return nil, nil, err
	}

	usernameList := []string{}
	propertyList := []string{}
	ok := true
	for _, info := range infoList {
		// 主机模型的info内容比inst模型的info内容多封装了一层，需要将内容提取出来。
		if objID == common.BKInnerObjIDHost {
			info, ok = info[common.BKInnerObjIDHost].(map[string]interface{})
			if !ok {
				err = fmt.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, rid: %s", objID, rid)
				blog.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, rid: %s", objID, rid)
				return nil, nil, err
			}
		}
		for _, item := range attrRsp.Data.Info {
			propertyList = append(propertyList, item.PropertyID)
			if info[item.PropertyID] != nil {
				username, ok := info[item.PropertyID].(string)
				if !ok {
					err = fmt.Errorf("failed to cast %s instance info from interface{} to string, rid: %s", objID, rid)
					blog.Errorf("failed to cast %s instance info from interface{} to string, rid: %s", objID, rid)
					return nil, nil, err
				}
				usernameList = append(usernameList, username)
			}
		}
	}
	propertyList = util.RemoveDuplicatesAndEmpty(propertyList)
	userList := util.RemoveDuplicatesAndEmpty(usernameList)
	// get username from esb
	usernameMap, err := s.getUsernameFromEsb(c, userList)
	if err != nil {
		blog.ErrorJSON("get username map from ESB failed, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, err
	}

	return usernameMap, propertyList, nil
}

func (s *Service) getUsernameFromEsb(c *gin.Context, userList []string) (map[string]string, error) {
	var defErr errors.DefaultCCErrorIf
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	userListStr := strings.Join(userList, ",")
	usernameMap := map[string]string{}
	if userList != nil && len(userList) != 0 {
		params := make(map[string]string)
		params["exact_lookups"] = userListStr
		params["fields"] = "username,display_name"
		user := plugins.CurrentPlugin(c, s.Config.Version)

		userListEsb, errNew := user.GetUserList(c, s.Config.ConfigMap)
		if errNew != nil {
			blog.ErrorJSON("get user list from ESB failed, err: %s, rid: %s", errNew.ToCCError(defErr).Error(), rid)
			userListEsb = []*metadata.LoginSystemUserInfo{}
			return nil, errNew.ToCCError(defErr)
		}

		for _, userInfo := range userListEsb {
			username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
			usernameMap[userInfo.EnName] = username
		}
		return usernameMap, nil
	}
	return usernameMap, nil
}
