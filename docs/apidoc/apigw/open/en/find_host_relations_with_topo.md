### Description

Query host relationship information based on the business topology instance node (Permission: Business access
permission)

### Parameters

| Name        | Type   | Required | Description                                                                                                                                                                                                     |
|-------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| page        | dict   | Yes      | Query conditions                                                                                                                                                                                                |
| fields      | array  | Yes      | List of host properties, control which fields of the host information should be returned. Please fill in according to your needs. It can be bk_biz_id, bk_host_id, bk_module_id, bk_set_id, bk_supplier_account |
| bk_obj_id   | string | Yes      | Model ID of the topology node, it can be a custom hierarchical model ID, set, module, etc., but cannot be a business                                                                                            |
| bk_inst_ids | array  | Yes      | List of instance IDs of the topology node, supports up to 50 instance nodes                                                                                                                                     |
| bk_biz_id   | int    | Yes      | Business ID                                                                                                                                                                                                     |

#### page

| Name  | Type   | Required | Description                                      |
|-------|--------|----------|--------------------------------------------------|
| start | int    | Yes      | Record start position                            |
| limit | int    | Yes      | Number of records per page, maximum value is 500 |
| sort  | string | No       | Sorting field                                    |

### Request Example

```json
{
    "bk_biz_id": 1,
    "page": {
        "start": 0,
        "limit": 10
    },
    "fields": [
        "bk_module_id",
        "bk_host_id"
    ],
    "bk_obj_id": "province",
    "bk_inst_ids": [10,11]
}
```

### Response Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "data":  {
      "count": 1,
      "info": [
          {
              "bk_host_id": 2,
              "bk_module_id": 51
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

#### data

| Name  | Type  | Description                   |
|-------|-------|-------------------------------|
| count | int   | Number of records             |
| info  | array | Host relationship information |

#### info

| Name         | Type | Description |
|--------------|------|-------------|
| bk_host_id   | int  | Host ID     |
| bk_module_id | int  | Module ID   |
