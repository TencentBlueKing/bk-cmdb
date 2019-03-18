package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

type RequestContext struct {
	// http header
	Header http.Header
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

	Metadata metadata.Metadata
}

type parseStream struct {
	RequestCtx *RequestContext
	Attribute  *meta.AuthAttribute
	err        error
	action     meta.Action
}

func newParseStream(rc *RequestContext) (*parseStream, error) {
	if nil == rc {
		return nil, errors.New("request context is nil")
	}

	return &parseStream{RequestCtx: rc}, nil
}

// parse is used to parse the auth attribute from RequestContext.
func (ps *parseStream) Parse() (*meta.AuthAttribute, error) {
	if ps.err != nil {
		return nil, ps.err
	}

	ps.validateAPI().
		validateVersion().
		validateResourceAction().
		validateUserAndSupplier().
		hostRelated().
		topology().
		topologyLatest().
		netCollectorRelated().
		processRelated().
		eventRelated().
		// finalizer must be at the end of the check chains.
		finalizer()

	if ps.err != nil {
		return nil, ps.err
	}

	return ps.Attribute, nil
}

func (ps *parseStream) validateAPI() *parseStream {
	if ps.err != nil {
		return ps
	}

	if ps.RequestCtx.Elements[0] != "api" {
		ps.err = errors.New("unsupported api format")
	}

	return ps
}

func (ps *parseStream) validateVersion() *parseStream {
	if ps.err != nil {
		return ps
	}

	version := ps.RequestCtx.Elements[1]
	if version != "v3" {
		ps.err = fmt.Errorf("unsupported version %s", version)
		return ps
	}
	ps.Attribute.APIVersion = version

	return ps
}

func (ps *parseStream) validateResourceAction() *parseStream {
	if ps.err != nil {
		return ps
	}

	action := ps.RequestCtx.Elements[2]
	switch action {
	case "find":
		ps.action = meta.Find
	case "findMany":
		ps.action = meta.FindMany

	case "create":
		ps.action = meta.Create
	case "createMany":
		ps.action = meta.CreateMany

	case "update":
		ps.action = meta.Update
	case "updateMany":
		ps.action = meta.UpdateMany

	case "delete":
		ps.action = meta.Delete
	case "deleteMany":
		ps.action = meta.DeleteMany

	default:
		ps.action = meta.Unknown
		// to compatible api that is not this kind of format,
		// this err will not be set, but it will be set when
		// all the api is normalized.

		// TODO: uncomment this err code.
		// ps.err = fmt.Errorf("unsupported action %s", action)
		return ps
	}

	return ps
}

// user and supplier account must be set in the http
// request header, otherwise, an error will be occur.
func (ps *parseStream) validateUserAndSupplier() *parseStream {
	if ps.err != nil {
		return ps
	}

	// validate user header at first.
	user := ps.RequestCtx.Header.Get(common.BKHTTPHeaderUser)
	if len(user) == 0 {
		ps.err = fmt.Errorf("request lost header: %s", common.BKHTTPHeaderUser)
		return ps
	}
	ps.Attribute.User.UserName = user

	// validate the supplier account now.
	supplier := ps.RequestCtx.Header.Get(common.BKHTTPOwnerID)
	if len(supplier) == 0 {
		ps.err = fmt.Errorf("request lost header: %s", common.BKHTTPOwnerID)
		return ps
	}
	ps.Attribute.User.SupplierAccount = supplier

	return ps
}

// finalizer is to find whether a url resource has been matched or not.
func (ps *parseStream) finalizer() *parseStream {
	if ps.err != nil {
		return ps
	}
	ps.err = errors.New("unsupported resource operation")
	return ps
}

func (ps *parseStream) hitRegexp(reg *regexp.Regexp, httpMethod string) bool {
	return reg.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == httpMethod
}

func (ps *parseStream) hitPattern(pattern, httpMethod string) bool {
	return pattern == ps.RequestCtx.URI && ps.RequestCtx.Method == httpMethod
}
