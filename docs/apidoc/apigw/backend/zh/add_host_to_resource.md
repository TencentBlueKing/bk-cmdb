### 描述

新增主机到资源池(权限：主机池主机创建权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述    |
|---------------------|--------|----|-------|
| bk_supplier_account | string | 否  | 开发商账号 |
| host_info           | dict   | 是  | 主机信息  |
| bk_biz_id           | int    | 否  | 业务ID  |

#### host_info

| 参数名称                 | 参数类型   | 必选 | 描述                     |
|----------------------|--------|----|------------------------|
| bk_host_innerip      | string | 是  | 主机内网ip                 |
| import_from          | string | 否  | 主机导入来源,以api方式导入为3      |
| bk_cloud_id          | int    | 否  | 管控区域ID，不填则添加到默认管控区域0   |
| bk_addressing        | string | 否  | 寻址方式，不填默认为静态寻址方式static |
| bk_host_name         | string | 否  | 主机名称                   |
| bk_asset_id          | string | 否  | 固资编号                   |
| bk_created_at        | string | 否  | 创建时间                   |
| bk_updated_at        | string | 否  | 更新时间                   |
| bk_created_by        | string | 否  | 创建人                    |
| bk_updated_by        | string | 否  | 更新人                    |
| bk_cloud_inst_id     | string | 否  | 云主机实例ID                |
| bk_cloud_vendor      | string | 否  | 云厂商                    |
| bk_cloud_host_status | string | 否  | 云主机状态                  |
| bk_comment           | string | 否  | 备注                     |
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
| operator             | string | 否  | 主要维护人                  |
| bk_bak_operator      | string | 否  | 备份维护人                  |

**注意：此处的输入参数仅对必填以及系统内置的参数做了说明，其余需要填写的参数取决于用户自己定义的主机属性字段，参数值的设置参考主机的属性字段配置
**

### 调用示例

```json
{
  "bk_biz_id": 3,
  "host_info": {
    "0": {
      "bk_host_innerip": "127.0.0.1",
      "bk_host_name": "host02",
      "bk_cloud_id": 0,
      "import_from": "3",
      "bk_addressing": "dynamic",
      "bk_asset_id": "udschdfhebv",
      "bk_created_at": "",
      "bk_updated_at": "",
      "bk_created_by": "admin",
      "bk_updated_by": "admin",
      "bk_cloud_inst_id": "1",
      "bk_cloud_vendor": "15",
      "bk_cloud_host_status": "2",
      "bk_comment": "canway-host",
      "bk_cpu": 8,
      "bk_cpu_architecture": "x86",
      "bk_cpu_module": "Intel(R) X87",
      "bk_disk": 195,
      "bk_host_outerip": "12.0.0.1",
      "bk_host_innerip_v6": "0000:0000:0000:0000:0000:0000:0000:0234",
      "bk_host_outerip_v6": "0000:0000:0000:0000:0000:0000:0000:0345",
      "bk_isp_name": "1",
      "bk_mac": "00:00:00:00:00:02",
      "bk_mem": 32155,
      "bk_os_bit": "64-bit",
      "bk_os_name": "linux redhat",
      "bk_os_type": "1",
      "bk_os_version": "7.8",
      "bk_outer_mac": "00:00:00:00:00:02",
      "bk_province_name": "110000",
      "bk_service_term": 6,
      "bk_sla": "1",
      "bk_sn": "abcdsd3252425",
      "bk_state": "测试中",
      "bk_state_name": "CN",
      "operator": "admin",
      "bk_bak_operator": "admin"
    }
  }
}
```

示例中host_info的"0"表示行数，可按顺序递增

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "success": [
      "0"
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
