<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class  ModuleHostConfigModel extends Cc_Model {
    
    public function _construct() {
        parent::_construct();
    }

    /*
     * @根据业务获取主机
     * @param appId
     * @return array()
     */
    public function getHostByAppID($appId) {
        $this->db->select('HostID');
        $this->db->from('ModuleHostConfig');
        $this->db->where('ApplicationID', $appId);
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

    /**
    * @todo 添加主机模块关系
    * @param $hostId int 主机Id
    * @param $moduleId int 模块Id
    * @param $setId int groupId
    * @param $appId int 业务Id
    * @return int or boolean
    */
    public function addModuleHostConfig($hostId, $moduleId, $setId, $appId){
        if(($type = 'hostId' && !$hostId) || ($type = 'moduleId' && !$moduleId) || ($type = 'setId' && !$setId) || ($type='appId' && !$appId)){
            $this->_errInfo = $type.'非法';
            return false;
        }

        $data = array();
        $data['HostID'] = $hostId;
        $data['ModuleID'] = $moduleId;
        $data['SetID'] = $setId;
        $data['ApplicationID'] = $appId;

        $result = $this->db->insert('ModuleHostConfig', $data);
        if(!$result){
            $this->errInfo = '添加主机模块关系,失败';
            $err = $this->db->error();
            CCLog::LogErr($this->_errInfo.', mysql_errno: '. $err['code'] .',mysql_error: '. $err['message']);
            return false;
        }

        return $this->db->insert_id();
    }

    /**
     * @id查询主机
     * @param $hostId 主机Id
     * @param $moduleId 集群Id
     * @param $setId Group自增Id
     * @param $appId 业务Id
     * @return 主机id、集群id、groupid、业务id映射关系
     */
    public function getModuleHostConfigById($hostId=array(), $moduleId=array(), $setId=array(), $appId=array()) {
        $hostIdArr     = !is_array($hostId) ? explode(',', $hostId) : array(intval($hostId));
        $moduleIdArr   = !is_array($moduleId) ? explode(',', $moduleId) : array(intval($moduleId));
        $setIdArr      = !is_array($setId) ? explode(',', $setId) : array(intval($setId));
        $appIdArr      = !is_array($appId) ? explode(',', $appId) : array(intval($appId));
        
        $hostIdArr != array('') && $this->db->where_in('HostID', $hostIdArr);
        $moduleIdArr != array('') && $this->db->where_in('ModuleID', $moduleIdArr);
        $setIdArr != array('') && $this->db->where_in('SetID', $setIdArr);
        $appIdArr != array('') && $this->db->where_in('ApplicationID', $appIdArr);
        $query = $this->db->get('ModuleHostConfig');
        return $query->num_rows() ?  $query->result_array() : array();
    }

    /**
     * @id查询主机
     * @param $moduleId 模块Id
     * @return hostIdArr
     */
    public function getHostIdByModuleId($moduleId) {
        $this->db->select('HostID');
        $this->db->from('ModuleHostConfig');
        $this->db->where('ModuleID', $moduleId);
        $query = $this->db->get();
        return $query->num_rows() ?  $query->result_array() : array();
    }

    /**
     * @todo 统计业务下的主机
     */
    public function statHostCountByApp($appIdArr, $emptyModuleIdArr) {
        !is_array($appIdArr) && $appIdArr = array($appIdArr);
        !is_array($emptyModuleIdArr) && $emptyModuleIdArr = array($emptyModuleIdArr);
        $this->db->select('count(distinct(HostID)) as cnt, ApplicationID');
        $this->db->from('cc_ModuleHostConfig');
        $this->db->where_in('ApplicationID', $appIdArr);
        $this->db->where_not_in('ModuleID', $emptyModuleIdArr);
        $this->db->group_by('ApplicationID');
        $query = $this->db->get();

        if(0 == $query->num_rows()) {
            return array();
        }
        $resultArr = array();
        $result = $query->result_array();
        foreach($result as $re) {
            $resultArr[$re['ApplicationID']] = $re['cnt'];
        }
        return $resultArr;
    }

    /**
     * @统计setId下的主机
     * @param $appID
     */
    public function statHostCountBySetID($appId) {
        $this->db->select('count(distinct(HostID)) as cnt,SetID');
        $this->db->from('cc_ModuleHostConfig');
        $this->db->where('ApplicationID', $appId);
        $this->db->group_by('SetID');
        $query = $this->db->get();
        if(0 == $query->num_rows()) {
            return array();
        }
        $resultArr = array();
        $result = $query->result_array();
        foreach($result as $re) {
            $resultArr[$re['SetID']] = $re['cnt'];
        }
        return $resultArr ;
    }

    /**
     * @统计ModuleID下的主机
     * @param  $appID
     */
    public function statHostCountByModuleID($appId) {
        $this->db->select('count(HostID) as cnt,ModuleID');
        $this->db->from('cc_ModuleHostConfig');
        $this->db->where('ApplicationID', $appId);
        $this->db->group_by('ModuleID');
        $query = $this->db->get();
        if(0 == $query->num_rows()) {
            return array();
        }
        $resultArr = array();
        $result = $query->result_array();
        foreach($result as $re) {
            $resultArr[$re['ModuleID']] = $re['cnt'];
        }
        return $resultArr;
    }

    /**
     * @统计指定业务、集群、模块下的主机
     */
    public function getHostCountById($moduleId, $setId, $appId) {
        $this->db->select('count( distinct(HostID)) as cnt');
        $this->db->from('cc_ModuleHostConfig');

        if(!empty($moduleId)){
            $this->db->where('ModuleID', $moduleId);
        }

        if(!empty($setId)){
            $this->db->where('SetID', $setId);
        }

        $query = $this->db->get();
        if(0 == $query->num_rows) {
            return 0;
        }
        $result = $query->result_array();
        return $result[0]['cnt'];
    }

    /**
     * @根据ip和业务id获取主机
     */
    public function getHostsByIpAndAppId($ip, $appId, $fields='ModuleHostConfig.ApplicationID,ModuleHostConfig.SetID,ModuleHostConfig.ModuleID,ModuleHostConfig.HostID,AssetID,HostName,DeviceClass,Operator,BakOperator,InnerIP,OuterIP,Status,CreateTime,HardMemo,Region,OSName,IdcName') {
        $this->db->select($fields);
        $this->db->from('ModuleHostConfig');
        $this->db->join('HostBase', 'ModuleHostConfig.HostID=HostBase.HostID', 'INNER');
        $this->db->where('ModuleHostConfig.ApplicationID', $appId);

        if(! empty($ip)) {
            $this->db->group_start('', ' AND ');
            $this->db->where_in('InnerIP', $ip);
            $this->db->or_where_in('OuterIP', $ip);
            $this->db->group_end();
        }

        $query = $this->db->get();
        return $query ? $query->result_array() : array();
    }

    /**
     * @根据集群Id和业务Id查询主机
     */
    public function getHostsBySetIdAndAppId($setIdArr, $appId, $fields='ModuleHostConfig.ApplicationID,ModuleHostConfig.SetID,ModuleHostConfig.ModuleID,ModuleHostConfig.HostID,AssetID,HostName,DeviceClass,Operator,BakOperator,InnerIP,OuterIP,Status,CreateTime,Mem,HardMemo,OSName,IdcName,Region') {
        $this->db->select($fields);
        $this->db->from('ModuleHostConfig');
        $this->db->join('HostBase', 'ModuleHostConfig.HostID=HostBase.HostID', 'INNER');
        $this->db->where_in('SetID', $setIdArr);
        $this->db->where('ModuleHostConfig.ApplicationID', $appId);
        $query = $this->db->get();
        return $query ? $query->result_array() : array();
    }

    /**
     * @根据模块Id和业务Id查询主机
     */
    public function getHostsByModuleIdAndAppId($moduleId, $appId, $fields='ModuleHostConfig.ApplicationID,ModuleHostConfig.SetID,ModuleHostConfig.ModuleID,ModuleHostConfig.HostID,AssetID,HostName,DeviceClass,Operator,BakOperator,InnerIP,OuterIP,Status,CreateTime,Mem,HardMemo,Source,OSName,IdcName,Region') {
        $this->db->select($fields);
        $this->db->from('ModuleHostConfig');
        $this->db->join('HostBase', 'ModuleHostConfig.HostID=HostBase.HostID', 'INNER');
        $this->db->where_in('ModuleID', $moduleId);
        $this->db->where('ModuleHostConfig.ApplicationID', $appId);
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

    /**
     * @获取模块下所有的主机Id
     */
    public function getHostIdInModuleIdArr($moduleIdArr, $fields=null) {
        if(empty($moduleIdArr)){
            return array();
        }

        if(is_null($fields)){
            $this->db->select('HostID');
        }else{
            $this->db->select($fields);
        }
        $this->db->where_in('ModuleID', $moduleIdArr);
        $this->db->distinct();
        $query = $this->db->get('ModuleHostConfig');

        return $query->num_rows() ? $query->result_array() : array();
    }

}