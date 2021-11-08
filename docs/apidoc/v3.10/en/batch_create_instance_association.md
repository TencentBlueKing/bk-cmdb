### Functional description

batch create instance association(v3.10.2+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field          | Type   | Required | Description                                     |
| -------------- | ------ | -------- | ----------------------------------------------- |
| bk_obj_id      | string | Yes      | source object id                                |
| bk_asst_obj_id | string | Yes      | target object id                                |
| bk_obj_asst_id | string | Yes      | unique id of association between object         |
| details        | array  | Yes      | contents of association that need to be created， maximum number of 
associations is 200 |

##### details:

| Field           | Type   | Required | Description              |
| --------------- | ------ | -------- | ------------------------ |
| bk_inst_id      | int | Yes      | source model instance id |
| bk_asst_inst_id | int | Yes      | target model instance id |

### Request Parameters Example

```json
{
	"bk_obj_id":"bk_switch",
  	"bk_asst_obj_id":"host",
	"bk_obj_asst_id": "bk_switch_belong_host",
 	"details":[
	  { 
		"bk_inst_id": 11,
 		"bk_asst_inst_id": 21
 	  },
 	  { 
 		"bk_inst_id": 12,
 		"bk_asst_inst_id": 22
 	  }
	]
}
```

### Return Result Example

```json
{
 	"result": true,
	"code": 0,
 	"message": "success",
	"data": {
        "success_created": {
            "0":73
        },
        "error_msg":{
            "1":"the association inst is not exist"
        }
    }
}
```

### Return Result Parameters Description

#### data

| Field           | Type   | Description                                                  |
| --------------- | ------ | ------------------------------------------------------------ |
| success_created | map | key is index of request's detail array, value is success created association id |
| error_msg       | map | key is index of request's detail array, value is error message |

