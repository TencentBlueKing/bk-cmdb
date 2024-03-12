### Function Description

Query Cluster

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type   | Required | Description                                                  |
| ------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_supplier_account | string | No       | Developer account                                            |
| bk_biz_id           | int    | Yes      | Business ID                                                  |
| fields              | array  | Yes      | Query fields, all fields are attributes defined in the set, including preset fields and user-defined fields |
| condition           | dict   | Yes      | Query condition, all fields are attributes defined in the set, including preset fields and user-defined fields |
| page                | dict   | Yes      | Paging condition                                             |

#### page

| Field | Type   | Required | Description           |
| ----- | ------ | -------- | --------------------- |
| start | int    | Yes      | Record start position |
| limit | int    | Yes      | Limit per page        |
| sort  | string | No       | Sorting field         |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 2,
    "fields": [
        "bk_set_name"
    ],
    "condition": {
        "bk_set_name": "test"
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_set_name"
    }
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 1,
        "info": [
            {
                "bk_set_name": "test",
                "default": 0
            }
        ]
    }
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure        |
| message    | string | Error message returned in case of request failure            |
| data       | object | Request returned data                                        |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |

#### data

| Field | Type  | Description                                                  |
| ----- | ----- | ------------------------------------------------------------ |
| count | int   | Number of data elements                                      |
| info  | array | Result set, where all fields are attributes defined in the set |

#### info

| Field                | Type   | Description                                                |
| -------------------- | ------ | ---------------------------------------------------------- |
| bk_set_name          | string | Cluster name                                               |
| default              | int    | 0-normal cluster, 1-built-in module set, default is 0      |
| bk_biz_id            | int    | Business ID                                                |
| bk_capacity          | int    | Design capacity                                            |
| bk_parent_id         | int    | Parent node ID                                             |
| bk_set_id            | int    | Cluster ID                                                 |
| bk_service_status    | string | Service status: 1/2(1: open, 2: closed)                    |
| bk_set_desc          | string | Cluster description                                        |
| bk_set_env           | string | Environment type: 1/2/3(1: test, 2: experience, 3: formal) |
| create_time          | string | Creation time                                              |
| last_time            | string | Update time                                                |
| bk_supplier_account  | string | Developer account                                          |
| description          | string | Description of the data                                    |
| set_template_version | array  | Current version of cluster template                        |
| set_template_id      | int    | Cluster template ID                                        |
| bk_created_at        | string | Creation time                                              |
| bk_updated_at        | string | Update time                                                |
| bk_created_by        | string | Creator                                                    |

**Note: The return value here only describes the system's built-in property fields. The rest of the return value depends on the user-defined property fields.**