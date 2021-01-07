/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service_test

import (
	"net/http"
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/source_controller/coreservice/app/options"
	"configcenter/src/source_controller/coreservice/service"
	"configcenter/src/storage/dal/mongo"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
)

var defaultHeader = func() http.Header {

	header := http.Header{}
	header.Set(common.BKHTTPCCRequestID, "test_req_id")
	header.Set(common.BKHTTPHeaderUser, "test_user")
	header.Set(common.BKHTTPOwnerID, "test_owner")
	header.Set(common.BKHTTPLanguage, "en")
	return header
}()

func startCoreService(t *testing.T, ip string, port uint) {

	// create a logics service
	coreService := service.New()

	// register the server hander
	bonServer := backbone.Server{
		ListenAddr: ip,
		ListenPort: port,
		Handler:    restful.NewContainer().Add(coreService.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	// set backbone config
	bonC := &backbone.Config{
		Server: bonServer,
	}

	// new server instance
	engine, err := backbone.NewMockBackbone(bonC)
	require.NoError(t, err)

	errE, err := errors.New("../../../../resources/errors/")
	require.NoError(t, err)
	lan, err := language.New("../../../../resources/language/")
	require.NoError(t, err)
	engine.CCErr = errE
	engine.Language = lan
	// set the config
	coreService.SetConfig(options.Config{
		Mongo: mongo.Config{
			Address:  "cc:cc@localhost:27010,localhost:27011,localhost:27012,localhost:27013",
			User:     "cc",
			Password: "cc",
			Database: "cmdb",
		},
	}, engine, engine.CCErr, engine.Language)

	return
}
