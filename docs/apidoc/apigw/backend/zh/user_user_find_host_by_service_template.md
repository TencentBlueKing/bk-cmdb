### 描述

获取服务模板下的主机 (版本：v3.8.6，权限：业务访问权限)

### 输入参数

| 参数名称                    | 参数类型   | 必选 | 描述                       |
|-------------------------|--------|----|--------------------------|
| bk_biz_id               | int    | 是  | 业务ID                     |
| bk_service_template_ids | array  | 是  | 服务模板ID列表，最多可填500个        |
| bk_module_ids           | array  | 否  | 模块ID列表, 最多可填500个         |
| fields                  | array  | 是  | 主机属性列表，控制返回结果的模块信息里有哪些字段 |
| page                    | object | 是  | 分页信息                     |

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数,最大500 |

### 调用示例

```json
{
    "bk_biz_id": 5,
    "bk_service_template_ids": [
        48,
        49
    ],
    "bk_module_ids": [
        65,
        68
    ],
    "fields": [
        "bk_host_id",
        "bk_cloud_id"
    ],
    "page": {
        "start": 0,
        "limit": 10
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "count": 6,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 3
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 4
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 7
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 8
            }
        ]
    }
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称  | 参数类型  | 描述     |
|-------|-------|--------|
| count | int   | 记录条数   |
| info  | array | 主机实际数据 |

#### data.info

| 参数名称                 | 参数类型   | 描述                |
|----------------------|--------|-------------------|
| bk_host_name         | string | 主机名               |
| bk_host_innerip      | string | 内网IP              |
| bk_host_id           | int    | 主机ID              |
| bk_cloud_id          | int    | 管控区域              |
| import_from          | string | 主机导入来源,以api方式导入为3 |
| bk_asset_id          | string | 固资编号              |
| bk_cloud_inst_id     | string | 云主机实例ID           |
| bk_cloud_vendor      | string | 云厂商               |
| bk_cloud_host_status | string | 云主机状态             |
| bk_comment           | string | 备注                |
| bk_cpu               | int    | CPU逻辑核心数          |
| bk_cpu_architecture  | string | CPU架构             |
| bk_cpu_module        | string | CPU型号             |
| bk_disk              | int    | 磁盘容量（GB）          |
| bk_host_outerip      | string | 主机外网IP            |
| bk_host_innerip_v6   | string | 主机内网IPv6          |
| bk_host_outerip_v6   | string | 主机外网IPv6          |
| bk_isp_name          | string | 所属运营商             |
| bk_mac               | string | 主机内网MAC地址         |
| bk_mem               | int    | 主机名内存容量（MB）       |
| bk_os_bit            | string | 操作系统位数            |
| bk_os_name           | string | 操作系统名称            |
| bk_os_type           | string | 操作系统类型            |
| bk_os_version        | string | 操作系统版本            |
| bk_outer_mac         | string | 主机外网MAC地址         |
| bk_province_name     | string | 所在省份              |
| bk_service_term      | int    | 质保年限              |
| bk_sla               | string | SLA级别             |
| bk_sn                | string | 设备SN              |
| bk_state             | string | 当前状态              |
| bk_state_name        | string | 所在国家              |
| operator             | string | 主要维护人             |
| bk_bak_operator      | string | 备份维护人             |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
