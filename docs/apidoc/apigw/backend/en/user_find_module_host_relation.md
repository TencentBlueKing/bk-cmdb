### Description

Query the relationship between hosts and modules based on the module ID (version: v3.8.7, permission: Business access
permission)

### Parameters

| Name          | Type   | Required | Description                                                                     |
|---------------|--------|----------|---------------------------------------------------------------------------------|
| bk_biz_id     | int    | Yes      | Business ID                                                                     |
| bk_module_ids | array  | Yes      | Module ID array, up to 200                                                      |
| module_fields | array  | Yes      | Module attribute list, control which fields to return in the module information |
| host_fields   | array  | Yes      | Host attribute list, control which fields to return in the host information     |
| page          | object | Yes      | Pagination parameters                                                           |

#### page

| Name  | Type | Required | Description                                 |
|-------|------|----------|---------------------------------------------|
| start | int  | No       | Record start position, default is 0         |
| limit | int  | Yes      | Number of records per page, maximum is 1000 |

**Note: The relationship between hosts and modules under a module may be returned multiple times, and the paging method
is based on the host ID sorting.**

### Request Example

```json
{
    "bk_biz_id": 1,
    "bk_module_ids": [
        1,
        2,
        3
    ],
    "module_fields": [
        "bk_module_id",
        "bk_module_name"
    ],
    "host_fields": [
        "bk_host_innerip",
        "bk_host_id"
    ],
    "page": {
        "start": 0,
        "limit": 500
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
    "count": 2,
    "relation": [
      {
        "host": {
          "bk_host_id": 1,
          "bk_host_innerip": "127.0.0.1",
        },
        "modules": [
          {
            "bk_module_id": 1,
            "bk_module_name": "m1",
          },
          {
            "bk_module_id": 2,
            "bk_module_name": "m2",
          }
        ]
      },
      {
        "host": {
          "bk_host_id": 2,
          "bk_host_innerip": "127.0.0.2",
        },
        "modules": [
          {
            "bk_module_id": 3,
            "bk_module_name": "m3",
          }
        ]
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

#### Explanation of data field:

| Name     | Type  | Description                      |
|----------|-------|----------------------------------|
| count    | int   | Number of records                |
| relation | array | Actual data of hosts and modules |

#### Explanation of data.relation field:

| Name    | Type   | Description                                             |
|---------|--------|---------------------------------------------------------|
| host    | object | Host data                                               |
| modules | array  | Information about the modules to which the host belongs |

#### Explanation of data.relation.host field:

| Name                 | Type   | Description                               |
|----------------------|--------|-------------------------------------------|
| bk_host_name         | string | Host name                                 |
| bk_host_innerip      | string | Private IP of the host                    |
| bk_host_id           | int    | Host ID                                   |
| bk_cloud_id          | int    | Control area                              |
| import_from          | string | Host import source, imported via API is 3 |
| bk_asset_id          | string | Fixed asset number                        |
| bk_cloud_inst_id     | string | Cloud host instance ID                    |
| bk_cloud_vendor      | string | Cloud vendor                              |
| bk_cloud_host_status | string | Cloud host status                         |
| bk_comment           | string | Comment                                   |
| bk_cpu               | int    | Number of logical cores of CPU            |
| bk_cpu_architecture  | string | CPU architecture                          |
| bk_cpu_module        | string | CPU model                                 |
| bk_disk              | int    | Disk capacity (GB)                        |
| bk_host_outerip      | string | Outer IP of the host                      |
| bk_host_innerip_v6   | string | Inner IPv6 of the host                    |
| bk_host_outerip_v6   | string | Outer IPv6 of the host                    |
| bk_isp_name          | string | ISP name                                  |
| bk_mac               | string | Inner MAC address of the host             |
| bk_mem               | int    | Memory capacity of the host (MB)          |
| bk_os_bit            | string | OS bit                                    |
| bk_os_name           | string | OS name                                   |
| bk_os_type           | string | OS type                                   |
| bk_os_version        | string | OS version                                |
| bk_outer_mac         | string | Outer MAC address of the host             |
| bk_province_name     | string | Province name                             |
| bk_service_term      | int    | Warranty period                           |
| bk_sla               | string | SLA level                                 |
| bk_sn                | string | Device SN                                 |
| bk_state             | string | Current status                            |
| bk_state_name        | string | Country of the host                       |
| operator             | string | Main maintainer                           |
| bk_bak_operator      | string | Backup maintainer                         |

Explanation of data.relation.modules field:

| Name                | Type    | Description                                                |
|---------------------|---------|------------------------------------------------------------|
| bk_module_id        | int     | Module ID                                                  |
| bk_module_name      | string  | Module name                                                |
| default             | int     | Indicates the module type                                  |
| create_time         | string  | Creation time                                              |
| bk_set_id           | int     | Cluster ID                                                 |
| bk_bak_operator     | string  | Backup maintainer                                          |
| bk_biz_id           | int     | Business ID                                                |
| bk_module_type      | string  | Module type                                                |
| bk_parent_id        | int     | Parent node ID                                             |
| bk_supplier_account | string  | Developer account                                          |
| last_time           | string  | Update time                                                |
| host_apply_enabled  | bool    | Whether to enable automatic application of host properties |
| operator            | string  | Main maintainer                                            |
| service_category_id | integer | Service category ID                                        |
| service_template_id | int     | Service template ID                                        |
| set_template_id     | int     | Cluster template ID                                        |
| bk_created_at       | string  | Creation time                                              |
| bk_updated_at       | string  | Update time                                                |
| bk_created_by       | string  | Creator                                                    |

**Note: The returned value here only explains the built-in property fields. The rest of the returned values depend on
the user's own defined property fields.**
