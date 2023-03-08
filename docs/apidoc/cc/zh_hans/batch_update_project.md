### 功能描述

更新项目(版本：v3.10.23+，权限：项目的更新权限)

### 请求参数

{{ common_args_desc }}


#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| ids | array| 是 |在cc中的id唯一标识数组,一次限制最大传200个|
| data |  object | 是 | 包含需要更新的字段|

#### data

| 字段                 |  类型      | 必选   | 描述                                                                                                            |
|--------------------|------------|--------|---------------------------------------------------------------------------------------------------------------|
| bk_project_name    |  string  | 否     | 项目名称                                                                                                          |
| bk_project_desc    |  string  | 否     | 项目描述                                                                                                          |
| bk_project_type    |  enum  | 否     | 项目类型，可选值："mobile_game"(手游)、"pc_game"(端游)、"web_game"(页游)、"platform_prod"(平台产品)、"support_prod"(支撑产品)、"other"(其他) |
| bk_project_sec_lvl | enum   | 否     | 保密级别，可选值："public"(公开)、"private"(私有)、"classified"(机密)                                                          |
| bk_project_owner   |  string  | 否     | 项目负责人                                                                                                         |
| bk_project_team    | array   |  否    | 所属团队                                                                                                          |
| bk_project_icon    |  string  | 否     | 项目图标     |
| bk_status          |  string  | 否     | 项目状态，可选值："enable"(启用)、"disabled"(未启用)                                                                         |


### 请求参数示例

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "ids":[
        1, 2, 3
    ],   
    "data": {
        "bk_project_name": "test",
        "bk_project_desc": "test project",
        "bk_project_type": "mobile_game",
        "bk_project_sec_lvl": "public",
        "bk_project_owner": "admin",
        "bk_project_team": [1, 2],
        "bk_status": "enable",
        "bk_project_icon": "https://127.0.0.1/file/png/11111"
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |
