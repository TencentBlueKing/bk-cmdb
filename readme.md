![](docs/resource/img/BKCN.jpg)
---
[![license](https://img.shields.io/badge/license-mit-brightgreen.svg?style=flat)](https://github.com/Tencent/bk-cmdb/blob/master/ LICENSE)
[![Release Version](https://img.shields.io/badge/release-3.1.0-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/releases)
[![Build Status](https://travis-ci.org/Tencent/bk-cmdb.svg?branch=master)](https://travis-ci.org/Tencent/bk-cmdb)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/pulls)

[(English Documents Available)](readme_en.md)

> **Important Note**: The `master` branch may be in an unstable or unavailable state during development*.
Please use [releases] (https://github.com/tencent/bk-cmdb/releases) instead of `master` to get stable binaries.

Blue Whale Configuration Platform (Blue Whale CMDB) is an enterprise configuration management service based on operation and maintenance scenarios.

This open source is a new refactored version 3.0. Compared with 2.0, it provides a new custom model management. Users can not only easily expand the built-in model attributes, but also add models and associations at any time according to different enterprise needs, and integrate network, middleware, virtual resources, etc. Management of the CMDB. In addition, it adds more new features that meet the needs of the scene: machine data snapshots, automatic data discovery, active event change notifications, more granular rights management, and scalable business topologies.

In terms of technology construction, the core of the new version focuses on resources. We divide the atomic resources managed by CMDB into three types: host, process and general object, and construct atomic operation layers for these resources. Above these atomic operations, we built a layer of the scene closer to the user's operation, and the scene layer completes the user's request by combining operations on different resources.

This refactoring uses golang as the development language. Compared with the 2.0 version, the system's operating efficiency has been greatly improved. In addition, the microservice architecture design is adopted, and the deployment of the system can support the traditional mode and the container mode.

The open source version will be consistent with the blue whale configuration platform version built into the Blue Whale Community Edition and Enterprise Edition and will be updated simultaneously. Our goal is to create a unified configuration management platform that is compatible with different industries and different architectures. It is the industry's leading free and open source CMDB with good versatility and ease of use. We welcome interested colleagues to participate.



## Overview
* [Architecture Design] (docs/overview/architecture.md)
* [code directory] (docs/overview/code_framework.md)
* [Design Concept] (docs/overview/design.md)
* [Usage Scenario] (docs/overview/usecase.md)

## Features
* Topology host management: host basic attributes, host snapshot data, host affiliation management
* Organizational structure management: scalable business-based organizational structure management
* Model management: can manage built-in models such as business, cluster, host, etc., as well as custom models.
* Process Management: Module-based host process management
* Event registration and push: Provide callback-based event registration and push
* Universal Rights Management: Flexible User Group Based Rights Management
* Operational audit: auditing and backtracking of user operational behavior

For a detailed description of the above functions, please refer to [Function Description] (http://bk.tencent.com/document/bkprod/000120.html)

## Experience

[Experience Blue Whale cmdb] (docs/overview/experience.md)

## Getting started
* [Download and Compile] (docs/overview/source_compile.md)
* [Installation Deployment] (docs/overview/installation.md)
* [API Usage Notes] (docs/apidoc/readme.md)

## Version plan
* [Version Iteration Rule] (docs/VERSION.md)

## Support
1. Refer to the bk-cmdb installation documentation [Installation Documentation] (docs/overview/installation.md)
2. Read [source code] (https://github.com/Tencent/bk-cmdb/tree/master)
3. Read [wiki] (https://github.com/Tencent/bk-cmdb/wiki/cmdb-3.0) or ask for help
4. Learn about the Blue Whale community: [Blue Whale Community Edition Exchange 1 Group] (https://jq.qq.com/?_wv=1027&k=5zk8F7G)
5. Contact us, technical exchange QQ group:

![qq](docs/resource/img/qq.png)

## Contributing
For the bk-cmdb branch management, issues, and pr specifications, read the [bk-cmdb Contributing Guide] (docs/CONTRIBUTING.md).

## FAQ

https://github.com/Tencent/bk-cmdb/wiki/FAQ

## License
Bk-cmdb is based on the MIT protocol. Please refer to [LICENSE](LICENSE) for details.
