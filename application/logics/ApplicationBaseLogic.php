<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ApplicationBaseLogic extends Cc_Logic {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 获取用户有权限的App
     * @return  app数组
     */
    public function getUserApp() {
        try {
            $appId = $this->session->userdata('appId');
            if (!$appId) {
                return array();
            }
            $this->load->model('ApplicationBaseModel');
            $app = $this->ApplicationBaseModel->getAppById($appId);
            return $app;
        } catch (Exception $e) {
            CCLog::LogErr($e->getMessage());
            $this->_errInfo = '获取业务异常';
            return array();
        }
    }

    /**
     * 获取业务下的主机个数
     * @param $appIdArr 业务ID的数组
     * @return  业务下的主机个数
     */
    public function getHostNumByAppIdArr($appIdArr) {
        try {
            $this->load->model('HostBaseModel');
            $appHostNum = $this->HostBaseModel->StatHostByApp($appIdArr);
            $result = array();
            if ($appHostNum) {
                foreach ($appHostNum as $ahn) {
                    $result[$ahn['ApplicationID']] = $ahn['cnt'];
                }
            }

            foreach ($appIdArr as $appId) {
                if (!isset($result[$appId])) {
                    $result[$appId] = 0;
                }
            }
            return $result;
        } catch (Exception $e) {
            CCLog::LogErr("获取业务下主机 exception:" . $e->getMessage());
        }
    }

    /*
     * 获取用户有权限的业务
     * @param userName  用户名
     * @param Company   公司名称
     * @return  业务ID数组
     */
    public function getAppIdOnLogin($userName) {
        $role = $this->session->userdata('role');
        try {
            $this->load->model('ApplicationBaseModel');
            $appId = array();
            if ('admin' == $role) {
                $apps = $this->ApplicationBaseModel->getAppByCompany();
            } else {
                $apps = $this->ApplicationBaseModel->getAppIdOnLogin($userName);
            }

            $appId = array_column($apps, 'ApplicationID');
            return array_unique($appId);
        } catch (Exception $e) {
            CCLog::LogErr('getAppIdOnLogin exception:' . $e->getMessage());
            $this->_errInfo = '获取有权限业务ID异常';
            return array();
        }
    }

    /**
     * 增加默认业务
     * @param $company 公司名
     * @return  业务ID数组
     */
    public function addDefaultApp($company) {
        try {
            $this->load->model('ApplicationBaseModel');
            $appId = $this->ApplicationBaseModel->addDefaultApp($company);

            if (!$appId) {
                $this->_errInfo = $this->ApplicationBaseModel->_errInfo;
                return false;
            }

            $nowTime = date('Y-m-d H:i:s');
            $data = array();
            $data['ApplicationID'] = $appId;
            $data['ChnName'] = '';
            $data['Default'] = 1;
            $data['ParentID'] = 0;
            $data['SetName'] = DEFAULT_SET_NAME;
            $data['CreateTime'] = $nowTime;
            $data['LastTime'] = $nowTime;
            $this->load->model('SetBaseModel');
            $setId = $this->SetBaseModel->addSet($data);

            if (!$setId) {
                $this->_errInfo = $this->SetBaseModel->_errInfo;
                return false;
            }

            $moduleNames = explode(',', DEFAULT_MODULE_NAME);
            foreach ($moduleNames as $mn) {
                $data = array();
                $data['ApplicationID'] = $appId;
                $data['Default'] = 1;
                $data['ModuleName'] = $mn;
                $data['SetID'] = $setId;
                $data['CreateTime'] = $nowTime;
                $data['LastTime'] = $nowTime;
                $this->load->model('ModuleBaseModel');
                $moduleId = $this->ModuleBaseModel->addModule($data);

                if (!$moduleId) {
                    $this->_errInfo = $this->ModuleBaseModel->_errInfo;
                    return false;
                }
            }

            return array('AppID' => $appId, 'SetID' => $setId, 'ModuleID' => $moduleId);
        } catch (Exception $e) {
            CCLog::LogErr('add Default App exception:' . $e->getMessage());
            $this->_errInfo = '增加默认业务异常';
            return false;
        }
    }

    /**
     * 删除业务
     * @param $appId  业务ID
     */
    public function deleteApp($appId) {
        try {
            $this->load->model('ModuleHostConfigModel');
            $this->load->model('ApplicationBaseModel');

            $result = array();
            $log = array();
            $hostIdArr = $this->ModuleHostConfigModel->getHostByAppID($appId);
            $app = current($this->ApplicationBaseModel->getAppById($appId));
            if ($app['Default'] == 1) {
                $result['success'] = False;
                $result['errCode'] = 'default_application_failure';
                return $result;
            }
            if (0 != count($hostIdArr)) {
                $result['success'] = False;
                $result['errCode'] = 'application_host_exsit';
                return $result;
            }
            $this->ApplicationBaseModel->deleteApp($appId);

            $this->updateAppInSession();           //刷新session中数据
            $log['ApplicationID'] = $appId;
            $log['OpType'] = '删除';
            $log['OpTarget'] = '业务';
            $log['OpContent'] = '业务名:[' . $app['ApplicationName'] . ']';
            CCLog::addOpLogArr($log);
            $result['success'] = true;
            return $result;
        } catch (Exception $e) {
            CCLog::LogErr("删除业务exception:" . $e->getMessage());
            return $this->getOutput(false, 'application_delete_error');
        }
    }

    /**
     * 新增业务post请求处理
     * @param $appName 业务名
     * @return json
     */
    public function addApplication($appName, $type, $maintainers, $productStr, $level, $lifeCycle, $createTime = '', $creator = '') {
        $this->load->model('ApplicationBaseModel');
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');

        try {
            $nowTime = date('Y-m-d H:i:s');
            $company = $this->session->userdata('company');
            $uin = $this->session->userdata('username');
            $companyName = $this->session->userdata('company');
            $companyId = $this->session->userdata('company_id');
            $appResult = $this->ApplicationBaseModel->getAppByAppNameAndCompany($appName);
            $result = array();
            if ($appResult) {
                return $this->getOutput(false, 'samename_application_exsit');
            }
            $data = array('ApplicationName' => $appName,
                          'Creator' => empty($creator) ? $uin : $creator,
                          'CreateTime' => empty($createTime) ? $nowTime : $createTime,
                          'Default' => 0,
                          'DeptName' => empty($companyName) ?  COMPANY_NAME: $companyName ,
                          'Display' => 1,
                          'Type' => $type,
                          'Level' => $level,
                          'LastTime' => $nowTime,
                          'Owner' => empty($company) ?  COMPANY_NAME: $company ,
                          'Maintainers' => $maintainers,
                          'LifeCycle' => $lifeCycle,
                          'ProductPm' => $productStr,
                          'CompanyID' => 0);

            $appId = $this->ApplicationBaseModel->addApplication($data);
            $set = array('SetName' => DEFAULT_SET_NAME, 'ApplicationID' => $appId, 'Default' => 1);
            $setId = $this->SetBaseModel->addSet($set);
            $defaultModuleArr = explode(',', DEFAULT_MODULE_NAME);

            foreach ($defaultModuleArr as $moduleName) {
                $module = array('SetID' => $setId,
                                'ApplicationID' => $appId,
                                'ModuleName' => $moduleName,
                                'LastTime' => $nowTime,
                                'Default' => 1);
                $this->ModuleBaseModel->addModule($module);
            }

            if(!is_cli()) {
                $this->updateAppInSession();    //刷新session中数据
            }
            $result['success'] = true;
            $result['appId'] = $appId;
            return $result;
        } catch (Exception $e) {
            CCLog::LogErr("新增业务exception:" . $e->getMessage());
            return $this->getOutput(false, 'new_application_error');
        }
    }

    /**
     * 编辑业务信息
     * @param appId,$maintainers,appName
     * @return json
     */
    public function editApplication($appId, $maintainers, $appName) {
        $this->load->model('ApplicationBaseModel');

        try {
            $app = current($this->ApplicationBaseModel->getAppById($appId));
            $company = $this->session->userdata('company');
            $data = array('ApplicationName' => $appName,
                          'Maintainers' => $maintainers,
                          'ApplicationID' => $appId,
                          'LastTime' => date('Y-m-d h:i:s'));

            $checkResult = $this->ApplicationBaseModel->getAppByAppNameAndAppID($appName, $appId);
            if (!$checkResult) {
                return $this->getOutput(false, 'samename_application_exsit');
            }

            $this->ApplicationBaseModel->editApplication($data, $app['ApplicationName'], $app['Maintainers']);

            $this->updateAppInSession();    //刷新session中数据
            $result = array();
            $result['success'] = true;
            return $result;
        } catch (Exception $e) {
            CCLog::LogErr("修改项目exception:" . $e->getMessage());
            return $this->getOutput(false, 'edit_application_fail');
        }

    }

    /**
     * 根据Id获取业务
     * @param AppId
     * @return 业务信息
     */
    public function getAppById($appId) {
        try {
            if (!$appId) {
                return array();
            }
            $this->load->model('ApplicationBaseModel');
            return $this->ApplicationBaseModel->getAppById($appId);
        } catch (Exception $e) {
            CCLog::LogErr("getAppById exception:" . $e->getMessage());
            $this->_errInfo = '根据Id获取业务异常';
            return array();
        }
    }

    /**
     * 根据Id获取业务
     * @param company 公司名
     * @return 业务信息
     */
    public function getAppByCompany() {
        try {
            $this->load->model('ApplicationBaseModel');
            return $this->ApplicationBaseModel->getAppByCompany();
        } catch (Exception $e) {
            CCLog::LogErr("getAppByCompany exception" . $e->getMessage());
            $this->_errInfo = '获取公司业务异常';
            return array();
        }
    }

    /**
     * 根据Id获取业务
     * @param company
     * @return 业务信息
     */
    public function getResPoolByCompany($company) {
        $this->load->model('ApplicationBaseModel');
        try {
            if (!$company) {
                return array();
            }
            return $this->ApplicationBaseModel->getResPoolByCompany($company);
        } catch (Exception $e) {
            CCLog::LogErr("getResPoolByCompany exception" . $e->getMessage());
            $this->_errInfo = '获取资源池错误';
            return array();
        }
    }

    /**
     * 更新session中的业务字段
     */
    public function updateAppInSession() {
        try {
        $this->load->library('session');

        $userName = $this->session->userdata('username');

        $data = array();
        $appIdArr = $this->getAppIdOnLogin($userName);
        $isNewUser = count($appIdArr) == 1 ? 1 : 0;
        $data['appId'] = $appIdArr;
        $this->session->set_userdata($data);
        $appInfo = $this->getUserApp();
        $appHostCount = $this->getHostNumByAppIdArr($appIdArr);
        foreach ($appInfo as $value) {
            if ($value['Default'] != 1) {
                $data['app'][$value['ApplicationID']]['ApplicationID'] = $value['ApplicationID'];
                $data['app'][$value['ApplicationID']]['ApplicationName'] = $value['ApplicationName'];
                $data['app'][$value['ApplicationID']]['Owner'] = $value['Owner'];
                $data['app'][$value['ApplicationID']]['CompanyID'] = $value['CompanyID'];
                $data['app'][$value['ApplicationID']]['ApplicationHostCount'] = $appHostCount[$value['ApplicationID']];
            }
        }

        $this->session->set_userdata($data);

        if (!empty($data['app'])) {
            $data['defaultApp'] = end($data['app']);
            $this->session->set_userdata($data);
            $expires = time() + $this->config->item('sess_expiration') * 1000;
            $this->input->set_cookie('defaultAppId', $data['defaultApp']['ApplicationID'], $expires);
            $this->input->set_cookie('defaultAppName', $data['defaultApp']['ApplicationName'], $expires);
        } else {
            $data['defaultApp'] = array();
            $data['app'] = array();
            $this->session->set_userdata($data);
            setcookie('defaultAppId');
            setcookie('defaultAppName');
        }
        }catch (Exception $e) {
            CCLog::LogErr("getResPoolByCompany exception" . $e->getMessage());
            $this->_errInfo = '获取资源池错误';
            return array();
        }
    }

    /**
     * 根据业务集群模块统计主机
     * @param $appName 业务名
     * @param $company 公司Id
     * @return $result
     */
    public function getAppSetModuleHostStat($appName, $company) {
        $appId = intval($this->input->get_post('ApplicationID'));
        $this->load->model('ApplicationBaseModel');
        try {
            $result = $this->ApplicationBaseModel->getAppSetModuleHostStat($appId, $appName, $company);
            if (!$result) {
                $this->_errInfo = $this->ApplicationBaseModel->_errInfo;
            }
            return $result;
        } catch(Exception $e) {
            CCLog::LogErr("getAppSetModuleHostStat exception" . $e->getMessage());
            return array();
        }
    }

    /*
     * 获取所有业务列表
     */
    public function getAppList() {
        $this->load->model('ApplicationBaseModel');
        try {
            $result =  $this->ApplicationBaseModel->getAppList();

            Utility::filterArrayFields($result, array('Maintainers','ProductPm'));
            return $result;
        } catch(Exception $e) {
            CCLog::LogErr("getAppList exception" . $e->getMessage());
            return array();
        }
    }

    /*
     * 获取业务的分布拓扑树
     */
    public function getAppSetModuleTreeByAppId($appId) {
        $this->load->model('ApplicationBaseModel');
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');
        try {

            $appInfoArr = $this->ApplicationBaseModel->getAppById($appId, '*');
            if(!$appInfoArr) {
                return array();
            }
            Utility::filterArrayFields($appInfoArr, array('Maintainers','ProductPm'));
            $appInfo = $appInfoArr[0];
            $setInfo = $this->SetBaseModel->getSetById(array(), $appId);
            if(!$setInfo) {
                $appInfo['Children'] = array();
                return $appInfo;
            }

            $setIdArr = array_column($setInfo, 'SetID');
            $setId2Set = array();
            foreach($setInfo as $_set){
                $_set['Children'] = array();
                $setId2Set[$_set['SetID']]  = $_set;
            }

            $moduleInfo = $this->ModuleBaseModel->getModuleById(array(), array(), $appId);
            if(!$moduleInfo){
                $appInfo['Children'] = $setInfo;
                return $appInfo;
            }

            $moduleId2module = array();
            foreach($moduleInfo as $_module){
                $moduleId2module[$_module['ModuleID']] = $_module;
            }

            $moduleIdArr = array_keys($moduleId2module);
            $moduleHostNum = $this->getHostNumInModule($moduleIdArr);
            foreach($moduleId2module as $_mid=>&$_module){
                $_module['HostNum'] = isset($moduleHostNum[$_mid]) ? $moduleHostNum[$_mid] : 0;
            }
            foreach($moduleId2module as $_mid=>$_m){
                $_sid = $_m['SetID'];
                if(isset($setId2Set[$_sid]) && isset($moduleId2module[$_mid]))
                {
                    $setId2Set[$_sid]['Children'][] = $moduleId2module[$_mid];
                }
            }
            $appInfo['Children'] = array_values($setId2Set);
            return $appInfo;
        } catch(Exception $e) {
            CCLog::LogErr("getAppSetModuleTreeByAppId exception" . $e->getMessage());
            return array();
        }
    }

    /**
     *获取模块下的主机个数
     */
    private function getHostNumInModule($moduleIdArr) {
        $this->load->model('ModuleHostConfigModel');

        try {
            if(empty($moduleIdArr)) {
                return array();
            }

            $moduleIdKv = array();
            foreach($moduleIdArr as $mo) {
                $moduleIdKv[$mo] = 0;
            }

            $moduleHostIdArr = $this->ModuleHostConfigModel->getHostIdInModuleIdArr($moduleIdArr, 'HostID,ModuleID');
            if(!$moduleHostIdArr) {
                return array();
            }

            foreach($moduleHostIdArr as $hm) {
                $moduleIdKv[$hm['ModuleID']]++;
            }

            return $moduleIdKv;
        } catch(Exception $e) {
            CCLog::LogErr("getHostNumInModule exception" . $e->getMessage());
            return array();
        }
    }

    /**
     *获取用户有权限的app
     */
    public function getAppByUin($userName) {
        $this->load->model('ApplicationBaseModel');
        try {
            $appCompany =  $this->ApplicationBaseModel->getAppByCompany();
            Utility::filterArrayFields($appCompany , array('Maintainers','ProductPm'));

            $appUin =  $this->ApplicationBaseModel->getAppIdOnLogin($userName, '*');
            Utility::filterArrayFields($appUin , array('Maintainers','ProductPm'));

            return array('appCompany'=>$appCompany, 'appUin'=>$appUin);
        } catch(Exception $e) {
            CCLog::LogErr("getAppByUin exception" . $e->getMessage());
            return array();
        }
    }

}