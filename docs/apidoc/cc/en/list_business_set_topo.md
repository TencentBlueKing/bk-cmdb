### Functional description

Query business set topology (v3.10.12+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_set_id    |  int    | yes | Business set ID|
| bk_parent_obj_id | string |yes| The parent object ID of the query model is required|
| bk_parent_id     |  int    | yes | The parent ID of the query model is required|

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

### Return Result Example

```json
{
  "result":true,
  "code":0,
  "message":"",
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

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| data    |  array |Data returned by request                           |
| request_id    |  string |Request chain id    |

#### data

| Name    | Type   | Description              |
| ------- | ------ | --------------- |
| bk_obj_id  | string   | Model object ID|
| bk_inst_id    |  int    | Model instance ID   |
| bk_inst_name | string |Model instance name   |
| default    |  int |Model instance classification    |


