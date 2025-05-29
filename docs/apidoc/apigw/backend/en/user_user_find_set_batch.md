### Description

Batch obtain the property details of specified clusters under the specified business based on the business ID and the
list of cluster instance IDs, and the list of properties to be returned. (v3.8.6)

### Parameters

| Name      | Type  | Required | Description                                                                                       |
|-----------|-------|----------|---------------------------------------------------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                                                                       |
| bk_ids    | array | Yes      | List of cluster instance IDs, i.e., bk_set_id list, up to 500                                     |
| fields    | array | Yes      | List of cluster properties, control which fields are included in the returned cluster information |

### Request Example

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        11,
        12
    ],
    "fields": [
        "bk_set_id",
        "bk_set_name",
        "create_time"
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [
        {
            "bk_set_id": 12,
            "bk_set_name": "ss1",
            "create_time": "2020-05-15T22:15:51.67+08:00",
            "default": 0
        },
        {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "create_time": "2020-05-12T21:04:36.644+08:00",
            "default": 0
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | array  | Data returned by the request                                       |

#### data

| Name                 | Type   | Description                                                  |
|----------------------|--------|--------------------------------------------------------------|
| bk_set_name          | string | Cluster name                                                 |
| default              | int    | 0-normal cluster, 1-built-in module collection, default is 0 |
| bk_biz_id            | int    | Business ID                                                  |
| bk_capacity          | int    | Design capacity                                              |
| bk_parent_id         | int    | Parent node ID                                               |
| bk_set_id            | int    | Cluster ID                                                   |
| bk_service_status    | string | Service status: 1/2 (1: open, 2: closed)                     |
| bk_set_desc          | string | Cluster description                                          |
| bk_set_env           | string | Environment type: 1/2/3 (1: test, 2: experience, 3: formal)  |
| create_time          | string | Creation time                                                |
| last_time            | string | Update time                                                  |
| bk_supplier_account  | string | Developer account                                            |
| description          | string | Description information of the data                          |
| set_template_version | array  | Current version of the cluster template                      |
| set_template_id      | int    | Cluster template ID                                          |
| bk_created_at        | string | Creation time                                                |
| bk_updated_at        | string | Update time                                                  |

**Note: The returned values here only explain the system-built property fields. The rest of the returned values depend
on the user's own defined property fields.**
