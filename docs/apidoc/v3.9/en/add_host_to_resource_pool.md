### Functional description

add host to resource pool

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field               |  Type        | Required |  Description                                           |
|---------------------|--------------|----------|--------------------------------------------------------|
| bk_supplier_account | string       | Yes      | Supplier account                                       |
| host_info           | object array | Yes      | The info of the hosts to be added to the resource pool |
| directory           | int          | No       | The directory ID the hosts to be added to              |

#### host_info
| Field           | Type   | Required | Description                     |
|-----------------|--------|----------|---------------------------------|
| bk_host_innerip | string | Yes      | Host inner IP                   |
| bk_host_name    | string | No       | Host name, or other properties  |
| operator        | string | No       | Maintainer, or other properties |
| bk_comment      | string | No       | Comment, or other properties    |


### Request Parameters Example

```json
{
    "bk_supplier_account": "0",
    "host_info": [
        {
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host1",
            "operator": "admin"
        },
        {
            "bk_host_innerip": "",
            "bk_comment": "comment"
        }
    ],
    "directory": 1
}
```


### Return Result Example

```json
{
  "result": false,
  "code": 1110004,
  "message": "Failed to create host",
  "permission": null,
  "data": {
    "success": [
      {
        "index": 0,
        "bk_host_id": 11,
      }
    ],
    "error": [
      {
        "index": 1,
        "error_message": "'bk_host_innerip' unassigned",
      }
    ]
  }
}
```

#### response

| Field   | Type   | Description                                            |
| ------- | ------ | ------------------------------------------------------ |
| result  | bool   | request success or failed. true:success；false: failed |
| code    | int    | error code. 0: success, >0: something error            |
| message | string | error info description                                 |
| data    | object | response data                                          |

#### data description

| Field   | Type  | Description                            |
| ------- | ----- | -------------------------------------- |
| success | array | successfully added hosts' info         |
| error   | array | unsuccessfully added hosts' error info |

#### success description

| Field      | Type | Description                              |
| ---------- | ---- | ---------------------------------------- |
| index      | int  | successfully added hosts' index in array |
| bk_host_id | int  | host ID of the host                      |

#### error description

| Field         | Type   | Description                                |
| ------------- | ------ | ------------------------------------------ |
| index         | int    | unsuccessfully added hosts' index in array |
| error_message | string | error message of the failure               |