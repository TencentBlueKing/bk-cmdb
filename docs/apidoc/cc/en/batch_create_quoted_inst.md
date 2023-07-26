### Function description

batch create quoted model instance (version: v3.10.30+, permission: update permission of the source model instance)

### Request parameters

{{ common_args_desc }}

#### Interface parameters

| Field          | Type         | Required | Description                                          |
|----------------|--------------|----------|------------------------------------------------------|
| bk_obj_id      | string       | yes      | source model id                                      |
| bk_property_id | string       | yes      | source model quoted property id                      |
| data           | object array | yes      | instance data to be created, the maximum limit is 50 |

#### data[n]

| Field       | Type   | Required                                           | Description                                                                                                                  |
|-------------|--------|----------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|
| bk_inst_id  | int64  | no                                                 | source model instance id, if not set, the created instance should be associated with source model instance using create_inst |
| name        | string | depends on the "isrequired" config of the property | name, this is only an example, actual fields is defined by quoted model properties                                           |
| operator    | string | depends on the "isrequired" config of the property | operator, this is only an example, actual fields is defined by quoted model properties                                       | 
| description | string | depends on the "isrequired" config of the property | description, this is only an example, actual fields is defined by quoted model properties                                    |

### Request parameter examples

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "data": [
    {
      "bk_inst_id": 123,
      "name": "test",
      "operator": "user",
      "description": "test instance"
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
  "permission": null,
  "data": {
    "ids": [
      1,
      2
    ]
  },
  "request_id": "dsda1122adasadadada2222"
}
```

**Note:**

- The order of the array of ids in the returned data remains the same as the order of the array data in the parameters.

### Return result parameter description

#### response

| Name       | Type   | Description                                                                                         |
|------------|--------|-----------------------------------------------------------------------------------------------------|
| result     | bool   | The success or failure of the request. true: the request was successful; false: the request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                            |
| message    | string | The error message returned by the failed request.                                                   |
| permission | object | Permission information                                                                              |
| request_id | string | request_chain_id                                                                                    |
| data       | object | data returned by the request                                                                        |

#### data

| Name | Type  | Description                                        |
|------|-------|----------------------------------------------------|
| ids  | array | unique identifier array of created instances in cc |
