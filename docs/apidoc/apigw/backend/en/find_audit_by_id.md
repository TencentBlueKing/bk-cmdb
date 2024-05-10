### Description

Get detailed information based on the audit ID (Permission: Operation audit query permission)

### Parameters

| Name | Type  | Required | Description                         |
|------|-------|----------|-------------------------------------|
| id   | array | Yes      | Audit ID array, up to 200 at a time |

### Request Example

```python
{
    "id":[95,118]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": [
        {
            "id": 95,
            "audit_type": "host",
            "bk_supplier_account": "0",
            "user": "admin",
            "resource_type": "host",
            "action": "update",
            "operate_from": "user",
            "operation_detail": {
                "details": {
                    "pre_data": {
                        "bk_asset_id": "",
                        "bk_bak_operator": "",
                        "bk_cloud_host_status": null,
                        "bk_cloud_id": 0,
                        "bk_cloud_inst_id": "",
                        "bk_cloud_vendor": null,
                        "bk_comment": "",
                        "bk_cpu": null,
                        "bk_cpu_mhz": null,
                        "bk_cpu_module": "",
                        "bk_disk": null,
                        "bk_host_id": 4,
                        "bk_host_innerip": "1.1.1.1",
                        "bk_host_name": "",
                        "bk_host_outerip": "",
                        "bk_isp_name": null,
                        "bk_mac": "",
                        "bk_mem": null,
                        "bk_os_bit": "",
                        "bk_os_name": "",
                        "bk_os_type": null,
                        "bk_os_version": "",
                        "bk_outer_mac": "",
                        "bk_province_name": null,
                        "bk_service_term": null,
                        "bk_sla": null,
                        "bk_sn": "",
                        "bk_state": null,
                        "bk_state_name": null,
                        "bk_supplier_account": "0",
                        "create_time": "2020-10-21T18:49:14.342+08:00",
                        "docker_client_version": "",
                        "docker_server_version": "",
                        "import_from": "1",
                        "last_time": "2020-10-21T18:49:14.342+08:00",
                        "operator": "",
                        "test1": null,
                        "test2": null
                    },
                    "cur_data": null,
                    "update_fields": {
                        "test1": "2020-10-01 00:00:00"
                    }
                },
                "bk_obj_id": "host"
            },
            "operation_time": "2020-10-21 18:49:48",
            "bk_biz_id": 1,
            "resource_id": 4,
            "resource_name": "1.1.1.1"
        },
        {
            "id": 118,
            "audit_type": "host",
            "bk_supplier_account": "0",
            "user": "admin",
            "resource_type": "host",
            "action": "delete",
            "operate_from": "user",
            "operation_detail": {
                "details": {
                    "pre_data": {
                        "bk_asset_id": "",
                        "bk_bak_operator": "",
                        "bk_cloud_host_status": null,
                        "bk_cloud_id": 0,
                        "bk_cloud_inst_id": "",
                        "bk_cloud_vendor": null,
                        "bk_comment": "",
                        "bk_cpu": null,
                        "bk_cpu_mhz": null,
                        "bk_cpu_module": "",
                        "bk_disk": null,
                        "bk_host_id": 4,
                        "bk_host_innerip": "1.1.1.1",
                        "bk_host_name": "",
                        "bk_host_outerip": "",
                        "bk_isp_name": null,
                        "bk_mac": "",
                        "bk_mem": null,
                        "bk_os_bit": "",
                        "bk_os_name": "",
                        "bk_os_type": null,
                        "bk_os_version": "",
                        "bk_outer_mac": "",
                        "bk_province_name": null,
                        "bk_service_term": null,
                        "bk_sla": null,
                        "bk_sn": "",
                        "bk_state": null,
                        "bk_state_name": null,
                        "bk_supplier_account": "0",
                        "create_time": "2020-10-21T18:49:14.342+08:00",
                        "docker_client_version": "",
                        "docker_server_version": "",
                        "import_from": "1",
                        "last_time": "2020-10-21T18:49:48.569+08:00",
                        "operator": "",
                        "test1": "2020-10-01T00:00:00+08:00",
                        "test2": null
                    },
                    "cur_data": null,
                    "update_fields": null
                },
                "bk_obj_id": "host"
            },
            "operation_time": "2020-10-21 19:02:30",
            "bk_biz_id": 1,
            "resource_id": 4,
            "resource_name": "1.1.1.1"
        }
    ]
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

| Name                | Type   | Description          |
|---------------------|--------|----------------------|
| id                  | int    | Audit ID             |
| audit_type          | string | Operation audit type |
| bk_supplier_account | string | Developer account    |
| user                | string | Operator             |
| resource_type       | string | Resource type        |
| action              | string | Operation type       |
| operate_from        | string | Source platform      |
| operation_detail    | object | Operation details    |
| operation_time      | string | Operation time       |
| bk_biz_id           | int    | Business id          |
| resource_id         | int    | Resource id          |
| resource_name       | string | Resource name        |
| rid                 | string | Request chain ID     |

#### operation_detail

| Name      | Type   | Description   |
|-----------|--------|---------------|
| details   | object | Detailed data |
| bk_obj_id | string | Model type    |

#### details

| Name          | Type   | Description    |
|---------------|--------|----------------|
| pre_data      | object | Previous data  |
| cur_data      | object | Current data   |
| update_fields | object | Updated fields |
