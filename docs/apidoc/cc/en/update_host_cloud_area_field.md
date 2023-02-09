### Functional description

Update the host's cloud  area  field based on the host id list and cloud  area id

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required	   | Description                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            |  int  |no   | Business ID |
| bk_cloud_id         |  int  |yes   | Cloud area ID |
| bk_host_ids         |  array  | yes      | Host IDs, up to 2000|


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_ids": [43, 44], 
    "bk_cloud_id": 27,
    "bk_biz_id": 1
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": ""
}
```

### Return result instance cloud area + intranet IP duplicate

```python
{
  "result": false,
  "code": 1199014,
  "message": "data uniqueness check failed, bk_host_innerip repeated",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null
}
```

### Return result instance too many hosts for one operation
```python
{
  "result": false,
  "code": 1199077,
  "message": "the number of records per operation exceeds the maximum limit: 2000",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null
}
```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |
