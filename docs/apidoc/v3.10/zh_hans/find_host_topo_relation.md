### 功能描述

获取主机与拓扑的关系

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
 bk_biz_id| int| 是|业务ID|
| bk_set_ids|array | 否| 集群ID列表，最多200条|
| bk_module_ids|array | 否| 模块ID列表，最多500条| 
| bk_host_ids|array | 否| 主机ID列表，最多500条| 
| page| object| 是|分页信息|

#### page 字段说明

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
|start|int|否|获取数据偏移位置|
|limit|int|是|过去数据条数限制，建议 为200|

### 请求参数示例

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "page":{
        "start":0,
        "limit":10
    },
    "bk_biz_id":2,
    "bk_set_ids": [1, 2],
    "bk_module_ids": [23, 24],
    "bk_host_ids": [25, 26]
}
```

### 返回结果示例

```python

{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "count": 10,
    "info": [
        {
        "bk_biz_id": 3,
        "bk_host_id": 5,
        "bk_module_id": 54,
        "bk_set_id": 10,
        "bk_supplier_account": "0"
        },
        .....
    ]
}
```



### 返回结果参数说明


#### data 字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| count| int| 记录条数 |
| info| object array |  业务下主机与集群，模块，集群的数据详情列表 |


#### data.info 字段说明：
| 名称  | 类型  | 说明 |
|---|---|---|
| bk_biz_id | int | 业务ID |
| bk_set_id | int | 集群ID |
| bk_module_id | int | 模块ID |
| bk_host_id | int | 主机ID |
