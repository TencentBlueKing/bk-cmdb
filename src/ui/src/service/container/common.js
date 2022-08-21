/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { CONTAINER_OBJECTS, WORKLOAD_TYPES } from '@/dictionary/container'
import containerClusterService from '@/service/container/cluster'
import containerNamespaceService from '@/service/container/namespace'
import containerWorkloadService from '@/service/container/workload'

// 根据workload具体类型判断是否为workload
export const isWorkload = type => Object.values(WORKLOAD_TYPES).includes(type)

// 获取容器节点大类型
export const getContainerNodeType = type => (isWorkload(type) ? CONTAINER_OBJECTS.WORKLOAD : type)

// 与传统模型字段类型的映射，原则上在交互形态完全一致的情况下才可以转换
export const typeMapping = {
  string: 'singlechar',
  numeric: 'float',
  mapString: 'map',
  array: 'array',
  object: 'object',
  timestamp: 'time'
}

export const getPropertyType = type => typeMapping[type] || type

export const getPropertyName = (id, objId, locale) => {
  const lang = locale === 'en' ? 'en' : 'zh'
  return propertyNameI18n[objId]?.[id]?.[lang]
}

export const getContainerInstanceService = (objId) => {
  const services = {
    [CONTAINER_OBJECTS.CLUSTER]: containerClusterService,
    [CONTAINER_OBJECTS.NAMESPACE]: containerNamespaceService,
    [CONTAINER_OBJECTS.WORKLOAD]: containerWorkloadService
  }

  return services[objId]
}

export const propertyNameI18n = {
  [CONTAINER_OBJECTS.CLUSTER]: {
    name: {
      zh: '集群名称',
      en: 'clustername'
    },
    scheduling_engine: {
      zh: '调度引擎',
      en: 'enginetype'
    },
    uid: {
      zh: '集群ID',
      en: 'clusterid'
    },
    xid: {
      zh: 'TKE集群ID',
      en: 'systemid'
    },
    version: {
      zh: '集群版本',
      en: 'version'
    },
    network_type: {
      zh: '网络类型',
      en: 'networktype'
    },
    region: {
      zh: '所属地域',
      en: 'region'
    },
    vpc: {
      zh: 'VPC',
      en: 'vpcid'
    },
    network: {
      zh: '集群网络',
      en: 'clusternetwork'
    },
    type: {
      zh: '集群类型',
      en: 'clustertype'
    }
  },
  [CONTAINER_OBJECTS.NAMESPACE]: {
    name: {
      zh: '命名空间名称',
      en: 'Name'
    },
    cluster_uid: {
      zh: '集群ID',
      en: 'clusterid'
    },
    labels: {
      zh: '命名空间标签',
      en: 'Labels'
    },
    resource_quotas: {
      zh: '命名空间资源限制',
      en: 'Resource Quotas'
    }
  },
  [CONTAINER_OBJECTS.WORKLOAD]: {
    name: {
      zh: '工作负载名称',
      en: 'Name'
    },
    namespace: {
      zh: '所属命名空间',
      en: 'Namespace'
    },
    strategy_type: {
      zh: '升级策略',
      en: 'StrategyType'
    },
    labels: {
      zh: '工作负载标签',
      en: 'Labels'
    },
    selector: {
      zh: '工作负载选择器',
      en: 'Selector'
    },
    replicas: {
      zh: '工作负载实例数',
      en: 'Replicas'
    },
    min_ready_seconds: {
      zh: '最小就绪时间',
      en: 'MinReadySeconds'
    },
    rolling_update_strategy: {
      zh: '滚动更新策略',
      en: 'RollingUpdateStrategy'
    }
  },
  [CONTAINER_OBJECTS.NODE]: {
    name: {
      zh: '节点名称',
      en: 'Name'
    },
    roles: {
      zh: '节点角色',
      en: 'Roles'
    },
    labels: {
      zh: '节点标签',
      en: 'Labels'
    },
    creation_timestamp: {
      zh: '节点创建时间',
      en: 'CreationTimestamp'
    },
    taints: {
      zh: '节点污点',
      en: 'Taints'
    },
    unschedulable: {
      zh: '是否可调度',
      en: 'Unschedulable'
    },
    internal_ip: {
      zh: '节点内网IP',
      en: 'InternalIP'
    },
    external_ip: {
      zh: '节点外网IP',
      en: 'ExternalIP'
    },
    hostname: {
      zh: '节点主机名',
      en: 'Hostname'
    },
    runtime: {
      zh: '运行时组件',
      en: 'ContainerRuntime'
    },
    kube_proxy: {
      zh: 'Kube-proxy代理模式',
      en: 'kubeProxy'
    },
    pod_cidr: {
      zh: '节点Pod地址范围',
      en: 'PodCIDR'
    }
  },
  [CONTAINER_OBJECTS.POD]: {
    name: {
      zh: 'Pod名称',
      en: 'Name'
    },
    namespace: {
      zh: '所属命名空间',
      en: 'Namespace'
    },
    priority: {
      zh: 'Pod优先级',
      en: 'Priority'
    },
    node_name: {
      zh: '指定节点调度',
      en: 'NodeName'
    },
    start_time: {
      zh: 'Pod启动时间',
      en: 'StartTime'
    },
    labels: {
      zh: 'Pod标签',
      en: 'Labels'
    },
    ip: {
      zh: 'Pod容器网络IP',
      en: 'IP'
    },
    ips: {
      zh: 'Pod容器网络IPs',
      en: 'IPs'
    },
    controlled_by: {
      zh: '所属副本控制器',
      en: 'ControlledBy'
    },
    containers: {
      zh: 'Pod包含容器',
      en: 'Containers'
    },
    qos_class: {
      zh: 'Pod服务质量',
      en: 'QoSClass'
    },
    volumes: {
      zh: 'Pod卷信息',
      en: 'Volumes'
    },
    node_selectors: {
      zh: '将Pod指派给节点',
      en: 'Node-Selectors'
    },
    tolerations: {
      zh: 'Pod污点',
      en: 'Tolerations'
    }
  },
  [CONTAINER_OBJECTS.CONTAINER]: {
    name: {
      zh: '名称',
      en: 'Name'
    },
    container_uid: {
      zh: '容器ID',
      en: 'Container ID'
    },
    image: {
      zh: '镜像信息',
      en: 'Image'
    },
    ports: {
      zh: '容器端口',
      en: 'Ports'
    },
    host_ports: {
      zh: '主机端口映射',
      en: 'Host Ports'
    },
    args: {
      zh: '启动参数',
      en: 'Args'
    },
    started: {
      zh: '启动时间',
      en: 'Started'
    },
    limits: {
      zh: '资源限制',
      en: 'Limits'
    },
    requests: {
      zh: '申请资源大小',
      en: 'Requests'
    },
    liveness: {
      zh: '存活探针',
      en: 'Liveness'
    },
    environment: {
      zh: '环境变量',
      en: 'Environment'
    },
    mounts: {
      zh: '挂载卷',
      en: 'Mounts'
    }
  }
}
