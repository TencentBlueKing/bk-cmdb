### Functional description

Query host locks from host id list (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------|
|id_list|  array| yes | Host ID list|


### Request Parameters Example

```python
{
   "bk_app_code": "esb_test",
   "bk_app_secret": "xxx",
   "bk_username": "xxx",
   "bk_token": "xxx",
   "id_list":[1, 2]
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
    "data": {
        1: true,
        2: false
    }
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data
| Field      | Type      | Description         |
|-----------|-----------|--------------|
| data |object| The data returned by the request, key is ID, and value is locked|
