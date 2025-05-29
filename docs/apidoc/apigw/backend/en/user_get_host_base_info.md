### Description

Get Host Basic Information Details (Permission: Host Pool Host View Permission)

### Parameters

| Name                | Type   | Required | Description                                               |
|---------------------|--------|----------|-----------------------------------------------------------|
| bk_supplier_account | string | No       | Developer account                                         |
| bk_host_id          | int    | Yes      | Host identity ID, i.e., the value of the bk_host_id field |

### Request Example

```json
{
    "bk_host_id": 10000
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_property_id": "bk_host_name",
            "bk_property_name": "host name",
            "bk_property_value": "centos7"
        },
        ......
        {
            "bk_property_id": "bk_host_id",
            "bk_property_name": "host ID",
            "bk_property_value": "10000"
        }
    ],
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Request returned data                                               |

#### data

| Name              | Type   | Description    |
|-------------------|--------|----------------|
| bk_property_id    | string | Property ID    |
| bk_property_name  | string | Property name  |
| bk_property_value | string | Property value |

**Note**

- If the host's property field is of table type, the returned bk_property_value is null. To query the value of the table
  type field, use the list_quoted_inst interface. Documentation
  link: [list_quoted_inst](https://github.com/TencentBlueKing/bk-cmdb/blob/v3.12.x/docs/apidoc/cc/zh_hans/list_quoted_inst.md)
