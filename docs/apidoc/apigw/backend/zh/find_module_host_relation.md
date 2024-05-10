### 描述

根据模块ID查询主机和模块的关系(版本：v3.8.7，权限：业务访问权限)

### 输入参数

| 参数名称          | 参数类型   | 必选 | 描述                     |
|---------------|--------|----|------------------------|
| bk_biz_id     | int    | 是  | 业务ID                   |
| bk_module_ids | array  | 是  | 模块ID数组，最多200条          |
| module_fields | array  | 是  | 模块属性列表，控制返回结果的模块里有哪些字段 |
| host_fields   | array  | 是  | 主机属性列表，控制返回结果的主机里有哪些字段 |
| page          | object | 是  | 分页参数                   |

#### page

| 参数名称  | 参数类型 | 必选 | 描述            |
|-------|------|----|---------------|
| start | int  | 否  | 记录开始位置,默认值0   |
| limit | int  | 是  | 每页限制条数,最大1000 |

**注: 一个模块下的主机关系可能会拆分多次返回，分页方式是按主机ID排序进行分页。**

### 调用示例

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

### 响应示例

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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data 字段说明：

| 参数名称     | 参数类型  | 描述        |
|----------|-------|-----------|
| count    | int   | 记录条数      |
| relation | array | 主机和模块实际数据 |

#### data.relation 字段说明:

| 参数名称    | 参数类型   | 描述        |
|---------|--------|-----------|
| host    | object | 主机数据      |
| modules | array  | 主机所属的模块信息 |

#### data.relation.host 字段说明:

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

data.relation.modules 字段说明:

| 参数名称                | 参数类型    | 描述           |
|---------------------|---------|--------------|
| bk_module_id        | int     | 模块id         |
| bk_module_name      | string  | 模块名称         |
| default             | int     | 表示模块类型       |
| create_time         | string  | 创建时间         |
| bk_set_id           | int     | 集群id         |
| bk_bak_operator     | string  | 备份维护人        |
| bk_biz_id           | int     | 业务id         |
| bk_module_type      | string  | 模块类型         |
| bk_parent_id        | int     | 父节点的ID       |
| bk_supplier_account | string  | 开发商账号        |
| last_time           | string  | 更新时间         |
| host_apply_enabled  | bool    | 是否启用主机属性自动应用 |
| operator            | string  | 主要维护人        |
| service_category_id | integer | 服务分类ID       |
| service_template_id | int     | 服务模版ID       |
| set_template_id     | int     | 集群模板ID       |
| bk_created_at       | string  | 创建时间         |
| bk_updated_at       | string  | 更新时间         |
| bk_created_by       | string  | 创建人          |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
