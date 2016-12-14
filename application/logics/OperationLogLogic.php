<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class OperationLogLogic extends Cc_Logic {
    public function __construct() {
        parent::__construct();
    }

    /**
     * @查询操作日志
     * @param $operator  string 操作者username
     * @param $appId     int 业务Id
     * @param $opTime    string 操作时间
     * @param $opType    string 操作类型(新增、删除、查询)
     * @param $opTarget  string 操作对象(业务、大区、模块、主机...)
     * @param $opContent string 操作内容
     * @param $clientIp  string 操作者Ip
     * @param $start     int 数据起始下标
     * @param $limit     int 数据偏移量
     * @param $orderBy   string sql查询排序字段名
     * @param $direction string 排序方式(asc、desc)
     * @return array 日志信息
     */
    public function getOperationLog($operator = '', $appId = 0, $opTime = '', $opType = '', $opTarget = '', $opContent = '', $clientIp = '', $start = 0, $limit = 20, $orderBy = '', $direction = 'DESC') {
        try {
            $this->load->model('OperationLogModel');
            $operationLog = $this->OperationLogModel->getOperationLog($operator, $appId, $opTime, $opType, $opTarget, $opContent, $clientIp, $start, $limit, $orderBy, $direction);
            return $operationLog;
        } catch (Exception $e) {
            CCLog::LogErr('OperationLogLogic->getOperationLog:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @获取用户操作日志
     * @return array
     */
    public function getUserOperationLog($operator = '', $appId = 0, $opType = '', $opTarget = '', $opContent = '', $clientIp = '', $startTime = '', $endTime = '', $start = 0, $limit = 0, $orderBy = '', $direction = 'DESC') {
        try {
            $this->load->model('OperationLogModel');
            $operationLog = $this->OperationLogModel->getUserOperationLog($operator, $appId, $opType, $opTarget, $opContent, $clientIp, $startTime, $endTime, $start, $limit, $orderBy, $direction);

            return $operationLog;
        } catch (Exception $e) {
            CCLog::LogErr('OperationLogLogic->getUserOperationLog:' . $e->getMessage());
            return array();
        }
    }
}