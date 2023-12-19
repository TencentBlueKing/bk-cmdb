/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

// Package jwt defines jwt related logics
package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/header/util"

	"github.com/golang-jwt/jwt/v4"
)

var handler = &jwtHandler{
	Enabled: false,
}

// GetHandler get jwt handler
func GetHandler() Handler {
	return handler
}

// Handler is jwt handler for http header
type Handler interface {
	Parse(header http.Header) (http.Header, error)
	Sign(header http.Header) (http.Header, error)
}

// Init the jwt handler.
func Init(prefix string) error {
	configKey := "jwt"
	if prefix != "" {
		configKey = prefix + ".jwt"
	}

	if !cc.IsExist(configKey) {
		return nil
	}

	conf := new(config)
	if err := cc.UnmarshalKey(configKey, conf); err != nil {
		return fmt.Errorf("parse jwt config by prefix %s failed, err: %v", prefix, err)
	}

	if !conf.Enabled {
		blog.Infof("%s jwt is disabled", prefix)
		return nil
	}
	handler.Enabled = true

	if conf.PublicKey != "" {
		jwtPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(conf.PublicKey))
		if err != nil {
			return fmt.Errorf("parse jwt public key %s failed, err: %v", conf.PublicKey, err)
		}

		handler.PublicKey = jwtPublicKey
	}

	if conf.PrivateKey != "" {
		jwtPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(conf.PrivateKey))
		if err != nil {
			return fmt.Errorf("parse jwt private key %s failed, err: %v", conf.PrivateKey, err)
		}

		handler.PrivateKey = jwtPrivateKey
	}

	return nil
}

type config struct {
	// Enabled is the flag to enable jwt authorization
	Enabled bool `mapstructure:"enabled"`
	// PublicKey is the jwt public key
	PublicKey string `mapstructure:"publicKey"`
	// PrivateKey is the jwt private key
	PrivateKey string `mapstructure:"privateKey"`
}

// jwtHandler used to parse requests from blueking api-gateway.
type jwtHandler struct {
	// Enabled is the flag to enable jwt authorization
	Enabled bool
	// PublicKey is the public key to parse jwt token from blueking api-gateway http request
	PublicKey *rsa.PublicKey
	// PrivateKey is the private key to parse jwt token from blueking api-gateway http request
	PrivateKey *rsa.PrivateKey
}

// Parse jwt info from api-gateway header to cc header
func (j *jwtHandler) Parse(header http.Header) (http.Header, error) {
	if !j.Enabled {
		// compatible for esb request
		// TODO remove this when esb is not supported
		header = util.ConvertLegacyHeader(header)
		return header, nil
	}

	jwtToken := httpheader.GetBkJWT(header)
	if len(jwtToken) == 0 {
		// compatible for esb request
		header = util.ConvertLegacyHeader(header)
		// TODO returns error when esb is not supported
		return header, nil
	}

	token, err := j.parseToken(jwtToken)
	if err != nil {
		return nil, err
	}

	if err = token.validate(); err != nil {
		return nil, err
	}

	header.Set("Content-Type", "application/json")
	httpheader.SetUser(header, token.User.UserName)
	httpheader.SetAppCode(header, token.App.AppCode)

	return header, nil
}

// Sign jwt info by cc header
func (j *jwtHandler) Sign(header http.Header) (http.Header, error) {
	if !j.Enabled {
		return header, nil
	}

	if j.PrivateKey == nil {
		return nil, errors.New("jwt private key is not set")
	}

	info := &claims{
		App: &app{
			Version:  1,
			AppCode:  httpheader.GetAppCode(header),
			Verified: true,
		},
		User: &user{
			Version:  1,
			UserName: httpheader.GetUser(header),
			Verified: true,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, info)

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := token.SignedString(j.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("sign jwt info failed, err: %v", err)
	}

	httpheader.SetBkJWT(header, tokenStr)

	return header, nil
}

// claims blueking api gateway jwt struct
type claims struct {
	App  *app  `json:"app"`
	User *user `json:"user"`
	jwt.RegisteredClaims
}

// app blueking application info
type app struct {
	Version  int64  `json:"version"`
	AppCode  string `json:"app_code"`
	Verified bool   `json:"verified"`
}

// validate app
func (a *app) validate() error {
	if !a.Verified {
		return errors.New("app is not verified")
	}
	return nil
}

// user blueking user info
type user struct {
	Version  int64  `json:"version"`
	UserName string `json:"username"`
	// Verified 用户是否通过认证，免用户认证的时候是false，其他情况下是true，因此我们目前不需要关注这个字段
	Verified bool `json:"verified"`
}

// validate user
func (u *user) validate() error {
	if u.UserName == "" {
		return errors.New("user is not set")
	}
	return nil
}

// validate claims
func (c *claims) validate() error {
	if c.App == nil {
		return errors.New("app info is not set")
	}

	if err := c.App.validate(); err != nil {
		return err
	}

	if c.User == nil {
		return errors.New("user info is not set")
	}

	if err := c.User.validate(); err != nil {
		return err
	}

	return nil
}

// parseToken parse jwt token
func (j *jwtHandler) parseToken(token string) (*claims, error) {
	if j.PublicKey == nil {
		return nil, errors.New("jwt public key is not set")
	}

	tokenClaims, err := jwt.ParseWithClaims(token, new(claims), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if tokenClaims == nil {
		return nil, errors.New("jwt parse result is empty")
	}

	info, ok := tokenClaims.Claims.(*claims)
	if !ok {
		return nil, fmt.Errorf("token claims type %T is invalid", tokenClaims.Claims)
	}

	if !tokenClaims.Valid {
		return nil, errors.New("token claims is invalid")
	}

	return info, nil
}
