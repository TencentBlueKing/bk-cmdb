### Description

Create a new business (Permission: Business Creation Permission)

### Parameters

| Name                | Type   | Required | Description                                             |
|---------------------|--------|----------|---------------------------------------------------------|
| bk_supplier_account | string | Yes      | Developer account                                       |
| bk_biz_name         | string | Yes      | Business name                                           |
| bk_biz_maintainer   | string | Yes      | Operation and maintenance personnel                     |
| bk_biz_productor    | string | No       | Product personnel                                       |
| bk_biz_developer    | string | No       | Development personnel                                   |
| bk_biz_tester       | string | No       | Test personnel                                          |
| operator            | string | No       | Operator                                                |
| bk_created_at       | string | No       | Creation time                                           |
| bk_updated_at       | string | No       | Update time                                             |
| bk_created_by       | string | No       | Creator                                                 |
| bk_updated_by       | string | No       | Last updater                                            |
| life_cycle          | string | No       | Lifecycle: Testing (1), Live (2, default), Shutdown (3) |
| time_zone           | string | No       | Time zone                                               |
| language            | string | No       | Language, "1" for Chinese, "2" for English              |

**Note: The input parameters here only explain the required and system-built parameters. The rest of the parameters to
be filled in depend on the user's own defined attribute fields.**

### Request Example

```json
{
  "bk_biz_name": "main-a1",
  "bk_biz_maintainer": "admin",
  "bk_biz_productor": "admin",
  "bk_biz_developer": "admin",
  "bk_biz_tester": "admin",
  "time_zone": "Asia/Shanghai",
  "bk_created_at": "",
  "bk_updated_at": "",
  "bk_created_by": "admin",
  "bk_updated_by": "admin",
  "language": "1",
  "operator": "admin",
  "life_cycle": "2"
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "bk_biz_developer": "admin",
    "bk_biz_id": 5,
    "bk_biz_maintainer": "admin",
    "bk_biz_name": "main-a1",
    "bk_biz_productor": "admin",
    "bk_biz_tester": "admin",
    "bk_created_at": "2023-11-14T16:51:02.168+08:00",
    "bk_created_by": "admin",
    "bk_supplier_account": "0",
    "bk_updated_at": "2023-11-14T16:51:02.168+08:00",
    "create_time": "2023-11-14T16:51:02.168+08:00",
    "default": 0,
    "language": "1",
    "last_time": "2023-11-14T16:51:02.168+08:00",
    "life_cycle": "2",
    "operator": "admin",
    "time_zone": "Asia/Shanghai"
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |

#### data

| Name                | Type   | Description                                |
|---------------------|--------|--------------------------------------------|
| bk_biz_id           | int    | Business ID                                |
| bk_biz_name         | string | Business name                              |
| bk_biz_maintainer   | string | Operation and maintenance personnel        |
| bk_biz_productor    | string | Product personnel                          |
| bk_biz_developer    | string | Development personnel                      |
| bk_biz_tester       | string | Test personnel                             |
| time_zone           | string | Time zone                                  |
| language            | string | Language, "1" for Chinese, "2" for English |
| bk_supplier_account | string | Developer account                          |
| create_time         | string | Creation time                              |
| last_time           | string | Update time                                |
| default             | int    | Business type                              |
| operator            | string | Main maintainer                            |
| life_cycle          | string | Business lifecycle                         |
| bk_created_at       | string | Creation time                              |
| bk_updated_at       | string | Update time                                |
| bk_created_by       | string | Creator                                    |
