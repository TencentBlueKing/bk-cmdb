## 概述
一个用于将前端/API的通用查询参数解析为mongodb过滤条件的后端模块。该通用查询参数满足以下条件：
- 仅支持查询相应数据类型所支持的字段，且用于查询的值和字段类型需要匹配，防止使用错误的参数进行查询返回无数据造成误解。如主机只可以用主机的字段进行查询，且主机名只能用string类型查询，使用其它字段或不匹配的类型进行查询会报错
- 支持查询的字段类型包括基本类型（数值/bool值/字符串/时间）和复杂结构（数组和对象）
- 支持字段查询参数的组合，如多个字段的查询条件的 `与` 、 `或` 关系和查询嵌套结构体的字段等
- 兼容之前版本的 [querybuilder查询参数](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md) ，便于后续迁移

## 格式
通用查询条件为数据的属性字段过滤规则的组合，用于根据属性字段搜索数据。该参数为以下两种过滤规则类型，可以嵌套。具体支持的过滤规则如下：

### 组合过滤规则
由其它规则组合而成的过滤规则，组合的节点间支持逻辑与/或关系

| 字段        | 类型     | 必选  | 描述                              |
|-----------|--------|-----|---------------------------------|
| condition | string | 是   | 组合查询条件，支持 `AND` 和 `OR` 两种方式     |
| rules     | array  | 是   | 查询规则，可以是 `组合查询参数` 或 `原子查询参数` 类型 |

### 原子过滤规则
基础的过滤规则，表示用于对某一个字段进行过滤的规则。任何过滤规则都直接是原子过滤规则, 或由多个原子过滤规则组合而成

| 名称       | 类型                            | 必填  | 说明  |
|----------|-------------------------------|-----|-----|
| field    | string                        | 是   | 字段名 |
| operator | string                        | 是   | 操作符 | 
| value    | 不同的field和operator对应不同的value格式 | 否   | 操作数 |

#### operator 字段详细说明
##### 通用操作符
- equal
    + 含义：匹配字段值和`value`相等的数据
    + value格式：基本数据类型(数值/bool值/字符串)
- not_equal
    + 含义：匹配字段值和`value`不等的的数据
    + value格式：基本数据类型(数值/bool值/字符串)

##### 集合操作符
- in
  + 含义：匹配字段值在`value`指定的集合中的数据
  + value格式：基本数据类型(数值/bool值/字符串)组成的数组，类型需要一致
- not_in
  + 含义：匹配字段值不在`value`指定的集合中的数据
  + value格式：基本数据类型(数值/bool值/字符串)组成的数组，类型需要一致

##### 数字操作符
> 支持的field包括数值类型和时间戳类型的字段

- less
    + 含义：匹配字段值小于`value`的数据
    + value格式：数值
- less_or_equal
    + 含义：匹配字段值小于或等于`value`的数据
    + value格式：数值
- greater
    + 含义：匹配字段值大于`value`的数据
    + value格式：数值
- greater_or_equal
    + 含义：匹配字段值大于或等于`value`的数据
    + value格式：数值

##### 时间操作符
> 支持field为时间类型的字段，不包括时间戳类型的字段

> 目前value支持时间戳格式和cc时间格式（示例："2006-01-02 15:04:05" 和 "2006-01-02T15:04:05+08:00"）

- datetime_less
    + 含义：匹配字段值表示的时间早于`value`的数据
    + value格式：时间戳格式的数值类型和cc时间格式的字符串类型
- datetime_less_or_equal
    + 含义：匹配字段值表示的时间早于或等于`value`的数据
    + value格式：时间戳格式的数值类型和cc时间格式的字符串类型
- datetime_greater
    + 含义：匹配字段值表示的时间晚于`value`的数据
    + value格式：时间戳格式的数值类型和cc时间格式的字符串类型
- datetime_greater_or_equal
    + 含义：匹配字段值表示的时间晚于或等于`value`的数据
    + value格式：时间戳格式的数值类型和cc时间格式的字符串类型

##### 字符串操作符
- begins_with
    + 含义：匹配字段值是以`value`开头的字符串的数据，该操作符大小写敏感
    + value格式：非空字符串
- begins_with_insensitive
    + 含义：匹配字段值是以`value`开头的字符串的数据，该操作符大小写不敏感
    + value格式：非空字符串
- not_begins_with
    + 含义：匹配字段值不是以`value`开头的字符串的数据，该操作符大小写敏感
    + value格式：非空字符串
- not_begins_with_insensitive
    + 含义：匹配字段值不是以`value`开头的字符串的数据，该操作符大小写不敏感
    + value格式：非空字符串
- contains
    + 含义：匹配字段值包含`value`的字符串的数据。注：为了兼容querybuilder的contains操作符，该操作符大小写不敏感
    + value格式：非空字符串
- contains_sensitive
  + 含义：匹配字段值包含`value`的字符串的数据，该操作符大小写敏感
  + value格式：非空字符串
- not_contains
    + 含义：匹配字段值不包含`value`的字符串的数据，该操作符大小写敏感
    + value格式：非空字符串
- not_contains_insensitive
    + 含义：匹配字段值不包含`value`的字符串的数据，该操作符大小写不敏感
    + value格式：非空字符串
- ends_with
    + 含义：匹配字段值是以`value`结尾的字符串的数据，该操作符大小写敏感
    + value格式：非空字符串
- ends_with_insensitive
    + 含义：匹配字段值是以`value`结尾的字符串的数据，该操作符大小写不敏感
    + value格式：非空字符串
- not_ends_with
    + 含义：匹配字段值不是以`value`结尾的字符串的数据，该操作符大小写敏感
    + value格式：非空字符串
- not_ends_with_insensitive
    + 含义：匹配字段值不是以`value`结尾的字符串的数据，该操作符大小写不敏感
    + value格式：非空字符串

##### 数组操作符
- is_empty
  + 含义：匹配字段值是空数组的数据
  + value格式：不接受参数
- is_not_empty
  + 含义：匹配字段值是非空数组的数据
  + value格式：不接受参数
- size
  + 含义：匹配字段值是长度为`value`的数组的数据
  + value格式：数值

##### 空值操作符
- is_null
    + 含义：匹配字段值是`null`的数据
    + value格式：不接受参数
- is_not_null
  + 含义：匹配字段值不是`null`的数据
    + value格式：不接受参数

##### 字段存在状态操作符
- exist
    + 含义：匹配包含字段`field`的数据
    + value格式：不接受参数
- not_exist
    + 含义：匹配不包含字段`field`的数据
    + value格式：不接受参数

##### 复杂结构操作符
- filter_object
    + 含义：匹配字段值满足`value`对应的过滤规则的数据
    + value格式：过滤object类型字段值中的子字段的过滤规则
- filter_array
    + 含义：匹配不包含字段`field`的数据
    + value格式：过滤array类型字段值中的元素的过滤规则，其下层级的原子过滤条件的`field`支持用`element`表示匹配任意一个数组元素，用数组下标表示匹配指定元素

## 示例
- 查询条件示例：
``` json
{
    "condition": "and",
    "rules": [
        {
            "field": "test",
            "operator": "equal",
            "value": 1
        },
        {
            "condition": "or",
            "rules": [
                {
                    "field": "test1",
                    "operator": "filter_array",
                    "value": {
                        "condition": "and",
                        "rules": [
                            {
                                "field": "element",
                                "operator": "filter_object",
                                "value": {
                                    "field": "test2",
                                    "operator": "not_equal",
                                    "value": "a"
                                }
                            },
                            {
                                "field": "element",
                                "operator": "filter_object",
                                "value": {
                                    "condition": "and",
                                    "rules": [
                                        {
                                            "field": "test2",
                                            "operator": "in",
                                            "value": [
                                                "b",
                                                "c"
                                            ]
                                        },
                                        {
                                            "field": "test3",
                                            "operator": "not_equal",
                                            "value": false
                                        }
                                    ]
                                }
                            }
                        ]
                    }
                },
                {
                    "field": "test4",
                    "operator": "datetime_less",
                    "value": 100
                }
            ]
        }
    ]
}
```

- 用上述条件查询出的数据示例：
``` json
{
    "test": 1,
    "test1": [
        {
            "test2": "d",
            "test3": false
        },
        {
            "test2": "b",
            "test3": true
        }
    ],
    "test4": "2006-01-02 15:04:05"
}
```
