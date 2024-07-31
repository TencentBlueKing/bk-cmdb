# 不归属任何一种分类的表

## cc_System

#### 作用

存放系统相关配置信息，如数据库版本信息、标签页和页底自定义信息等，无明确表字段结构

#### 表结构

- 数据库初始化数据

| 字段                  | 类型       | 描述      |
|---------------------|----------|---------|
| _id                 | ObjectId | 数据唯一ID  |
| type                | String   | 类型      |
| current_version     | String   | 当前版本    |
| distro              | String   | 发行版     |
| distro_version      | String   | 发行版版本   |
| init_version        | String   | 初始化版本   |
| init_distro_version | String   | 初始化发行版本 |

- 系统最大拓扑层级、标签、底注、字段验证规则等配置信息数据

| 字段          | 类型       | 描述     |
|-------------|----------|--------|
| _id         | ObjectId | 数据唯一ID |
| create_time | ISODate  | 创建时间   |
| last_time   | ISODate  | 最后更新时间 |
| config      | String   | 配置数据   |

- gse路由注册相关配置数据

| 字段        | 类型         | 描述                   |
|-----------|------------|----------------------|
| _id       | ObjectId   | 数据唯一ID               |
| host_snap | NumberLong | gse数据入库的stream_to_id |

## cc_DelArchive

#### 作用

用于归档被删除的数据

#### 表结构

| 字段          | 类型       | 描述     |
|-------------|----------|--------|
| _id         | ObjectId | 数据唯一ID |
| oid         | String   | 事件ID   |
| coll        | String   | 所操作的表  |
| detail      | Object   | 操作数据详情 |
| create_time | ISODate  | 创建时间   |
| last_time   | ISODate  | 最后更新时间 |

## cc_idgenerator

#### 作用

存放生成的唯一id信息

#### 表结构

| 字段          | 类型        | 描述        |
|-------------|-----------|-----------|
| _id         | String    | 唯一ID对应的数据 |
| SequenceID  | NumberInt | 生成的唯一ID   |
| create_time | ISODate   | 创建时间      |
| last_time   | ISODate   | 最后更新时间    |