/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"bytes"
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/constant"
	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	authpb "github.com/TencentBlueKing/bk-cmdb/pkg/proto/auth-server"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// newGrpcMux creates a new http mux using grpc gateway.
func (s *Service) newGrpcMux(ctx context.Context) (*runtime.ServeMux, error) {
	// create grpc mux
	mux := runtime.NewServeMux(
		// defines the protobuf to json marshal and unmarshal options
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &grpcGWMarshaller{
			JSONPb: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:     true,
					UseEnumNumbers:    true,
					EmitDefaultValues: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		runtime.WithIncomingHeaderMatcher(grpcGWHeaderMatcher),
		// defines the error handlers
		runtime.WithErrorHandler(grpcGWErrHandler),
		runtime.WithRoutingErrorHandler(grpcGWRoutingErrHandler))

	// register grpc gateway http handlers
	registerHandlers := map[config.ServiceName]func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error{
		config.AuthServer: authpb.RegisterAuthHandler,
	}

	for service, registerHandler := range registerHandlers {
		if err := registerHandler(ctx, mux, s.grpcClients[service]); err != nil {
			log.Error(ctx, "register grpc gateway http handler failed", "service", service, log.E(err))
			return nil, err
		}
	}

	return mux, nil
}

// grpcGWMarshaller is the protobuf to json grpcGWMarshaller for grpc gateway.
type grpcGWMarshaller struct {
	*runtime.JSONPb
}

// Marshal marshals the protobuf message to json, wrapper the raw message with cmdb base response.
func (m *grpcGWMarshaller) Marshal(v any) ([]byte, error) {
	dataBytes, err := m.JSONPb.Marshal(v)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString(`{"data":`)
	buf.Write(dataBytes)
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

var grpcHeaderMap = map[string]struct{}{
	constant.UserHeader:         {},
	constant.AppCodeHeader:      {},
	constant.TenantHeader:       {},
	constant.HTTPLanguageHeader: {},
}

// grpcGWHeaderMatcher is the grpc gateway incoming header matcher.
func grpcGWHeaderMatcher(s string) (string, bool) {
	_, exists := grpcHeaderMap[s]
	if !exists {
		return "", false
	}
	return s, true
}

// grpcGWErrHandler is the error handler for grpc gateway.
func grpcGWErrHandler(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter,
	r *http.Request, err error) {

	kt := kit.NewKitFromHeader(ctx, r.Header)
	st, ok := status.FromError(err)
	if !ok {
		_ = rest.APIError(kt, err).Render(w)
		return
	}

	statusCode := runtime.HTTPStatusFromCode(st.Code())
	errCode := cerr.GetErrCodeByHTTPStatus(statusCode)
	err = cerr.NewError(errCode, st.Message())
	_ = rest.APIErrorWithStatus(kt, err, statusCode).Render(w)
}

// grpcGWRoutingErrHandler is the http routing error handler for grpc gateway.
func grpcGWRoutingErrHandler(ctx context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter,
	_ *http.Request, httpStatus int) {

	kt := kit.GetGrpcKit(ctx)
	errCode := cerr.GetErrCodeByHTTPStatus(httpStatus)
	err := cerr.NewError(errCode, http.StatusText(httpStatus))
	_ = rest.APIErrorWithStatus(kt, err, httpStatus).Render(w)
}
