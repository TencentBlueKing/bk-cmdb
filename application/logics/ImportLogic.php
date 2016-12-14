<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ImportLogic extends Cc_Logic {

    private $_custometerFields = array('Customer001', 'Customer002', 'Customer003','Customer004', 'Customer005', 'Customer006','Customer007', 'Customer008', 'Customer009', 'Customer010','Customer011', 'Customer012', 'Customer013','Customer014', 'Customer015', 'Customer016','Customer017', 'Customer018', 'Customer019', 'Customer020','Customer021', 'Customer022', 'Customer023','Customer024', 'Customer025', 'Customer026','Customer027', 'Customer028', 'Customer029', 'Customer030','Customer031', 'Customer032', 'Customer033','Customer034', 'Customer035', 'Customer036','Customer037', 'Customer038', 'Customer039', 'Customer040','Customer041', 'Customer042', 'Customer043','Customer044', 'Customer045', 'Customer046','Customer047', 'Customer048', 'Customer049', 'Customer050');

    public function __construct() {
        parent::__construct();
    }

    /**
    * @导入机器
    * @return boolean 成功 or 失败
    */
    public function importPrivateHostByExcel() {
        set_time_limit(0);
        $company = $this->session->userdata('company');
        $this->load->model('ApplicationBaseModel');
        $app = $this->ApplicationBaseModel->getResPoolByCompany($company);

        $appId = $app['ApplicationID'];
        $appName = $app['ApplicationName'];

        /*处理导入的字段映射关系*/
        $filename = $this->input->post('filename');
        $post = $this->input->post();
        $exportFields = array();
        unset($post['filename']);
        foreach ($post as $key => $value) {
            if(isset($post[$key . '_' . $value])) {
                $exportFields[$post[$key . '_' . $value]] = $key;
            }
        }

        /*如果不存在，则认为是自定义*/
        $this->load->logic('HostBaseLogic');
        $hostProperty = $this->HostBaseLogic->getHostPropertyByType('keyName');
        $hostCustomerProperty = $this->HostBaseLogic->getHostPropertyByOwner($company, 'keyName');
        $hostProperty = array_merge($hostProperty, $hostCustomerProperty);
        foreach ($this->_custometerFields as $key => $value) {
            if(in_array($value, array_keys($hostCustomerProperty))) {
                unset($this->_custometerFields[$key]);
            }
        }
        
        foreach ($exportFields as $key => $value) {
            if($key && !in_array($key, array_keys($hostProperty))) {
                if(!in_array($key, $hostProperty)) {
                    $hostTableField = array_shift($this->_custometerFields);
                    if($hostTableField) {
                        $hostCustomerPropertyData = array(
                            'PropertyKey'=>$key,
                            'PropertyName'=>$key,
                            'HostTableField'=>$hostTableField,
                            'Group'=>'Customer',
                            'Owner'=>$company,
                            'CreateTime'=>date('Y-m-d H:i:s'),
                            'LastTime'=>date('Y-m-d H:i:s')
                        );

                        $this->load->model('HostCustomerPropertyModel');
                        $this->HostCustomerPropertyModel->addHostCustomerProperty($hostCustomerPropertyData);
                        $exportFields[$hostTableField] = $value;
                    }
                }
                unset($exportFields[$key]);
            }
        }
        
        $data = ExcelUtility::readExcel($filename);
        if(!isset($data['内网IP'])) {
            $this->_errInfo = $this->config->item('innerip_column_not_found_failure')->Info;
            return false;
        }

        $success = $errInfo = array();
        $this->load->model('ModuleBaseModel');
        $module = $this->ModuleBaseModel->getModuleByName(DEFAULT_MODULE_NAME, array(), $appId);
        if(!$module) {
            $this->_errInfo = '`appName='. $appName . $this->config->item('lack_of_empty_module')->Info;
            return false;
        }

        $moduleId = $module[0]['ModuleID'];
        $setId    = $module[0]['SetID'];
        if(!$moduleId || !$setId || !$appId) {
            $this->_errInfo = '`appName='. $appName . $this->config->item('lack_of_setid_or_moduleid')->Info;
            return false;
        }

        $importInnerIp = array();
        $error = array();
        for($i=0, $hostNum = count($data['内网IP']); $i<$hostNum; $i++) {
            if(strlen($data['内网IP'][$i]) === 0) {
                continue;
            }

            if(isset($exportFields['HostName'])) {
                $data[$exportFields['HostName']][$i] = trim($data[$exportFields['HostName']][$i]);
            }
            if(isset($exportFields['OuterIP'])) {
                $data[$exportFields['OuterIP']][$i] = trim($data[$exportFields['OuterIP']][$i]);
            }
            
            $data['内网IP'][$i] = trim($data['内网IP'][$i]);
            
            if(in_array($data['内网IP'][$i], $importInnerIp)) {
                $error[$i] = '第'. ($i+1) . '行`内网IP:' . $data['内网IP'][$i] . '`和前面的重复了';
                continue;
            } else {
                if($data['内网IP'][$i]) {
                    $importInnerIp[] = $data['内网IP'][$i];
                }
            }

            if(strlen($data['内网IP'][$i]) === 0) {
                $error[$i] =  '第'. ($i+1) .'行`内网IP`不能为空';
                continue;
            }

            $innerIp = explode(',', $data['内网IP'][$i]);
            foreach($innerIp as $_ip) {
                if(!filter_var($_ip, FILTER_VALIDATE_IP)) {
                    $error[$i] =  '第'. ($i+1) .'行' . $_ip . '非法';
                    continue;
                }
            }

            if(isset($exportFields['OuterIP']) && isset($data[$exportFields['OuterIP']][$i]) && strlen($data[$exportFields['OuterIP']][$i])>0) {
                $outerIp = explode(',', $data[$exportFields['OuterIP']][$i]);
                foreach($outerIp as $_ip) {
                    if(!filter_var($_ip, FILTER_VALIDATE_IP)) {
                        $error[$i] =  '第'. ($i+1) .'行' . $_ip . '非法';
                        continue;
                    }
                }
            }
        }

        $this->load->logic('BaseParameterDataLogic');

        /*是否已经存在一样的内网IP*/
        $condition = array();
        $hosts = array();
        $exHost = array();
        $condition['InnerIP'] = $importInnerIp;
        $this->load->logic('HostBaseLogic');
        if(count($importInnerIp) > 0) {
            $hosts = $this->HostBaseLogic->getHostByCondition($condition);
            if(count($hosts) > 0) {
                foreach($hosts as $host) {
                    $exHost[$host['HostID']] = $host['InnerIP'];
                }
            }
        }
        
        $prefix = 'pc';
        $totalNums = count($data['内网IP']);
        $update = $insert = $successIp = array();
        for($i=0, $hostNum = count($data['内网IP']); $i<$hostNum; $i++) {
            $hostInfo = array();
            if(isset($error[$i])) {
                $errInfo[] = $error[$i];
                continue;
            }

            /*自定义导入字段*/
            foreach ($exportFields as $key => $value) {
                if($key){
                    $hostInfo[$key] = isset($data[$value][$i]) ? htmlspecialchars($data[$value][$i]) : '';
                }
            }

            $hostInfo['InnerIP'] = isset($data['内网IP'][$i]) ? htmlspecialchars($data['内网IP'][$i]) : '';
            $hostInfo['Source'] = 3;
            $hostInfo['CreateTime'] = date('Y-m-d H:i:s');
            $hostInfo['LastTime'] = date('Y-m-d H:i:s');

            if(strlen($hostInfo['InnerIP']) === 0) {
                continue;
            }

            if(in_array($hostInfo['InnerIP'], $exHost)) {
                $exHostArr = array_flip($exHost);
                $this->load->model('HostBaseModel');
                unset($hostInfo['CreateTime']);
                $updateHostRes = $this->HostBaseModel->updateHostById($hostInfo, $exHostArr[$hostInfo['InnerIP']]);
                if(!$updateHostRes) {
                    $errInfo[] = $this->HostBaseModel->_errInfo;
                    continue;
                }

                $update[] = $i;
            }else {
                $this->load->model('HostBaseModel');
                $addHostRes = $this->HostBaseModel->AddHost($hostInfo);
                if(!$addHostRes) {
                    $errInfo[] = $this->HostBaseModel->_errInfo;
                    continue;
                }


                $this->load->model('ModuleHostConfigModel');
                $addRes = $this->ModuleHostConfigModel->addModuleHostConfig($addHostRes, $moduleId, $setId, $appId);

                if(!$addRes) {
                    $errInfo[] = '第'. ($i+1) .'行' . $hostInfo['InnerIP'] . '到空闲机，失败';
                    continue;
                }
                $insert[] = $i;
            }

            $success[] = ($i+1);
            $successIp[] = $hostInfo['InnerIP'];
        }

        $result = count($success) === $totalNums ? TRUE : FALSE;
        
        if(count($success) > 0) {
            $log = array();
            $log['ApplicationID'] = $appId;
            if($result){
                $this->_errInfo = count($insert) > 0 ?  count($insert) . '条导入成功;' : '';
                $this->_errInfo .= count($update) > 0 ?  count($update) . '条更新成功;' : '';
                $log['OpContent'] = implode(' | ', $successIp);
            }else{
                $this->_errInfo = '第' . implode(',', $success) . '行导入成功;';
                $this->_errInfo = count($insert) > 0 ?  count($insert) . '条导入成功;' : '';
                $this->_errInfo .= count($update) > 0 ?  count($update) . '条更新成功;' : '';
                $this->_errInfo .= '，失败详细如下！<ul class="import-error-list"><li>'.implode("</li><li>", $errInfo).'</li></ul>';
                $log['OpContent'] = implode(' | ', $successIp) .' | '. implode(' | ', $errInfo);
            }

            $log['OpName'] = '导入私有云机器';
            $log['OpResult'] = '成功';
            $log['OpTarget'] = '主机';
            $log['OpType'] = '添加';
            CCLog::addOpLogArr($log);
        } else {
            $this->_errInfo = '导入失败！<ul class="import-error-list"><li>'.implode("</li><li>", $errInfo).'</li></ul>';
            
            $log = array();
            $log['ApplicationID'] = $appId;
            $log['OpContent'] = implode(' | ', $successIp) .' | '. implode(' | ', $errInfo);
            $log['OpName'] = '导入私有云机器';
            $log['OpResult'] = '失败';
            $log['OpTarget'] = '主机';
            $log['OpType'] = '添加';
            CCLog::addOpLogArr($log);
        }

        return $result;
    }

    /**
     * @导入机器时读取表头
     */
    public function getImportPrivateHostTableFieldsByExcel() {
        set_time_limit(0);
        $data = array();
        $company = $this->session->userdata('company');
        $this->load->model('ApplicationBaseModel');
        $app = $this->ApplicationBaseModel->getResPoolByCompany($company);

        if(!$app) {
            return $this->getOutput(false, 'resource_app_not_exist');
        }
        
        $appId = $app['ApplicationID'];
        $appName = $app['ApplicationName'];

        $time = date("YmdHis" , time());
        $fileName = $this->session->userdata('username') . $time;
        $config['upload_path'] = $this->config->item('upload_path') . '/importPrivateHostByExcel/';
        $config['allowed_types'] = 'xls|xlsx|csv';
        $config['file_name'] = $fileName;
        $this->load->library('upload', $config);
        
        if (!$this->upload->do_upload("importPrivateHost")) {//"importPrivateHostByExcel"页面的文件域
            return $this->getOutput(false, 'upload_file_error_failure');
        }
        
        $file = $this->upload->data();
        $fileContent = ExcelUtility::readExcel($file['full_path']);
        $keys = array_keys($fileContent);
        $tableField = array();
        foreach ($keys as $key => $value) {
            if(strlen($value) > 16) {
                return $this->getOutput(false, 'table_header_field_beyond');
            }

            if($value && !in_array($value, array('内网IP'))) {
                $temVar = array();
                $temVar['id'] = $key;
                $temVar['name'] = $value;
                $tableField[] = $temVar;
            }
        }

        if(!in_array('内网IP', $keys)) {
            return $this->getOutput(false, 'table_header_order');
        }

        $data['keys'] = $tableField;
        $this->load->logic('HostBaseLogic');
        $hostProperty = $this->HostBaseLogic->getHostPropertyByType('keyName');
        $hostCustomerProperty = $this->HostBaseLogic->getHostPropertyByOwner($company, 'keyName');
        $hostProperty = array_merge($hostProperty, $hostCustomerProperty);
        foreach ($hostProperty as $key => $value) {
            if(!in_array($key, array('InnerIP', 'Source', 'HostID', 'DeadLineTime', 'ModuleName', 'SetName', 'CreateTime'))) {
                $temVar = array();
                $temVar['id'] = $key;
                $temVar['text'] = $value;
                $data['fields'][] = $temVar;
            }
        }
        
        $data['filename'] = $file['full_path'];
        $data['name'] = 'importToCC';
        $data['success'] = true;
        return $data;
    }
}