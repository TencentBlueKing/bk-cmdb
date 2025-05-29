### 描述

更新主机属性(权限：对于已分配到业务下的主机需要业务主机编辑权限，对于主机池主机需要主机池主机编辑权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述           |
|---------------------|--------|----|--------------|
| bk_host_id          | string | 是  | 主机ID，多个以逗号分隔 |
| bk_host_name        | string | 否  | 主机名          |
| bk_comment          | string | 否  | 备注           |
| bk_cpu              | int    | 否  | CPU逻辑核心数     |
| bk_cpu_architecture | string | 否  | CPU架构        |
| bk_cpu_module       | string | 否  | CPU型号        |
| bk_disk             | int    | 否  | 磁盘容量（GB）     |
| bk_host_outerip     | string | 否  | 主机外网IP       |
| bk_host_outerip_v6  | string | 否  | 主机外网IPv6     |
| bk_isp_name         | string | 否  | 所属运营商        |
| bk_mac              | string | 否  | 主机内网MAC地址    |
| bk_mem              | int    | 否  | 主机名内存容量（MB）  |
| bk_os_bit           | string | 否  | 操作系统位数       |
| bk_os_name          | string | 否  | 操作系统名称       |
| bk_os_type          | string | 否  | 操作系统类型       |
| bk_os_version       | string | 否  | 操作系统版本       |
| bk_outer_mac        | string | 否  | 主机外网MAC地址    |
| bk_province_name    | string | 否  | 所在省份         |
| bk_sla              | string | 否  | SLA级别        |
| bk_sn               | string | 否  | 设备SN         |
| bk_state            | string | 否  | 当前状态         |
| bk_state_name       | string | 否  | 所在国家         |
| operator            | string | 否  | 主要维护人        |
| bk_bak_operator     | string | 否  | 备份维护人        |

**注意：此处仅对系统内置可编辑的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段**

### 调用示例

```json
{
  "bk_host_id": "1,2,3",
  "bk_host_name": "test",
  "bk_comment": "canway-host-101",
  "bk_cpu": 16,
  "bk_cpu_architecture": "arm",
  "bk_cpu_module": "Intel(R) 2.00GHz",
  "bk_disk": 120,
  "bk_host_outerip": "12.0.0.3",
  "bk_host_outerip_v6": "0000:0000:0000:0000:0000:0000:0000:0248",
  "bk_isp_name": "3",
  "bk_mac": "00:00:00:00:00:56",
  "bk_mem": 36666,
  "bk_os_bit": "32-bit",
  "bk_os_name": "ubuntu",
  "bk_os_type": "4",
  "bk_os_version": "7.9.1",
  "bk_outer_mac": "00:00:00:00:00:56",
  "bk_province_name": "440000",
  "bk_sla": "2",
  "bk_sn": "abcd3252425",
  "bk_state": "备用机",
  "bk_state_name": "BE",
  "operator": "admin",
  "bk_bak_operator": "admin"
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
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
