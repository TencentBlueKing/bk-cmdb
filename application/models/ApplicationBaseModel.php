<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ApplicationBaseModel extends Cc_Model {

	public function __construct(){
		parent::__construct();
	}

    /*
     *根据业务ID查询业务
     * @param 业务ID数组
     * @return 业务数组
     */
	public function getAppById($appId = array(), $fields='*') {
		if(!is_array($appId)){
			if(strpos(',', $appId) !== false){
				$appId = explode(',', $appId);
			}else{
				$appId = (array)intval($appId);
			}
		}
        $this->db->select($fields);
		$this->db->where_in('ApplicationID', $appId);
		$query = $this->db->get('ApplicationBase');

		return $query && $query->num_rows()>0 ? $query->result_array() : array();
	}

    /*
     * 查询有权限的业务
     * @param   userName,company
     * @return  业务数组
     */
	public function getAppIdOnLogin($userName, $fields = 'ApplicationID') {
        $userLike = '_'.$userName.'_';
		$this->db->select($fields);
        $this->db->where('Creator',$userName);  //userName创建的业务

        $this->db->or_group_start();
        $this->db->where('Default',1);
        $this->db->group_end();

        $this->db->or_group_start();
		$this->db->like('Maintainers', $userLike);  //用户有业务运维权限的业务
        $this->db->group_end();
		$query = $this->db->get('ApplicationBase');

        return $query->num_rows() ? $query->result_array() : array();
	}


    /*
     * @添加默认业务
     * @param company
     * @return 业务ID
     */
	public function addDefaultApp($company) {
		$this->db->select('ApplicationID');
		$this->db->where('Owner', $company);
        $this->db->where('Default', 1);
		$query = $this->db->get('ApplicationBase');

        $companyId = $this->session->userdata('company_id');
        $companyId = is_null($companyId) ? 0 : $companyId;
		if( $query && $query->num_rows()==0) {
			$query->free_result();
			$data = array();
			$data['ApplicationName'] = DEFAULT_APP_NAME;
			$data['Creator'] = $company;
			$data['CreateTime'] = date('Y-m-d H:i:s');
			$data['Default'] = 1;
			$data['Display'] = 1;
			$data['Level'] = 2;
			$data['Maintainers'] = $company;
			$data['Owner'] = $company;
			$data['Type'] = 0;
            $data['Source'] = 0;
            $data['CompanyID'] = $companyId;

			$result = $this->db->insert('ApplicationBase', $data);

			if(!$result){
				$this->_errInfo = '添加默认业务失败!';
				$err = $this->db->error();
				CCLog::LogErr('添加默认业务失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
			}
			return $this->db->insert_id();
		}

		$this->_errInfo = '默认业务已存在!';
		return false;
	}

    /**
     * @查询资源池
     * @param company
     * @return 业务数组
     */
    public function getResPoolByCompany($company, $fields='*') {
        $this->db->select($fields);
        $this->db->from('ApplicationBase');
        $this->db->where('Owner',$company);
        $this->db->where('Default',1);
        $query = $this->db->get();
        return $query->num_rows() ? $query->row_array() : array();
    }

    /**
     * @func 删除业务
     * @param appId
     * @return 业务数组
     */
    public function deleteApp($appId) {
        $this->db->where('ApplicationID',$appId);
        $this->db->delete('ApplicationBase');
    }

    /**
     * @func 查询是否有同名业务
     * @param appId
     * @return 业务数组
     */
    public function getAppByAppNameAndCompany($appName) {
        $this->db->where('ApplicationName', $appName);
        $query = $this->db->get('ApplicationBase');
        return $query->num_rows() ? $query->result_array() : false;
    }

    /**
     * @查询同名业务名是否可用
     * @param appId
     * @return 业务数组
     */
    public function getAppByAppNameAndAppID($appName, $appId) {
        $this->db->where('ApplicationName', $appName);
        $this->db->where('ApplicationID <>', $appId);
        $result = $this->db->get('ApplicationBase');
        return $result->num_rows() ? false : true ;
    }

    /*
     * 新增业务
     * @param 业务数组
     * @return 业务Id
     */
    public function addApplication($data) {
        $this->db->insert('ApplicationBase', $data);
        $appId = $this->db->insert_id();
        $log = array();
        $log['ApplicationID'] = $appId;
        $log['OpType']='新增';
        $log['OpTarget']='业务';
        $log['OpContent'] = '业务名：['.$data['ApplicationName'].']';
        CCLog::addOpLogArr($log);
        return $appId;
    }

    /**
     * @更新业务信息
     * @param data,appName,maintainers
     * @return 无
     */
    public function editApplication($data, $appName, $maintainers) {
        $this->db->where('ApplicationID', $data['ApplicationID']);
        $this->db->update('ApplicationBase', $data);
        $log= array();
        $log['OpType'] = '更新';
        $log['OpTarget'] = '业务';
        $log['ApplicationID'] = $data['ApplicationID'];
        $log['OpContent'] = '[业务名'.$appName.'->'.$data['ApplicationName'].'][运维列表'.str_replace('_','',$maintainers).'->'.str_replace('_','',$data['Maintainers']).']';
        CCLog::addOpLogArr($log);
    }

    /**
     * @查询公司的所有业务
     * @param userName,Company
     * @return  业务ID数组
     */
    public function getAppByCompany($source = '') {

        $this->db->from('ApplicationBase');
        $source && $this->db->where('Source', $source);
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

    /*
     * 根据业务名查询业务
     * @param appName
     * @return  业务数组
     */
    public function getAppByName($appName) {
        $this->db->where('ApplicationName', $appName);
        $query = $this->db->get('ApplicationBase');
        return $query && $query->num_rows() > 0 ? $query->result_array() : array();
    }

    /**
     * @查询业务、集群、模块下主机统计
     * @param appId,appName,company
     * @return 统计数组
     */
    public function getAppSetModuleHostStat($appId = array(), $appName, $company) {
        $appId = is_array($appId) ? $appId : explode(',', $appId);
        if($appId!==array('') && !empty($appId)){
            $this->db->where_in('ApplicationID', $appId);
        }

        if(isset($appName) && !empty($appName)) {
            $this->db->where('ApplicationName', $appName);
        }

        if(isset($company) && !empty($company)) {
            $this->db->where('Owner', $company);
        }

        $this->db->order_by('Owner', 'desc');
        $query = $this->db->get('cc_ApplicationBase');
        $appResult = !$query ? array() : $query->result_array();

        if($appId!==array('') && !empty($appId)){
            $this->db->where_in('ApplicationID', $appId);
        }
        $query = $this->db->get('cc_SetBase');
        $setResult = !$query ? array() : $query->result_array();
        $sets = array();
        foreach($setResult as $_sv) {
            if(!isset($sets[$_sv['ApplicationID']])){
                $sets[$_sv['ApplicationID']] = 1;
            }else{
                $sets[$_sv['ApplicationID']]++;
            }
        }

        if($appId!==array('') && !empty($appId)){
            $this->db->where_in('ApplicationID', $appId);
        }
        $query = $this->db->get('cc_ModuleBase');
        $moduleResult = !$query ? array() : $query->result_array();
        $modules = array();
        foreach($moduleResult as $_mv){
            if(!isset($modules[$_mv['ApplicationID']])){
                $modules[$_mv['ApplicationID']] = 1;
            }else{
                $modules[$_mv['ApplicationID']]++;
            }
        }

        if($appId!==array('') && !empty($appId)){
            $this->db->where_in('ApplicationID', $appId);
        }

        $query = $this->db->get('cc_ModuleHostConfig');
        $moduleHostResult = !$query ? array() : $query->result_array();
        $moduleHosts = array();
        foreach($moduleHostResult as $_mhv){
            if(!isset($moduleHosts[$_mhv['ApplicationID']])){
                $moduleHosts[$_mhv['ApplicationID']] = 1;
            }else{
                $moduleHosts[$_mhv['ApplicationID']]++;
            }
        }

        foreach($appResult as $_ak=>$_av){
            $appResult[$_ak]['HostNum'] = isset($moduleHosts[$_av['ApplicationID']]) ? $moduleHosts[$_av['ApplicationID']] : 0;
            $appResult[$_ak]['ModuleNum'] = isset($modules[$_av['ApplicationID']]) ? $modules[$_av['ApplicationID']] : 0;
            $appResult[$_ak]['SetNum'] = isset($sets[$_av['ApplicationID']]) ? $sets[$_av['ApplicationID']] : 0;
        }
        return $appResult;
    }

    /**
     * @查询业务列表
     */
    public function getAppList($fields = '*') {
        $this->db->select($fields);
        $query = $this->db->get('ApplicationBase');
        return $query->num_rows() ? $query->result_array() : array();
    }

}