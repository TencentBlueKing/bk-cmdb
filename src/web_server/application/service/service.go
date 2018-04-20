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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/core/cc/config"
	"configcenter/src/common/http/httpserver/webserver"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	confCenter "configcenter/src/web_server/application/config"
	"configcenter/src/web_server/application/logics"
	"configcenter/src/web_server/application/middleware"
	"configcenter/src/web_server/application/rdiscover"
	webCommon "configcenter/src/web_server/common"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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
		confData = ccWeb.cfCenter.GetConfigureCxt()

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

	site := config["site.domain_url"] + "/"
	version := config["api.version"]
	loginURL := config["site.bk_login_url"]
	appCode := config["site.app_code"]
	check_url := config["site.check_url"]
	sessionName := config["session.name"]
	skipLogin := config["session.skip"]
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
	go func() {
		store, rediserr := sessions.NewRedisStore(10, "tcp", redisIp+":"+redisPort, redisSecret, []byte("secret"))
		if rediserr != nil {
			panic(rediserr)
		}
		ccWeb.httpServ.Use(sessions.Sessions(sessionName, store))
		ccWeb.httpServ.Use(middleware.Cors())

		ccWeb.RegisterActions(a.Wactions)
		middleware.APIAddr = rdapi.GetRdAddrSrvHandle(types.CC_MODULE_APISERVER, a.AddrSrv)
		ccWeb.httpServ.Use(middleware.ValidLogin(loginURL, appCode, site, check_url, apiSite, skipLogin, multipleOwner))
		ccWeb.httpServ.Static("/static", static)
		blog.Info(static)
		ccWeb.httpServ.LoadHTMLFiles(static + "/index.html") //("static/index.html")
		ccWeb.httpServ.GET("/", func(c *gin.Context) {
			session := sessions.Default(c)
			role := session.Get("role")
			userName, _ := session.Get("userName").(string)
			language, _ := session.Get("language").(string)
			apiSite, err := a.AddrSrv.GetServer(types.CC_MODULE_APISERVER)
			if nil != err {
				blog.Errorf("api server not start %s", err.Error())
			}
			userPriviApp, rolePrivilege, modelPrivi, sysPrivi := logics.GetUserAppPri(apiSite, userName, common.BKDefaultOwnerID, language)
			var strUserPriveApp, strRolePrivilege, strModelPrivi, strSysPrivi string
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

			session.Set("userPriviApp", string(strUserPriveApp))
			session.Set("rolePrivilege", string(strRolePrivilege))

			session.Set("modelPrivi", string(strModelPrivi))
			session.Set("sysPrivi", string(strSysPrivi))
			session.Save()

			c.HTML(200, "index.html", gin.H{
				"site":        site,
				"version":     version,
				"role":        role,
				"curl":        curl,
				"userName":    userName,
				"agentAppUrl": agentAppUrl,
			})
		})

		// MetricServer
		conf := metric.Config{
			ModuleName: types.CC_MODULE_PROCCONTROLLER,
		}
		metricActions := metric.NewMetricController(conf, ccWeb.HealthMetric)
		for _, metricAction := range metricActions {
			ccWeb.httpServ.GET(metricAction.Path, func(c *gin.Context) {
				metricAction.HandlerFunc(c.Writer, c.Request)
			})
		}

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
	fmt.Println(actions)
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
	a := api.GetAPIResource()

	// check mongo
	mongoHealthy := metric.HealthItem{Name: "mongo"}
	if err := a.InstCli.Ping(); err != nil {
		mongoHealthy.IsHealthy = false
		mongoHealthy.Message = err.Error()
	} else {
		mongoHealthy.IsHealthy = true
	}
	meta.Items = append(meta.Items, mongoHealthy)

	// check redis
	redisHealthy := metric.HealthItem{Name: "redis"}
	if err := a.CacheCli.Ping(); err != nil {
		redisHealthy.IsHealthy = false
		redisHealthy.Message = err.Error()
	} else {
		redisHealthy.IsHealthy = true
	}
	meta.Items = append(meta.Items, redisHealthy)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "proccontroller is not healthy"
			break
		}
	}

	return meta
}
