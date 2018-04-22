# CMDB 3.0 二次开发框架说明

## Inputer 接口声明

接口声明如下：

``` golang

// Inputer is the interface that must be implemented by every Inputer.
type Inputer interface {

    // Name the Inputer description.
    // This information will be printed when the Inputer is abnormal, which is convenient for debugging.
    Name() string

    // Input should not be blocked
    Input() interface{}
}

```

Inputer 是必须要自己实现的接口。
> 1. Name 方法返回此Inputer的名字，如果Inputer运行过程中出现错误，框架会在输出的错误日志里用调用此方法，为了便于调试建议使用方返回唯一的名字。
> 2. Input 是Inputer 接口唯一向框架输出数据的接口。开发者需要在此方法的实现里面自己组织需要的数据，并且是非阻塞实现，并将经过以下三个方法进行处理后的数据返回：
>> - CreateTransaction
>> - CreateTimingTransaction
>> - CreateCommonEvent

## API List

### 常规Inputer 注册，此方法注册的Inputer仅会被框架执行一次
> 1. 方法：RegisterInputer(inputer input.Inputer, putter output.Puter, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) 
> 2. 参数：
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - putter：自定义的 output.Putter接口实现。
>> - exceptionFunc：异常回调方法。在框架执行Inputer出现异常的时候会调用此方法。
> 3. 返回值：
>> - input.InputerKey：Inputer 成功注册如框架后，框架会为此Inputer生成一个唯一的Key。
>> - error：注册Inputer失败后的错误信息。

### 需要定时执行的Inputer注册，此方法注册的Inputer会被框架定时执行
> 1. 方法：RegisterTimingInputer(inputer input.Inputer, putter output.Puter, frequency time.Duration, exceptionFunc input.ExceptionFunc) (input.InputerKey, error) 
> 2. 参数：
>> - inputer：所有实现了input.Inputer接口的对象实例。
>> - putter：自定义的 output.Putter接口实现。
>> - frequency：执行此Inputer 的时间间隔。
>> - exceptionFunc：异常回调方法。在框架执行Inputer出现异常的时候会调用此方法。
> 3. 返回值：
>> - input.InputerKey：Inputer 成功注册如框架后，框架会为此Inputer生成一个唯一的Key。
>> - error：注册Inputer失败后的错误信息。

### 创建事务对象，被此对象包装过的对象会被归类为一个事务，执行过程不会被打断。
> 1. 方法：CreateTransaction() input.Transaction
> 2. 参数：
>> - 无输入参数
> 3. 返回值：
>> - input.Transaction：事务对象，可以容纳所有实现了Saver接口的方法。


### 创建定时事务对象，被此对象包装过的对象会被归类为一个事务，执行过程不会被打断。
> 1. 方法：CreateTimingTransaction(duration time.Duration) input.Transaction
> 2. 参数：
>> - duration：次事务被执行的时间间隔。
> 3. 返回值：
>> - input.Transaction：事务对象，可以容纳所有实现了Saver接口的方法。

### 创建普通事件
> 1. 方法：CreateCommonEvent(saver types.Saver) interface{} 
> 2. 参数：
>> - saver: 所有实现了Saver接口的方法。（框架层面提供的：inst, model 系列的接口均有实现saver接口）
> 3. 返回值：
>> - 包装后的对象。


### 创建业务对象
> 1. 方法：CreateBusiness(supplierAccount string)(inst.Inst, error) 
> 2. 参数：
>> - supplierAccount：开发商ID。
> 3. 返回值：
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建业务失败会返回错误。

### 创建业务对象
> 1. 方法：CreateSet() (inst.Inst, error)
> 2. 参数：
>> - 无输入参数
> 3. 返回值：
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建集群失败会返回错误。

### 创建模块对象
> 1. 方法：CreateModule() (inst.Inst, error)
> 2. 参数：
>> - 无输入参数
> 3. 返回值：
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建模块失败会返回错误。

### 创建普通对象
> 1. 方法：CreateCommonInst(target model.Model) (inst.Inst, error)
> 2. 参数：
>> - target：用于指明是创建的实例所属的模型定义。
> 3. 返回值：
>> - inst.Inst： 实例接口对象，包含对当前实例数据进行维护的接口。
>> - error: 如果创建实例失败会返回错误。

### 创建模型分类对象
> 1. 方法：CreateClassification() model.Classification 
> 2. 参数：
>> - 无输入参数
> 3. 返回值：
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。

### 按照模型分类的名字进行模糊查找，返回所有名字与输入的名字相似的分类对象的迭代器。
> 1. 方法：FindClassificationsLikeName(name string) (model.ClassificationIterator, error)
> 2. 参数：
>> - name：分类的名字
> 3. 返回值：
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。

### 按照条件进行精确查找，返回所有符合条件的分类对象的迭代器
> 1. 方法：FindClassificationsByCondition(condition types.MapStr) (model.ClassificationIterator, error)
> 2. 参数：
>> - name：分类的名字
> 3. 返回值：
>> - model.Classification：模型分类对象，通过此对象可以对该分类下的模型进行管理。
>> - error: 如果创建实例失败会返回错误。

## 应用示例

``` golang 
package example

import (
    "configcenter/src/framework/api"
    "configcenter/src/framework/core/types"
    "fmt"
)

func init() {

    _, sender, _ := api.CreateCustomOutputer("example_output", func(data types.MapStr) error {
        fmt.Println("outputer:", data)
        return nil
    })

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
func (cli *myInputer) Input() interface{} {
    fmt.Println("my_inputer")

    // 1. 返回 MapStr对象，此方法用于有Inputer绑定了自定义Outputer的时候使用，内置Outputer不采用此方法传递数据。
    /**
    return types.MapStr{
        "test": "outputer",
        "hoid": "",
    }
    */

    // 此方法仅用于内置Outputer 的数据返回
    // 1. 构建模型分类
    // 2. 通过模型分类构建model
    // 3. 通过model 构建模型属性
    // 4. 利用包装器对要返回的数据做处理。
    cls := api.CreateClassification()

    model := cls.CreateModel()
    attr := model.CreateAttribute()
    attr.SetName("test")

    return api.CreateCommonEvent(cls)

}

```