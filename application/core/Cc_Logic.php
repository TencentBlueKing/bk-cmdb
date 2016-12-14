<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Cc_Logic extends CI_Logic {
    public $_errInfo = '';

    public function __construct() {
        $this->_errInfo = '';
    }

    /**
     * @function   构造json消息串，用于传递请求处理结果
     * @param $isSuccess      请求是否处理成功
     * @param string $errCodeKey  错误消息定义
     * @return array  消息传递体
     */
    public function getOutput($isSuccess, $errCodeKey='') {
        $ret = array();
        $ret['success'] = $isSuccess;

        if (! $isSuccess) {
            $ret['errCode'] = $errCodeKey;
        }

        return $ret;
    }
}