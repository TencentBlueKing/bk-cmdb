<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

include_once(APPPATH.'libraries/auth/AuthBase.php');

class AuthWeb extends AuthBase{

	private $_params = array();

	private static $_whiteList = array('/', '/app/add', '/host/hostquery', '/host/importprivatehostbyexcel', '/host/importrendhostbyexcel','/host/quickimport', '/host/quickbuy', '/topology/rsyncfromusp','/app/sycproject','/host/setdefaultfield','/welcome/setdefaultcom');

	private $_controller = '';

	private $_action = '';

	public function __construct() {
		parent::__construct();
		$segments = $this->_ciObj->uri->rsegments;
		$this->_controller = $segments[1];
		$this->_action = $segments[2];

//        $this->checkSysPermissions();  //系统有权限访问的cgi

	}

	/*
	 * @校验token,防csrf
	 */
    private function authtoken()
    {
        if($this->_ciObj->input->is_ajax_request()) {
            $sessionToken = $this->_ciObj->session->userdata('token');
            $cookieToken = $this->_ciObj->input->get_request_header('Token', TRUE);
            if(!isset($sessionToken) || !isset($cookieToken) || $sessionToken != $cookieToken) {
                header('Timeout:true');
				die(json_encode(array('success'=>false)));
            }
        }
    }

	/*
	 * @校验系统权限
	 */
    private function checkSysPermissions() {
        $this->_ciObj->load->database();
        $urlArr = array_reverse($this->_ciObj->uri->segment_array());

        if(empty($UrlArr)) {
            $controller = 'index';
            $action = 'index';
        } else if(1 === count($UrlArr)) {
			$controller = strtolower($UrlArr[0]);
            $action = 'index';
        } else {
            $action = strtolower($UrlArr[0]);
            $controller = strtolower($UrlArr[1]);
        }

        $this->_ciObj->db->where('Controller', $controller);
        $this->_ciObj->db->where('Action', $action);

        $nums = $this->_ciObj->db->get('cc_SysPermissions')->num_rows();
        if(0===$nums) {
            show_404();
        }
    }

	/*
	 * @初始化参数
	 */
	private function initParams(){
		foreach($_GET as $_k=>$_v) {
			$this->setParams($_k, $_v);
		}
		foreach($_POST as $_k=>$_v) {
			$this->setParams($_k, $_v);
		}

        if(empty($_GET) && empty($_POST )) {
            return true;
        }
		$this->auth();
	}

	/*
	 * @鉴权入口
	 */
	private function auth() {
		if($this->doAuth()) {
			return true;
		}
		if ($this->_ciObj->input->is_ajax_request()) {
			$result = array();
			$result['success'] = false;
			$result['errCode'] = $this->_ciObj->config->item('auth_no_permission')->Code;
			$result['errInfo'] = $this->_ciObj->config->item('auth_no_permission')->Info;
			exit(json_encode($result));
		}

		header('Location:/');         //框架测试阶段，暂时关闭跳转
		exit;
	}

	/*
	 * @鉴权实体函数
	 */
	private function doAuth(){
		$appId = $this->_ciObj->session->userdata('appId') ? $this->_ciObj->session->userdata('appId') : array();
		$appId = is_array($appId) ? $appId : explode(',',$appId);
		if($this->isInWhiteList()){
			return true;
		}elseif(!$this->getParams('ApplicationID') || array_diff(explode(',', $this->getParams('ApplicationID')), $appId)){
			return false;
		}

		/*
		 *注意不要随意变换if代码段的先后顺序。可能造成鉴权绕过。
		*/

		$keys = array_keys($this->getParams());
		
		if(in_array('HostID', $keys)){
			return $this->validHostId($this->getParams('HostID'), $this->getParams('ApplicationID'));
		}

		if(in_array('HostIP', $keys)){
			return $this->validHostIp($this->getParams('HostIP'), $this->getParams('ApplicationID'));
		}

		if(in_array('ModuleID', $keys)){
			return $this->validModuleBase($this->getParams('ModuleID'), $this->getParams('ApplicationID'));
		}

		if(in_array('SetID', $keys)){
			return $this->validSetBase($this->getParams('SetID'), $this->getParams('ApplicationID'));
		}

		return true;
	}

	/*
	 * @获取key对应的值
	 */
	private function getParams($item='') {
		if(!$item){
			return $this->params;
		}
		return isset($this->params[$item]) ? $this->params[$item] : false;
	}

	/*
	 * @设置key对应的值
	 */
	private function setParams($_k, $_v) {
		return $this->params[$_k] = $_v;
	}

	/*
	 * @是否不验证业务Id
	 */
	private function isInWhiteList() {
		return in_array(strtolower('/'.$this->_controller.'/'.$this->_action), self::$_whiteList);
	}
}