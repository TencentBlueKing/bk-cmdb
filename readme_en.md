![](docs/resource/img/BKEN.png)
---
[![license](https://img.shields.io/badge/license-mit-brightgreen.svg?style=flat)](https://github.com/Tencent/bk-cmdb/blob/master/LICENSE)
[![Release Version](https://img.shields.io/badge/release-3.0.8-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/releases)
[![Build Status](https://travis-ci.org/Tencent/bk-cmdb.svg?branch=master)](https://travis-ci.org/Tencent/bk-cmdb)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/pulls)

> **Note**: The `master` branch may be in an *unstable or even broken state* during development.
To get stable binaries, please use [releases](https://github.com/tencent/bk-cmdb/releases) instead of the `master` branch.

BlueKing Configuration Management DataBase (bk-cmdb) is an enterprise level configuration management service database. 

bk-cmdb provides users a complete new way to self-define model management. Not only can users easily expand the built-in model attributes, but also add models and associations at any time according to different enterprise requirements, and incorporate networks, middleware, and virtual resources into the bk-cmdb. It also adds more new features that meet the needs of certain requirements. For example: machine data snapshots, automatic data discovery, active push of event changes, more fine-grained permission management, and scalable service topology, etc. 

The main focus of the new release is resources. We divide the atomic resources managed by CMDB into hosts, processes, and generic objects, and build an atomic operation layer on these resources. On top of these atomic operations, we built a scenario layer that is closer to user operations. The scenario layer completes users’ requests through a combination of operations of different resources.

bk-cmdb uses golang. It is of high performance and easy to develop. In addition, it adopts microservice architecture design, which has strong scalability, is easily monitored and supports smooth upgrade. Besides traditional methods, deployment through Docker is supported.

Source code of BlueKing’s Community Edition and Enterprise Edition is consistent and synchronized upon changes. Our goal is to create a unified configuration management platform that is compatible with different industries and different architectures, and to become industry-leading free and open-source CMDB with good versatility and ease of use. We welcome participation of interested developers.

## Overview
* [Architecture Design (In Chinese)](docs/overview/architecture.md)
* [Code Directory (In Chinese)](docs/overview/code_framework.md)
* [Design Philosophy (In Chinese)](docs/overview/design.md)
* [Use Case (In Chinese)](docs/overview/usecase.md)

## Features
* Topological host management: basic attributes, snapshot data, ownership management
* Organizational Structure Management: Scalable Business-Based Organizational Structure Management
* Model management: Management of business, cluster, host and other built-in models, and customizable model management.
* Process Management: Module-based host process management
* Event registration and push: provide callback-based event registration and push
* Universal Rights Management: Flexible User Group Based Permission Management
* Operation Audit: Auditing and Backtracking of User Operational Behavior

If you want to know more about the above features, please refer to the [Feature Description (In Chinese)](http://bk.tencent.com/document/bkprod/000120.html)

## Experience
[Hands-on Training of BlueKing CMDB Docker Deployment](docs/wiki/container-support.md)

## Getting started
* [Download and Compilation (In Chinese)](docs/overview/source_compile.md)
* [Installation and Deployment (In Chinese)](docs/overview/installation.md)
* [API Instructions (In Chinese)](docs/apidoc/readme.md)

## Version plan
* [Version iteration rules (In Chinese)](docs/VERSION.md)

## Support
1. Refer to the bk-cmdb installation document [Installation Documentation (In Chinese)](docs/overview/installation.md)
2. Read [source (In Chinese)](https://github.com/Tencent/bk-cmdb/tree/master)
3. Read the wiki (In Chinese)(https://github.com/Tencent/bk-cmdb/wiki/cmdb-3.0) or ask for help
4. Learn about BlueKing Community related information: [Blue Whale Community Edition 1 Group](https://jq.qq.com/?_wv=1027&k=5zk8F7G)
5. Contact us, technical exchange QQ group:

![qq](docs/resource/img/qq.png)

## Contributing
For bk-cmdb branch management, issues, and pr specifications, read the [bk-cmdb Contributing Guide (In Chinese)](docs/CONTRIBUTING.md).
[Tencent Open Source Incentive Plan](https://opensource.tencent.com/contribution) aims to encourage developers’ participation and contribution. Welcome everybody.

## FAQ
https://github.com/Tencent/bk-cmdb/wiki/FAQ

## License
Bk-cmdb is based on the MIT License. Please refer to [LICENSE](LICENSE) for details.
