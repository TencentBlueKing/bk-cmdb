<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class HostBaseModel extends Cc_Model {

	public function _construct() {
		parent::_construct();
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
	public function getHostById($hostId = array(), $moduleId = array(), $setId = array(), $appId = array(), $start = '', $limit = '', $orderby = '', $direction = 'DESC') {
		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$setId = is_array($setId) ? $setId : explode(',', $setId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->database('db');
		$this->db->from('cc_HostBase');

		if(count($hostId) > 0) {
			$this->db->where_in('cc_HostBase.HostID', $hostId);
		}

		$this->db->select('cc_HostBase.*, group_concat(distinct cc_ModuleBase.ModuleName) as ModuleName, group_concat(distinct cc_SetBase.SetName) as SetName, group_concat(distinct cc_ApplicationBase.ApplicationName) as ApplicationName, group_concat(distinct cc_ApplicationBase.ApplicationID) as ApplicationID, group_concat(distinct cc_ApplicationBase.Owner) as Owner');	
		$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.HostID=cc_HostBase.HostID', 'LEFT');
		$this->db->join('cc_ApplicationBase', 'cc_ApplicationBase.ApplicationID=cc_ModuleHostConfig.ApplicationID', 'LEFT');
		$this->db->join('cc_SetBase', 'cc_SetBase.SetID=cc_ModuleHostConfig.SetID', 'LEFT');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID', 'LEFT');

		if($moduleId!==array('') && !empty($moduleId)) {
			$this->db->where_in('cc_ModuleHostConfig.ModuleID', $moduleId);
		}

		if($setId!==array('') && !empty($setId)) {
			$this->db->where_in('cc_ModuleHostConfig.SetID', $setId);
		}

		if($appId!==array('') && !empty($appId)) {
			$this->db->where_in('cc_ModuleHostConfig.ApplicationID', $appId);
		}

		$this->db->group_by('cc_HostBase.HostID');

		if($start !== '' && $limit !== '') {
			$this->db->limit($limit, $start);
		}

		if($orderby != '') {
			$this->db->order_by($orderby, $direction);
		}
		$query = $this->db->get();
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/**
	 * @id查询主机数
	 * @param hostId 主机ID
	 * @param moduleId 模块ID
	 * @param setId 大区ID
	 * @param appId 业务ID
	 * @param start 起始下标
	 * @param limit 欲取主机行数
	 * @param orderBy 排序字段
	 * @param direction 排序规则，desc/asc
	 * @return int 主机数量
	 */
	public function getHostCountById($hostId = array(), $moduleId = array(), $setId = array(), $appId = array()) {
		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$setId = is_array($setId) ? $setId : explode(',', $setId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->database('db');
		$this->db->from('cc_HostBase');

		if(count($hostId) > 0) {
			$this->db->where_in('cc_HostBase.HostID', $hostId);
		}

		if($moduleId != array('') || $setId != array('') || $appId != array('')) {

			$this->db->select('cc_HostBase.HostID');
			$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.HostID=cc_HostBase.HostID', 'LEFT');
			if($moduleId != array('') && !empty($moduleId)) {
				$this->db->where_in('cc_ModuleHostConfig.ModuleID', $moduleId);
			}

			if($setId != array('') && !empty($setId)) {
				$this->db->where_in('cc_ModuleHostConfig.SetID', $setId);
			}

			if($appId != array('') && !empty($appId)) {
				$this->db->where_in('cc_ModuleHostConfig.ApplicationID', $appId);
			}
		}

		$this->db->group_by('cc_HostBase.HostID');
		$query = $this->db->get();
		return $query ? $query->num_rows() : 0;
	}

	/**
	 * @ip查询主机
	 * @param hostIp 主机ip
	 * @param moduleId 模块ID
	 * @param setId 大区ID
	 * @param appId 业务ID
	 * @param start 起始下标
	 * @param limit 欲取主机行数
	 * @param orderBy 排序字段
	 * @param direction 排序规则，desc/asc
	 * @return array() 主机信息数组
	 */
	public function getHostByIp($hostIp = array(), $moduleId = array(), $setId = array(), $appId = array(), $start = '', $limit = '', $orderby = '', $direction = 'DESC') {

		$hostIp = is_array($hostIp) ? $hostIp : explode(',', $hostIp);
		$moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$setId = is_array($setId) ? $setId : explode(',', $setId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->database('db');
		$this->db->from('cc_HostBase');

		if(count($hostIp) > 0) {
			$this->db->group_start();
			$this->db->where_in('cc_HostBase.InnerIP', $hostIp);
			$this->db->or_where_in('cc_HostBase.OuterIP', $hostIp);
			$this->db->group_end();
		}

		$this->db->select('cc_HostBase.*, group_concat(distinct cc_ModuleBase.ModuleName) as ModuleName, group_concat(distinct cc_SetBase.SetName) as SetName, group_concat(distinct cc_ApplicationBase.ApplicationName) as ApplicationName, group_concat(distinct cc_ApplicationBase.ApplicationID) as ApplicationID, group_concat(distinct cc_ApplicationBase.Owner) as Owner');	
		$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.HostID=cc_HostBase.HostID', 'LEFT');
		$this->db->join('cc_ApplicationBase', 'cc_ApplicationBase.ApplicationID=cc_ModuleHostConfig.ApplicationID', 'LEFT');
		$this->db->join('cc_SetBase', 'cc_SetBase.SetID=cc_ModuleHostConfig.SetID', 'LEFT');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID', 'LEFT');

		if($moduleId != array('') && !empty($moduleId)) {
			$this->db->where_in('cc_ModuleHostConfig.ModuleID', $moduleId);
		}

		if($setId != array('') && !empty($setId)) {
			$this->db->where_in('cc_ModuleHostConfig.SetID', $setId);
		}

		if($appId != array('') && !empty($appId)) {
			$this->db->where_in('cc_ModuleHostConfig.ApplicationID', $appId);
		}

		$this->db->group_by('cc_HostBase.HostID');

		if($start !== '' && $limit !== '') {
			$this->db->limit($limit, $start);
		}

		if($orderby != '') {
			$this->db->order_by($orderby, $direction);
		}
		$query = $this->db->get();
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/**
	 * @查询HostBase表的字段
	 * @return array() 主机表字段的数组
	 */
	public function getAllowedFields() {
		$default_fields = array('HostID','AssetID', 'AutoRenew', 'BakOperator', 'BandWidth', 'CreateTime', 'DeadLineTime', 'DeviceClass', 'HardMemo', 'HostName', 'IDCArea', 'InnerIP', 'Memo', 'Operator', 'OSName', 'OuterIP', 'ServerType', 'Status', 'ZoneID', 'ZoneName', 'LastTime');

		$this->load->database('db');
		$fileds = $this->db->list_fields('cc_HostBase');

		$columns = array();
		foreach($fileds as $field) {
			if(!in_array($field, $columns)) {
				$columns[] = $field;
			}
		}

		return !empty($columns) ? $columns : $default_fields;
	}

	/**
	 * @多条件查询主机
	 * @param condition 查询条件
	 * @param distinct 是否去重
	 * @return array() 主机信息数组
	 */
	public function getHostByCondition($condition, $distinct = true, $start = 0, $limit = 0) {
		$appId = array();
		if(isset($condition['ApplicationID'])) {
			$appId = is_array($condition['ApplicationID']) ? $condition['ApplicationID'] : explode(',', $condition['ApplicationID']);
		}

		$setId = array();
		if(isset($condition['SetID'])) {
			$setId = is_array($condition['SetID']) ? $condition['SetID'] : explode(',', $condition['SetID']);
		}

		$moduleId = array();
		if(isset($condition['ModuleID'])) {
			$moduleId = is_array($condition['ModuleID']) ? $condition['ModuleID'] : explode(',', $condition['ModuleID']);
		}

		$assetId = array();
		if(isset($condition['AssetID'])) {
			$assetId = is_array($condition['AssetID']) ? $condition['AssetID'] : explode(',', $condition['AssetID']);
		}

		$HostName = array();
		if(isset($condition['HostName'])) {
			$HostName = is_array($condition['HostName']) ? $condition['HostName'] : explode(',', $condition['HostName']);
		}

		$SN = array();
		if(isset($condition['SN'])) {
			$SN = is_array($condition['SN']) ? $condition['SN'] : explode(',', $condition['SN']);
		}

		$hostId = array();
		if(isset($condition['HostID'])) {
			$hostId = is_array($condition['HostID']) ? $condition['HostID'] : explode(',', $condition['HostID']);
		}

		$innerIp = array();
		if(isset($condition['InnerIP'])) {
			$innerIp = is_array($condition['InnerIP']) ? $condition['InnerIP'] : explode(',', $condition['InnerIP']);
		}

		$outerIp = array();
		if(isset($condition['OuterIP'])) {
			$outerIp = is_array($condition['OuterIP']) ? $condition['OuterIP'] : explode(',', $condition['OuterIP']);
		}

		$IfInnerIPexact = '';
		if(isset($condition['IfInnerIPexact'])) {
			$IfInnerIPexact = $condition['IfInnerIPexact'];
		}
		
		$IfOuterexact = '';
		if(isset($condition['IfOuterexact'])) {
			$IfOuterexact = $condition['IfOuterexact'];
		}

		$CreateTime = '';
		if(isset($condition['CreateTime'])) {
			$CreateTime = $condition['CreateTime'];
		}

		$DeadLineTime = '';
		if(isset($condition['DeadLineTime'])) {
			$DeadLineTime = $condition['DeadLineTime'];
		}

		unset($condition['ApplicationID'], $condition['SetID'], $condition['ModuleID'], $condition['AssetID'], $condition['SN'], $condition['HostName'], $condition['HostID'], $condition['InnerIP'], $condition['OuterIP'], $condition['IfInnerIPexact'], $condition['IfOuterexact'], $condition['CreateTime'], $condition['DeadLineTime']);

		$this->load->database('db');
		$this->db->from('cc_HostBase');

		foreach($condition as $_k=>$_v) {
			$this->db->where_in('cc_HostBase.'.$_k, $_v);
		}
		
		$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.HostID=cc_HostBase.HostID', 'LEFT');
		$this->db->join('cc_ApplicationBase', 'cc_ApplicationBase.ApplicationID=cc_ModuleHostConfig.ApplicationID', 'LEFT');
		$this->db->join('cc_SetBase', 'cc_SetBase.SetID=cc_ModuleHostConfig.SetID', 'LEFT');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID', 'LEFT');

		if(count($hostId)>0 && $hostId!=array('')) {
			$this->db->where_in('cc_HostBase.HostID', $hostId);
		}

		if(count($assetId)>0 && $assetId!=array('')) {
			if(count($assetId)>1) {
				$this->db->where_in('cc_HostBase.AssetID', $assetId);
			}else {
				$this->db->like('cc_HostBase.AssetID', $assetId[0]);
			}
		}

		if(count($HostName)>0 && $HostName!=array('')) {
			if(count($HostName)>1) {
				$this->db->where_in('cc_HostBase.HostName', $HostName);
			}else {
				$this->db->like('cc_HostBase.HostName', $HostName[0]);
			}
		}

		if(count($SN)>0 && $SN!=array('')) {
			if(count($SN)>1) {
				$this->db->where_in('cc_HostBase.SN', $SN);
			}else {
				$this->db->like('cc_HostBase.SN', $SN[0]);
			}
		}

		if(count($innerIp)>0 && $innerIp != array('')) {
			if($IfInnerIPexact) {
				$this->db->where_in('cc_HostBase.InnerIP', $innerIp);
			}else {
				if(isset($IfInnerIPexact)) {
					$this->db->like('cc_HostBase.InnerIP', $innerIp[0]);
					if(count($innerIp) > 1) {
						foreach ($innerIp as $key => $value) {
							if($key > 0) {
								$this->db->or_like('cc_HostBase.InnerIP', $value);
							}
						}
					}
				}else {
					$this->db->where_in('cc_HostBase.InnerIP', $innerIp);
				}
			}
		}
		
		if(count($outerIp) > 0 && $outerIp != array('')) {
			if($IfOuterexact) {
				$this->db->where_in('cc_HostBase.OuterIP', $outerIp);
			}else {
				if(isset($IfOuterexact)) {
					$this->db->like('cc_HostBase.OuterIP', $outerIp[0]);
					if(count($outerIp) > 1) {
						foreach ($outerIp as $key => $value) {
							if($key > 0) {
								$this->db->or_like('cc_HostBase.OuterIP', $value);
							}
						}
					}
				}else {
					$this->db->where_in('cc_HostBase.OuterIP', $outerIp);
				}
			}
		}

		if($CreateTime) {
			$this->db->where('cc_HostBase.CreateTime >=', $CreateTime . ' 00:00:00');
			$this->db->where('cc_HostBase.CreateTime <=', $CreateTime . ' 23:59:59');
		}

		if($DeadLineTime) {
			$this->db->where('cc_HostBase.DeadLineTime >=', $DeadLineTime . ' 00:00:00');
			$this->db->where('cc_HostBase.DeadLineTime <=', $DeadLineTime . ' 23:59:59');
		}

		if(count($moduleId)>0 && $moduleId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.ModuleID', $moduleId);
		}

		if(count($setId)>0 && $setId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.SetID', $setId);
		}

		if(count($appId)>0 && $appId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.ApplicationID', $appId);
		}

		if($distinct) {
			$this->db->select('cc_HostBase.*, group_concat(distinct cc_ModuleBase.ModuleID) as ModuleID, group_concat(distinct cc_ModuleBase.ModuleName) as ModuleName, group_concat(distinct cc_SetBase.SetID) as SetID, group_concat(distinct cc_SetBase.SetName) as SetName, group_concat(distinct cc_ApplicationBase.ApplicationID) as ApplicationID, group_concat(distinct cc_ApplicationBase.ApplicationName) as ApplicationName, group_concat(distinct cc_ApplicationBase.Owner) as Owner');
			$this->db->group_by('cc_HostBase.HostID');
		}else {
			$this->db->select('cc_HostBase.*, cc_ModuleBase.ModuleName as ModuleName, cc_SetBase.SetName as SetName, cc_ApplicationBase.ApplicationName as ApplicationName, cc_ApplicationBase.Owner as Owner');	
		}
		$this->db->order_by('cc_HostBase.InnerIP asc');
		if($start+$limit >0) {
			$this->db->limit($limit, $start);
		}

		$query = $this->db->get();

		return $query ? $query->result_array() : array();
	}

	/**
	 * @多条件查询主机数
	 * @param condition 查询条件
	 * @param distinct 是否去重
	 * @return int 主机数
	 */
	public function getHostCountByCondition($condition, $distinct = true) {
		$appId = array();
		if(isset($condition['ApplicationID'])) {
			$appId = is_array($condition['ApplicationID']) ? $condition['ApplicationID'] : explode(',', $condition['ApplicationID']);
		}

		!$appId && $appId[] = $this->input->cookie('defaultAppId');

		$setId = array();
		if(isset($condition['SetID'])) {
			$setId = is_array($condition['SetID']) ? $condition['SetID'] : explode(',', $condition['SetID']);
		}

		$moduleId = array();
		if(isset($condition['ModuleID'])) {
			$moduleId = is_array($condition['ModuleID']) ? $condition['ModuleID'] : explode(',', $condition['ModuleID']);
		}

		$assetId = array();
		if(isset($condition['AssetID'])) {
			$assetId = is_array($condition['AssetID']) ? $condition['AssetID'] : explode(',', $condition['AssetID']);
		}

		$HostName = array();
		if(isset($condition['HostName'])) {
			$HostName = is_array($condition['HostName']) ? $condition['HostName'] : explode(',', $condition['HostName']);
		}

		$SN = array();
		if(isset($condition['SN'])) {
			$SN = is_array($condition['SN']) ? $condition['SN'] : explode(',', $condition['SN']);
		}

		$hostId = array();
		if(isset($condition['HostID'])) {
			$hostId = is_array($condition['HostID']) ? $condition['HostID'] : explode(',', $condition['HostID']);
		}

		$innerIp = array();
		if(isset($condition['InnerIP'])) {
			$innerIp = is_array($condition['InnerIP']) ? $condition['InnerIP'] : explode(',', $condition['InnerIP']);
		}

		$outerIp = array();
		if(isset($condition['OuterIP'])) {
			$outerIp = is_array($condition['OuterIP']) ? $condition['OuterIP'] : explode(',', $condition['OuterIP']);
		}
		
		$IfInnerIPexact = '';
		if(isset($condition['IfInnerIPexact'])) {
			$IfInnerIPexact = $condition['IfInnerIPexact'];
		}
		
		$IfOuterexact = '';
		if(isset($condition['IfOuterexact'])) {
			$IfOuterexact = $condition['IfOuterexact'];
		}

		$CreateTime = '';
		if(isset($condition['CreateTime'])) {
			$CreateTime = $condition['CreateTime'];
		}

		$DeadLineTime = '';
		if(isset($condition['DeadLineTime'])) {
			$DeadLineTime = $condition['DeadLineTime'];
		}

		unset($condition['ApplicationID'], $condition['SetID'], $condition['ModuleID'], $condition['AssetID'], $condition['SN'], $condition['HostName'], $condition['HostID'], $condition['InnerIP'], $condition['OuterIP'], $condition['start'], $condition['limit'], $condition['IfInnerIPexact'], $condition['IfOuterexact'], $condition['CreateTime'], $condition['DeadLineTime']);

		$this->load->database('db');
		$this->db->from('cc_HostBase');

		foreach($condition as $_k=>$_v) {
			$this->db->where_in('cc_HostBase.'.$_k, $_v);
		}
		
		$this->db->join('cc_ModuleHostConfig', 'cc_ModuleHostConfig.HostID=cc_HostBase.HostID', 'LEFT');
		$this->db->join('cc_ApplicationBase', 'cc_ApplicationBase.ApplicationID=cc_ModuleHostConfig.ApplicationID', 'LEFT');
		$this->db->join('cc_SetBase', 'cc_SetBase.SetID=cc_ModuleHostConfig.SetID', 'LEFT');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID', 'LEFT');

		if(count($hostId)>0 && $hostId!=array('')) {
			$this->db->where_in('cc_HostBase.HostID', $hostId);
		}

		if(count($assetId)>0 && $assetId!=array('')) {
			if(count($assetId)>1) {
				$this->db->where_in('cc_HostBase.AssetID', $assetId);
			}else {
				$this->db->like('cc_HostBase.AssetID', $assetId[0]);
			}
		}

		if(count($HostName)>0 && $HostName!=array('')) {
			if(count($HostName)>1) {
				$this->db->where_in('cc_HostBase.HostName', $HostName);
			}else {
				$this->db->like('cc_HostBase.HostName', $HostName[0]);
			}
		}

		if(count($SN)>0 && $SN != array('')) {
			if(count($SN)>1) {
				$this->db->where_in('cc_HostBase.SN', $SN);
			}else {
				$this->db->like('cc_HostBase.SN', $SN[0]);
			}
		}

		if(count($innerIp)>0 && $innerIp!=array('')) {
			if($IfInnerIPexact) {
				$this->db->where_in('cc_HostBase.InnerIP', $innerIp);
			}else {
				if(isset($IfInnerIPexact)) {
					$this->db->like('cc_HostBase.InnerIP', $innerIp[0]);
					if(count($innerIp) > 1) {
						foreach ($innerIp as $key => $value) {
							if($key > 0) {
								$this->db->or_like('cc_HostBase.InnerIP', $value);
							}
						}
					}
				}else {
					$this->db->where_in('cc_HostBase.InnerIP', $innerIp);
				}
			}
		}
		
		if(count($outerIp)>0 && $outerIp!=array('')) {
			if($IfOuterexact) {
				$this->db->where_in('cc_HostBase.OuterIP', $outerIp);
			}else {
				if(isset($IfOuterexact)) {
					$this->db->like('cc_HostBase.OuterIP', $outerIp[0]);
					if(count($outerIp) > 1) {
						foreach ($outerIp as $key => $value) {
							if($key > 0) {
								$this->db->or_like('cc_HostBase.OuterIP', $value);
							}
						}
					}
				}else {
					$this->db->where_in('cc_HostBase.OuterIP', $outerIp);
				}
			}
		}

		if($CreateTime) {
			$this->db->where('cc_HostBase.CreateTime >=', $CreateTime . ' 00:00:00');
			$this->db->where('cc_HostBase.CreateTime <=', $CreateTime . ' 23:59:59');
		}

		if($DeadLineTime) {
			$this->db->where('cc_HostBase.DeadLineTime >=', $DeadLineTime . ' 00:00:00');
			$this->db->where('cc_HostBase.DeadLineTime <=', $DeadLineTime . ' 23:59:59');
		}

		if(count($moduleId)>0 && $moduleId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.ModuleID', $moduleId);
		}

		if(count($setId)>0 && $setId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.SetID', $setId);
		}

		if(count($appId)>0 && $appId != array('')) {
			$this->db->where_in('cc_ModuleHostConfig.ApplicationID', $appId);
		}

		if($distinct) {
			$this->db->select('cc_HostBase.HostID');
			$this->db->group_by('cc_HostBase.HostID');
		}else {
			$this->db->select('cc_HostBase.HostID');	
		}

		$query = $this->db->get();
		return $query ? $query->num_rows() : 0;
	}

	/**
	 * @id查询主机相关Id
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
	public function getHostAllIdById($hostId = array(), $moduleId = array(), $setId = array(), $appId = array()) {

		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$setId = is_array($setId) ? $setId : explode(',', $setId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);
		$this->load->database('db');

		if($hostId != array('') && !empty($hostId)) {
			$this->db->where_in('HostID', $hostId);
		}

		if($moduleId != array('') && !empty($moduleId)) {
			$this->db->where_in('ModuleID', $moduleId);
		}

		if($setId != array('') && !empty($setId)) {
			$this->db->where_in('SetID', $setId);
		}

		if($appId != array('') && !empty($appId)) {
			$this->db->where_in('ApplicationID', $appId);
		}

		$query = $this->db->get('cc_ModuleHostConfig');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/**
	* @Id更新主机属性
	* @param $stdProperty array 主机属性数组，下标为HostBase表字段名，值为需要更新的值
	* @param $hostId int 主机Id
	* @return boolean 成功 or 失败
	*/
	public function updateHostById($stdProperty, $hostId){
		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$this->load->database('db');
		$fields = $this->db->list_fields('HostBase');

		$columns = array();
		foreach($fields as $field) {
			if(!in_array($field, $columns) && !in_array($field, array('HostID'))) {
				$columns[] = $field;
			}
		}

		$propertyNames = array_keys($stdProperty);
		$diff = array_diff($propertyNames, $columns);
		if($diff) {
			$this->_errInfo = '不支持修改字段['. implode(',', $diff) .']';
			return false;
		}

		$this->db->where_in('HostID', $hostId);
		$query = $this->db->get('cc_HostBase');

		if(!$query || $query->num_rows()==0) {
			$this->_errInfo = '查询主机，失败！没有主机[id='. implode(',', $hostId) .']';
			CCLog::LogErr($this->_errInfo);
			return false;
		}

		$result = $query->result_array();
		$innerIp = array_column($result, 'InnerIP');

		$opContent = '更新主机[id='. implode(',', $innerIp) .']属性：';
		foreach($query->result_array() as $_h) {
			foreach($propertyNames as $_pn) {
				$opContent .= $_h[$_pn] .'->'. $stdProperty[$_pn] .' | ';
			}
		}
		$opContent = trim($opContent, ' | ');

		$this->load->database();
		$this->db->where_in('HostID', $hostId);
		$return = $this->db->update('cc_HostBase', $stdProperty);

		if(!$return) {
			$this->_errInfo = '更新失败！';
			$err = $this->db->error();
			CCLog::LogErr('sql执行失败！mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
		}

		$log = array();
		$log['ApplicationID'] = 0;
		$log['OpContent'] = $opContent;
		$log['OpFrom'] = 0;
		$log['OpName'] = '修改主机属性';
		$log['OpResult'] = $return ? 1 : 0;
		$log['OpTarget'] = '主机';
		$log['OpType'] = '更新';
		CCLog::addOpLogArr($log);

		return $return;
	}

	/**
	* @转移一台主机到多个模块
	* @param $hostId 主机ID
	* @param $moduleId 模块ID
	* @param $appId 业务ID
	* @param $isIncrement 是否增量更新
	* @return boolean 成功 or 失败
	*/
	public function modSingleHostToMultiModule($hostId, $moduleId, $appId, $isIncrement) {
		$hostId = intval($hostId);
		$moduleId = is_array($moduleId) ? $moduleId : explode(',', $moduleId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->model('ModuleBaseModel');
		$module = $this->ModuleBaseModel->getModuleByHostId($hostId);

		foreach($module as $_m) {
			if($_m['ModuleName']==='空闲机' || !$isIncrement) {	
				$delRes = $this->delSingleHostFromSingleModule($hostId, $_m['ModuleID']);

				if(!$delRes) {
					return false;
				}
			}
		}

		$module = $this->ModuleBaseModel->getModuleById($moduleId);
		$moduleId2Info = array();
		foreach($module as $_m) {
			$moduleId2Info[$_m['ModuleID']] = $_m;
		}

		foreach($moduleId as $_mid) {
			$addRes = $this->addSingleHostToSingleModule($hostId, $_mid);

			if(!$addRes) {
				return false;
			}
		}
		return true;
	}

	/**
	* @转移多台主机到一个模块
	* @param $hostId 主机ID
	* @param $moduleId 模块ID
	* @param $appId 业务ID
	* @param $isIncrement 是否增量更新
	* @return boolean 成功 or 失败
	*/
	public function modMultiHostToSingleModule($hostId, $moduleId, $appId, $isIncrement) {
		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$moduleId = intval($moduleId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->model('ModuleBaseModel');
		$module = $this->ModuleBaseModel->getModuleByHostId($hostId);
		
		$hostId2ModuleInfo = array();
		foreach($module as $_m) {
			$hostId2ModuleInfo[$_m['HostID']][] = $_m;
		}

		$module = $this->ModuleBaseModel->getModuleById($moduleId);
		$host2ModuleNameSuc = array();
		foreach($hostId as $_h) {
            if(isset($hostId2ModuleInfo[$_h])) {
                foreach($hostId2ModuleInfo[$_h] as $_m) {
                    if($_m['ModuleName'] === '空闲机' || !$isIncrement) {
                        $delRes = $this->delSingleHostFromSingleModule($_h, $_m['ModuleID']);

                        if(!$delRes) {
                            return false;
                        }
                    }
                }
            }


			$addRes = $this->addSingleHostToSingleModule($_h, $module[0]['ModuleID']);
			if(!$addRes) {
				return false;
			}
		}
		return true;
	}

	/**
	* @主机操作-快速分配页面，分配主机资源
	* @param hostId 主机Id
	* @param moduleId 业务的空闲机模块Id
	* @param appId 业务Id
	* @return boolean 分配成功or失败
	*/
	public function quickDistribute($hostId, $moduleId, $appId){
		$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
		$moduleId = intval($moduleId);
		$appId = is_array($appId) ? $appId : explode(',', $appId);

		$this->load->model('ModuleBaseModel');
		$module = $this->ModuleBaseModel->getModuleByHostId($hostId);
		
		$hostId2ModuleInfo = array();
		foreach($module as $_m) {
			$hostId2ModuleInfo[$_m['HostID']][] = $_m;
		}

		$module = $this->ModuleBaseModel->getModuleById($moduleId);
		$host2ModuleNameSuc = array();
		foreach($hostId as $_h) {
	
			$moduleId = $module[0]['ModuleID'];
			$setId = $module[0]['SetID'];
			$appId = $module[0]['ApplicationID'];

			$this->db->select('cc_ModuleHostConfig.HostID,cc_ModuleBase.ModuleID,cc_ModuleBase.ModuleName,cc_ApplicationBase.ApplicationName');
			$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID');
			$this->db->join('cc_ApplicationBase', 'cc_ApplicationBase.ApplicationID=cc_ModuleHostConfig.ApplicationID');
			$this->db->where('cc_ModuleHostConfig.HostID', $_h);
			$query = $this->db->get('cc_ModuleHostConfig');

			$result = $query && $query->num_rows()>0 ? $query->result_array() : array();

			$isDistributed = false;//是否已存在关联关系
			foreach($result as $_r) {
				if($_r['ModuleID'] == $moduleId) {//已存在关联关系
					$isDistributed = true;
					continue;
				}

				if($_r['ApplicationName']=='资源池') {
					$delRes = $this->delSingleHostFromSingleModule($_h, $_r['ModuleID']);
					if(!$delRes) {
						return false;
					}
				}
			}

			if($isDistributed) {//已存在关联关系
				continue;
			}

			$addRes = $this->addSingleHostToSingleModule($_h, $moduleId);
			if(!$addRes) {
				return false;
			}
		}
		return true;
	}

	/**
	* @移除一台主机到一个模块的关联关系
	* @param $hostId 主机ID
	* @param $moduleId 模块ID
	* @return boolean 成功 or 失败
	*/
	private function delSingleHostFromSingleModule($hostId, $moduleId) {
		$hostId = intval($hostId);
		$moduleId = intval($moduleId);

		$this->load->database('db');
		$this->db->select('cc_HostBase.HostID,cc_HostBase.InnerIP,cc_ModuleHostConfig.HostID,cc_ModuleBase.ModuleName,cc_ModuleBase.ApplicationID');
		$this->db->from('cc_ModuleHostConfig');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID');
		$this->db->join('cc_HostBase', 'cc_HostBase.HostID=cc_ModuleHostConfig.HostID');
		$this->db->where('cc_ModuleHostConfig.HostID', $hostId);
		$this->db->where('cc_ModuleHostConfig.ModuleID', $moduleId);
		$query = $this->db->get();

		if(!$query || $query->num_rows() == 0) {
			$this->_errInfo = '主机模块关系不存在！';
			CCLog::LogErr('主机模块关系不存在! sql:'. $this->db->last_query());
			return false;
		}

		$result = $query->result_array();
		$appId = $result[0]['ApplicationID'];
		$hostIp = $result[0]['InnerIP'];
		$oldModuleName = array_column($result, 'ModuleName');


		$this->db->where('HostID', $hostId);
		$this->db->where('ModuleID', $moduleId);
		$delRes = $this->db->delete('cc_ModuleHostConfig');

		if(!$delRes) {
			$this->_errInfo = '解除主机模块关系，失败！';
			CCLog::LogErr('解除主机模块关系，失败！ sql:'. $this->db->last_query());

			$log = array();
			$log['ApplicationID'] = $appId;
			$log['OpContent'] = '主机['. $hostIp .']移出模块['. implode(',', $oldModuleName) .']';
			$log['OpFrom'] = 0;
			$log['OpName'] = '主机移出模块';
			$log['OpResult'] = 0;
			$log['OpTarget'] = '主机';
			$log['OpType'] = '更新';
			CCLog::addOpLogArr($log);
			return false;
		}


		$log = array();
		$log['ApplicationID'] = $appId;
		$log['OpContent'] = '主机['. $hostIp .']移出模块['. implode(',', $oldModuleName) .']';
		$log['OpFrom'] = 0;
		$log['OpName'] = '主机移出模块';
		$log['OpResult'] = 1;
		$log['OpTarget'] = '主机';
		$log['OpType'] = '更新';
		CCLog::addOpLogArr($log);
		return true;
	}

	/**
	* @添加一台主机到一个模块的关联关系
	* @param $hostId 主机ID
	* @param $moduleId 模块ID
	* @return boolean 成功 or 失败
	*/
	private function addSingleHostToSingleModule($hostId, $moduleId) {
		$hostId = intval($hostId);
		$moduleId = intval($moduleId);

		$this->load->database('db');

		$this->db->select('InnerIP');
		$this->db->where('HostID', $hostId);
		$query = $this->db->get('cc_HostBase');
		if(!$query || $query->num_rows() == 0) {
			$this->_errInfo = '主机不存在！';
			CCLog::LogErr('主机不存在！! sql:'. $this->db->last_query());
			return false;
		}
		$result = $query->row_array();
		$hostIp = $result['InnerIP'];

		$this->db->select('ModuleName,ApplicationID,SetID');
		$this->db->where('ModuleID', $moduleId);
		$query = $this->db->get('cc_ModuleBase');
		if(!$query || $query->num_rows() == 0) {
			$this->_errInfo = '模块不存在！';
			CCLog::LogErr('模块不存在！! sql:'. $this->db->last_query());
			return false;
		}

		$result = $query->row_array();
		$appId = $result['ApplicationID'];
		$setId = $result['SetID'];
		$newModuleName = $result['ModuleName'];

		$this->db->select('cc_ModuleHostConfig.HostID,cc_ModuleBase.ModuleID,cc_ModuleBase.ModuleName');
		$this->db->join('cc_ModuleBase', 'cc_ModuleBase.ModuleID=cc_ModuleHostConfig.ModuleID');
		$this->db->where('cc_ModuleHostConfig.HostID', $hostId);
		$this->db->where('cc_ModuleHostConfig.ApplicationID', $appId);
		$query = $this->db->get('cc_ModuleHostConfig');

		$result = $query && $query->num_rows()>0 ? $query->result_array() : array();
		foreach($result as $_r) {
			if($_r['ModuleName'] == '空闲机') {
				$delRes = $this->delSingleHostFromSingleModule($hostId, $_r['ModuleID']);
				if(!$delRes) {
					return false;
				}
			}

			if($_r['ModuleID'] == $moduleId) {
				return true;
			}
		}

		$data = array();
		$data['HostID'] = $hostId;
		$data['ModuleID'] = $moduleId;
		$data['SetID'] = $setId;
		$data['ApplicationID'] = $appId;
		$addRes = $this->db->insert('cc_ModuleHostConfig', $data);

		if(!$addRes) {
			$this->_errInfo = '主机移入模块，失败！';
			CCLog::LogErr('主机移入模块，失败！ sql:'. $this->db->last_query());

			$log = array();
			$log['ApplicationID'] = $appId;
			$log['OpContent'] = '主机['. $hostIp .']移入模块['. $newModuleName .']';
			$log['OpFrom'] = 0;
			$log['OpName'] = '主机移入模块';
			$log['OpResult'] = 0;
			$log['OpTarget'] = '主机';
			$log['OpType'] = '更新';
			CCLog::addOpLogArr($log);
			return false;
		}


		$log = array();
		$log['ApplicationID'] = $appId;
		$log['OpContent'] = '主机['. $hostIp .']移入模块['. $newModuleName .']';
		$log['OpFrom'] = 0;
		$log['OpName'] = '主机移入模块';
		$log['OpResult'] = 1;
		$log['OpTarget'] = '主机';
		$log['OpType'] = '更新';
		CCLog::addOpLogArr($log);
		return true;
	}
    /**
     * @根据固资号查询主机
     * @param $AssetID
     * @return true or false
     */
    public function getHostByAssetID($AssetID, $fields='*') {
        if(!is_array($AssetID)) {
            $AssetID = array($AssetID);
        }
        $this->load->database();
        $this->db->select($fields);
        $this->db->from('cc_HostBase');
        $this->db->where_in('AssetID', $AssetID);
        $query  = $this->db->get();
        if(0 == $query->num_rows()) {
            return FALSE;
        }
        return $query->result_array();
    }


    /**
     * @根据AssetID更新主机信息
     * @param $AssetID，主机信息
     * @return true or false
     */
    public function updateHostByAssetID($AssetID,$data) {
        $this->db->where('AssetID', $AssetID);
        $this->db->update('cc_HostBase', $data);
        return TRUE;
    }


    /**
     * @批量插入主机
     * @param 主机信息
     * @return HostID
     */
    public function AddHost($data) {
        $this->db->insert('cc_HostBase', $data);
        $HostID = $this->db->insert_id();
        if(!$HostID) {
        	CCLog::LogErr('插入主机['. $data['AssetID'] .']失败, sql:'. $this->db->last_query());
    		$this->_errInfo = '插入主机['. $data['AssetID'] .']失败';
    		return false;
    	}
        return $HostID;
    }

    /**
     * @按业务维度统计主机个数
     * @param 无
     * @return 主机个数管理数组
     */
    public function StatHostByApp($AppIDArr = array()) {
        $this->load->database();
        $this->db->select('count(distinct(cc_ModuleHostConfig.HostID)) as cnt,ApplicationID');
        $this->db->from('cc_ModuleHostConfig');
        $this->db->group_by('ApplicationID');
        if(0 != count($AppIDArr)) {
            $this->db->where_in('ApplicationID', $AppIDArr);
        }
        $query = $this->db->get();
        if(0 == $query->num_rows()) {
            return False;
        }
        return $query->result_array();
    }

    /**
     * @按开发商维度
     * @param 无
     * @return 主机个数管理数组
     */
    public function StatHostByOwner() {
        $this->load->database();
        $this->db->select('count(distinct(cc_ModuleHostConfig.HostID)) as cnt,cc_ApplicationBase.Owner');
        $this->db->from('cc_ModuleHostConfig');
        $this->db->join('cc_ApplicationBase','cc_ModuleHostConfig.ApplicationID=cc_ApplicationBase.ApplicationID');
        $this->db->group_by('cc_ApplicationBase.Owner');
        $query = $this->db->get();
        if(0 == $query->num_rows()) {
            return False;
        }
        return $query->result_array();
    }

    /**
     * @批量的删除主机
     * @param   HostIDArr
     * @return  无
     */
    public function RecovHostByHostIDArr($HostIDArr) {
        if(0 == count($HostIDArr)) {
            return;
        }
        $this->db->where_in('HostID',$HostIDArr);
        $this->db->delete('cc_HostBase');
    }

    /**
    * @删除主机
    * @param $hostId 主机Id
    * @param $moduleId 模块Id
    * @param @setId groupId
    * @param @appId 业务Id
    */
    public function deleteHostById($hostId = array(), $moduleId = array(), $setId = array(), $appId = 0) {
    	$hostId 	= !is_array($hostId) ? explode(',', $hostId) : array(intval($hostId));
    	$moduleId 	= !is_array($moduleId) ? explode(',', $moduleId) : array(intval($moduleId));
    	$setId 		= !is_array($setId) ? explode(',', $setId) : array(intval($setId));
    	$appId 		= intval($appId);

    	if(empty($appId)) {
    		$this->_errInfo = '缺少appId';
    		return false;
    	}

    	if(empty($hostId) && empty($moduleId) && empty($setId)) {
    		$this->_errInfo = 'hostId、moduleId、setId不能同时为空';
    		return false;
    	}

    	/**
    	* 待移除主机模块关系的HostID
    	*/
    	$result = $this->getHostAllIdById($hostId, $moduleId, $setId, $appId);

    	$hostId2ModuleId = array();
    	foreach($result as $_r) {
    		$hostId2ModuleId[$_r['HostID']][] =  $_r['ModuleID'];
    		if(count($hostId2ModuleId[$_r['HostID']])>1) {
    			$this->_errInfo = '主机属于多个模块，不能直接移至空闲机模块';
    			return false;
    		}
    	}

    	$hostId_pre = array();
    	foreach($result as $_r) {
    		if(!in_array($_r['HostID'], $hostId_pre)){
    			$hostId_pre[] = $_r['HostID'];
    		}
    	}

    	if($hostId != array('') && !empty($hostId)) {
    		$this->db->where_in('HostID', $hostId);
    	}

    	if($moduleId != array('') && !empty($moduleId)) {
    		$this->db->where_in('ModuleID', $moduleId);
    	}

    	if($setId != array('') && !empty($setId)) {
    		$this->db->where_in('SetID', $setId);
    	}

    	$this->db->where('ApplicationID', $appId);

    	$delRes = $this->db->delete('cc_ModuleHostConfig');
    	if(!$delRes) {
    		$this->_errInfo = '删除主机，失败!';
    		CCLog::LogErr('删除主机，失败！message: '. $this->db->last_query());
    		return false;
    	}

    	/**
		 * 不再属于其他模块的主机，转入空闲机模块
    	*/
    	$result = $this->getHostAllIdById($hostId_pre, $moduleId, $setId, $appId);
    	$hostId_suf = array();
    	foreach($result as $_r) {
    		if(!in_array($_r['HostID'], $hostId_suf)) {
    			$hostId_suf[] = $_r['HostID'];
    		}
    	}

    	$host_2_empty_module = array_diff($hostId_pre, $hostId_suf);
    	if(!empty($host_2_empty_module)) {
	    	$this->load->model('ModuleBaseModel');
	    	$resPool = $this->ModuleBaseModel->getResPoolIDByAppID($appId);
	    	if(!$resPool || !isset($resPool['ModuleID'])) {
	    		$this->_errInfo = '缺少空闲机模块!';
    			CCLog::LogErr('缺少空闲机模块！appId: '. print_r($appId, true));
    			return false;
	    	}

	    	$result = true;
	    	foreach($host_2_empty_module as $_h) {
	    		$result = $result && $this->addSingleHostToSingleModule($_h, $resPool['ModuleID']);
	    		if(!$result) {
	    			return false;
	    		}
	    	}
	    }

	    return true;
    }

    /**
    * @上交主机
    * @param $hostId 主机Id
    * @param $moduleId 模块Id
    * @param @setId groupId
    * @param @appId 业务Id
    * @return boolean 成功 or 失败
    */
    public function resHostModule($hostId = array(), $moduleId = array(), $setId = array(), $appId = 0) {
    	$hostId 	= !is_array($hostId) ? explode(',', $hostId) : array(intval($hostId));
    	$moduleId 	= !is_array($moduleId) ? explode(',', $moduleId) : array(intval($moduleId));
    	$setId 		= !is_array($setId) ? explode(',', $setId) : array(intval($setId));
    	$appId 		= intval($appId);

    	if(empty($appId)) {
    		$this->_errInfo = '缺少appId';
    		return false;
    	}

    	if(empty($hostId) && empty($moduleId) && empty($setId)) {
    		$this->_errInfo = 'hostId、moduleId、setId不能同时为空';
    		return false;
    	}

    	/**
    	* 待移除主机模块关系的HostID
    	*/
    	$result = $this->getHostAllIdById($hostId, $moduleId, $setId, $appId);
    	$hostId_pre = array();
    	foreach($result as $_r) {
    		if(!in_array($_r['HostID'], $hostId_pre)) {
    			$hostId_pre[] = $_r['HostID'];
    		}
    	}

    	if($hostId != array('') && !empty($hostId)) {
    		$this->db->where_in('HostID', $hostId);
    	}

    	if($moduleId != array('') && !empty($moduleId)) {
    		$this->db->where_in('ModuleID', $moduleId);
    	}

    	if($setId != array('') && !empty($setId)) {
    		$this->db->where_in('SetID', $setId);
    	}

    	$this->db->where('ApplicationID', $appId);

    	$delRes = $this->db->delete('cc_ModuleHostConfig');
    	if(!$delRes) {

    		$this->_errInfo = '删除主机，失败!';
    		CCLog::LogErr('删除主机，失败！message: '. $this->db->last_query());
    		return false;
    	}

    	/**
		 * 不再属于开发商所有业务下其他模块的主机，转入空闲机模块
    	*/
    	$result = $this->getHostAllIdById($hostId_pre, $moduleId, $setId, array());

    	$this->load->model('ApplicationBaseModel');
		$resApp = $this->ApplicationBaseModel->getResPoolByCompany($this->session->userdata('company'));
		if(!$resApp || !isset($resApp['ApplicationID'])) {
			$this->_errInfo = '缺少资源池业务!';
			CCLog::LogErr('缺少资源池业务! company: '. $this->session->userdata('company'));
			return false;
		}

		$this->load->model('ModuleBaseModel');
    	$resPool = $this->ModuleBaseModel->getResPoolIDByAppID($resApp['ApplicationID']);
    	if(!$resPool || !isset($resPool['ModuleID'])) {
    		$this->_errInfo = '`资源池业务`缺少`空闲机`模块!';
			CCLog::LogErr('`资源池业务`缺少`空闲机`模块! appId: '. $resApp['ApplicationID']);
			return false;
    	}

    	$hostId_suf = array();
    	foreach($result as $_r) {
    		if(!in_array($_r['HostID'], $hostId_suf)) {
    			$hostId_suf[] = $_r['HostID'];
    		}
    	}

    	/*host_2_res_empty_module：需要进入`资源池-空闲机`的主机Id*/
    	$host_2_res_empty_module = array_diff($hostId_pre, $hostId_suf);
    	if(!empty($host_2_res_empty_module)) {
	    	return $this->modMultiHostToSingleModule($host_2_res_empty_module, $resPool['ModuleID'], $resApp['ApplicationID'], false);
	    }

	    return true;
    }

    /**
     * @根据HostID更新主机信息
     * @param $hostId 主机Id
     * @param $data 主机信息
     * @return true or false
     */
    public function updateHostBaseByHostId($data, $hostId) {
    	$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
    	$this->load->database('db');
    	$this->db->where_in('HostID', $hostId);
    	$query = $this->db->get('HostBase');

    	if(!$query || $query->num_rows()==0) {
    		$this->_errInfo = '主机不存在';
    		CCLog::LogErr('主机不存在,message: '. $this->db->last_query());
    		return false;
    	}

    	$propertyNames = array_keys($data);
    	$opContent = '';
    	foreach($query->result_array() as $_h) {
    		foreach($propertyNames as $_p) {
    			$opContent .= $_p .':'. $_h[$_p] .'->'. $data[$_p] .' | ';
    		}
    	}

        $this->db->where_in('HostID', $hostId);
        $result = $this->db->update('cc_HostBase',$data);

        if(!$result) {
        	$this->_errInfo = '缺少空闲机模块!';
			CCLog::LogErr('缺少空闲机模块！appId: '. print_r($appId, true));
        }

        $log = array();
		$log['ApplicationID'] = $appId;
		$log['OpContent'] = $opContent;
		$log['OpFrom'] = 0;
		$log['OpName'] = '批量修改主机基础属性';
		$log['OpResult'] = $result ? 1 : 0;
		$log['OpTarget'] = '主机';
		$log['OpType'] = '更新';
		CCLog::addOpLogArr($log);
    }

	/**
	 * @查询所有主机数量
	 * @param
	 * @return int
	 */
	public function getHostCount() {
		$this->load->database('db');
		$this->db->from('cc_HostBase');
		$query = $this->db->get();

		return $query->num_rows();
	}

	/**
	 * @根据主机类型来源查询主机
	 * @param $Source
	 * @return int
	 */
	public function getHostBySource($Source = '') {
		$this->load->database('db');
		if($Source) {
			$this->db->where('Source', $Source);
		}

		$this->db->from('cc_HostBase');
		$query = $this->db->get();

		return $query->result_array();
	}

	/**
     * @插入主机
     * @param 主机信息
     * @return HostID
     */
    public function addHostBase($data, $appId) {	
    	if(!$data['AssetID']) {
    		$this->_errInfo = '`固资编号`不能为空!';
    		return false;
    	}
    	$this->db->select('HostID');
    	$this->db->where('AssetID', $data['AssetID']);
    	$query = $this->db->get('cc_HostBase');
    	if(!$query) {
    		$this->_errInfo = '`固资编号`['. $data['AssetID'] .']查询错误';
    		return false;
    	}

    	if($query->num_rows()>0) {
    		$result = $query->result_array();
    		$hostId = array_column($result, 'HostID');

    		$appId = is_array($appId) ? $appId : explode(',', $appId);
    		$this->db->select('ApplicationID');
    		$this->db->where_in('ApplicationID', $appId);
    		$this->db->where_in('HostID', $hostId);
    		$query = $this->db->get('cc_ModuleHostConfig');

    		if(!$query) {
    			$this->_errInfo = '`固资编号`['. $data['AssetID'] .']查询错误';
    			return false;
    		}

    		if($query->num_rows()>0) {
	    		$this->_errInfo = '`固资编号`['. $data['AssetID'] .']已存在';
	    		return false;
	    	}
    	}

        $this->db->insert('cc_HostBase',$data);
        $HostID = $this->db->insert_id();
        if(!$HostID) {
        	CCLog::LogErr('插入主机['. $data['AssetID'] .']失败, sql:'. $this->db->last_query());
    		$this->_errInfo = '插入主机['. $data['AssetID'] .']失败';
    		return false;
    	}
        return $HostID;
    }

	/**
	 * @删除主机
	 * @param $hostId 主机Id
	 * @param @appId 业务Id
	 * @return boolean 成功 or 失败
	 */
	public function deleteHostApplicationById($hostId = array(), $appId = 0) {
		$hostId 	= is_array($hostId) ? $hostId : explode(',', $hostId);
		$appId 		= intval($appId);

		if(empty($appId)) {
			$this->_errInfo = '缺少appId';
			return false;
		}

		if(empty($hostId)) {
			$this->_errInfo = 'hostId不能为空';
			return false;
		}

        $this->db->select('InnerIP,OuterIp,Source,HostID');
        $this->db->from('cc_HostBase');
        $this->db->where_in('HostID', $hostId);
        $query = $this->db->get();
        if(0 == $query->result_array()) {
            return false;
        }
        $HostInfo = $query->result_array();

		$this->db->where_in('HostID', $hostId);
		$delHostRes = $this->db->delete('cc_HostBase');

		if(!$delHostRes) {
			$this->_errInfo = '删除主机失败';
			return false;
		}else {
			$this->db->where_in('HostID', $hostId);
			$this->db->where('ApplicationID', $appId);
			$delHostAppRes = $this->db->delete('cc_ModuleHostConfig');
		}

		if(!$delHostAppRes) {
			$this->_errInfo = '删除主机业务关系失败';
			return false;
		}
		
		return true;
	}

	/**
     * @根据业务ID+IP更新VIP
     * @param $data 主机信息
     * @return true or false
     */
    public function updateVIPByHostId($data, $hostId) {
    	$hostId = is_array($hostId) ? $hostId : explode(',', $hostId);
    	$this->load->database('db');
    	$this->db->where_in('HostID', $hostId);
    	$query = $this->db->get('HostBase');

    	if(!$query || $query->num_rows() == 0) {
    		$this->_errInfo = '主机不存在';
    		CCLog::LogErr('主机不存在,message: '. $this->db->last_query());
    		return false;
    	}

        $this->db->where_in('HostID', $hostId);
        $result = $this->db->update('cc_HostBase',$data);

        if(!$result) {
        	$this->_errInfo = '更新VIP失败!';
			return false;
        }

        return true;
    }

    /**
     * @字段名称获取相应字段的下拉值
     * @param $type
     * @return array()
     */
    public function getFieldSelectByType($type) {
    	if(!in_array($type, array('DeviceClass','Region','Status','OSName'))) {
    		return array();
    	}

    	$this->load->database('db');
    	$this->db->group_by($type);
    	$this->db->select($type);
    	$query = $this->db->get('HostBase');

        $result = $query->result_array();

        if(!$result) {
        	$this->_errInfo = '没有数据!';
			return false;
        }

        return $result;
    }

	/*
     * @根据主机Id获取主机
     * @return set属性关联数组
     */
	public function getHostsByHostId($hostIdArr, $fields='HostID,AssetID,HostName,DeviceClass,Operator,BakOperator,InnerIP,OuterIP,State,ServerType,CreateTime,DeadLineTime,Memo,HardMemo,IDCArea,AutoRenew,BandWidth,OSName') {
		if(!is_array($hostIdArr)) {
			$hostIdArr = array($hostIdArr);
		}

		$this->db->select($fields);
		$this->db->where_in('HostID', $hostIdArr);
		$query = $this->db->get('HostBase');

		return $query ? $query->result_array() : array();
	}
}