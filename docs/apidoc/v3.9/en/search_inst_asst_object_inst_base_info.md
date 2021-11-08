
### Functional description

search instance association topology

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
| fields         |  array   | No     | Specifies the field of the query. The parameter is any attribute of the service. If the field information is not filled in, the system returns all fields of the service.|
| condition      |  dict    | No     | query condition|
| page           |  dict    | No     | page condition |

#### condition

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
| bk_obj_id |  string    | Yes     | instance object ID |
| bk_inst_id|  int    |  Yes    | instacne id |
|association_obj_id|string|  Yes  | The model ID of the associated object, returning the instance basic data (bk_inst_id, bk_inst_name) associated with the bk_inst_id instance in the association_obj_id model|
|is_target_object| bool |  No |whether bk_obj_id is the target model, default false, the source model in the association, and No is the target model|

#### page

| Field                |  Type      | Required	   |  Description                       |
|---------------------|------------|--------|-----------------------------|
#### page params

| Field                 |  Type      | Required	   |  Description       | 
|--------|------------|--------|------------|
|start|int|No|get the data offset location|
|limit|int|No|The number of data points in the past is limited, default value 20, max:200|

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "condition": {
        "bk_obj_id":"bk_switch", 
		"bk_inst_id":12, 
		"association_obj_id":"host", 
		"is_target_object":true, 
    },
    "page": {
        "start": 0,
        "limit": 10,
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
    "data": {
        "count": 4,
        "info": [
            {
                "bk_inst_id": 1,
                "bk_inst_name": "127.0.0.3"
            }
        ],
        "page": {
            "start": 0,
            "limit": 1
        }
    }
}
```

### Return Result Parameters Description

#### data

| Field      | Type      | Description         |
|-----------|-----------|--------------|
| count| int| the num of record |
| info| object array |  the associated object, instance basic data of the instance association model（bk_inst_id,bk_inst_name） |
| page| object| page info|

#### data.info ：
| Field      | Type      | Description         |
|-----------|-----------|--------------|
| bk_inst_id | int | instance id |
| bk_inst_name | string  | instance nae  | 

##### data.info.bk_inst_id,data.info.bk_inst_name :

The values corresponding to different models bk_inst_id,bk_inst_name

| model| module   | bk_inst_id   | bk_inst_name |
|---|---|---|---|
|business | bk_biz_id | bk_biz_name|
|set | bk_set_id | bk_set_name|
|module | bk_module_id | bk_module_name|
|process | bk_process_id | bk_process_name|
|host | bk_host_id | bk_host_inner_ip|
|object | bk_inst_id | bk_inst_name|



#### data.page 

| Field       | Type     | Description         |
|------------|----------|--------------|
|start|int|server obtains the data offset position this time|
|limit|int|server client return data limit|

