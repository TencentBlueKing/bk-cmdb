<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class ModuleHostConfigLogic extends Cc_Logic  {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 根据IP和业务Id查询主机
     */
    public function getHostsByIPAndAppId($IP, $appId) {

        $this->load->model('ModuleHostConfigModel');
        try {
            $hostArr = $this->ModuleHostConfigModel->getHostsByIpAndAppId($IP, $appId);

            if(!$hostArr){
                return $hostArr;
            }

            $ModuleIdArr = array();
            foreach($hostArr as $key=>$value){
                if(!in_array($value['ModuleID'], $ModuleIdArr)){
                    $ModuleIdArr[] = $value['ModuleID'];
                }
            }
            $modules = $this->buildModuleInfo($moduleId, $setId, $appId);
            $sets = $this->buildSetInfo($setId, $appId);
            $app = $this->buildAppInfo($appId);

            return $this->buildHostInfo($hostArr, $modules, $sets, $app);
        } catch (Exception $e) {
            CCLog::LogErr("getHostsByIPAndAppId exception" . $e->getMessage());
            return array();
        }
    }

    /**
     * 根据业务Id和集群Id查询主机
     */
    public function getHostsBySetIdAndAppId($setIdArr, $appId) {

        $this->load->model('ModuleHostConfigModel');
        try {
            $hosts = $this->ModuleHostConfigModel->getHostsBySetIdAndAppId($setIdArr, $appId);
            if(!$hosts) {
                return $hosts;
            }

            $moduleIdArr = array();
            $modules = $this->buildModuleInfo($moduleIdArr, $setId, $appId);
            $sets = $this->buildSetInfo($setId, $appId);
            $app = $this->buildAppInfo($appId);

            return $this->buildHostInfo($hosts, $modules, $sets, $app);
        } catch (Exception $e) {
            CCLog::LogErr("getHostsBySetIdAndAppId exception" . $e->getMessage());
            return array();
        }
    }

    /**
     * 根据业务Id和模块Id查询主机
     */
    public function getHostsByModuleIdAndAppId($moduleIdArr, $appId) {

        $this->load->model('ModuleHostConfigModel');
        try {
            $hosts = $this->ModuleHostConfigModel->getHostsByModuleIdAndAppId($moduleIdArr, $appId);

            if(!$hosts) {
                return $hosts;
            }

            $setIdArr = array();
            $modules = $this->buildModuleInfo($moduleIdArr, $setIdArr, $appId);
            $sets = $this->buildSetInfo($setIdArr, $appId);
            $app = $this->buildAppInfo($appId);

            return $this->buildHostInfo($hosts, $modules, $sets, $app);
        } catch (Exception $e) {
            CCLog::LogErr("getHostsByModuleIdAndAppId exception" . $e->getMessage());
            return array();
        }
    }

    /**
     * 构建业务信息
     */
    private function buildAppInfo($appId){
        $applicationModel = $this->load->model('ApplicationBaseModel');
        try {
            $app = $this->ApplicationBaseModel->getAppById($appId, 'ApplicationID,ApplicationName');

            foreach($app as $key=>$value){
                $app[$value['ApplicationID']] = $value['ApplicationName'];
                unset($app[$key]);
            }
            return $app;
        } catch (Exception $e) {
            CCLog::LogErr("getHostsByModuleIdAndAppId exception" . $e->getMessage());
            return array();
        }
    }

    /**
     * 构建集群信息
     */
    private function buildSetInfo($setId, $appId){
        $this->load->model('SetBaseModel');
        $sets = $this->SetBaseModel->getSetById($setId, $appId, 'SetID,SetName');
        $result = array();
        foreach($sets as $key=>$value){
            $return[$value['SetID']] = $value['SetName'];
        }

        return $result;
    }

    /**
     * 构建模块信息
     */
    private function buildModuleInfo(&$moduleId, &$setId, $appId){
        $this->load->model('ModuleBaseModel');
        $modules = $this->ModuleBaseModel->getModuleById($moduleId, $setId, $appId, 'ApplicationID,SetID,ModuleID,ModuleName');

        foreach($modules as $key=>$value){
            if(!$setId) {
                $setId[] = $value['SetID'];
            }
            $modules[$value['ModuleID']] = $value['ModuleName'];
            unset($modules[$key]);
        }

        return $modules;
    }

    /**
     * 构建主机信息
     */
    private function buildHostInfo($hosts, $modules, $sets, $app){
        foreach($hosts as $key=>$value){
            $value['ApplicationName'] = isset($app[$value['ApplicationID']]) ? $app[$value['ApplicationID']] : '';
            $value['SetName'] = isset($sets[$value['SetID']]) ? $sets[$value['SetID']] : '';
            $value['ModuleName'] = isset($modules[$value['ModuleID']]) ? $modules[$value['ModuleID']] : '';
            $hosts[$key] = $value;
        }

        return $hosts;
    }

}