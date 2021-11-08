### Functional description

find host snapshot in batch (v3.8.6)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_ids  | int array  | Yes     | bk_host_id array，the max length is 200 |
| fields  |  string array   | Yes     | host snapshot property list, the specified snapshot property feilds will be returned <br>supported fields: bk_host_id,bk_all_ips|

### Request Parameters Example

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

### Return Result Example

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
