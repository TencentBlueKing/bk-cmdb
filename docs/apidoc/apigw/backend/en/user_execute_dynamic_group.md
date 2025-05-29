### Description

Query and retrieve data based on the specified dynamic grouping rule (Version: v3.9.6, Permission: Business access
permission)

### Parameters

| Name            | Type   | Required | Description                                                                                                                                                                                                                        |
|-----------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id       | int    | Yes      | Business ID                                                                                                                                                                                                                        |
| id              | string | Yes      | Dynamic grouping primary key ID                                                                                                                                                                                                    |
| fields          | array  | Yes      | Host attribute list, controls which fields are returned in the host, accelerates interface requests and reduces network traffic transmission. If the target resource does not have the specified fields, the field will be ignored |
| disable_counter | bool   | No       | Whether to return the total number of records, default is to return                                                                                                                                                                |
| page            | object | Yes      | Pagination settings                                                                                                                                                                                                                |

#### page

| Name  | Type   | Required | Description                                            |
|-------|--------|----------|--------------------------------------------------------|
| start | int    | Yes      | Record start position                                  |
| limit | int    | Yes      | Number of records per page, maximum is 200             |
| sort  | string | No       | Retrieval sorting, default is to sort by creation time |

### Request Example

```json
{
    "bk_biz_id": 1,
    "disable_counter": true,
    "id": "XXXXXXXX",
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_host_name"
    ],
    "page":{
        "start": 0,
        "limit": 10
    }
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
      "count": 1,
      "info": [
        {
          "bk_cloud_id": 0,
          "bk_host_id": 2,
          "bk_host_innerip": "127.0.0.1",
          "bk_host_name": "host12"
        },
        {
          "bk_cloud_id": 0,
          "bk_host_id": 9,
          "bk_host_innerip": "127.0.0.2",
          "bk_host_name": "host111"
        }
      ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description                                                                                                                                                                                                                      |
|-------|-------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| count | int   | Total number of records that the current rule can match (used for callers to pre-page, the actual number of records returned in a single request and whether the data is all pulled is based on the JSON Array parsing quantity) |
| info  | array | Dict array, actual data of hosts, returns host's own attribute information when dynamic grouping is host query, returns set information when dynamic grouping is set query                                                       |

#### data.info -- Group target is host

| Name                 | Type   | Description                          |
|----------------------|--------|--------------------------------------|
| bk_host_name         | string | Host name                            |
| bk_host_innerip      | string | Inner IP                             |
| bk_host_id           | int    | Host ID                              |
| bk_cloud_id          | int    | Control area                         |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number                   |
| bk_cloud_inst_id     | string | Cloud host instance ID               |
| bk_cloud_vendor      | string | Cloud vendor                         |
| bk_cloud_host_status | string | Cloud host status                    |
| bk_comment           | string | Comment                              |
| bk_cpu               | int    | CPU logical cores                    |
| bk_cpu_architecture  | string | CPU architecture                     |
| bk_cpu_module        | string | CPU model                            |
| bk_disk              | int    | Disk capacity (GB)                   |
| bk_host_outerip      | string | Host outer IP                        |
| bk_host_innerip_v6   | string | Host inner IPv6                      |
| bk_host_outerip_v6   | string | Host outer IPv6                      |
| bk_isp_name          | string | ISP name                             |
| bk_mac               | string | Host inner MAC address               |
| bk_mem               | int    | Host memory capacity (MB)            |
| bk_os_bit            | string | Operating system bit number          |
| bk_os_name           | string | Operating system name                |
| bk_os_type           | string | Operating system type                |
| bk_os_version        | string | Operating system version             |
| bk_outer_mac         | string | Host outer MAC address               |
| bk_province_name     | string | Province name                        |
| bk_service_term      | int    | Warranty period                      |
| bk_sla               | string | SLA level                            |
| bk_sn                | string | Device SN                            |
| bk_state             | string | Current status                       |
| bk_state_name        | string | Country                              |
| operator             | string | Main maintainer                      |
| bk_bak_operator      | string | Backup maintainer                    |

#### data.info -- Group target is set

| Name                 | Type   | Description                                                 |
|----------------------|--------|-------------------------------------------------------------|
| bk_set_name          | string | Set name                                                    |
| default              | int    | 0-ordinary set, 1-built-in module set, default is 0         |
| bk_biz_id            | int    | Business id                                                 |
| bk_capacity          | int    | Design capacity                                             |
| bk_parent_id         | int    | Parent node ID                                              |
| bk_set_id            | int    | Set id                                                      |
| bk_service_status    | string | Service status: 1/2 (1: open, 2: closed)                    |
| bk_set_desc          | string | Set description                                             |
| bk_set_env           | string | Environment type: 1/2/3 (1: test, 2: experience, 3: formal) |
| create_time          | string | Creation time                                               |
| last_time            | string | Update time                                                 |
| bk_supplier_account  | string | Developer account                                           |
| description          | string | Data description                                            |
| set_template_version | array  | Set template's current version                              |
| set_template_id      | int    | Set template ID                                             |

**Note: The return value here only explains the system's built-in attribute fields, and the rest of the return values
depend on the user's self-defined attribute fields**
