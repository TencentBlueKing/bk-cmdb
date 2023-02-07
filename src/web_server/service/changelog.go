package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/version"
	"configcenter/src/thirdparty/hooks"
	webCommon "configcenter/src/web_server/common"

	"github.com/gin-gonic/gin"
)

// VersionListItem single version data
type VersionListItem struct {
	Version    string `json:"version"`
	UpdateTime string `json:"time"`
	IsCurrent  bool   `json:"is_current"`
}

// VersionListResult the return body of the version list data
type VersionListResult struct {
	metadata.BaseResp `json:",inline"`
	Data              []VersionListItem `json:"data"`
}

// VersionDetailResult the return body of the version detail data
type VersionDetailResult struct {
	metadata.BaseResp `json:",inline"`
	Data              string `json:"data"`
}

// GetVersionList  get all cmdb version info from changelog path directory
func (s *Service) GetVersionList(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	language := webCommon.GetLanguageByHTTPRequest(c)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	changelogPath, err := cc.String("webServer.changelogPath.ch")
	if err != nil {
		blog.Errorf("configuration file missing [%s] configuration item, err: %v, rid: %s",
			"webServer.changelogPath.ch", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommConfMissItem,
			ErrMsg: defErr.CCErrorf(common.CCErrCommConfMissItem, "webServer.changelogPath.ch").Error(),
		})
		return
	}

	files, err := os.ReadDir(changelogPath)
	if err != nil {
		blog.Errorf("failed to read directory: %s, err: %v, rid: %s", changelogPath, err, rid)
		c.JSON(http.StatusOK, VersionListResult{
			BaseResp: metadata.BaseResp{
				Result:      false,
				Code:        0,
				ErrMsg:      err.Error(),
				Permissions: nil,
			},
			Data: nil,
		})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusOK, VersionListResult{
			BaseResp: metadata.BaseResp{
				Result:      true,
				Code:        0,
				ErrMsg:      "",
				Permissions: nil,
			},
			Data: nil,
		})
		return
	}

	versionInfoList := getVersionInfoList(files)

	c.JSON(http.StatusOK, VersionListResult{
		BaseResp: metadata.BaseResp{
			Result:      true,
			Code:        0,
			ErrMsg:      "",
			Permissions: nil,
		},
		Data: versionInfoList,
	})
}

// GetVersionDetail get specific version changelog
func (s *Service) GetVersionDetail(c *gin.Context) {
	header := c.Request.Header
	rid := util.GetHTTPCCRequestID(header)
	language := webCommon.GetLanguageByHTTPRequest(c)

	option := new(metadata.ChangelogDetailConfigOption)
	if err := json.NewDecoder(c.Request.Body).Decode(option); err != nil {
		blog.Errorf("get version detail failed, decode body: %+v, err: %v, rid: %s", c.Request.Body, err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrCommJSONUnmarshalFailed,
			ErrMsg: "json unmarshal error",
		})
		return
	}

	changelogPath, errCode, err := getChangelogPath(language)
	if err != nil {
		blog.Errorf("failed to get changelog path, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   errCode,
			ErrMsg: fmt.Sprintf("failed to get changelog path, error: %v", err),
		})
		return
	}

	versionFilePath, err := getVersionFilePath(changelogPath, option.Version)
	if err != nil {
		// 找不到指定的版本数据则返回对应错误
		blog.Errorf("the changelog file for %s could not be found, err: %v, rid: %s", option.Version, err, rid)
		c.JSON(http.StatusOK, VersionDetailResult{
			BaseResp: metadata.BaseResp{
				Result:      false,
				Code:        0,
				ErrMsg:      err.Error(),
				Permissions: nil,
			},
			Data: "",
		})
		return
	}

	versionData, err := os.ReadFile(versionFilePath)
	if err != nil {
		blog.Errorf("failed to open file, err: %v, rid: %s", err, rid)
		c.JSON(http.StatusOK, metadata.BaseResp{
			Result: false,
			Code:   common.CCErrWebOpenFileFail,
			ErrMsg: fmt.Sprintf("failed to open file, error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, VersionDetailResult{
		BaseResp: metadata.BaseResp{
			Result:      true,
			Code:        0,
			ErrMsg:      "",
			Permissions: nil,
		},
		Data: string(versionData),
	})
	return
}

// getVersionInfoList get all cmdb version info from changelog files
func getVersionInfoList(files []os.DirEntry) []VersionListItem {
	versionInfoList := make([]VersionListItem, 0)
	for _, file := range files {
		fileVersion, updateTime := getFileVersion(file.Name())
		if fileVersion == "" {
			continue
		}
		// 找出当前版本并把IsCurrent设为true
		versionInfoList = append(versionInfoList, VersionListItem{
			Version:    fileVersion,
			UpdateTime: updateTime,
			IsCurrent:  fileVersion == getCurrentVersion(),
		})
	}
	return versionInfoList
}

// getVersionFilePath gets the version log path specified in the request body
// eg: if the version in the request body is v3.10.aa, v3.10.23-rc, v3.10.22-alpha, return "" and error
func getVersionFilePath(changelogPath string, version string) (string, error) {
	versionRegex := "^v\\d+\\.\\d+\\.\\d+$"
	if hooks.GetVersionRegexHook() != "" {
		versionRegex = hooks.GetVersionRegexHook()
	}
	matched, err := regexp.MatchString(versionRegex, version)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("version: " + version + " does not conform to the version number format")
	}

	files, err := os.ReadDir(changelogPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: " + changelogPath)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no files in " + changelogPath)
	}

	for _, file := range files {
		fileVersion, _ := getFileVersion(file.Name())
		if fileVersion == "" {
			continue
		}
		if fileVersion == version {
			return filepath.Join(changelogPath, file.Name()), nil
		}
	}
	return "", fmt.Errorf("no changelog file for " + version + " in " + changelogPath)
}

// getCurrentVersion Get the current version
// eg: returns the version according to the existing tag
//     release-v3.10.22-alpha1, return "v3.10.22"
//     release-v3.10.18_alpha1, return "v3.10.18"
//     release-v3.10.16, return "v3.10.16"
//     release-v3.10.x_feature-agent-id_alpha, return ""
func getCurrentVersion() string {
	currentVersion := version.CCVersion
	currentVersionRegex := "v\\d+\\.\\d+\\.\\d+(-|_|$)"
	if hooks.GetCurrentVersionRegexHook() != "" {
		currentVersionRegex = hooks.GetCurrentVersionRegexHook()
	}
	reg := regexp.MustCompile(currentVersionRegex)
	currentVersion = reg.FindString(currentVersion)
	if currentVersion == "" {
		return ""
	}
	//版本日志文件由产品于验收通过版本（rc版本）时一起出
	// 去后缀操作：
	// 用于在产品进行功能验证时保证当前版本号（带后缀）与不带后缀的版本号的版本日志能匹配得上。
	// 例如：当前版本号为v3.10.23-rc，与之对应的版本日志的版本号为v3.10.23
	if strings.Index(currentVersion, "-") != -1 || strings.Index(currentVersion, "_") != -1 {
		return currentVersion[:len(currentVersion)-1]
	}
	return currentVersion
}

// getFileVersion get version and updateTime in filename
// eg: test.md, _test.md, test_test.txt; return "", ""
//     vaa.bb.cc_2006-01-02.md, v3.10.22_2022-02-29.md; return "", ""
//     v3.10.23-rc_2006-01-02.md, v3.10.22-alpha_2006-01-02.md; return "", ""
//     v3.10.22_2022-03-18.md; return v3.10.22, 2022-03-18
func getFileVersion(filename string) (string, string) {
	matched, err := regexp.MatchString("^.+_.+\\.md$", filename)
	if err != nil {
		blog.Errorf("match the changelog file name failed, err: %v", err)
		return "", ""
	}
	if !matched {
		return "", ""
	}

	filename = filename[:len(filename)-3]
	// 判断文件名中版本号的格式是否符合要求
	// 版本日志文件名不会出现版本号带后缀的情况
	versionRegex := "^v\\d+\\.\\d+\\.\\d+$"
	if hooks.GetVersionRegexHook() != "" {
		versionRegex = hooks.GetVersionRegexHook()
	}
	matched, err = regexp.MatchString(versionRegex, strings.Split(filename, "_")[0])
	if err != nil {
		blog.Errorf("matches version in the changelog file name failed, err: %v", err)
		return "", ""
	}
	if !matched {
		return "", ""
	}

	// 判断文件名中的发布时间的格式是否符合要求
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		blog.Errorf("matches updateTime in the changelog file name failed, err: %v", err)
		return "", ""
	}
	_, err = time.ParseInLocation("2006-01-02", strings.Split(filename, "_")[1], local)
	if err != nil {
		blog.Errorf("the updateTime in %s.md does not conform to date format, err: %v", filename, err)
		return "", ""
	}

	fileVersion := strings.Split(filename, "_")[0]
	updateTime := strings.Split(filename, "_")[1]
	return fileVersion, updateTime
}

// getChangelogPath get changelogPath from the common config
func getChangelogPath(language string) (string, int, error) {
	var (
		confItem      string
		changelogPath string
		err           error
	)
	switch common.LanguageType(language) {
	case common.Chinese:
		confItem = "webServer.changelogPath.ch"
		changelogPath, err = cc.String(confItem)
		if err != nil {
			return "", common.CCErrCommConfMissItem, fmt.Errorf("configuration file missing [%s] configuration item"+
				", err: %v", confItem, err)
		}
	case common.English:
		confItem = "webServer.changelogPath.en"
		changelogPath, err = cc.String(confItem)
		if err != nil {
			return "", common.CCErrCommConfMissItem, fmt.Errorf("configuration file missing [%s] configuration item"+
				", err: %v", confItem, err)
		}
	default:
		return "", common.CCErrCommParamsInvalid, fmt.Errorf("'language' data parameter verification does not pass")
	}
	return changelogPath, 0, err
}
