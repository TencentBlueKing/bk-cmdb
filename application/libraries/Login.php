<?php if(!defined('BASEPATH')) exit('No direct script access allowed');

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

/**
 * CC 登录鉴权类
 *
 * 登录态校验，登入，登出
 */
class Login {

	public $ci = NULL;
    private $_errorMsg = '';

	public function __construct(){
		$this->ci = & get_instance();
		$this->ci->load->library('session');
	}

    /**
     * 获取当前用户
     * @return array
     */
    public function getCurrentUser(){
        $userData = array();
        $userData['username'] = $this->ci->session->userdata('username');
        $userData['chinese_name'] = $this->ci->session->userdata('chinese_name');
        $userData['company'] = $this->ci->session->userdata('company');
        $userData['tel'] = $this->ci->session->userdata('tel');
        $userData['email'] = $this->ci->session->userdata('email');
        $userData['role'] = $this->ci->session->userdata('role');
        $userData['cc_token'] = $this->ci->session->userdata('cc_token');
        return $userData;
    }

    /**
     * 判断当前用户是否登录
     * @return bool
     */
    public function isAuthed(){
        $userInfo = $this->getCurrentUser();
        if(isset($userInfo['cc_token'])){
            $ccToken = $userInfo['cc_token'];
        }else{
            return $this->loginUser();
        }
        # 校验用户的cc_token是否有效
        return $this->isCcTokenValid($ccToken);
    }

    /**
     * 判断当前用户的CC TOKEN是否有效
     * @param null $ccToken
     * @return bool
     */
    public function isCcTokenValid($ccToken=""){
        # token解密
        $this->ci->load->library('encryption');
        $result = $this->ci->encryption->decrypt($ccToken);
        if(empty($result)){
            $this->_errorMsg = 'token非法';
            return false;
        }
        # 获取token有效期
        $tokenInfo = explode("|", $result);
        if(empty($tokenInfo) || count($tokenInfo) < 4){
            $this->_errorMsg = 'token非法';
            return false;
        }
        $expire = intval($tokenInfo[0]);
        $nowTime = time();
        # token有效期已过，或者token有效期大于当前时间24小时， 都为无效
        if(($nowTime - $expire > 0) || ($expire - $nowTime > 86400)){
            $this->_errorMsg = '登录态已过期';
            return false;
        }
        # 校验用户名和密码
        $username = $tokenInfo[1];
        $password = $tokenInfo[2];
        $userInfo = $this->isValidUserPassword($username, $password);
        if(!empty($userInfo)){
            $userInfo['Password'] = $password;
        }else{
            $this->_errorMsg = '用户名密码错误';
            return false;
        }
        if(is_numeric($userInfo['TokenExpire'])
            && ($userInfo['TokenExpire'] > 0)
            && ($expire < intval($userInfo['TokenExpire']))){
            $this->_errorMsg = '登录态已过期';
            return false;
        }
        return $userInfo;
    }

    /**
     * 判断用户的用户名和密码是否有效
     * @param $username
     * @param $password
     * @return bool $userInfo
     */
    public function isValidUserPassword($username, $password){
        if(empty($username)){
            return false;
        }
        # 从DB查询出用户信息
        $this->ci->load->model('UserModel');
        try {
            $userInfo = $this->ci->UserModel->getUserByUsername($username);
            if(empty($userInfo)){
                return false;
            }
            $passwordHash = $userInfo['Password'];
            if($passwordHash == $password){
                # 如果用户批量导入， 那么初始密码校验也可以通过
                return $userInfo;
            }
            $result = password_verify($password, $passwordHash);
            if(empty($result)){
                return false;
            }
            return $userInfo;
        }catch (Exception $ex){
            return false;
        }
    }

    /**
     * 生成CC TOKEN， 用于cookie记住登录用户
     * @param $username
     * @param $password
     * @return bool
     */
    public function genCcToken($username, $password){
        $salt = $this->genRandomPassword();
        $nowTime = time();
        $plainToken = ($nowTime + 86400) . '|' . $username . '|' . $password . '|' . $salt;
        $this->ci->load->library('encryption');
        $token = $this->ci->encryption->encrypt($plainToken);
        return $token;
    }

    /**
     * 生成随机密码
     * @return string
     */
    public function genRandomPassword($length = 8){
        $randomPassword = '';

        # 密码字符集
        $passwordCharSet = array('a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
        'i', 'j', 'k', 'l','m', 'n', 'o', 'p', 'q', 'r', 's',
        't', 'u', 'v', 'w', 'x', 'y','z', 'A', 'B', 'C', 'D',
        'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L','M', 'N', 'O',
        'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y','Z',
        '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!',
        '@','#', '$', '%', '^', '&', '*', '(', ')', '-', '_',
        '[', ']', '{', '}', '<', '>', '~', '`', '+', '=', ',',
        '.', ';', ':', '/', '?', '|');

        $passwordKeys = array_rand($passwordCharSet, $length);
        for($i = 0; $i < $length; $i++){
            $randomPassword .= $passwordCharSet[$passwordKeys[$i]];
        }
        return $randomPassword;
    }

    /**
     * 将用户在app登录
     * @param null $userInfo
     * @return bool
     */
    public function loginUser(){
        $username = $this->ci->input->get_post("UserName");
        $password = $this->ci->input->get_post("Password");
        $userInfo = $this->isValidUserPassword($username, $password);
        if(empty($userInfo)){
            # 尝试cc_token登录
            $ccToken = $this->ci->input->cookie('cc_token');
            $userInfo = $this->isCcTokenValid($ccToken);
            if(empty($userInfo)){
                $this->_errorMsg = '用户名密码错误';
                return false;
            }else{
                $username = $userInfo['UserName'];
                $password = $userInfo['Password'];
            }
        }
        # 给用户生成cc token
        $token = $this->genCcToken($username, $password);
        # 更新session中用户信息
        $userData = array();
        $userData['user_id'] = $userInfo['id'];
        $userData['username'] = $userInfo['UserName'];
        $userData['chinese_name'] = $userInfo['ChName'];
        $userData['company'] = $userInfo['Company'];
        $userData['company_id'] = COMPANY_ID;
        $userData['tel'] = $userInfo['Tel'];
        $userData['email'] = $userInfo['Email'];
        $userData['role'] = $userInfo['Role'];
        $userData['cc_token'] = $token;
        $this->ci->session->set_userdata($userData);
        # 记录cc token到cookie
        $this->ci->input->set_cookie('cc_token', $token, 86400, COOKIE_DOMAIN);//token每天登录一次

        # 给用户初始化业务信息
        $this->ci->load->Logic('ApplicationBaseLogic');
        $appIdArr = $this->ci->ApplicationBaseLogic->getAppIdOnLogin($userData['username']);
        # 没有默认业务则新增默认业务
        if(empty($appIdArr)) {
            $dftApp = $this->ci->ApplicationBaseLogic->addDefaultApp($userData['company']);
            if(empty($dftApp)){
                CCLog::LogInfo('为用户' . $userData['username'] .'新增默认业务失败');
                $this->_errorMsg = '为用户' . $userData['username'] .'新增默认业务失败';
                return false;
            }
            $appIdArr = array($dftApp['AppID']);
        }
        # 如果用户只有一个默认业务， 那么为新用户
        $newUser = count($appIdArr) == 1 ? 1 : 0;
        $userData['appId'] = $appIdArr;
        $userData['newUser'] = $newUser;
        $userData['app'] = array();
        $this->ci->session->set_userdata($userData);

        $appInfoArr = $this->ci->ApplicationBaseLogic->getUserApp();
        $appHostCount = $this->ci->ApplicationBaseLogic->getHostNumByAppIdArr($appIdArr);
        $dftCookieAppId = $ccToken = $this->ci->input->cookie('cc_token');
        $dftCookieAppId = (empty($dftCookieAppId)) ? '' : $dftCookieAppId;
        $dftAppInfo = array();
        foreach ($appInfoArr as $appInfo) {
            $appId = $appInfo['ApplicationID'];
            $appName = $appInfo['ApplicationName'];
            if(($appInfo['Default'] == 0) && (empty($dftAppInfo) && ($dftCookieAppId == $appInfo['ApplicationID']))) {
                $dftAppInfo['ApplicationID'] = $appId;
                $dftAppInfo['ApplicationName'] = $appName;
            }

            if($appInfo['Default'] != 1){
                $userData['app'][$appId]['ApplicationID'] = $appId;
                $userData['app'][$appId]['ApplicationName'] = $appName;
                $userData['app'][$appId]['ApplicationHostCount'] = $appHostCount[$appId];
            }
        }
        # 如果没有cookie中记录的默认app， 那么使用第一个做为用户的默认业务
        if(empty($dftAppInfo)){
            $dftAppInfo = reset($userData['app']);
        }
        $userData['defaultApp'] = $dftAppInfo;
        $expires = time() + $this->ci->config->item('sess_expiration')*1000;
        $this->ci->input->set_cookie('defaultAppId', $dftAppInfo['ApplicationID'], $expires);
        $this->ci->input->set_cookie('defaultAppName', $dftAppInfo['ApplicationName'], $expires);
        $this->ci->session->set_userdata($userData);

        # 为ajax请求生成csrf验证token
        $this->generateToken();

        # 记录登录日志
        CCLog::addLogin();
        return true;
    }

    /**
     * 清除掉用户的app登录态
     * @return bool
     */
    public function logoutUser(){
        $curUser = $this->getCurrentUser();
        $userName = $curUser['username'];
        # 更新用户的token失效时间， 把当前的token都进行失效
        $this->ci->load->model('UserModel');
        $this->ci->UserModel->updateUserTokenExpire($userName, time() + 86400);
        # 清除cookie
        $this->ci->load->helper('cookie');
        delete_cookie('token');    //csrf token
        delete_cookie('cc_token', COOKIE_DOMAIN); //cc token
        # 清除session
        $this->ci->session->sess_destroy();
        return true;
    }

    /**
     * 用户修改密码
     * @param null $userInfo
     * @return bool
     */
    public function changePassword($username, $newPassword){
        try{
            $userInfo = $this->getCurrentUser();
            $role = $userInfo['role'];
            $curUser = $userInfo['username'];
            if(($role != 'admin') && ($curUser != $username)){
                $this->_errorMsg = '您无权修改此用户的密码';
                return false;
            }
            # 存储新密码
            $this->ci->load->model('UserModel');
            $userInfo = $this->ci->UserModel->getUserByUsername($username);
            if(empty($userInfo)){
                $this->_errorMsg = '用户(' . $username . ')不存在， 不能修改密码';
                return false;
            }
            $userInfo['Password'] = password_hash($newPassword, PASSWORD_DEFAULT);
            if(!password_verify($newPassword, $userInfo['Password'])){
                $this->_errorMsg = '修改密码失败';
                return false;
            }
            if(!$this->ci->UserModel->saveUser($userInfo)){
                $this->_errorMsg = '保存密码失败';
                return false;
            }
        }catch (Exception $ex){
            $this->_errorMsg = '修改密码失败';
            return false;
        }
        return true;
    }

    /**
     * 用户登录态无效，跳转
     */
    public function redirectUnauthed(){
        $this->ci->load->helper('url');
        $loginUrl = LOGIN_URL;
        $callbackUrl = uri_string();
        $jumpUrl = $loginUrl;
        if($callbackUrl != 'account/logout'){
            $jumpUrl .= '?cburl=' . urlencode($callbackUrl);
        }
        redirect($jumpUrl, 'location', 301);
    }

	/**
	 * 用户登录时，生成csrf token
	 */
	private function generateToken(){
		try{
			$rand = mt_rand(0, mt_getrandmax());
			$rand .= $this->ci->input->ip_address();
			$token = md5(uniqid($rand, TRUE));
			$this->ci->input->set_cookie('token', $token, $this->ci->config->item('sess_expiration')*200);//cookie过期时间与会话过期时间一致
			$this->ci->session->set_userdata(array('token'=>$token));
		}catch(Exception $e){
			CCLog::LogErr('generateToken exception'.$e->getMessage());
		}
	}

    public function getErrMsg(){
        return $this->_errorMsg;
    }
}