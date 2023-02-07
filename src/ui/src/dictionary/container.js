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

import { BUILTIN_MODELS } from './model-constants'

export const CONTAINER_OBJECTS = Object.freeze({
  CLUSTER: 'cluster',
  NAMESPACE: 'namespace',
  FOLDER: 'folder',
  WORKLOAD: 'workload',
  POD: 'pod',
  CONTAINER: 'container',
  NODE: 'node'
})

export const CONTAINER_OBJECT_NAMES = Object.freeze({
  [CONTAINER_OBJECTS.CLUSTER]: {
    FULL: 'Cluster',
    SHORT: 'C'
  },
  [CONTAINER_OBJECTS.NAMESPACE]: {
    FULL: 'Namespace',
    SHORT: 'Ns'
  },
  [CONTAINER_OBJECTS.FOLDER]: {
    FULL: 'Folder',
    SHORT: 'F'
  },
  [CONTAINER_OBJECTS.WORKLOAD]: {
    FULL: 'Workload',
    SHORT: 'Wl'
  },
  [CONTAINER_OBJECTS.POD]: {
    FULL: 'Pod',
    SHORT: 'P'
  },
  [CONTAINER_OBJECTS.CONTAINER]: {
    FULL: 'Container',
    SHORT: 'Cn'
  },
  [CONTAINER_OBJECTS.NODE]: {
    FULL: 'Node',
    SHORT: 'Nd'
  }
})

export const CONTAINER_OBJECT_LEVELS = Object.freeze({
  [CONTAINER_OBJECTS.CLUSTER]: {
    NEXT: CONTAINER_OBJECTS.NAMESPACE,
    PREV: BUILTIN_MODELS.BUSINESS
  },
  [CONTAINER_OBJECTS.NAMESPACE]: {
    NEXT: CONTAINER_OBJECTS.WORKLOAD,
    PREV: CONTAINER_OBJECTS.CLUSTER
  },
  [CONTAINER_OBJECTS.FOLDER]: {
    NEXT: BUILTIN_MODELS.HOST,
    PREV: CONTAINER_OBJECTS.CLUSTER
  },
  [CONTAINER_OBJECTS.WORKLOAD]: {
    NEXT: BUILTIN_MODELS.HOST,
    PREV: CONTAINER_OBJECTS.NAMESPACE
  }
})

export const WORKLOAD_TYPES = Object.freeze({
  DEPLOYMENT: 'deployment',
  STATEFUL_SET: 'statefulSet',
  DAEMON_SET: 'daemonSet',
  GAME_STATEFUL_SET: 'gameStatefulSet',
  GAME_DEPLOYMENT: 'gameDeployment',
  CRON_JOB: 'cronJob',
  JOB: 'job',
  PODS: 'pods'
})

export const WORKLOAD_OBJECT_NAMES = Object.freeze({
  [WORKLOAD_TYPES.DEPLOYMENT]: {
    FULL: 'Deployment',
    SHORT: 'Dp'
  },
  [WORKLOAD_TYPES.STATEFUL_SET]: {
    FULL: 'StatefulSet',
    SHORT: 'Ss'
  },
  [WORKLOAD_TYPES.DAEMON_SET]: {
    FULL: 'DaemonSet',
    SHORT: 'Ds'
  },
  [WORKLOAD_TYPES.GAME_STATEFUL_SET]: {
    FULL: 'Game StatefulSet',
    SHORT: 'Gs'
  },
  [WORKLOAD_TYPES.GAME_DEPLOYMENT]: {
    FULL: 'Game Deployments',
    SHORT: 'Gd'
  },
  [WORKLOAD_TYPES.CRON_JOB]: {
    FULL: 'CronJob',
    SHORT: 'Cj'
  },
  [WORKLOAD_TYPES.JOB]: {
    FULL: 'Job',
    SHORT: 'J'
  },
  [WORKLOAD_TYPES.PODS]: {
    FULL: 'Pod',
    SHORT: 'P'
  }
})

export const TOPO_MODE_KEYS = Object.freeze({
  CONTAINER: 'container',
  BIZ_NODE: 'bizNode',
  NORMAL: 'normal',
  NONE: 'none'
})

export const MIX_SEARCH_MODES = Object.freeze({
  LIKE_CONTAINER: 'likeContainer',
  LIKE_NORMAL: 'likeNormal',
  UNKNOW: 'unknow'
})

export const CONTAINER_OBJECT_PROPERTY_KEYS = Object.freeze({
  [CONTAINER_OBJECTS.CLUSTER]: {
    ID: 'bk_cluster_id',
    NAME: 'bk_cluster_name'
  },
  [CONTAINER_OBJECTS.NAMESPACE]: {
    ID: 'bk_namespace_id',
    NAME: 'bk_namespace_name'
  },
  [CONTAINER_OBJECTS.FOLDER]: {
    ID: 'bk_folder_id',
    NAME: 'bk_folder_name'
  },
  [CONTAINER_OBJECTS.WORKLOAD]: {
    ID: 'bk_workload_id',
    NAME: 'bk_workload_name'
  },
  [CONTAINER_OBJECTS.POD]: {
    ID: 'bk_pod_id',
    NAME: 'bk_pod_name'
  },
  [CONTAINER_OBJECTS.CONTAINER]: {
    ID: 'bk_container_id',
    NAME: 'bk_container_name'
  }
})

export const CONTAINER_OBJECT_INST_KEYS = Object.freeze({
  [CONTAINER_OBJECTS.CLUSTER]: {
    ID: 'id',
    NAME: 'name'
  },
  [CONTAINER_OBJECTS.NAMESPACE]: {
    ID: 'id',
    NAME: 'name'
  },
  [CONTAINER_OBJECTS.FOLDER]: {
    ID: 'id',
    NAME: 'name'
  },
  [CONTAINER_OBJECTS.WORKLOAD]: {
    ID: 'id',
    NAME: 'name'
  },
  [CONTAINER_OBJECTS.POD]: {
    ID: 'id',
    FULL_ID: `${CONTAINER_OBJECTS.POD}_id`,
    NAME: 'name',
    FULL_NAME: `${CONTAINER_OBJECTS.POD}_name`,
  },
  [CONTAINER_OBJECTS.NODE]: {
    ID: 'id',
    FULL_ID: `${CONTAINER_OBJECTS.NODE}_id`,
    NAME: 'name',
    FULL_NAME: `${CONTAINER_OBJECTS.NODE}_name`
  },
  [CONTAINER_OBJECTS.CONTAINER]: {
    ID: 'id',
    FULL_ID: `${CONTAINER_OBJECTS.CONTAINER}_id`,
    NAME: 'name',
    FULL_NAME: `${CONTAINER_OBJECTS.CONTAINER}_name`
  }
})
