### Functional description

create set

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | No     | supplier account code |
| bk_biz_id      | int     | Yes     | Business ID |
| data           | dict    | Yes     | data |

#### data

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_parent_id        |  int     | Yes     | the parent inst identifier |
| bk_set_name         |  string  | Yes     | set name |
| default             |  int     | No     | 0-normal set, 1-built-in set, default is 0 |
| set_template_id     |  int     | No     | set template ID, required when need to create set using set template |

**Note: Other optional fields are defined by the model**

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 1,
    "data": {
        "bk_parent_id": 1,
        "bk_set_name": "test-set",
        "bk_set_desc": "test-set",
        "bk_capacity": 1000,
        "description": "description",
        "set_template_id": 1
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
        "bk_biz_id": 11,
        "bk_capacity": 1000,
        "bk_parent_id": 11,
        "bk_service_status": "1",
        "bk_set_desc": "test-set",
        "bk_set_env": "3",
        "bk_set_id": 4780,
        "bk_set_name": "test-set",
        "bk_supplier_account": "0",
        "create_time": "2022-02-22T20:34:01.386+08:00",
        "default": 0,
        "description": "description",
        "last_time": "2022-02-22T20:34:01.386+08:00",
        "set_template_id": 11,
        "set_template_version": null
     }
}
```

### Return Result Parameters Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | The success or failure of the request. true: the request was successful; false: the request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| data | object | The data returned by the request.
| permission | object | Permission information |
| request_id | string | Request chain id |

#### data

| field | type | description |
| -----------|-----------|--------------|
| bk_biz_id | int | business_id |
| bk_capacity | int | design_capacity |
|bk_parent_id | int | ID of the parent node |
| bk_set_id | int | cluster_id |
| bk_service_status | string | Service status:1/2(1:open,2:closed) |
|bk_set_desc|string|cluster_description|
| bk_set_env | string | Environment type:1/2/3(1:test,2:experience,3:official) |
|bk_set_name|string|cluster_name|
| create_time | string | creation time |
| last_time | string | update_time |
| bk_supplier_account | string | developer_account |
| default | int | 0-general cluster, 1-built-in module set, default is 0 |
| description | string | Description information for the data |
| set_template_version | array | The current version of the cluster template |
| set_template_id| int | Cluster template ID |
