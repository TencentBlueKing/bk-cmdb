### Description

Get business idle hosts, faulty hosts, and modules to be recycled based on the business ID.

### Parameters

| Name      | Type | Required | Description |
|-----------|------|----------|-------------|
| bk_biz_id | int  | Yes      | Business ID |

### Request Example

```python
{
    "bk_biz_id":0
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
    "bk_set_id": 2,
    "bk_set_name": "idle pool",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "idle host",
        "default": 1,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "fault host",
        "default": 2,
        "host_apply_enabled": false
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "recycle host",
        "default": 3,
        "host_apply_enabled": false
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### Explanation of data Parameters

| Name        | Type   | Description                                                                                   |
|-------------|--------|-----------------------------------------------------------------------------------------------|
| bk_set_id   | int64  | Instance ID of the set to which idle hosts, faulty hosts, and modules to be recycled belong   |
| bk_set_name | string | Instance name of the set to which idle hosts, faulty hosts, and modules to be recycled belong |
| module      | array  | Information about idle hosts, faulty hosts, and modules to be recycled                        |

#### Explanation of module Parameters

| Name               | Type   | Description                                                          |
|--------------------|--------|----------------------------------------------------------------------|
| bk_module_id       | int    | Instance ID of idle hosts, faulty hosts, or modules to be recycled   |
| bk_module_name     | string | Instance name of idle hosts, faulty hosts, or modules to be recycled |
| default            | int    | Indicates the module type                                            |
| host_apply_enabled | bool   | Whether to enable automatic application of host properties           |
