# 服务模板功能相关表

## cc_ServiceTemplate

#### 作用

服务模板信息表

#### 表结构

| 字段                  | 类型         | 描述            |
|---------------------|------------|---------------|
| _id                 | ObjectId   | 数据唯一ID        |
| bk_biz_id           | NumberLong | 业务id          |
| id                  | NumberLong | 自增id          |
| name                | String     | 模板名称          |
| service_category_id | NumberLong | 服务类别id        |
| creator             | String     | 创建人           |
| modifier            | String     | 更新人           |
| host_apply_enabled  | Boolean    | 是否开启了主机属性自动应用 |
| bk_supplier_account | String     | 开发商ID         |
| create_time         | ISODate    | 创建时间          |
| last_time           | ISODate    | 最后更新时间        |

## cc_ServiceTemplateAttr

#### 作用

存放服务模板属性配置信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 自增id   |
| bk_biz_id           | NumberLong | 业务id   |
| service_template_id | NumberLong | 服务模板id |
| bk_attribute_id     | NumberLong | 属性字段id |
| bk_property_value   | String     | 属性字段值  |
| creator             | String     | 创建者    |
| modifier            | String     | 更新者    |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

## cc_ProcessTemplate

#### 作用

进程模板信息表

#### 表结构

| 字段                   | 类型         | 描述       |
|----------------------|------------|----------|
| _id                  | ObjectId   | 数据唯一ID   |
| id                   | NumberLong | 自增id     |
| creator              | String     | 创建人      |
| modifier             | String     | 更新人      |
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

| 字段       | 类型      | 描述   |
|----------|---------|------|
| enable   | Boolean | 是否启用 |
| ip       | String  | 绑定ip |
| port     | String  | 绑定端口 |
| protocol | String  | 通信方式 |