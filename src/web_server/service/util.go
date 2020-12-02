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
	webCommon "configcenter/src/web_server/common"
	"configcenter/src/web_server/middleware/user/plugins"

	"github.com/gin-gonic/gin"
	"github.com/holmeswang/contrib/sessions"
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

//应依照"bk_obj_id"和"bk_property_type":"objuser"查询"cc_ObjAttDes"集合(已建立索引),得到"bk_property_id"的值;
//然后以它的值为key,取得Info中的value,然后以value作为param访问ESB,得到其中文名。
func (s *Service) getUserMapFromESBNew(c *gin.Context, objID string, infoList []mapstr.MapStr) (map[string]string, []string, error) {
	rid := util.GetHTTPCCRequestID(c.Request.Header)
	var defErr errors.DefaultCCErrorIf
	cond := metadata.QueryCondition{
		Fields: []string{metadata.AttributeFieldPropertyID},
		Condition: map[string]interface{}{
			metadata.AttributeFieldObjectID:     objID,
			metadata.AttributeFieldPropertyType: common.FieldTypeUser,
		},
	}
	attrRsp, err := s.CoreAPI.CoreService().Model().ReadModelAttr(c, c.Request.Header, objID, &cond)
	if nil != err {
		blog.Errorf("failed to request the object controller, err: %s, rid: %s", err.Error(), rid)
		return nil, nil, err
	}
	if !attrRsp.Result {
		blog.Errorf("failed to search the object(%s), err: %s, rid: %s", objID, attrRsp.ErrMsg, rid)
		return nil, nil, err
	}

	usernameList := []string{}
	propertyList := []string{}
	for _, info := range infoList {
		if objID == common.BKInnerObjIDHost {
			info = info[common.BKInnerObjIDHost].(map[string]interface{})
		}
		for _, item := range attrRsp.Data.Info {
			propertyList = append(propertyList, item.PropertyID)
			if info[item.PropertyID] != nil {
				usernameList = append(usernameList, info[item.PropertyID].(string))
			}
		}
	}
	propertyList = util.RemoveDuplicatesAndEmpty(propertyList)
	userList := util.RemoveDuplicatesAndEmpty(usernameList)
	userListStr := strings.Join(userList, ",")

	usernameMap := map[string]string{}

	if len(userList) != 0 {
		params := make(map[string]string)
		params["exact_lookups"] = userListStr
		params["fields"] = "username,display_name"
		user := plugins.CurrentPlugin(c, s.Config.Version)
		//如果是skip-auth模式，这里mock一个返回值吧。
		userListEsb, skipLogin := mockUserList(c, rid)
		if skipLogin != true {
			var errNew *errors.RawErrorInfo
			userListEsb, errNew = user.GetUserList(c, s.Config.ConfigMap)
			if errNew != nil {
				blog.ErrorJSON("get user list from ESB failed, err: %s, rid: %s", errNew.ToCCError(defErr).Error(), rid)
				userListEsb = []*metadata.LoginSystemUserInfo{}
				return nil, nil, err
			}
		}
		for _, userInfo := range userListEsb {
			username := fmt.Sprintf("%s(%s)", userInfo.EnName, userInfo.CnName)
			usernameMap[userInfo.EnName] = username
			//usernameMap["admin"] = "admin(admin)"
			//usernameMap["Daniel-Wu"] = "Daniel-Wu(吴彦祖)"
		}
	}
	return usernameMap, propertyList, nil
}

func mockUserList(c *gin.Context, rid string) ([]*metadata.LoginSystemUserInfo, bool) {
	session := sessions.Default(c)
	skipLogin := session.Get(webCommon.IsSkipLogin)
	skipLogins, ok := skipLogin.(string)
	if ok && "1" == skipLogins {
		blog.V(5).Infof("use skip login flag: %v, rid: %s", skipLogin, rid)
		adminData := []*metadata.LoginSystemUserInfo{
			{
				CnName: "admin",
				EnName: "admin",
			},
			{
				CnName: "吴彦祖",
				EnName: "Daniel-Wu",
			},
		}
		return adminData, true
	} else {
		return nil, false
	}
}
