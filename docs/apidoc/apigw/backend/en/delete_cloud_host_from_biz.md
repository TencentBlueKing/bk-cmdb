### Description

Delete cloud hosts from the idle machine cluster of a business (Dedicated interface for cloud host management, Version:
v3.10.19+, Permission: Business host editing permission)

### Parameters

| Name        | Type  | Required | Description                                                                                                                            |
|-------------|-------|----------|----------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id   | int   | Yes      | Business ID                                                                                                                            |
| bk_host_ids | array | Yes      | Array of cloud host IDs to be deleted. The array length is at most 200, and a batch of hosts can only succeed or fail at the same time |

**Note: This interface can only delete cloud hosts. Filling in the IDs of other ordinary hosts will result in an error
that the host does not exist. The bk_cloud_host_identifier field in the properties of a cloud host is true, while for
other ordinary hosts, it is false. Cloud hosts can be added using cloud host-related interfaces such as (
add_cloud_host_to_biz Add cloud hosts to the business's idle machine module).**

### Request Example

```json
{
    "bk_biz_id": 123,
    "bk_host_ids": [
        1,
        2
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
