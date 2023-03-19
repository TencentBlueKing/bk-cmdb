### Functional description

This interface is used to query the concise relationship information of the upper and lower levels (models) directly associated with an instance of a certain level (model) in the business topology. (v3.10.1+)


If the business topology level is business, Department (user-defined business level), set and module from top to bottom. Then:


1. You can query the relationship information of the direct superior Department to which a set belongs upward;


2. The module relationship information directly associated with the set can be queried downward.


Conversely, you can not directly query the module relationship contained in a Department of a custom level instance through Department, because Department and module are not directly associated.


### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| src_biz_obj  | string  |yes     | In business hierarchy, the model ID of source hierarchy can be "biz," user-defined hierarchy model ID(bk_obj_id),"set," and "module." |
| src_ids  | array  |yes     | List of instance IDs represented by src_biz_obj, with a list length in the range of [1200]|
| dest_biz_obj  | string  |yes     | The business hierarchy model directly (immediately) associated with src_biz_obj.  Where the business ("biz") As an exception, dest_biz_obj can be "biz" for any src_biz_obj. But the two are not allowed to be the same. |
| page  | object  |yes     | Paging configuration information returned by queried data|

#### Page field Description

| Field| Type   | Required| Description                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | yes | Record start position, starting from 0         |
| limit | int    | yes | Limit bars per page, Max. 500|
| sort | string    | Unavailable   | This field is sorted by the identity ID of the associated (dest_biz_obj) by default in the interface. Please do not set this field|



### Request Parameters Example

```json
{
    "bk_app_code": "xxx",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### Data description
| Field      | Type      | Description      |
|-----------|------------|------------|
| bk_biz_id | int   | The business ID to which this instance belongs     |
| src_id | int   | It is consistent with the ID list entered by the src_ids in the parameter. Represents the instance ID of the input query model|
| dest_id | int| The instance ID directly associated with the model corresponding to dest_biz_obj in the parameter and the instance corresponding to src_ids|

Note:

1. If it is a downward query (query from the level to the lower level), it is judged that the method of paging and pulling data is that the returned data array list is empty.


2. In the case of an upward query (from a low level to a high level), the interface can return all query results at once, provided that the value of page.limit is>= the length of src_ids.
