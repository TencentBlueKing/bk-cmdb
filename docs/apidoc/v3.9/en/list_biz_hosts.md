### Functional description

list hosts under special business

### Request Parameters


#### General Parameters

| Field         | Type   | Required | Description                                                  |
| ------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_app_code   | string | Yes      | APP ID                                                       |
| bk_app_secret | string | Yes      | APP Secret(APP TOKEN), which can be got via BlueKing Developer Center -&gt; Click APP ID -&gt; Basic Info |
| bk_token      | string | No       | Current user login token, bk_token or bk_username must be valid, bk_token can be got by Cookie |
| bk_username   | string | No       | Current user username, APP in the white list, can use this field to specify the current user |

#### Interface Parameters

| Field                | Type   | Required | Description                                                  |
| -------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_supplier_account  | string | No       | supplier account code                                        |
| bk_biz_id            | int    | Yes      | Business ID                                                  |
| set_cond             | array  | No       | cluster search condition **NOTICE: set_cond and bk_set_ids can't be used at the same time** |
| bk_set_ids           | array  | No       | cluster ID, length must be less than 200 **NOTICE: set_cond and bk_set_ids can't be used at the same time** |
| module_cond          | array  | No       | module search condition **NOTICE: module_cond and bk_module_ids can't be used at the same time** |
| bk_module_ids        | array  | No       | module ID, length must be less than 500   **NOTICE: module_cond and bk_module_ids can't be used at the same time** |
| page                 | dict   | Yes      | search condition                                             |
| host_property_filter | dict   | No       | host property filter                                         |
| fields               | array  | Yes      | host property list, the specified host property feilds will be returned <br>it can speed up the request and reduce the network payload |


#### host_property_filter
host property filter is a combined of atom filter rule, combine operator could be `AND` or `OR`, nested up to 2 levels。
atom rule has three fields: `field`, `operator`, `value`

| Field    | Type   | Required | Description |
| -------- | ------ | -------- | ----------- |
| field    | string | Yes      | field       |
| operator | string | No       | operator    |
| value    | -      | No       | value       |

reference: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### set_cond

| Field    | Type   | Required | Description            |
| -------- | ------ | -------- | ---------------------- |
| field    | string | Yes      | Value of set field     |
| operator | string | Yes      | value : $regex $eq $ne |
| value    | string | Yes      | Value of set field     |

#### module_cond

| Field    | Type   | Required | Description            |
| -------- | ------ | -------- | ---------------------- |
| field    | string | Yes      | Value of module field  |
| operator | string | Yes      | value : $regex $eq $ne |
| value    | string | Yes      | Value of module field  |


#### page

| Field | Type   | Required | Description            |
| ----- | ------ | -------- | ---------------------- |
| start | int    | Yes      | start record           |
| limit | int    | Yes      | page limit, max is 500 |
| sort  | string | No       | the field for sort     |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_host_id"
    },
    "set_cond": [
        {
            "field": "bk_set_name",
            "operator": "$eq",
            "value": "set1"
        }
    ],
    "bk_biz_id": 3,
    "bk_module_ids": [54,56],
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "127.0.0.1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_os_type",
                        "operator": "not_in",
                        "value": [
                            "3"
                        ]
                    },
                    {
                        "field": "bk_cloud_id",
                        "operator": "equal",
                        "value": 0
                    }
                ]
            }
        ]
    }
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "count": 2,
    "info": [
      {
        "bk_cloud_id": 0,
        "bk_host_id": 1,
        "bk_host_innerip": "192.168.15.18",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 2,
        "bk_host_innerip": "192.168.15.4",
        "bk_mac": "",
        "bk_os_type": null
      }
    ]
  }
}
```

### Return Result Parameters Description

#### data

| Field | Type  | Description       |
| ----- | ----- | ----------------- |
| count | int   | the num of record |
| info  | array | host data         |

#### data.info
| Field            | Type   | Description       |
| ---------------- | ------ | ----------------- |
| bk_isp_name      | string | telecom operators |
| bk_sn            | string | device SN         |
| operator         | string | maintainer        |
| bk_outer_mac     | string | outer MAC         |
| bk_state_name    | string | country           |
| bk_province_name | string | province          |
| import_from      | string | import from       |
| bk_sla           | string | SLA level         |
| bk_service_term  | int    | warranty          |
| bk_os_type       | string | os type           |
| bk_os_version    | string | os version        |
| bk_os_bit        | int    | os bits           |
| bk_mem           | string | memory capacity   |
| bk_mac           | string | mac address       |
| bk_host_outerip  | string | outer ip          |
| bk_host_name     | string | hostname          |
| bk_host_innerip  | string | inner ip          |
| bk_host_id       | int    | host id           |
| bk_disk          | int    | disk capacity     |
| bk_cpu_module    | string | CPU module        |
| bk_cpu_mhz       | int    | CPU hz            |
| bk_cpu           | int    | CPU count         |
| bk_comment       | string | comment           |
| bk_cloud_id      | int    | cloud area id     |
| bk_bak_operator  | string | backup maintainer |
| bk_asset_id      | string | device id         |
