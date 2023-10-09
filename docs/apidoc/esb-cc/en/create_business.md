### Functional description

New business

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_supplier_account | string     | No     | supplier account code |
| data           | dict    | Yes     | Data |

#### data

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| bk_biz_name       |  string  | Yes     | Business Name |
| bk_biz_maintainer |  string  | Yes     | the maintainers |
| bk_biz_productor  |  string  | Yes     | the productor |
| bk_biz_developer  |  string  | Yes     | the developer |
| bk_biz_tester     |  string  | Yes     | the tester |
| time_zone         |  string  | Yes     | time zone |
| language          |  string  | Yes     | language: "1" represent Chinese, "2" represent English |

**Note: The input parameters here only describe the required parameters and the built-in parameters of the system. The other parameters that need to be filled in depend on the attribute fields defined by the user.**

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "data": {
        "bk_biz_name": "cc_app_test",
        "bk_biz_maintainer": "admin",
        "bk_biz_productor": "admin",
        "bk_biz_developer": "admin",
        "bk_biz_tester": "admin",
        "time_zone": "Asia/Shanghai",
        "language": "1"
    }
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "bk_biz_developer": "admin",
        "bk_biz_id": 8852,
        "bk_biz_maintainer": "admin",
        "bk_biz_name": "cc_app_test",
        "bk_biz_productor": "admin",
        "bk_biz_tester": "admin",
        "bk_supplier_account": "0",
        "create_time": "2022-02-22T20:10:14.295+08:00",
        "default": 0,
        "language": "1",
        "last_time": "2022-02-22T20:10:14.295+08:00",
        "life_cycle": "2",
        "operator": null,
        "time_zone": "Asia/Shanghai"
    }
}
```
### Return Result Parameters Description

#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | Whether the request was successful or not. true:request successful; false request failed.
| code | int | The error code. 0 means success, >0 means failure error.
| message | string | The error message returned by the failed request.
| data | object | The data returned by the request.
| permission | object | Permission information |
| request_id | string | Request chain id |

#### data

| field | type | description |
| -----------|-----------|--------------|
| bk_biz_id | int | business id |
| bk_biz_name | string | business name |
| bk_biz_maintainer | string | Operations and maintenance personnel |
| bk_biz_productor | string | Product Personnel |
| bk_biz_developer | string | Developer |
| bk_biz_tester | string | Testers |
| time_zone | string | time zone |
| language | string | language, "1" for Chinese, "2" for English |
| bk_supplier_account | string | Developer account |
| create_time | string | Create time |
| last_time | string | update_time |
|default | int | Indicates business type |
| operator | string | Primary maintainer |
|life_cycle |string | business_cycle |
