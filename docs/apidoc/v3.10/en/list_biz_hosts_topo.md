### Functional description

list hosts under business with their topology information, can filter by host and set and module condition

### Request Parameters


#### General Parameters

| Field         | Type   | Required | Description                                                  |
| ------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_app_code   | string | Yes      | APP ID                                                       |
| bk_app_secret | string | Yes      | APP Secret(APP TOKEN), which can be got via BlueKing Developer Center -&gt; Click APP ID -&gt; Basic Info |
| bk_token      | string | No       | Current user login token, bk_token or bk_username must be valid, bk_token can be got by Cookie |
| bk_username   | string | No       | Current user username, APP in the white list, can use this field to specify the current user |
| fields        | array  | Yes      | host property list, the specified host property feilds will be returned <br>it can speed up the request and reduce the network payload |

#### Interface Parameters

| Field                  | Type  | Required | Description                                                  |
| ---------------------- | ----- | -------- | ------------------------------------------------------------ |
| bk_biz_id              | int   | Yes      | Business ID                                                  |
| page                   | dict  | Yes      | paging search condition                                      |
| set_property_filter    | dict  | No       | set property filter                                          |
| module_property_filter | dict  | No       | module property filter                                       |
| host_property_filter   | dict  | No       | host property filter                                         |
| fields                 | array | Yes      | host property list, the specified host property feilds will be returned <br>it can speed up the request and reduce the network payload |


#### set_property_filter
set property filter is a combined of atom filter rule, combine operator could be `AND` or `OR`, nested up to 2 levels。
atom rule has three fields: `field`, `operator`, `value`

| Field    | Type   | Required | Description |
| -------- | ------ | -------- | ----------- |
| field    | string | Yes      | field       |
| operator | string | No       | operator    |
| value    | -      | No       | value       |

reference: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>


#### module_property_filter
module property filter is a combined of atom filter rule, combine operator could be `AND` or `OR`, nested up to 2 levels。
atom rule has three fields: `field`, `operator`, `value`

| Field    | Type   | Required | Description |
| -------- | ------ | -------- | ----------- |
| field    | string | Yes      | field       |
| operator | string | No       | operator    |
| value    | -      | No       | value       |

reference: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>


#### host_property_filter
host property filter is a combined of atom filter rule, combine operator could be `AND` or `OR`, nested up to 2 levels。
atom rule has three fields: `field`, `operator`, `value`

| Field    | Type   | Required | Description |
| -------- | ------ | -------- | ----------- |
| field    | string | Yes      | field       |
| operator | string | No       | operator    |
| value    | -      | No       | value       |

reference: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### page

| Field | Type | Required | Description                      |
| ----- | ---- | -------- | -------------------------------- |
| start | int  | Yes      | start record                     |
| limit | int  | Yes      | page limit, maximum value is 500 |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "page": {
        "start": 0,
        "limit": 10
    },
    "bk_biz_id": 3,
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "set_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_set_name",
                "operator": "not_equal",
                "value": "test"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_set_id",
                        "operator": "in",
                        "value": [
                            1,
                            2,
                            3
                        ]
                    },
                    {
                        "field": "bk_service_status",
                        "operator": "equal",
                        "value": "1"
                    }
                ]
            }
        ]
    },
    "module_property_filter": {
        "condition": "OR",
        "rules": [
            {
                "field": "bk_module_name",
                "operator": "equal",
                "value": "test"
            },
            {
                "condition": "AND",
                "rules": [
                    {
                        "field": "bk_module_id",
                        "operator": "not_in",
                        "value": [
                            1,
                            2,
                            3
                        ]
                    },
                    {
                        "field": "bk_module_type",
                        "operator": "equal",
                        "value": "1"
                    }
                ]
            }
        ]
    },
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
    "count": 3,
    "info": [
      {
        "host": {
          "bk_cloud_id": 0,
          "bk_host_id": 1,
          "bk_host_innerip": "192.168.15.18",
          "bk_mac": "",
          "bk_os_type": null
        },
        "topo": [
          {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "module": [
              {
                "bk_module_id": 56,
                "bk_module_name": "m1"
              }
            ]
          }
        ]
      },
      {
        "host": {
          "bk_cloud_id": 0,
          "bk_host_id": 2,
          "bk_host_innerip": "192.168.15.4",
          "bk_mac": "",
          "bk_os_type": null
        },
        "topo": [
          {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "module": [
              {
                "bk_module_id": 56,
                "bk_module_name": "m1"
              }
            ]
          }
        ]
      },
      {
        "host": {
          "bk_cloud_id": 0,
          "bk_host_id": 3,
          "bk_host_innerip": "192.168.15.12",
          "bk_mac": "",
          "bk_os_type": null
        },
        "topo": [
          {
            "bk_set_id": 10,
            "bk_set_name": "idle pool",
            "module": [
              {
                "bk_module_id": 54,
                "bk_module_name": "idle host"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### Return Result Parameters Description

#### data

| Field | Type  | Description                        |
| ----- | ----- | ---------------------------------- |
| count | int   | the num of record                  |
| info  | array | host data and topology information |

#### data.info
| Field | Type  | Description       |
| ----- | ----- | ----------------- |
| host  | dict  | the num of record |
| topo  | array | host data         |

#### data.info.host
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

#### data.info.topo
| Field       | Type   | Description          |
| ----------- | ------ | -------------------- |
| bk_set_id   | int    | host set ID          |
| bk_set_name | string | host set name        |
| module      | array  | host set module info |

#### data.info.topo.module
| Field          | Type   | Description |
| -------------- | ------ | ----------- |
| bk_module_id   | int    | module ID   |
| bk_module_name | string | module name |
