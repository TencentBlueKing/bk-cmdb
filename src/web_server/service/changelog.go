package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

type VersionInfo struct {
	Version    string `json:"version"`
	UpdateTime string `json:"time"`
	IsCurrent  bool   `json:"is_current"`
}

type versionDataResult struct {
	Result    bool        `json:"result"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

// GetVersionList get all cmdb versions
func (s *Service) GetVersionList(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)

	var (
		versionInfoList []VersionInfo
		isCurrent       bool
	)

	changelogPath, err := cc.ChangeLogPath("ChangeLogPath")
	if err != nil {
		blog.Errorf("get changelog path failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommGetCommConf,
			ErrMsg: "get changelog path failed",
		})
		return
	}

	files, err := os.ReadDir(changelogPath["ch"])
	if err != nil {
		blog.Errorf("failed to open file, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrWebOpenFileFail,
			ErrMsg: fmt.Sprintf("failed to open file, error: %v", err),
		})
		return
	}

	for _, file := range files {
		filename := file.Name()[:len(file.Name())-3]
		// 找出当前版本并把IsCurrent设为true
		if strings.Split(filename, "_")[0] == version.CCVersion {
			isCurrent = true
		} else {
			isCurrent = false
		}
		versionInfoList = append(versionInfoList, VersionInfo{
			Version:    strings.Split(filename, "_")[0],
			UpdateTime: strings.Split(filename, "_")[1],
			IsCurrent:  isCurrent,
		})
	}

	c.JSON(200, versionDataResult{
		Result:    true,
		Code:      0,
		Message:   "success",
		RequestId: rid,
		Data:      versionInfoList,
	})
}

// GetVersionDetail get specific version changelog
func (s *Service) GetVersionDetail(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	var (
		changeLogPath   string
		versionFilePath string
	)

	changelogDetailConfig := new(metadata.ChangelogDetailConfig)
	if err := json.NewDecoder(c.Request.Body).Decode(changelogDetailConfig); err != nil {
		blog.Errorf("get version detail failed, decode body err: %v, body: %+v, rid: %s", err, c.Request.Body, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: "json unmarshal error",
		})
		return
	}

	changeLogPathMap, err := cc.ChangeLogPath("ChangeLogPath")
	if err != nil {
		blog.Errorf("get changelog path failed, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommGetCommConf,
			ErrMsg: "get changelog path failed",
		})
		return
	}

	switch language {
	case "zh-cn":
		changeLogPath = changeLogPathMap["ch"]
	case "en":
		changeLogPath = changeLogPathMap["en"]
	default:
		result := metadata.ResponseDataMapStr{
			BaseResp: metadata.BaseResp{
				Result: false,
				Code:   common.CCErrCommParamsInvalid,
				ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, "language").Error(),
			},
		}
		c.JSON(http.StatusOK, result)
		return
	}

	// 获取请求体中指定的版本日志
	files, err := os.ReadDir(changeLogPath)
	if err != nil {
		blog.Errorf("failed to open file, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrWebOpenFileFail,
			ErrMsg: fmt.Sprintf("failed to open file, error: %v", err),
		})
		return
	}
	for _, file := range files {
		filename := file.Name()[:len(file.Name())-3]
		if strings.Split(filename, "_")[0] == changelogDetailConfig.Version {
			versionFilePath = filepath.Join(changeLogPath, file.Name())
		}
	}
	// 找不到指定的版本数据则返回version参数错误
	if versionFilePath == "" {
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommParamsInvalid,
			ErrMsg: defErr.Errorf(common.CCErrCommParamsInvalid, "version").Error(),
		})
		return
	}

	versionData, err := os.ReadFile(versionFilePath)
	if err != nil {
		blog.Errorf("failed to open file, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusBadRequest, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrWebOpenFileFail,
			ErrMsg: fmt.Sprintf("failed to open file, error: %v", err),
		})
		return
	}

	c.JSON(200, versionDataResult{
		Result:    true,
		Code:      0,
		Message:   "success",
		RequestId: rid,
		Data:      string(versionData),
	})
	return
}
