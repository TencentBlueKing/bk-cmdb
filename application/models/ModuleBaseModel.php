<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ModuleBaseModel extends Cc_Model {

	public function __construct(){
		parent::__construct();
	}

	/*
     * @数据库中通过Id查询模块信息
     * @param setId,appId,moduleId
     * @return 数组
     */
	public function getModuleById($moduleId=array(), $setId=array(), $appId=array()){
		$moduleIdArr = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$setIdArr = is_array($setId) ? $setId : explode(',', $setId);
		$appIdArr = is_array($appId) ? $appId : explode(',', $appId);

		if($moduleId!=array('') && !empty($moduleId)) {
			$this->db->where_in('ModuleID', $moduleIdArr);
		}

		if($setId!=array('') && !empty($setId)) {
			$this->db->where_in('SetID', $setIdArr);
		}

		if($appId!=array('') && !empty($appId)) {
			$this->db->where_in('ApplicationID', $appIdArr);
		}

		$query = $this->db->get('cc_ModuleBase');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/*
     * @通过模块名获取模块
     * @param setId,appId,moduleName
     * @return 数组
     */
	public function getModuleByName($moduleName=array(), $setId=array(), $appId=array()) {
		$moduleNameArr = is_array($moduleName) ? $moduleName : explode(',', $moduleName);
		$setIdArr = is_array($setId) ? $setId : explode(',', $setId);
		$appIdArr = is_array($appId) ? $appId : explode(',', $appId);

		if(empty($appId) || $appId==array('')) {
			return array();
		}

		if($moduleName!=array('')){
			$this->db->where_in('ModuleName', $moduleNameArr);
		}

		if($setId!=array('') && !empty($setId)){
			$this->db->where_in('SetID', $setIdArr);
		}

		$this->db->where_in('ApplicationID', $appIdArr);

		$query = $this->db->get('cc_ModuleBase');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/*
     * @添加模块
     * @param data
     * @return bool
     */
	public function addModule($data) {
		if(!isset($data['ApplicationID']) || !isset($data['SetID']) || !isset($data['ModuleName'])){
			$this->_errInfo = 'ApplicationID/SetID/ModuleName缺一不可';
			return false;
		}

		$this->db->select('ModuleID');
		$this->db->where('ModuleName', $data['ModuleName']);
		$this->db->where('SetID', $data['SetID']);
		$query = $this->db->get('cc_ModuleBase');

		if($query && $query->num_rows() == 0){
			$query->free_result();

			$result = $this->db->insert('cc_ModuleBase', $data);
			if(!$result){
				$this->_errInfo = '添加模块失败!';
				$err = $this->db->error();
				CCLog::LogErr('添加模块失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
                return false;
			}

            $log['ApplicationID'] = $data['ApplicationID'];
            $log['OpContent'] = '模块名:['.$data['ModuleName'].']';
            $log['OpTarget'] = '模块';
            $log['OpType'] = '新增';
            CCLog::addOpLogArr($log);
			return $this->db->insert_id();
		}

		$this->_errInfo = '同名模块已存在[ModuleName='. $data['ModuleName'] .']!';
		return false;
	}

	/*
     * @通过主机Id查询模块
     * @param hostId数组
     * @return 模块数组
     */
	public function getModuleByHostId($hostId=array()) {
		$hostIdArr = is_array($hostId) ? $hostId : explode(',', $hostId);

		if($hostIdArr == array('')) {
			return array();
		}

		$this->db->select('cc_ModuleBase.*, cc_ModuleHostConfig.HostID');
		$this->db->from('cc_ModuleBase');
		$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.ModuleID=cc_ModuleBase.ModuleID');
		$this->db->where_in('cc_ModuleHostConfig.HostID', $hostIdArr);
		$query = $this->db->get();

		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/*
     * @删除模块
     * @param moduleId,setId,appId
     * @return bool
     */
	public function deleteModuleById($moduleId=array(), $setId=array(), $appId=0) {
		$moduleIdArr 	= !is_array($moduleId) ? explode(',', $moduleId) : array(intval($moduleId));
    	$setIdArr 		= !is_array($setId) ? explode(',', $setId) : array(intval($setId));
    	$appId 		= intval($appId);

		$this->db->where_in('ModuleID', $moduleIdArr);
		$setIdArr != array('') && !empty($setId) && $this->db->where_in('SetID', $setIdArr);
		$appId != array('') && !empty($appId) && $this->db->where('ApplicationID', $appId);
		$query = $this->db->get('cc_ModuleBase');
		if(!$query || $query->num_rows() == 0){
			$this->_errInfo = '模块[id='. implode(',', $moduleIdArr) .']不存在';
			return false;
		}
		$module = $query->result_array();

		$this->db->where_in('ModuleID', $moduleIdArr);
		$setId != array('') && !empty($setId) &&  $this->db->where_in('SetID', $setId);
		$appId != array('') && !empty($appId) && $this->db->where('ApplicationID', $appId);
		$result = $this->db->delete('cc_ModuleBase');

		if(!$result) {
			$this->_errInfo = '删除模块失败!';
			$err = $this->db->error();
			CCLog::LogErr('删除模块失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
		}

		$log = array();
		$log['ApplicationID'] = $appId;
		$log['OpContent'] = '模块名['. implode(',', array_column($module, 'ModuleName')) .']';
		$log['OpFrom'] = 0;
		$log['OpName'] = '删除模块';
		$log['OpResult'] = $result ? 1 : 0;
		$log['OpTarget'] = '模块';
		$log['OpType'] = '删除';
		CCLog::addOpLogArr($log);
		return false;
	}

	/*
     * @编辑模块
     * @param moduleId,setId,appId
     * @return bool
     */
    public function editModule($appId, $setId, $moduleId, $moduleName, $operator, $bakOperator) {
        $this->db->select('ModuleID');
        $this->db->from('cc_ModuleBase');
        if(isset($SetID)) {
            $this->db->where('SetID',$SetID);
        }
        $this->db->where('ModuleID <> ', $moduleId);
        $this->db->where('ModuleName', $moduleName);
		$this->db->where('SetID', $setId);
        $query = $this->db->get();
        if(0 != $query->num_rows()) {
            return array('success'=>false, 'errorcode'=>'same_module_name');
        }
        $this->db->select('ModuleName, Operator, BakOperator');
        $this->db->from('cc_ModuleBase');
        $this->db->where('ModuleID', $moduleId);
        $query = current($this->db->get()->result_array());

        $nowTime = date("Y-m-d h:i:s");
        $data = array();

        if(!empty($moduleName)) {
            $data['ModuleName'] = $moduleName;
        }
        if(!empty($operator)) {
            $data['Operator'] = $operator;
        }
        if(!empty($bakOperator)) {
            $data['BakOperator'] = $bakOperator;
        }

        if(0==count($data)) {
            return array('success'=>true);
        }
        $data['LastTime'] = $nowTime;
        $this->db->where('ModuleID',$moduleId);
        $result = $this->db->update('cc_ModuleBase',$data);

        /*更新失败*/
        if(!$result) {
            throw new Exception("更新模块失败!");
        }
		/*更新成功*/
        $log['ApplicationID'] = $appId;
        $content = '';
        if( $moduleName != $query['ModuleName'] && !empty($moduleName)) {
            $content.=",[模块名:".$query['ModuleName'].'->'.$moduleName.']';
        }
        if( $operator != $query['Operator'] && !empty($operator)) {
            $content.=",[维护人:".$query['Operator'].'->'.$operator.']';
        }
        if( $bakOperator != $query['BakOperator'] && !empty($bakOperator))  {
            $content.=",[备份维护人:".$query['BakOperator'].'->'.$bakOperator.']';
        }

        $log['OpContent'] = '模块'.$content;
        $log['OpResult'] = $result ? 1 : 0;
        $log['OpTarget'] = '模块';
        $log['OpType'] = '修改';
        CCLog::addOpLogArr($log);
        return array('success'=>true);
    }

	/**
     * @查询所有的非默认模块
     * @param appId
     * @return 数组
     */
    public function listModuleNotDefault($appId) {
        $this->db->select('*');
        $this->db->from('cc_ModuleBase');
        $this->db->where('ApplicationID', $appId);
        $this->db->where('Default', 0);
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

	/**
	 * @查询业务下的所有模块
	 * @param appId
	 * @return 模块名称数组
	 */
	public function getModulesNameByappId($appId, $setIdArr = array()) {
		if(!is_array($setIdArr)) {
			$setIdArr = array($setIdArr);
		}
		$this->db->select('distinct (ModuleName)');
		$this->db->from('ModuleBase');
		$this->db->where('ApplicationID', $appId);
		if(!empty($setIdArr)) {
			$this->db->where_in('SetID', $setIdArr);
		}

		$query = $this->db->get();
		return $query->num_rows() ? $query->result_array() : array();
	}

	/**
	 * @根据业务Id和模块名查询模块
	 */
	public function getModulesIdByAppIdAndModuleName($appId, $setId, $moduleName) {
		if(!is_array($setId)) {
			$setId = array($setId);
		}

		$this->db->select('distinct (ModuleID)');
		$this->db->from('ModuleBase');
		$this->db->where('ApplicationID', $appId);
		if(!empty($setId))  {
			$this->db->where_in('SetID', $setId);
		}
		if(!empty($moduleName)) {
			$this->db->where_in('ModuleName', $moduleName);
		}

		$query = $this->db->get();
		return $query ? $query->result_array() : array();
	}

	/**
     * @查询资源池业务空闲机模块
     * @param appId
     * @return  返回空闲机模块字段
     */
    public function getResPoolIDByAppID($appId, $fields = '*') {
        $this->load->database('db');
        $this->db->select($fields);
        $this->db->from('cc_ModuleBase');
        $this->db->where('ApplicationID', $appId);
        $this->db->where('Default',1);
        $this->db->where('ModuleName','空闲机');
        $query = $this->db->get();
        if(0 == $query->num_rows()){
            return FALSE;
        }
        return $query->row_array();
    }

}