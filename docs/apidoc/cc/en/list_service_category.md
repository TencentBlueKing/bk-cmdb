### Function Description

Query the list of service categories, based on the business ID, including commonly used service categories.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type | Required | Description |
| --------- | ---- | -------- | ----------- |
| bk_biz_id | int  | Yes      | Business ID |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1
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
    "count": 20,
    "info": [
      {
        "bk_biz_id": 0,
        "id": 16,
        "name": "Apache",
        "bk_root_id": 14,
        "bk_parent_id": 14,
        "is_built_in": true
      },
      {
        "bk_biz_id": 0,
        "id": 19,
        "name": "Ceph",
        "bk_root_id": 18,
        "bk_parent_id": 18,
        "is_built_in": true
      },
      {
        "bk_biz_id": 1,
        "id": 1,
        "name": "Default",
        "bk_root_id": 1,
        "is_built_in": true
      }
    ]
  }
}
```

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by the request                                 |

#### data Field Explanation

| Field | Type  | Description              |
| ----- | ----- | ------------------------ |
| count | int   | Total number of records  |
| info  | array | List of returned results |

#### info Field Explanation

| Field               | Type   | Description                |
| ------------------- | ------ | -------------------------- |
| id                  | int    | Service category ID        |
| name                | string | Service category name      |
| bk_root_id          | int    | Root service category ID   |
| bk_parent_id        | int    | Parent service category ID |
| is_built_in         | bool   | Whether it is built-in     |