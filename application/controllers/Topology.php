<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Topology extends Cc_Controller {

    private static $_indexTopo = '拓扑管理';

    private static $_topoList = '集群模块列表';

    public function __construct() {
        parent::__construct();
    }

    /**
     * topo首页
     */
    public function index() {
        $this->load->logic('TopologyLogic');
        $this->load->logic('ApplicationBaseLogic');
        $this->load->logic('SetBaseLogic');
        $this->load->Logic('ModuleBaseLogic');

        $data = array();
        $app = $this->session->userdata('defaultApp');
        $appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;

        $appIdArr = $this->session->userdata('appId');     //用户业务Id数组
        $isFirstApp = 0;
        $defaultSetId = 0;

        $appInfo = current($this->ApplicationBaseLogic->getAppById($appId));

        $level = $appInfo['Level'];
        $topo = $this->TopologyLogic->getTree4TopoIndex($appId, $level, $defaultSetId);
        $set = $this->SetBaseLogic->getSetById(array(), $appId);

        if ($level == 3) {
            $emptys = !empty($topo) ? 1 : 0;
        } else {
            $emptys = !empty($topo[0]['items']) ? 1 : 0;
        }
        $data['appId'] = $appId;
        $data['appName'] = $appInfo['ApplicationName'];
        $data['topo'] = json_encode($topo);
        $data['Level'] = $level;
        $data['deSetID'] = $defaultSetId;
        $data['Default'] = $appInfo['Default'];
        $data['emptys'] = $emptys;
        $data['firstapp'] = $isFirstApp;
        $data['header'] = self::$_indexTopo;
        $data['active'] = '/topology';
        $data['subactive'] = '/index';
        $this->load->library('layout');
        $this->layout->view('topo/index', $data);
    }

    /**
     * 刷新topo数据
     * @return json
     */
    public function getTopData() {
        $this->load->logic('TopologyLogic');
        $this->load->logic('ApplicationBaseLogic');
        $this->load->logic('SetBaseLogic');

        $data = array();
        $appId = intval($this->input->get_post('ApplicationID', true));
        if (!$appId) {
            $app = $this->session->userdata('defaultApp');
            $appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
        }

        $appInfo = current($this->ApplicationBaseLogic->getAppById($appId));
        $defaultSetId = 0;
        $data['appId'] = $appId;

        $level = $appInfo['Level'];
        $data['appName'] = $appInfo['ApplicationName'];
        $topo = $this->TopologyLogic->getTree4TopoIndex($appId, $level, $defaultSetId);
        $emptys = count($topo) ? 1 : 0;

        $data['topo'] = json_encode($topo);
        $data['Level'] = $level;
        $data['deSetID'] = $defaultSetId;
        $data['Default'] = $appInfo['Default'];
        $data['emptys'] = $emptys;
        $set = $this->SetBaseLogic->getSetById(array(), $appId);
        $data['header'] = self::$_indexTopo;
        $this->output->set_output(json_encode($data));
    }

    /*
     * 列表页面
     */
    public function topoList() {
        $this->load->Logic('UserLogic');
        $this->load->Logic('ApplicationBaseLogic');

        $result = $this->UserLogic->getUserList();
        $kvResult = array();
        $listArr = array();
        foreach ($result as $uin => $name) {
            $kvResult[$uin] = $uin . "($name)";
            $listArr[] = array('id' => $uin, 'text' => $uin . "($name)");
        }

        $app = $this->session->userdata('defaultApp');
        $appId = isset($app['ApplicationID']) ? $app['ApplicationID'] : 0;
        $appInfo = current($this->ApplicationBaseLogic->getAppById($appId));
        $data = array();
        $data['Level'] = $appInfo['Level'];
        $data['kv'] = json_encode($kvResult);
        $data['list'] = json_encode($listArr);
        $data['header'] = self::$_topoList;
        $data['active'] = '/topology';
        $data['subactive'] = '/topolist';
        $this->load->library('layout');
        $this->layout->view('topo/list', $data);
    }

    /**
     * 获取集群列表
     * @return 集群列表json
     */
    public function setList() {
        $appId = $this->input->post('ApplicationID');
        $this->load->Logic('SetBaseLogic');
        $setArr = $this->SetBaseLogic->listset($appId);

        $result = array();
        $result['success'] = true;
        $result['data'] = $setArr;
        $result['total'] = count($setArr);
        return $this->output->set_output(json_encode($result));
    }

    /*
     * 查询模块列表
     * @return 模块列表json
     */
    public function moduleList() {
        $appId = intval($this->input->post('ApplicationID'));
        $this->load->logic('ModuleBaseLogic');

        $modules = $this->ModuleBaseLogic->listModule($appId);
        $result = array();
        $result['success'] = true;
        $result['data'] = $modules;
        $result['total'] = count($modules);
        $this->output->set_output(json_encode($result));
    }
}