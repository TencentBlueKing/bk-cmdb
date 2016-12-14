<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH.'config/api_error_const.php';

class Api_Controller extends CI_Controller {

    public $errInfo;

    public function __construct() {
        parent::__construct();
    }

    /**
     * @输出json消息
     * @param result json消息体 array 结构
     * @return 输出json消息
     */
    public function outputJson($result) {
        $out = json_encode($result);
        $contentType = 'application/json';
        return $this->output->set_content_type($contentType)->set_output($out);
    }

    /**
     * @输出成功消息
     * @param data 返回数据
     * @return 输出json消息
     */
    public function outSuccess($data) {
        $result = array('code'=>0, 'data'=>$data);
        return $this->outputJson($result);
    }

    /**
     * @输出错误信息
     * @param message 返回消息
     * @param extMsg 消息附带信息
     * @return 输出json消息
     */
    public function outFailure($code, $msg, $extMsg) {
        $result = array('code' => $code, 'msg' => $msg, 'extmsg' => $extMsg);
        return $this->outputJson($result);
    }
}