### Functional description

 Get detailed information based on audit ID

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field | Type | Required | Description |
|-----------|------------|--------|------------|
| id |int array |Yes | Audit id array

### Request Parameters Example
```python
{
    "id":[95,118]
}
```

### Return Result Example

```json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
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

### Return Result Parameters Description

#### data

| Field | Type | Description |
|-----------|-----------|--------------|
|  | array | Operation audit record information |

# single audit details

| Field | Type | Description |
|-----------|-----------|--------------|
| id | int | Audit ID |
| audit_type | string | Operation audit type |
| bk_supplier_account | string | Developer account |
| user | string | Operator |
| resource_type | string | Resource type |
| action | string | Action type |
| operate_from | string | Source platform |
| operation_detail | object | Operation details |
| operation_time | string | Operation time |
| bk_biz_id | int | business id |
| resource_id | int | resource id |
| resource_name | string | Resource name |
