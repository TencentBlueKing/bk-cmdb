### 描述

新增云主机到业务的空闲机模块 (云主机管理专用接口, 版本: v3.10.19+, 权限：业务主机编辑权限)

### 输入参数

| 参数名称      | 参数类型  | 必选 | 描述                                  |
|-----------|-------|----|-------------------------------------|
| bk_biz_id | int   | 是  | 业务ID                                |
| host_info | array | 是  | 新增的云主机信息，数组长度最多为200，一批主机仅可同时成功或同时失败 |

#### host_info

主机信息，其中管控区域ID、内网IP、云厂商、云主机实例ID字段为必填字段，其它字段为主机模型中定义的属性字段。在此仅展示部分字段示例，其它字段请按需填写

| 参数名称                 | 参数类型   | 必选 | 描述                        |
|----------------------|--------|----|---------------------------|
| bk_cloud_id          | int    | 是  | 管控区域ID                    |
| bk_host_innerip      | string | 是  | IPv4格式的主机内网IP，多个IP之间用逗号分隔 |
| bk_cloud_vendor      | array  | 是  | 云厂商                       |
| bk_cloud_inst_id     | array  | 是  | 云主机实例ID                   |
| bk_addressing        | string | 否  | 寻址方式，云主机的寻址方式都是static     |
| bk_host_name         | string | 否  | 主机名，也可以为其它属性              |
| operator             | string | 否  | 主要维护人，也可以为其它属性            |
| bk_comment           | string | 否  | 备注，也可以为其它属性               |
| import_from          | string | 否  | 主机导入来源,以api方式导入为3         |
| bk_asset_id          | string | 否  | 固资编号                      |
| bk_created_at        | string | 否  | 创建时间                      |
| bk_updated_at        | string | 否  | 更新时间                      |
| bk_created_by        | string | 否  | 创建人                       |
| bk_updated_by        | string | 否  | 更新人                       |
| bk_cloud_host_status | string | 否  | 云主机状态                     |
| bk_cpu               | int    | 否  | CPU逻辑核心数                  |
| bk_cpu_architecture  | string | 否  | CPU架构                     |
| bk_cpu_module        | string | 否  | CPU型号                     |
| bk_disk              | int    | 否  | 磁盘容量（GB）                  |
| bk_host_outerip      | string | 否  | 主机外网IP                    |
| bk_host_innerip_v6   | string | 否  | 主机内网IPv6                  |
| bk_host_outerip_v6   | string | 否  | 主机外网IPv6                  |
| bk_isp_name          | string | 否  | 所属运营商                     |
| bk_mac               | string | 否  | 主机内网MAC地址                 |
| bk_mem               | int    | 否  | 主机名内存容量（MB）               |
| bk_os_bit            | string | 否  | 操作系统位数                    |
| bk_os_name           | string | 否  | 操作系统名称                    |
| bk_os_type           | string | 否  | 操作系统类型                    |
| bk_os_version        | string | 否  | 操作系统版本                    |
| bk_outer_mac         | string | 否  | 主机外网MAC地址                 |
| bk_province_name     | string | 否  | 所在省份                      |
| bk_service_term      | int    | 否  | 质保年限                      |
| bk_sla               | string | 否  | SLA级别                     |
| bk_sn                | string | 否  | 设备SN                      |
| bk_state             | string | 否  | 当前状态                      |
| bk_state_name        | string | 否  | 所在国家                      |
| bk_bak_operator      | string | 否  | 备份维护人                     |

### 调用示例

```json
{
  "bk_biz_id": 123,
  "host_info": [
    {
      "bk_cloud_id": 0,
      "bk_host_innerip": "127.0.0.1",
      "bk_cloud_vendor": "2",
      "bk_cloud_inst_id": "45515",
      "bk_host_name": "host1",
      "operator": "admin",
      "bk_comment": "comment"
    },
    {
      "bk_cloud_id": 0,
      "bk_host_innerip": "127.0.0.2",
      "bk_cloud_vendor": "2",
      "bk_cloud_inst_id": "45656",
      "bk_host_name": "host2",
      "operator": "admin",
      "bk_comment": "comment"
    }
  ]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "ids": [
      1,
      2
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称 | 参数类型  | 描述           |
|------|-------|--------------|
| ids  | array | 创建成功的主机的ID数组 |
