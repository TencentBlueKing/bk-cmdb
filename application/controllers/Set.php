<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Set extends Cc_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 根据setId获取set信息
     * @return set数组
     */
    public function getSetInfoById() {
        $appID = intval($this->input->post('ApplicationID', true));
        $setID = intval($this->input->post('SetID', true));

        $this->load->Logic('SetBaseLogic');
        $set = $this->SetBaseLogic->getSetById($setID, $appID);
        $data = array();
        $data['set'] = current($set);
        $data['success'] = true;
        return $this->output->set_output(json_encode($data));
    }

    /**
     * 获取业务下所有集群信息
     * @param ApplicationID
     * @return set信息
     */
    public function getAllSetInfo() {
        $appId = intval($this->input->post('ApplicationID', true));
        $this->load->Logic('SetBaseLogic');
        $set = $this->SetBaseLogic->getSetById(array(), $appId);

        $data = array();
        $data['set'] = $set;
        $data['success'] = true;
        return $this->output->set_output(json_encode($data));
    }

    /**
     * 新建集群
     * @return json
     */
    public function newSet() {
        $appId = intval($this->input->post('ApplicationID', true));
        $setName = trim(htmlspecialchars($this->input->post('SetName', true)));
        $chnName = trim(htmlspecialchars($this->input->post('ChnName', true)));
        $enviType = intval($this->input->post('EnviType', true));
        $serviceStatus = intval($this->input->post('ServiceStatus', true));
        $capacity = $this->input->post('Capacity', true);
        $des = trim(htmlspecialchars($this->input->post('Des', true)));
        $openStatus = trim(htmlspecialchars($this->input->post('Openstatus', true)));
        if ($appId === 0) {
            return $this->outputJson(FALSE, 'applicationid_validation_failure');
        }
        if (mb_strlen($setName) === 0 || mb_strlen($setName) > 64) {
            $this->outputJson(FALSE, 'setname_length_error');
            return;
        }

        if (mb_strlen($chnName) > 32) {
            $this->outputJson(FALSE, 'setch_length_error');
            return;
        }

        if (mb_strlen($openStatus) > 16) {
            $this->outputJson(FALSE, 'set_Openstatus_length_error');
            return;
        }

        if (strlen($serviceStatus) == 0 || mb_strlen($serviceStatus) > 16) {
            $this->outputJson(FALSE, 'setch_length_error');
            return;
        }

        if (mb_strlen($des) > 250) {
            $this->outputJson(FALSE, 'des_outof_length');
            return;
        }
        $this->load->Logic('SetBaseLogic');
        $result = $this->SetBaseLogic->newSet($appId, $setName, $chnName, $enviType, $serviceStatus, $capacity, $des, $openStatus);
        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        return $this->outputJson($result['success'], $errCode);
    }

    /*
     * 新建集群
     * @return json
     */
    public function editSet() {
        $appId = intval($this->input->post('ApplicationID', true));
        $setId = $this->input->post('SetID', true);
        $setName = trim(htmlspecialchars($this->input->post('SetName', true)));
        $chnName = trim(htmlspecialchars($this->input->post('ChnName', true)));
        $enviType = intval($this->input->post('SetEnviType', true));
        $serviceStatus = intval(htmlspecialchars($this->input->post('ServiceStatus', true)));
        $capacity = $this->input->post('Capacity', true);
        $des = trim(htmlspecialchars($this->input->post('Des', true)));
        $openStatus = trim(htmlspecialchars($this->input->post('Openstatus', true)));
        if ($appId == 0 || $setId == 0) {
            return $this->outputJson(FALSE, 'appid_or_setid_inneed');
        }

        if (isset($setName) && (mb_strlen($setName) > 64)) {
            $this->outputJson(FALSE, 'setname_length_error');
            return;
        }

        if (isset($chnName) && (mb_strlen($chnName) > 32)) {
            $this->outputJson(FALSE, 'setch_length_error');
            return;
        }

        if (isset($openStatus) && mb_strlen($openStatus) > 16) {
            $this->outputJson(FALSE, 'set_Openstatus_length_error');
            return;
        }

        if (isset($ServiceStatus) && (strlen($ServiceStatus) == 0 || mb_strlen($ServiceStatus) > 16)) {
            $this->outputJson(FALSE, 'setservice_status_error');
            return;
        }
        if (mb_strlen($des) > 250) {
            $this->outputJson(FALSE, 'des_outof_length');
            return;
        }
        $this->load->logic('SetBaseLogic');
        $result = $this->SetBaseLogic->editSet($appId, $setId, $setName, $chnName, $enviType, $serviceStatus, $capacity, $des, $openStatus);
        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        return $this->outputJson($result['success'], $errCode);
    }

    /*
     * 删除集群
     * @param ApplicationID,SetID
     * @return json
     */
    public function delSet() {
        $setId = intval($this->input->post('SetID', true));
        $appId = intval($this->input->post('ApplicationID', true));

        if ($setId == 0 || $appId == 0) {
            $this->outputJson(false, 'appid_or_setid_inneed');
            return;
        }
        $this->load->Logic('SetBaseLogic');
        $result = $this->SetBaseLogic->delSet($appId, $setId);
        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        return $this->outputJson($result['success'], $errCode);
    }

    /*
     * 克隆集群
     * @return json
     */
    public function cloneSet() {
        $appId = intval($this->input->post('ApplicationID', true));
        $setId = intval($this->input->post('SetID', true));
        $setNameStr = $this->input->post('SetName', true);
        $this->load->Logic('SetBaseLogic');

        $setNameArr = array_unique(explode("\n", $setNameStr));
        foreach($setNameArr as $key=>&$set) {
            if(empty($set)) {
                unset($setNameArr[$key]);
            }
        }

        $setNameArr = array_values($setNameArr);
        $result = $this->SetBaseLogic->cloneSet($appId, $setId, $setNameArr);

        if (0 == count($result)) {
            return $this->outputJson(true);
        } else {
            $failSetStr = implode(',', $result);
            $errMsg = '集群 【' . $failSetStr . ' 】重复，请重新输入需要克隆的集群';
            $result = array();
            $result['success'] = false;
            $result['errCode'] = '2011';
            $result['errInfo'] = $errMsg;
            return $this->output->set_output(json_encode($errMsg));
        }
    }
}