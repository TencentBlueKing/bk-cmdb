### Description

Add hosts to the specified resource pool based on the host list information (Permission: Host pool host creation
permission)

### Parameters

| Name      | Type         | Required | Description                |
|-----------|--------------|----------|----------------------------|
| host_info | object array | Yes      | Host information           |
| directory | int          | No       | Resource pool directory ID |

#### host_info

| Name                 | Type   | Required | Description                                                     |
|----------------------|--------|----------|-----------------------------------------------------------------|
| bk_host_innerip      | string | Yes      | Host's inner IP                                                 |
| bk_cloud_id          | int    | Yes      | Control area ID                                                 |
| bk_addressing        | string | No       | Addressing method, default is static addressing method "static" |
| bk_host_name         | string | No       | Host name, can also be other attributes                         |
| operator             | string | No       | Primary maintainer, can also be other attributes                |
| bk_comment           | string | No       | Remark, can also be other attributes                            |
| bk_cloud_vendor      | array  | No       | Cloud vendor                                                    |
| bk_cloud_inst_id     | array  | No       | Cloud host instance ID                                          |
| import_from          | string | No       | Host import source, 3 for API import                            |
| bk_asset_id          | string | No       | Fixed asset number                                              |
| bk_created_at        | string | No       | Creation time                                                   |
| bk_updated_at        | string | No       | Update time                                                     |
| bk_created_by        | string | No       | Creator                                                         |
| bk_updated_by        | string | No       | Updater                                                         |
| bk_cloud_host_status | string | No       | Cloud host status                                               |
| bk_cpu               | int    | No       | CPU logical cores                                               |
| bk_cpu_architecture  | string | No       | CPU architecture                                                |
| bk_cpu_module        | string | No       | CPU model                                                       |
| bk_disk              | int    | No       | Disk capacity (GB)                                              |
| bk_host_outerip      | string | No       | Host outer IP                                                   |
| bk_host_innerip_v6   | string | No       | Host inner IPv6                                                 |
| bk_host_outerip_v6   | string | No       | Host outer IPv6                                                 |
| bk_isp_name          | string | No       | ISP name                                                        |
| bk_mac               | string | No       | Host inner MAC address                                          |
| bk_mem               | int    | No       | Host memory capacity (MB)                                       |
| bk_os_bit            | string | No       | OS bit                                                          |
| bk_os_name           | string | No       | OS name                                                         |
| bk_os_type           | string | No       | OS type                                                         |
| bk_os_version        | string | No       | OS version                                                      |
| bk_outer_mac         | string | No       | Host outer MAC address                                          |
| bk_province_name     | string | No       | Province name                                                   |
| bk_service_term      | int    | No       | Warranty period                                                 |
| bk_sla               | string | No       | SLA level                                                       |
| bk_sn                | string | No       | Device SN                                                       |
| bk_state             | string | No       | Current state                                                   |
| bk_state_name        | string | No       | Country name                                                    |
| bk_bak_operator      | string | No       | Backup maintainer                                               |

**Note: The control area ID and inner IP fields are required fields. Other fields are attribute fields defined in the
host model. Only partial fields are shown here, please fill in other fields as needed.

### Request Example

```json
{
    "host_info": [
        {
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host1",
            "bk_cloud_id": 0,
            "operator": "admin",
            "bk_addressing": "dynamic",
            "bk_comment": "comment"
        },
        {
            "bk_host_innerip": "127.0.0.2",
            "bk_host_name": "host2",
            "operator": "admin",
            "bk_comment": "comment"
        }
    ],
    "directory": 1
}
```

### Response Example

```json
{
  "result": false,
  "code": 0,
  "message": "success",
  "data": {
      "success": [
          {
              "index": 0,
              "bk_host_id": 6
          }
      ],
      "error": [
          {
              "index": 1,
              "error_message": "'bk_cloud_id' unassigned"
          }
      ]
  },
  "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| data       | object | Data returned by the request                                      |
| permission | object | Permission information                                            |

#### data Field Description

| Name    | Type  | Description                                       |
|---------|-------|---------------------------------------------------|
| success | array | Array of successfully added host information      |
| error   | array | Array of host information that failed to be added |

#### success Field Description

| Name       | Type | Description                          |
|------------|------|--------------------------------------|
| index      | int  | Index of the successfully added host |
| bk_host_id | int  | ID of the successfully added host    |

#### error Field Description

| Name          | Type   | Description                               |
|---------------|--------|-------------------------------------------|
| index         | int    | Index of the host that failed to be added |
| error_message | string | Failure reason                            |
