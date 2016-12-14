<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Welcome extends Cc_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 默认页面，包含主机汇总和统计数据
     */
    public function index() {
        $data = $this->buildPageDataArr($this->lang->line('my_pools'), '/welcome/index');

        /*获取主机数据及相关统计数据*/
        $data['defaultApp'] = $this->session->userdata('defaultApp');
        $data['appCount'] = count($this->session->userdata('app'));

        $host = $this->_getHost($data['defaultApp']['ApplicationID']);
        $data['hostCount'] = $host['hostCount'];
        $data['hostIDC'] = $host['hostIDC'];

        if ($data['appCount'] == 1) {
            $this->load->logic('SetBaseLogic');
            $group = $this->SetBaseLogic->getSetById(array(), $data['defaultApp']['ApplicationID']);
            $data['GroupCount'] = count($group) - 1;
        }

        /*获取当前用户当前业务下的操作记录*/
        $this->load->logic('OperationLogLogic');
        $data['operationLog'] = $this->OperationLogLogic->getOperationLog($this->session->userdata('username'), $data['defaultApp']['ApplicationID'], '', '', '', '', '', 0, 10, 'OpTime', 'desc');
        $data['emptyPoolHostCount'] = $this->getEmptyHostCount($data['defaultApp']['ApplicationID']);

        $this->load->library('Layout');
        $this->layout->view('welcome', $data);
    }

    /**
     * 设置session中的默认开发商
     * @return json
     */
    public function setDefaultCom() {
        $companyCode = $this->input->get_post('company_code', true);
        $companyName = $this->input->get_post('company_name', true);
        $companyId = $this->input->get_post('company_id', true);

        $data = array();
        $data['company'] = $companyCode;
        $data['company_id'] = $companyId;
        $this->session->set_userdata($data);
        return $this->outputJson(true);
    }

    /**
     * 设置session中的默认业务
     * @return 是否设置默认业务成功
     */
    public function setDefaultApp() {
        $appId = intval($this->input->get_post('ApplicationID', true));
        $this->load->logic('ApplicationBaseLogic');
        $app = $this->ApplicationBaseLogic->getAppById($appId);
        if (!$app) {
            $this->outputJson(false, 'app_not_exist');
            return false;
        }

        $data = array();
        $data['defaultApp']['ApplicationID'] = $app[0]['ApplicationID'];
        $data['defaultApp']['ApplicationName'] = $app[0]['ApplicationName'];

        $companyList = $this->session->userdata('company_list');
        $companyCode = $app[0]['Owner'];
        $CompanyId = $app[0]['CompanyID'];
        $CompanyName = isset($companyList) ? $companyList[$companyCode]['company_name'] : '';
        $data['company'] = $companyCode;
        $data['company_id'] = $CompanyId;

        $this->session->set_userdata($data);
        $expires = $this->config->item('sess_expiration') * 1000;
        $this->input->set_cookie('defaultAppId', $data['defaultApp']['ApplicationID'], $expires);
        $this->input->set_cookie('defaultAppName', $data['defaultApp']['ApplicationName'], $expires);

        return $this->outputJson(true);
    }

    /**
     * 帮助页面
     */
    public function help() {
        $data = $this->buildPageDataArr($this->lang->line('help'), '/welcome/index');
        $this->load->library('Layout');
        $this->layout->view('help', $data);
    }

    /**
     * 获取当前业务空闲机下主机数量
     * @return 空闲机主机数量的数组
     */
    private function getEmptyHostCount($appId) {
        $this->load->logic('ModuleBaseLogic');
        $module = $this->ModuleBaseLogic->getModuleByName(array($this->lang->line('empty_pools')), array(), $appId);

        if (count($module) == 0) {
            return 0;
        }

        $this->load->logic('HostBaseLogic');
        $host = $this->HostBaseLogic->getHostById(array(), $module[0]['ModuleID'], array(), explode(',', $appId));
        return count($host);
    }

    /**
     * 获取当前业务下主机数量以及过期分布
     * @return 主机数量和过期分布的数组
     */
    private function _getHost($ApplicationID) {
        $data = array();
        $host = array();
        $hostIDC = array();
        $parameterData = array();
        $this->load->logic('HostBaseLogic');
        $host = $this->HostBaseLogic->getHostById(array(), array(), array(), explode(',', $ApplicationID));
        $data['hostCount'] = count($host);
        $this->load->logic('BaseParameterDataLogic');
        $parameterData = $this->BaseParameterDataLogic->getBaseParameterDataByDataType('region');
        foreach ($host as $value) {
            if (!isset($parameterData[$value['Region']])) {
                $parameterData[$value['Region']] = $this->lang->line('guang_zhou');
            }
            $hostIDC[$parameterData[$value['Region']]]['city'] = $parameterData[$value['Region']];
            $hostIDC[$parameterData[$value['Region']]]['IPNums'][] = $value['InnerIP'];
        }

        $data['hostIDC'] = array();
        foreach ($hostIDC as $IDC) {
            $IdcArr = array();
            $IdcArr['city'] = $IDC['city'];
            $IdcArr['IPNums'] = count($IDC['IPNums']);
            $data['hostIDC'][] = $IdcArr;
        }
        $data['hostIDC'] = json_encode($data['hostIDC'], true);

        return $data;
    }

}
