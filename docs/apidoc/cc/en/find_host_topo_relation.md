### Function Description

Get the relationship between hosts and topology (Permission: Business access permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field         | Type   | Required | Description                    |
| ------------- | ------ | -------- | ------------------------------ |
| bk_biz_id     | int    | Yes      | Business ID                    |
| bk_set_ids    | array  | No       | List of cluster IDs, up to 200 |
| bk_module_ids | array  | No       | List of module IDs, up to 500  |
| bk_host_ids   | array  | No       | List of host IDs, up to 500    |
| page          | object | Yes      | Page information               |

#### page Field Description

| Field | Type | Required | Description                                          |
| ----- | ---- | -------- | ---------------------------------------------------- |
| start | int  | No       | Data offset position                                 |
| limit | int  | Yes      | Number of records per page, recommended value is 200 |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "page":{
        "start":0,
        "limit":10
    },
    "bk_biz_id":2,
    "bk_set_ids": [1, 2],
    "bk_module_ids": [23, 24],
    "bk_host_ids": [25, 26]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "data": {
        "count": 2,
        "data": [
            {
                "bk_biz_id": 2,
                "bk_host_id": 2,
                "bk_module_id": 2,
                "bk_set_id": 2,
            },
            {
                "bk_biz_id": 1,
                "bk_host_id": 1,
                "bk_module_id": 1,
                "bk_set_id": 1,
            }
        ],
        "page": {
            "limit": 10,
            "start": 0
        }
    },
    "message": "success",
    "permission": null,
    "request_id": "f5a6331d4bc2433587a63390c76ba7bf"
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |

#### data Field Explanation:

| Field  | Type         | Description                                                  |
| ----- | ------------ | ------------------------------------------------------------ |
| count | int          | Number of records                                            |
| data  | object array | Details list of data for hosts and clusters, modules, clusters under the business |
| page  | object       | Page                                                         |

#### data.data Field Explanation:

| Field                | Type   | Description      |
| ------------------- | ------ | ---------------- |
| bk_biz_id           | int    | Business ID      |
| bk_set_id           | int    | Cluster ID       |
| bk_module_id        | int    | Module ID        |
| bk_host_id          | int    | Host ID          |

#### data.page Field Explanation:

| Field  | Type | Description                |
| ----- | ---- | -------------------------- |
| start | int  | Data offset position       |
| limit | int  | Number of records per page |