### Description

Query cluster templates based on business ID.

### Parameters

| Name             | Type   | Required | Description                   |
|------------------|--------|----------|-------------------------------|
| bk_biz_id        | int    | Yes      | Business ID                   |
| set_template_ids | array  | No       | Array of cluster template IDs |
| page             | object | No       | Pagination information        |

#### Explanation of the page field

| Name  | Type   | Required | Description                                   |
|-------|--------|----------|-----------------------------------------------|
| start | int    | No       | Starting position of the record               |
| limit | int    | No       | Number of records per page, maximum 1000      |
| sort  | string | No       | Sorting field, '-' indicates descending order |

### Request Example

```json
{
  "bk_biz_id": 10,
  "set_template_ids":[1, 11],
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "-name"
  }
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 2,
    "info": [
      {
        "id": 1,
        "name": "zk1",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:09:23.859+08:00",
        "last_time": "2020-03-25T18:59:00.167+08:00",
        "bk_supplier_account": "0"
      },
      {
        "id": 11,
        "name": "q",
        "bk_biz_id": 10,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:10:05.176+08:00",
        "last_time": "2020-03-16T15:10:05.176+08:00",
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

#### data

| Name  | Type  | Description      |
|-------|-------|------------------|
| count | int   | Total count      |
| info  | array | Returned results |

#### info

| Name                | Type   | Description           |
|---------------------|--------|-----------------------|
| id                  | int    | Cluster template ID   |
| name                | array  | Cluster template name |
| bk_biz_id           | int    | Business ID           |
| creator             | string | Creator               |
| modifier            | string | Last modifier         |
| create_time         | string | Creation time         |
| last_time           | string | Update time           |
| bk_supplier_account | string | Supplier account      |
