### 功能描述

更新容器集群属性字段(v3.10.23+，权限:容器集群的编辑权限)
### 请求参数

{{ common_args_desc }}

#### 接口参数
| 字段                |  类型      | 必选   |  描述      |
|---------------------|------------|--------|------------|
| bk_biz_id    |  int  | 是     | 业务ID|
| ids           | array        | 否     | cluster在cmdb中的唯一ID列表|
| data         | object     | 是     | 需要更新的数据 |

#### data

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| scheduling_engine |  string  | 否  | 调度引擎 |
| xid |  string  | 否   | 关联集群ID |
| version   |  string  | 否   | 集群版本 |
| network_type   |  string  | 否   | 网络类型 |
| region |  string  | 否    | 地域|
| vpc |  string  | 否    | vpc网络|
| network |  array  | 否    | 集群网络|
| type |  string  | 否     | 集群类型 |

**注意：**
- 一次性更新集群数量不超过100个。

### 请求参数示例

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "ids":[
        1
    ],
    "data":{
        "scheduling_engine":"engine1",
        "version":"1.20.6",
        "network_type":"underlay",
        "region":"xxx",
        "vpc":"xxx",
        "network":"127.0.0.0/21",
        "type":"public-cluster"
    }
}
```

### 返回结果示例

```json
 {
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
    "data": null
}
```
### 返回结果参数说明

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 无数据返回                          |
