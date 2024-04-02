### Function Description

Query business set topology (Version: v3.10.12+, Permission: Business set access)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field            | Type   | Required | Description                                |
| ---------------- | ------ | -------- | ------------------------------------------ |
| bk_biz_set_id    | int    | Yes      | Business set ID                            |
| bk_parent_obj_id | string | Yes      | ID of the parent object to query the model |
| bk_parent_id     | int    | Yes      | ID of the parent to query the model        |

### Request Parameters Example

```json
{
  "bk_app_code":"esb_test",
  "bk_app_secret":"xxx",
  "bk_username":"xxx",
  "bk_token":"xxx",
  "bk_biz_set_id":3,
  "bk_parent_obj_id":"bk_biz_set_obj",
  "bk_parent_id":344
}
```

### Response Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission":null,
  "data":[
    {
      "bk_obj_id":"bk_biz_set_obj",
      "bk_inst_id":5,
      "bk_inst_name":"xxx",
      "default":0
    },
    {
      "bk_obj_id":"bk_biz_set_obj",
      "bk_inst_id":6,
      "bk_inst_name":"xxx",
      "default":0
    }
  ],
  "request_id": "dsda1122adasadadada2222"
}
```

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| data       | array  | Data returned by the request                                 |
| request_id | string | Request chain ID                                             |

#### data

| Field         | Type   | Description                   |
| ------------ | ------ | ----------------------------- |
| bk_obj_id    | string | Model object ID               |
| bk_inst_id   | int    | Model instance ID             |
| bk_inst_name | string | Model instance name           |
| default      | int    | Model instance classification |