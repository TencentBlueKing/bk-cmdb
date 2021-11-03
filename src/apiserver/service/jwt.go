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
	"fmt"

	"github.com/golang-jwt/jwt"
)

type App struct {
	Version   int64  `json:"version"`
	AppCode   string `json:"app_code"`
	BKAppCode string `json:"bk_app_code"`
	Verified  bool   `json:"verified"`
}

type User struct {
	Version    int64  `json:"version"`
	UserName   string `json:"username"`
	BKUserName string `json:"bk_username"`
	Verified   bool   `json:"verified"`
}

type Claims struct {
	App  App  `json:"app"`
	User User `json:"user"`
	jwt.RegisteredClaims
}

// ParseToken parse token
func ParseToken(token, jwtSecret string) (*Claims, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(jwtSecret)) // 解析公钥
	if err != nil {
		return nil, err
	}
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
