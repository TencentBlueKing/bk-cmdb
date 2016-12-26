
DROP TABLE IF EXISTS `cc_ApplicationBase`;
CREATE TABLE `cc_ApplicationBase` (
  `ApplicationID` int(11) NOT NULL auto_increment,
  `ApplicationName` varchar(64) NOT NULL default '',
  `Creator` varchar(16) NOT NULL default '',
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `Default` int(1) NOT NULL default '0',
  `DeptName` varchar(64) NOT NULL default '',
  `Description` varchar(256) NOT NULL default '',
  `Display` int(1) NOT NULL default '1',
  `GroupName` varchar(64) NOT NULL default '',
  `LifeCycle` varchar(16) NOT NULL default '',
  `Maintainers` varchar(512) default NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `Level` int(1) NOT NULL default '2',
  `Owner` varchar(16) NOT NULL default '',
  `ProductPm` varchar(128) NOT NULL default '',
  `Type` int(1) NOT NULL default '0',
  `Source` varchar(16) NOT NULL default '',
  `CompanyID` int(11) NOT NULL default '0',
  `BusinessDeptName` varchar(64) NOT NULL default '',
  PRIMARY KEY  (`ApplicationID`),
  KEY `i_ApplicationName` (`ApplicationName`),
  KEY `i_Creator` (`Creator`),
  KEY `i_Owner` (`Owner`),
  KEY `i_Maintainers` (`Maintainers`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='业务基础信息表';


DROP TABLE IF EXISTS `cc_BaseParameterData`;
CREATE TABLE `cc_BaseParameterData` (
  `ParameterID` int(11) NOT NULL auto_increment,
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `DataType` varchar(50) NOT NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `ParameterCode` varchar(50) NOT NULL default '',
  `ParameterName` varchar(50) NOT NULL default '',
  `ParentCode` varchar(50) NOT NULL default '',
  PRIMARY KEY  (`ParameterID`),
  KEY `i_DatraType` (`DataType`),
  KEY `i_ParameterCode` (`ParameterCode`),
  KEY `i_ParameterName` (`ParameterName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;


DROP TABLE IF EXISTS `cc_HostBase`;
CREATE TABLE `cc_HostBase` (
  `HostID` int(11) NOT NULL auto_increment,
  `AssetID` varchar(64) NOT NULL default '',
  `BakOperator` varchar(16) NOT NULL default '',
  `Cpu` int(3) NOT NULL default '0',
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `Description` varchar(256) NOT NULL default '',
  `DeviceClass` varchar(32) NOT NULL default '',
  `HardMemo` varchar(512) NOT NULL default '',
  `HostName` varchar(32) NOT NULL default '',
  `IdcName` varchar(128) NOT NULL default '',
  `InnerIP` varchar(128) NOT NULL default '',
  `LastTime` timestamp NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `Mem` int(5) NOT NULL default '0',
  `Operator` varchar(16) NOT NULL default '',
  `OSName` varchar(32) NOT NULL default '',
  `OuterIP` varchar(128) NOT NULL default '',
  `PosCode` int(3) NOT NULL default '0',
  `Region` varchar(8) NOT NULL default '',
  `ServerRack` varchar(16) NOT NULL default '',
  `SN` varchar(32) NOT NULL default '',
  `Source` varchar(16) NOT NULL default '',
  `Status` varchar(32) NOT NULL default '',
  `ZoneID` int(11) NOT NULL default '0',
  `ZoneName` varchar(100) NOT NULL default '',
  `GseProxy` int(1) default '0',
  `Extend001` varchar(255) NOT NULL default '',
  `Extend002` varchar(255) NOT NULL default '',
  `Extend003` varchar(255) NOT NULL default '',
  `Extend004` varchar(255) NOT NULL default '',
  `Extend005` varchar(255) NOT NULL default '',
  `Customer001` varchar(255) NOT NULL default '',
  `Customer002` varchar(255) NOT NULL default '',
  `Customer003` varchar(255) NOT NULL default '',
  `Customer004` varchar(255) NOT NULL default '',
  `Customer005` varchar(255) NOT NULL default '',
  `Customer006` varchar(255) NOT NULL default '',
  `Customer007` varchar(255) NOT NULL default '',
  `Customer008` varchar(255) NOT NULL default '',
  `Customer009` varchar(255) NOT NULL default '',
  `Customer010` varchar(255) NOT NULL default '',
  `Customer011` varchar(255) NOT NULL default '',
  `Customer012` varchar(255) NOT NULL default '',
  `Customer013` varchar(255) NOT NULL default '',
  `Customer014` varchar(255) NOT NULL default '',
  `Customer015` varchar(255) NOT NULL default '',
  `Customer016` varchar(255) NOT NULL default '',
  `Customer017` varchar(255) NOT NULL default '',
  `Customer018` varchar(255) NOT NULL default '',
  `Customer019` varchar(255) NOT NULL default '',
  `Customer020` varchar(255) NOT NULL default '',
  `Customer021` varchar(255) NOT NULL default '',
  `Customer022` varchar(255) NOT NULL default '',
  `Customer023` varchar(255) NOT NULL default '',
  `Customer024` varchar(255) NOT NULL default '',
  `Customer025` varchar(255) NOT NULL default '',
  `Customer026` varchar(255) NOT NULL default '',
  `Customer027` varchar(255) NOT NULL default '',
  `Customer028` varchar(255) NOT NULL default '',
  `Customer029` varchar(255) NOT NULL default '',
  `Customer030` varchar(255) NOT NULL default '',
  `Customer031` varchar(255) NOT NULL default '',
  `Customer032` varchar(255) NOT NULL default '',
  `Customer033` varchar(255) NOT NULL default '',
  `Customer034` varchar(255) NOT NULL default '',
  `Customer035` varchar(255) NOT NULL default '',
  `Customer036` varchar(255) NOT NULL default '',
  `Customer037` varchar(255) NOT NULL default '',
  `Customer038` varchar(255) NOT NULL default '',
  `Customer039` varchar(255) NOT NULL default '',
  `Customer040` varchar(255) NOT NULL default '',
  `Customer041` varchar(255) NOT NULL default '',
  `Customer042` varchar(255) NOT NULL default '',
  `Customer043` varchar(255) NOT NULL default '',
  `Customer044` varchar(255) NOT NULL default '',
  `Customer045` varchar(255) NOT NULL default '',
  `Customer046` varchar(255) NOT NULL default '',
  `Customer047` varchar(255) NOT NULL default '',
  `Customer048` varchar(255) NOT NULL default '',
  `Customer049` varchar(255) NOT NULL default '',
  `Customer050` varchar(255) NOT NULL default '',
  PRIMARY KEY  (`HostID`),
  KEY `i_AssetID` (`AssetID`),
  KEY `i_BakOperator` (`BakOperator`),
  KEY `i_HostName` (`HostName`),
  KEY `i_InnerIP` (`InnerIP`),
  KEY `i_Operator` (`Operator`),
  KEY `i_OuterIP` (`OuterIP`),
  KEY `i_Source` (`Source`),
  KEY `i_SN` (`SN`),
  KEY `i_CreateTime` (`CreateTime`),
  KEY `i_DeviceClass` (`DeviceClass`),
  KEY `i_OSName` (`OSName`),
  KEY `i_Region` (`Region`),
  KEY `i_Status` (`Status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='主机基础信息表';



DROP TABLE IF EXISTS `cc_HostCustomerProperty`;
CREATE TABLE `cc_HostCustomerProperty` (
  `ID` int(11) NOT NULL auto_increment,
  `PropertyKey` varchar(25) NOT NULL,
  `PropertyName` varchar(25) NOT NULL,
  `Group` varchar(16) NOT NULL,
  `HostTableField` varchar(16) NOT NULL,
  `Owner` varchar(16) NOT NULL,
  `CreateTime` datetime NOT NULL,
  `LastTime` datetime NOT NULL,
  PRIMARY KEY  (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



DROP TABLE IF EXISTS `cc_HostPropertyClassify`;
CREATE TABLE `cc_HostPropertyClassify` (
  `ID` int(11) NOT NULL auto_increment,
  `PropertyKey` varchar(25) NOT NULL,
  `PropertyName` varchar(25) NOT NULL,
  `Group` varchar(16) NOT NULL default '',
  `HostTableField` varchar(16) NOT NULL default '',
  `Order` int(2) NOT NULL default '0',
  `CreateTime` datetime NOT NULL,
  `LastTime` datetime NOT NULL,
  PRIMARY KEY  (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



DROP TABLE IF EXISTS `cc_HostSource`;
CREATE TABLE `cc_HostSource` (
  `ID` int(11) NOT NULL auto_increment,
  `SourceCode` varchar(32) character set utf8 NOT NULL default '''''',
  `SourceName` varchar(128) character set utf8 NOT NULL default '''''',
  `IsPublic` int(1) NOT NULL default '1',
  `CompanyCode` varchar(16) character set utf8 NOT NULL default '''''',
  `CreateTime` datetime NOT NULL,
  `LastTime` datetime NOT NULL,
  PRIMARY KEY  (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



DROP TABLE IF EXISTS `cc_ModuleBase`;
CREATE TABLE `cc_ModuleBase` (
  `ModuleID` int(11) NOT NULL auto_increment,
  `ApplicationID` int(11) NOT NULL default '0',
  `BakOperator` varchar(16) NOT NULL default '',
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `Default` int(1) NOT NULL default '0',
  `Description` varchar(256) NOT NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `ModuleName` varchar(64) NOT NULL default '',
  `Operator` varchar(16) NOT NULL default '',
  `SetID` int(11) NOT NULL default '0',
  PRIMARY KEY  (`ModuleID`),
  KEY `i_ApplicationID` (`ApplicationID`),
  KEY `i_ModuleName` (`ModuleName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='模块基础信息表';


DROP TABLE IF EXISTS `cc_ModuleHostConfig`;
CREATE TABLE `cc_ModuleHostConfig` (
  `ID` int(11) NOT NULL auto_increment,
  `ApplicationID` int(11) NOT NULL default '0',
  `Description` varchar(256) NOT NULL default '',
  `HostID` int(11) NOT NULL default '0',
  `ModuleID` int(11) NOT NULL default '0',
  `SetID` int(11) NOT NULL default '0',
  PRIMARY KEY  (`ID`),
  KEY `i_ApplicationID` (`ApplicationID`),
  KEY `i_HostID` (`HostID`),
  KEY `i_ModuleID` (`ModuleID`),
  KEY `i_SetID` (`SetID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='模块与主机绑定关系表';


DROP TABLE IF EXISTS `cc_OperationLog`;
CREATE TABLE `cc_OperationLog` (
  `ID` int(11) NOT NULL auto_increment,
  `ApplicationID` int(11) NOT NULL default '0',
  `CompanyCode` varchar(16) NOT NULL default '',
  `Description` varchar(256) NOT NULL default '',
  `ExecTime` double NOT NULL default '0',
  `ClientIP` varchar(16) NOT NULL default '',
  `OpContent` text NOT NULL,
  `OpFrom` int(1) NOT NULL default '0',
  `Operator` varchar(16) NOT NULL default '',
  `OpName` varchar(32) NOT NULL default '',
  `OpTarget` varchar(32) NOT NULL default '',
  `OpResult` int(1) NOT NULL default '0',
  `OpTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `OpType` varchar(16) NOT NULL default '',
  `WebSys` varchar(16) NOT NULL default '',
  PRIMARY KEY  (`ID`),
  KEY `i_OpTime` (`OpTime`),
  KEY `i_Operator` (`Operator`),
  KEY `i_OpName` (`OpName`),
  KEY `i_ApplicationID` (`ApplicationID`),
  KEY `i_OpResult` (`OpResult`),
  KEY `i_OpTarget` (`OpTarget`),
  KEY `i_OpType` (`OpType`),
  KEY `i_WebSys` (`WebSys`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户操作日志表';


DROP TABLE IF EXISTS `cc_SetBase`;
CREATE TABLE `cc_SetBase` (
  `SetID` int(11) NOT NULL auto_increment,
  `ApplicationID` int(11) NOT NULL default '0',
  `Default` int(1) NOT NULL default '0',
  `Capacity` int(11) unsigned default '0',
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `ChnName` varchar(32) NOT NULL default '',
  `Description` varchar(256) NOT NULL default '',
  `EnviType` varchar(16) NOT NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `ParentID` int(11) NOT NULL default '0',
  `SetName` varchar(64) NOT NULL default '',
  `ServiceStatus` varchar(16) NOT NULL default '',
  `Openstatus` varchar(16) NOT NULL default '',
  PRIMARY KEY  (`SetID`),
  KEY `i_ApplicationID` (`ApplicationID`),
  KEY `i_SetID` (`SetID`),
  KEY `i_SetName` (`SetName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='大区基础信息表';


DROP TABLE IF EXISTS `cc_SetProperty`;
CREATE TABLE `cc_SetProperty` (
  `ID` int(11) NOT NULL auto_increment,
  `PropertyType` varchar(32) NOT NULL default '',
  `PropertyCode` varchar(32) NOT NULL default '',
  `PropertyName` varchar(128) NOT NULL default '',
  `CreateTime` datetime NOT NULL,
  `LastTime` datetime NOT NULL,
  PRIMARY KEY  (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



DROP TABLE IF EXISTS `cc_SysPermissions`;
CREATE TABLE `cc_SysPermissions` (
  `ID` int(11) NOT NULL auto_increment,
  `Action` varchar(32) NOT NULL default '',
  `Controller` varchar(32) NOT NULL default '',
  `Folder1` varchar(32) NOT NULL default '',
  `Folder2` varchar(32) NOT NULL default '',
  `Admin` int(1) NOT NULL default '1',
  `Qcloud` int(1) NOT NULL default '1',
  `Tencent` int(1) NOT NULL default '1',
  PRIMARY KEY  (`ID`),
  KEY `i_Action` (`Action`),
  KEY `i_Controller` (`Controller`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='系统权限表';



DROP TABLE IF EXISTS `cc_UrlVisitLog`;
CREATE TABLE `cc_UrlVisitLog` (
  `ID` int(11) NOT NULL auto_increment,
  `Action` varchar(16) NOT NULL default '',
  `Controller` varchar(16) NOT NULL default '',
  `ClientIP` varchar(16) NOT NULL default '',
  `CompanyCode` varchar(16) NOT NULL default '',
  `Description` varchar(256) NOT NULL default '',
  `Folder1` varchar(16) NOT NULL default '',
  `Folder2` varchar(16) NOT NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `ServerIP` varchar(16) NOT NULL default '',
  `UserName` varchar(16) NOT NULL default '',
  `WebSys` varchar(16) NOT NULL default '',
  PRIMARY KEY  (`ID`),
  KEY `i_UserName` (`UserName`),
  KEY `i_CompanyCode` (`CompanyCode`),
  KEY `i_Controller` (`Folder1`,`Folder2`,`Controller`,`Action`),
  KEY `i_LastTime` (`LastTime`),
  KEY `i_WebSys` (`WebSys`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='url访问日志表';


DROP TABLE IF EXISTS `cc_User`;
CREATE TABLE `cc_User` (
  `id` int(10) unsigned NOT NULL auto_increment,
  `UserName` varchar(128) NOT NULL default '',
  `Password` varchar(128) NOT NULL default '',
  `ChName` varchar(256) NOT NULL default '',
  `Company` varchar(128) NOT NULL default '',
  `Tel` varchar(45) NOT NULL,
  `QQ` varchar(45) NOT NULL default '',
  `Email` varchar(128) NOT NULL default '',
  `Role` enum('admin','user') NOT NULL default 'user',
  `Status` enum('ok','disabled') NOT NULL default 'ok',
  `TokenExpire` int(11) NOT NULL default '0',
  PRIMARY KEY  (`id`),
  UNIQUE KEY `UserName_UNIQUE` (`UserName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `cc_UserCustom`;
CREATE TABLE `cc_UserCustom` (
  `UserName` varchar(16) NOT NULL default '',
  `DefaultApplication` int(11) NOT NULL default '0',
  `DefaultColumn` text NOT NULL ,
  `DefaultPageSize` int(2) NOT NULL default '20',
  `DefaultField` varchar(512) NOT NULL default '' COMMENT '主机查询字段',
  `DefaultCon` text NOT NULL  COMMENT '主机查询条件',
  `Description` varchar(256) NOT NULL default '',
  `SetGseCol` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY  (`UserName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户定制表';


DROP TABLE IF EXISTS `cc_UserLoginLog`;
CREATE TABLE `cc_UserLoginLog` (
  `ID` int(10) unsigned NOT NULL auto_increment,
  `ClientIP` varchar(16) NOT NULL default '',
  `Description` varchar(256) NOT NULL default '',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `ServerIP` varchar(16) NOT NULL default '',
  `UserName` varchar(16) NOT NULL default '',
  `UserAgent` varchar(512) NOT NULL default '',
  `WebSys` varchar(16) NOT NULL default '',
  `CompanyCode` varchar(16) NOT NULL default '',
  PRIMARY KEY  (`ID`),
  KEY `i_CompanyCode` (`WebSys`),
  KEY `i_UserName` (`UserName`),
  KEY `i_LastTime` (`LastTime`),
  KEY `i_WebSys` (`WebSys`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户登录日志表';


DROP TABLE IF EXISTS `cc_UserPermissions`;
CREATE TABLE `cc_UserPermissions` (
  `ID` int(11) NOT NULL auto_increment,
  `ChName` varchar(32) NOT NULL default '',
  `CreateTime` datetime NOT NULL default '1970-01-01 00:00:00',
  `Description` varchar(256) NOT NULL default '',
  `GroupID` int(1) NOT NULL default '2',
  `GroupName` varchar(16) NOT NULL default '业务运维',
  `LastTime` timestamp NOT NULL default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP,
  `OwnerName` varchar(64) NOT NULL default '',
  `OwnerUin` varchar(16) NOT NULL default '',
  `ParentUin` varchar(16) NOT NULL default '',
  `UserName` varchar(16) NOT NULL default '',
  `UserType` varchar(16) NOT NULL default '',
  PRIMARY KEY  (`ID`),
  KEY `i_GroupID` (`GroupID`),
  KEY `i_OwnerUin` (`OwnerUin`),
  KEY `i_ParentUin` (`ParentUin`),
  KEY `i_UserName` (`UserName`),
  KEY `i_UserType` (`UserType`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户信息表';

