# 实例相关表

## cc_ObjectBase_0_pub_{模型id}

#### 作用

模型实例表，CMDB3.10及3.10以上版本使用

#### 表结构

| 字段                  | 类型         | 描述            |
|---------------------|------------|---------------|
| _id                 | ObjectId   | 数据唯一ID        |
| bk_admin_ip         | String     | 管理IP（模型自定义字段） |
| bk_biz_status       | String     | 运营状态          |
| bk_model            | String     | 设备型号          |
| bk_asset_id         | String     | 固资编号          |
| bk_inst_name        | String     | 实例名称          |
| bk_os_detail        | String     | 操作系统详情        |
| bk_inst_id          | NumberLong | 实例id          |
| bk_detail           | String     | 详细描述          |
| bk_func             | String     | 用途            |
| bk_obj_id           | String     | 模型id          |
| bk_operator         | String     | 维护人           |
| bk_vendor           | String     | 厂商            |
| bk_supplier_account | String     | 开发商ID         |
| create_time         | ISODate    | 创建时间          |
| last_time           | ISODate    | 最后更新时间        |

**注意**：此处以”防火墙“模型实例举例，模型实例字段取决于用户在模型中定义的属性字段

## cc_InstAsst_0_pub_{模型id}

#### 作用

存放实例关联数据，CMDB3.10及3.10以上版本

#### 表结构

| 字段                  | 类型         | 描述       |
|---------------------|------------|----------|
| _id                 | ObjectId   | 数据唯一ID   |
| id                  | NumberLong | 自增id     |
| bk_inst_id          | NumberLong | 实例id     |
| bk_obj_id           | String     | 模型id     |
| bk_asst_inst_id     | NumberLong | 关联实例id   |
| bk_asst_obj_id      | String     | 关联模型id   |
| bk_supplier_account | String     | 开发商ID    |
| bk_obj_asst_id      | String     | 模型关联关系id |
| bk_asst_id          | String     | 关联类型id   |