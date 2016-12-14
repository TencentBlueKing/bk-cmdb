<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ModuleBaseLogic extends Cc_Logic {
    public function __construct() {
        parent::__construct();
    }

    /**
     * id查询module
     * @param $moduleId 模块Id
     * @param $setId    大区Id
     * @param $appId    业务Id
     * @return array() module信息数组
     */
    public function getModuleById($moduleId = array(), $setId = array(), $appId = array()) {
        $this->load->model('ModuleBaseModel');
        try {
            $module = $this->ModuleBaseModel->getModuleById($moduleId, $setId, $appId);

            return $module;
        } catch (Exception $e) {
            CCLog::LogErr("getModuleById exception".$e->getMessage());
            $this->_errInfo = '根据Id获取模块异常';
            return array();
        }
    }

    /**
     * 模块名查询module
     * @param moduleName 模块名
     * @param setId      大区ID
     * @param appId      业务ID
     * @return array() module信息数组
     */
    public function getModuleByName($moduleName = array(), $setId = array(), $appId = array()) {
        $this->load->model('ModuleBaseModel');
        try {
            $module = $this->ModuleBaseModel->getModuleByName($moduleName, $setId, $appId);

            return $module;
        } catch (Exception $e) {
            CCLog::LogErr('get Module by Name exception'.$e->getMessage());
            $this->_errInfo = '根据模块名获取模块异常';
            return array();
        }
    }

    /**
     * 新增模块
     */
    public function addModule($appId, $setId, $moduleName, $operator, $bakOperator) {
        $this->load->model('ModuleBaseModel');
        try {
            $datetime = date("Y-m-d H:i:s");
            $data = array('ApplicationID' => $appId,
                          'SetID' => $setId,
                          'ModuleName' => $moduleName,
                          'Default' => 0,
                          'CreateTime' => $datetime,
                          'LastTime' => $datetime,
                          'Operator' => $operator,
                          'BakOperator' => $bakOperator);
            $result = $this->ModuleBaseModel->addModule($data);
            if ($result) {
                return $this->getOutput(true);
            } else {
                return $this->getOutput(false, 'same_module_name');
            }
        } catch (Exception $e) {
            CCLog::LogErr('add module error:' . $e->getMessage());
            return $this->getOutput(false, 'new_module_error');
        }
    }

    /**
     * 修改模块属性
     * @param $moduleName 模块名
     * @param $setId      大区ID
     * @param $appId      业务ID
     */
    public function editModule($appId, $setId, $moduleId, $moduleName, $operator, $bakOperator) {
        $this->load->model('ModuleBaseModel');
        try {
            $moduleIdArr = is_array($moduleId) ? $moduleId : array($moduleId);
            foreach ($moduleIdArr as $moduleId) {
                $moduleId = intval($moduleId);
                $result = $this->ModuleBaseModel->editModule($appId, $setId, $moduleId, $moduleName, $operator, $bakOperator);
            }
            $errCode = isset($result['errorcode']) ? $result['errorcode'] : '';
            return $this->getOutput($result['success'], $errCode);
        } catch (Exception $e) {
            CCLog::LogErr('edit module error:' . $e->getMessage());
            return $this->getOutput(false, 'edit_module_error');
        }
    }

    /**
     * 删除模块
     */
    public function delModule($appId, $moduleId, $setId) {
        $this->load->model('ModuleBaseModel');
        $this->load->model('ModuleHostConfigModel');
        try {
            $host = $this->ModuleHostConfigModel->getHostIdByModuleId($moduleId);
            if ($host) {
                return $this->getOutput(false, 'module_exsit_host');
            }
            $result = $this->ModuleBaseModel->deleteModuleById($moduleId, $setId, $appId);
            return $this->getOutput(true);
        } catch (Exception $e) {
            CCLog::LogErr('del module error:' . $e->getMessage());
            return $this->getOutput(false, 'delete_module_fail');
        }
    }

    /**
     * 查询模块列表
     * @return 模块列表
     */
    public function listModule($appId) {
        $this->load->model('ModuleBaseModel');
        $this->load->model('ModuleHostConfigModel');
        $this->load->model('SetBaseModel');
        try {
            $modules = $this->ModuleBaseModel->listModuleNotDefault($appId);
            if (empty($modules)) {
                return array();
            }
            $moduleHostCount = $this->ModuleHostConfigModel->StatHostCountByModuleID($appId);

            $setArr = $this->SetBaseModel->getSetById(array(), $appId);
            $setKv = array();
            foreach ($setArr as $set) {
                $setKv[$set['SetID']] = $set['SetName'];
            }
            foreach ($modules as &$mhc) {
                $mhc['SetName'] = isset($setKv[$mhc['SetID']]) ? $setKv[$mhc['SetID']] : '';
                $mhc['HostCount'] = isset($moduleHostCount[$mhc['ModuleID']]) ? $moduleHostCount[$mhc['ModuleID']] : 0;
            }
            return $modules;
        } catch (Exception $e) {
            CCLog::LogErr("list not default module Err" . $e->getMessage());
            return $this->getOutput(false, 'list_module_not_default');
        }
    }

    /**
     * 根据业务Id查询所有模块
     * @param $appName 业务名
     * @param $company 公司Id
     * @return 数组
     */
    public function getModulesByAppId($appId){
        $this->load->model('ModuleBaseModel');
        try {
            $result = $this->ModuleBaseModel->getModulesNameByAppID($appId);
            return $result;
        } catch (Exception $e) {
            CCLog::LogErr("getModuleByAppId exception" . $e->getMessage());
            return array();
        }
    }
}