### Functional description

Delete service class based on service class ID

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| id            |  int  |yes   | Service class ID|
| bk_biz_id                  |  int        | yes  | Business ID |

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "id": 6
}
```

### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "28807929b7af4fcd9b834fd200ceb2ad"
}
```

### Return Result Parameters Description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|
