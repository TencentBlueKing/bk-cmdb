### Function description

Add hosts to the service idle machine
- This interface ensures that hosts are either added successfully or fail at the same time(v3.10.25+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| field | type | mandatory | description |
| -----------|------------|--------|------------|
| bk_host_list | array | Yes | Host information (array length is limited to 200 at a time) |
| bk_biz_id | int | yes | business_id |

#### bk_host_list(host-related fields)

| field | type | required | description |
| -----------|------------|--------|------------|
| bk_host_innerip | string | yes | host_internal_ip |
| bk_cloud_id | int | Yes | cloud_region_id |
| bk_addressing | string | Yes | Addressing method, "static", "dynamic" |
| operator | string | No | Primary maintainer | 
...

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_host_list": [
        {
            "bk_host_innerip": "10.0.0.1",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        },
        {
            "bk_host_innerip": "10.0.0.2",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        }
    ]
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
        "bk_host_ids": [
            1,
            2
        ]
    }
}
```
### Return Result Parameters Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.|
| code | int | The error code. 0 means success, >0 means failure error.|
| message | string | The error message returned by the failed request.|
| data | object | The data returned by the request.|
| permission | object | Permission information |
| request_id | string | Request chain id |

#### data
| field | type | description |
| -----------|-----------|--------------|
| bk_host_ids | array | host_id of the host |
