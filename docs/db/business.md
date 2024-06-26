# 业务下相关资源表

## cc_ServiceCategory

#### 作用

服务分类信息表

#### 表结构

| 字段                  | 类型         | 描述      |
|---------------------|------------|---------|
| _id                 | ObjectId   | 数据唯一ID  |
| id                  | NumberLong | 分类id    |
| name                | String     | 分类名称    |
| bk_root_id          | NumberLong | 根节点id   |
| bk_parent_id        | NumberLong | 父节点id   |
| is_built_in         | Boolean    | 是否为系统内置 |
| bk_biz_id           | NumberLong | 业务id    |
| bk_supplier_account | String     | 开发商ID   |

## cc_HostFavourite

#### 作用

存放收藏主机搜索条件信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| bk_biz_id           | NumberLong | 业务id   |
| id                  | NumberLong | 自增id   |
| info                | String     | ip查询条件 |
| name                | String     | 收藏的名称  |
| count               | NumberLong | 数量     |
| user                | String     | 创建者    |
| type                | String     | 收藏类型   |
| query_params        | String     | 通用查询条件 |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

## cc_ServiceInstance

#### 作用

服务实例信息表

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 自增id   |
| bk_biz_id           | NumberLong | 业务id   |
| name                | String     | 服务实例名  |
| labels              | Object     | 标签     |
| service_template_id | NumberLong | 服务模板id |
| bk_host_id          | NumberLong | 主机id   |
| bk_module_id        | NumberLong | 模块id   |
| creator             | String     | 创建人    |
| modifier            | String     | 更新人    |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

#### labels 字段结构示例

| 字段  | 类型     | 描述                      |
|-----|--------|-------------------------|
| env | String | 标签名称，由用户自定义，此处 env 仅为例子 |

## cc_Process

#### 作用

进程实例信息表

#### 表结构

| 字段                   | 类型         | 描述       |
|----------------------|------------|----------|
| _id                  | ObjectId   | 数据唯一ID   |
| service_instance_id  | NumberLong | 服务实例id   |
| start_cmd            | String     | 启动命令     |
| user                 | String     | 启动用户     |
| bk_start_param_regex | String     | 启动参数匹配规则 |
| restart_cmd          | String     | 重启命令     |
| face_stop_cmd        | String     | 强制停止命令   |
| work_path            | String     | 工作路径     |
| description          | String     | 备注       |
| priority             | NumberLong | 启动优先级    |
| bk_process_name      | String     | 进程别名     |
| bk_start_check_secs  | NumberLong | 启动等待时长   |
| bk_func_name         | String     | 进程名称     |
| reload_cmd           | String     | 进程重载命令   |
| pid_file             | String     | PID文件路径  |
| auto_start           | Boolean    | 是否自动拉起   |
| bk_biz_id            | NumberLong | 业务ID     |
| proc_num             | NumberLong | 启动数量     |
| stop_cmd             | String     | 停止命令     |
| bk_process_id        | NumberLong | 进程id     |
| timeout              | NumberLong | 操作超时时长   |
| bind_info            | Array      | 绑定信息     |
| bk_supplier_account  | String     | 开发商ID    |
| create_time          | ISODate    | 创建时间     |
| last_time            | ISODate    | 最后更新时间   |

#### bind_info 字段结构示例

| 字段              | 类型         | 描述         |
|-----------------|------------|------------|
| enable          | Boolean    | 是否启用       |
| template_row_id | NumberLong | 在进程模板中的行id |
| ip              | String     | 绑定ip       |
| port            | String     | 绑定端口       |
| protocol        | String     | 通信方式       |

## cc_ProcessInstanceRelation

#### 作用

存放进程实例关联关系信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| bk_biz_id           | NumberLong | 业务id   |
| bk_process_id       | NumberLong | 进程id   |
| service_instance_id | NumberLong | 服务实例id |
| process_template_id | NumberLong | 进程模板id |
| bk_host_id          | NumberLong | 主机id   |
| bk_supplier_account | String     | 开发商ID  |