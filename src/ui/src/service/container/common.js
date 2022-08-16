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
  object: 'object'
}

export const getPropertyType = type => typeMapping[type] || type

export const getPropertyName = (id, objId, locale) => {
  const lang = locale === 'en' ? 'en' : 'zh'
  return propertyNameI18n[objId]?.[id]?.[lang]
}

export const propertyNameI18n = {
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
  }
}
