# 模型相关表

## cc_ObjDes

#### 作用

模型表，存放模型信息

#### 表结构

| 字段                   | 类型         | 描述               |
|----------------------|------------|------------------|
| _id                  | ObjectId   | 数据唯一ID           |
| position             | String     | 模型坐标，用于模型关系拓扑图中  |
| modifier             | String     | 修改者              |
| id                   | NumberLong | 自增id             |
| bk_obj_id            | String     | 模型唯一标识           |
| bk_obj_name          | String     | 模型名称             |
| description          | String     | 描述信息             |
| bk_classification_id | String     | 模型所属分组           |
| bk_obj_icon          | String     | 模型logo           |
| bk_ishidden          | Boolean    | 是否是隐藏模型如进程、服务模型等 |
| ispre                | Boolean    | 是否为系统预置          |
| bk_ispaused          | Boolean    | 是否停用             |
| creator              | String     | 创建者              |
| obj_sort_number      | NumberLong | 模型排序序号           |
| bk_supplier_account  | String     | 开发商ID            |
| create_time          | ISODate    | 创建时间             |
| last_time            | ISODate    | 最后更新时间           |

## cc_AsstDes

#### 作用

存放模型间的关联类型

#### 表结构

| 字段                  | 类型         | 描述              |
|---------------------|------------|-----------------|
| _id                 | ObjectId   | 数据唯一ID          |
| id                  | NumberLong | 自增ID            |
| src_des             | String     | 源模型到到关联目标模型的描述  |
| dest_des            | String     | 目标模型到源模型的描述     |
| direction           | String     | 方向：源指向目标、无方向、双向 |
| ispre               | String     | 是否为系统预置         |
| bk_asst_name        | String     | 关联类型名称          |
| bk_asst_id          | NumberLong | 关联类型唯一标识        |
| bk_supplier_account | String     | 开发商ID           |

## cc_ModelQuoteRelation

#### 作用

存放模型与引用类型关联关系信息

#### 表结构

| 字段                  | 类型       | 描述     |
|---------------------|----------|--------|
| _id                 | ObjectId | 数据唯一ID |
| dest_model          | String   | 引用模型   |
| src_model           | String   | 源模型    |
| bk_property_id      | String   | 模型属性id |
| type                | String   | 模型属性类型 |
| bk_supplier_account | String   | 开发商ID  |

## cc_ObjClassification

#### 作用

存放模型分组信息

#### 表结构

| 字段                     | 类型         | 描述     |
|------------------------|------------|--------|
| _id                    | ObjectId   | 数据唯一ID |
| bk_classification_id   | String     | 分组唯一标识 |
| bk_classification_name | String     | 分组名    |
| bk_classification_type | String     | 分组类型   |
| bk_classification_icon | String     | 分组logo |
| id                     | NumberLong | 自增id   |
| bk_supplier_account    | String     | 开发商ID  |

## cc_PropertyGroup

#### 作用

模型属性字段分组信息表

#### 表结构

| 字段                  | 类型         | 描述       |
|---------------------|------------|----------|
| _id                 | ObjectId   | 数据唯一ID   |
| bk_biz_id           | NumberLong | 业务id     |
| id                  | NumberLong | 自增id     |
| bk_obj_id           | String     | 分组所属模型id |
| ispre               | Boolean    | 是否为系统预置  |
| bk_group_id         | String     | 分组唯一标识   |
| bk_group_name       | String     | 分组名称     |
| bk_group_index      | NumberLong | 分组索引     |
| bk_isdefault        | Boolean    | 是否是默认分组  |
| is_collapse         | Boolean    | 是否折叠     |
| bk_supplier_account | String     | 开发商ID    |
| create_time         | ISODate    | 创建时间     |
| last_time           | ISODate    | 最后更新时间   |

## cc_ObjAsst

#### 作用

存放模型关联信息

#### 表结构

| 字段                  | 类型         | 描述                    |
|---------------------|------------|-----------------------|
| _id                 | ObjectId   | 数据唯一ID                |
| bk_asst_obj_id      | String     | 关联模型名称                |
| id                  | NumberLong | 自增id                  |
| bk_obj_id           | String     | 模型名称                  |
| bk_asst_id          | String     | 关联类型                  |
| bk_obj_asst_id      | String     | 关联标识                  |
| bk_obj_asst_name    | String     | 关联名称                  |
| ispre               | Boolean    | 是否为系统预置               |
| mapping             | String     | 源-目标约束，有n:n、1:n、1:1三种 |
| on_delete           | String     | 是否可删除                 |
| bk_supplier_account | String     | 开发商ID                 |

## cc_ObjAttDes

#### 作用

存放模型属性字段信息

#### 表结构

| 字段                  | 类型         | 描述                        |
|---------------------|------------|---------------------------|
| _id                 | ObjectId   | 数据唯一ID                    |
| bk_isapi            | Boolean    | 是否为API参数                  |
| creator             | String     | 创建人                       |
| bk_property_name    | String     | 属性字段名称                    |
| placeholder         | String     | 提示信息                      |
| editable            | Boolean    | 在实例中是否可编辑                 |
| bk_issystem         | Boolean    | 是否仅为系统内部使用                |
| option              | —          | 用户自定义内容，存储的内容及数据格式由字段类型决定 |
| id                  | NumberLong | 自增id                      |
| bk_obj_id           | String     | 模型id                      |
| bk_property_id      | String     | 属性字段id                    |
| bk_property_group   | String     | 属性字段分组                    |
| unit                | String     | 单位                        |
| isrequired          | Boolean    | 是否必填                      |
| isreadonly          | Boolean    | 是否只读                      |
| bk_property_index   | NumberLong | 属性字段排序号                   |
| ispre               | Boolean    | 是否为系统预置                   |
| bk_property_type    | String     | 属性字段类型                    |
| bk_biz_id           | NumberLong | 业务id，不为0则代表是业务自定义属性       |
| bk_template_id      | NumberLong | 字段模板id                    |
| default             | —          | 字段默认值，存储的内容及数据格式由字段类型决定   |
| description         | String     | 描述信息                      |
| ismultiple          | Boolean    | 是否可多选                     |
| isonly              | Boolean    | 是否唯一                      |
| bk_supplier_account | String     | 开发商ID                     |
| create_time         | ISODate    | 创建时间                      |
| last_time           | ISODate    | 最后更新时间                    |