## how it design

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
过滤规则解析方法，从map[string]interface{}数据中解析出一个过滤规则实例

## Operator 详细说明
### 通用操作符
- OperatorEqual    ("equal")
- OperatorNotEqual ("not_equal")

### 数组操作符
- OperatorIn    ("in")
    + 含义：匹配记录字段值是否在指定集合中
- OperatorNotIn ("not_in")
    + 含义：匹配记录字段值不在指定集合中
- OperatorIsEmpty    ("is_empty")
    + 含义：匹配记录字段值为空数组
- OperatorIsNotEmpty ("is_not_empty")
    + 含义：匹配记录字段值为非空数组
	
### 数字操作符
- OperatorLess           ("less")
    + 含义：匹配记录字段值 < `{Value}`
- OperatorLessOrEqual    ("less_or_equal")
    + 含义：匹配记录字段值 <= `{Value}`
- OperatorGreater        ("greater")
    + 含义：匹配记录字段值 > `{Value}`
- OperatorGreaterOrEqual ("greater_or_equal")
    + 含义：匹配记录字段值 >= `{Value}`

### 时间操作符
- OperatorDatetimeLess           ("datetime_less")
    + 含义：匹配记录字段值表示的时间早于 < `{Value}`
- OperatorDatetimeLessOrEqual    ("datetime_less_or_equal")
    + 含义：匹配记录字段值表示的时间不晚于 <= `{Value}`
- OperatorDatetimeGreater        ("datetime_greater")
    + 含义：匹配记录字段值表示的时间晚于 < `{Value}`
- OperatorDatetimeGreaterOrEqual ("datetime_greater_or_equal")
    + 含义：匹配记录字段值表示的时间不早于 >= `{Value}`

### 字符串操作符
- OperatorBeginsWith    ("begins_with")
    + 含义：匹配记录字段值是以`{Value}`开头的字符串
- OperatorNotBeginsWith ("not_begins_with")
    + 含义：匹配记录字段值不是以`{Value}`开头的字符串
- OperatorContains      ("contains")
    + 含义：匹配记录字段值包含`{Value}`的字符串
- OperatorNotContains   ("not_contains")
    + 含义：匹配记录字段值不包含`{Value}`的字符串
- OperatorsEndsWith     ("ends_with")
    + 含义：匹配记录字段值是以`{Value}`结尾的字符串
- OperatorNotEndsWith   ("not_ends_with")
    + 含义：匹配记录字段值不是以`{Value}`结尾的字符串


### 空值操作符
- OperatorIsNull    ("is_null")
    + 含义：匹配记录字段值为 `null`
- OperatorIsNotNull ("is_not_null")
    + 含义：匹配记录字段值不为 `null`

### 字段存在状态操作符
- OperatorExist    ("exist")
    + 含义：匹配记录包含字段 `{Field}`
- OperatorNotExist ("not_exist")
    + 含义：匹配记录不包含字段 `{Field}`

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
          "field": "category",
          "operator": "equal",
          "value": 1
        }
      ]
    }
  ]
}
```