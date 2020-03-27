# auth extensions

## why need it?
auth 模块提供了到iam的统一访问接口，如果直接运用auth接口对接到各个场景层或APIServer，将会嵌入非常多的代码到其中，比如生成auth接口需要的数据格式，其中最复杂的逻辑是获取资源的分层信息，
主机对应的父层级信息有模块，集群，自定义mainline层级以及业务。为了尽量降低对业务逻辑的干扰，extensions模块用于实现在auth接口层的一层封装，场景层只需要传入资源ID即可实现组装出正确的
auth数据格式并进行资源注册或者鉴权。采用extensions模块有如下有点：
- 场景层/APIServer只需在适当位置加入对应函数调用, 无额外逻辑，保证了场景层代码主干清晰。
- extensions模块封装了对接iam的数据转换逻辑，目前对接时，这块相对不稳定，可能需要多次调整，对于这类调整，在extension模块中只需修改对应资源类型的拼装逻辑，对变更起到了隔离作用。
- 如上一步说描述，extensions模块可以方便的对接到其它类型的权限中心，而不用修改权限中心代码。

## what is it?
extensions模块是auth模块的一个功能扩展，使得它可以方便对接到各个场景层。比如向iam注册主机的接口，只需要传入hostid信息，extensions模块即可组装出最终iam需要的数据结构。


## design
extensions模块所有对接权限中心的接口均采用批量的形式，为降低问题复杂度，extension模块默认批量注册的资源属于同一个业务（全局资源的业务ID假设是0）。
extensions 被设计成一种松散的组织方式，各种资源类型之间尽量避免复用，目的是降低后期调整方案时的测试难度。

## code structure
- `types.go` 定义 extensions 模块用到的基本数据结构。
	+ `AuthManager` 类为extensions模块对外接口的封装
	+ `BusinessSimply` 以及 `SetSimply` 等类用于从db数据中提取extension模块必要的数据
- 其它`.go`文件为对接其它各个场景的接口封装，比如 `host.go` 封装了对主机的注册鉴权等相关接口及实现。


## special resource id design
### host
#### Resource_id format
- common format: `business/{business}:set/{set}:module/{module}:host/{host_id}`


#### Authorization for add host
- resource_id: `business/{business}:set/{set}:module/{module}:host`
- action: transferhost

####  Authorization for update/read/delete host
- resource_id: `business/{business}:set/{set}:module/{module}:host/{host_id}`
- action: update/read/delete


#### Deal with one host belong to multiple modules
- respect as multiple iam resource
- ex: host_id belong to module1 and module 2, then there will be two resource
    + resource_id: `business/{business}:set/{set}:module/{module1}:host/{host_id}`
    + resource_id: `business/{business}:set/{set}:module/{module2}:host/{host_id}`
