### Description

Update Business Information (Permission: Business Edit Permission)

### Parameters

| Name                | Type   | Required | Description                                             |
|---------------------|--------|----------|---------------------------------------------------------|
| bk_supplier_account | string | Yes      | Developer account                                       |
| bk_biz_id           | int    | Yes      | Business ID                                             |
| bk_biz_name         | string | No       | Business name                                           |
| bk_biz_developer    | string | No       | Developer                                               |
| bk_biz_maintainer   | string | No       | Maintainer                                              |
| bk_biz_productor    | string | No       | Productor                                               |
| bk_biz_tester       | string | No       | Tester                                                  |
| operator            | string | No       | Operator                                                |
| life_cycle          | string | No       | Life cycle: Testing(1), Online(2, default), Shutdown(3) |
| language            | string | No       | Language, "1" for Chinese, "2" for English              |

**Note: The parameter here only explains the system-built editable parameters, and the rest of the parameters to be
filled depend on the user's own defined attribute fields.**

### Request Example

```json
{
  "bk_biz_name": "cc_app_test",
  "bk_biz_maintainer": "admin",
  "bk_biz_productor": "admin",
  "bk_biz_developer": "admin",
  "bk_biz_tester": "admin",
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
  "data": null
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
