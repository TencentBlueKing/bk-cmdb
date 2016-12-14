<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class OperationLogModel extends Cc_Model {
	public function __construct(){
		parent::__construct();
	}

	/**
	 * @获取操作日志
	 * @return array
	 */
	public function getOperationLog($operator='', $appId=0, $opTime='', $opType='', $opTarget='', $opContent='', $clientIp='', $start=0, $limit=0, $orderBy='', $direction='DESC'){
		if($operator!=''){
			$this->db->where('Operator', $operator);
		}

		if($appId!=''){
			$this->db->where('ApplicationID', $appId);
		}

		if($this->session->userdata('company_code')!=''){
			$this->db->where('CompanyCode', $this->session->userdata('company_code'));
		}

		if($opTime!=''){
			$time = strtotime($opTime);
			$this->db->where('opTime >=', $time);
			$this->db->where('opTime <=', $time+86400);
		}

		if($opType!=''){
			$this->db->where('OpType', $opType);
		}

		if($opTarget!=''){
			$this->db->where('OpTarget', $opTarget);
		}

		if($opContent!=''){
			$this->db->like('OpContent', $opContent);
		}

		if($clientIp!=''){
			$this->db->where('ClientIP', $clientIp);
		}

		if($orderBy!=''){
			$this->db->order_by($orderBy, $direction);
		}

		if($limit>0){
			$this->db->limit($limit, $start);
		}

		$query = $this->db->get('OperationLog');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

	/**
	 * @获取操作用户数
	 * @param $webSys
	 * @param $start
	 * @param $end
	 * @return int
	 */
	public function getOperationUserCount($webSys='', $start='', $end=''){

		if($webSys) {
			$this->db->where("WebSys", $webSys);
		}

		if($start) {
			$this->db->where("OpTime >=", $start);
		}

		if($end) {
			$this->db->where("OpTime <=", $end);
		}

		$this->db->group_by("Operator");
		$query = $this->db->get('cc_OperationLog');
		return $query && $query->num_rows()>0 ? $query->num_rows() : 0;
	}

	/**
	 * @获取用户操作日志
	 * @return array
	 */
	public function getUserOperationLog($operator='', $appId=0, $opType='', $opTarget='', $opContent='', $clientIp='', $starttime='', $endtime='',$start=0, $limit=0, $orderBy='', $direction='DESC'){
		if($operator!=''){
			$this->db->like('Operator', $operator);
		}

		if($appId!=''){
			$this->db->where('ApplicationID', $appId);
		}

		if($starttime!=''){
			$this->db->where('opTime >=', $starttime);
		}

		if($endtime!=''){
			$this->db->where('opTime <=', $endtime . ' 23:59:59');
		}

		if($opType!=''){
			$this->db->like('OpType', $opType);
		}

		if($opTarget!=''){
			$this->db->like('OpTarget', $opTarget);
		}

		if($opContent!=''){
			$this->db->like('OpContent', $opContent);
		}

		if($clientIp!=''){
			$this->db->where('ClientIP', $clientIp);
		}

		if($orderBy!=''){
			$this->db->order_by($orderBy, $direction);
		}else{
			$this->db->order_by('OpTime', $direction);
		}

		if($limit>0){
			$this->db->limit($limit, $start);
		}

		$query = $this->db->get('OperationLog');
		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}
}