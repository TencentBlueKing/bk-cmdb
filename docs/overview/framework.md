# CMDB 3.0 二次开发框架使用指南

**注：处于开发状态，文档内容在正式发布前不保证兼容。**

## Inputer 接口声明

接口声明如下：

``` golang

// Inputer is the interface that must be implemented by every Inputer.
type Inputer interface {

	// Description the Inputer description.
	// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
	Name() string

	// Run execute the user logics
	Run() interface{}

	// Stop stop the run function
	Stop() error
}


```

Inputer 是必须要自己实现的接口。
> 1. Name 方法返回此Inputer的名字，如果Inputer运行过程中出现错误，框架会在输出的错误日志里用调用此方法，为了便于调试建议使用方返回唯一的名字。
> 2. Run 是Inputer 的运行方法，执行开发者的代码入口方法。
> 3. Stop 是Inputer的停止方法，如果Inputer是长期运行的，需要实现Stop方法来终止运行。

## API List


### Inputer 注册 API

#### 异常回调方法声明

``` golang

// ExceptionFunc the exception callback
type ExceptionFunc func(data interface{}, errMsg error)

```

#### 常规Inputer 注册，此方法注册的Inputer仅会被框架执行一次
> 方法：RegisterInputer(inputer input.Inputer, exceptionFunc input.ExceptionFunc)
> 
> 
> 参数：
> 
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - exceptionFunc：异常回调方法,在框架执行Inputer出现异常的时候会调用此方法，如果不需要此信息可以传入nil。
>
> 返回值：
> 无

#### 需要定时执行的Inputer注册，此方法注册的Inputer会被框架定时执行
> 方法：RegisterTimingInputer(inputer input.Inputer, frequency time.Duration, exceptionFunc input.ExceptionFunc)
>  
> 参数：
> 
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - frequency：执行此Inputer 的时间间隔。
>> - exceptionFunc：异常回调方法,在框架执行Inputer出现异常的时候会调用此方法，如果不需要此信息可以传入nil。
>
> 返回值：
> 无


### 业务管理 API

#### 业务管理包装器

> 类型：BusinessWrapper
> 方法列表：

``` golang
// SetValue 此方法用于对业务的字段进行赋值， key 业务字段，val 字段取值。
// 如果出错会有错误信息返回。
// 如有新增的自定义字段需要赋值的时候可以采用此方法。
SetValue(key string, val interface{}) error

// Save 执行保存逻辑。在保存时候会根据业务的唯一性配置监测当前业务是否存在，
// 如果存在则仅执行更新操作，如果当前业务不存在，则执行新建操作。
// 如果出错会有错误信息返回。
Save() error

// SetDeveloper 配置当前业务的开发人员，如果出错会有错误信息返回。
SetDeveloper(developer string) error 

// GetDeveloper 获取当前业务的开发人员，如果获取失败会有错误信息返回。
GetDeveloper() (string, error)

// SetMaintainer 配置当前业务的运维人员，如果配置失败会有错误信息返回。
SetMaintainer(maintainer string) error

// GetMaintainer 获取当前业务配置的运维人员，如果获取失败会有错误信息返回。
GetMaintainer() (string, error)

// SetName 配置当前业务的业务名，如果配置失败会有错误信息返回。
SetName(name string) error

// GetName 获取当前业务的业务名
// 如果获取失败会返回错误信息。
GetName() (string, error) 

// SetProductor 设置当前业务的产品人员
// 如果设置失败会有错误信息返回
SetProductor(productor string) error

// GetProductor 获取当前业务的产品人员
// 如果获取失败会有错误信息返回。
GetProductor() (string, error) 

// SetTester 设置当前业务的测试人员
// 如果设置失败会有错误信息返回。
SetTester(tester string) error 

// GetTester 设置当前业务的测试人员
// 如果获取失败会有错误信息返回。
GetTester() (string, error)

// SetSupplierAccount 设置当前业务的开发商ID
// 如果设置失败会有错误信息返回。
SetSupplierAccount(supplierAccount string) error 

// GetSupplierAccount 获取当前业务的开发商ID
// 如果获取失败会有错误信息返回。
GetSupplierAccount() (string, error)

// SetLifeCycle 设置当前业务的声明周期
// 此处只能取值以下枚举值：
// 测试中 BusinessLifeCycleTesting  
// 已上线 BusinessLifeCycleOnLine 
// 已停用 BusinessLifeCycleStopped
// 如果出错会返回错误信息
SetLifeCycle(lifeCycle string) error 

// GetLifeCycle 获取当前业务的状态
// 如果出错会返回错误信息
GetLifeCycle() (string, error)

// SetOperator 设置操作人员
// 如果出错会有错误返回
SetOperator(operator string) error

// GetOperator 获取操作人员
// 如果获取出错会有错误信息返回
GetOperator() (string, error)
```

#### 业务迭代器
> 类型：BusinessIteratorWrapper
> 方法列表：

``` golang
// Next 迭代读取每一个数据对象
// 当数据读取到最后一条的时候 error 会返回io.EOF
Next() (*BusinessWrapper, error)

// ForEach 遍历业务集合
ForEach(callback func(business *BusinessWrapper) error) error
```

#### 创建业务对象
> 方法：CreateBusiness(supplierAccount string) (*BusinessWrapper, error)
> 
> 参数：
>> - supplierAccount：开发商ID。
>
> 返回值：
> 
>> - BusinessWrapper： 业务管理对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建业务失败会返回错误。

#### 按照名字业务名进行查找
> 方法：FindBusinessLikeName(supplierAccount, businessName string) (*BusinessIteratorWrapper, error) 

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - businessName：业务名
>
> 返回值：
> 
>> - BusinessIteratorWrapper： 业务迭代器。
>> - error: 错误信息。

#### 按照条件对业务进行搜索
> 方法：FindBusinessByCondition(supplierAccount string, cond common.Condition) (*BusinessIteratorWrapper, error)

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - cond：查询条件
>
> 返回值：
> 
>> - BusinessIteratorWrapper： 业务迭代器。
>> - error: 错误信息。


### 集群管理 API

#### 集群管理包装器

> 类型：SetWrapper
> 方法列表：

``` golang
// SetValue 此方法用于对集群的字段进行赋值， key 集群字段，val 字段取值。
// 如果出错会有错误信息返回。
// 如有新增的自定义字段需要赋值的时候可以采用此方法。
SetValue(key string, val interface{}) error

// SetDescription 设置集群的描述信息
// 如果配置失败会返回错误信息
SetDescription(intro string) error 

// SetMark 设置集群的备注信息
SetMark(desc string) error 

// SetEnv 设置集群的环境
// 取值仅可以是以下枚举：
//  测试 SetEnvTesting      
//  体验 SetEnvGuest           
//  正式 SetEnvNormal             
SetEnv(env string) error

// GetEnv 获取集群的环境信息
// 获取失败会返回错误信息
GetEnv() (string, error)

// SetServiceStatus 设置服务的状态
// 取值只可以是以下枚举：
//  开放 SetServiceOpen
//  关闭 SetServiceClose
SetServiceStatus(status string) error 

// GetServiceStatus 获取服务状态
GetServiceStatus() (string, error)

// SetCapacity 设置集群的设计容量
SetCapacity(capacity int64) error

// GetCapacity 获取集群设计容量
GetCapacity() (int, error)

// SetBusinessID 设置集群所属的业务ID，调用此方法会同步将当前集群的父节点设置为传入的业务。
SetBusinessID(businessID int64) error

// GetBusiness 获取集群所属的业务ID，在调用Save 和Update 之后此处返回的仅是最后一次被更新的业务的ID
GetBusinessID() (int, error)

// SetSupplierAccount 设置集群所属的开发商ID
SetSupplierAccount(supplierAccount string) error

// GetSupplierAccount 获取集群所属的开发商ID
GetSupplierAccount() (string, error) 


// GetID 获取集群的ID，在调用Save 和Update 之后此处返回的仅是最后一次被更新的集群的ID
GetID() (int, error)

// SetParent 设置当前节点的父实例节点，只有在当前集群的父实例不是业务的时候才需要设置此参数。
SetParent(parentInstID int64) error

// SetName 设置集群的名字
SetName(name string) error

// GetName 获取集群的名字
GetName() (string, error) 

// Save 保存集群信息。在保存之前会监测当前集群信息是否已经存在，
// 如果存在则仅执行更新操作，如果不存在则执行新建操作。
Save() error
```

#### 集群迭代器
> 类型：SetIteratorWrapper
> 方法列表：

``` golang
// Next 迭代读取每一个数据对象
// 当数据读取完毕后，error 会返回 io.EOF
Next() (*SetWrapper, error)

// ForEach 对集群的集合进行迭代遍历
ForEach(callback func(set *SetWrapper) error) error
```

#### 创建集群对象
> 方法：CreateSet(supplierAccount string) (*SetWrapper, error) 
> 
> 参数：
> 
>> - - supplierAccount：开发商ID
>
> 返回值：
> 
>> - SetWrapper： 集群数据对象
>> - error：异常信息

#### 按照集群名进行查找
> 方法：FindSetLikeName(supplierAccount, setName string) (*SetIteratorWrapper, error) 

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - setName：集群名
>
> 返回值：
> 
>> - SetIteratorWrapper： 集群迭代器
>> - error: 错误信息

#### 按照条件对集群进行搜索
> 方法：FindSetByCondition(supplierAccount string, cond common.Condition) (*SetIteratorWrapper, error)

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - cond：查询条件
>
> 返回值：
> 
>> - SetIteratorWrapper： 集群迭代器。
>> - error: 错误信息。


### 模块管理 API

#### 模块管理包装器

> 类型：ModuleWrapper
> 方法列表：

``` golang
// SetValue 用于对模块自定义字段的值进行配置
SetValue(key string, val interface{}) error

// SetOperator 设置主要维护人
SetOperator(operator string) error

// GetOperator 获取主要维护人
GetOperator() (string, error) 

// SetBakOperator 设置备份维护人
SetBakOperator(bakOperator string) error 

// GetBakOperator 获取备份维护人
GetBakOperator() (string, error)

// SetTopo 设置模块的层级
SetTopo(bizID, setID int64) error

// GetBusinessID 获取业务ID
GetBusinessID() (int, error)

// SetSupplierAccount 设置开发商ID
SetSupplierAccount(supplierAccount string) error

// GetSupplierAccount 获取开发商ID
GetSupplierAccount() (string, error) 

// SetSetID 设置模块所属的集群
SetSetID(setID int64) error

// SetName 设置模块的名字
SetName(name string) error

// GetName 获取模块的名字
GetName() (string, error)

// GetID 获取模块的ID, 在调用Save 和Update 之后此处返回的仅是最后一次被更新的模块的ID
GetID() (int, error) 

// Save 保存模块的信息。
// 如果模块信息已经存在，则仅执行更新操作。
// 如果模块信息不存在，则执行新建操作。
Save() error

```

#### 模块迭代器
> 类型：ModuleIteratorWrapper
> 方法列表：

``` golang
// Next 迭代遍历每一个数据对象，如果遍历完集合 error 会返回io.EOF
Next() (*ModuleWrapper, error)

// ForEach 循环遍历模块集合，并将每个模块传递给回调函数。
ForEach(callback func(set *ModuleWrapper) error) error
```

#### 创建模块对象
> 方法：CreateModule(supplierAccount string) (*ModuleWrapper, error) 
> 
> 参数：
> 
>> - supplierAccount：开发商ID
>
> 返回值：
> 
>> - ModuleWrapper：模块数据对象
>> - error：异常信息

#### 按照模块名进行查找
> 方法：FindModuleLikeName(supplierAccount, moduleName string) (*ModuleIteratorWrapper, error) 

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - moduleName：集群名
>
> 返回值：
> 
>> - ModuleIteratorWrapper： 集群迭代器
>> - error: 错误信息

#### 按照条件对模块进行搜索
> 方法：FindModuleByCondition(supplierAccount string, cond common.Condition) (*ModuleIteratorWrapper, error)

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - cond：查询条件
>
> 返回值：
> 
>> - ModuleIteratorWrapper： 模块迭代器。
>> - error: 错误信息。


### 云区域管理 API

#### 云区域管理包装器

> 类型：PlatWrapper
> 方法列表：

``` golang
// SetValue 设置云区域自定义字段 的值 key 字段名 value 字段值
SetValue(key string, val interface{}) error
// SetSupplierAccount 设置云区域所属开发商 
SetSupplierAccount(supplierAccount string) error
// GetSupplierAccount 获取云区域开发商
GetSupplierAccount() (string, error)
// GetID 获取云区域ID,在调用Save 和Update 之后此处返回的仅是最后一次被更新的云区域的ID
GetID() (int, error) 
// SetName 设置云区域名字
SetName(name string) error

// GetName 获取云区域名字
GetName() (string, error) 
```

#### 云区域迭代器
> 类型：PlatIteratorWrapper
> 方法列表：

``` golang
// Next 迭代遍历每一个数据对象，如果遍历完集合 error 会返回io.EOF
Next() (*PlatWrapper, error)

// ForEach 循环遍历模块集合，并将每个模块传递给回调函数。
ForEach(callback func(plat *PlatWrapper) error) error
```

#### 创建云区域对象
> 方法：CreatePlat(supplierAccount string) (*PlatWrapper, error) 
> 
> 参数：
> 
>> - supplierAccount：开发商ID
>
> 返回值：
> 
>> - PlatWrapper：模块数据对象
>> - error：异常信息

#### 按照云区域名进行查找
> 方法：FindPlatLikeName(supplierAccount, moduleName string) (*PlatIteratorWrapper, error) 

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - moduleName：集群名
>
> 返回值：
> 
>> - PlatIteratorWrapper： 云区域迭代器
>> - error: 错误信息

#### 按照条件对模块进行搜索
> 方法：FindPlatByCondition(supplierAccount string, cond common.Condition) (*PlatIteratorWrapper, error)

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - cond：查询条件
>
> 返回值：
> 
>> - PlatIteratorWrapper： 云区域迭代器。
>> - error: 错误信息。

### 主机管理 API

#### 主机转移接口方法声明

> 类型：TransferInterface

> 方法列表:

``` golang

// MoveToModule 将主机转移到主机当前所在业务的其他模块，isIncrement true 主机原来所在模块不会被改变，false 会从原来所在模块中删除
MoveToModule(newModuleIDS []int64, isIncrement bool) error

// MoveToFaultModule 将主机移动到主机当前所在业务的故障机模块
MoveToFaultModule() error

// MoveToIdleModule 将主机移动到主机当前所在业务的空闲机模块
MoveToIdleModule() error

// MoveToResourcePools 将主机移动到资源池
MoveToResourcePools() error

// MoveToBusinessIdleModuleFromResourcePools 将主机从资源池分配到业务空闲机模块
MoveToBusinessIdleModuleFromResourcePools(bizID int64) error

// MoveToAnotherBusinessModules 将主机从当前业务转移到另一个业务的给定模块下
MoveToAnotherBusinessModules(bizID int64, moduleID int64) error

// ResetBusinessHosts 将主机从给定的模块和集群下清空，转移至业务的空闲机下
ResetBusinessHosts(setID, moduleID int64) error

```

#### 主机管理包装器（用于查询接口返回的数据结构）

> 类型：FinderHostWrapper

> 方法列表：


``` golang

// GetBizs 获取业务信息
GetBizs() ([]*BusinessWrapper, error) 

// GetSets 获取业务信息
GetSets() ([]*SetWrapper, error)

// GetModules 获取模块信息
GetModules() ([]*ModuleWrapper, error) 

其余方法与HostWrapper一致

```

#### 主机管理包装器

> 类型：HostWrapper
> 方法列表：

``` golang

// Transfer 返回主机转移操作方法
Transfer() inst.TransferInterface

// SetModuleIDS 设置主机所属业务的模块ID,HostAppendModule 表示追加所属模块，HostReplaceModule 表示替换所属模块
SetModuleIDS(moduleIDS []int64, act HostModuleActionType) 

// SetBusiness 设置主机所属的业务
SetBusiness(bizID int64)

// SetTopo 设置主机所属的业务及业务下的模块ID, act 取值，HostAppendModule 表示追加所属模块，HostReplaceModule 表示替换所属模块
SetTopo(bizID int64, setName, moduleName string, act HostModuleActionType) error

// SetValue 为自定义字段进行复制，key 字段名，val 字段的值
SetValue(key string, val interface{}) error

// GetModel 获取主机所对应的模型定义对象
GetModel() model.Model 

// SetBakOperator 设置备份维护人
SetBakOperator(bakOperator string) error

// GetBakOperator 获取备份维护人
GetBakOperator() (string, error)

// SetOsBit 设置操作系统位数
SetOsBit(osbit string) error 

// GetOsBit 获取操作系统位数
GetOsBit() (string, error)

// SetSLA 设置SLA安全级别
// 取值尽可以为以下枚举值之一：
//    HostSLALevel1            
//    HostSLALevel2            
//    HostSLALevel3            
SetSLA(sla string) error

// GetSLA 获取SLA安全界别
GetSLA() (string, error)

// SetCloudID 设置云区域ID
SetCloudID(cloudID int64) error

// GetCloudID 获取云区域ID
GetCloudID() (int, error)

// SetInnerIP 设置内网IP
SetInnerIP(innerIP string) error

// GetInnerIP 获取内网IP
GetInnerIP() (string, error)

// SetOpeartor 设置主维护人
SetOperator(operator string) error

// GetOperator 获取主维护人
GetOperator() (string, error) 

// SetCPU 设置CPU逻辑核心数
SetCPU(cpu int64) error 

// GetCPU 获取CPU逻辑核心数
GetCPU() (int, error)

// SetCPUMhz 设置CPU频率
SetCPUMhz(cpuMhz int64) error

// GetCPUMhz 获取CPU频率
GetCPUMhz() (int, error)

// SetOSType 获取OS类型
// 取值仅可以是以下列表之一：
//    HostOSTypeLinux         
//    HostOSTypeWindows       
SetOsType(osType string) error 

// GetOsType 获取操作系统类型
GetOsType() (string, error)

// SetOuterIP 设置外网IP
SetOuterIP(outerIP string) error

// GetOuterIP 获取外网IP
GetOuterIP() (string, error)

// SetAssetID 设置固资编号
SetAssetID(assetID string) error

// GetAssetID 获取固资编号
GetAssetID() (string, error) 

// SetInnerMac 设置内网Mac 地址
SetInnerMac(mac string) error 

// GetInnerMac 获取内网Mac 地址
GetInnerMac() (string, error)

// SetOuterMac 设置内网Mac 地址
SetOuterMac(mac string) error 

// GetOuterMac 获取内网Mac 地址
GetOuterMac() (string, error)

// SetSN 设备SN
SetSN(sn string) error

// GetSN 获取设备SN
GetSN() (string, error)

// SetCPUModule 设置CPU型号
SetCPUModule(cpuModule string) error

// GetCPUModule 获取CPU型号
GetCPUModule() (string, error) 

// SetName 设置主机名
SetName(hostName string) error

// GetName 获取主机名
GetName() (string, error) 

// SetServiceTerm 设置质保年限
SetServiceTerm(serviceTerm int64) error

// GetServiceTerm 获取质保年限
GetServiceTerm() (int, error)

// SetComment 设置备注
SetComment(comment string) error

// GetComment 获取备注信息
GetComment() (string, error)

// SetMem 设置内存容量
SetMem(mem int64) error

// GetMem 获取内存容量
GetMem() (int, error) 

// SetDisk 设置磁盘容量
SetDisk(disk int64) error 

// GetDisk 获取磁盘容量
GetDisk() (int, error)

// SetOsName 设置操作系统名
SetOsName(osName string) error

// GetOsName 获取操作系统名
GetOsName() (string, error)

// SetOsVersion 设置操作系统版本
SetOsVersion(osVersion string) error 
// GetOsVersion 获取操作系统版本
GetOsVersion() (string, error) 

// Save 保存主机信息。
// 如果主机已经存在，则仅执行更新操作。
// 如果主机不存在，则仅执行新建操作。
Save() error
```

#### 主机迭代器
> 类型：HostIteratorWrapper
> 方法列表：

``` golang
// Next 迭代获取主机数据对象，如果已经遍历完所有数据，那么error 会返回io.EOF
Next() (*HostWrapper, error)

// ForEach 遍历主机数据对象集合
ForEach(callback func(host *HostWrapper) error) error
```

#### 创建主机对象
> 方法：CreateHost(supplierAccount string) (*HostWrapper, error) 
> 
> 参数：
> 
>> - supplierAccount：开发商ID
>
> 返回值：
> 
>> - HostWrapper：主机数据对象
>> - error：异常信息

#### 按照主机名进行查找
> 方法：FindHostLikeName(supplierAccount, hostName string) (*HostIteratorWrapper, error) 

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - hostName：集群名
>
> 返回值：
> 
>> - hostIteratorWrapper： 主机迭代器
>> - error: 错误信息

#### 按照条件对主机进行搜索
> 方法：FindHostByCondition(supplierAccount string, cond common.Condition) (*HostIteratorWrapper, error)

> 
> 参数：
>> - supplierAccount：开发商ID。
>
>> - cond：查询条件
>
> 返回值：
> 
>> - HostIteratorWrapper： 主机迭代器。
>> - error: 错误信息。



### 通用实例管理 API


#### 通用实例

> 类型：Inst
> 方法列表：

``` golang

// GetModel 获取当前实例多对应的模型
GetModel() model.Model

// IsMainLine 用于判断当前实例是否是主线模型
// [开发中]当前不可用 
IsMainLine() bool

// GetAssociationModels 获取当前实例直接关联的所有模型的集合
// [开发中]当前不可用 
GetAssociationModels() ([]model.Model, error)

// GetInstID 获取当前实例的实例ID
GetInstID() int

// GetInstName 获取当前实例的实例名
GetInstName() string

// SetValue 为实例的字段赋值，key 实例的字段， value 字段的取值
SetValue(key string, value interface{}) error

// GetValues 获取当前实例的所有字段的及对应的值
GetValues() (types.MapStr, error)

// GetAssociationsByModleID 获取当前实例关联的某个模型的所有实例的集合
// [开发中]当前不可用
GetAssociationsByModleID(modleID string) ([]Inst, error)

// GetAllAssociations 获取当前实例直接关联的所有实例的集合
// [开发中]当前不可用
GetAllAssociations() (map[model.Model][]Inst, error)

// SetParent 设置当前实例的父实例
SetParent(parentInstID int) error

// GetParent 获取当前实例所在拓扑结构中所有的父节点的集合
// [开发中]当前不可用
GetParent() ([]Topo, error)

// GetChildren 获取当前实例所在拓扑结构中所有的子节点的结合
// [开发中]当前不可用
GetChildren() ([]Topo, error)
```

#### 创建普通对象
> 方法：CreateCommonInst(target model.Model) (inst.Inst, error)
> 
> 参数：
> 
>> - target：用于指明是创建的实例所属的模型定义。
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建实例失败会返回错误。


#### 按照实例名进行查找
> 方法：FindInstsLikeName(target model.Model, instName string) (*inst.Iterator, error) 

> 
> 参数：
>> - target：目标实例的模型
>
>> - instName：实例名
>
> 返回值：
> 
>> - inst.Iterator： 实例迭代器
>> - error: 错误信息

#### 按照条件对实例进行搜索
> 方法：FindInstsByCondition(target model.Model, cond common.Condition) (inst.Iterator, error)

> 
> 参数：
>> - target：目标实例的模型
>
>> - cond：查询条件
>
> 返回值：
> 
>> - inst.Iterator： 实例迭代器。
>> - error: 错误信息。
>> 

### 模型管理 API

#### 获取模型
> 方法：GetModel(supplierAccount, classificationID, objID string) (model.Model, error) 
> 
> 参数：
> 
>> - supplierAccount：开发商ID
> 
>> - classificationID：模型分类ID
>> 
>> - objID：模型ID
>> 
> 返回值：
> 
>> - model.Model： 模型对象
>> - error: 查询失败会返回异常信息

#### 创建模型分类对象
> 方法：CreateClassification(name string) model.Classification 
> 
> 参数：
> 
>> - name：分类的名字
>> 
> 返回值：
> 
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。

### 按照模型分类的名字进行模糊查找，返回所有名字与输入的名字相似的分类对象的迭代器。
> 方法：FindClassificationsLikeName(name string) (model.ClassificationIterator, error)
> 
> 参数：
> 
>> - name：分类的名字
> 返回值：
> 
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。

### 按照条件进行精确查找，返回所有符合条件的分类对象的迭代器
> 方法：FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error)
> 
> 参数：
> 
>> - name：分类的名字
> 返回值：
> 
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。


### 事件订阅 API

#### 事件回调方法声明

```` golang
// EventCallbackFunc the event deal function
type EventCallbackFunc func(evn []*Event) error
````

#### 取消事件订阅

> 方法：UnRegisterEvent(eventKey types.EventKey)
> 
> 参数：
> 
>> - eventKey：注册事件后返回的Key

>> 
> 返回值：
> 
>> 无

#### 订阅主机信息变更事件

> 方法：RegisterEventHost(eventFunc types.EventCallbackFunc) types.EventKey
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key


#### 订阅业务信息变更事件

> 方法：RegisterEventBusiness(eventFunc types.EventCallbackFunc) types.EventKey
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key

#### 订阅模块信息变更事件

> 方法：RegisterEventModule(eventFunc types.EventCallbackFunc) types.EventKey 
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key


#### 订阅主机身份信息变更事件

> 方法：RegisterEventHostIdentifier(eventFunc types.EventCallbackFunc) types.EventKey 
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key


#### 订阅集群信息变更事件

> 方法：RegisterEventSet(eventFunc types.EventCallbackFunc) types.EventKey
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key

#### 订阅通用模型实例信息变更事件

> 方法：RegisterEventInst(eventFunc types.EventCallbackFunc) types.EventKey
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key

#### 订阅模块转移信息变更事件

> 方法：RegisterEventModuleTransfer(eventFunc types.EventCallbackFunc) types.EventKey 
> 
> 参数：
> 
>> - eventFunc：事件回调方法

>> 
> 返回值：
> 
>> - 注册事件关联的Key
