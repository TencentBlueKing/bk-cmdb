# CMDB 3.0 二次开发框架使用指南


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
// SetValue 此方法用于对业务的字段进行复制， key 业务字段，val 字段取值。
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
SetValue(key string, val interface{}) error
SetDescription(intro string) error 
SetMark(desc string) error 
SetEnv(env string) error
GetEnv() (string, error)
SetServiceStatus(status string) error 
GetServiceStatus() (string, error)
SetCapacity(capacity int64) error
GetCapacity() (int, error)
SetBussinessID(businessID int64) error
GetBusinessID() (int, error)
SetSupplierAccount(supplierAccount string) error
GetSupplierAccount() (string, error) 
SetID(id string) error 
GetID() (string, error)
SetParent(parentInstID int64) error
SetName(name string) error
GetName() (string, error) 
Save() error
```

#### 集群迭代器
> 类型：SetIteratorWrapper
> 方法列表：

``` golang
Next() (*SetWrapper, error)
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
SetValue(key string, val interface{}) error
SetOperator(operator string) error
GetOperator() (string, error) 
SetBakOperator(bakOperator string) error 
GetBakOperator() (string, error)
SetBussinessID(businessID int64) error
GetBusinessID() (int, error)
SetSupplierAccount(supplierAccount string) error
GetSupplierAccount() (string, error) 
SetParent(parentInstID int64) error
SetName(name string) error
GetName() (string, error)
SetID(id string) error
GetID() (string, error) 
Save() error
```

#### 模块迭代器
> 类型：ModuleIteratorWrapper
> 方法列表：

``` golang
Next() (*ModuleWrapper, error)
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


### 主机管理 API

#### 主机管理包装器

> 类型：HostWrapper
> 方法列表：

``` golang
SetValue(key string, val interface{}) error
GetModel() model.Model 
SetBakOperator(bakOperator string) error
GetBakOperator() (string, error)
SetOsBit(osbit string) error 
GetOsBit() (string, error)
SetSLA(sla string) error
GetSLA() (string, error)
SetCloudID(cloudID int64) error
GetCloudID() (int, error)
SetInnerIP(innerIP string) error
GetInnerIP() (string, error)
SetOperator(operator string) error
GetOperator() (string, error) 
SetStateName(stateName string) error
GetStateName() (string, error)
SetCPU(cpu int64) error 
GetCPU() (int, error)
SetCPUMhz(cpuMhz int64) error
GetCPUMhz() (int, error)
SetOsType(osType string) error 
GetOsType() (string, error)
SetOuterIP(outerIP string) error
GetOuterIP() (string, error)
SetAssetID(assetID string) error
GetAssetID() (string, error) 
SetMac(mac string) error 
GetMac() (string, error)
SetProvinceName(provinceName string) error 
GetProvinceName() (string, error)
SetSN(sn string) error
GetSN() (string, error)
SetCPUModule(cpuModule string) error
GetCPUModule() (string, error) 
SetName(hostName string) error
GetName() (string, error) 
SetISPName(ispName string) error
GetISPName() (string, error) 
SetServiceTerm(serviceTerm int64) error
GetServiceTerm() (int, error)
SetComment(comment string) error
GetComment() (string, error)
SetMem(mem int64) error
GetMem() (int, error) 
SetDisk(disk int64) error 
GetDisk() (int, error)
SetOsName(osName string) error
GetOsName() (string, error)
SetOsVersion(osVersion string) error 
GetOsVersion() (string, error) 
Save() error
```

#### 主机迭代器
> 类型：HostIteratorWrapper
> 方法列表：

``` golang
Next() (*HostWrapper, error)
ForEach(callback func(set *HostWrapper) error) error
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
GetModel() model.Model
IsMainLine() bool
GetAssociationModels() ([]model.Model, error)
GetInstID() int
GetInstName() string
SetValue(key string, value interface{}) error
GetValues() (types.MapStr, error)
GetAssociationsByModleID(modleID string) ([]Inst, error)
GetAllAssociations() (map[model.Model][]Inst, error)
SetParent(parentInstID int) error
GetParent() ([]Topo, error)
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

