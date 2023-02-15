### 功能描述

新建容器集群(v3.10.23+，权限:容器集群的创建权限)

### 请求参数

{{ common_args_desc }}


#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | 是     | 业务ID|
| name    |  string  | 是     | 集群名称|
| scheduling_engine |  string  | 否  | 调度引擎 |
| uid   |  string  | 是   | 集群自有ID|
| xid |  string  | 否   | 关联集群ID |
| version   |  string  | 否   | 集群版本 |
| network_type   |  string  | 否   | 网络类型 |
| region |  string  | 否    | 地域|
| vpc |  string  | 否    | vpc网络|
| network |  array  | 否    | 集群网络|
| type |  string  | 否     | 集群类型 |

### 请求参数示例

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "name":"cluster",
    "scheduling_engine":"k8s",
    "uid":"xxx",
    "xid":"xxx",
    "version":"1.1.0",
    "network_type":"underlay",
    "region":"xxx",
    "vpc":"xxx",
    "network":[
        "127.0.0.0/21"
    ],
    "type":"public-cluster"
}
```

### 返回结果示例

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "id":1
    },
    "request_id":"87de106ab55549bfbcc46e47ecf5bcc7"
}
```
### 返回结果参数说明

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| data    | object | 请求返回的数据      |
| request_id    | string | 请求链ID    |

### data

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| id  | int   |  创建的容器集群ID |
