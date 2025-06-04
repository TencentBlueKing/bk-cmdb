### Description

Retrieve the list of service instances bound to a host based on the host ID.

### Parameters

| Name       | Type   | Required | Description                                                        |
|------------|--------|----------|--------------------------------------------------------------------|
| bk_biz_id  | int    | Yes      | Business ID                                                        |
| bk_host_id | int    | Yes      | Host ID to retrieve service instance information bound to the host |
| page       | object | No       | Query conditions                                                   |

#### page

| Name  | Type | Required | Description                             |
|-------|------|----------|-----------------------------------------|
| start | int  | Yes      | Record start position                   |
| limit | int  | Yes      | Number of records per page, maximum 500 |

### Request Example

```python
{
  "bk_biz_id": 1,
  "page": {
    "start": 0,
    "limit": 1
  },
  "bk_host_id": 26
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
    "count": 1,
    "info": [
       {
          "bk_biz_id": 1,
          "id": 1,
          "name": "test",
          "labels": {
              "test1": "1"
          },
          "service_template_id": 32,
          "bk_host_id": 26,
          "bk_module_id": 12,
          "creator": "admin",
          "modifier": "admin",
          "create_time": "2021-12-31T03:11:54.992Z",
          "last_time": "2021-12-31T03:11:54.992Z",
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
