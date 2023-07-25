### Functional description

list quoted model instances by condition (version: v3.10.30+, permission: view permission of the source model instance)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field          | Type   | Required | Description                                                                                                                                                |
|----------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id      | string | yes      | source model id                                                                                                                                            |
| bk_property_id | string | yes      | source model quoted property id                                                                                                                            |
| filter         | object | no       | query filter of the quoted instance                                                                                                                        |
| fields         | array  | no       | quoted instance property list, which controls which fields are in the returned result, which can accelerate interface requests and reduce network traffic. |
| page           | object | yes      | page condition                                                                                                                                             |

#### filter

This parameter is the filter rule to search for container based on its attribute fields. This parameter supports the
following two filter rules types. The combined filter rules can be nested with the maximum nesting level of 2. The
specific supported filter rule types are as follows:

##### combined filter rule

This filter rule type defines filter rules composed of other rules, the combined rules support logic and/or
relationships

| Field     | Type   | Required | Description                                                                |
|-----------|--------|----------|----------------------------------------------------------------------------|
| condition | string | yes      | query criteria, support `AND` and `OR`                                     |
| rules     | array  | yes      | query rules, can be of `combined filter rule` or `atomic filter rule` type |

##### atomic filter rule

This filter rule type defines basic filter rules, which represent rules for filtering a field. Any filter rule is either
directly an atomic filter rule, or a combination of multiple atomic filter rules

| Field    | Type                                                                 | Required | Description                                                                                                          |
|----------|----------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------|
| field    | string                                                               | yes      | container's field                                                                                                    |
| operator | string                                                               | yes      | operator, optional values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | different fields and operators correspond to different value formats | yes      | operand                                                                                                              |

Assembly rules can refer to: <https://github.com/TencentBlueKing/bk-cmdb/blob/v3.10.x/pkg/filter/README.md>

#### page

| Field        | Type   | Required | Description                                                                                                                                                                                                        |
|--------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | yes      | Record start position                                                                                                                                                                                              |
| limit        | int    | yes      | Limit per page, maximum 500                                                                                                                                                                                        |
| sort         | string | no       | Sort the field                                                                                                                                                                                                     |
| enable_count | bool   | yes      | The flag defining Whether to get the the number of query objects. If this flag is true, then the request is to get the quantity. The remaining fields must be initialized, start is 0, and limit is: 0, sort is "" |

### Request Parameters Example

#### Query Detail Request Parameters Example

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "filter": {
    "condition": "AND",
    "rules": [
      {
        "field": "name",
        "operator": "not_equal",
        "value": "test"
      },
      {
        "condition": "OR",
        "rules": [
          {
            "field": "operator",
            "operator": "not_in",
            "value": [
              "me"
            ]
          },
          {
            "field": "bk_inst_id",
            "operator": "equal",
            "value": 123
          }
        ]
      }
    ]
  },
  "fields": [
    "name",
    "description"
  ],
  "page": {
    "start": 0,
    "limit": 2,
    "sort": "name",
    "enable_count": false
  }
}
```

#### Query Quantity Request Parameters Example

```json
{
  "bk_app_code": "code",
  "bk_app_secret": "secret",
  "bk_username": "xxx",
  "bk_token": "xxxx",
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "filter": {
    "field": "name",
    "operator": "equal",
    "value": "test"
  },
  "page": {
    "enable_count": true
  }
}
```

### Return Result Example

#### Query Detail Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": {
    "count": 0,
    "info": [
      {
        "name": "test1",
        "description": "test instance 1"
      },
      {
        "name": "test2",
        "description": "test instance 2"
      }
    ]
  }
}
```

#### Query Quantity Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": {
    "count": 5,
    "info": []
  }
}
```

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |

#### data

| Field | Type  | Description                                                     |
|-------|-------|-----------------------------------------------------------------|
| count | int   | Number of containers                                            |
| info  | array | Container list, only returns the fields that is set in `fields` |

#### info

| Field       | Type   | Description                                                                               |
|-------------|--------|-------------------------------------------------------------------------------------------|
| name        | string | name, this is only an example, actual fields is defined by quoted model properties        |
| description | string | description, this is only an example, actual fields is defined by quoted model properties |
