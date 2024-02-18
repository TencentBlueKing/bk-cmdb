### Description

Host query without business information (Permission: Host pool host view permission)

### Parameters

| Name                 | Type   | Required | Description                                                                                                                                                                                   |
|----------------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id            | int    | No       | Business ID                                                                                                                                                                                   |
| page                 | object | Yes      | Query conditions                                                                                                                                                                              |
| host_property_filter | object | No       | Host property combination query conditions                                                                                                                                                    |
| fields               | array  | No       | List of host properties, controls which fields are returned in the results, speeding up the interface request and reducing network traffic. If not filled, all fields are returned by default |

#### host_property_filter

This parameter is a combination of host property field filtering rules, used to search for hosts based on host property
fields. The combination supports AND and OR two ways, can be nested, with a maximum nesting of 2 levels. The filtering
rule is a quadruple `field`, `operator`, `value`.

| Name      | Type   | Required | Description               |
|-----------|--------|----------|---------------------------|
| condition | string | No       | Combined query conditions |
| rules     | array  | No       | Rules                     |

#### rules

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can be referred
to: [QueryBuilder README](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Name  | Type | Required | Description                             |
|-------|------|----------|-----------------------------------------|
| start | int  | Yes      | Record start position                   |
| limit | int  | Yes      | Number of records per page, maximum 500 |

### Request Example

```json
{
    "page": {
        "start": 0,
        "limit": 3
    },
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "127.0.0.1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_os_type",
                        "operator": "not_in",
                        "value": [
                            "3"
                        ]
                    },
                    {
                        "field": "bk_cloud_id",
                        "operator": "equal",
                        "value": 0
                    }
                ]
            }
        ]
    }
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 30,
    "info": [
      {
        "bk_cloud_id": 0,
        "bk_host_id": 1,
        "bk_host_innerip": "192.168.15.18",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 2,
        "bk_host_innerip": "192.168.15.4",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 3,
        "bk_host_innerip": "192.168.15.12",
        "bk_mac": "",
        "bk_os_type": null
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description       |
|-------|-------|-------------------|
| count | int   | Number of records |
| info  | array | Actual host data  |

#### data.info

| Name            | Type   | Description             |
|-----------------|--------|-------------------------|
| bk_cloud_id     | int    | Cloud control area      |
| bk_host_id      | int    | Host ID                 |
| bk_host_innerip | string | Internal IP of the host |
| bk_mac          | string | Host MAC address        |
| bk_os_type      | string | Operating system type   |
