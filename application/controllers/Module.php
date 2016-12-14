<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Module extends Cc_Controller {

    public function __construct() {
        parent::__construct();
        $this->load->logic('ModuleBaseLogic');
    }

    /**
     * @根据Id获取模块信息
     * @return json
     */
    public function getModuleInfoById() {
        $this->load->logic('ModuleBaseLogic');

        $appId = intval($this->input->get_post('ApplicationID', true));
        $moduleId = intval($this->input->get_post('ModuleID', true));

        $module = $this->ModuleBaseLogic->getModuleById($moduleId, array(), $appId);
        $data = array();
        $data['module'] = $module[0];
        $data['success'] = TRUE;
        $this->output->set_output(json_encode($data));
        return;
    }

    /**
     * @新建模块
     * @return json
     */
    public function newModule() {
        $this->load->Logic('ModuleBaseLogic');

        $appId = intval($this->input->post('ApplicationID', true));
        $setId = intval($this->input->post('SetID', true));
        $moduleName = trim(htmlspecialchars($this->input->post('ModuleName', true)));
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);

        if (mb_strlen($moduleName) == 0 || mb_strlen($moduleName) > 60) {
            $this->outputJson(FALSE, 'module_length_error');
        }
        if ($appId == 0 || $setId == 0) {
            $this->outputJson(FALSE, 'appid_or_setid_inneed');
        }
        $result = $this->ModuleBaseLogic->addModule($appId, $setId, $moduleName, $operator, $bakOperator);
        $errCode = isset($result['errCode']) ? ($result['errCode']) : '';
        return $this->outputJson($result['success'], $errCode);
    }

    /*
     * @编辑模块
     * @return json
     */
    public function editModule() {
        $this->load->Logic('ModuleBaseLogic');

        $appId = intval($this->input->post('ApplicationID', true));
        $moduleId = $this->input->post('ModuleID', true);
        $moduleName = trim(htmlspecialchars($this->input->post('ModuleName', true)));
        $setId = intval($this->input->post('SetID', true));
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);

        if (mb_strlen($moduleName) == 0 || mb_strlen($moduleName) > 60) {
            $this->outputJson(FALSE, 'module_length_error');
        }
        if ($appId == 0 || empty($appId)) {
            $this->outputJson(FALSE, 'appid_or_moduleid_inneed');
        }

        $result = $this->ModuleBaseLogic->editModule($appId, $setId, $moduleId, $moduleName, $operator, $bakOperator);

        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        return $this->outputJson($result['success'], $errCode);
    }

    /*
     * @删除模块
     * @return json
     */
    public function delModule() {
        $this->load->Logic('ModuleBaseLogic');

        $appId = intval($this->input->post('ApplicationID', true));
        $moduleId = intval($this->input->post('ModuleID', true));
        $setId = intval($this->input->post('SetID', true));

        if ($appId == 0 || $moduleId == 0 || $setId == 0) {
            $this->outputJson(FALSE, 'appid_or_moduleid_or_setid_inneed');
        }

        $result = $this->ModuleBaseLogic->delModule($appId, $moduleId, $setId);
        $errCode = isset($result['errCode']) ? $result['errCode'] : '';

        $this->outputJson($result['success'], $errCode);
    }
}