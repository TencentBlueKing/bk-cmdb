<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH . 'core/Api_Controller.php';

class App extends Api_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 查询所有业务
     */
    public function getAppList() {
        $this->load->Logic('ApplicationBaseLogic');
        $appInfo = $this->ApplicationBaseLogic->getAppList();   //获取业务信息

        return $this->outSuccess($appInfo);
    }

    /**
     * 查询用户有权限的业务列表
     */
    public function getAppByUin() {
        $userName = $this->input->post('userName') ? $this->input->post('userName') : 0;

        if ($userName === 0) {
            $obj = $this->config->item(CC_UINORCOM_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('UserLogic');

        $appInfo = $this->ApplicationBaseLogic->getAppByUin($userName);
        $userInfo = $this->UserLogic->getUserInfo($userName);
        $appCompany = $appInfo['appCompany'];
        $appUin = $appInfo['appUin'];
        if ('admin' == $userInfo['Role']) {

            if (0 == count($appCompany)) {
                $obj = $this->config->item(CC_API_NOEST_APP);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            } else if (1 == count($appCompany)) {
                $obj = $this->config->item(CC_API_ONLYDEF_APP);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            }
            return $this->outSuccess($appCompany);

        } else {

            if (0 == count($appCompany)) {
                $obj = $this->config->item(CC_API_NOEST_APP);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            }

            if (1 == count($appCompany) && 1 == count($appUin)) {
                $obj = $this->config->item(CC_API_ONLYDEF_APP);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            }

            if (count($appCompany) > 1 && 1 == count($appUin)) {
                $obj = $this->config->item(CC_API_ONRIG_APP);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            }

            return $this->outSuccess($appUin);
        }
    }
}