### 功能描述

跨业务转移主机，只能将源业务空闲机池集群中的主机转移到目标业务的空闲机池集群(v3.10.27+，权限：主机转移到其他业务)

### 请求参数

#### 通用参数

{{ common_args_desc }}


#### 接口参数

| 字段           | 类型  | 必选  | 描述                                                     |
| ------------- | ----- | ---- | ------------------------------------------------------- |
| src_bk_biz_id | int   | 是   | 要转移的主机所属的业务ID                                    |
| bk_host_id    | array | 是   | 要转移的主机id列表，最大长度为500                            |
| dst_bk_biz_id | int   | 是   | 主机要转移到的业务ID                                       |
| bk_module_id  | int   | 是   | 主机要转移到的模块ID，该模块ID必须为下空闲机池set下的模块ID。 |

### 请求参数示例

```json
{
    "bk_app_code": "xxx",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "src_bk_biz_id": 2,
    "dst_bk_biz_id": 3,
    "bk_host_id": [
        9,
        10
    ],
    "bk_module_id": 10
}
```

### 返回结果示例

```json
{
    "result":true,
    "code":0,
    "data":null,
    "message":"success",
    "permission":null,
    "request_id":"87de106ab55549bfbcc46e47ecf5bcc7"
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
