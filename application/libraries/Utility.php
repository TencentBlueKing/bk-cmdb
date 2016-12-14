<?php if(!defined('BASEPATH')) exit('No direct script access allowed');

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

/**
 * 工具库
 */
class Utility  {

	public static $errInfo = '';

    /*
     * @获取客户端IP
     */
	public static function get_client_ip() {
		if(!empty($_SERVER['HTTP_CLIENT_IP'])) {
			$ip = $_SERVER['HTTP_CLIENT_IP'];
		}elseif (!empty($_SERVER['HTTP_X_FORWARDED_FOR'])) {
			$ip = $_SERVER['HTTP_X_FORWARDED_FOR'];
		}else {
			$ip = isset($_SERVER['REMOTE_ADDR']) ? $_SERVER['REMOTE_ADDR'] : '';
		}
		return $ip;
	}

    /*
     * @获取服务端IP
     */
	public static function get_server_ip() {
		if(isset ( $_SERVER )){
			if ($_SERVER ['SERVER_ADDR']) {
				$serverIp = $_SERVER ['SERVER_ADDR'];
			}else {
				$serverIp = $_SERVER ['LOCAL_ADDR'];
			}
		}else {
			$serverIp = getenv ( 'SERVER_ADDR' );
		}
		return $serverIp;
	}

	/*
	 * curl http请求
	 * @param $url 请求地址
	 * @param $data 请求参数数组array('a'=>1,'b'=>2)
	 * @param $type 请求方式，默认get请求
	 * @return array or boolean
	*/
	public static function http($url, $data, $type='get') {
		$parseUrl = parse_url($url);
		$path = rtrim($parseUrl['path'],'/');
		$key = 'signature';
		if(strpos($url, 'compapi') === false){
			$path = trim(trim($url, 'http://'),'/');
			$key = 'Signature';
		}

		$signature = self::genSignature($path, $data);
        $data[$key] = urlencode($signature);

        $ch = curl_init();
        if($type === 'get') {
            $requestStr = '';
			foreach($data as $dk=>$dv) {
				$requestStr .= '&'. $dk .'='. $dv;
			}
            $url = trim($url, '?') .'?' . trim($requestStr, '&');
        }elseif($type === 'newpost') {
            curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));
            curl_setopt($ch, CURLOPT_HTTPHEADER, array(
                'Content-Type: application/json',
                'Content-Length: ' . strlen(json_encode($data)))
            );
            $url .= '?bk_nonce=' . $data['bk_nonce'] . '&bk_timestamp=' . $data['bk_timestamp'];
        }else {
            curl_setopt($ch, CURLOPT_POST, 1);
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        }

        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_TIMEOUT, 60);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        $result = curl_exec($ch);
        CCLog::LogInfo('curl url:'.print_r($url, true));
        CCLog::LogInfo('curl result:'.print_r($result, true));
        if($result === false) {
            CCLog::LogErr('curl url:'.print_r($url, true));
            CCLog::LogErr('curl result:'.print_r($result, true));
            self::$errInfo = 'error: '.curl_error($ch) .'###request info: '. print_r(curl_getinfo($ch), true);
            return false;
        }

        curl_close($ch);
        $data = json_decode($result, true);

        if($data === false){
            CCLog::LogErr("curl result:".print_r($result, true));
            CCLog::LogInfo("curl url:".print_r($url, true));
            self::$errInfo = 'error: json_decode error!###request info: '. print_r(curl_getinfo($result), true);
            return false;
        }
        return $data;
    }

    /*
     * @生成签名
     */
	public static function genSignature($path, $data){
		$requestStr = '';
		foreach($data as $dk=>$dv){
			$requestStr .= '&'. $dk .'='. urldecode($dv);
		}
		$requestStr = trim($requestStr, '&');
		$message = 'GET'.$path.'/?'.$requestStr;
		return base64_encode(hash_hmac('sha1', $message, BKAPI_APP_SECRET, true));
	}

	/**
	* @计算excel列名
	* @param num int 列数
	* @return arrray 读取的数据
	*/
	public static function getCol($num){
        /*
         * 递归方式实现根据列数返回列的字母标识
        */
        $arr = array(0=>'Z', 1=>'A', 2=>'B', 3=>'C', 4=>'D', 5=>'E', 6=>'F', 7=>'G', 8=>'H', 9=>'I', 10=>'J', 11=>'K', 12=>'L', 13=>'M', 14=>'N', 15=>'O', 16=>'P', 17=>'Q', 18=>'R', 19=>'S', 20=>'T', 21=>'U', 22=>'V', 23=>'W', 24=>'X', 25=>'Y', 26=>'Z');
        if ($num == 0) {
            return '';
        }
        return self::getCol((int)(($num - 1) / 26)) . $arr[$num % 26];
    }

    /*
     * @todo 时间转换
     */
    public static function tranTime($time) {
        $rTime = date("Y-m-d H:i:s", strtotime($time));
        $hTime = date("H:i", strtotime($time));
        $time = time() - strtotime($time);
        if ($time < 60) {
            $str = '刚刚';
        } elseif ($time < 60 * 60) {
            $min = floor($time/60);
            $str = $min.'分钟前';
        } elseif ($time < 60 * 60 * 24) {
            $h = floor($time/(60*60));
            $str = $h.'小时前 '.$hTime;
        } elseif ($time < 60 * 60 * 24 * 3) {
            $d = floor($time/(60*60*24));
            if($d==1){
                $str = '昨天 '.$rTime;
            }else{
                $str = '前天 '.$rTime;
            }
        }else{
            $str = $rTime;
        }

        return $str;
    }

    /*
     * @纯净版的curl请求，
     * @param $url 请求url，url可带query_string
     * @param $data array post请求的body
     * @return array or false 平台返回的数据或者false
     */
    public static function post($url, $data,$method=true) {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        if($method) {
            curl_setopt($ch, CURLOPT_POST, 1);
            curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        }
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        $result = curl_exec($ch);
        CCLog::LogInfo("curl url:".print_r($url, true));
        CCLog::LogInfo("curl result:".print_r($result, true));
        if($result === false){
            CCLog::LogErr("curl url:".print_r($url, true));
            CCLog::LogErr("curl result:".print_r($result, true));
            return false;
        }
        curl_close($ch);
        $data = json_decode($result, true);
        if($data === false){
            CCLog::LogErr("curl result:".print_r($result, true));
            CCLog::LogInfo("curl url:".print_r($url, true));
            return false;
        }
        return $data;
    }

    /*
     * @时间格式验证
     * @param  string
     * @return fool
     */
    public function validateDate($date, $format = 'Y-m-d H:i:s') {
        $d = DateTime::createFromFormat($format, $date);
        return $d && $d->format($format) == $date;
    }

    /*
     * @过滤掉字段中的_
     */
    public static function filterArrayFields(&$inputArr, $fieldsArr) {
        foreach($inputArr as &$input) {
            foreach($fieldsArr as $key) {
                isset($input[$key]) && $input[$key] = str_replace('_', '', $input[$key]);
            }
        }
    }

    /*
     * @获取数字字段
     */
    public static function getNumericInArray(&$inputArr) {
        foreach($inputArr as $key=>$value) {  //封装方法
            if(floatval($value) > 0) {
                $inputArr[$key] = floatval($value);
            }else {
                unset($inputArr[$key]);
            }
        }
    }


/*
     * @字符串验证
     * @param  string
     * @return bool
     */
    public static function validateStringLen($string, $min=0, $max=0) {
        # 判断是否为字符串
        if(!is_string($string)){
            return false;
        }
        # 判断长度是否合法
        $stringLen = mb_strlen($string);
        if($min > 0 && ($stringLen < $min)){
            return false;
        }
        if($max > 0 && ($stringLen > $max)){
            return false;
        }
        return true;
    }

    /*
     * @字符串正则验证
     * @param  string
     * @return bool
     */
    public static function validateStringReg($string, $reg) {
        # 判断是否为字符串
        if(empty($string) || (!is_string($string))){
            return false;
        }
        if(preg_match($reg, $string) < 1){
            return false;
        }
        return true;
    }

    /*
     * @输入值验证
     * @param  string
     * @return bool
     */
    public static function validateValue($input, $values) {
        return in_array($input, $values);
    }

    /*
     * @输入校验
     * @param  string
     * @return bool
     */
    public static function validateInput($input, $config) {
        $validateType = self::array_get($config, 'type', 'len');
        switch($validateType){
            case 'len':
                $max = self::array_get($config, 'max', 0);
                $min = self::array_get($config, 'min', 0);
                return self::validateStringLen($input, $min, $max);
            case 'reg':
                $reg = self::array_get($config, 'reg', '/./is');
                return self::validateStringReg($input, $reg);
            case 'value':
                $values = self::array_get($config, 'values', array());
                return self::validateValue($input, $values);
            default:
                return false;
        }
    }

    /*
     * @获取数组值， 如果未定义的key则返回默认值
     * @param  string
     * @return bool
     */
    public static function array_get($array, $key, $dft) {
        return isset($array[$key]) ? $array[$key] : $dft;
    }

}
