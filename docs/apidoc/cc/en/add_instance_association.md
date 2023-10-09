### Functional description

Add an association relationship between model instances.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description          |
|----------------------|------------|--------|-----------------------------|
| bk_obj_asst_id           |  string     |  yes   | The unique id of the Association between models|
| bk_inst_id           |  int64     |  yes  | Source model instance id|
| bk_asst_inst_id           |  int64     |  yes  | Target model instance id|
| metadata           | object     | yes | meta data             |


metadata params

| Field                 | Type      | Required	   | Description         |
|---------------------|------------|--------|-----------------------------|
| label           |  string map     |  yes  | Tag information|


label params

| Field                 | Type      | Required	   | Description         |
|---------------------|------------|--------|-----------------------------|
| bk_biz_id           |  string      |  yes  | Business id |

### Request Parameters Example

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_asst_id": "bk_switch_belong_bk_host",
    "bk_inst_id": 11,
    "bk_asst_inst_id": 21,
    "metadata":{
        "label":{
            "bk_biz_id":"1"
        }
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": 1038
    },
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
}

```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| data    |  object |Data returned by request                           |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |

#### data

| Field       | Type     | Description         |
|------------|----------|--------------|
|id| int64| New instance Association identity id|

