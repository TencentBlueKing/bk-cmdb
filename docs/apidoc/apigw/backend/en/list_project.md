### Description

Query projects (Version: v3.10.23+, Permission: View permission for the project)

### Parameters

| Name   | Type   | Required | Description                                                                                                                               |
|--------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------------|
| filter | object | No       | Query conditions                                                                                                                          |
| fields | array  | No       | Property list, controls which fields are returned in the result, speeding up interface requests and reducing network traffic transmission |
| time_condition | object | No | Query criteria for querying business by time |
| page   | object | Yes      | Pagination information                                                                                                                    |

#### filter Field Explanation

Attribute field filtering rules, used to search data based on attribute fields. This parameter supports the following
two types of filtering rules, and combination filtering rules can be nested, with a maximum of two levels. The specific
supported filtering rule types are as follows:

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
| field    | string                                                             | Yes      | Field of the container                                                                                                           |
| operator | string                                                             | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | Different field and operator correspond to different value formats | No       | Operation value                                                                                                                  |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### time_condition

| Field   | Type   | Required| Description              |
|-------|--------|-----|--------------------|
| oper  | string |Yes| Operator, currently only and is supported|
| rules | array  |Yes| Time query criteria         |

#### page

| Name         | Type   | Required | Description                                                                                                                                                                                                                 |
|--------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| start        | int    | Yes      | Record start position                                                                                                                                                                                                       |
| limit        | int    | Yes      | Number of records per page, maximum is 500                                                                                                                                                                                  |
| sort         | string | No       | Sorting field                                                                                                                                                                                                               |
| enable_count | bool   | Yes      | Flag to indicate whether to obtain the number of query objects. If this flag is true, it means that this request is to obtain the count. In this case, other fields must be initialized, start is 0, limit is 0, sort is "" |

### Request Example

#### Get Detailed Information Request Parameters

```json
{
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "bk_status",
                "operator": "equal",
                "value": "enable"
            }
        ]
    },
    "time_condition": {
      "oper": "and",
      "rules": [
        {
          "field": "create_time",
          "start": "2021-05-13 01:00:00",
          "end": "2021-05-14 01:00:00"
        }
      ]
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "id",
        "enable_count": false
    }
}
```

#### Get Count Request Example

```json
{
    "filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "equal",
                "value": 1
            },
            {
                "field": "bk_status",
                "operator": "equal",
                "value": "enable"
            }
        ]
    },
    "time_condition": {
      "oper": "and",
      "rules": [
        {
          "field": "create_time",
          "start": "2021-05-13 01:00:00",
          "end": "2021-05-14 01:00:00"
        }
      ]
    },
    "page": {
        "enable_count":true
    }
}
```

### Response Example

#### Response for Detailed Information Interface

```json
{
    "result": true,
    "code": 0,
    "data": {
        "count": 0,
        "info": [
            {	
               "id": 1,
               "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
               "bk_project_name": "test",
               "bk_project_code": "test",
               "bk_project_desc": "test project",
               "bk_project_type": "mobile_game",
               "bk_project_sec_lvl": "public",
               "bk_project_owner": "admin",
               "bk_project_team": [1, 2],
               "bk_status": "enable",
               "bk_project_icon": "https://127.0.0.1/file/png/11111",
               "bk_supplier_account": "0",
               "create_time": "2022-12-22T11:22:17.504+08:00",
               "last_time": "2022-12-22T11:23:31.728+08:00"
            }
        ]
    },
    "message": "success",
    "permission": null,
}
```

#### Count Response Example

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "count":1,
        "info":[
        ]
    },
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

#### data.info

| Name                | Type   | Description                                                                                                                                                                                   |
|---------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                  | int    | Unique identifier for the project in CC                                                                                                                                                       |
| bk_project_id       | string | Project ID                                                                                                                                                                                    |
| bk_project_name     | string | Project name                                                                                                                                                                                  |
| bk_project_code     | string | Project code (English name)                                                                                                                                                                   |
| bk_project_desc     | string | Project description                                                                                                                                                                           |
| bk_project_type     | enum   | Project type, optional values: "mobile_game" (mobile game), "pc_game" (PC game), "web_game" (web game), "platform_prod" (platform product), "support_prod" (support product), "other" (other) |
| bk_project_sec_lvl  | enum   | Security level, optional values: "public" (public), "private" (private), "classified" (classified)                                                                                            |
| bk_project_owner    | string | Project owner                                                                                                                                                                                 |
| bk_project_team     | array  | Team to which the project belongs                                                                                                                                                             |
| bk_project_icon     | string | Project icon                                                                                                                                                                                  |
| bk_status           | string | Project status, optional values: "enable" (enabled), "disabled" (disabled)                                                                                                                    |
| bk_supplier_account | string | Developer account                                                                                                                                                                             |
| create_time         | string | Creation time                                                                                                                                                                                 |
| last_time           | string | Update time                                                                                                                                                                                   |
