<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class Init extends CI_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     *  初始化数据
     */
    public function initUserData() {

        $this->load->Logic('UserLogic');
        $this->load->Logic('ApplicationBaseLogic');
        $this->load->Logic('SetBaseLogic');
        $this->load->Logic('ModuleBaseLogic');
        $this->load->database();

        $moduleQuery = $this->db->get('cc_ModuleBase');
        if(0 != $moduleQuery->num_rows()) {
            echo "skip init data".PHP_EOL;
            return;
        }
        $InitHostPropertyArr = array("INSERT INTO `cc_HostPropertyClassify` VALUES ('1', 'AssetID', '固资编号', 'basic', 'AssetID','12', '2016-02-24 11:26:57', '2016-02-24 18:00:57')",
                                     "INSERT INTO `cc_HostPropertyClassify` VALUES ('7', 'DeviceClass', '设备类型', 'basic', 'DeviceClass','11', '2016-02-24 17:24:04', '2016-02-24 18:01:24')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('8', 'HostName', '主机名称', 'basic', 'HostName', '6','2016-02-24 17:26:00', '2016-02-24 18:01:48')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('9', 'Status', '运行状态', 'basic', 'Status', '7','2016-02-24 18:02:23', '2016-02-25 14:45:38')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('10', 'Operator', '维护人', 'basic', 'Operator','3', '2016-02-24 18:02:41', '2016-02-24 18:03:14')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('11', 'BakOperator', '备份维护人', 'basic', 'BakOperator','4', '2016-02-24 18:03:37', '2016-02-24 18:03:37')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('12', 'InnerIP', '内网IP', 'basic', 'InnerIP','1', '2016-02-24 18:04:01', '2016-02-24 18:04:01')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('13', 'OuterIP', '外网IP', 'basic', 'OuterIP','2', '2016-02-24 18:04:31', '2016-02-24 18:04:31')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('14', 'OSName', '操作系统', 'basic', 'OSName','7', '2016-02-24 18:04:53', '2016-02-24 18:04:53')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('15', 'Description', '备注', 'basic', 'Description','13', '2016-02-24 18:05:10', '2016-02-24 18:05:10')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('16', 'ZoneName', '可用区', 'basic', 'ZoneName','15', '2016-02-24 18:05:39', '2016-02-24 18:05:39')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('17', 'ZoneID', '可用区ID', 'basic', 'ZoneID', '14','2016-02-24 18:06:07', '2016-02-24 18:06:07')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('47', 'CreateTime', '入库时间', 'basic', 'CreateTime', '17','2016-02-24 19:11:25', '2016-02-24 19:11:25')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('49', 'Region', '机房城市', 'basic', 'Region','16', '2016-02-24 19:12:21', '2016-02-24 19:12:21')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('51', 'Cpu', 'Cpu', 'basic', 'Cpu', '8','2016-02-24 19:13:12', '2016-02-24 19:13:12')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('52', 'Mem', '内存', 'basic', 'Mem', '9','2016-02-24 19:13:37', '2016-02-24 19:13:37')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('60', 'HostID', '主机ID', 'basic', 'HostID','0', '2016-02-24 19:16:54', '2016-02-24 19:16:54')",
                                    "INSERT INTO `cc_HostPropertyClassify` VALUES ('72', 'ModuleName', '模块名称', 'basic', 'ModuleName','5', '2016-02-24 18:02:41', '2016-02-24 18:02:41')");
        $InitSetProperty = "INSERT INTO `cc_SetProperty` VALUES (1,'SetEnviType','1','测试','2015-12-28 18:19:28','2015-12-28 18:19:28'),
                                                                 (2,'SetEnviType','2','体验','2015-12-28 18:19:38','2015-12-28 18:19:38'),
                                                                  (3,'SetEnviType','3','正式','2015-12-28 18:19:48','2015-12-28 18:19:48'),
                                                                  (4,'SetServiceStatus','0','关闭','2015-12-28 18:20:03','2015-12-28 18:20:03'),
                                                                  (5,'SetServiceStatus','1','开放','2015-12-28 18:20:14','2015-12-28 18:20:14');";
        echo "begin truncate table\n";
        $sql = 'truncate table cc_HostPropertyClassify';
        $this->db->query($sql);
        $sql = 'truncate table cc_HostCustomerProperty';
        $this->db->query($sql);
        $sql = 'truncate table cc_User';
        $this->db->query($sql);
        $sql = 'truncate table cc_ApplicationBase';
        $this->db->query($sql);
        $sql = 'truncate table cc_SetBase';
        $this->db->query($sql);
        $sql = 'truncate table cc_ModuleBase';
        $this->db->query($sql);
        $sql = 'truncate table cc_SetProperty';
        $this->db->query($sql);
        echo "end truncate table\n";

        echo "begin create classify host property \n";
        foreach($InitHostPropertyArr as $sql) {
            echo $sql."\n";
            $this->db->query($sql);
        }
        echo "end create classify host property  \n";

        echo "begin create set property \n";
        echo $sql."\n";
        $this->db->query($InitSetProperty);
        echo "end create set property  \n";

        echo "begin create admin user \n";
        $this->UserLogic->addDefaultUser();     //增加默认用户
        echo "end create admin user \n";

        echo "begin create resource pool \n";
        $this->ApplicationBaseLogic->addDefaultApp(COMPANY_NAME);       //增加默认业务
        echo "end create resource pool \n";

        echo "begin create example app \n";
        $appName = '示例业务';
        $type = 0;
        $maintainers = '_admin_';
        $productStr = '_admin_';
        $level = 3;
        $lifeCycle = '公测';
        $createTime = date('Y-m-d H:i:s');
        $creator = 'admin';
        $appResult = $this->ApplicationBaseLogic->addApplication($appName, $type, $maintainers, $productStr, $level, $lifeCycle, $createTime , $creator );       //增加示例业务
        $appId = $appResult['appId'];
        echo "end create example app \n";



        echo "begin create example set \n";
        $setName = '示例集群';
        $chnName = '示例集群';
        $enviType = 0;
        $serviceStatus =0;
        $capacity = 0;
        $des = '';
        $openStatus = 0;
        $setId = 0;
        $this->SetBaseLogic->newSet($appId, $setName, $chnName, $enviType, $serviceStatus, $capacity, $des, $openStatus, $setId);       //增加示例集群
        echo "end create example set \n";

        echo "begin create example module \n";
        $moduleName = '示例模块';
        $operator = 'admin';
        $bakOperator = 'admin';
        $this->ModuleBaseLogic->addModule($appId, $setId, $moduleName, $operator, $bakOperator);       //增加示例模块
        echo "end create example module \n";

    }


}