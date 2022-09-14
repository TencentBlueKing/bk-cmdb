### 功能描述

 根据审计ID获取详细信息

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| id     |  array    |是      | 审计id数组,一次限制最大传200个                                             |  

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id":[95,118]
}

```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### 返回结果参数说明

#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误   |
| message | string | 请求失败返回的错误信息                   |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                          |

#### data

| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
|    id |      int  |    审计ID  |
|   audit_type  |     string   |   操作审计类型   |
|   bk_supplier_account  |    string    | 开发商账号     |
|   user  |      string  |    操作人  |
|   resource_type  |    string    |   资源类型   |
|  action   |    string    |    操作类型  |
|    operate_from |    string    |   来源平台   |
|  operation_detail   |     object     |  操作细节    |
| operation_time    |     string   |    操作时间  |
|  bk_biz_id   |       int |    业务id  |
| resource_id    |     int   |    资源id  |
|   resource_name  |     string   |  资源名称    |
|   rid  |     string   |  请求链id    |

#### operation_detail
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
|    details |      object  |    详细数据   |
|    bk_obj_id |      string  |    模型类型   |

#### details
| 字段      | 类型      | 描述         |
|-----------|-----------|--------------|
|    pre_data |      object  |    之前数据   |
|   cur_data  |     object   |   现在数据   |
|   update_fields  |     object   |   更新的字段   |
