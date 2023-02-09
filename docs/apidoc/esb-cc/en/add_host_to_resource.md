### Functional description

add host to resource

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account |  string     | No     | supplier account code |
| host_info      |  dict    | Yes     | host info |
| bk_biz_id      |  int     | No     | Business ID   |

#### host_info

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_host_innerip |  string   | Yes     | host inner ip |
| import_from     |  string   | Yes     | host import source |
| bk_cloud_id     |  int      | Yes     | Cloud area ID |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "host_info": {
        "0": {
            "bk_host_innerip": "10.0.0.1",
            "bk_cloud_id": 0,
            "import_from": "3"
        }
    }
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {}
}
```

### Return Result Parameters Description

#### response

| Field     | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |
| permission    | object | permission Information    |
| request_id    | string | request chain id    |
