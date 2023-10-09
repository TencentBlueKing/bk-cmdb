### Functional description

Lock the host according to the id list of the host, and add a new host lock. If the host has already been locked, it will also prompt that the locking is successful (v3.8.6).

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------|
|id_list|  int array| yes | Host ID list|


### Request Parameters Example

```python
{
   "bk_app_code": "esb_test",
   "bk_app_secret": "xxx",
   "bk_username": "xxx",
   "bk_token": "xxx",
   "id_list":[1, 2, 3]
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null,
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| data    |  object |Data returned by request                           |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |