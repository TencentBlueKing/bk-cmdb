<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class UserCustomLogic extends Cc_Logic {
    public function __construct() {
        parent::__construct();
    }

    /**
     * 查询用户定制数据
     * @return array 用户定制数据
     */
    public function getUserCustom() {
        try {
            $userName = $this->session->userdata('username');

            if (strlen($userName) === 0) {
                return array();
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->getUserCustom($userName);
            if (!$result) {
                return array();
            }

            return $result;
        } catch (Exception $e) {
            CCLog::LogErr($e->getMessage());
            $this->_errInfo = 'getUserCustom exception';
            return array();
        }
    }

    /**
     * 查询用户定制数据
     * @return array 用户定制列
     */
    public function getUserCustomColumn() {
        try {
            $userName = $this->session->userdata('username');

            if (strlen($userName) === 0) {
                return array();
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->getUserCustom($userName);
            if (!$result || empty($result['DefaultColumn'])) {
                return array('InnerIP', 'OuterIP', 'SetName', 'ModuleName', 'HostName', 'SN', 'ApplicationName');
            }

            return json_decode($result['DefaultColumn']);
        } catch (Exception $e) {
            CCLog::LogErr('getUserCustomColumn exception'.$e->getMessage());
            $this->_errInfo = 'getUserCustomColumn exception';
            return array();
        }
    }

    /**
     * 查询用户定制数据
     * @return int 用户默认appId
     */
    public function getUserCustomApp() {
        try {
            $userName = $this->session->userdata('username');

            if (strlen($userName) === 0) {
                return array();
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->getUserCustom($userName);
            if (!$result) {
                return 0;
            }

            return $result['DefaultApplication'];
        } catch (Exception $e) {
            CCLog::LogErr('getUserCustomApp exception'.$e->getMessage());
            $this->_errInfo = 'getUserCustomApp exception';
            return array();
        }
    }

    /**
     * 查询用户定制数据
     * @return 用户定制每页数据条数
     */
    public function getUserCustomPageSize() {
        try {
            $userName = $this->session->userdata('username');

            if (strlen($userName) === 0) {
                return array();
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->getUserCustom($userName);
            if (!$result) {
                return 20;
            }

            return $result['DefaultPageSize'];
        } catch (Exception $e) {
            CCLog::LogErr('getUserCustomPageSize exception'.$e->getMessage());
            $this->_errInfo = 'getUserCustomPageSize exception';
            return array();
        }
    }

    /**
     * 设置用户定制数据
     * @return boolean 成功 or 失败
     */
    public function setUserCustom() {
        try {
            $userName = $this->session->userdata('username');
            if (strlen($userName) === 0) {
                return false;
            }

            $data = array();
            intval($this->input->get_post('ApplicaitonID')) && $data['DefaultApplication'] = intval($this->input->get_post('ApplicaitonID'));

            $defaultColumn = $this->input->get_post('DefaultColumn');
            if ($defaultColumn && is_array($defaultColumn)) {
                foreach ($defaultColumn as &$_dc) {
                    if ($_dc == 'checkbox') {
                        unset($_dc);
                    }
                }
                $data['DefaultColumn'] = json_encode($defaultColumn);
            }

            intval($this->input->get_post('DefaultPageSize')) && $data['DefaultPageSize'] = intval($this->input->get_post('DefaultPageSize'));
            $this->input->get_post('DefaultCon') && $data['DefaultCon'] = $this->input->get_post('DefaultCon');

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->setUserCustom($data, $userName);

            if (!$result) {
                $this->_errInfo = $this->UserCustomModel->_errInfo;
            }

            return $result;
        } catch (Exception $e) {
            CCLog::LogErr('setUserCustom exception'.$e->getMessage());
            $this->_errInfo = 'setUserCustom exception';
            return false;
        }
    }

    /**
     * 设置用户定制数据
     * @return boolean 成功 or 失败
     */
    public function setUserCustomByUserName($data) {
        try {
            $userName = $this->session->userdata('username');
            if (strlen($userName) === 0) {
                return false;
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->setUserCustom($data, $userName);

            if (!$result) {
                $this->_errInfo = $this->UserCustomModel->_errInfo;
            }

            return $result;
        } catch (Exception $e) {
            CCLog::LogErr('setUserCustomByUserName exception'.$e->getMessage());
            $this->_errInfo = 'setUserCustomByUserName exception';
            return false;
        }
    }

    /**
     * 查询用户定制数据
     * @return array 用户定制查询字段
     */
    public function getUserCustomDefaultField() {
        try {
            $userName = $this->session->userdata('username');

            if (strlen($userName) === 0) {
                return array();
            }

            $this->load->model('UserCustomModel');
            $result = $this->UserCustomModel->getUserCustom($userName);
            if (!$result || empty($result['DefaultField'])) {
                return array();
            }

            return json_decode($result['DefaultField'], true);
        } catch (Exception $e) {
            CCLog::LogErr('getUserCustomDefaultField exception'.$e->getMessage());
            $this->_errInfo = 'getUserCustomDefaultField exception';
            return array();
        }
    }
}