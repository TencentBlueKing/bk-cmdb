### Functional description

Delete the object model attribute. You can delete the business user-defined field

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field| Type       | Required   | Description                         |
|-------|-------------|--------|-------------------------------|
| id    |  int         | no     | The unique identification ID of the deleted data record|


### Request Parameters Example

```python

{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id" : 0
}
```


### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
