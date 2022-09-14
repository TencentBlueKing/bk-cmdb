### Functional description

According to the unique identity id of the model instance Association relationship.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters
| Field                 | Type      | Required	   | Description          |
|----------------------|------------|--------|-----------------------------|
| id           |  int     |  yes | Unique identity id of the model instance Association             |
| bk_obj_id    |  string    |  yes | Source or target model id of the model instance Association (v3.10+)|

### Request Parameters Example

``` json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "id": 1
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}

```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|------------|----------|--------------|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

