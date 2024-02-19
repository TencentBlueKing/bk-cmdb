### Function Description

Query basic information of instance-associated model instances

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| fields    | array  | No       | Specify the fields to be queried. The parameter can be any attribute of the business. If the field information is not filled in, the system will return all fields of the business |
| condition | object | Yes      | Query conditions                                             |
| page      | object | No       | Pagination conditions                                        |

#### condition

| Field              | Type   | Required | Description                                                  |
| ------------------ | ------ | -------- | ------------------------------------------------------------ |
| bk_obj_id          | string | Yes      | Model ID of the instance                                     |
| bk_inst_id         | int    | Yes      | Instance ID                                                  |
| association_obj_id | string | Yes      | Model ID of the associated object. Returns the basic data (bk_inst_id, bk_inst_name) of instances associated with association_obj_id model and bk_inst_id instance |
| is_target_object   | bool   | No       | Whether bk_obj_id is the target model. Default is false, which means it is the source model in the association relationship; otherwise, it is the target model |

#### page

| Field | Type | Required | Description                                                  |
| ----- | ---- | -------- | ------------------------------------------------------------ |
| start | int  | No       | Record start position, default value is 0                    |
| limit | int  | No       | Number of records per page, default value is 20, maximum is 200 |

### Request Parameter Example

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

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |

#### data

| Field  | Type         | Description                                                  |
| ----- | ------------ | ------------------------------------------------------------ |
| count | int          | Number of records                                            |
| info  | object array | Model ID of the associated object. Basic data of instances associated with the instance model (bk_inst_id, bk_inst_name) |
| page  | object       | Pagination information                                       |

#### data.info Field Explanation：

| Field         | Type   | Description   |
| ------------ | ------ | ------------- |
| bk_inst_id   | int    | Instance ID   |
| bk_inst_name | string | Instance name |

##### data.info.bk_inst_id, data.info.bk_inst_name Field Explanation

Values corresponding to different model bk_inst_id, bk_inst_name

| Model        | bk_inst_id    | bk_inst_name     |
| ------------ | ------------- | ---------------- |
| Business     | bk_biz_id     | bk_biz_name      |
| Cluster      | bk_set_id     | bk_set_name      |
| Module       | bk_module_id  | bk_module_name   |
| Process      | bk_process_id | bk_process_name  |
| Host         | bk_host_id    | bk_host_inner_ip |
| Common Model | bk_inst_id    | bk_inst_name     |