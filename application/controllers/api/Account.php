<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH . 'core/Api_Controller.php';

class Account extends Api_Controller {
    public function __construct() {
        parent::__construct();
        $this->load->library('login');
    }

    /**
     * 判断用户登录鉴权
     */
    public function authed() {
        $ccToken = $this->input->cookie('cc_token');
        $response = array('success' => true);
        if (empty($ccToken)) {
            $ccToken = $this->input->get_post('cc_token');
        }
        $result = $this->login->isCcTokenValid($ccToken);
        if (empty($result)) {
            $response['success'] = false;
            $response['message'] = $this->login->getErrMsg();
        } else {
            unset($result['Password']);
            $response['user'] = $result;
        }
        $this->output->set_output(json_encode($response));
    }

    /**
     * 判断用户登录鉴权
     */
    public function logout() {
        $ccToken = $this->input->cookie('cc_token');
        $response = array('success' => true);
        if (empty($ccToken)) {
            $ccToken = $this->input->get_post('cc_token');
        }
        $result = $this->login->isCcTokenValid($ccToken);
        if (empty($result)) {
            $response['success'] = false;
            $response['message'] = $this->login->getErrMsg();
        } else {
            $this->login->logoutUser();
        }
        $this->output->set_output(json_encode($response));
    }
}