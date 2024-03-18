# 不归属任何一种分类的表

## cc_System

#### 作用

存放系统相关配置信息，如数据库版本信息、标签页和页底自定义信息等，无明确表字段结构

#### 表结构

| 字段  | 类型       | 描述     |
|-----|----------|--------|
| _id | ObjectId | 数据唯一ID |
|     |          |        |

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

## cc_Subscription

#### 作用

事件订阅信息表，cmdb3.9及以下版本

#### 表结构

| 字段                  | 类型         | 描述        |
|---------------------|------------|-----------|
| _id                 | ObjectId   | 数据唯一ID    |
| subscription_id     | NumberLong | 事件订阅id    |
| subscription_name   | String     | 事件订阅名称    |
| system_name         | String     | 系统名称      |
| callback_url        | String     | 回调url     |
| confirm_mode        | String     | 确认模型      |
| confirm_pattern     | String     | 确认方式      |
| time_out            | NumberLong | 超时时间（秒）   |
| subscription_form   | String     | 订阅来源      |
| operator            | String     | 创建人       |
| bk_supplier_account | String     | 开发商ID，默认0 |
| create_time         | ISODate    | 创建时间      |
| last_time           | ISODate    | 最后更新时间    |