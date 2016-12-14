#蓝鲸智云配置平台接口说明	

##环境说明
* 环境配置详见【README.md】
##输入参数
* 接口都采用post form-data的方式传递参数
##输出参数
* 1、输出格式统一为json
* 2、正确返回json的形式一般是这样的：
```json
{
    code: 0
}
```
* 3、错误返回json的形式一般是这样的：
```json
{
    "code": "0004",
    "msg": "ApplicationID Illegal",
    "extmsg": "非法的业务Id"
}
```
##4、接口说明

###查询业务列表
* 接口地址：/api/App/getapplist
* 请求方式：POST
* 参数列表：
 无
* 参数说明：
* 返回数据格式：
```json
{
  "code": 0, 
  "data": [
    {
      "ApplicationID": "996", 
      "ApplicationName": "资源池", 
      "Creator": "示例公司", 
      "CreateTime": "2016-03-07 15:14:37", 
      "Default": "1", 
      "DeptName": "", 
      "Description": ""
      .....
    }
  ]
}
```



###IP查询主机
* 接口地址：/api/Host/gethostlistbyip
* 请求方式：POST
* 参数列表：
 ApplicationID	业务ID	必选
 IP	主机IP(内网IP或外网IP)	必选

* 返回数据格式：
```json
{
    "code": 0,
    "data": [
        {
            "ApplicationID": "997",
            "SetID": "1187",
            "ModuleID": "2105",
            "HostID": "1",
            "AssetID": "pc-8caade00023acfa122a5f36aee26ae98",
            "HostName": "a",
            "DeviceClass": "",
            "Operator": "",
            "BakOperator": "",
            "InnerIP": "10.126.84.10",
            "OuterIP": "10.126.84.10",
            "Status": "",
            "CreateTime": "1970-01-01 00:00:00",
            "HardMemo": "",
            "Region": "",
            "OSName": "",
            "IdcName": "",
            "ApplicationName": "bbb",
            "SetName": "",
            "ModuleName": "空闲机"
        }
    ]
}
```
###分布模块ID查询主机
* 接口地址：/api/Host/getmodulehostlist
* 请求方式：POST
* 参数列表：
 ApplicationID	业务ID	必选
 ModuleID	分布模块ID	必选
* 返回数据格式：
```json
{
    "code": 0,
    "data": [
        {
            "ApplicationID": "997",
            "SetID": "1187",
            "ModuleID": "2105",
            "HostID": "2",
            "AssetID": "pc-e23675fcbe954635f8b860a172b423c3",
            "HostName": "b",
            "DeviceClass": "",
            "Operator": "",
            "BakOperator": "",
            "InnerIP": "10.126.84.11",
            "OuterIP": "10.126.84.11",
            "Status": "",
            "CreateTime": "1970-01-01 00:00:00",
            "Mem": "0",
            "HardMemo": "",
            "Source": "3",
            "OSName": "",
            "IdcName": "",
            "Region": "",
            "ApplicationName": "bbb",
            "SetName": "",
            "ModuleName": "空闲机"
        }
    ]
}
```
###分布SetID查询主机
* 接口地址：/api/Host/getsethostlist
* 请求方式：POST
* 参数列表：
 ApplicationID	业务ID	必选
 SetID	分布SetID	必选
* 返回数据格式：
```json
{
    "code": 0,
    "data": [
        {
            "ApplicationID": "997",
            "SetID": "1187",
            "ModuleID": "2105",
            "HostID": "2",
            "AssetID": "pc-e23675fcbe954635f8b860a172b423c3",
            "HostName": "b",
            "DeviceClass": "",
            "Operator": "",
            "BakOperator": "",
            "InnerIP": "10.126.84.11",
            "OuterIP": "10.126.84.11",
            "Status": "",
            "CreateTime": "1970-01-01 00:00:00",
            "HardMemo": "",
            "Region": "",
            "OSName": "",
            "IdcName": "",
            "ApplicationName": "bbb",
            "SetName": "",
            "ModuleName": "空闲机"
        }
    ]
}
```
###业务ID查询主机
* 接口地址：/api/Host/getapphostlist
* 请求方式：POST
* 参数列表：
ApplicationID	业务ID	必选

* 返回数据格式：
```json
{
    "code": 0,
    "data": [
        {
            "ApplicationID": "997",
            "SetID": "1187",
            "ModuleID": "2105",
            "HostID": "2",
            "AssetID": "pc-e23675fcbe954635f8b860a172b423c3",
            "HostName": "b",
            "DeviceClass": "",
            "Operator": "",
            "BakOperator": "",
            "InnerIP": "10.126.84.11",
            "OuterIP": "10.126.84.11",
            "Status": "",
            "CreateTime": "1970-01-01 00:00:00",
            "HardMemo": "",
            "Region": "",
            "OSName": "",
            "IdcName": "",
            "ApplicationName": "bbb",
            "SetName": "",
            "ModuleName": "空闲机"
        }
    ]
}
```
###根据appid查询业务的分布拓扑树
* 请求地址：/api/TopSetModule/getappsetmoduletreebyappid 
* 输入参数说明：
 ApplicationID	业务ID	必选
* 响应信息：
成功：
```json
{
    "code": 0,
    "data": {
        "ApplicationID": "997",
        "ApplicationName": "bbb",
        "Creator": "admin",
        "CreateTime": "2016-03-07 15:43:28",
        "Default": "0",
        "DeptName": "示例公司",
        "Description": "",
        "Display": "1",
        "GroupName": "",
        "LifeCycle": "公测",
        "Maintainers": "owenlxu;admin",
        "LastTime": "2016-03-08 05:09:42",
        "Level": "2",
        "Owner": "示例公司",
        "ProductPm": "admin",
        "Type": "0",
        "Source": "",
        "CompanyID": "0",
        "BusinessDeptName": "",
        "Children": [
            {
                "SetID": "1187",
                "ApplicationID": "997",
                "Default": "1",
                "Capacity": "0",
                "CreateTime": "1970-01-01 00:00:00",
                "ChnName": "",
                "Description": "",
                "EnviType": "",
                "LastTime": "2016-03-07 15:43:28",
                "ParentID": "0",
                "SetName": "空闲机池",
                "ServiceStatus": "",
                "Openstatus": "",
                "Children": [
                    {
                        "ModuleID": "2105",
                        "ApplicationID": "997",
                        "BakOperator": "",
                        "CreateTime": "1970-01-01 00:00:00",
                        "Default": "1",
                        "Description": "",
                        "LastTime": "2016-03-07 15:43:28",
                        "ModuleName": "空闲机",
                        "Operator": "",
                        "SetID": "1187",
                        "HostNum": 0
                    },
                    {
                        "ModuleID": "2119",
                        "ApplicationID": "997",
                        "BakOperator": "admin",
                        "CreateTime": "2016-03-08 18:36:48",
                        "Default": "0",
                        "Description": "",
                        "LastTime": "2016-03-08 18:36:48",
                        "ModuleName": "aaaa",
                        "Operator": "admin",
                        "SetID": "1187",
                        "HostNum": 0
                    }
                ]
            }
        ]
    }
}
```
###查询业务下的所有模块
* 请求地址：/api/Module/getmodules
* 输入参数说明：
* ApplicationID	业务ID	必选
* 响应信息：
成功：
```json
{
    "code": 0,
    "data": [
        "空闲机",
        "aaaa"
    ]
}
```
###查根据set属性查询主机
* 请求地址：/api/set/gethostsbyproperty
* 输入参数说明：
ApplicationID	业务ID	必选
SetID	大区ID，选填传多个，分割	选填
SetEnviType Set环境类型，选填，传多个,分割	选填
SetServiceStatus Set开放状态，选填，传多个,分割	选填
ModuleName	模块名，选填，传多个,分割	选填
* 响应信息：
成功：
```json
{
    "code": 0,
    "data": [
        {
            "InnerIP": "10.126.84.11",
            "OuterIP": "10.126.84.11",
            "Source": "3",
            "HostID": "2"
        }
    ]
}
```
###查根据set属性查询模块
* 请求地址：/api/Set/getmodulesbyproperty
* 输入参数说明：
ApplicationID	业务ID	必选
SetID	大区ID，选填传多个,分割	选填
SetEnviType Set环境类型，选填，传多个,分割	选填
SetServiceStatus Set开放状态，选填，传多个,分割	选填
* 响应信息：
成功：
```json
{
    "code": 0,
    "data": [
        "空闲机",
        "aaaa"
    ]
}
```
###获取所有set属性
* 请求地址：/api/Set/getsetproperty
* 输入参数说明：
无
* 响应信息：
```json
{
    "code": 0,
    "data": {
        "SetEnviType": [
            {
                "Property": "2",
                "value": "开放4"
            },
            {
                "Property": "3",
                "value": "开放4"
            }
        ],
        "SetServiceStatus": [
            {
                "Property": "0",
                "value": "开放4"
            },
            {
                "Property": "1",
                "value": "开放5"
            },
            {
                "Property": "5",
                "value": "开放67"
            },
            {
                "Property": "6",
                "value": "开放779"
            }
        ]
    }
}
```

###根据set属性获取set
* 请求地址：/api/Set/getsetsbyproperty
* 输入参数说明：
	ApplicationID	业务ID	必选
	SetEnviType	Set类型	非必选
	SetServiceStatus	Set服务状态	非必选
* 响应信息：
```json
{
    "code": 0,
    "data": [
        {
            "SetID": "1187",
            "SetName": "空闲机池"
        }
    ]
}
```

###获取userName有权限的业务
* 请求地址：/api/App/getappbyuin
* 输入参数说明：
	userName	API账号	必选
* 响应信息：
 请求数据成功：
```json
{
    "code": 0,
    "data": [
        {
            "ApplicationID": "996",
            "ApplicationName": "资源池",
            "Creator": "示例公司",
            "CreateTime": "2016-03-07 15:14:37",
            "Default": "1",
            "DeptName": "",
            "Description": "",
            "Display": "1",
            "GroupName": "",
            "LifeCycle": "",
            "Maintainers": "示例公司",
            "LastTime": "2016-03-07 15:14:37",
            "Level": "2",
            "Owner": "示例公司",
            "ProductPm": "",
            "Type": "0",
            "Source": "0",
            "CompanyID": "0",
            "BusinessDeptName": ""
        }
    ]
}
```
###新增加主机
* 请求地址：/api/Host/addHost
* 输入参数说明：
	InnerIP	内网IP	必选
	OuterIP	外网IP	选填
	HostName	主机名	选填
	Operator	操作者	选填
	BakOperator	备份操作者	选填
* 响应信息：

请求数据成功：
```json
 {
  "code": 0, 
  "data": [ ]
 }
```
###删除主机
* 请求地址：/api/Host/delHost
* 输入参数说明：
InnerIP	内网IP	必选
* 响应信息：

 请求数据成功：
```json
 {
  "code": 0, 
  "data": [ ]
 }
```

###编辑主机
* 请求地址：/api/Host/editHost
* 输入参数说明：
InnerIP	内网IP	必选
OuterIP	外网IP	选填
HostName	主机名	选填
Operator	操作者	选填
BakOperator	备份操作者	选填
* 响应信息：

 请求数据成功：
 
```json
 {
  "code": 0, 
  "data": [ ]
 }
```
###新增加模块
* 请求地址：/api/Module/addModule
* 输入参数说明：
AppName	业务名	必选
SetName	集群名	必选
ModuleName	模块名	必选
Operator	操作者	必选
BakOperator	备份操作者	必选
* 响应信息：

 请求数据成功：
```json
 {
  "code": 0, 
  "data": [ ]
 }
```
###删除模块
* 请求地址：/api/Module/delModule
* 输入参数说明：
AppName	业务名	必选
SetName	集群名	必选
ModuleName	模块名	必选
* 响应信息：

 请求数据成功：
```json
 {
  "code": 0, 
  "data": [ ]
 }
```
###编辑模块
* 请求地址：/api/Module/delModule
* 输入参数说明：
AppName	业务名	必选
SetName	集群名	必选
ModuleName	现模块名	必选
newModuleName	新的模块名	选填
Operator	维护人	选填
BakOperator	备份维护人	选填
* 响应信息：

 请求数据成功：
  ```json
 {
  "code": 0, 
  "data": [ ]
}
```


五、错误信息列表

code	msg	extMsg
* 0001
ApplicationID Illegal
非法的业务id

* 0002
no app
没有业务

* 0004
user name illegal
用户id非法

* 0005
only default app exsit
只有默认业务

* 0006
only right to app
没权利访问业务

* 0011
no host
没有主机

* 0016
invalid input ApplicationID
业务id输入无效


