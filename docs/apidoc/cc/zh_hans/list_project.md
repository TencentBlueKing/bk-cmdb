### 功能描述

查询项目(版本：v3.10.23+，权限：项目的查看权限)

### 请求参数

{{ common_args_desc }}


#### 接口参数
- 通用字段：

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| filter      | object  | 否   | 查询条件  |
| fields     | array  | 否     | 属性列表，控制返回结果里有哪些字段，能够加速接口请求和减少网络流量传输   |
| page       | object | 是     | 分页信息 |

#### filter 字段说明

属性字段过滤规则，用于根据属性字段搜索数据。该参数支持以下两种过滤规则类型，其中组合过滤规则可以嵌套，且最多嵌套2层。具体支持的过滤规则类型如下：

##### 组合过滤规则

由其它规则组合而成的过滤规则，组合的规则间支持逻辑与/或关系

| 字段        | 类型     | 必选  | 描述                              |
|-----------|--------|-----|---------------------------------|
| condition | string | 是   | 组合查询条件，支持 `AND` 和 `OR` 两种方式     |
| rules     | array  | 是   | 查询规则，可以是 `组合过滤规则` 或 `原子过滤规则` 类型 |

##### 原子过滤规则

基础的过滤规则，表示对某一个字段进行过滤的规则。任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则组合而成

| 名称       | 类型                            | 必选  | 说明                                                                                                |
|----------|-------------------------------|-----|---------------------------------------------------------------------------------------------------|
| field    | string                        | 是   | container的字段                                                                                      |
| operator | string                        | 是   | 操作符，可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between | 
| value    | 不同的field和operator对应不同的value格式 | 否   | 操作值                                                                                               |

组装规则可参考: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### page

| 字段  | 类型   | 必选 | 描述                 |
| ----- | ------ | ---- | -------------------- |
| start | int    | 是   | 记录开始位置         |
| limit | int    | 是   | 每页限制条数，最大500 |
| sort  | string | 否   | 排序字段             |
| enable_count |  bool  | 是  | 是否获取查询对象数量的标记。如果此标记为true那么表示此次请求是获取数量，此时其余字段必须为初始化值，start为0，limit为:0，sort为"" |

### 请求参数示例

### 获取详细信息请求参数
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "bk_status",
                "operator": "equal",
                "value": "enable"
            }
        ]
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "id",
        "enable_count": false
    }
}
```
### 获取数量请求示例
```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "bk_status",
                "operator": "equal",
                "value": "enable"
            }
        ]
    },
    "page": {
        "enable_count":true
    }
}
```
### 返回结果示例
### 详细信息接口响应
```json
{
    "result": true,
    "code": 0,
    "data": {
        "count": 0,
        "info": [
            {	
               "id": 1,
               "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
               "bk_project_name": "test",
               "bk_project_code": "test",
               "bk_project_desc": "test project",
               "bk_project_type": "mobile_game",
               "bk_project_sec_lvl": "public",
               "bk_project_owner": "admin",
               "bk_project_team": [1, 2],
               "bk_status": "enable",
               "bk_project_icon": "https://127.0.0.1/file/png/11111",
               "bk_supplier_account": "0",
               "create_time": "2022-12-22T11:22:17.504+08:00",
               "last_time": "2022-12-22T11:23:31.728+08:00"
            }
        ]
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 获取数量返回结果示例

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":1,
        "info":[
        ]
    },
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

#### data
| 字段  | 类型  | 描述         |
| ----- | ----- | ------------ |
| count | int   | 记录条数     |
| info  | array | 实际数据，仅返回fields里设置了的字段 |

#### data.info
| 字段                  | 类型  | 描述      |
|---------------------| ----- | --------- |
| id                  |  int     | 在cc中项目的唯一标识|
| bk_project_id       |  string  | 项目id|
| bk_project_name     |  string  | 项目名称|
| bk_project_code     |  string  | 项目英文名|
| bk_project_desc     |  string  | 项目描述|
| bk_project_type     |  enum  | 项目类型，可选值："mobile_game"(手游)、"pc_game"(端游)、"web_game"(页游)、"platform_prod"(平台产品)、"support_prod"(支撑产品)、"other"(其他)|
| bk_project_sec_lvl  | enum     | 保密级别，可选值："public"(公开)、"private"(私有)、"classified"(机密)|
| bk_project_owner    |  string  | 项目负责人|
| bk_project_team     | array    | 所属团队|
| bk_project_icon     |  string  | 项目图标     |
| bk_status           |  string  | 项目状态，可选值："enable"(启用)、"disabled"(未启用)|
| bk_supplier_account | string   | 开发商账号 |
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
