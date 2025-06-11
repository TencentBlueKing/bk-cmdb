### Description

Query the list of instances referencing the model (Version: v3.10.30+, Permission: When passing the business id, it
means querying the list of instances referencing the model from the perspective of the business. When the model is a
business, check the business view permission; otherwise, check the access permission of the corresponding model
instances.)

### Parameters

| Name           | Type         | Required | Description                                                                                                                                                                        |
|----------------|--------------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id      | string       | No       | Business id                                                                                                                                                                        |
| bk_obj_id      | string       | Yes      | Source model ID                                                                                                                                                                    |
| bk_property_id | string       | Yes      | ID of the property in the source model that references this model                                                                                                                  |
| filter         | object       | No       | Query conditions for instances referencing this model                                                                                                                              |
| fields         | string array | No       | List of properties of instances referencing this model, controls which fields are returned in the result, speeding up interface requests and reducing network traffic transmission |
| page           | object       | Yes      | Pagination information                                                                                                                                                             |

#### Explanation of the filter Field

Filtering rules for attributes of the model being referenced, used to search data based on attribute fields. This
parameter supports the following two types of filtering rules, and combination filtering rules can be nested, with a
maximum of two levels. The specific supported filtering rule types are as follows:

##### Combination Filtering Rules

Filtering rules composed of other rules, supporting logical AND/OR relationships between combined rules

| Name      | Type   | Required | Description                                                                           |
|-----------|--------|----------|---------------------------------------------------------------------------------------|
| condition | string | Yes      | Combined query condition, supports both `AND` and `OR`                                |
| rules     | array  | Yes      | Query rules, can be of type `Combination Filtering Rules` or `Atomic Filtering Rules` |

##### Atomic Filtering Rules

Basic filtering rules, indicating the rule for filtering a certain field. Any filtering rule is directly an atomic
filtering rule or composed of multiple atomic filtering rules

| Name     | Type                                                               | Required | Description                                                                                                                      |
|----------|--------------------------------------------------------------------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string                                                             | Yes      | Property field of the model being referenced                                                                                     |
| operator | string                                                             | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | Different field and operator correspond to different value formats | No       | Operation value                                                                                                                  |

Assembly rules can refer to: https://github.com/TencentBlueKing/bk-cmdb/blob/v3.10.x/pkg/filter/README.md

#### Explanation of the page Field

| Name         | Type   | Required | Description                                                                                                                                                                                                                 |
|--------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                                                                                                                                       |
| limit        | int    | Yes      | Number of records per page, maximum is 500                                                                                                                                                                                  |
| sort         | string | No       | Sorting field                                                                                                                                                                                                               |
| enable_count | bool   | Yes      | Flag to indicate whether to obtain the number of query objects. If this flag is true, it means that this request is to obtain the count. In this case, other fields must be initialized, start is 0, limit is 0, sort is "" |

### Request Example

#### Example of Detailed Information Request Parameters

```json
{
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

#### Example of Count Request Parameters

```json
{
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

### Response Example

#### Example of Detailed Information Response

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
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

#### Example of Count Response

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 5,
    "info": []
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description                                               |
|-------|-------|-----------------------------------------------------------|
| count | int   | Number of records                                         |
| info  | array | Actual data, only the fields set in `fields` are returned |

#### info

| Name        | Type   | Description                                                                            |
|-------------|--------|----------------------------------------------------------------------------------------|
| name        | string | Name, this is just an example, the actual fields depend on the model properties        |
| description | string | Description, this is just an example, the actual fields depend on the model properties |
