### 功能描述

根据主机实例ID列表和想要获取的主机快照属性列表批量获取主机快照 (v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_ids  | int array  | 是     | 主机实例ID列表, 即bk_host_id列表，最多可填200个 |
| fields  |  string array   | 是     | 主机快照属性列表，控制返回结果的主机快照信息里有哪些字段<br>目前支持字段有：bk_host_id,bk_all_ips|

### 请求参数示例

```json
{
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
