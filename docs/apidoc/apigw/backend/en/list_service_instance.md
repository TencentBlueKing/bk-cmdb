### Description

Query service instance list based on business ID, with the option to include module ID and other information in the
query.

### Parameters

| Name         | Type      | Required | Description                                                                                           |
|--------------|-----------|----------|-------------------------------------------------------------------------------------------------------|
| bk_biz_id    | int       | Yes      | Business ID                                                                                           |
| bk_module_id | int       | No       | Module ID                                                                                             |
| bk_host_ids  | int array | No       | List of host IDs, supports up to 1000 host IDs                                                        |
| selectors    | int       | No       | Label filtering function, operator options: `=`, `!=`, `exists`, `!`, `in`, `notin`                   |
| page         | object    | No       | Pagination parameters                                                                                 |
| search_key   | string    | No       | Name filtering parameter, can be filled with characters included in the process name for fuzzy search |

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
  "bk_module_id": 56,
  "bk_host_ids":[26,10],
  "search_key": "",
  "selectors": [{
    "key": "key1",
    "operator": "notin",
    "values": ["value1"]
  },{
    "key": "key1",
    "operator": "in",
    "values": ["value1", "value2"]
  }]
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
        "id": 72,
        "name": "t1",
        "bk_host_id": 26,
        "bk_module_id": 62,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-06-20T22:46:00.69+08:00",
        "last_time": "2019-06-20T22:46:00.69+08:00",
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
