### 1. 描述

定义数据库表需要业务索引


### 2. 功能

1. 索引的定义
2. 索引处理的辅助函数， eg: 对比和查找


### 3. 索引定义

#### 3.1 association.go 
分表后，实例关联关系表的索引

#### 3.2 instance.go 

分表后，通用模型实例表的索引

#### 3.3 collections 目录

用来定义数据库中表的索引。每一个go文件表示一个数据库表中的索引的定义。

#### 3.4 注册表中的索引

```
func init() {
	// 先注册未规范化的索引，如果索引出现冲突旧，删除未规范化的索引
	registerIndexes("数据库的表名", 索引定义)
}
```


#### 3.4 索引定义 

注意索引的名字
索引定义 = []types.Index{}

20210321 之前存在的的表，都有索引的定义。 

- comm 前缀的定义
- deprecated 前缀的定义 


#### 3.4.1 comm 前缀的索引定义

规范化后的索引名字
注意： 索引的名字一定要以CCLogicUniqueIdxNamePrefix或common.CCLogicIndexNamePrefix为前缀

同步索引到数据库db 表中索引的功能， 会拉取db表中已经定义的索引。 
筛选出索引名字以CCLogicUniqueIdxNamePrefix或common.CCLogicIndexNamePrefix前缀的索引， 筛选出来的db表中的索引与cc代码定义索引按照名字做差异对比。（
- 筛选出来的db表中的索引（按前缀筛选）比cc代码定义（comm 前缀的索引定义）多的索引进行删除。原因该索引已经被删除，但是在db中存在）
- 筛选出来的db表中的索引（按前缀筛选）比cc代码定义comm 前缀的索引定义缺少的索引进行创建。（原因该索引在cc中存在，但是在db中不存在）
- 索引名字在db 表中和cc定义（comm 前缀的索引定义）中同时存在，对比索引的定义是否相同，如果不相同，删除db表中的索引， 使用cc定义的索引新建 （原因索引在cc和db定义不同，需要重新建立索引）


#### 3.4.1 deprecated 前缀的索引定义

规范化前的索引定义

注意：
1. mongodb 索引没有没有修改的概念，如果要调整deprecated 前缀的索引定义中的索引，请删除后，在comm 前缀的索引定义中新加
2. 未规范化的索引，无法通过规则来筛选出那些是cc定义索引或者是用户手动新加的索引，所以再deprecatedindexname.go文件中记录了所有需要管理索引名字。
deprecatedindexname.go 不允许修改。 



同步索引到数据库db 表中索引的功能， 会拉取db表中已经定义的索引与cc中定义的索引（deprecated 前缀的索引定义）进行差异对比。

- 找出在db表中不存在cc中定义（deprecated 前缀的索引定义）的索引 进行创建。
- 找出在db表中存在cc中定义（deprecated 前缀的索引定义）的索引进行索引定义内容的对比，如果不想等，删除db中的索引，在新建索引。

- 找到表在deprecatedindexname.go文件中定义需要待管理索引的名字集合与cc中定义（deprecated 前缀的索引定义）的索引进行差异对比，找出只在deprecatedindexname.go存在的索引名字。 去数据库中删除索引。 



举例：

2021-03-24日xx表中的索引
|cc deprecated前缀的索引定义|            需要cc管理索引 |       db 中的索引||
|:--|--:|:--|:--|
| idx_id |           idx_id |              idx_id||
| idx_test |         idx_test |            |db中不存在需要新建|
| bk_obj_id_1|       bk_obj_id_1|          bk_obj_id_1||
| idx_inst_id|       bk_obj_id_1|          bk_inst_id_1|cc 定义索引bk_inst_id，但是数据库中是inst_id,需要删除数据库中的数据，在创建对应的索引|

2021-03-25日xx表中的索引，由于形态调整需要删除idx_test索引。
|cc deprecated前缀的索引定义|            需要cc管理索引 |       db 中的索引||
|:--|--:|:--|:--|
| idx_id |           idx_id |              idx_id||
|  |  idx_test |     idx_test       |idx_test索引需要管理，定义已经不存在，但是db还有该索引，需要删除|
| bk_obj_id_1|       bk_obj_id_1|          bk_obj_id_1||
| idx_inst_id|       bk_obj_id_1|          bk_inst_id_1|cc 定义索引bk_inst_id，但是数据库中是inst_id,需要删除数据库中的数据，在创建对应的索引|









