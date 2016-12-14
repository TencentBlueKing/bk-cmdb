<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class SetBaseLogic extends Cc_Logic {


    private $_setEnviTypeArr = array();      //set环境类型参数

    private $_setServiceStatusArr = array(); //set服务状态参数

    private $_setIdArr = array();            //SetID数组

    private $_moduleNameArr = array();      //模块名数组

    public function __construct() {
        parent::__construct();
    }

    /**
     * 根据setId获取set信息
     * @return 数组
     */
    public function getSetById($setId = array(), $appId = array()) {
        $this->load->model('SetBaseModel');
        try {
            $set = $this->SetBaseModel->getSetById($setId, $appId);
            return $set;
        } catch (Exception $e) {
            CCLog::LogErr("getSetById exception" . $e->getMessage());
            $this->_errInfo = '根据Id获取set错误';
            return array();
        }
    }

    /**
     * 根据setId删除set
     * @return bool
     */
    public function delSetById($setId, $appId) {
        $this->load->model('ModuleBaseModel');
        $this->load->model('ModuleHostConfigModel');
        $this->load->model('HostBaseModel');
        $this->load->model('SetBaseModel');

        try {
            $modules = $this->ModuleBaseModel->getModuleById(array(), $setId, $appId);
            $moduleId = array();
            foreach ($modules as $_m) {
                if (!in_array($_m['ModuleID'], $moduleId)) {
                    $moduleId[] = $_m['ModuleID'];
                }
            }

            $moduleHostConfig = $this->ModuleHostConfigModel->getModuleHostConfigById(array(), $moduleId, $setId, $appId);
            $hostId = array();
            foreach ($moduleHostConfig as $_mhc) {
                if (!in_array($_mhc['HostID'], $hostId)) {
                    $hostId[] = $_mhc['HostID'];
                }
            }

            //删除主机资源
            $delHostRes = $this->HostBaseModel->deleteHostById($hostId, $moduleId, $setId, $appId);
            if (!$delHostRes) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
                return false;
            }

            //删除模块资源
            $delModuleRes = $this->ModuleBaseModel->deleteModuleById($moduleId, $setId, $appId);
            if (!$delModuleRes) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
                return false;
            }

            //删除集群
            $delSetRes = $this->SetBaseModel->delSetById($setId, $appId);
            if (!$delSetRes) {
                $this->_errInfo = $this->SetBaseModel->_errInfo;
                return false;
            }
            return true;
        } catch (Exception $e) {
            CCLog::LogErr("delSetById exception" . $e->getMessage());
            $this->_errInfo = 'delSetById exception';
            return array();
        }
    }

    /*
     * 新建集群
     * @param setId,appId
     * @return json
     */
    public function newSet($appId, $setName, $chnName, $enviType, $serviceStatus, $capacity, $des, $openStatus, &$setId = '') {
        try {
            $this->load->model('SetBaseModel');
            $nowTime = date('Y-m-d h:i:s');
            $data = array('ApplicationID' => $appId,
                          'SetName' => $setName,
                          'Default' => 0,
                          'ParentID' => 0,
                          'EnviType' => $enviType,
                          'ServiceStatus' => $serviceStatus,
                          'Capacity' => $capacity,
                          'ChnName' => $chnName,
                          'LastTime' => $nowTime,
                          'Description' => $des,
                          'Openstatus' => $openStatus,
                          'CreateTime' => $nowTime);
            $result = $this->SetBaseModel->addSet($data);
            $setId = $result;
            $errCode = $result ? '' : 'same_set_exist';
            $reBool = $result ? true : false;
            return $this->getOutput($reBool, $errCode);
        } catch (Exception $e) {
            CCLog::LogErr("new set exception:" . $e->getMessage());
            return $this->getOutput(false, 'new_set_error');
        }
    }

    /*
     * 修改集群
     * @param setId,appId
     * @return json
     */
    public function editSet($appId, $setId, $setName, $chnName, $envType, $serviceStatus, $capacity, $des, $openStatus) {
        $this->load->model('SetBaseModel');
        try {
            $nowTime = date("Y-m-d H:i:s");

            $data = array();
            $data = array('ApplicationID' => $appId,
                          'Default' => 0,
                          'ParentID' => 0,
                          'LastTime' => $nowTime);
            $joinData = array();

            $joinData['SetName'] = $setName;
            $joinData['EnviType'] = $envType;
            $joinData['ServiceStatus'] = $serviceStatus;
            $joinData['Capacity'] = $capacity;
            $joinData['ChnName'] = $chnName;
            $joinData['Description'] = $des;
            $joinData['Openstatus'] = $openStatus;

            foreach($joinData as $key=>$val) {
                if(!empty($val)) {
                    $data[$key] = $val;
                }
            }
            $data['ServiceStatus'] = $serviceStatus;
            $setIdArr = is_array($setId) ? $setId : array($setId);
            foreach ($setIdArr as $setId) {
                $setId = intval($setId);
                $result = $this->SetBaseModel->editSet($setId, $data, $appId);
            }

            $errCode = $result ? '' : 'edit_set_error';
            $reBool = $result ? true : false;
            return $this->getOutput($reBool, $errCode);
        } catch (Exception $e) {
            CCLog::LogErr("new set Err" . $e->getMessage());
            return $this->getOutput(false, 'edit_set_error');
        }
    }

    /*
     * 删除集群
     * @param setId,appId
     * @return json
     */
    public function delSet($appId, $setId) {
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');
        try {
            $setInfo = $this->SetBaseModel->getHostBySetID($setId);
            if ($setInfo) {
                return $this->getOutput(false, 'set_exsit_host');
            }

            $moduleIdInfo = $this->SetBaseModel->getModuleBySetId($setId, 'ModuleID');
            if ($moduleIdInfo) {
                $ModuleIDInArr = array_column($moduleIdInfo, 'ModuleID');
                $this->ModuleBaseModel->deleteModuleById(implode(',', $ModuleIDInArr), $setId, $appId);
            }
            $result = $this->SetBaseModel->delSetById($setId, $appId);
            if (!$result) {
                return $this->getOutput(false, 'delete_set_error');
            }
            return $this->getOutput(true);
        } catch (Exception $e) {
            CCLog::LogErr("delete set exception " . $e->getMessage());
            return $this->getOutput(false, 'edit_set_error');
        }
    }

    /*
     * 获取set列表
     * @param setId,appId
     * @return json
     */
    public function listSet($appId) {
        try {
            $this->load->model('SetBaseModel');
            $this->load->logic('BaseParameterDataLogic');
            $sets = $this->SetBaseModel->listSetNotDefault($appId);
            if (0 == count($sets)) {
                return array();
            }

            $setIdArr = array_column($sets, 'SetID');
            $moduleCountInSetId = $this->SetBaseModel->getModuleCountBySetId($setIdArr);

            $setEnviTypeArr = $this->BaseParameterDataLogic->getBaseParameterDataByDataType('SetEnviType');
            $SetServiceStatusArr = $this->BaseParameterDataLogic->getBaseParameterDataByDataType('SetServiceStatus');

            $ModuleCountKv = array();
            foreach ($moduleCountInSetId as $mc) {
                $moduleCountKv[$mc['SetID']] = $mc['cnt'];
            }

            foreach ($sets as &$set) {
                $set['EnviType'] = isset($set['EnviType']) && !empty($set['EnviType']) ? $set['EnviType'] : 3;
                $set['ServiceStatus'] = isset($set['ServiceStatus']) && $set['ServiceStatus'] != '' ? $set['ServiceStatus'] : 1;
                $set['Capacity'] = isset($set['Capacity']) && !empty($set['Capacity']) ? $set['Capacity'] : 0;
                $set['ModuleNum'] = isset($ModuleCountKv[$set['SetID']]) ? $ModuleCountKv[$set['SetID']] : 0;
            }
            return $sets;
        } catch (Exception $e) {
            CCLog::LogErr("list not default set exception" . $e->getMessage());
            return $this->getOutput(false, 'list_set_not_default');
        }
    }

    /*
     * 克隆集群
     * @param setId 集群Id
     * @param appId 业务Id
     * @param setNameArr 集群名数组
     * @return json
     */
    public function cloneSet($appId, $setId, $setNameArr) {
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');
        try {
            $set = $this->SetBaseModel->getSetById($setId, $appId);
            $moduleArr = $this->ModuleBaseModel->getModuleById(array(), $setId, $appId);
            $nowTime = date("Y-m-d H:i:s");
            $setInfo = current($set);
            $failArr = array();
            /**添加新的模块**/
            foreach ($setNameArr as $setName) {
                $setInfo['SetName'] = trim(htmlspecialchars($setName));
                $setInfo['LastTime'] = $nowTime;
                $setInfo['CreateTime'] = $nowTime;
                unset($setInfo['SetID']);
                $setId = $this->SetBaseModel->addSet($setInfo);
                if (!$setId) {
                    $failArr[] = $setName;
                    continue;
                }
                foreach ($moduleArr as $module) {
                    $module['SetID'] = $setId;
                    $module['LastTime'] = $nowTime;
                    $module['CreateTime'] = $nowTime;
                    unset($module['ModuleID']);
                    $this->ModuleBaseModel->addModule($module);
                }
            }
            return $failArr;
        } catch (Exception $e) {
            CCLog::LogErr("clone set exception" . $e->getMessage());
            return $this->getOutput(false, 'clone_set_error');
        }
    }

    /*
     * 根据set属性获取主机
     */
    public function getHostsBySetProperty($setEnviType, $setServiceStatus, $appId, $setId, $moduleName) {
        $this->load->model('SetPropertyModel');
        $this->load->model('ModuleBaseModel');
        $this->load->model('ModuleHostConfigModel');
        $this->load->model('HostBaseModel');
        try {
            $this->getQueryParams($setEnviType, $setServiceStatus, $setId, $moduleName);

            $setIdArr = $this->_setIdArr;
            $setServiceStatusArr = $this->_setServiceStatusArr;
            $setEnviTypeArr = $this->_setEnviTypeArr;
            $moduleNameArr = $this->_moduleNameArr;


            $setArr = $this->SetPropertyModel->getSetsByPropertyAndSetId($appId, $setIdArr, $setServiceStatusArr, $setEnviTypeArr);
            if(empty($setArr)) {
                return array();
            }
            $setIdArr = array_column($setArr, 'SetID');

            $moduleArr = $this->ModuleBaseModel->getModulesIdByAppIdAndModuleName($appId, $setIdArr, $moduleNameArr);

            if(empty($moduleArr)) {
                return array();
            }
            $moduleIdArr = array_column($moduleArr, 'ModuleID');
            $hostArr = $this->ModuleHostConfigModel->getHostIdInModuleIDArr($moduleIdArr);
            if(empty($hostArr)) {
                return array();
            }
            $hostIdArr = array_column($hostArr, 'HostID');

            $hostArr =  $this->HostBaseModel->getHostsByHostId($hostIdArr, 'InnerIP,OuterIP,Source,HostID');
            return $hostArr;
        } catch (Exception $e) {
            CCLog::LogErr("getHostsBySetPropertyexception" . $e->getMessage());
            return $this->getOutput(false, 'clone_set_error');
        }
    }

    /*
     * 获取set属性
     * @return set属性关联数组
     */
    public function getSetProperty() {
        $this->load->model('SetPropertyModel');
        try{
            $setProperty = $this->SetPropertyModel->getSetProperty();
            if(!$setProperty) {
                return array();
            }
            $setPropertyType = array_column($setProperty, 'PropertyType');
            $setPropertyType = array_unique($setPropertyType);
            $setProperyKv = array();

            foreach($setPropertyType as $sp) {
                $setProperyKv[$sp] = array();
            }

            foreach($setProperty as $sp) {
                if(isset($setProperyKv[$sp['PropertyType']])) {
                    $setProperyKv[$sp['PropertyType']][] = array('Property'=>$sp['PropertyCode'], 'value'=>$sp['PropertyName']);
                }
            }

            return $setProperyKv;
        } catch (Exception $e) {
            CCLog::LogErr("getHostsBySetProperty exception" . $e->getMessage());
            return $this->getOutput(false, 'clone_set_error');
        }

    }

    /*
     * 根据set属性获取set
     * @return json
     */
    public function getSetsBySetProperty($setEnviType, $setServiceStatus, $appId) {
        $this->load->model('SetPropertyModel');
        try{
            return $this->SetPropertyModel->getSetsByProperty($appId, $setServiceStatus, $setEnviType);
        } catch (Exception $e) {
            CCLog::LogErr("getSetsBySetProperty exception" . $e->getMessage());
            return $this->getOutput(false, 'clone_set_error');
        }
    }

    /*
     * 根据set属性获取模块
     * @return json
     */
    public function getModulesBySetProperty($setEnviType, $setServiceStatus, $appId, $setId) {
        $this->load->model('SetPropertyModel');
        $this->load->model('ModuleBaseModel');

        try{
            $this->getQueryParams($setEnviType, $setServiceStatus, $setId);

            $setIdArr = $this->_setIdArr;
            $setServiceStatusArr = $this->_setServiceStatusArr;
            $setEnviTypeArr = $this->_setEnviTypeArr;

            $setArr = $this->SetPropertyModel->getSetsByPropertyAndsetId($appId, $setIdArr, $setServiceStatusArr, $setEnviTypeArr);
            if(empty($setArr)) {
                return array();
            }

            $setIDArr = array_column($setArr, 'SetID');
            $modules = $this->ModuleBaseModel->getModulesNameByAppId($appId, $setIdArr);
            $moduleNameArr = array_column($modules,'ModuleName');
            return $moduleNameArr;
        } catch (Exception $e) {
            CCLog::LogErr("getSetsBySetProperty exception" . $e->getMessage());
            return $this->getOutput(false, 'clone_set_error');
        }
    }

    /**
     * 获取请求参数
     * @return json
     */
    private function getQueryParams($setEnviType, $setServiceStatus, $setId, $moduleName = '') {
        /*过滤set环境*/
        if(empty($SetEnviType)) {
            $setEnviTypeArr = array();
        } else {
            $setEnviTypeArr = explode(',', $setEnviType);
            foreach($setEnviTypeArr as &$sen) {
                $sen = intval($sen);
            }
        }
        $this->_setEnviTypeArr = $setEnviTypeArr;
        /*过滤服务状态*/
        if(empty($SetServiceStatus)) {
            $setServiceStatusArr = array();
        } else {
            $setServiceStatusArr = explode(',', $SetServiceStatus);
            foreach($setServiceStatusArr as &$ser) {
                $ser = intval($ser);
            }
        }
        $this->_setServiceStatusArr = $setServiceStatusArr;

        /*过滤setId*/
        if( empty($setId)) {
            $setIdArr = array();
        } else {
            $setIdArr = explode(',', $setId);
            foreach($setIdArr as &$sid) {
                $sid = intval($sid);
            }
        }
        $this->_setIdArr = $setIdArr;
        /*过滤模块名*/
        if( empty($moduleName)) {
            $moduleNameArr = array();
        } else {
            $moduleNameArr = explode(',', $moduleName);
        }
        $this->_moduleNameArr = $moduleNameArr;
    }
}