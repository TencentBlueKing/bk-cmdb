### Description

This interface is used to query concise relationship information directly related to the upper and lower levels (models)
of an instance in the business topology. (v3.10.1+)

If the business topology levels are from top to bottom: business, department (custom business level), cluster, module.
Then:

1. Upwards, you can query the relationship information of the directly superior **department** of a certain cluster.
2. Downwards, you can query the module relationship information directly associated with that cluster.

Conversely, you cannot directly query the relationship between a custom level instance **department** and the modules it
contains, as departments and modules are not directly associated.

### Parameters

| Name         | Type   | Required | Description                                                                                                                                                              |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| src_biz_obj  | string | Yes      | In the business level, the model ID of the source level, which can be "biz", the model ID of the custom level (bk_obj_id), "set", or "module".                           |
| src_ids      | array  | Yes      | A list of instance IDs representing src_biz_obj, with a list length ranging from [1, 200].                                                                               |
| dest_biz_obj | string | Yes      | The business level model directly (closely) related to src_biz_obj. For business ("biz"), any src_biz_obj's dest_biz_obj can be "biz". However, they cannot be the same. |
| page         | object | Yes      | Query the paging configuration information returned by the data                                                                                                          |

#### Explanation of the page field

| Name  | Type   | Required       | Description                                                                                                                      |
|-------|--------|----------------|----------------------------------------------------------------------------------------------------------------------------------|
| start | int    | Yes            | Record starting position, starting from 0                                                                                        |
| limit | int    | Yes            | Number of records per page, maximum 500                                                                                          |
| sort  | string | Not applicable | This field is set by default in the interface to sort by the ID of the associated (dest_biz_obj) identity. Do not set this field |

### Request Example

```json
{
    "src_biz_obj": "biz",
    "src_ids":[3,302],
    "dest_biz_obj":"nation",
    "page":{
        "start": 0,
        "limit": 2
    }
}
```

### Response Example

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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### Explanation of data

| Name      | Type | Description                                                                                                                     |
|-----------|------|---------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int  | The business ID to which the instance belongs                                                                                   |
| src_id    | int  | Consistent with the ID list of the src_ids input in the parameter, representing the ID of the instance queried by the parameter |
| dest_id   | int  | Corresponding to the model of dest_biz_obj and the instance directly associated with src_ids                                    |

Note:

1. If it is a downward query (from higher level to lower level), the method to judge whether the data fetching is
   complete is that the data array list returned is empty.
2. If it is an upward query (from lower level to higher level), this interface can return all query results at once. The
   condition is that the value of page.limit must be >= the length of src_ids.
