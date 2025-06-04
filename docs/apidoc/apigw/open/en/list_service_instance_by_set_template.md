### Description

Retrieve the list of service instances based on the cluster template ID.

### Parameters

| Name            | Type   | Required | Description           |
|-----------------|--------|----------|-----------------------|
| bk_biz_id       | int    | Yes      | Business ID           |
| set_template_id | int    | Yes      | Cluster template ID   |
| page            | object | Yes      | Pagination parameters |

#### page

| Name  | Type | Required | Description                                                        |
|-------|------|----------|--------------------------------------------------------------------|
| start | int  | Yes      | Record start position                                              |
| limit | int  | Yes      | Number of records per page, maximum 500, recommended to set to 200 |

### Request Example

```python
{
  "bk_biz_id": 1,
  "set_template_id": 1,
  "page": {
    "start": 0,
    "limit": 10
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
    "data": {
        "count": 2,
        "info": [
            {
                "bk_biz_id": 3,
                "id": 1,
                "name": "10.0.0.1_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 1,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T02:46:25.002Z",
                "last_time": "2020-10-09T02:46:25.002Z",
                "bk_supplier_account": "0"
            },
            {
                "bk_biz_id": 3,
                "id": 3,
                "name": "127.0.122.2_lgh-process-1",
                "labels": null,
                "service_template_id": 50,
                "bk_host_id": 3,
                "bk_module_id": 59,
                "creator": "admin",
                "modifier": "admin",
                "create_time": "2020-10-09T03:04:19.859Z",
                "last_time": "2020-10-09T03:04:19.859Z",
                "bk_supplier_account": "0"
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data Field Explanation

| Name  | Type  | Description              |
|-------|-------|--------------------------|
| count | int   | Total number of records  |
| info  | array | List of returned results |

#### info Field Explanation

| Name                | Type   | Description                |
|---------------------|--------|----------------------------|
| id                  | int    | Service instance ID        |
| name                | string | Service instance name      |
| bk_biz_id           | int    | Business ID                |
| bk_module_id        | int    | Module ID                  |
| bk_host_id          | int    | Host ID                    |
| creator             | string | Creator of this data       |
| modifier            | string | Last modifier of this data |
| create_time         | string | Creation time              |
| last_time           | string | Update time                |
| bk_supplier_account | string | Supplier account           |
| service_template_id | int    | Service template ID        |
| labels              | map    | Label information          |
