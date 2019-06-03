package core

import (
	"context"
	"gopkg.in/redis.v5"
	"net/http"

	"configcenter/src/auth/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

type ContextParams struct {
	context.Context
	Header          http.Header
	SupplierAccount string
	User            string
	ReqID           string
	Error           errors.DefaultCCErrorIf
	Lang            language.DefaultCCLanguageIf
}

type Operation struct {
	*backbone.Engine
	header      http.Header
	ccErr       errors.DefaultCCErrorIf
	ccLang      language.DefaultCCLanguageIf
	user        string
	ownerID     string
	cache       *redis.Client
	AuthManager *extensions.AuthManager
}

type StatisticalOperation interface {
	CreateInnerChart(ctx ContextParams, inner string) ([]metadata.ExceptionResult, error)
}
