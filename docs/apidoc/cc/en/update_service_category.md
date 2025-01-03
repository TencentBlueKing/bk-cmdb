### Function Description

Update service category (Currently, only the name field can be updated. Permission: Service Category Editing Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description           |
| --------- | ------ | -------- | --------------------- |
| id        | int    | Yes      | Service category ID   |
| name      | string | Yes      | Service category name |
| bk_biz_id | int    | Yes      | Business ID           |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "id": 3,
  "name": "222"
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "f87f975a8f4a44ddbf6606ec432741a6",
    "data": {
        "bk_biz_id": 3,
        "id": 22,
        "name": "api",
        "bk_root_id": 21,
        "bk_parent_id": 21,
        "is_built_in": false
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
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Updated service category information                         |

#### data

| Field                | Type   | Description                      |
| ------------------- | ------ | -------------------------------- |
| bk_biz_id           | int    | Business ID                      |
| id                  | int    | Service category ID              |
| name                | string | Service category name            |
| bk_root_id          | int    | Root service category ID         |
| bk_parent_id        | int    | Parent service category ID       |
| is_built_in         | bool   | Whether it is a built-in service |