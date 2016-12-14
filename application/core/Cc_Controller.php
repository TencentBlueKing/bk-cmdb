<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Cc_Controller extends CI_Controller {
    public $_errInfo;
    protected $p_noAuthActions;

    public function __construct() {
        parent::__construct();
        $this->load->library('parser');
        $this->lang->load('zh_cn');
        $this->_errInfo = '';

        # 登录验证
        $this->setNoAuthActions();
        $segments = $this->uri->rsegments;
		$action = $segments[2];
        if(!in_array($action, $this->p_noAuthActions)){
            $this->load->library('login');
            if(!$this->login->isAuthed()){
                $this->login->redirectUnauthed();
            }
        }
    }

    /**
     * 子类覆盖此方法， 把不需要登录校验的action配置在p_noAuthActions中。
     */
    protected function setNoAuthActions(){
        $this->p_noAuthActions = array();
    }

    /**
     * @构造json消息串，用于传递请求处理结果。
     * @param $isSuccess       请求是否处理成功
     * @param string $errCodeKey  错误消息定义
     * @return array  消息传递体
     */
    public function outputJson($isSuccess, $errCodeKey='') {
        $ret = array();
        $ret['success'] = $isSuccess;

        if (! $isSuccess) {
            $this->config->item($errCodeKey);
            $ret['errCode'] = $this->config->item($errCodeKey)->Code;
            $ret['errInfo'] = $this->config->item($errCodeKey)->Info;
            $ret['message'] = $this->config->item($errCodeKey)->Info;
        }

        return $this->output->set_output(json_encode($ret));
    }

    /**
     * @构造json消息串，用于传递请求处理结果。
     * @param $isSuccess       请求是否处理成功
     * @param string $message  上一层错误信息
     * @return array  消息传递体
     */
    public function outputJsonByMessage($isSuccess, $message = '') {
        $ret = array();
        $ret['success'] = $isSuccess;
        $ret['message'] = $message;

        return $this->output->set_output(json_encode($ret));
    }

    /**
     * 构造数组，用于记录页面数据
     * @param $header
     * @param $activePage
     * @return array
     */
    protected function buildPageDataArr($header, $activePage) {
        $data = array();
        $data['header'] = $header;
        $data['active'] = $activePage;

        return $data;
    }

    /**
     * @构造返回数据的json消息串
     * @return array  消息传递体
     */
    public function outputDataJson($info = array(), $total = 0) {
        $data['data'] = $info;
        $data['total'] = $total;

        return $this->output->set_output(json_encode($data));
    }

    /**
     * @返回用户展示格式
     * @return json
     */
    public function getUserDisplayName($user, $key) {
        return isset($user[$key]) ? $key . '(' . $user[$key] . ')' : $key;
    }
}