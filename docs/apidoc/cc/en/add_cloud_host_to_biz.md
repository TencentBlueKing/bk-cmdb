### Function description

add cloud host to biz idle module (cloud host management dedicated interface, version: v3.10.19+, permission: edit business host)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| field     | type         | mandatory | description                                                                                                               |
|-----------|--------------|-----------|---------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int          | yes       | business id                                                                                                               |
| host_info | array | yes       | to be added cloud host information, array length is limited to 200, these hosts can only succeed or fail at the same time |

#### host_info

host information fields, cloud area ID and inner IP fields are required, other fields are attribute fields defined in host model. Only field examples are shown here, please fill in other fields as needed.

| field           | type   | required | description                                                       |
|-----------------|--------|----------|-------------------------------------------------------------------|
| bk_host_innerip | string | yes      | host inner ip in IPv4 format, multiple ips are seperated by comma |
| bk_cloud_id     | int    | yes      | cloud area id                                                     |
| bk_host_name    | string | no       | host name, or any other property                                  |
| operator        | string | no       | host main maintainer, or other attributes                         |
| bk_comment      | string | no       | comment, or other attributes                                      |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 123,
    "host_info": [
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host1",
            "operator": "admin",
            "bk_comment": "comment"
        },
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.2",
            "bk_host_name": "host2",
            "operator": "admin",
            "bk_comment": "comment"
        }
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "ids": [
            1,
            2
        ]
    }
}
```

### Return Result Parameters Description

#### response

| name       | type   | description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| data       | object | The data returned by the request.                                                         |
| permission | object | Permission information                                                                    |
| request_id | string | Request chain id                                                                          |

#### data
| field | type      | description                        |
|-------|-----------|------------------------------------|
| ids   | array | successfully created host id array |