<?php

class CCLog {
    private static $_startTime;    //开始时间

    private static $_endTime;        //结束时间

    private static $_excTime;    //执行时间

    private static $_companyCode;   //公司code

    private static $_userName;    //操作用户

    private static $_clientIP;    //用户IP

    private static $_ciObj;        //CI对象

    public function __construct() {
        self::$_ciObj = &get_instance();
        self::$_ciObj->load->library('session');
        self::$_companyCode = self::$_ciObj->session->userdata('company') ? self::$_ciObj->session->userdata('company') : '';
        self::$_userName = self::$_ciObj->session->userdata('username') ? self::$_ciObj->session->userdata('username') : '';
        self::$_clientIP = Utility::get_client_ip();
        self::$_startTime = gettimeofday();
    }

    /*
     * @记录错误日志
     */
    public static function LogErr($message) {
        return log_message('error', '[' . self::$_companyCode . ']  [' . self::$_userName . ']  [' . self::$_clientIP . ']  ' . $message);
    }

    /*
     * @记录Debug日志
     */
    public static function LogDebug($message) {
        return log_message('debug', '[' . self::$_companyCode . ']  [' . self::$_userName . ']  [' . self::$_clientIP . ']  ' . $message);
    }

    /*
     * @记录info日志
     */
    public static function LogInfo($message) {
        return log_message('info', '[' . self::$_companyCode . ']  [' . self::$_userName . ']  [' . self::$_clientIP . ']  ' . $message);
    }

    /*
     * @记录一般日志
     */
    public static function Log($message) {
        return log_message('all', '[' . self::$_companyCode . ']  [' . self::$_userName . ']  [' . self::$_clientIP . ']  ' . $message);
    }

    /*
     * @记录登陆日志
     */
    public static function addLogin() {
        require_once(APPPATH . '/libraries/xhlog/DbLog.php');

        $log = array(
                'UserName' => self::$_userName,
                'CompanyCode' => self::$_companyCode,
                'ClientIP' => Utility::get_client_ip(),
                'ServerIP' => Utility::get_server_ip(),
                'LastTime' => date('Y-m-d H:i:s'),
                'UserAgent' => $_SERVER['HTTP_USER_AGENT']);

        DbLog::setDb(self::$_ciObj->db, 'cc_UserLoginLog');
        DbLog::log($log);
    }

    /*
     * @记录页面访问日志
     */
    public static function addUrlVisit() {
        require_once(APPPATH . '/libraries/xhlog/DbLog.php');
        $urlArr = array_reverse(self::$_ciObj->uri->segment_array());

        if (empty($urlArr)) {
            $controller = 'index';
            $ction = 'index';
        } else if (1 === count($urlArr)) {
            $controller = strtolower($urlArr[0]);
            $action = 'index';
        } else {
            $action = strtolower($urlArr[0]);
            $controller = strtolower($urlArr[1]);
            $folder2 = isset($urlArr[2]) ? strtolower($urlArr[2]) : '';
            $folder1 = isset($urlArr[3]) ? strtolower($urlArr[3]) : '';
        }

        $log = array(
                'UserName' => self::$_userName,
                'Controller' => $controller,
                'Action' => $action,
                'Folder1' => $folder1,
                'Folder2' => $folder2,
                'ClientIP' => Utility::get_client_ip(),
                'ServerIP' => Utility::get_server_ip(),
                'LastTime' => date('Y-m-d H:i:s'),
                'CompanyCode' => self::$_companyCode);
        DbLog::setDb(self::$_ciObj->db, 'cc_UrlVisitLog');
        DbLog::log($log);
    }

    /*
     * @添加用户操作日志
     */
    public static function addOpLogArr($opArr) {
        require_once(APPPATH . '/libraries/xhlog/DbLog.php');
        self::$_ciObj->load->database();

        $opTime = date("Y-m-d H:i:s", $_SERVER['REQUEST_TIME']);
        $operator = self::$_userName;
        $companyCode = self::$_companyCode;

        self::$_endTime = gettimeofday();
        self::$_excTime = (self::$_endTime['sec'] - self::$_startTime['sec']) + (self::$_endTime['usec'] - self::$_startTime['usec']) / 1000;
        self::$_startTime = gettimeofday();

        $opResult = isset($opArr['OpResult']) ? $opArr['OpResult'] : 1;
        $opFrom = isset($opArr['opFrom']) ? $opArr['OpFrom'] : 0;
        $appId = isset($opArr['ApplicationID']) ? $opArr['ApplicationID'] : -1;
        $opName = isset($opArr['opName']) ? $opArr['OpName'] : $opArr['OpType'] . $opArr['OpTarget'];

        $data = array(
                    'OpTime' => $opTime,
                    'Operator' => $operator,
                    'ApplicationID' => $appId,
                    'CompanyCode' => $companyCode,
                    'OpTarget' => $opArr['OpTarget'],
                    'OpContent' => $opArr['OpContent'],
                    'OpName' => $opName,
                    'OpType' => $opArr['OpType'],
                    'ExecTime' => abs(self::$_excTime),
                    'OpResult' => $opResult,
                    'OpFrom' => $opFrom);
        DbLog::setDb(self::$_ciObj->db, 'cc_OperationLog');
        DbLog::log($data);
    }

}
