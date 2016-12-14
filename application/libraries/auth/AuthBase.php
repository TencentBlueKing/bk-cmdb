<?php if(!defined('BASEPATH')) exit('No direct script access allowed');

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class AuthBase{

	public $_ciObj = NULL;

	public function __construct() {
		$this->_ciObj = &get_instance();
		$this->_ciObj->load->library('session');
	}

	/*
	 * @校验集群
	 */
	public function validSetBase($setId, $appId) {
		$this->_ciObj->load->logic('SetBaseLogic');
		$set = $this->_ciObj->SetBaseLogic->getSetById($setId, $appId);
		return !empty($set);
	}

	/*
	 * @校验模块
	 */
	public function validModuleBase($moduleId, $appId) {
		$this->_ciObj->load->logic('ModuleBaseLogic');
		$module = $this->_ciObj->ModuleBaseLogic->getModuleById($moduleId, array(), $appId);
		return !empty($module);
	}

	/*
	 * @校验主机Id
	 */
	public function validHostId($hostId, $appId) {
		$this->_ciObj->load->logic('HostBaseLogic');
		$host = $this->_ciObj->HostBaseLogic->getHostById($hostId, array(), array(), $appId);
		return !empty($host);
	}

	/*
	 * @校验主机IP
	 */
	public function validHostIp($hostIp, $appId) {
		$this->_ciObj->load->logic('HostBaseLogic');
		$host = $this->_ciObj->HostBaseLogic->getHostByIp($hostIp, array(), array(), $appId);
		return !empty($host);
	}
}