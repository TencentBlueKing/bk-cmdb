<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class SetBaseModel extends Cc_Model{

	public function __construct(){
		parent::__construct();
	}

    /*
     * @数据库中通过Id查询集群
     * @param setId,appId
     * @return 数组
     */
	public function getSetById($setId=array(), $appId=array()){
		$setIdArr = is_array($setId) ? $setId : explode(',', $setId);
		$appIdArr = is_array($appId) ? $appId : explode(',', $appId);

		if($setId!=array('') && !empty($setId)) {
			$this->db->where_in('SetID', $setIdArr);
		}

		if($appId!=array('') && !empty($appId)) {
			$this->db->where_in('ApplicationID', $appIdArr);
		}

		$query = $this->db->get('SetBase');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}


    /*
     * @数据库插入集群记录
     * @param 集群信息数组
     * @return bool
     */
	public function addSet($data) {
		$this->db->select('SetID');
		$this->db->where('SetName', $data['SetName']);
		$this->db->where('ApplicationID', $data['ApplicationID']);
		$query = $this->db->get('SetBase');

		if($query && $query->num_rows()==0){
			$query->free_result();
			$result = $this->db->insert('SetBase', $data);

			if(!$result) {
				$this->_errInfo = '添加集群失败!';
				$err = $this->db->error();
				CCLog::LogErr('添加集群失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
                return false;
			}
            $log['ApplicationID'] = $data['ApplicationID'];
            $log['OpContent'] = '集群名['. $data['SetName'] .']';
            $log['OpTarget'] = '集群';
            $log['OpType'] = '新增';
            CCLog::addOpLogArr($log);
			return $this->db->insert_id();
		}

		$this->_errInfo = '同名集群已存在[SetName='. $data['SetName'] .']!';
		return false;
	}

    /*
     * @数据库中删除集群
     * @param 集群信息数组
     * @return bool
     */
	public function delSetById($setId=array(), $appId=0){
		$setIdArr 	= !is_array($setId) ? explode(',', $setId) : array(intval($setId));
    	$appId 	= intval($appId);

		$this->db->where_in('SetID', $setIdArr);
		$appId!=array('') && !empty($appId) && $this->db->where('ApplicationID', $appId);
		$query = $this->db->get('SetBase');
		if(!$query || $query->num_rows()==0){
			$this->_errInfo = 'Group[id='. implode(',', $setId) .']不存在';
			return false;
		}
		$set = $query->result_array();

		$this->db->where_in('SetID', $setIdArr);
		$appId != array('') && !empty($appId) && $this->db->where('ApplicationID', $appId);
		$result = $this->db->delete('SetBase');

		if(!$result){
			$this->_errInfo = '删除集群失败！';
			$err = $this->db->error();
			CCLog::LogErr('删除集群失败！! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
		}

        $log = array();
		$log['ApplicationID'] = $appId;
		$log['OpContent'] = '删除集群，集群名['. implode(',', array_column($set, 'SetName')) .']';
		$log['OpFrom'] = 0;
		$log['OpName'] = '删除集群';
		$log['OpResult'] = $result ? 1 : 0;
		$log['OpTarget'] = '集群';
		$log['OpType'] = '删除';
		CCLog::addOpLogArr($log);
		return $result;
	}

    /*
     * @编辑集群
     * @param setId,data,appId
     * @return bool
     */
    public function editSet($setId, $data, $appId)
    {
        if(isset($data['SetName'])) {
            $this->db->select('SetID');
            $this->db->where('ApplicationID', $appId);
            $this->db->where('SetID <> ', $setId);
            $this->db->where('SetName', $data['SetName']);
            $query = $this->db->get('SetBase');
            if(0 != $query->num_rows()) {
                return false;
            }
        }

        $this->db->select('*');
        $this->db->from('SetBase');
        $this->db->where('SetID', $setId);
        $set = $this->db->get()->row_array();

        $this->db->where('SetID', $setId);
        $result = $this->db->update('SetBase', $data);
        $log['ApplicationID'] = $appId;
        $content = '';
        if(isset($data['SetName']) && $data['SetName'] != $set['SetName']) {
            $content.=",[集群名:".$set['SetName'].'->'.$data['SetName'].']';
        }
        if(isset($data['ChnName']) && $data['ChnName']!=$set['ChnName']) {
            $content.=",[中文名:".$set['ChnName'].'->'.$data['ChnName'].']';
        }
        if(isset($data['EnviType']) && $data['EnviType']!=$set['EnviType']) {
            $content.=",[环境类型:".$set['EnviType'].'->'.$data['EnviType'].']';
        }
        if(isset($data['ServiceStatus']) && $data['ServiceStatus']!=$set['ServiceStatus']) {
            $content.=",[服务状态:".$set['ServiceStatus'].'->'.$data['ServiceStatus'].']';
        }
        if(isset($data['Description']) && $data['Description']!=$set['Description']) {
            $content.=",[描述:".$set['Description'].'->'.$data['Description'].']';
        }
        $log['OpContent'] = "集群".$content;
        $log['OpResult'] = $result ? 1 : 0;
        $log['OpTarget'] = '集群';
        $log['OpType'] = '修改';
        CCLog::addOpLogArr($log);
        return true;
    }

    /*
     * @数据库中根据Id获取主机信息
     * @param setId
     * @return 主机数组
     */
    public function getHostBySetID($setId) {
        $this->db->select('HostID');
        $this->db->from('ModuleHostConfig');
        $this->db->where('SetID',$setId);
        $query = $this->db->get('');
        return $query->num_rows() ? $query->result_array() : array();
    }

    /*
     * @根据集群Id获取模块
     * @param setId
     * @return 主机数组
     */
    public function getModuleBySetId($setId, $fileds = '*') {
        $this->db->select($fileds);
        $this->db->from('ModuleBase');
        $this->db->where('SetID', $setId);
        $query = $this->db->get('');
        return $query->num_rows() ? $query->result_array() : array();
    }

    /*
     * @查询业务下非默认集群个数
     * @param 集群Id数组
     * @return 主机数组
     */
    public function listSetNotDefault($appId) {
        $this->db->select('*');
        $this->db->from('cc_SetBase');
        $this->db->where('ApplicationID', $appId);
        $this->db->where('Default', 0);
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

    /*
     * @查询集群下的模块个数
     * @param 集群Id数组
     * @return 主机数组
     */
    public function getModuleCountBySetId($setIdArr) {
        $this->db->select('count(ModuleID) as cnt,SetID');
        $this->db->from('cc_ModuleBase');
        $this->db->where_in('SetID', $setIdArr);
        $this->db->group_by('SetID');
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }
}