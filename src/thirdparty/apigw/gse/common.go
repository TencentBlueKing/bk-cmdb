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

package gse

import (
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/apigwutil"
)

// ListAgentStateRequest use to list gse agent status request
type ListAgentStateRequest struct {
	AgentIDList []string `json:"agent_id_list"`
}

// ListAgentStateResp use to list gse agent status response
type ListAgentStateResp struct {
	apigwutil.ApiGWBaseResponse
	Data []ListAgentStateData `json:"data"`
}

// ListAgentStateData the data in list agent state response
type ListAgentStateData struct {
	BKAgentID  string `json:"bk_agent_id"`
	BKCloudID  int64  `json:"bk_cloud_id"`
	Version    string `json:"version"`
	RunMode    int64  `json:"run_mode"`
	StatusCode int    `json:"status_code"`
}

// AsyncPushFileRequest use to push file to host request
type AsyncPushFileRequest struct {
	TimeoutSeconds int64   `json:"timeout_seconds,omitempty"`
	AutoMkdir      bool    `json:"auto_mkdir,omitempty"`
	Tasks          []*Task `json:"tasks"`
}

// Task the task about push file
type Task struct {
	FileName    string   `json:"file_name"`
	StoreDir    string   `json:"store_dir"`
	FileContent string   `json:"file_content"`
	Owner       string   `json:"owner"`
	Right       int32    `json:"right"`
	AgentIDList []string `json:"agent_id_list"`
}

// AsyncPushFileResp use to push file to host response
type AsyncPushFileResp struct {
	apigwutil.ApiGWBaseResponse
	Data AsyncPushFileData `json:"data"`
}

// AsyncPushFileData the data in push file response
type AsyncPushFileData struct {
	Result AsyncPushFileResult `json:"result"`
}

// AsyncPushFileResult the push file result
type AsyncPushFileResult struct {
	TaskID string `json:"task_id"`
}

// GetTransferFileResultRequest the request about get push file result
type GetTransferFileResultRequest struct {
	TaskID      string   `json:"task_id"`
	AgentIDList []string `json:"agent_id_list"`
}

// GetTransferFileResultResp the response about get push file result
type GetTransferFileResultResp struct {
	apigwutil.ApiGWBaseResponse
	Data GetTransferFileResultData `json:"data"`
}

// GetTransferFileResultData the data in the response about get push file result
type GetTransferFileResultData struct {
	Result []GetTransferFileResult `json:"result"`
}

// GetTransferFileResult the result in the response about get push file data
type GetTransferFileResult struct {
	ErrorCode int64                        `json:"error_code"`
	ErrorMsg  string                       `json:"error_msg"`
	Content   GetTransferFileResultContent `json:"content"`
}

// GetTransferFileResultContent the content about each host push file result
type GetTransferFileResultContent struct {
	Protover       int64   `json:"protover"`
	Mode           int64   `json:"mode"`
	Type           string  `json:"type"`
	Progress       float64 `json:"progress"`
	Size           int64   `json:"size"`
	Speed          float64 `json:"speed"`
	StartTime      int64   `json:"start_time"`
	EndTime        int64   `json:"end_time"`
	SourceAgentID  string  `json:"source_agent_id"`
	DestAgentID    string  `json:"dest_agent_id"`
	SourceFileDir  string  `json:"source_file_dir"`
	SourceFileName string  `json:"source_file_name"`
	DestFileDir    string  `json:"dest_file_dir"`
	DestFileName   string  `json:"dest_file_name"`
	Status         int64   `json:"status"`
	StatusInfo     string  `json:"status_info"`
}

// AddStreamToResp add streamto response
type AddStreamToResp struct {
	apigwutil.ApiGWBaseResponse
	Data *metadata.GseConfigAddStreamToResult `json:"data"`
}

// QueryStreamToResp query streamto response
type QueryStreamToResp struct {
	apigwutil.ApiGWBaseResponse
	Data []metadata.GseConfigAddStreamToParams `json:"data"`
}

// AddRouteResp add route response
type AddRouteResp struct {
	apigwutil.ApiGWBaseResponse
	Data *metadata.GseConfigAddRouteResult `json:"data"`
}

// QueryRouteResp query route response
type QueryRouteResp struct {
	apigwutil.ApiGWBaseResponse
	Data []metadata.GseConfigChannel `json:"data"`
}
