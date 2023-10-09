### Functional description

Query cloud region

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description       |
|----------------------|------------|--------|-------------|
|condition| object| no | Query criteria|
| page|  object| yes | Paging information|


#### condition
| Field                 | Type      | Required   | Description       |
|----------------------|------------|--------|-------------|
|bk_cloud_id| int| no | Cloud area ID |
|bk_cloud_name| string| no | Cloud area name |

#### Page field Description

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
|start| int| no | Get data offset position|
|limit| int| yes | Limit on the number of data pieces in the past, 200 is recommended|


### Request Parameters Example

``` python
{

    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_cloud_id": 12,
        "bk_cloud_name" "aws",
    },
    "page":{
        "start":0,
        "limit":10
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
  "data": {
    "count": 10,
    "info": [
         {
            "bk_cloud_id": 0,
            "bk_cloud_name": "aws",
            "bk_supplier_account": "0",
            "create_time": "2019-05-20T14:59:48.354+08:00",
            "last_time": "2019-05-20T14:59:48.354+08:00"
        },
        .....
    ]
   
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

| Name| Type| Description|
|---|---|---|
| count|  int| Number of records|
| info|  array |List information of queried cloud area|

#### Data.info Field Description:
| Name| Type| Description|
|---|---|---|
| bk_cloud_id | int |Cloud area ID|
| bk_cloud_name | string  |Cloud area name|
| create_time | string |Settling time|
| last_time | string |Last modified time|



