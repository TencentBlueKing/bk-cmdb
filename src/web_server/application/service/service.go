/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/httpserver/webserver"
	"configcenter/src/common/language"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/storage/redisclient"
	confCenter "configcenter/src/web_server/application/config"
	"configcenter/src/web_server/application/logics"
	"configcenter/src/web_server/application/middleware"
	"configcenter/src/web_server/application/rdiscover"
	webCommon "configcenter/src/web_server/common"
)

//CCAPIServer define data struct of bcs ccapi server
type CCWebServer struct {
	conf     *config.CCAPIConfig
	httpServ *gin.Engine
	rd       *rdiscover.RegDiscover
	cfCenter *confCenter.ConfCenter
}

func NewCCWebServer(conf *config.CCAPIConfig) (*CCWebServer, error) {
	s := &CCWebServer{}

	//config
	s.conf = conf
	addr, _ := s.conf.GetAddress()
	port, _ := s.conf.GetPort()

	s.httpServ = gin.Default()
	a := api.NewAPIResource()
	a.SetConfig(s.conf)
	a.InitWaction()

	//RDiscover
	s.rd = rdiscover.NewRegDiscover(s.conf.RegDiscover, addr, port, false)
	a.AddrSrv = s.rd
	a.APIAddr = rdapi.GetRdAddrSrvHandle(types.CC_MODULE_APISERVER, a.AddrSrv)

	//	a.Lang = language.New()

	//ConfCenter
	s.cfCenter = confCenter.NewConfCenter(s.conf.RegDiscover)
	return s, nil
}

//Stop the ccapi server
func (ccWeb *CCWebServer) Stop() error {
	return nil
}

//Start the web server
func (ccWeb *CCWebServer) Start() error {
	chErr := make(chan error, 2)

	// configure center
	go func() {
		err := ccWeb.cfCenter.Start()
		blog.Errorf("configure center module start failed!. err:%s", err.Error())
		chErr <- err
	}()

	/// fetch config of itselft
	var confData []byte
	_ = confData
	for {
		confData = ccWeb.cfCenter.GetConfigureCtx()

		if confData == nil {
			blog.Warnf("fail to get configure, will get again")
			time.Sleep(time.Second * 2)
			continue
		} else {
			blog.Infof("get configure. ctx(%s)", string(confData))
			break
		}
	}

	a := api.NewAPIResource()
	config, _ := a.ParseConf(confData)

	// load the language resource
	if dirPath, ok := config["language.res"]; ok {
		if res, err := language.New(dirPath); nil != err {
			blog.Error("failed to create language object, error info is  %s ", err.Error())
			chErr <- err
		} else {
			a.Lang = res
		}
	} else {
		for {
			langCtx := ccWeb.cfCenter.GetLanguageResCxt()
			if langCtx == nil {
				blog.Warnf("fail to get language package, will get again")
				time.Sleep(time.Second * 2)
				continue
			} else {
				languageif := language.NewFromCtx(langCtx)
				a.Lang = languageif
				blog.Info("lanugage package loaded")
				break
			}
		}
	}

	// load the errors resource
	if dirPath, ok := config["erros.res"]; ok {
		if res, err := errors.New(dirPath); nil != err {
			blog.Error("failed to create errors object, error info is  %s ", err.Error())
			chErr <- err
		} else {
			a.Error = res
		}
	} else {
		for {
			errCtx := ccWeb.cfCenter.GetErrorResCxt()
			if errCtx == nil {
				blog.Warnf("fail to get errors package, will get again")
				time.Sleep(time.Second * 2)
				continue
			} else {
				errIf := errors.NewFromCtx(errCtx)
				a.Error = errIf
				blog.Info("lanugage erros loaded")
				break
			}
		}
	}

	site := config["site.domain_url"] + "/"
	version := config["api.version"]
	loginURL := config["site.bk_login_url"]
	appCode := config["site.app_code"]
	check_url := config["site.check_url"]
	sessionName := config["session.name"]
	skipLogin := config["session.skip"]
	defaultlanguage := config["session.defaultlanguage"]
	if "" == defaultlanguage {
		defaultlanguage = "zh-cn"
	}
	apiSite, _ := a.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
	static := config["site.html_root"]
	webCommon.ResourcePath = config["site.resources_path"]
	redisIp := config["session.host"]
	redisPort := config["session.port"]
	redisSecret := config["session.secret"]
	multipleOwner := config["session.multiple_owner"]
	agentAppUrl := config["app.agent_app_url"]
	redisSecret = strings.TrimSpace(redisSecret)
	curl := fmt.Sprintf(loginURL, appCode, site)

	redisCli, err := redisclient.NewRedis(redisIp, redisPort, "", redisSecret, "0")
	if nil != err {
		blog.Errorf("connect redis error %s", err.Error())
		return err
	}
	err = redisCli.Open()
	if nil != err {
		blog.Errorf("connect redis error %s", err.Error())
		return err
	}
	a.CacheCli = redisCli
	go func() {
		store, rediserr := sessions.NewRedisStore(10, "tcp", redisIp+":"+redisPort, redisSecret, []byte("secret"))
		if rediserr != nil {
			panic(rediserr)
		}
		ccWeb.httpServ.Use(sessions.Sessions(sessionName, store))
		ccWeb.httpServ.Use(middleware.Cors())

		ccWeb.RegisterActions(a.Wactions)
		middleware.APIAddr = rdapi.GetRdAddrSrvHandle(types.CC_MODULE_APISERVER, a.AddrSrv)
		ccWeb.httpServ.Use(middleware.ValidLogin(loginURL, appCode, site, check_url, apiSite, skipLogin, multipleOwner, defaultlanguage))
		ccWeb.httpServ.Static("/static", static)
		blog.Info(static)
		ccWeb.httpServ.LoadHTMLFiles(static + "/index.html") //("static/index.html")
		// MetricServer
		conf := metric.Config{
			ModuleName:    types.CC_MODULE_WEBSERVER,
			ServerAddress: ccWeb.conf.AddrPort,
		}
		metricActions := metric.NewMetricController(conf, ccWeb.HealthMetric)
		for _, metricAction := range metricActions {
			newmetricAction := metricAction
			ccWeb.httpServ.GET(newmetricAction.Path, func(c *gin.Context) {
				newmetricAction.HandlerFunc(c.Writer, c.Request)
			})
		}
		ccWeb.httpServ.GET("/", func(c *gin.Context) {
			session := sessions.Default(c)
			role := session.Get("role")
			userName, _ := session.Get("userName").(string)
			language, _ := session.Get("language").(string)
			apiSite, err := a.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
			if nil != err {
				blog.Errorf("api server not start %s", err.Error())
			}
			userPriviApp, rolePrivilege, modelPrivi, sysPrivi, mainLineObjIDArr := logics.GetUserAppPri(apiSite, userName, common.BKDefaultOwnerID, language)
			var strUserPriveApp, strRolePrivilege, strModelPrivi, strSysPrivi, mainLineObjIDStr string
			if nil == userPriviApp {
				strUserPriveApp = ""
			} else {
				cstrUserPriveApp, _ := json.Marshal(userPriviApp)
				strUserPriveApp = string(cstrUserPriveApp)
			}

			if nil == rolePrivilege {
				strRolePrivilege = ""
			} else {
				cstrRolePrivilege, _ := json.Marshal(rolePrivilege)
				strRolePrivilege = string(cstrRolePrivilege)
			}
			if nil == modelPrivi {
				strModelPrivi = ""
			} else {
				cstrModelPrivi, _ := json.Marshal(modelPrivi)
				strModelPrivi = string(cstrModelPrivi)
			}
			if nil == sysPrivi {
				strSysPrivi = ""
			} else {
				cstrSysPrivi, _ := json.Marshal(sysPrivi)
				strSysPrivi = string(cstrSysPrivi)
			}

			mainLineObjIDB, _ := json.Marshal(mainLineObjIDArr)
			mainLineObjIDStr = string(mainLineObjIDB)

			session.Set("userPriviApp", string(strUserPriveApp))
			session.Set("rolePrivilege", string(strRolePrivilege))

			session.Set("modelPrivi", string(strModelPrivi))
			session.Set("sysPrivi", string(strSysPrivi))
			session.Set("mainLineObjID", string(mainLineObjIDStr))
			session.Save()

			//set cookie
			appIDArr := make([]string, 0)
			for key, _ := range userPriviApp {
				appIDArr = append(appIDArr, strconv.FormatInt(key, 10))
			}
			appIDStr := strings.Join(appIDArr, "-")
			c.SetCookie("bk_privi_biz_id", appIDStr, 24*60*60, "", "", false, false)

			c.HTML(200, "index.html", gin.H{
				"site":        site,
				"version":     version,
				"role":        role,
				"curl":        curl,
				"userName":    userName,
				"agentAppUrl": agentAppUrl,
			})
		})

		ip, _ := ccWeb.conf.GetAddress()
		port, _ := ccWeb.conf.GetPort()
		portStr := strconv.Itoa(int(port))
		addr := ip + ":" + portStr
		err := ccWeb.httpServ.Run(addr)

		blog.Error("http listen and serve failed! err:%s", err.Error())
		chErr <- err
	}()

	//start rdiscover
	go func() {
		err := ccWeb.rd.Start()
		blog.Errorf("rdiscover start failed! err:%s", err.Error())
		chErr <- err
	}()

	select {
	case err := <-chErr:
		blog.Error("exit! err:%s", err.Error())
		return err
	}

	return nil
}

func (ccWeb *CCWebServer) RegisterActions(actions []*webserver.Action) {
	//fmt.Println(actions)
	for _, action := range actions {
		switch action.Verb {
		case "GET":
			ccWeb.httpServ.GET(action.Path, action.Handler)
		case "POST":
			ccWeb.httpServ.POST(action.Path, action.Handler)
		case "PUT":
			ccWeb.httpServ.PUT(action.Path, action.Handler)
		case "DELETE":
			ccWeb.httpServ.DELETE(action.Path, action.Handler)
		case "OPTIONS":
			ccWeb.httpServ.OPTIONS(action.Path, action.Handler)
		default:
			blog.Error("unrecognized action verb: %s", action.Verb)
		}
	}
}

// HealthMetric check netservice is health
func (ccWeb *CCWebServer) HealthMetric() metric.HealthMeta {
	meta := metric.HealthMeta{IsHealthy: true}

	// check zk
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityServicediscover, ccWeb.rd.Ping()))
	// check dependence
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CC_MODULE_APISERVER, metric.CheckHealthy(middleware.APIAddr())))

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "webserver is not healthy"
			break
		}
	}

	return meta
}
