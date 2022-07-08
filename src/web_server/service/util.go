package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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

// 依照"bk_obj_id"和"bk_property_type":"objuser"查询"cc_ObjAttDes"集合,得到"bk_property_id"的值;
// 然后以它的值为key,取得Info中的value,然后以value作为param访问ESB,得到其中文名。
func (s *Service) getUsernameMapWithPropertyList(c *gin.Context, objID string, infoList []mapstr.MapStr) (
	map[string]string, []string, error) {
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

	usernameList := []string{}
	propertyList := []string{}
	ok := true
	for _, info := range infoList {
		// 主机模型的info内容比inst模型的info内容多封装了一层，需要将内容提取出来。
		if objID == common.BKInnerObjIDHost {
			info, ok = info[common.BKInnerObjIDHost].(map[string]interface{})
			if !ok {
				err = fmt.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, "+
					"rid: %s", objID, rid)
				blog.Errorf("failed to cast %s instance info from interface{} to map[string]interface{}, rid: %s",
					objID, rid)
				return nil, nil, err
			}
		}
		for _, item := range attrRsp.Info {
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
	defErr := s.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(c.Request.Header))
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	usernameMap := map[string]string{}

	if userList != nil && len(userList) != 0 {
		params := make(map[string]string)
		params["fields"] = "username,display_name"
		user := plugins.CurrentPlugin(c, s.Config.LoginVersion)
		// 处理请求的用户数据，每100个用户组成一个字符串进行用户数据的获取
		userListStr := s.getUserListStr(userList)
		var wg sync.WaitGroup
		var lock sync.RWMutex
		var firstErr errors.CCErrorCoder
		pipeline := make(chan bool, 10)
		userListEsb := make([]*metadata.LoginSystemUserInfo, 0)
		for _, subStr := range userListStr {
			pipeline <- true
			wg.Add(1)
			go func(subStr string) {
				defer func() {
					wg.Done()
					<-pipeline
				}()
				params["exact_lookups"] = subStr
				userListEsbSub, errNew := user.GetUserList(c, params)
				if errNew != nil {
					blog.Errorf("get user list from ESB failed, err: %s, rid: %s", errNew.ToCCError(defErr), rid)
					firstErr = errNew.ToCCError(defErr)
					return
				}

				lock.Lock()
				userListEsb = append(userListEsb, userListEsbSub...)
				lock.Unlock()
			}(subStr)
		}
		wg.Wait()

		if firstErr != nil {
			return nil, firstErr
		}

		for _, userInfo := range userListEsb {
			username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
			usernameMap[userInfo.EnName] = username
		}
		return usernameMap, nil
	}
	return usernameMap, nil
}

// getUserListStr get user list str
func (s *Service) getUserListStr(userList []string) []string {
	userListLen := len(userList)
	userListStr := make([]string, 0)
	const getUserMaxCount = 100
	if userListLen <= getUserMaxCount {
		userStr := strings.Join(userList, ",")
		userListStr = append(userListStr, userStr)
		return userListStr
	}

	for i := 0; i < userListLen; i += getUserMaxCount {
		if i+getUserMaxCount < userListLen {
			subUserList := userList[i : i+getUserMaxCount]
			userStr := strings.Join(subUserList, ",")
			userListStr = append(userListStr, userStr)
		} else {
			subUserList := userList[i:userListLen]
			userStr := strings.Join(subUserList, ",")
			userListStr = append(userListStr, userStr)
		}
	}

	return userListStr
}

// getDepartment search department detail and return a id-fullname map
func (s *Service) getDepartment(c *gin.Context, objID string) ([]metadata.DepartmentItem, []string, error) {

	rid := util.GetHTTPCCRequestID(c.Request.Header)
	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeOrganization,
		},
	}
	attrRsp, err := s.CoreAPI.CoreService().Model().ReadModelAttr(c, c.Request.Header, objID, &cond)
	if err != nil {
		blog.Errorf("search object[%s] attribute failed, err: %v, rid: %s", objID, err, rid)
		return nil, nil, err
	}

	if len(attrRsp.Info) == 0 {
		return make([]metadata.DepartmentItem, 0), make([]string, 0), nil
	}

	propertyList := make([]string, 0)
	for _, item := range attrRsp.Info {
		propertyList = append(propertyList, item.PropertyID)
	}

	department, err := s.Logics.GetDepartment(c, s.Config)
	if err != nil {
		blog.Errorf("get department failed, err: %v, rid: %s", err, rid)
		return nil, nil, err
	}

	return department.Results, propertyList, nil
}
