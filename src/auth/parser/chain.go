package parser

import (
	"errors"
	"fmt"
	"sync"

	"configcenter/src/auth"
)

type RequestContext struct {
	// http method
	Method string
	// request's url path
	URI string
	// elements parsed from url, started with api field
	// 0: api field
	// 1: version field
	// 2: action field
	// >=3: resource fields
	Elements []string
	// http request body contents.
	Body []byte

	Attribute *auth.Attribute
}

// next means if the next chain func should be called.
type chainFunc func(ctx *RequestContext) (next bool, err error)

var chains []chainFunc
var once sync.Once

func init() {
	once.Do(func() {
		// do not change the chain sequences.
		chains = []chainFunc{
			validateAPI,
			handleVersion,
		}
	})
}

func validateAPI(ctx *RequestContext) (bool, error) {
	if ctx.Elements[0] != "api" {
		return false, errors.New("unsupported api format")
	}
	return true, nil
}

func handleVersion(ctx *RequestContext) (bool, error) {
	version := ctx.Elements[1]
	if version != "v3" {
		return false, fmt.Errorf("unsupported version %s", version)
	}
	ctx.Attribute.Resource.APIVersion = version
	return true, nil
}
