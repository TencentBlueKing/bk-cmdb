### Functional description

Tag a service instance based on the service instance id and the tag set 

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
|instance_ids| array| yes | Service instance ID, the max length is 100|
|labels| object| yes | Label added|
| bk_biz_id  | int     | yes     | Business ID |

#### Labels Field Description
- key verification rules:`^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- Value verification rule:`^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
  }
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

### Return result description

#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |


