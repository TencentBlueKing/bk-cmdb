package service

import (
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/storage/mock"
	"github.com/emicklei/go-restful"
	"net/http"
	"net/http/httptest"
	"strings"
)

func NewRestfulTestCase(data string) (*Service, *restful.Request, *restful.Response) {
	req, resp := NewRestfulRequestResponse(data)
	return NewService(), req, resp
}

func NewRestfulRequestResponse(data string) (*restful.Request, *restful.Response) {
	bodyReader := strings.NewReader(data)
	httpRequest, _ := http.NewRequest("POST", "/", bodyReader)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set(common.BKHTTPOwnerID, "")
	request := &restful.Request{Request: httpRequest}

	resp := &restful.Response{ResponseWriter: httptest.NewRecorder()}

	resp.SetRequestAccepts("application/json;application/xml")
	resp.WriteHeader(200)
	resp.Write([]byte("ok"))

	return request, resp
}

func NewService() *Service {
	core := &backbone.Engine{
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
	}
	return &Service{Core: core, Instance: &mock.MockDB{}, Cache: nil}
}
