### 描述

根据主机列表信息添加主机到指定id的资源池(权限：主机池主机创建权限)

### 输入参数

| 参数名称      | 参数类型         | 必选 | 描述      |
|-----------|--------------|----|---------|
| host_info | object array | 是  | 主机信息    |
| directory | int          | 否  | 资源池目录ID |

#### host_info

| 参数名称                 | 参数类型   | 必选 | 描述                     |
|----------------------|--------|----|------------------------|
| bk_host_innerip      | string | 是  | 主机内网ip                 |
| bk_cloud_id          | int    | 是  | 管控区域id                 |
| bk_addressing        | string | 否  | 寻址方式，不填默认为静态寻址方式static |
| bk_host_name         | string | 否  | 主机名，也可以为其它属性           |
| operator             | string | 否  | 主要维护人，也可以为其它属性         |
| bk_comment           | string | 否  | 备注，也可以为其它属性            |
| bk_cloud_vendor      | array  | 否  | 云厂商                    |
| bk_cloud_inst_id     | array  | 否  | 云主机实例ID                |
| import_from          | string | 否  | 主机导入来源,以api方式导入为3      |
| bk_asset_id          | string | 否  | 固资编号                   |
| bk_created_at        | string | 否  | 创建时间                   |
| bk_updated_at        | string | 否  | 更新时间                   |
| bk_created_by        | string | 否  | 创建人                    |
| bk_updated_by        | string | 否  | 更新人                    |
| bk_cloud_host_status | string | 否  | 云主机状态                  |
| bk_cpu               | int    | 否  | CPU逻辑核心数               |
| bk_cpu_architecture  | string | 否  | CPU架构                  |
| bk_cpu_module        | string | 否  | CPU型号                  |
| bk_disk              | int    | 否  | 磁盘容量（GB）               |
| bk_host_outerip      | string | 否  | 主机外网IP                 |
| bk_host_innerip_v6   | string | 否  | 主机内网IPv6               |
| bk_host_outerip_v6   | string | 否  | 主机外网IPv6               |
| bk_isp_name          | string | 否  | 所属运营商                  |
| bk_mac               | string | 否  | 主机内网MAC地址              |
| bk_mem               | int    | 否  | 主机名内存容量（MB）            |
| bk_os_bit            | string | 否  | 操作系统位数                 |
| bk_os_name           | string | 否  | 操作系统名称                 |
| bk_os_type           | string | 否  | 操作系统类型                 |
| bk_os_version        | string | 否  | 操作系统版本                 |
| bk_outer_mac         | string | 否  | 主机外网MAC地址              |
| bk_province_name     | string | 否  | 所在省份                   |
| bk_service_term      | int    | 否  | 质保年限                   |
| bk_sla               | string | 否  | SLA级别                  |
| bk_sn                | string | 否  | 设备SN                   |
| bk_state             | string | 否  | 当前状态                   |
| bk_state_name        | string | 否  | 所在国家                   |
| bk_bak_operator      | string | 否  | 备份维护人                  |

**注意：上述参数中管控区域ID和内网IP字段为必填字段，其它字段为主机模型中定义的属性字段。在此仅展示部分字段示例，其它字段请按需填写

### 调用示例

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

### 响应示例

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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data 字段说明

| 参数名称    | 参数类型  | 描述          |
|---------|-------|-------------|
| success | array | 添加成功的主机信息数组 |
| error   | array | 添加失败的主机信息数组 |

#### success 字段说明

| 参数名称       | 参数类型 | 描述        |
|------------|------|-----------|
| index      | int  | 添加成功的主机下标 |
| bk_host_id | int  | 添加成功的主机ID |

#### error 字段说明

| 参数名称          | 参数类型   | 描述        |
|---------------|--------|-----------|
| index         | int    | 添加失败的主机下标 |
| error_message | string | 失败原因      |
