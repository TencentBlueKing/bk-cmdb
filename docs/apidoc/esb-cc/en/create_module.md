### Functional description

create module

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | No     | supplier account code |
| bk_biz_id      | int     | Yes     | Business ID |
| bk_set_id      | int     | Yes     | the set id |
| data           | dict    | Yes     | Data |

#### data

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_parent_id      | int     | Yes     | the parent inst id |
| bk_module_name    | string  | Yes     | Module name |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 1,
    "bk_set_id": 10,
    "data": {
        "bk_parent_id": 10,
        "bk_module_name": "test"
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
    "data": {
        "bk_bak_operator": null,
        "bk_biz_id": 1,
        "bk_module_id": 37825,
        "bk_module_name": "test",
        "bk_module_type": "1",
        "bk_parent_id": 10,
        "bk_set_id": 10,
        "bk_supplier_account": "0",
        "create_time": "2022-02-22T20:25:19.049+08:00",
        "default": 0,
        "host_apply_enabled": false,
        "last_time": "2022-02-22T20:25:19.049+08:00",
        "operator": null,
        "service_category_id": 2,
        "service_template_id": 0,
        "set_template_id": 0
    }
}
```
### Return Result Parameters Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| data | object | The data returned by the request.
| permission | object | Permission information |
| request_id | string | Request chain id |

#### data

| field | type | description |
| -----------|-----------|--------------|
| bk_bak_operator | string | backup maintainer |
| bk_module_id | int | model id |
|bk_biz_id | int | business id|
| bk_module_id | int | module id |
| bk_module_name | string | module name |
|bk_module_type|string|module type|
|bk_parent_id|int|the id of the parent|
| bk_set_id | int | set id |
| bk_supplier_account | string | developer_account |
| create_time | string | creation_time |
| last_time | string | update_time |
|default | int | Indicates the module type |
| host_apply_enabled |bool | whether to enable automatic application of host attributes |
| operator | string | primary maintainer |
|service_category_id|int|service category id|
|service_template_id|int|service template id|
| set_template_id | int | set template id |
