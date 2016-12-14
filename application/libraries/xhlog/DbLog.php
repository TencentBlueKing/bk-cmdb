<?php

class DbLog {
    public static $db;
    public static $key;
    public static $tbl;

    public static $log;
    public static $current = -1;

    /*
     * 设置数据库
     */
    public static function setDb($db, $tbl) {
        self::$db = $db;
        self::$tbl = $tbl;
        self::$key = $tbl;
    }

    /*
     * 记录日志
     */
    public static function log($log) {
        self::$current = self::$current + 1;
        self::$log[self::$key][self::$current] = $log;
        return true;
    }

    /*
     * 记录多条日志
     */
    public static function mutiLog($mutiLogs) {
        foreach ($mutiLogs as $line => $log) {
            self::$current = self::$current + $line + 1;
            self::$log[self::$key][self::$current] = $log;
        }
        return true;
    }

    /*
     * 写数据库日志
     */
    public static function flush() {
        if (empty(self::$log)) {
            return false;
        }
        foreach (self::$log as $tbl => $log) {
            if (empty($log)) {
                continue;
            }
            self::$db->insert_batch($tbl, $log);
        }
        return true;
    }
}

register_shutdown_function(array('DbLog', 'flush'));