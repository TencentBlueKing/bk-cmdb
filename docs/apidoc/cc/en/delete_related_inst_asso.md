### Functional description

 Deletes an Association between instances based on the ID of the instance Association relationship. (Effective Version: 3.5.40)

### Request Parameters

{{ common_args_desc }}


#### Interface Parameters

| Field| Type     | Required| Description             |
| :--- | :------- | :--- | :--------------- |
| id   |  int |yes   | ID of the instance Association (note: Identity ID of non-model instance), up to 500|
| bk_obj_id | string |yes| The model unique name of the relationship source model|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id":[1,2],
    "bk_obj_id": "abc"
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": null
}
```

### Return Result Parameters Description

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|