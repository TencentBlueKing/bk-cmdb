### Functional description

search inst by object

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account |  string  | No     | supplier account code |
| bk_obj_id           |  string  | Yes     | the object id |
| fields              |  array   | No     | need to show |
| condition           |  dict    | No     | search condition |
| page                |  dict    | No     | page condition |

#### page

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| start    |  int    | Yes     | start record |
| limit    |  int    | Yes     | page limit, max is 200 |
| sort     |  string | No     | the field for sort |

#### fields parameter description

The parameter is all the fields defined by the model corresponding to the target instance of the query


#### condition parameter description

The condition parameter is all the fields defined by the model corresponding to the target instance of the query

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_obj_id": "xxx",
    "fields": [
    ],
    "condition": {
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": ""
    }
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 4,
        "info": [
            {
                "bk_inst_id": 0,
                "bk_inst_name": "default area",
                "bk_supplier_account": "123456789"
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
| count     | int       | the inst count |
| info      | array     | the set of instances of the model being queried |

#### info

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| bk_inst_id         | int       | inst area ID |
| bk_inst_name       | string    | the inst name |
| bk_supplier_account | string    | supplier account code |
