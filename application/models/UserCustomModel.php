<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class UserCustomModel extends Cc_Model{
	public function __construct(){
		parent::__construct();
	}

	/**
	 * @todo 查询用户定制数据
	 * @return array() 定制信息数组
	 */
	public function getUserCustom($userName){
		$this->load->database('db');
		$this->db->where('UserName', $userName);
		$query = $this->db->get('UserCustom');

		return $query && $query->num_rows()>0 ? $query->row_array() : array();
	}

	/**
	 * 设置用户定制数据
	 * @return boolean 成功 or 失败
	 */
	public function setUserCustom($data, $userName){
		$userCustom = $this->getUserCustom($userName);

		$this->load->database();
		if(empty($userCustom)){
			$data['UserName'] = $userName;
			$addRes = $this->db->insert('UserCustom', $data);
			if(!$addRes){
				$this->_errInfo = '添加用户定制, 失败!';
				$err = $this->db->error();
				CCLog::LogErr($this->_errInfo.'mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
				return false;
			}
			return true;
		}

		$this->db->where('UserName', $userName);
		$upRes = $this->db->update('UserCustom', $data);

		if(!$upRes){
			$this->_errInfo = '更新用户定制, 失败!';
			$err = $this->db->error();
			CCLog::LogErr($this->_errInfo.'mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
			return false;
		}
		return true;
	}
}