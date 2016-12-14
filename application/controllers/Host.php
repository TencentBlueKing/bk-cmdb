<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Host extends Cc_Controller {

	public function __construct(){
		parent::__construct();
	}

	/**
	* @主机查询页面
	*/
	public function hostQuery(){
		$company = $this->session->userdata('company');
		$data = $this->buildPageDataArr($this->lang->line('host_query'), '/host/hostQuery');

		$this->load->logic('UserLogic');
		$user = $this->UserLogic->getUserList();
		$userList = array();
		foreach ($user as $user=>$name) {
			if($user && !in_array($user, array_column($userList, 'id'))) {
				$userList[] = array('id'=>$user, 'text'=>$name);
			}
		}

		$data['userList'] = json_encode($userList);

		$tablesFields = array();
		$tablesFields[] = array('field'=>'checkbox','title'>'#','menu'=>false,'width'=>30,'template'=>'<input type="checkbox" #:data.Checked# value="#:data.HostID#" class="c-grid-checkbox"/>');
		$this->load->logic('HostBaseLogic');
		$hostProperty = $this->HostBaseLogic->getHostPropertyByType('keyName');
		$hostCustomerProperty = $this->HostBaseLogic->getHostPropertyByOwner($company, 'keyName');
		$hostProperty = array_merge($hostProperty, $hostCustomerProperty);
		foreach ($hostProperty as $key => $property) {
			$temVar = array();
			$temVar['field'] = $key;
			$temVar['title'] = $property;
			$temVar['width'] = '140';
			if($key == 'InnerIP') {
				$temVar['template'] = '<a href="javascript:void(0)" class="a-innerip" title="#:data.InnerIP#">#:data.InnerIP#</a>';
			}
			$tablesFields[] = $temVar;
		}

		$data['OSName'] = $this->getHostFieldSelect('OSName');

		$appId = intval($this->input->get_post('ApplicationID'));
		if(!$appId){
			$app = $this->session->userdata('defaultApp');
			$appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
		}
		$data['appId'] = $appId;
		$this->load->logic('TopologyLogic');
		$topo = $this->TopologyLogic->getTopoTree4view($appId);
		$data['level'] = isset($topo['topo'][0]['lvl']) ? $topo['topo'][0]['lvl'] : 3;
		$data['topo'] = json_encode($topo);

		$this->load->logic('SetBaseLogic');
		$data['set_select'] = $this->SetBaseLogic->getSetById(array(), $appId);

		$this->load->logic('ModuleBaseLogic');
		$data['module_select'] = $this->ModuleBaseLogic->getModuleById(array(), array(), $appId);

		$this->load->logic('UserCustomLogic');
        $columns = $this->UserCustomLogic->getUserCustomColumn();
        $columns = $columns ? $columns : array('InnerIP','OuterIP', 'SetName', 'ModuleName');
		$data['columns'] = json_encode($columns);
		$data['DefaultField'] = $this->UserCustomLogic->getUserCustomDefaultField();
		$data['InnerIP'] = '';
		$innerIp = trim(trim($this->input->get_post('InnerIP')), ',');
		if($innerIp != NULL && strlen($innerIp)>0){
			$data['InnerIP'] = $innerIp;
		}

		$data['OuterIP'] = '';
		$outerIp = trim(trim($this->input->get_post('OuterIP')), ',');
		if($outerIp != NULL && strlen($outerIp)>0){
			$data['OuterIP'] = $outerIp;
		}

		$data['AssetID'] = '';
		$assetID = trim(trim($this->input->get_post('AssetID')), ',');
		if($assetID != NULL && strlen($assetID)>0){
			$data['AssetID'] = $assetID;
		}

        $hostFiles = array();
        $necFields = array('InnerIP','OuterIP','SetName','ModuleName');
        foreach($hostProperty as $key=>$hf){
            if(!in_array($key,$necFields)){
                $hostFiles[$key] = $hf;
            }
        }
        $data['HostFields'] = $hostFiles;
        $data['customerQueryFields'] = $hostFiles;
        $data['NecFields'] = $necFields;
        $data['hostPropertyField'] = $hostProperty;
        $data['tablesFields'] = json_encode($tablesFields);
		$this->load->library('layout');
		$this->layout->view('host/hostQuery', $data);
	}

	/**
	* @主机字段下拉
	*/
	private function getHostFieldSelect($type) {
		$list = array();
		$this->load->logic('HostBaseLogic');
		$name = $this->HostBaseLogic->getFieldSelectByType($type);
		foreach($name as $key=>$value) {
			if($value[$type]) {
				$list[] = array('id'=>$value[$type], 'text'=>$value[$type]);
			}
		}
		
		return json_encode($list);
	}

	/**
	* @主机查询页面，树结构
	*/
	public function getTopoTree4view(){
		$data = array();
		$appId = intval($this->input->get_post('ApplicationID'));
		if(!$appId) {
			$app = $this->session->userdata('defaultApp');
			$appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
		}

		$this->load->logic('TopologyLogic');
		$result = json_encode($this->TopologyLogic->getTopoTree4view($appId));

		return $this->output->set_output($result);
	}

	/**
	* @主机查询页面，点击树节点，查询主机
	*/
	public function getHostById() {
		$appId = $this->input->get_post('ApplicationID');
		if(!$appId) {
			$app = $this->session->userdata('defaultApp');
			$appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
			if(!$appId) {
				return $this->outputDataJson();
			}
		}
		$setId = $this->input->get_post('SetID');
		$moduleId = $this->input->get_post('ModuleID');

		$page = intval($this->input->get_post('page'));
		$pageSize = intval($this->input->get_post('pageSize'));

		$expires = time() + $this->config->item('sess_expiration') * 100;
		$condition = array();
		$appId && $condition['ApplicationID'] = $appId;
		$setId && $condition['SetID'] = $setId;
		$moduleId && $condition['ModuleID'] = $moduleId;

		$start = ($page-1) * $pageSize;
		$limit = $pageSize;

        $this->load->logic('HostBaseLogic');
        $total = $this->HostBaseLogic->getHostCountById(array(), $moduleId, $setId, $appId);
        $hostInfo = array();
        if($total > 0) {
            $hosts = $this->HostBaseLogic->getHostById(array(), $moduleId, $setId, $appId);
            $hostInfo = $this->HostBaseLogic->getHostById(array_column($hosts, 'HostID'), array(), array(), array(), $start, $limit);
        }

		$parameterData = array();
		$this->load->logic('BaseParameterDataLogic');
		$parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();

		#用户key=>value#
		$this->load->logic('UserLogic');
        $userListKv = $this->UserLogic->getUserList();

		foreach($hostInfo as $key=>$host) {
			if(isset($parameterData[$host['Source']])) {
				$hostInfo[$key]['Checked'] = '';
				$hostInfo[$key]['Source'] = $parameterData[$host['Source']];
			}

			if(isset($userListKv[$host['Operator']])) {
				$hostInfo[$key]['Operator'] = $this->getUserDisplayName($userListKv, $hostInfo[$key]['Operator']);
			}

			if(isset($userListKv[$host['BakOperator']])) {
				$hostInfo[$key]['BakOperator'] = $this->getUserDisplayName($userListKv, $hostInfo[$key]['BakOperator']);
			}
 		}

		return $this->outputDataJson($hostInfo, $total);
	}

	/**
	* @主机查询页面，通用查询
	*/
	public function getHostByCondition() {
		$this->load->logic('HostBaseLogic');
		$condition = array();
		$appId = $this->input->get_post('ApplicationID');
		if(!$appId) {
			$app = $this->session->userdata('defaultApp');
			$appId[] = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
			if(!$appId) {
				return $this->outputDataJson();
			}
		}
		$condition['ApplicationID'] = $appId;

		$setId = $this->input->get_post('SetID');
		$setId && $condition['SetID'] = $setId;

		$moduleId = $this->input->get_post('ModuleID');
		$moduleId && $condition['ModuleID'] = $moduleId;

		$innerIp = $this->input->get_post('InnerIP');
		$innerIp && $condition['InnerIP'] = explode(',', trim(trim($innerIp), ','));

		$outerIp = $this->input->get_post('OuterIP');
		$outerIp && $condition['OuterIP'] = explode(',', trim(trim($outerIp), ','));

		$ifOuterexact = $this->input->get_post('IfOuterexact');
		$ifOuterexact && $condition['IfOuterexact'] = trim($ifOuterexact);

		$ifInnerIPexact = $this->input->get_post('IfInnerIPexact');
		$ifInnerIPexact && $condition['IfInnerIPexact'] = trim($ifInnerIPexact);

		$hostProperty = $this->HostBaseLogic->getHostPropertyByType('keyName');
		$company = $this->session->userdata('company');
		$hostCustomerProperty = $this->HostBaseLogic->getHostPropertyByOwner($company, 'keyName');
		$hostProperty = array_merge($hostProperty, $hostCustomerProperty);
        $filpFields = array('InnerIP','OuterIP','SetName','ModuleName');
        foreach($hostProperty as $key=>$value) {
            if(!in_array($key, $filpFields)) {
                $stdProperty = $this->input->get_post($key);
				$stdProperty && $condition[$key] = explode(',', trim(trim($stdProperty), ','));
            }
        }
		
		$page = intval($this->input->get_post('page'));
		$pageSize = intval($this->input->get_post('pageSize'));
		$start = ($page-1) * $pageSize;
		$limit = $pageSize;
		$expires = time()+$this->config->item('sess_expiration') * 100;
		$total = $this->HostBaseLogic->getHostCountByCondition($condition, true);
		$hostInfo = array();
		if($total > 0) {
			$hosts = $this->HostBaseLogic->getHostByCondition($condition, true);
			$hostInfo = $this->HostBaseLogic->getHostById(array_column($hosts, 'HostID'), array(), array(), array(), $start, $limit);
		}

		$parameterData = array();
		$this->load->logic('BaseParameterDataLogic');
		$parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();
		#用户key=>value#
		$this->load->logic('UserLogic');
        $userListKv = $this->UserLogic->getUserList();
		foreach($hostInfo as $key=>$host) {
			if(isset($parameterData[$host['Source']])) {
				$hostInfo[$key]['Checked'] = '';
				$hostInfo[$key]['Source'] = $parameterData[$host['Source']];
			}

			if(isset($userListKv[$host['Operator']])) {
				$hostInfo[$key]['Operator'] = $this->getUserDisplayName($userListKv, $hostInfo[$key]['Operator']);
			}

			if(isset($userListKv[$host['BakOperator']])) {
				$hostInfo[$key]['BakOperator'] = $this->getUserDisplayName($userListKv, $hostInfo[$key]['BakOperator']);
			}
		}
		
		return $this->outputDataJson($hostInfo, $total);
	}

	/**
	* @主机查询页，转移选中主机
	*/
	public function modHostModule() {
		$appId = intval($this->input->get_post('ApplicationID'));
		$moduleId = $this->input->get_post('ModuleID');
		$hostId = $this->input->get_post('HostID');
		$isIncrement = !!$this->input->get_post('IsIncrement');//默认覆盖更新

		if(!$appId) {
			return $this->outputJson(false, 'lack_of_application_id');
		}

		if(!$moduleId) {
			return $this->outputJson(false, 'lack_of_module_id');
		}

		if(!$hostId) {
			return $this->outputJson(false, 'lack_of_host_id');
		}

		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->modHostModule($hostId, $moduleId, $appId, $isIncrement);
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}

		return $this->outputJsonByMessage(true, $this->lang->line('transfor_success'));
	}

	/**
	* @主机查询页，删除选中主机
	*/
	public function delHostModule() {

		$appId = intval($this->input->get_post('ApplicationID'));
		$setId = $this->input->get_post('SetID');
		$moduleId = $this->input->get_post('ModuleID');
		$hostId = $this->input->get_post('HostID');

		if(!$appId) {
			return $this->outputJson(false, 'lack_of_application_id');
		}

		$type = '';
		if((($type='setId') && !$setId) && (($type='moduleId') && !$moduleId) && (($type='hostId') && !$hostId)) {
			return $this->outputJsonByMessage(false, $this->lang->line('lack_of') . $type);
		}

		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->deleteHost($hostId, $moduleId, $setId, $appId);
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}

		return $this->outputJsonByMessage(true, $this->lang->line('transfor_empty_module_success'));
	}

	/**
	* @主机查询页，上缴选中主机
	*/
	public function resHostModule() {
		$this->load->logic('ApplicationBaseLogic');
		$appId = intval($this->input->get_post('ApplicationID'));
		$setId = $this->input->get_post('SetID');
		$moduleId = $this->input->get_post('ModuleID');
		$hostId = $this->input->get_post('HostID');

		if(!$appId) {
			return $this->outputJson(false, 'lack_of_application_id');
		}

		$type = '';
		if((($type='setId') && !$setId) && (($type='moduleId') && !$moduleId) && (($type='hostId') && !$hostId)) {
			return $this->outputJsonByMessage(false, $this->lang->line('lack_of') . $type);
		}

		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->resHostModule($hostId, $moduleId, $setId, $appId);
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}

		$appId = $this->session->userdata('appId');
		$appHostCount = array();
		$appHostCount = $this->ApplicationBaseLogic->getHostNumByAppIdArr($appId);
		$app = $this->session->userdata('app');
		$sessionData = array();
		foreach ($app as $key=>$value) {
			$sessionData['app'][$value['ApplicationID']]['ApplicationID'] = $value['ApplicationID'];
			$sessionData['app'][$value['ApplicationID']]['ApplicationName'] = $value['ApplicationName'];
            $sessionData['app'][$value['ApplicationID']]['CompanyID'] = $this->session->userdata('company_id');
            $sessionData['app'][$value['ApplicationID']]['Owner'] = $this->session->userdata('company');
			$sessionData['app'][$value['ApplicationID']]['ApplicationHostCount'] = $appHostCount[$value['ApplicationID']];
		}
		$this->session->set_userdata($sessionData);

		return $this->outputJsonByMessage(true, $this->lang->line('turn_in_success'));
	}

	/**
	* @主机查询页，修改选中主机，修改主机信息
	*/
	public function updateHostInfo() {
		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->updateHost();
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}
		
		return $this->outputJsonByMessage(true, $this->lang->line('modfiy_success'));
	}

	/**
	 * @host detail Page for this controller.
	 */
	public function details() {
		$host = array();
		$company = $this->session->userdata('company');
		$this->load->logic('HostBaseLogic');
		$appId = intval($this->input->get_post('ApplicationID'));
		if(!$appId) {
			$app = $this->session->userdata('defaultApp');
			$appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
			if(!$appId) {
				return $this->outputJson(false, 'lack_of_application_id');
			}
		}

		$condition = array();
		$hostID = $this->input->get_post('HostID');
		$assetID = $this->input->get_post('AssetID');
		$innerIP = $this->input->get_post('InnerIP');
		$outerIP = $this->input->get_post('OuterIP');
		if(isset($hostID)) {
			$hostId = intval($this->input->get_post('HostID'));
		}elseif(isset($assetID) || isset($innerIP) || isset($outerIP)) {
			$condition['ApplicationID'] = $appId;
			if($assetID) {
				$condition['AssetID'] = explode(',', $assetID);
			}

			if($innerIP) {
				$condition['InnerIP'] = explode(',', $innerIP);
			}

			if($outerIP) {
				$condition['OuterIP'] = explode(',', $outerIP);
			}

			$this->load->logic('HostBaseLogic');
			$host = $this->HostBaseLogic->getHostByCondition($condition);
			if(count($host) > 0) {
				$hostId = $host[0]['HostID'];
			}else {
				return $this->outputJson(false, 'app_no_host');
			}
		}

		if(!$hostId) {
			return $this->outputJson(false, 'lack_of_host_id');
		}

		$hostProperty = $this->HostBaseLogic->getHostPropertyByType('groupKey');
		$hostPropertykeyName = $this->HostBaseLogic->getHostPropertyByType('keyName');
		$hostCustomerProperty = $this->HostBaseLogic->getHostPropertyByOwner($company , 'groupKey');
		$hostCustomerkeyName = $this->HostBaseLogic->getHostPropertyByOwner($company , 'keyName');
		$hostPropertykeyName = array_merge($hostPropertykeyName, $hostCustomerkeyName);
		$host = $host ? $host : $this->HostBaseLogic->getHostById($hostId, array(), array(), $appId);
		$parameterData = array();
		$this->load->logic('BaseParameterDataLogic');
		$parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();
		$header = $hostPropertykeyName;
		$this->load->logic('UserCustomLogic');
		$userCustomCon = $this->UserCustomLogic->getUserCustom();
		if(isset($userCustomCon['DefaultCon']) && $userCustomCon['DefaultCon']) {
			$userCustomCon = json_decode($userCustomCon['DefaultCon'], true);
		}else {
			$userCustomCon = array_keys($header);
			$data['DefaultCon'] = json_encode($userCustomCon);
			$this->UserCustomLogic->setUserCustomByUserName($data);
		}

		#用户key=>value#
		$this->load->logic('UserLogic');
        $userListKv = $this->UserLogic->getUserList();
		
		$data = array();
		$hostInfo = array();
		$i = $k = 0;
		$hostInfo['ccattribute'] = array();
		foreach ($hostPropertykeyName as $key => $value) {
			if(in_array($key, $hostProperty['basic'])) {
				$basis = array();
				$basis['id'] = $i;
				$basis['text'] = $header[$key];
				if(in_array($key, $userCustomCon)) {
					$basis['selected'] = true;
				}
				
				#主机展示的字段翻译#
				if($key == 'Source') {
					$basis['content'] = isset($parameterData[$host[0][$key]]) ? $parameterData[$host[0][$key]] : $host[0][$key];
				}elseif($key == 'AutoRenew') {
					$basis['content'] = $host[0][$key] ? $this->lang->line('have_open') : $this->lang->line('no_open');
				}elseif($key == 'Operator' || $key == 'BakOperator') {
					$basis['content'] = $this->getUserDisplayName($userListKv, $host[0][$key]);
				}else {
					$basis['content'] = $host[0][$key];
				}

				if(in_array($key, array('InnerSwitchPort', 'OuterSwitchPort', 'AssetID', 'HardMemo', 'OSName', 'Description'))) {
					$basis['tips'] = true;
				}
				$basis['key'] = $key;
				$hostInfo['basis'][] = $basis;
				$i += 1;
			}

			if(in_array($key, $hostCustomerProperty['Customer'])) {
				$ccattribute = array();
				$ccattribute['id'] = $k;
				$ccattribute['text'] = $header[$key];
				if(in_array($key, $userCustomCon)) {
					$ccattribute['selected'] = true;
				}
				$ccattribute['content'] = $host[0][$key];
				$ccattribute['key'] = $key;
				$hostInfo['ccattribute'][] = $ccattribute;
				$k += 1;
			}
		}
		
		$data['host'] = json_encode($hostInfo);
		$data['parameterData'] = $parameterData;
        $this->load->view('host/details', $data);
	}

	/**
	 * @设置用户主机详情展示字段
	 * @return json
	 */
	public function setDefaultField() {
		$key = $this->input->get_post('key',true);
		$type = $this->input->get_post('type',true);
		$field = $this->input->get_post('field',true);
		if(!$field) {
			$field = 'DefaultCon';
		}
		$this->load->logic('UserCustomLogic');
		$userCustomCon = $this->UserCustomLogic->getUserCustom();
		$userCustomCon = json_decode($userCustomCon[$field], true);
		if(!$userCustomCon) {
			$userCustomCon = array('InnerIP','OuterIP', 'SetName', 'ModuleName');
		}
		if($type == 'a') {
			if(!in_array($key, $userCustomCon)) {
				$userCustomCon[] = $key;
			}
		}elseif($type == 'd') {
			if(in_array($key, $userCustomCon)) {
				unset($userCustomCon[array_search($key, $userCustomCon)]);
			}
		}
		$data[$field] = json_encode($userCustomCon);
		$this->UserCustomLogic->setUserCustomByUserName($data);
		return $this->outputJsonByMessage(true, $this->lang->line('setint_success'));
	}

	/**
	* @快速导入页面
	*/
	public function quickImport() {

		$company = $this->session->userdata('company');
		$this->load->logic('ApplicationBaseLogic');
		$apps = $this->ApplicationBaseLogic->getAppByCompany();

		$appId = array_column($apps, 'ApplicationID');
		$this->load->logic('HostBaseLogic');
		$host = $this->HostBaseLogic->getHostById(array(), array(), array(), $appId);
		$parameterData = array();
		$this->load->logic('BaseParameterDataLogic');
		$parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();
		$data = $this->buildPageDataArr($this->lang->line('quick_import'), '/host/quickImport');
		
		$data['company'] = $company;
		$this->load->library('layout');
		$this->layout->view('host/quickImport', $data);
	}

	/**
	* @快速导入页面查询主机
	*/
	public function getHost4QuickImport() {

		$company = $this->session->userdata('company');
		$this->load->logic('ApplicationBaseLogic');
		$apps = $this->ApplicationBaseLogic->getAppByCompany();

		if(count($apps)<2) {
			header('location:'.BASE_URL.'/app/index');
			exit;
		}
		
		$isDistributed = $this->input->get_post('IsDistributed')==='true' ? true : false;
		$condition['ApplicationID'] = array_column($apps, 'ApplicationID');
		$this->load->logic('HostBaseLogic');
		$hosts = $this->HostBaseLogic->getHostByCondition($condition, true);

		$hostList = array();
		foreach($hosts as $_h) {
			$_h['Checked'] = '';
			if($isDistributed) {
				if($_h['ApplicationName'] !== $this->lang->line('pools')) {
					$hostList[] = $_h;
				}
			}else {
				if($_h['ApplicationName'] === $this->lang->line('pools')) {
					$hostList[] = $_h;
				}
			}
		}

		$result = array();
		$result['success'] = true;
		$result['data'] = $hostList;
		$result['total'] = count($hostList);
		return $this->output->set_output(json_encode($result));
	}

	/**
	* @excel导入主机
	*/
	public function importPrivateHostByExcel() {
		set_time_limit(0);
		$this->load->logic('ImportLogic');
		$result = $this->ImportLogic->importPrivateHostByExcel();
		
		$data = array();
		$data['success'] = !!$result;
		$data['errInfo'] = $this->ImportLogic->_errInfo;
		$data['name'] = 'importToCC';
		return $this->output->set_output('<script>parent.window.uploadCallback('.json_encode($data).');</script>');
	}

	/**
	* @excel导入私有云主机到CC时读取表头
	*/
	public function getImportPrivateHostTableFieldsByExcel() {
		set_time_limit(0);
		$this->load->logic('ImportLogic');
		$data = $this->ImportLogic->getImportPrivateHostTableFieldsByExcel();
		if(!$data['success']) {
			$data['errInfo'] = $this->config->item($data['errCode'])->Info;
		}
		
		return $this->output->set_output('<script>parent.window.uploadCallbackToHostField('.json_encode($data).');</script>');
	}

	/**
	* @快速导入页，转移选中主机
	*/
	public function quickDistribute() {
		set_time_limit(0);
		$this->load->logic('ApplicationBaseLogic');

		$appId = $this->input->get_post('ApplicationID');
		$toAppId = intval($this->input->get_post('ToApplicationID'));
		$moduleId = $this->input->get_post('ModuleID');
		$hostId = $this->input->get_post('HostID');

		if(!$appId) {
			return $this->outputJson(false, 'lack_of_application_id');
		}

		if(!$moduleId) {
			$this->load->logic('ModuleBaseLogic');
			$module = $this->ModuleBaseLogic->getModuleByName($this->lang->line('empty_pools'), array(), $toAppId);
			if(!$module) {
				return $this->outputJson(false, 'lack_of_module_id');
			}

			$moduleId = $module[0]['ModuleID'];
		}

		if(!$hostId) {
			return $this->outputJson(false, 'lack_of_host_id');
		}

		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->quickDistribute($hostId, $moduleId, $toAppId);
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}

		$appId = $this->session->userdata('appId');
		$appHostCount = $this->ApplicationBaseLogic->getHostNumByAppIdArr($appId);
		$app = $this->session->userdata('app');
		$sessionData = array();
		foreach ($app as $key=>$value) {
			$sessionData['app'][$value['ApplicationID']]['ApplicationID'] = $value['ApplicationID'];
			$sessionData['app'][$value['ApplicationID']]['ApplicationName'] = $value['ApplicationName'];
			$sessionData['app'][$value['ApplicationID']]['ApplicationHostCount'] = $appHostCount[$value['ApplicationID']];
		}
		$this->session->set_userdata($sessionData);
		return $this->outputJsonByMessage(true, $this->lang->line('transfor_success'));
	}

	/**
	* @主机导出
	*/
	public function hostExport() {
		set_time_limit(0);

		$this->load->logic('ExportLogic');
		$result = $this->ExportLogic->exportHostToExcel();
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->ExportLogic->_errInfo);
		}
	}

	/**
	 * @删除私有云空闲机（默认业务）下主机
	 */
	public function delPrivateDefaultApplicationHost() {
		$appId = $this->input->get_post('ApplicationID');
		$hostId = $this->input->get_post('HostID');

		if(!$appId) {
			return $this->outputJson(false, 'lack_of_application_id');
		}

		$this->load->logic('ApplicationBaseLogic');
		$app = $this->ApplicationBaseLogic->getAppById($appId);
		if(count($app) != 1 || !isset($app[0]['Default']) || $app[0]['Default'] != 1) {
			return $this->outputJson(false, 'no_default_app_cannot_delete_host');
		}

		if(!$hostId){
			return $this->outputJson(false, 'lack_of_host_id');
		}

		$host = array();
		$this->load->logic('HostBaseLogic');
		$host = $this->HostBaseLogic->getHostById($hostId);
		foreach($host as $value) {
			if($value['Source'] == 1) {
				return $this->outputJson(false, 'cannot_delete_host');
			}
		}

		$hostApp = array();
		$this->load->logic('HostBaseLogic');
		$hostApp = $this->HostBaseLogic->getHostAppGroupModuleRealtionByHostIDs($hostId);

		$appId = is_array($appId) ? $appId : explode(',', $appId);
		$appId = array_unique($appId);
		$hostAppId = array_unique(array_column($hostApp, 'ApplicationID'));
		
		if(array_diff($appId, $hostAppId) || array_diff($hostAppId, $appId)) {
			return $this->outputJson(false, 'delete_host_belong_to_other');
		}

		$this->load->logic('HostBaseLogic');
		$result = $this->HostBaseLogic->RealDeleteHost($hostId, $appId);
		if(!$result) {
			return $this->outputJsonByMessage(false, $this->HostBaseLogic->_errInfo);
		}
		
		return $this->outputJsonByMessage(true, $this->lang->line('delete_success'));
	}
}