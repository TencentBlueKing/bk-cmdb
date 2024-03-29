# Bug提单
如果你在使用cmdb的过程中幸运的发现了一个bug，社区欢迎你通过[提交issue](https://github.com/TencentBlueKing/bk-cmdb/issues/new)
将这个bug反馈给社区来进行定位和修复。针对你提交的bug，还需要你提供一些具体的信息来帮助社区来定位、修复这个问题。

社区需要你提供下面这些基础信息：
* 问题本身的具体内容，可以是一段日志，也可以是web页面上弹出的提示截图等。
* 你使用的版本，也即版本号。可以对某个组件执行version命令来获取，如 `./cmdb_apiserver --version`。
* 问题复现的具体流程。

针对复杂的问题，如果还需要提供其它更为详细的信息，还请你能够及时提供。你提供的信息越多，我们就能越快的定位、解决问题。

为了方便大家提issue，社区提供一个issue的模板帮助大家梳理这些信息。示例如下：

```markdown
问题描述
===========
<这里写问题描述>


重现方法
================
<列出如何重现的方法或操作步骤>


**重要提醒**: 请优先尝试最新发布的版本 (发布清单： https://github.com/TencentBlueKing/bk-cmdb/releases), 如果问题不能在最新发布的版本里重现，说明此问题已经被修复。


关键信息
=========

**重要提醒**: 这些关键信息会辅助我们快速定位问题。

请提供以下信息:

 - [x] bk-cmdb   版本 (发布版本号 或 git tag): `<示例： v3.0.6-alpha 或者 git sha. 请不要使用 "最新版本" 或 "当前版本"等无法准确定位代码版本的语句描述>`
 - [ ] Redis     版本: `<示例： 3.2.11>`
 - [ ] MongoDB   版本: `<示例： 2.8.0>`
 - [ ] ZooKeeper 版本: `<示例： 3.4.11>`
 - [ ] 操作系统      : `<示例： Centos 5 (x64)>`
 - [ ] bk-cmdb 异常日志
```

如果你不确定你发现的是Bug还是Feature，也可以直接提issue，我们来一起看一下。

## 解决
一旦你提交的问题被社区确定为是一个bug，修复的方式有很多种：
  - 如果你对cmdb的源码很熟悉，可直接将bug修复后，将修改的代码以PR的方式提交给社区，在经过review通过后，合入主库。
  - 如果你对Golang比较熟悉，想自己修复这个bug，但不知道怎么修复。没关系，社区可以进行指导，帮助你修复这个问题，最终也是通过PR的方式合入主库。
  - 由社区来修复这个bug。

bug修复后，你可以基于这个分支自行编译代码再部署使用。也可以等社区release新的版本后再使用。
