### Functional description

Remove the tag from the service instance under the specified service according to the service id, the service instance id, and the tag to be removed

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| instance_ids            | array  | Yes   | service inscc/apidocs/en/list_service_instance_detail.mdtances ID array, the max length is 500 |
| keys            | array  | Yes   | key of lables to be remove |
| bk_biz_id            |  int  |yes   | Business ID |


### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "instance_ids": [60, 62],
  "keys": ["value1", "value3"]
}
```

### Return Result Example

```python
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

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|
