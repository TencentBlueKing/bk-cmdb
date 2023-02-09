### Functional description

Create a cloud area based on the cloud area name

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description       |
|----------------------|------------|--------|-------------|
| bk_cloud_name  | string     | yes     | Cloud area name |

### Request Parameters Example

``` python
{
    
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_cloud_name": "test1"
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
        "created": {
            "origin_index": 0,
            "id": 6
        }
    }
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field          | Type     | Description     |
|---------------|----------|----------|
| created      |  object   | Create successfully, return message|


#### data.created

| Name    | Type   | Description       |
|---------|--------|------------|
| origin_index|  int |The result order of the corresponding request|
| id|  int |Cloud zone id, bk_Cloud_id|


