### Functional description

update host properties in batches(can't update host's cloud area property)

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field               |  Type        | Required |  Description                                                   |
|---------------------|--------------|----------|----------------------------------------------------------------|
| bk_supplier_account | string       | No       | Supplier account                                               |
| update              | object array | Yes      | The updated properties and values for the host, limited to 500 |

#### update
| Field      | Type   | Required | Description                                                                  |
|------------|--------|----------|------------------------------------------------------------------------------|
| properties | object | Yes      | The updated properties for the host, can't update host's cloud area property |
| bk_host_id | int    | Yes      | The host ID, for host update                                                 |

#### properties
| Field        | Type   | Required | Description                                                                     |
|--------------|--------|----------|---------------------------------------------------------------------------------|
| bk_host_name | string | No       | Host name, or other properties, can't update host's cloud area property         |
| operator     | string | No       | Maintainer, or other properties, can't update host's cloud area property        |
| bk_comment   | string | No       | Comment, or other properties, can't update host's cloud area property           |
| bk_isp_name  | string | No       | Telecom operators, or other properties, can't update host's cloud area property |


### Request Parameters Example

```json
{
    "bk_supplier_account":"0",
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


### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```
