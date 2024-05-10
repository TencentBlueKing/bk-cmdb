### 描述

根据指定动态分组规则查询获取数据 (版本：v3.9.6，权限：业务访问权限)

### 输入参数

| 参数名称            | 参数类型   | 必选 | 描述                                                            |
|-----------------|--------|----|---------------------------------------------------------------|
| bk_biz_id       | int    | 是  | 业务ID                                                          |
| id              | string | 是  | 动态分组主键ID                                                      |
| fields          | array  | 是  | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输,目标资源不具备指定的字段时该字段将被忽略 |
| disable_counter | bool   | 否  | 是否返回总记录条数，默认返回                                                |
| page            | object | 是  | 分页设置                                                          |

#### page

| 参数名称  | 参数类型   | 必选 | 描述               |
|-------|--------|----|------------------|
| start | int    | 是  | 记录开始位置           |
| limit | int    | 是  | 每页限制条数,最大200     |
| sort  | string | 否  | 检索排序， 默认按照创建时间排序 |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "disable_counter": true,
    "id": "XXXXXXXX",
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_host_name"
    ],
    "page":{
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
      "count": 1,
      "info": [
        {
          "bk_cloud_id": 0,
          "bk_host_id": 2,
          "bk_host_innerip": "127.0.0.1",
          "bk_host_name": "host12"
        },
        {
          "bk_cloud_id": 0,
          "bk_host_id": 9,
          "bk_host_innerip": "127.0.0.2",
          "bk_host_name": "host111"
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

| 参数名称  | 参数类型  | 描述                                                                 |
|-------|-------|--------------------------------------------------------------------|
| count | int   | 当前规则能匹配到的总记录条数（用于调用者进行预分页，实际单次请求返回数量以及数据是否全部拉取完毕以JSON Array解析数量为准） |
| info  | array | dict数组，主机实际数据, 当动态分组为host查询时返回host自身属性信息,当动态分组为set查询时返回set信息       |

#### data.info -- 分组目标为主机

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

#### data.info -- 分组目标为集群

| 参数名称                 | 参数类型   | 描述                         |
|----------------------|--------|----------------------------|
| bk_set_name          | string | 集群名称                       |
| default              | int    | 0-普通集群，1-内置模块集合，默认为0       |
| bk_biz_id            | int    | 业务id                       |
| bk_capacity          | int    | 设计容量                       |
| bk_parent_id         | int    | 父节点的ID                     |
| bk_set_id            | int    | 集群id                       |
| bk_service_status    | string | 服务状态:1/2(1:开放,2:关闭)        |
| bk_set_desc          | string | 集群描述                       |
| bk_set_env           | string | 环境类型：1/2/3(1:测试,2:体验,3:正式) |
| create_time          | string | 创建时间                       |
| last_time            | string | 更新时间                       |
| bk_supplier_account  | string | 开发商账号                      |
| description          | string | 数据的描述信息                    |
| set_template_version | array  | 集群模板的当前版本                  |
| set_template_id      | int    | 集群模板ID                     |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**
