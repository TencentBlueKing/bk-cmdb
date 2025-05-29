### Description

Query basic information of instance-associated model instances

### Parameters

| Name      | Type   | Required | Description                                                                                                                                                                        |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| fields    | array  | No       | Specify the fields to be queried. The parameter can be any attribute of the business. If the field information is not filled in, the system will return all fields of the business |
| condition | object | Yes      | Query conditions                                                                                                                                                                   |
| page      | object | No       | Pagination conditions                                                                                                                                                              |

#### condition

| Name               | Type   | Required | Description                                                                                                                                                        |
|--------------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id          | string | Yes      | Model ID of the instance                                                                                                                                           |
| bk_inst_id         | int    | Yes      | Instance ID                                                                                                                                                        |
| association_obj_id | string | Yes      | Model ID of the associated object. Returns the basic data (bk_inst_id, bk_inst_name) of instances associated with association_obj_id model and bk_inst_id instance |
| is_target_object   | bool   | No       | Whether bk_obj_id is the target model. Default is false, which means it is the source model in the association relationship; otherwise, it is the target model     |

#### page

| Name  | Type | Required | Description                                                     |
|-------|------|----------|-----------------------------------------------------------------|
| start | int  | No       | Record start position, default value is 0                       |
| limit | int  | No       | Number of records per page, default value is 20, maximum is 200 |

### Request Example

```python
{
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

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name  | Type         | Description                                                                                                              |
|-------|--------------|--------------------------------------------------------------------------------------------------------------------------|
| count | int          | Number of records                                                                                                        |
| info  | object array | Model ID of the associated object. Basic data of instances associated with the instance model (bk_inst_id, bk_inst_name) |
| page  | object       | Pagination information                                                                                                   |

#### data.info Field Explanation：

| Name         | Type   | Description   |
|--------------|--------|---------------|
| bk_inst_id   | int    | Instance ID   |
| bk_inst_name | string | Instance name |

##### data.info.bk_inst_id, data.info.bk_inst_name Field Explanation

Values corresponding to different model bk_inst_id, bk_inst_name

| Model        | bk_inst_id    | bk_inst_name     | 
|--------------|---------------|------------------|
| Business     | bk_biz_id     | bk_biz_name      |
| Cluster      | bk_set_id     | bk_set_name      |
| Module       | bk_module_id  | bk_module_name   |
| Process      | bk_process_id | bk_process_name  |
| Host         | bk_host_id    | bk_host_inner_ip |
| Common Model | bk_inst_id    | bk_inst_name     |
