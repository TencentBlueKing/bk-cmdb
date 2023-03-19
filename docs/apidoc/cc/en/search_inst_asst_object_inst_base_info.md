

### Functional description

Query instance Association model instance basic information

### Request Parameters

{{ common_args_desc }}


#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| fields         |   array   | no     | Specify the fields to query. The parameter is any attribute of the business. If you do not fill in the field information, the system will return all the fields of the business|
| condition      |   object    | no     | Query criteria|
| page           |   object    | no     | Paging condition|

#### condition

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_obj_id |  string    | yes  | Instance model ID|
| bk_inst_id|   int    | yes | Instance ID|
|association_obj_id| string| yes | The model ID of the associated object, which returns the instance basic data (bk_inst_id,bk_inst_name) associated with the bk_inst_id instance of the Association_obj_id model|
|is_target_object|  bool |no| Whether bk_obj_id is the target model, the default is false, the source model in the Association relationship, otherwise, it is the target model|

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | no      | Record start position, default 0|
| limit    |   int    | no     | Limit bars per page, default 20, maximum 200|


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_obj_id":"bk_switch", 
		"bk_inst_id":12, 
		"association_obj_id":"host", 
		"is_target_object":true 
    },
    "page": {
        "start": 0,
        "limit": 10
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
        "count": 4,
        "info": [
            {
                "bk_inst_id": 1,
                "bk_inst_name": "127.0.0.3"
            }
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
| info|  object array |Model ID of associated object, instance basic data of instance associated model (bk_inst_id,bk_inst_name)|
| page|  object| Paging information|

#### Data.info Field Description:
| Name| Type| Description|
|---|---|---|
| bk_inst_id | int |Instance ID|
| bk_inst_name | string  |Instance name|

##### Data.info.BK_inst_id, data.info.BK_inst_name field descriptions

Values corresponding to bk_inst_id, bk_inst_name for different models

| Model   |  bk_inst_id   |  bk_inst_name |
|---|---|---|
|Business|  bk_biz_id | bk_biz_name|
|Set|  bk_set_id | bk_set_name|
|Module|  bk_module_id | bk_module_name|
|Process|  bk_process_id | bk_process_name|
|Host|  bk_host_id | bk_host_inner_ip|
|Universal model|  bk_inst_id | bk_inst_name|

