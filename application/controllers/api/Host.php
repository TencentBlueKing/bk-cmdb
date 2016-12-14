<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH . 'core/Api_Controller.php';

class Host extends Api_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 根据setId获取主机列表
     */
    public function getSetHostList() {
        $appId = intval($this->input->post('ApplicationID'));
        $setIdArr = explode(',', $this->input->post('SetID'));

        Utility::getNumericInArray($setIdArr);

        if ($appId === 0) {
            $obj = $this->config->item(API_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        if (empty($setIdArr)) {
            return $this->outSuccess(array());
        }

        $this->load->Logic('ModuleHostConfigLogic');
        $hosts = $this->ModuleHostConfigLogic->getHostsBySetIdAndAppId($setIdArr, $appId);

        return $this->outSuccess($hosts);
    }

    /**
     * 根据模块Id获取主机
     */
    public function getModuleHostList() {
        $appId = intval($this->input->post('ApplicationID'));
        $moduleIdArr = explode(',', $this->input->post('ModuleID'));

        Utility::getNumericInArray($moduleIdArr);

        if ($appId === 0) {
            $obj = $this->config->item(CC_API_APPID_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        if (empty($moduleIdArr)) {
            return $this->outSuccess(array());
        }

        $this->load->Logic('ModuleHostConfigLogic');
        $hosts = $this->ModuleHostConfigLogic->getHostsByModuleIdAndAppId($moduleIdArr, $appId);

        if (!$hosts) {
            $obj = $this->config->item(CC_API_NOEST_HOST);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        return $this->outSuccess($hosts);
    }

    /**
     * 根据IP获取主机
     */
    public function getHostListByIP() {
        $appId = intval($this->input->post('ApplicationID'));
        $ipArr = explode(',', $this->input->post('IP'));

        foreach ($ipArr as $key => $value) {
            if (filter_var($value, FILTER_VALIDATE_IP)) {
                $ipArr[$key] = $value;
            } else {
                unset($ipArr[$key]);
            }
        }

        if (0 === $appId) {
            $obj = $this->config->item(CC_API_APPID_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        if (empty($ipArr)) {
            return $this->outSuccess(array());
        }

        $this->load->Logic('ModuleHostConfigLogic');
        $hosts = $this->ModuleHostConfigLogic->getHostsByIPAndAppId($ipArr, $appId);

        $hostIdArr = array();
        foreach ($hosts as $value) {
            $hostIdArr[] = $value['HostID'];
        }

        if (!$hosts) {
            $obj = $this->config->item(CC_API_NOEST_HOST);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        return $this->outSuccess($hosts);
    }

    /**
     * 根据业务获取主机
     */
    public function getAppHostList() {
        $appId = intval($this->input->post('ApplicationID'));
        if (0 === $appId) {
            $obj = $this->config->item(CC_API_APPID_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }

        $this->load->Logic('ModuleHostConfigLogic');
        $hosts = $this->ModuleHostConfigLogic->getHostsByIPAndAppId(array(), $appId);
        $this->outSuccess($hosts);
        return;
    }

    /**
     * 新增主机
     */
    public function addHost() {
        $innerIP = $this->input->post('InnerIP', true);
        $outerIP = $this->input->post('OuterIP',true);
        $hostName = $this->input->post('HostName', true);
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);
        #内网IP非法
        if (!$innerIP || !filter_var($innerIP, FILTER_VALIDATE_IP)) {
            $obj = $this->config->item(CC_API_INV_INPUT_IP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        #外网IP
        if($outerIP  &&  !filter_var($outerIP, FILTER_VALIDATE_IP)) {
            $obj = $this->config->item(CC_API_INV_INPUT_IP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('HostBaseLogic');
        $hosts = $this->HostBaseLogic->addHost($innerIP, $outerIP, $hostName, $operator, $bakOperator);
        if($hosts !== true ) {
            $obj = $this->config->item($hosts);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        return $this->outSuccess(array());
    }

    /**
     * 删除主机
     */
    public function delHost() {
        $innerIP = $this->input->post('InnerIP', true);
        #内网IP非法
        if (!$innerIP || !filter_var($innerIP, FILTER_VALIDATE_IP)) {
            $obj = $this->config->item(CC_API_INV_INPUT_IP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('HostBaseLogic');
        $hosts = $this->HostBaseLogic->getHostByIp($innerIP);
        if(empty($hosts)) {
            $obj = $this->config->item(CC_API_GET_IP_FAIL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $appId = $hosts[0]['ApplicationID'];
        $hostId = $hosts[0]['HostID'];
        $this->HostBaseLogic->RealDeleteHost($hostId, $appId);
        return $this->outSuccess(array());
    }

    /**
     * 编辑主机
     */
    public function editHost() {
        $innerIP = $this->input->post('InnerIP', true);
        $outerIP = $this->input->post('OuterIP',true);
        $hostName = $this->input->post('HostName', true);
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);
        #内网IP非法
        if (!$innerIP || !filter_var($innerIP, FILTER_VALIDATE_IP)) {
            $obj = $this->config->item(CC_API_INV_INPUT_IP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        #外网IP
        if($outerIP  &&  !filter_var($outerIP, FILTER_VALIDATE_IP)) {
            $obj = $this->config->item(CC_API_INV_INPUT_IP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('HostBaseLogic');
        $hosts = $this->HostBaseLogic->editHost($innerIP, $outerIP, $hostName, $operator, $bakOperator);
        if($hosts !== true ) {
            $obj = $this->config->item($hosts);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        return $this->outSuccess(array());
    }
}