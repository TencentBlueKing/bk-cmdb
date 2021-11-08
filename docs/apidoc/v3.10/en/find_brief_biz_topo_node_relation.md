### Functional description

this api is used to find the brief relations in a business topology from the source mainline object to the destination
mainline object with positive direction (as is from biz to set) or negative direction (as is from set to biz). (v3.
10.1+)

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field      | Type      | Required | Description                                                  |
| ---------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id  | int64       | Yes      | Business ID                                                  |
| src_biz_obj  | string  | Yes     | the source mainline object，can be one of the "biz"、custom level(bk_obj_id)、"set"、"module". |
| src_ids  | array int  | Yes     |  the instance id list of the source object(src_biz_obj), the length range is [1,200]|
| dest_biz_obj  | string  | Yes     | the destination mainline object, which should be the neighbour of the source object, except biz model. and it should not be same with the source object.|
| page  | object  | Yes     |  page information of the search result|

#### page description

| Field | Type   | Required | Description                                       |
| ----- | ------ | -------- | ------------------------------------------------- |
| start | int    | Yes       | start record, from 0.                                     |
| limit | int    | Yes       | page limit, maximum value is 500                 |
| sort | string    | do not set   | it has a default value of the destination object's instance id. |



### Request Parameters Example

```json
{
    "bk_app_code": "xxx",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
	"src_biz_obj": "biz",
	"src_ids":[3,302],
	"dest_biz_obj":"nation",
	"page":{
		"start": 0,
		"limit": 2
	}
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data":
    [
        {
            "bk_biz_id": 3,
            "src_id": 3,
            "dest_id": 3812
        },
        {
            "bk_biz_id": 302,
            "src_id": 302,
            "dest_id": 3813
        }
    ]
}
```

### Return Result Parameters Description

#### data description
| Field       | Type     | Description         |
|------------|----------|--------------|
| bk_biz_id | int   | business's instance id     |
| src_id | int   | the source object's instance id, which is one of the src_ids. |
| dest_id | int| the destination object's instance id, which is related with the src id. |

Note：

1. if you search from the top to bottom way in the business mainline topology, the way you check if you have already 
   found all the relations is that you get 0 data array in the response.


2. if you search from the bottom to top way in the business mainline topology, the response will return all the 
   relations if page.limit is >= the length of src_ids.
