# filter

filter 是一种查询表达式。

## 特性
1. 支持多种查询操作符。
2. 支持JSON字段操作符：=、in。
3. 支持多种 value 类型。
4. 支持嵌套。


## 函数功能说明
- Validate(opt *ExprOption) (hitErr error) - 用于校验Filter的合法性，ExprOption传入特定的限制参数。
- UnmarshalJSON(raw []byte) error - 自定义JSON序列化函数
- WithType() RuleType - 实现RuleFactory，嵌套表达式需要，返回表达式类型。
- RuleField() string - 实现RuleFactory，嵌套表达式需要，返回表达式字段。
注：
1. 需注意当使用JSON字段操作符时，字段名仅需要将嵌套字段通过 '.' 关联即可。e.g: "extension.vpc.id"

## 示例
1. 名称为cmdb，且年龄大于18岁。
```go

expr := ExpressionAnd(
    RuleEqual("name","cmdb"),
    RuleGreaterThan("age",18),
)

// same as
expr = &Expression{
    Op: And,
    Rules: []RuleFactory{
        &AtomRule{
            Field: "name",
            Op:    Equal.Factory(),
            Value: "cmdb",
        },
        &AtomRule{
            Field: "age",
            Op:    GreaterThan.Factory(),
            Value: 18,
        },
    },
}
```

2. 名称为cmdb，且年龄大于18或者身高小于1.8。
```go
expr := ExpressionAnd(
    RuleEqual("name","cmdb"),
    ExpressionOr(
        RuleGreaterThan("age",18),
        RuleLessThan("height",1.8),
    ),
)
```


3. 名称为cmdb，且 Extension Json字段中vpc的id为3。
```go
expr := ExpressionAnd(
    RuleEqual("name", "cmdb"),
    RuleJSONEqual("extension.vpc.id", 3),
)
```

4. 查询管理者中有cmdb的数据。
```go
expr := ExpressionAnd(
    RuleJSONContains("managers", "cmdb"),
)

```
