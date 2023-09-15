### Function description

batch delete quoted model instance (version: v3.10.30+, permission: update permission of the source model instance)

### Request parameters

{{ common_args_desc }}

#### Interface parameters

| Field          | Type        | Required | Description                                                        |
|----------------|-------------|----------|--------------------------------------------------------------------|
| bk_obj_id      | string      | yes      | source model id                                                    |
| bk_property_id | string      | yes      | source model quoted property id                                    |
| ids            | int64 array | yes      | id list of quoted instance to be deleted, the maximum limit is 500 |

### Request parameter examples

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "ids": [
    1,
    2
  ]
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "dsda1122adasadadada2222"
}
```

### Return result parameter description

#### response

| Name       | Type   | Description                                                                                         |
|------------|--------|-----------------------------------------------------------------------------------------------------|
| result     | bool   | The success or failure of the request. true: the request was successful; false: the request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                            |
| message    | string | The error message returned by the failed request.                                                   |
| permission | object | Permission information                                                                              |
| request_id | string | request_chain_id                                                                                    |
