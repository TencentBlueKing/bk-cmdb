<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class HostBaseLogic extends Cc_Logic {
    
    public function __construct(){
        parent::__construct();
    }

    /**
     * @id查询主机
     * @param hostId 主机ID
     * @param moduleId 模块ID
     * @param setId 大区ID
     * @param appId 业务ID
     * @param start 起始下标
     * @param limit 欲取主机行数
     * @param orderBy 排序字段
     * @param direction 排序规则，desc/asc
     * @return array() 主机信息数组
     */
    public function getHostById($hostId = array(), $moduleId = array(), $setId = array(), $appId = array(), $start = '', $limit = '', $orderBy = '', $direction = 'DESC') {
        try{
            $this->load->model('HostBaseModel');
            $host = $this->HostBaseModel->getHostById($hostId, $moduleId, $setId, $appId, $start, $limit, $orderBy, $direction);
            return $host;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic->getHostById:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @id查询主机数
     * @param hostId 主机ID
     * @param moduleId 模块ID
     * @param setId 大区ID
     * @param appId 业务ID
     * @return int 主机数
     */
    public function getHostCountById($hostId = array(), $moduleId = array(), $setId = array(), $appId = array()) {
        try{
            $this->load->model('HostBaseModel');
            $total = $this->HostBaseModel->getHostCountById($hostId, $moduleId, $setId, $appId);
            return $total;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic->getHostCountById:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @ip查询主机
     * @param hostIp 主机ip
     * @param appId 业务ID
     * @param start 起始下标
     * @param limit 欲取主机行数
     * @param orderBy 排序字段
     * @param direction 排序规则，desc/asc
     * @return array() 主机信息数组
     */
    public function getHostByIp($hostIp = array(), $moduleId = array(), $setId = array(), $appId = array(), $start = '', $limit = '', $orderBy = '', $direction = 'DESC') {
        try{
            $this->load->model('HostBaseModel');
            $host = $this->HostBaseModel->getHostByIp($hostIp, $moduleId, $setId, $appId, $start, $limit, $orderBy, $direction);
            return $host;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic->getHostByIp:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @多条件查询主机
     * @return array() 主机信息数组
     */
    public function getHostByCondition($condition, $distinct = true, $start = 0, $limit = 0) {
        try{
            $this->load->model('HostBaseModel');
            $allowedFields = $this->HostBaseModel->getAllowedFields();
            $allowedFields = array_merge($allowedFields, array('ApplicationID', 'SetID', 'ModuleID', 'IfInnerIPexact', 'IfOuterexact'));
            $diff = array_diff(array_keys($condition), $allowedFields);
            if($diff) {
                $this->_errInfo = '不支持查询条件['. implode(',', $diff) .']';
                return array();
            }

            $host = $this->HostBaseModel->getHostByCondition($condition, $distinct, $start, $limit);
            return $host;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic->getHostByCondition:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @多条件查询主机数
     * @param condition array 查询条件
     * @param distinct boolean 是否去重
     * @return int 主机数
     */
    public function getHostCountByCondition($condition, $distinct = true) {
        try{
            $this->load->model('HostBaseModel');
            $allowedFields = $this->HostBaseModel->getAllowedFields();
            $allowedFields = array_merge($allowedFields, array('ApplicationID', 'SetID', 'ModuleID', 'IfInnerIPexact', 'IfOuterexact'));
            $diff = array_diff(array_keys($condition), $allowedFields);
            if($diff) {
                $this->_errInfo = '不支持查询条件['. implode(',', $diff) .']';
                return array();
            }

            $total = $this->HostBaseModel->getHostCountByCondition($condition, $distinct);
            return $total;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic->getHostCountByCondition:' . $e->getMessage());
            return array();
        }
    }

    /**
     * @更新主机
     * @return array() 主机信息数组
     */
    public function updateHost() {
        try{
            $appId = $this->input->post('ApplicationID');
            if(!$appId) {
                $this->_errInfo = '业务ID非法';
            }

            $hostId = $this->input->post('HostID');
            if(!$hostId) {
                $this->_errInfo = '主机ID非法';
            }

            $stdProperty = $this->input->post('stdProperty', true);
            $stdProperty == NULL ? array() : (is_array($stdProperty) ? $stdProperty : json_decode($stdProperty));

            $cusProperty = $this->input->post('cusProperty', true);
            $cusProperty == NULL ? array() : (is_array($cusProperty) ? $cusProperty : json_decode($cusProperty));

            if(empty($stdProperty) && empty($cusProperty)) {
                $this->_errInfo = '请至少修改一个属性';
                return false;
            }

            $stdAllowedFields = array('HostName', 'BakOperator', 'Operator', 'Description', 'Source');
            $diff = array_diff(array_keys($stdProperty), $stdAllowedFields);
            if($diff) {
                $this->_errInfo = '不支持修改字段['. implode(',', $diff) .']';
                return false;
            }

            if(isset($stdProperty['Source'])) {
                $this->load->logic('BaseParameterDataLogic');
                $parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();

                if(!isset($parameterData[$stdProperty['Source']])) {
                    $this->_errInfo = '云供应商非法，必需为['. implode('、', array_keys($parameterData)) .']其中之一';
                    return false;
                }
            }

            if(!empty($stdProperty)) {
                $this->load->model('HostBaseModel');
                $upStd = $this->HostBaseModel->updateHostById($stdProperty, $hostId, $appId);
                if(!$upStd) {
                    $this->_errInfo = $this->HostBaseModel->_errInfo;
                    return false;
                }
            }

            return true;
        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
    * @转移主机所属模块
    * @return boolean 成功 or 失败
    */
    public function modHostModule($hostId, $moduleId, $appId, $isIncrement) {
        try{
            $this->load->model('HostBaseModel');
            $moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
            if(count($moduleId)>1) {
                $hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
                $result = true;
                foreach($hostId as $_h) {
                    $result = $result && $this->HostBaseModel->modSingleHostToMultiModule($_h, $moduleId, $appId, $isIncrement);
                    if(!$result) {
                        $this->_errInfo = $this->HostBaseModel->_errInfo;
                        return false;
                    }
                }

                return true;
            } else {
                $moduleId = $moduleId[0];
                $result = $this->HostBaseModel->modMultiHostToSingleModule($hostId, $moduleId, $appId, $isIncrement);
                if(!$result) {
                    $this->_errInfo = $this->HostBaseModel->_errInfo;
                }

                return $result;
            }

        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
    * @主机操作-快速分配页面，分配主机资源
    * @param hostId 主机Id
    * @param moduleId 业务的空闲机模块Id
    * @param appId 业务Id
    * @return boolean 分配成功or失败
    */
    public function quickDistribute($hostId, $moduleId, $appId) {
        try{
            $this->load->model('HostBaseModel');
            $moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
            $moduleId = $moduleId[0];
            $result = $this->HostBaseModel->quickDistribute($hostId, $moduleId, $appId);
            if(!$result) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
            }

            return $result;

        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
    * @转移主机所属模块
    * @return boolean 成功 or 失败
    */
    public function deleteHost($hostId, $moduleId, $setId, $appId) {
        try{
            $this->load->model('HostBaseModel');
            $result = $this->HostBaseModel->deleteHostById($hostId, $moduleId, $setId, $appId);

            if(!$result) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
            }

            return $result;
        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
    * @上缴主机
    * @return boolean 成功 or 失败
    */
    public function resHostModule($hostId, $moduleId, $setId, $appId) {
        try{
            $this->load->model('HostBaseModel');
            $result = $this->HostBaseModel->resHostModule($hostId, $moduleId, $setId, $appId);

            if(!$result) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
            }

            return $result;
        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
    * @根据业务id查询主机数量
    * @return int
    */
    public function getHostCountByApplicationIDs($ApplicationIDs) {
        $this->load->model('HostBaseModel');
        return $this->HostBaseModel->getHostCountById(array(), array(), array(), $ApplicationIDs);
    }

    /**
     * @根据主机id查询主机相关的业务集群模块
     * @return int
     */
    public function getHostAppGroupModuleRealtionByHostIDs($hostIDs) {
        $this->load->model('ModuleHostConfigModel');
        return $this->ModuleHostConfigModel->getModuleHostConfigById($hostIDs);
    }

    /**
     * @删除主机
     * @return boolean 成功 or 失败
     */
    public function RealDeleteHost($hostId, $appID) {
        try{
            $this->load->model('HostBaseModel');
            $result = $this->HostBaseModel->deleteHostApplicationById($hostId, $appID);

            if(!$result) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
            }

            return $result;
        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
     * @获取主机某些字段下拉值
     * @return array
     */
    public function getFieldSelectByType($type) {
        try{
            $this->load->model('HostBaseModel');
            $result = $this->HostBaseModel->getFieldSelectByType($type);

            if(!$result) {
                $this->_errInfo = $this->HostBaseModel->_errInfo;
            }

            return $result;
        }catch(Exception $e) {
            CCLog::LogErr($e->getMessage());
            return false;
        }
    }

    /**
     * @获取主机属性值
     * @param $type
     */
    public function getHostPropertyByType($type) {
        try{
            $data = array();
            $this->load->model('HostPropertyClassifyModel');
            $hostProperty = $this->HostPropertyClassifyModel->getAllHostProperty();
            foreach ($hostProperty as $key => $value) {
                if($type == 'keyName') {
                    $data[$value['HostTableField']] = $value['PropertyName'];
                }elseif($type == 'groupKey') {
                    $data[$value['Group']][] = $value['HostTableField'];
                }elseif($type ==  'nameKey') {
                    $data[$value['PropertyName']] = $value['HostTableField'];
                }
            }

            return $data;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic:getHostPropertyByType->' . $e->getMessage());

            return array();
        }
    }

    /**
     * @获取主机属性值
     * @param $owner
     * @param $type
     */
    public function getHostPropertyByOwner($owner, $type) {
        try{
            $data = array();
            $this->load->model('HostCustomerPropertyModel');
            $hostProperty = $this->HostCustomerPropertyModel->getHostPropertyByOwner($owner);
            foreach ($hostProperty as $key => $value) {
                if($type == 'keyName') {
                    $data[$value['HostTableField']] = $value['PropertyName'];
                }elseif($type == 'groupKey') {
                    $data[$value['Group']][] = $value['HostTableField'];
                }elseif($type ==  'nameKey') {
                    $data[$value['PropertyName']] = $value['HostTableField'];
                }
            }

            return $data;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic:getHostPropertyByOwner->' . $e->getMessage());

            return array();
        }
    }

    /**
     * @新增主机入库
     * @param $innerIP
     * @param $outerIP
     * @param $hostName
     * @param $operator
     * @param $bakOperator
     */
    public function addHost($innerIP, $outerIP, $hostName, $operator, $bakOperator) {
        try{
            $this->load->model('HostBaseModel');
            $this->load->model('ApplicationBaseModel');
            $this->load->model('ModuleBaseModel');
            $this->load->model('ModuleHostConfigModel');
            $result = $this->HostBaseModel->getHostByIp($innerIP);
            if(!empty($result)) {
                return CC_API_DUP_INPUT_IP;
            }
            $defaultApp = $this->ApplicationBaseModel->getResPoolByCompany(COMPANY_NAME);
            if(empty($defaultApp)) {
                return CC_API_ADD_IP_FAIL;
            }
            $appId = $defaultApp['ApplicationID'];
            $moduleInfo = $this->ModuleBaseModel->getResPoolIDByAppID($appId);
            if(empty($moduleInfo)) {
                return CC_API_ADD_IP_FAIL;
            }
            $setId = $moduleInfo['SetID'];
            $moduleId = $moduleInfo['ModuleID'];
            $nowTime = date('Y-m-d H:i:s');
            $data = array();
            $data['InnerIP'] = $innerIP;
            if(!empty($outerIP)) {
                $data['OuterIP'] = $outerIP;
            }
            if(!empty($hostName)) {
                $data['HostName'] = $hostName;
            }
            if(!empty($operator)) {
                $data['Operator'] = $operator;
            }
            if(!empty($operator)) {
                $data['Operator'] = $operator;
            }
            if(!empty($bakOperator)) {
                $data['BakOperator'] = $bakOperator;
            }
            $data['CreateTime'] = $nowTime;
            $data['LastTime'] = $nowTime;
            $hostId = $this->HostBaseModel->AddHost($data);
            if(!$hostId) {
                return CC_API_ADD_IP_FAIL;
            }
            $result = $this->ModuleHostConfigModel->addModuleHostConfig($hostId, $moduleId, $setId, $appId);
            return true;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic: add Host Err ->' . $e->getMessage());
            return array();
        }
    }

    /**
     * @新增主机入库
     * @param $innerIP
     * @param $outerIP
     * @param $hostName
     * @param $operator
     * @param $bakOperator
     */
    public function editHost($innerIP, $outerIP, $hostName, $operator, $bakOperator) {
        try{
            $this->load->model('HostBaseModel');
            $result = $this->HostBaseModel->getHostByIp($innerIP);
            if(empty($result)) {
                return CC_API_GET_IP_FAIL;
            }
            $hostId = $result[0]['HostID'];
            $nowTime = date('Y-m-d H:i:s');
            $data = array();
            $data['InnerIP'] = $innerIP;
            if(!empty($outerIP)) {
                $data['OuterIP'] = $outerIP;
            }
            if(!empty($hostName)) {
                $data['HostName'] = $hostName;
            }
            if(!empty($operator)) {
                $data['Operator'] = $operator;
            }
            if(!empty($operator)) {
                $data['Operator'] = $operator;
            }
            if(!empty($bakOperator)) {
                $data['BakOperator'] = $bakOperator;
            }
            if(empty($data)) {
                return true;
            }
            $data['LastTime'] = $nowTime;
            $result = $this->HostBaseModel->updateHostById($data, $hostId);
            return true;
        }catch(Exception $e) {
            CCLog::LogErr('HostBaseLogic: edit Host Err ->' . $e->getMessage());
            return array();
        }
    }
}