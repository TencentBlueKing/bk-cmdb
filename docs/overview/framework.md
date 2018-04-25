# CMDB 3.0 二次开发框架说明

**注：功能还在开发中，此文档仅供参考，在开发过程中，最终发布之前此文档内列举的接口有可能会微调。**

## Inputer 接口声明

接口声明如下：

``` golang

// Inputer is the interface that must be implemented by every Inputer.
type Inputer interface {

    // Name the Inputer description.
    // This information will be printed when the Inputer is abnormal, which is convenient for debugging.
    Name() string

    // Run the execution function
    Run() interface{}
	
    // Stop stop the Input function
    Stop() error

}

```

Inputer 是必须要自己实现的接口。
> 1. Name 方法返回此Inputer的名字，如果Inputer运行过程中出现错误，框架会在输出的错误日志里用调用此方法，为了便于调试建议使用方返回唯一的名字。
> 2. Input 是Inputer 的方法，执行开发者的代码入口方法。

## API List

### 创建构造查询条件对象，用于构建查询条件
> 方法：CreateCondition() *common.Condition
> 
> 参数：
> 
>> - tableName：表名。
>
> 返回值：
> 
>> - Condition：查询条件对象，用于构造查询条件。

### 常规Inputer 注册，此方法注册的Inputer仅会被框架执行一次
> 方法：RegisterInputer(inputer input.Inputer, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) 
> 
> 
> 参数：
> 
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - putter：自定义的output.Putter接口实现。
>> - exceptionFunc：异常回调方法。在框架执行Inputer出现异常的时候会调用此方法。
>
> 返回值：
>> - input.InputerKey：Inputer 成功注册如框架后，框架会为此Inputer生成一个唯一的Key。
>> - error：注册Inputer失败后的错误信息。

### 需要定时执行的Inputer注册，此方法注册的Inputer会被框架定时执行
> 方法：RegisterTimingInputer(inputer input.Inputer, putter output.Puter, frequency time.Duration, exceptionFunc input.ExceptionFunc) (input.InputerKey, error)
>  
> 参数：
> 
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - putter：自定义的 output.Putter接口实现。
>> - frequency：执行此Inputer 的时间间隔。
>> - exceptionFunc：异常回调方法。在框架执行Inputer出现异常的时候会调用此方法。
>
> 返回值：
>> - input.InputerKey：Inputer 成功注册如框架后，框架会为此Inputer生成一个唯一的Key。
>> - error：注册Inputer失败后的错误信息。



### 创建业务对象
> 方法：CreateBusiness(supplierAccount string)(inst.Inst, error) 
> 
> 参数：
>> - supplierAccount：开发商ID。
>
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建业务失败会返回错误。

### 创建集群对象
> 方法：CreateSet() (inst.Inst, error)
> 
> 参数：
> 
>> - 无输入参数
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建集群失败会返回错误。

### 创建模块对象
> 方法：CreateModule() (inst.Inst, error)
> 
> 参数：
> 
>>- 无输入参数
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建模块失败会返回错误。

### 创建主机对象
> 方法：CreateHost()(inst.Inst, error) 
> 
> 参数：
>> - 无输入参数
>
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建业务失败会返回错误。

### 创建普通对象
> 方法：CreateCommonInst(target model.Model) (inst.Inst, error)
> 
> 参数：
> 
>> - target：用于指明是创建的实例所属的模型定义。
> 返回值：
> 
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建实例失败会返回错误。

### 创建模型分类对象
> 方法：CreateClassification() model.Classification 
> 
> 参数：
> 
>> - 无输入参数
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


## 内置Outputer 接口声明

### 模型分类管理接口

``` golang

// ClassificationIterator the classification iterator
type ClassificationIterator interface {
	Next() (Classification, error)
}

// Classification the interface declaration for model classification
type Classification interface {
	types.Saver

	SetID(id string)
	SetName(name string)
	SetIcon(iconName string)
	GetID() string

	CreateModel() Model
	FindModelsLikeName(modelName string) (Iterator, error)
	FindModelsByCondition(condition *common.Condition) (Iterator, error)
}
```

> 模型分类对象可以看由以下API构建：
>> - CreateClassification() model.Classification
>> - FindClassificationsLikeName(name string) (model.ClassificationIterator, error) 
>> - FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error) 

### 模型管理接口

``` golang


// Iterator the model iterator
type Iterator interface {
	Next() (Model, error)
}

// Model the interface declaration for model maintence
type Model interface {
	types.Saver
	SetIcon(iconName string)
	GetIcon() string
	SetID(id string)
	GetID() string
	SetName(name string)
	GetName() string

	SetPaused()
	SetNonPaused()
	Paused() bool

	SetPosition(position string)
	GetPosition() string
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDescription(desc string)
	GetDescription() string
	SetCreator(creator string)
	GetCreator() string
	SetModifier(modifier string)
	GetModifier() string

	CreateAttribute() Attribute
	CreateGroup() Group

	FindAttributesLikeName(attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error)

	FindGroupsLikeName(groupName string) (GroupIterator, error)
	FindGroupsByCondition(condition *common.Condition) (GroupIterator, error)
}

```

> 模型管理对象实例必须由模型分类接口创建。


### 分组管理接口

``` golang
// Group the interface declaration for model maintence
type Group interface {
	types.Saver

	SetID(id string)
	GetID() string
	SetName(name string)
	SetIndex(idx int)
	GetIndex() int
	SetSupplierAccount(ownerID string)
	GetSupplierAccount() string
	SetDefault()
	SetNonDefault()
	Default() bool

	CreateAttribute() Attribute
	FindAttributesLikeName(attributeName string) (AttributeIterator, error)
	FindAttributesByCondition(condition *common.Condition) (AttributeIterator, error)
}

// GroupIterator the group iterator
type GroupIterator interface {
	Next() (Group, error)
}
``` 

> 模型属性分组必须由具体的模型对象构建。

### 模型属性管理接口

``` golang

// AttributeIterator the attribute iterator
type AttributeIterator interface {
	Next() (Attribute, error)
}

// Attribute the interface declaration for model attribute maintence
type Attribute interface {
	types.Saver

	SetID(id string)
	SetName(name string)
	SetUnit(unit string)
	SetPlaceholer(placeHoler string)
	SetEditable()
	SetNonEditable()
	Editable() bool
	SetRequired()
	SetNonRequired()
	Required() bool
	SetKey(isKey bool)
	SetOption(option string)
	SetDescrition(des string)
}

```

> 模型属性可以由具体的模型对象构建，也可以由具体的模型分组对象构建。

## 应用示例

``` golang 
package example

import (
    "configcenter/src/framework/api"
    "configcenter/src/framework/core/types"
    "fmt"
)

func init() {
    api.RegisterInputer(target, sender, nil)
}

var target = &myInputer{}

type myInputer struct {
}

// Description the Inputer description.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *myInputer) Name() string {
    return "name_myinputer"
}

// Input the input should not be blocked
func (cli *myInputer) Run() interface{} {
   
    // 1. 构建模型分类
    // 2. 通过模型分类构建model
    // 3. 通过model 构建模型属性
    cls := api.CreateClassification()
    cls.SetName("test_name")
    cls.Save()

    model := cls.CreateModel()
    attr := model.CreateAttribute()
    attr.SetName("test")

    return nil

}

func(cli *myInputer)Stop()error{
    return nil
}

```
