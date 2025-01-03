### Function Description

Create Service Category (Permission: Service Category Creation Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field        | Type   | Required | Description           |
| ------------ | ------ | -------- | --------------------- |
| name         | string | Yes      | Service category name |
| bk_parent_id | int    | No       | Parent node ID        |
| bk_biz_id    | int    | Yes      | Business ID           |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_parent_id": 0,
  "bk_biz_id": 1,
  "name": "test101"
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
    "bk_biz_id": 1,
    "id": 6,
    "name": "test101",
    "bk_root_id": 5,
    "bk_parent_id": 5,
    "is_built_in": false
  }
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Newly created service category information                   |

#### data

| Field               | Type    | Description                                                  |
| ------------------- | ------- | ------------------------------------------------------------ |
| id                  | integer | Service category ID                                          |
| root_id             | integer | Service category root node ID                                |
| parent_id           | integer | Service category parent node ID                              |
| is_built_in         | bool    | Whether it is a built-in node (built-in nodes cannot be edited) |
| bk_biz_id           | int     | Business ID                                                  |
| name                | string  | Service category name                                        |