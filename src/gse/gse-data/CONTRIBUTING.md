# 开发指南

开发前需先阅读了解项目的代码规范、日志规范、代码提交规范。

## 代码规范

### 目录文件命名规范

- 目录命名以横线(-)全小写形式，如third-party, 涉及第三方和跨平台依赖时可适当使用下划线格式;
- 文件命名以下划线(_)全小写形式，如xxxx_xxxx.h、xxxx_xxxx.c、xxxx_xxxx.hpp、xxxx_xxxx.cpp;

### 头文件多重包含保护

使用宏定义防止头文件被多重包含，命名格式: `_{项目名}_{模块名}_{子模块名(若有)}_{文件名}_{头文件后缀}_`,

`头文件后缀`: 为区分不同头文件后缀，故此需对应定义为`_XXXX_H_` 或 `_XXXX_HPP_`

示例:

```c
#ifndef _GSE_MAIN_NET_TCPCLI_H_
#define _GSE_MAIN_NET_TCPCLI_H_

    ......

#endif // _GSE_MAIN_NET_TCPCLI_H_
```

### 头文件的引用次序

将包含次序标准化可增强可读性、避免隐藏依赖。建议标准包含次序：`当前子模块头文件、C标准库、C++标准库、第三方库、项目内其他子模块头文件`

`注意`: 每个头文件区块必须使用换行分隔开，便于查看和代码格式化, 禁止将各种类型头文件混合到一个代码块内, 区块之间顺序可根据具体场景微调.

示例:

```c
// 当前子模块头文件
#include "local.hpp"

// C标准库
#include <unistd.h>
#include <stdlib.h>
#include <memory.h>

// C++标准库
#include <iostream>
#include <vector>
#include <string>

// 第三方库
#include <event/event.h>
#include <event/thread.h>
#include <event/buffer.h>

// 项目内其他子模块头文件
#include "common/common.hpp"
```

`注意`: 为保证唯一性, 头文件的命名应该基于所在项目源代码树的全路径, 避免使用UNIX特殊的快捷目录`.`(当前目录) 或 `..`(上级目录)

### 类型、函数、变量命名

- `类名`: 采用`大驼峰样式`, 即每个单词的首字母`大写`, 如`class BaseImpl;`;
- `成员函数`: 采用`大驼峰样式`, 即每个单词的首字母`大写`, 如`void FuncName();`;
- `成员变量`: 采用`m_`开头(C语言中struct可不加任何前缀)，采用`小驼峰样式`, 即单词的首字母`小写`, 其他单次首字母`大写`, 如`m_baseName;`;
- `普通变量`: 采用`小驼峰样式`, 如`int statusCode = 0;`;
- `常变量`: 采用`k`开头，后接`大驼峰样式`, 如`const int kDaysInAWeek = 7;`;
- `宏定义`: 采用`全大写样式`, 单词间用下划线相连, 如`#define MSG_SIZE 256`;
- `枚举定义`: 采用`全大写样式`, 单词间用下划线相连, 如`enum ProtoType { PT_TCP = 0 }`;
- `全局变量`: 采用`g_`开头(尽量不要使用全局变量), 之后采用`小驼峰样式`, 如`g_serviceName`;
- `命名空间`: 采用`全小写样式`, 如`namespace gse`;

### 注释

采用Doxgen支持的文档注释语法

[Documenting The Code](https://www.doxygen.nl/manual/docblocks.html#cppblock)

### 代码格式化插件

项目基于`clang-format`进行统一格式化，以此来约束代码结构，以下为常见IDE配置：

```json
{
    "Language":"Cpp",
    "BasedOnStyle":"LLVM",
    "UseTab":"Never",
    "TabWidth":4,
    "IndentWidth":4,
    "ColumnLimit":0,
    "MaxEmptyLinesToKeep":1,
    "AccessModifierOffset":-4,
    "IndentCaseLabels":"false",
    "FixNamespaceComments":"true",
    "DerivePointerAlignment":"true",
    "PointerAlignment":"Left",
    "BreakBeforeBraces":"Custom",
    "SpacesInAngles":"false",
    "AllowShortFunctionsOnASingleLine":"Inline",
    "BraceWrapping":{
        "AfterUnion":"true",
        "AfterStruct":"true",
        "AfterClass":"true",
        "AfterEnum":"true",
        "AfterFunction":"true",
        "AfterControlStatement":"true",
        "BeforeCatch":"true",
        "BeforeElse":"true",
        "AfterNamespace":"false"
    }
}
```

`注意`: SpacesInAngles配置项为可选项，默认关闭，根据需要可设置为true在模板参数列表<>中增加空格，用于解决低版本C++编译问题

## 日志规范

- 日志内容必须具备可读性，杜绝日志中出现无法理解的神秘数字、短语及缩写等非行业内通用的表达方式，必须以非系统开发人员可以理解的表达方式进行编写;
- 日志格式有利于日志管理系统收集，检索及分析;
- 日志内容采用英文语句编写，单词采用小写;
- 日志内容禁止出现有利于破解的消息，尽量避免日志中输出函数名称，变量名称等，方便破解软件的信息;
- 日志内容禁止出现敏感信息，包含但不限于用户名，密码，系统秘钥内容等;

### 日志格式

日志字段说明:

 * time: 打印日志时的时间戳，精确到毫秒级别;
 * level: 日志级别，可取值debug，info，warn，error，fatal, 系统日志推荐级别为info;
 * service: 产生日志的服务模块;
 * file: 输出日志的文件信息;
 * line: 代码行号;
 * process: 进程ID;
 * thread: 线程ID;
 * requestid: 用于对消息进行定位跟踪;
 * message：日志消息体，此处只能是字符串，采用英文语言，全小写表达，内容应该符合5W1H的规则, 即即 ***when、who、where、what、why、context***;

示例:

```text
{"time":"2020-07-27 14:20:25.128","level":"error","service":"task","process":0,"thread":1, "requestid":"1234567890","file":"dispatcher.cpp","line":666,"message":"the dispatcher can not find any worker to use, the total worker count is (10), already used (10), need to scale it"}
```

### 日志分类存储

- 系统日志【必须】：GSE系统自身运行的日志，如系统启动，定时检查 等系统自身运行行为相关的日志;
- 业务日志【必须】：API调用信息，业务逻辑信息，执行业务逻辑涉及的数据等（在满足保密和性能要求的前提下存储）;
- 第三方库导出的日志【按需】：第三方库提供日志接口导出的日志，此类日志按组件名字分文件存储，标准输出的日志如果无法区分日志的所属就统一到独立的日志文件存储;

最佳实践, 细分三个子目录存储日志:

* business：API 调用信息，业务逻辑信息，执行业务逻辑涉及的数据等（在满足保密和性能要求的前提下存储）, 采用宏`BUSINESS_LOG`打印;
* system：GSE系统自身运行的日志，如系统启动，定时检查等系统自身运行行为相关的日志, 采用宏`LOG_XXX`（DEBUG,INFO,ERROR,WARN)打印;
* stdout：通过标准输打印的日志，为开发人员使用, 采用宏`DEV_LOG` 打印;

## 代码提交规范

### 分支与Issues

请确保关键性的问题都关联其对应的Issues，并确保该Issues是否已经存在, 尽量将相同话题在一个Issue处理或进行关联，不要冗余创建Issues话题。

* 关键特性需要关联Issue;
* 代码提交需在单独的特性分支不要在master分支提交;

### Commit信息格式

Each commit message should include a **type**, a **scope**, a **subject** and a **issue**:

Commit信息需包含**type**, **scope**（若需要）, **subject**, **issue**(若有):

```
 <type>(<scope>): <subject>. issue #num
```

Commit信息不要超过100个字符，尽量保证内容的可读性, 比较好的提交示例如下,

```
 #8 FEAT(template): add new template. issue #202// 该提交说明新增了一个模板文件支持，type为'FEAT', scope为'template', subject描述了具体改动, 关联issue #202
 #7 FIX(dockerfile): fix an issue with source image in Dockerfile. issue #201 // 该提交说明修复了一个镜像问题，type为'FIX', scope为'dockerfile', subject描述了具体改动, 关联issue #201
 #6 DOCS(project): update available templates section. issue #200 // 该提交说明更新了项目文档，type为'DOCS', scope为'project', subject描述了具体改动, 关联issue #200
```

#### Commit Type类型说明

Commit Type必须是下面类型中的一个:

* **FEAT**: 新特性
* **FIX**: 修复问题
* **OPT**: 优化
* **DOCS**: 仅仅文档更新
* **STYLE**: 修改不影响代码逻辑的格式问题
* **REFACTOR**: 既不是新特性也不是修复问题的重构
* **TEST**: 更新测试相关
* **CHORE**: 修改构建、部署等相关内容

#### Commit Scope说明

Scope为可选，该信息需简短说明改动提交的关联内容， 如`project`, `etc`...

#### Commit Subject说明

Subject需简短准确描述改动内容:

* 准确适用时态语义: "change" not "changed" nor "changes"
* 首单次字母不要大写
* 结尾不要使用中英文句号，读个简短描述可以使用;分隔，但不建议单个Commit信息提交过多

## Protobuf协议规范

[grpc-gateway protobuf style guide](https://buf.build/docs/style-guide/#files-and-packages)

