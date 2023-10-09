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

import i18n from '@/i18n'

function t(key) {
  return i18n.t(key)
}
export const actions = {
  create: t('创建'),
  update: t('修改'),
  delete: t('删除'),
  archive: t('归档'),
  recover: t('还原'),
  assign_host: t('分配到业务'),
  unassign_host: t('归还主机'),
  transfer_host_module: t('转移模块')
}

export const types = {
  business: {
    business: t('业务')
  },
  business_resource: {
    service_template: t('服务模板'),
    set_template: t('集群模板'),
    service_category: t('服务分类'),
    dynamic_group: t('动态分组'),
    model_instance: t('模型实例'),
    set: t('集群'),
    module: t('模块'),
    process: t('进程')
  },
  host: {
    host: t('主机')
  },
  model: {
    model: t('模型'),
    model_association: t('模型关联'),
    model_group: t('模型分组'),
    model_unique: t('模型唯一校验')
  },
  model_instance: {
    model_instance: t('实例'),
    instance_association: t('关联关系')
  },
  association_kind: {
    association_kind: t('关联类型')
  },
  cloud_resource: {
    cloud_area: t('管控区域'),
    cloud_account: t('云账户'),
    cloud_sync_task: t('云同步任务')
  },
  biz_topology: {
    biz_topology: t('业务拓扑')
  },
  service_instance: {
    service_instance: t('服务实例')
  }
}
