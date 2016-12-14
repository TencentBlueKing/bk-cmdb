<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ExportLogic extends Cc_Logic {

    public function __construct() {
        parent::__construct();
    }

 /**
  * @主机信息导出到excel
  */
  public function exportHostToExcel() {
    $conditionArr = array();

    $appId = intval($this->input->get_post('ApplicationID'));
    $conditionArr['ApplicationID'] = $appId ? $appId : $this->input->cookie('defaultAppId');

    $setId = $this->input->get_post('SetID');
    $setId && $conditionArr['SetID'] = $setId;

    $moduleId = $this->input->get_post('ModuleID');
    $moduleId && $conditionArr['ModuleID'] = $moduleId;

    $hostId = $this->input->get_post('HostID');
    $hostId && $conditionArrv['HostID'] = $hostId;

    $innerIp = $this->input->get_post('InnerIP');
    $innerIp && $conditionArr['InnerIP'] = explode(',', $innerIp);

    $outerIp = $this->input->get_post('OuterIP');
    $outerIp && $conditionArr['OuterIP'] = explode(',', $outerIp);

    $assetId = $this->input->get_post('AssetID');
    $assetId && $conditionArr['AssetID'] = explode(',', $assetId);

    $operator = $this->input->get_post('Operator');
    $operator && $conditionArr['Operator'] = explode(',', $operator);

    $ifOuterExact = $this->input->get_post('IfOuterexact');
    $ifOuterExact && $conditionArr['IfOuterexact'] = trim($ifOuterExact);

    $ifInneripExact = $this->input->get_post('IfInnerIPexact');
    $ifInneripExact && $conditionArr['IfInnerIPexact'] = trim($ifInneripExact);

    /*导出自定义表头处理表头*/
    $this->load->logic('HostBaseLogic');
    $hostPropertyArr = $this->HostBaseLogic->getHostPropertyByType('nameKey');
    $owner = $this->session->userdata('company');
    $hostCustomerPropertyArr = $this->HostBaseLogic->getHostPropertyByOwner($owner, 'nameKey');
    $hostPropertyArr = array_merge($hostPropertyArr, $hostCustomerPropertyArr);

    $this->load->logic('UserCustomLogic');
    $userCustom = $this->UserCustomLogic->getUserCustom();
    if(empty($userCustom['DefaultColumn'])) {
      $header = $hostPropertyArr;
    } else {
      $field2HostFieldName = array_flip($hostPropertyArr);
      $DefaultColumn = json_decode($userCustom['DefaultColumn'], true);
      foreach ($DefaultColumn as $value) {
        if(isset($field2HostFieldName[$value])) {
          $header[$field2HostFieldName[$value]] = $value;
        }
      }
    }
    
    $this->load->model('HostBaseModel');
    $hosts = $this->HostBaseModel->getHostByCondition($conditionArr, true);
    $hostIdArr =  array_column($hosts, 'HostID');
    unset($hosts);
    $data = $this->HostBaseModel->getHostById($hostIdArr);
    $parameterData = array();
    $this->load->logic('BaseParameterDataLogic');
    $parameterData = $this->BaseParameterDataLogic->getSupportHostSourceKv();
    foreach ($data as $key=>$hostInfo) {
      if(isset($parameterData[$hostInfo['Source']])) {
        $data[$key]['Source'] = $parameterData[$hostInfo['Source']];
      }
    }
    
    $header = array_flip($header);
    foreach ($data as $key => $value) {
        foreach ($value as $k => $v) {
            if(!isset($header[$k])) {
                unset($data[$key][$k]);
            }
        }
    }
    
    setcookie('comdownload', 1);
    ExcelUtility::exportToExcel($header, $data, 'hostInfo');
    return true;
  }
}