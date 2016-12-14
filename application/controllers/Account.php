<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Account extends Cc_Controller {

    public static $indexPageCh = '用户管理';
    protected $p_adminActions = array();

    public function __construct() {
        parent::__construct();
        $this->load->library('login');
        $this->load->model('UserModel');

        # 校验管理员权限
        $this->setAdminActions();
        $segments = $this->uri->rsegments;
        $action = $segments[2];
        if (in_array($action, $this->p_adminActions)) {
            $this->checkAdminPermission();
        }
    }

    protected function checkAdminPermission() {
        $curUser = $this->login->getCurrentUser();
        $role = $curUser['role'];
        if ($role == 'admin') {
            return true;
        }
        # 普通用户可以操作自己
        $username = $this->input->get_post('UserName');
        if ($username === $curUser['username']) {
            return true;
        }
        $response = array();
        $response['success'] = false;
        $response['message'] = '您无权进行此操作，请联系管理员';
        echo json_encode($response);
        exit();
    }

    /**
     * 把不需要登录校验的action配置在p_noAuthActions中。
     */
    protected function setNoAuthActions() {
        $this->p_noAuthActions = array('login', 'doLogin', 'updateUserCompany', 'authed');
    }

    protected function setAdminActions() {
        $this->p_adminActions = array('changePassword', 'saveUser', 'delUser');
    }

    /**
     * 用户列表
     */
    public function index() {
        $this->load->library('layout');
        $data = $this->buildPageDataArr(self::$indexPageCh, '/account/index');
        $curUser = $this->login->getCurrentUser();
        $role = $curUser['role'];
        if ($role == 'admin') {
            $data['users'] = $this->UserModel->getUserList();
        } else {
            $data['users'] = array($this->UserModel->getUserByUsername($curUser['username']));
        }

        $this->layout->view('account/index', $data);

    }

    /**
     * 登录页面
     */
    public function login() {
        $callbackUrl = $this->input->get_post('cburl');
        $data = array('cburl' => $callbackUrl);
        $this->load->view('account/login', $data);

    }

    /**
     * 登录动作
     */
    public function doLogin() {
        if (!$this->login->loginUser()) {
            $data = array('message' => $this->login->getErrMsg());
            return $this->load->view('account/login', $data);
        }
        $callbackUrl = $this->input->get_post('cburl');
        if (empty($callbackUrl) || $callbackUrl == '') {
            $callbackUrl = BASE_URL . '/';
        }
        $this->load->helper('url');
        redirect($callbackUrl, 'location', 301);
    }

    /**
     * 修改密码
     */
    public function changePassword() {
        $params = ['Password' => array('type' => 'len', 'min' => MIN_PASSWORD_LENGTH),
                   'UserName' => array('type' => 'len', 'min' => MIN_PASSWORD_LENGTH)];
        $response = array('success' => true);
        foreach ($params as $paramKey => $validateConfig) {
            $params[$paramKey] = $this->input->get_post($paramKey);
            if (!Utility::validateInput($params[$paramKey], $validateConfig)) {
                $response['success'] = false;
                $response['message'] = '用户名和密码长度至少为' . MIN_PASSWORD_LENGTH;
                $this->output->set_output(json_encode($response));
                return;
            }
        }
        $result = $this->login->changePassword($params['UserName'], $params['Password']);
        if (!$result) {
            $response['success'] = false;
            $response['message'] = $this->login->getErrMsg();
        }
        $this->output->set_output(json_encode($response));
    }

    /**
     * 保存用户
     */
    public function saveUser() {
        $params = array('id' => array(), 'UserName' => array('type' => 'len', 'min' => MIN_PASSWORD_LENGTH),
                        'ChName' => array('type' => 'len', 'min' => 2), 'Tel' => array(), 'QQ' => array(),
                        'Email' => array(), 'Role' => array('type' => 'value', 'values' => array('admin', 'user')));
        $response = array('success' => true);
        foreach ($params as $paramKey => $validateConfig) {
            $params[$paramKey] = $this->input->get_post($paramKey);
            if (!Utility::validateInput($params[$paramKey], $validateConfig)) {
                $response['success'] = false;
                $response['message'] = '请正确填写用户信息';
                $this->output->set_output(json_encode($response));
                return;
            }
        }
        $curUser = $this->login->getCurrentUser();
        $role = $curUser['role'];
        if ($role != 'admin') {
            $params['Role'] = 'user';
        }
        $result = $this->UserModel->saveUser($params);
        if (empty($result)) {
            $response['success'] = false;
            $response['message'] = $this->UserModel->getErrMsg();
            $this->output->set_output(json_encode($response));
            return;
        }
        # 如果是新增用户， 那么设置初始密码
        if (empty($params['id'])) {
            $this->login->changePassword($result['UserName'], DEFAULT_PASSWORD);
        }
        $response['user'] = $result;
        $this->output->set_output(json_encode($response));
    }

    /**
     * 删除用户
     */
    public function delUser() {
        $delUserId = $this->input->get_post('id');
        $delUserId = (empty($delUserId)) ? 0 : $delUserId;
        $response = array('success' => true);
        $result = $this->UserModel->delUserById($delUserId);
        if (!$result) {
            $response['success'] = false;
            $response['message'] = $this->UserModel->getErrMsg();
        }
        $this->output->set_output(json_encode($response));
    }

    /**
     * 登出动作
     */
    public function logout() {
        $this->login->logoutUser();
        $this->login->redirectUnauthed();
    }

    /**
     * 更新用户表的公司信息
     */
    public function updateUserCompany() {
        $response = array('success' => true);
        $result = $this->UserModel->updateUserCompany();
        if (!$result) {
            $response['success'] = false;
            $response['message'] = $this->UserModel->getErrMsg();
        }
        $this->output->set_output(json_encode($response));
    }
}