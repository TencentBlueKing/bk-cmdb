### Description

Batch Update Host Properties Based on Host ID and Attributes (Cannot be used to update the control area field in host
properties, Permission: Business Host Editing Permission)

### Parameters

| Name       | Type   | Required | Description                                                                                                       |
|------------|--------|----------|-------------------------------------------------------------------------------------------------------------------|
| properties | object | Yes      | Properties and values to be updated for hosts, cannot be used to update the control area field in host properties |
| bk_host_id | int    | Yes      | Host ID used for updating                                                                                         |

#### properties

| Name         | Type   | Required | Description                                                                                                       |
|--------------|--------|----------|-------------------------------------------------------------------------------------------------------------------|
| bk_host_name | string | No       | Host name, can also be other properties, cannot be used to update the control area field in host properties       |
| operator     | string | No       | Main maintainer, can also be other properties, cannot be used to update the control area field in host properties |
| bk_comment   | string | No       | Remark, can also be other properties, cannot be used to update the control area field in host properties          |
| bk_isp_name  | string | No       | Affiliated ISP, can also be other properties, cannot be used to update the control area field in host properties  |

### Request Example

```json
{
    "update":[
      {
        "properties":{
          "bk_host_name":"batch_update",
          "operator": "admin",
          "bk_comment": "test",
          "bk_isp_name": "1"
        },
        "bk_host_id":46
      }
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
