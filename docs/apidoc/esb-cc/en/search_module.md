### Functional description

search module

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | No     | supplier account code |
| bk_biz_id      |  int     | Yes     | the business id |
| bk_set_id      |  int     | No     | Set ID |
| fields         |  array   | Yes     | search fields |
| condition      |  dict    | Yes     | search condition |
| page           |  dict    | Yes     | page condition |

#### page

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| start    |  int    | Yes     | start record |
| limit    |  int    | Yes     | page limit |
| sort     |  string | No     | the field for sort |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 2,
    "fields": [
        "bk_module_name",
        "bk_set_id"
    ],
    "condition": {
        "bk_module_name": "test"
    },
    "page": {
        "start": 0,
        "limit": 10
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 2,
        "info": [
            {
                "bk_module_name": "test",
                "bk_set_id": 11,
                "default": 0
            },
            {
                "bk_module_name": "test",
                "bk_set_id": 12,
                "default": 0
            }
        ]
    }
}
```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |
| permission    | object | permission Information    |
| request_id    | string | request chain id    |

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     | int       | the data item count |
| info      | array     | the data result array |

#### info

| Field | Type | Description |
|-----------|-----------|-----------|
| bk_module_name | string | module name |
| bk_set_id | int | set id |
|default | int | indicates the module type |