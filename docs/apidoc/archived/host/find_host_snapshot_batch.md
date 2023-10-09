### 功能描述

根据主机实例ID列表和想要获取的主机快照属性列表批量获取主机快照 (v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_ids  |  array  | 是     | 主机实例ID列表, 即bk_host_id列表，最多可填200个 |
| fields  |   array   | 是     | 主机快照属性列表，控制返回结果的主机快照信息里有哪些字段<br>目前支持字段有：bk_host_id,bk_all_ips|

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_ids": [
        1,
        2
    ],
    "fields": [
        "bk_host_id",
        "bk_all_ips"
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_all_ips": {
                "interface": [
                    {
                        "addrs": [
                            {
                                "ip": "192.xx.xx.xx"
                            },
                            {
                                "ip": "fe80::xx:xx:xx:xx"
                            }
                        ],
                        "mac": "52:xx:xx:xx:xx:xx"
                    },
                    {
                        "addrs": [
                            {
                                "ip": "192.xx.xx.xx"
                            }
                        ],
                        "mac": "02:xx:xx:xx:xx:xx"
                    }
                ]
            },
            "bk_host_id": 1
        },
        {
            "bk_all_ips": {
                "interface": [
                    {
                        "addrs": [
                            {
                                "ip": "172.xx.xx.xx"
                            },
                            {
                                "ip": "fe80::xx:xx:xx:xx"
                            }
                        ],
                        "mac": "52:xx:xx:xx:xx:xx"
                    },
                    {
                        "addrs": [
                            {
                                "ip": "192.xx.xx.xx"
                            }
                        ],
                        "mac": "02:xx:xx:xx:xx:xx"
                    }
                ]
            },
            "bk_host_id": 2
        }
    ]
}
```
### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |
