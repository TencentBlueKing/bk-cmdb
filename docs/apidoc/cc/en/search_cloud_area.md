### Function Description

Query control area (Permission: control area view permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description        |
| --------- | ------ | -------- | ------------------ |
| condition | object | No       | Query conditions   |
| page      | object | Yes      | Paging information |

#### condition

| Field         | Type   | Required | Description       |
| ------------- | ------ | -------- | ----------------- |
| bk_cloud_id   | int    | No       | Control area ID   |
| bk_cloud_name | string | No       | Control area name |

#### page Field Description

| Field | Type | Required | Description                                      |
| ----- | ---- | -------- | ------------------------------------------------ |
| start | int  | No       | Data offset position                             |
| limit | int  | Yes      | Number of data restrictions, recommended for 200 |

### Request Parameter Example

```python
{

    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_cloud_id": 12,
        "bk_cloud_name": "aws"
    },
    "page":{
        "start":0,
        "limit":10
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
    "count": 10,
    "info": [
         {
            "bk_cloud_id": 0,
            "bk_cloud_name": "aws",
            "bk_supplier_account": "0",
            "create_time": "2019-05-20T14:59:48.354+08:00",
            "last_time": "2019-05-20T14:59:48.354+08:00"
        },
        .....
    ]
   
  }
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |

#### data

| Field  | Type  | Description                             |
| ----- | ----- | --------------------------------------- |
| count | int   | Number of records                       |
| info  | array | Information about the control area list |

#### data.info Field Description：

| Field          | Type   | Description            |
| ------------- | ------ | ---------------------- |
| bk_cloud_id   | int    | Control area ID        |
| bk_cloud_name | string | Control area name      |
| create_time   | string | Creation time          |
| last_time     | string | Last modification time |