### 权限说明

* 管理员具有所有权限
*  权限配置  只有管理员有权限

* 展示
	* 主机管理完全展示，具体的鉴权后端去判断
	* 自定义模型实例管理的权限按照后端api返回判断
	* 资源池管理权限按照后端api返回判断
	* 后台配置权限按照后端api返回判断

### 获取角色绑定权限
* API:
GET /api/{version}/topo/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}
* API名称： get_role_privilege
* 功能说明：
	* 中文：获取角色绑定权限
	* English ：get role bind privilege
* input body:
无
* input参数说明:

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 |supplier account code |
| bk_obj_id| string| 是|无| 模型ID |  object ID |
| bk_property_id| string| 是| 无|模型对应用户角色属性ID| object property id|


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":[
        "hostupdate",
        "hosttrans",
        "topoupdate",
        "customapi",
        "proconfig"
    ]
}
```
* output 参数说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string数组| 请求返回的数据 |return data|

无权限时返回为空数组

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| hostupdate | string| 主机编辑|host update|
| hosttrans  | string| 主机转移|host transfer|
| topoupdate | string| 主机拓扑编辑|business topo update|
| customapi | string| 自定义api|user custom api|
| proconfig | string| 进程管理|process config|




###  绑定角色权限
* API: POST /api/{version}/topo/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}
* API名称： bind_role_privilege
* 功能说明：
	* 中文：绑定角色绑定权限
	* English ： bind  user privilege
* input body:
```json
[
    "hostupdate",
    "hosttrans",
    "topoupdate",
    "customapi",
    "proconfig"
]
```

* input 参数说明：

输入为空数组则不绑定任何权限

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
|---|---|---|---|---|---|
| hostupdate | string| 否|无| 主机编辑|host update|
| hosttrans  | string| 否|无| 主机转移|host transfer|
| topoupdate | string| 否|无| 主机拓扑编辑|business topo update|
| customapi | string| 否|无| 自定义api|user custom api|
| proconfig | string| 否|无| 进程管理|process config|

* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":""
}
```

* output参数说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|


###  新建用户分组
* API: POST /api/{version}/topo/privilege/group/{bk_supplier_account}
* API名称： create_user_group
* 功能说明：
	* 中文：创建用户分组
	* English ： create  user group
* input body:
```json
{
    "group_name":"管理员",
    "user_list":"owen;tt"
}
```

* input参数说明:

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| group_name| string| 是|无|分组名 | group name|
| user_list|string | 是| 无|分组用户列表，多个用;分割| user list|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":""
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|

###  更新用户分组
* API: PUT  /api/{version}/topo/privilege/group/{bk_supplier_account}/{group_id}
* API名称： update_user_group
* 功能说明：
	* 中文：更新用户分组
	* English ： update  user group
* input body:
```json
{
    "group_name":"管理员",
    "user_list":"owen;tt"
}
```

* input参数说明:

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| group_id| string| 是|无|分组ID  | group ID|
| group_name| string| 否|无|分组名 | group name|
| user_list|string | 否| 无|分组用户列表，多个用;分割| user list|


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":""
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|


###  查询用户分组
* API: POST /api/{version}/topo/privilege/group/{bk_supplier_account}/search
* API名称： search_user_group
* 功能说明：
	* 中文：查询用户分组
	* English ： search  user group
* input body:

```json
{
    "group_name":"管理员",
    "user_list":"owen;tt"
}
```
* input参数说明:

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| group_name| string| 否|无|分组名 | group name|
| user_list|string | 否| 无|分组用户列表，多个用;分割| user list|

body 为空对象时返回所有的分组
* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":[
        {
            "group_name":"管理员",
            "user_list":"owen;tt",
            "group_id":1
        }
    ]
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object array| 请求返回的数据 |return data|

data object 说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| group_name| string| 分组名 |group name|
| user_list| string| 用户列表 |user list|
| group_id| string | 分组ID |user group ID|



###  删除用户分组
* API: DELETE  /api/{version}/topo/privilege/group/{bk_supplier_account}/{bk_group_id}
* API名称： delete_user_group
* 功能说明：
	* 中文：删除用户分组
	* English ： delete  user group
* input body:
无
* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| group_id| string| 是|无|分组ID  | group ID|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":""
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|

###  查询分组权限
* API :  GET  /api/{version}/topo/privilege/group/detail/{bk_supplier_account}/{group_id}
* API名称： search_group_privilege
* 功能说明：
	* 中文： 查询分组权限
	* English ： search  group privilege
* input body:
无
* input参数说明


| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---|---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| group_id| string| 是|无|分组ID  | group ID|

* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":
        {
            "group_id":1,
            "sys_config":{
                "global_busi":[
                    "resource"
                ],
                "back_config":[
                    "event",
                    "model",
                    "audit"
                ]
            },
            "model_config":{
                "network":{
                    "router":[
                        "update",
                        "delete"
                    ]
                }
            }
        }
    
}

```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| group_id| string|分组ID  | group ID|
| sys_config | object | 系统配置 |system config|
| back_config|object | 后台配置 |back config|
| model_config| object| 模型配置 |model config|


sys_config  目前仅有global_busi   字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| resource| string|主机资源池 | host resource pool|

back_config 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| event   | string| 事件推送配置| event push config|
| model   | string| 模型配置| model config|
| audit   | string| 审计配置| audit config|

model_config字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| create | string| 新增| create|
| update | string|  编辑| update|
| delete| string| 删除| delete|
| search| string| 查询| search|



###  查询用户权限
* API:  GET  /api/{version}/topo/privilege/user/detail/{bk_supplier_account}/{user_name}
* API名称： get_user_privilege
* 功能说明：
	* 中文： 获取用户权限
	* English ： get  user privilege
* input body:
无
* intput字段说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---|---| ---|
| bk_supplier_account| string| 是|无|开发商账号 | supplier account code|
| user_name| string| 是|无|用户名  | user name|


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":
        {
            "bk_group_id":1,
            "sys_config":{
                "global_busi":[
                    "resource"
                ],
                "back_config":[
                    "event",
                    "model",
                    "audit"
                ]
            },
            "model_config":{
                "network":{
                    "router":[
                        "update",
                        "delete"
                    ]
                }
            }
        }
    
}

```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| group_id| string|分组ID  | group ID|
| sys_config | object | 系统配置 |system config|
| back_config|object | 后台配置 |back config|
| model_config| object| 模型配置 |model config|


sys_config  目前仅有global_busi   字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| resource| string| 是|无|主机资源池 | host resource pool|

back_config 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| event   | string| 事件推送配置| event push config|
| model   | string| 模型配置| model config|
| audit   | string| 审计配置| audit config|

model config字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| create | string| 新增| create|
| update | string|  编辑| update|
| delete| string| 删除| delete|
| search| string| 查询| search|


###  更新分组权限
* API: POST  /api/{version}/topo/privilege/group/detail/{bk_supplier_account}/{group_id}
* API名称： update_group_privilege
* 功能说明：
	* 中文： 更新分组权限
	* English ： update  group privilege
* input body:

```json
{
    "sys_config":{
        "global_busi":[
            "resource"
        ],
        "back_config":[
            "event",
            "model",
            "audit"
        ]
    },
    "model_config":{
        "network":{
            "router":[
                "update",
                "delete"
            ]
        }
    }
}

```
* input字段说明
sys_config  目前仅有global_busi , 字段说明为：


| 名称  | 类型  |必填| 默认值 | 说明 |Description|
|---|---|---|---|---|---|
| resource| string|否|无|主机资源池 | host resource pool|


back_config 字段说明：

| 名称  | 类型  |必填| 默认值 | 说明 |Description|
|---|---|---|---|---|---|
| event   | string| 否|无| 事件推送配置| event push config|
| model   | string| 否|无| 模型配置| model config|
| audit   | string| 否|无| 审计配置| audit config|


model_config字段说明：

| 名称  | 类型  |必填| 默认值 | 说明 |Description|
|---|---|---|---|---|---|
| create | string| 否|无| 新增| create|
| update | string| 否|无|  编辑| update|
| delete| string| 否|无| 删除| delete|
| search| string| 否|无| 查询| search|


*  output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":""
}

```


* output字段说明


| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|
