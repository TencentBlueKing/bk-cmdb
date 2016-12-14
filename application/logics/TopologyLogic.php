<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class TopologyLogic extends Cc_Logic {
    public function __construct() {
        parent::__construct();
    }

    /**
     * id查询拓扑树
     */
    public function getTopoTree4view($appId) {
        $this->load->model('ModuleHostConfigModel');
        $this->load->model('ApplicationBaseModel');
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');
        $this->load->model('HostBaseModel');
        try {

            $topo = array('topo' => array(), 'empty' => array());
            $app = $this->ApplicationBaseModel->getAppById($appId);
            if (!$app) {
                return $topo;
            }
            $topo['topo'][0]['id'] = $app[0]['ApplicationID'];
            $topo['topo'][0]['text'] = $app[0]['ApplicationName'];
            $topo['topo'][0]['type'] = 'application';
            $topo['topo'][0]['icon'] = 'c-icon icon-app application';
            $topo['topo'][0]['opened'] = true;
            $topo['topo'][0]['lvl'] = $app[0]['Level'];

            $set = $this->SetBaseModel->getSetById(array(), $appId);
            $module = $this->ModuleBaseModel->getModuleById(array(), array(), $appId);
            $host = $this->HostBaseModel->getHostAllIdById(array(), array(), array(), $appId);

            $emptyModuleId = 0;
            foreach ($module as $_m) {
                if ($_m['Default'] == 1) {
                    $emptyModuleId = $_m['ModuleID'];
                    break;
                }
            }
            
            $hostInApp = $this->ModuleHostConfigModel->statHostCountByApp($appId, $emptyModuleId);
            $hostInSet = $this->ModuleHostConfigModel->statHostCountBySetID($appId);
            $hostInModule = $this->ModuleHostConfigModel->statHostCountByModuleID($appId);

            $topo['topo'][0]['number'] = isset($hostInApp[$appId]) ? $hostInApp[$appId] : 0;
            $topo['topo'][0]['text'] = $topo['topo'][0]['text'].'<span class="host-topo-num">'.$topo['topo'][0]['number'].'</span>';
            $expand = isset($topo['topo'][0]['lvl']) && $topo['topo'][0]['lvl'] == 3 ? false : true;
            foreach ($set as $s) {
                $setItem = array();
                $setItem['id'] = $s['SetID'];
                $setItem['appId'] = $topo['topo'][0]['id'];
                $setItem['text'] = $s['SetName'];
                $setItem['icon'] = 'c-icon icon-group fa-hide set';
                $setItem['type'] = 'set';
                $setItem['opened'] = $expand;
                $setItem['number'] = isset($hostInSet[$s['SetID']]) ? $hostInSet[$s['SetID']] : 0;
                $number = $setItem['number'];
                $setItem['text'] = $s['SetName'].'<span class="host-topo-num">'.$number.'</span>';
                $setItem['children'] = array();
                foreach ($module as $m) {
                    $moduleItem = array();
                    if ($m['SetID'] === $s['SetID']) {
                        $moduleItem['id'] = $m['ModuleID'];
                        $moduleItem['appId'] = $topo['topo'][0]['id'];
                        $moduleItem['text'] = $m['ModuleName'];
                        $moduleItem['icon'] = 'c-icon icon-modal module';
                        $moduleItem['type'] = 'module';
                        $moduleItem['number'] = isset($hostInModule[$m['ModuleID']]) ? $hostInModule[$m['ModuleID']] : 0;
                        $number = $moduleItem['number'];
                        $moduleItem['text'] = $m['ModuleName'].'<span class="host-topo-num">'.$number.'</span>';
                        if ($m['Default'] == 1) {
                            $topo['empty'][] = $moduleItem;
                        } else {
                            $setItem['children'][] = $moduleItem;
                        }
                    }
                }
                if ($app[0]['Level'] == 3 && $s['Default'] != 1) {        //3级树不显示空闲机池
                    $topo['topo'][0]['children'][] = $setItem;
                } elseif ($app[0]['Level'] == 2 && $s['Default'] == 1) {  //2级树只显示空闲机池
                    $topo['topo'][0]['children'] = $setItem['children'];
                }
            }
            return $topo;
        } catch (Exception $e) {
            CCLog::LogErr('getTopoTree4view exception'.$e->getMessage());
            $this->_errInfo = 'getTopoTree4view exception';
            return array();
        }
    }

    /**
     * 获取业务对应topo树组装为前台可以展示的结构
     * @params appId
     * @return 拓扑结构
     */
    public function getTree4TopoIndex($appId, $level, &$defaultSetId = 0) {
        $this->load->model('SetBaseModel');
        $this->load->model('ModuleBaseModel');
        try {
            $setArr = $this->SetBaseModel->getSetById(array(), $appId);
            $moduleArr = $this->ModuleBaseModel->getModuleById(array(), array(), $appId);

            $setInfoArr = array();

            if ($level == 2) {
                $set = current($setArr);
                $defaultSetId = $set['SetID'];
//                $spriteCssClass = 'c-icon icon-group hide';
//                $setInfo = array();
//                $setInfo['id'] = $set['SetID'];
//                $setInfo['text'] = $set['SetName'];
//                $setInfo['icon'] = $spriteCssClass;
//                $setInfo['spriteCssClass'] = $spriteCssClass;
//                $setInfo['type'] = 'set';
//                $setInfo['opened'] = true;
//                $setInfo['number'] = 32;
//                $setInfo['children'] = array();
                foreach ($moduleArr as $module) {
                    if ($set['SetID'] == $module['SetID'] && $module['Default'] == 0) {
                        $moduleInfo = array();
                        $moduleInfo['id'] = $module['ModuleID'];
                        $moduleInfo['icon'] = 'c-icon icon-modal';
                        $moduleInfo['text'] = $module['ModuleName'];
                        $moduleInfo['text'] = "<span class='node-text'>".$moduleInfo['text']."</span>";
                        $moduleInfo['type'] = 'module';
                        $moduleInfo['number'] = 32;
                        $moduleInfo['operator'] = $module['Operator'];
                        $moduleInfo['bakoperator'] = $module['BakOperator'];
                        $setInfoArr[] = $moduleInfo;
//                        $setInfo['children'][] = $moduleInfo;
                    }
                }
//                $setInfoArr [] = $setInfo;
            } elseif ($level == 3) {
                foreach ($setArr as $set) {
                    $setInfo = array();
                    $setInfo['id'] = $set['SetID'];
                    $setInfo['text'] = $set['SetName'];
                    $setInfo['text'] = "<span class='node-text'>".$setInfo['text']."</span>"."<span class='creat-module-btn btn btn-success btn-xs'><i class='fa fa-plus'></i> 模块</span>";
                    $spriteCssClass = 'c-icon icon-group';
                    $setInfo['icon'] = $spriteCssClass;
                    $setInfo['spriteCssClass'] = $spriteCssClass;
                    $setInfo['type'] = 'set';
                    $setInfo['opened'] = false;
                    $setInfo['number'] = 32;
                    $setInfo['children'] = array();
                    foreach ($moduleArr as $module) {
                        if ($set['SetID'] == $module['SetID'] && $module['Default'] == 0) {
                            $moduleInfo = array();
                            $moduleInfo['id'] = $module['ModuleID'];
                            $moduleInfo['icon'] = 'c-icon icon-modal';
                            $moduleInfo['text'] = $module['ModuleName'];
                            $moduleInfo['text'] = "<span class='node-text'>".$moduleInfo['text']."</span>";
                            $moduleInfo['operator'] = $module['Operator'];
                            $moduleInfo['bakoperator'] = $module['BakOperator'];
                            $moduleInfo['type'] = 'module';
                            $moduleInfo['number'] = 32;
                            $setInfo['children'][] = $moduleInfo;
                        }
                    }
                    if ($set['Default'] == 0) {
                        $setInfoArr [] = $setInfo;
                    }
                }
            }
            return $setInfoArr;
        } catch (Exception $e) {
            CCLog::LogErr('getTree4TopoIndex exception'.$e->getMessage());
            $this->_errInfo = 'getTree4TopoIndex exception';
            return array();
        }
    }
}