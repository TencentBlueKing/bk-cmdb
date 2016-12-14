<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class UserModel extends Cc_Model{
	public function __construct(){
		parent::__construct();
        $this->load->database('db');
	}

	public function getUserByUsername($username=''){
        try{
            $this->db->where('UserName', $username);
		    $query = $this->db->get('cc_User');
		    if(empty($query) || $query->num_rows() < 1){
                return false;
            }
            $userInfo = $query->first_row('array');
        }catch (Exception $ex){
            return false;
        }
        return $userInfo;
    }

    public function getUserById($userId=0){
        try{
            $this->db->where('id', $userId);
		    $query = $this->db->get('cc_User');
		    if(empty($query) || $query->num_rows() < 1){
                return false;
            }
            $userInfo = $query->first_row('array');
        }catch (Exception $ex){
            return false;
        }
        return $userInfo;
    }

    public function getUserList(){
        $userListArr = array();
        try{
		    $query = $this->db->get('cc_User');
		    if(empty($query) || $query->num_rows() < 1){
                return $userListArr;
            }
            $userListArr = $query->result_array();
            foreach($userListArr as &$userInfo){
                unset($userInfo['Password']);
            }
        }catch (Exception $ex){
            return $userListArr;
        }
        return $userListArr;
    }

    public function saveUser($userInfo = array()){
        try{
            $userInfo['Company'] = COMPANY_NAME;
            isset($userInfo['id']) && $user = $this->getUserById($userInfo['id']);
            if(isset($userInfo['UserName'])){
                $userNameExist = $this->getUserByUsername($userInfo['UserName']);
                if($userNameExist and $userNameExist['id'] != $userInfo['id']){
                    $this->_errInfo = '用户[' . $userInfo['UserName'] . ']已存在';
                    return false;
                }
            }
            if(isset($userInfo['id']) && (!empty($user)) && ($userInfo['id'] != '')){
                # 更新user
                $this->db->where('id', $userInfo['id']);
                $result = $this->db->update('User', $userInfo);
            }else{
                # 新增user
                unset($userInfo['id']);
                $result = $this->db->insert('User', $userInfo);
                $userInfo['id'] = $this->db->insert_id();
            }
            if(!$result) {
                $this->_errInfo = '保存用户信息失败!';
                $err = $this->db->error();
                CCLog::LogErr('保存用户信息失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
                return false;
            }
        }catch (Exception $ex){
            $this->_errInfo = '保存用户信息失败!';
            CCLog::LogErr('保存用户信息失败! 异常: ' . $ex);
            return false;
        }
        return $userInfo;
    }

    public function updateUserTokenExpire($username, $expire){
        try{
            # 更新user
            $this->db->where('UserName', $username);
            $result = $this->db->update('User', array('TokenExpire' => $expire));

            if(!$result) {
                $err = $this->db->error();
                CCLog::LogErr('保存用户token失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
                return false;
            }
        }catch (Exception $ex){
            CCLog::LogErr('保存用户token失败! 异常: ' . $ex);
            return false;
        }
        return true;
    }

    public function delUserById($userId=0){
        try{
            # 更新user
            $this->db->where('id', $userId);
            $result = $this->db->delete('User');

            if(!$result) {
                $this->_errInfo = '删除用户信息失败!';
                $err = $this->db->error();
                CCLog::LogErr('删除用户信息失败! mysql_errno: '. $err['code'] .', mysql_error: '. $err['message']);
                return false;
            }
        }catch (Exception $ex){
            $this->_errInfo = '删除用户信息失败!';
            CCLog::LogErr('删除用户信息失败! 异常: ' . $ex);
            return false;
        }
        return true;
    }

    public function updateUserCompany(){
        try{
            # 更新user
            $userInfo = array('Company' => COMPANY_NAME);
            $result = $this->db->update('User', $userInfo);
        }catch (Exception $ex){
            $this->_errInfo = '更新用户公司信息失败!';
            CCLog::LogErr('更新用户公司信息失败! 异常: ' . $ex);
            return false;
        }
        return true;
    }

    public function getErrMsg(){
        return $this->_errInfo;
    }

}