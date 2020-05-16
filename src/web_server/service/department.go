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
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"

	"github.com/gin-gonic/gin"
)

// GetDepartment get department info
func (s *Service) GetDepartment(c *gin.Context) {
	rspBody := metadata.DepartmentResult{}
	data, err := s.Logics.GetDepartment(c, s.Config)
	if err != nil {
		blog.ErrorJSON("get department error. err:%s", err)
		rspBody.Result = false
		rspBody.Code = err.GetCode()
		rspBody.ErrMsg = err.Error()
	} else {
		rspBody.Data = data
		rspBody.Result = true

	}

	c.JSON(http.StatusOK, rspBody)
	return
}

// GetDepartmentProfile get department info
func (s *Service) GetDepartmentProfile(c *gin.Context) {
	rspBody := metadata.DepartmentProfileResult{}

	data, err := s.Logics.GetDepartmentProfile(c, s.Config)
	if err != nil {
		rspBody.Result = false
		rspBody.Code = err.GetCode()
		rspBody.ErrMsg = err.Error()
	} else {
		rspBody.Data = data
		rspBody.Result = true

	}

	c.JSON(http.StatusOK, rspBody)
	return
}
