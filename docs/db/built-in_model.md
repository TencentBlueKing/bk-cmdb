# 内置模型相关表

## cc_BizSetBase

#### 作用

存放业务集信息

#### 表结构

| 字段                  | 类型         | 描述                 |
|---------------------|------------|--------------------|
| _id                 | ObjectId   | 数据唯一ID             |
| bk_biz_set_id       | NumberLong | 业务集id              |
| bk_biz_set_name     | String     | 业务集名称              |
| bk_biz_set_desc     | String     | 业务集描述              |
| bk_biz_maintainer   | String     | 维护者                |
| bk_scope            | Object     | 业务集范围              |
| default             | NumberLong | 是否是默认业务集，1代表是，0代表否 |
| bk_supplier_account | String     | 开发商ID              |
| create_time         | ISODate    | 创建时间               |
| last_time           | ISODate    | 最后更新时间             |

**注意**：此处仅对业务集内置的模型字段做说明，业务集表结构字段取决于用户在业务集模型中定义的属性字段

#### bk_scope 字段结构示例

| 字段        | 类型      | 描述         |
|-----------|---------|------------|
| match_all | Boolean | 是否包含所有业务   |
| filter    | Object  | 所包含业务的条件规则 |

#### bk_scope.filter 字段结构示例

| 字段        | 类型     | 描述        |
|-----------|--------|-----------|
| condition | String | 规则操作符     |
| rules     | Array  | 过滤节点的范围规则 |

#### bk_scope.filter.rules 字段结构示例

过滤规则为三元组 `field`, `operator`, `value`

| 字段       | 类型     | 描述                          |
|----------|--------|-----------------------------|
| field    | string | 字段名                         |
| operator | string | 操作符                         |
| value    | -      | 操作数,不同的operator对应不同的value格式 |

## cc_HostBase

#### 作用

主机信息表

#### 表结构

| 字段                       | 类型           | 描述           |
|--------------------------|--------------|--------------|
| _id                      | ObjectId     | 数据唯一ID       |
| bk_os_name               | String       | 操作系统名称       |
| bk_addressing            | String       | 寻址方式         |
| import_from              | String       | 录入方式         |
| bk_sla                   | String       | SLA级别        |
| bk_cpu_module            | String       | CPU型号        |
| bk_host_innerip_v6       | String Array | 内网IPv6       |
| bk_agent_id              | String       | GSE Agent ID |
| bk_cloud_host_identifier | Boolean      | 云主机标识        |
| bk_os_bit                | String       | 操作系统位数       |
| bk_sn                    | String       | 设备SN         |
| bk_host_outerip          | String Array | 外网IP         |
| bk_service_term          | NumberLong   | 质保年限         |
| bk_state_name            | String       | 所在国家         |
| bk_mac                   | String       | 内网MAC地址      |
| bk_cloud_inst_id         | String       | 云主机实例ID      |
| bk_cpu_architecture      | String       | CPU架构        |
| bk_os_version            | String       | 操作系统版本       |
| operator                 | String Array | 主要维护人        |
| bk_os_type               | String       | 操作系统类型       |
| bk_host_id               | NumberLong   | 主机id         |
| bk_cloud_id              | NumberLong   | 管控区域id       |
| bk_comment               | String       | 备注           |
| bk_cloud_host_status     | NumberLong   | 云主机状态        |
| bk_bak_operator          | String Array | 备份维护人        |
| bk_province_name         | NumberLong   | 所在省份         |
| bk_isp_name              | String       | 所属运营商        |
| bk_disk                  | NumberLong   | 磁盘容量(GB)     |
| bk_outer_mac             | String       | 外网MAC        |
| bk_state                 | String       | 当前状态         |
| bk_host_outerip_v6       | String Array | 外网IPv6       |
| bk_host_name             | String       | 主机名称         |
| bk_asset_id              | String       | 固资编号         |
| bk_cloud_vendor          | String       | 云厂商          |
| bk_host_innerip          | String Array | 内网IP         |
| bk_cpu                   | NumberLong   | CPU逻辑核心数     |
| bk_mem                   | NumberLong   | 内存容量(MB)     |
| bk_supplier_account      | String       | 开发商ID        |
| create_time              | ISODate      | 创建时间         |
| last_time                | ISODate      | 最后更新时间       |

**注意**：此处仅对主机内置的模型字段做说明，主机表结构字段取决于用户在主机模型中定义的属性字段

## cc_ProjectBase

#### 作用

项目表

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 自增ID   |
| bk_project_id       | String     | 项目ID   |
| bk_project_code     | String     | 项目英文名  |
| bk_project_type     | String     | 项目类型   |
| bk_project_owner    | String     | 项目负责人  |
| bk_project_name     | String     | 项目名称   |
| bk_project_desc     | String     | 项目描述   |
| bk_project_sec_lvl  | String     | 保密级别   |
| bk_status           | String     | 项目状态   |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

**注意**：此处仅对项目内置的模型字段做说明，项目表结构字段取决于用户在项目模型中定义的属性字段