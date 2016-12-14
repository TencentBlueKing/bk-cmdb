<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH . 'core/Api_Controller.php';

class Module extends Api_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 查询所有模块
     */
    public function getModules() {
        $appId = intval($this->input->post('ApplicationID'));
        $this->load->Logic('ModuleBaseLogic');
        $result = $this->ModuleBaseLogic->getModulesByAppId($appId);
        $result = array_column($result, 'ModuleName');

        return $this->outSuccess($result);
    }

    /**
     * 新增模块
     */
    public function addModule() {
        $appName = $this->input->post('AppName', true);
        $setName = $this->input->post('SetName', true);
        $moduleName = trim(htmlspecialchars($this->input->post('ModuleName', true)));
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);
        if(!$appName || ! $setName || ! $moduleName || !$operator || !$bakOperator) {
            $obj = $this->config->item(CC_API_PARAMS_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('SetBaseLogic');
        $this->load->Logic('ModuleBaseLogic');
        $appInfo = $this->ApplicationBaseLogic->getAppList();
        $appId = 0;
        foreach($appInfo as $app) {
            if($app['ApplicationName'] == $appName) {
                $appId = $app['ApplicationID'];
            }
        }
        if(!$appId) {
            $obj = $this->config->item(CC_API_NOEST_APP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $setInfo = $this->SetBaseLogic->listSet($appId);
        $setId = 0;
        foreach($setInfo as $set) {
            if($set['SetName'] == $setName) {
                $setId = $set['SetID'];
            }
        }
        if(!$setId) {
            $obj = $this->config->item(CC_API_NOEST_SET);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $result = $this->ModuleBaseLogic->getModuleByName($moduleName, $setId, $appId);
        if(!empty($result)) {
            $obj = $this->config->item(CC_API_INV_DU_MOD);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->ModuleBaseLogic->addModule($appId, $setId, $moduleName, $operator, $bakOperator);
        return $this->outSuccess(array());
    }

    /**
     * 编辑模块
     */
    public function editModule() {
        $appName = $this->input->post('AppName', true);
        $setName = $this->input->post('SetName', true);
        $moduleName = trim(htmlspecialchars($this->input->post('ModuleName', true)));
        $newModuleName = trim(htmlspecialchars($this->input->post('newModuleName', true)));
        $operator = $this->input->post('Operator', true);
        $bakOperator = $this->input->post('BakOperator', true);
        if(!$appName || ! $setName || ! $moduleName || !$operator || !$bakOperator|| ! $newModuleName) {
            $obj = $this->config->item(CC_API_PARAMS_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('SetBaseLogic');
        $this->load->Logic('ModuleBaseLogic');
        $appInfo = $this->ApplicationBaseLogic->getAppList();
        $appId = 0;
        foreach($appInfo as $app) {
            if($app['ApplicationName'] == $appName) {
                $appId = $app['ApplicationID'];
            }
        }
        if(!$appId) {
            $obj = $this->config->item(CC_API_NOEST_APP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $setInfo = $this->SetBaseLogic->listSet($appId);
        $setId = 0;
        foreach($setInfo as $set) {
            if($set['SetName'] == $setName) {
                $setId = $set['SetID'];
            }
        }
        if(!$setId) {
            $obj = $this->config->item(CC_API_NOEST_SET);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $result = $this->ModuleBaseLogic->getModuleByName($moduleName, $setId, $appId);
        if(empty($result)) {
            $obj = $this->config->item(CC_API_INV_MODULE_NAME);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $moduleId = $result[0]['ModuleID'];
        if($newModuleName != $moduleName) {
            $checkResult = $this->ModuleBaseLogic->getModuleByName($newModuleName, $setId, $appId);
            if(!empty($checkResult)) {
                $obj = $this->config->item(CC_API_INV_DU_MOD);
                return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
            }
        }
        $this->ModuleBaseLogic->editModule($appId, $setId, $moduleId, $newModuleName, $operator, $bakOperator);
        return $this->outSuccess(array());
    }

    /**
     * 删除模块
     */
    public function delModule() {
        $appName = $this->input->post('AppName', true);
        $setName = $this->input->post('SetName', true);
        $moduleName = trim(htmlspecialchars($this->input->post('ModuleName', true)));
        if(!$appName || ! $setName || ! $moduleName ) {
            $obj = $this->config->item(CC_API_PARAMS_ILLEGAL);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('SetBaseLogic');
        $this->load->Logic('ModuleBaseLogic');
        $appInfo = $this->ApplicationBaseLogic->getAppList();
        $appId = 0;
        foreach($appInfo as $app) {
            if($app['ApplicationName'] == $appName) {
                $appId = $app['ApplicationID'];
            }
        }
        if(!$appId) {
            $obj = $this->config->item(CC_API_NOEST_APP);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $setInfo = $this->SetBaseLogic->listSet($appId);
        $setId = 0;
        foreach($setInfo as $set) {
            if($set['SetName'] == $setName) {
                $setId = $set['SetID'];
            }
        }
        if(!$setId) {
            $obj = $this->config->item(CC_API_NOEST_SET);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $result = $this->ModuleBaseLogic->getModuleByName($moduleName, $setId, $appId);
        if(empty($result)) {
            $obj = $this->config->item(CC_API_INV_MODULE_NAME);
            return $this->outFailure($obj->code, $obj->msg, $obj->extmsg);
        }
        $moduleId = $result[0]['ModuleID'];
        $this->ModuleBaseLogic->delModule($appId, $moduleId, $setId);
        return $this->outSuccess(array());
    }
}