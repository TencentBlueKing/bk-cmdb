<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class  UserLogic extends Cc_Logic {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 获取用户列表
     * @return 协作者与创建者信息数组
     */
    public function getUserList() {
        $data = array();
        $this->load->model('UserModel');
        try {
            $userList = $this->UserModel->getUserList();
            foreach($userList as $userInfo) {
                $key = $userInfo['UserName'];
                $data["$key"] = $userInfo['ChName'];//中间转换加引号防止过长
            }
            
            return $data;
        } catch (Exception $e) {
            CCLog::LogErr("getUserList exception:" . $e->getMessage());
            return FALSE;
        }
    }

    /**
     * 获取用户类型
     * @return 用户类型
     */
    public function getUserInfo($userName) {
        $this->load->model('UserModel');
        try {
            $userInfo = $this->UserModel->getUserByUsername($userName);
            return $userInfo;
        } catch (Exception $e) {
            CCLog::LogErr("getUserInfo exception:" . $e->getMessage());
            return FALSE;
        }
    }

    /**
     * 增加默认的admin用户
     */
    public function addDefaultUser() {
        $this->load->model('UserModel');
        try {
            $userInfo = array(  'UserName' => 'admin',
                                'Password' => '$2y$10$dBv5y.ArJAMkzJQ4f4X7Pe40thyfJujNYIcXhPpPrwn6rPRJzMNDy',
                                'ChName' => '公司管理员',
                                'Tel' => '13111112222',
                                'QQ' => '12345',
                                'Email' => 'admin@sample.com',
                                'Role' => 'admin',
                                'Status'=> 'ok');
            $this->UserModel->saveUser($userInfo);
        } catch (Exception $e) {
            CCLog::LogErr("addDefaultUser exception:" . $e->getMessage());
            return FALSE;
        }
    }

}