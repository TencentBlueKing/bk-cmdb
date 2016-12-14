<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class App extends Cc_Controller {

    private static $_indexPageCh = '业务管理';

    private static $_newPageCh = '新增业务';

    public function __construct() {
        parent::__construct();
    }

    /**
     * 业务管理引导页
     */
    public function index() {
        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('UserLogic');
        $this->load->library('layout');

        $data = $this->buildPageDataArr(self::$_indexPageCh, '/app/index');
        $userApp = $this->ApplicationBaseLogic->getUserApp();

        foreach ($userApp as $k => $ua) {
            if ($ua['Default'] == 1) {
                unset($userApp[$k]);
            }
        }

        if (0 == count($userApp)) {
            $this->layout->view('app/noApp', $data);
        } else {
            $appIdArr = array_column($userApp, 'ApplicationID');
            $viewData = array();
            $appHostNum = $this->ApplicationBaseLogic->getHostNumByAppIdArr($appIdArr);
            $userListKv = $this->UserLogic->getUserList();
            foreach ($userApp as &$ua) {
                $ua['Maintainers'] = str_replace('_', '', $ua['Maintainers']);
                $maintainerArr = explode(';', $ua['Maintainers']);
                foreach ($maintainerArr as $key => &$ma) {
                    if (isset($userListKv[$ma])) {
                        $ma = $this->getUserDisplayName($userListKv, $ma);
                    } else {
                        unset($maintainerArr[$key]);
                    }

                }
                $ua['Maintainers'] = implode(';', array_values($maintainerArr));
                $ua['Creator'] = $this->getUserDisplayName($userListKv, $ua['Creator']);
                $ua['CreateTime'] = substr($ua['CreateTime'], 0, 11);
                $ua['HostNum'] = $appHostNum[$ua['ApplicationID']];
                $viewData [] = $ua;
            }
            $data['app'] = $viewData;
            $this->layout->view('app/index', $data);
        }
    }

    /**
     * 新增业务页面
     */
    public function newApp() {
        $this->load->Logic('UserLogic');
        $this->load->library('layout');

        $data = $this->buildPageDataArr(self::$_newPageCh, '/app/index');
        $data['owner'] = $this->session->userdata('company');
        $data['uin'] = $this->session->userdata('username');

        $userList = $this->UserLogic->getUserList();
        $data['userList'] = array_keys($userList);
        $data['userNameList'] = array_values($userList);
        $uinArr = array();
        foreach ($userList as $key => $value) {
            $uinArr[] = array('id' => $key, 'text2' => $value, 'text' => $key);
        }

        $data['userListJ'] = json_encode($uinArr);
        $this->layout->view('app/newApp', $data);
    }

    /**
     * 新增业务请求处理
     */
    public function add() {
        $type = intval($this->input->post('Type', true));
        $level = intval($this->input->post('Level', true));
        $lifeCycle = $this->input->post('LifeCycle', true);
        $lifeCycle = isset($lifeCycle) ? $lifeCycle : '';

        $appName = trim(htmlspecialchars($this->input->post('ApplicationName', true)));
        if (mb_strlen($appName) > 32 || mb_strlen($appName) < 1) {
            $this->outputJson(false, 'application_length_error');
            return;
        }

        $maintainers = $this->input->post('Maintainers', true);
        if (count($maintainers) > 24 || empty($maintainers)) {
            $this->outputJson(false, 'maintainer_length_error');
            return;
        }
        $maintainerStr = implode('_;_', $maintainers);      //转存运维字段
        $maintainerStr = '_' . $maintainerStr . '_';

        $products = $this->input->post('ProducterList', true);
        if (count($products) > 8 || empty($products)) {
            $this->outputJson(false, 'producterlist_length_error');
            return;
        }
        $productStr = implode('_;_', $products);
        $productStr = '_' . $productStr . '_';

        $this->load->Logic('ApplicationBaseLogic');
        $result = $this->ApplicationBaseLogic->addApplication($appName, $type, $maintainerStr, $productStr, $level, $lifeCycle);
        $appIdArr = $this->session->userdata('appId');

        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        if ($result['success'] == true) {
            $data = array('success' => true);
            return $this->output->set_output(json_encode($data));
        } else {
            return $this->outputJson($result['success'], $errCode);
        }

    }

    /**
     * 删除业务请求
     */
    public function delete() {
        $appId = intval($this->input->post('ApplicationID', true));
        $this->load->logic('ApplicationBaseLogic');

        $delResult = $this->ApplicationBaseLogic->deleteApp($appId);
        $errCode = isset($delResult['errCode']) ? $delResult['errCode'] : '';
        return $this->outputJson($delResult['success'], $errCode);
    }

    /**
     * 修改业务请求
     */
    public function edit() {
        $appName = trim(htmlspecialchars($this->input->post('ApplicationName', true)));
        $appID = intval($this->input->post('ApplicationID', true));
        $maintainers = $this->input->post('Maintainers', true);
        $maintainerArr = array_unique(explode(';', $maintainers));

        if (mb_strlen($appName) > 32 || mb_strlen($appName) < 1) {
            $this->outputJson(false, 'application_length_error');
            return;
        }
        if (count($maintainerArr) > 24) {
            $this->outputJson(false, 'maintainer_length_error');
            return;
        }

        $maintainerStr = implode("_;_", $maintainerArr);
        $maintainerStr = '_' . $maintainerStr . '_';
        $this->load->Logic('ApplicationBaseLogic');
        $result = $this->ApplicationBaseLogic->editApplication($appID, $maintainerStr, $appName);
        $errCode = isset($result['errCode']) ? $result['errCode'] : '';
        return $this->outputJson($result['success'], $errCode);

    }

    /**
     * 获取运维列表
     * @return  json
     */
    public function getMaintainers() {
        $this->load->Logic('UserLogic');
        $result = $this->UserLogic->getUserList();
        if (!$result) {
            return $this->outputJson(false, 'get_maintainers_fail');
        }
        return $this->output->set_output(json_encode($result));
    }

}