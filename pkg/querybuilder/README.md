## What is it
一个用于前端/API的查询参数解析mongodb过滤条件的后端模块, 支持类 [jQuery-QueryBuilder](https://github.com/mistic100/jQuery-QueryBuilder) 形式的输入参数.

## 与 jQuery-QueryBuilder 的区别
- 不支持 type 字段，所有value字段均解析成基本字段类型
- 为支持基于时间类型的对比，扩展除了相应的操作符
	+ OperatorDatetimeLess           ("datetime_less")
	+ OperatorDatetimeLessOrEqual    ("datetime_less_or_equal")
	+ OperatorDatetimeGreater        ("datetime_greater")
	+ OperatorDatetimeGreaterOrEqual ("datetime_greater_or_equal")
- 不支持 `between` 和 `not_between` 运算符, 这类运算符可基于基本比较运算符组合实现

## How it implemented

### Rule
过滤规则的接口

- `GetDeep` 定义过滤规则的深度
- `Validate` 校验过滤规则是否有效
- `ToMgo` 转换成`mongodb`查询条件

### AtomRule
原子过滤规则，任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则按逻辑与/或组合而成

#### AtomRule.Field 字段
#### AtomRule.Operator 字段
#### AtomRule.Value 字段

### CombinedRule
组合过滤规则，组合的节点可以是原子过滤规则或组合过滤规则

### RuleParser
过滤规则解析方法，从`map[string]interface{}`数据中解析出一个过滤规则实例

## Operator 详细说明
### 通用操作符
- OperatorEqual    ("equal")
    + 含义：相等比较
    + Value格式： 基本数据类型(数值/bool值/字符串)
- OperatorNotEqual ("not_equal")
    + 含义：不等比较
    + Value格式： 基本数据类型(数值/bool值/字符串)

### 数组操作符
- OperatorIn    ("in")
    + 含义：匹配记录字段值是否在指定集合中
    + Value格式： 基本数据类型组成的数值，类型需要一致
- OperatorNotIn ("not_in")
    + 含义：匹配记录字段值不在指定集合中
    + Value格式： 基本数据类型组成的数值，类型需要一致
- OperatorIsEmpty    ("is_empty")
    + 含义：匹配记录字段值为空数组
    + Value格式： 不接受参数
- OperatorIsNotEmpty ("is_not_empty")
    + 含义：匹配记录字段值为非空数组
    + Value格式： 不接受参数
	
### 数字操作符
- OperatorLess           ("less")
    + 含义：匹配记录字段值 < `{Value}`
    + Value格式： 数值
- OperatorLessOrEqual    ("less_or_equal")
    + 含义：匹配记录字段值 <= `{Value}`
    + Value格式： 数值
- OperatorGreater        ("greater")
    + 含义：匹配记录字段值 > `{Value}`
    + Value格式： 数值
- OperatorGreaterOrEqual ("greater_or_equal")
    + 含义：匹配记录字段值 >= `{Value}`
    + Value格式： 数值

### 时间操作符
> 目前仅支持 `RFC3339` 时间格式
- OperatorDatetimeLess           ("datetime_less")
    + 含义：匹配记录字段值表示的时间早于 < `{Value}`
    + Value格式： `RFC3339` 格式字符串
- OperatorDatetimeLessOrEqual    ("datetime_less_or_equal")
    + 含义：匹配记录字段值表示的时间不晚于 <= `{Value}`
    + Value格式： `RFC3339` 格式字符串
- OperatorDatetimeGreater        ("datetime_greater")
    + 含义：匹配记录字段值表示的时间晚于 < `{Value}`
    + Value格式： `RFC3339` 格式字符串
- OperatorDatetimeGreaterOrEqual ("datetime_greater_or_equal")
    + 含义：匹配记录字段值表示的时间不早于 >= `{Value}`
    + Value格式： `RFC3339` 格式字符串

### 字符串操作符
- OperatorBeginsWith    ("begins_with")
    + 含义：匹配记录字段值是以`{Value}`开头的字符串
    + Value格式：非空字符串
- OperatorNotBeginsWith ("not_begins_with")
    + 含义：匹配记录字段值不是以`{Value}`开头的字符串
    + Value格式：非空字符串
- OperatorContains      ("contains")
    + 含义：匹配记录字段值包含`{Value}`的字符串
    + Value格式：非空字符串
- OperatorNotContains   ("not_contains")
    + 含义：匹配记录字段值不包含`{Value}`的字符串
    + Value格式：非空字符串
- OperatorsEndsWith     ("ends_with")
    + 含义：匹配记录字段值是以`{Value}`结尾的字符串
    + Value格式：非空字符串
- OperatorNotEndsWith   ("not_ends_with")
    + 含义：匹配记录字段值不是以`{Value}`结尾的字符串
    + Value格式：非空字符串

### 空值操作符
- OperatorIsNull    ("is_null")
    + 含义：匹配记录字段值为 `null`
    + Value格式：不接受参数
- OperatorIsNotNull ("is_not_null")
    + 含义：匹配记录字段值不为 `null`
    + Value格式：不接受参数

### 字段存在状态操作符
- OperatorExist    ("exist")
    + 含义：匹配记录包含字段 `{Field}`
    + Value格式：不接受参数
- OperatorNotExist ("not_exist")
    + 含义：匹配记录不包含字段 `{Field}`
    + Value格式：不接受参数

## demo
```json
{
  "condition": "AND",
  "rules": [
    {
      "field": "price",
      "operator": "less",
      "value": 10.25
    },
    {
      "condition": "OR",
      "rules": [
        {
          "field": "category",
          "operator": "equal",
          "value": 2
        },
        {
          "condition": "AND",
          "rules": [{
            "field": "name",
            "operator": "not_in",
            "value": ["a", "b"]
          }]
        }
      ]
    }
  ]
}
```

## TODO
- 考虑是否要提供接口与其它条件合并
    > 一种可选择的方案是，ToMgo之后由用户自行合并
- 如何提取某个字段的过滤条件呢？