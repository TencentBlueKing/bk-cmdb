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

import $http from '@/api'

export interface IPage {
  start: number
  limit: number
  sort?: string
  enable_count?: boolean
}

export interface IRule {
  field: string
  operator: 'not_equal' | 'not_in' | 'equal'
  value: any
}

export interface IFilter {
  condition: 'AND' | 'OR'
  rules: IRule[]
}

export interface IParams<T> {
  params: T
  config?: any
}

const actions = {
  // 获取实列数据
  async findmanyQuotedInstance(ctx, { params, config }: {
    params: {
      bk_obj_id: string
      bk_property_id: string
      filter?: IFilter
      page: IPage
    }
    config?: any
  }) {
    return $http.post('findmany/quoted/instance', params, config).catch(() => ({ info: [], count: 0 }))
  },
  // 创建业务集
  async createTableBizSet(ctx, { params, config }: IParams<any>) {
    return $http.post('table/create/biz_set', params, config)
  },
  // 创建业务
  async createTableBiz(ctx, { params, config }: IParams<any>) {
    // eslint-disable-next-line @typescript-eslint/naming-convention
    const { bk_supplier_account, ...data } = params
    return $http.post(`table/biz/${bk_supplier_account}`, data, config)
  },
  // 创建集合
  async createTableSet(ctx, { params, config }: IParams<any>) {
    return $http.post('', params, config)
  },
  // 创建模块
  async createTableModule(ctx, { params, config }: IParams<any>) {
    return $http.post('', params, config)
  },
  // 创建通用表格字段
  async createTableInstance(ctx, { params, config }: IParams<any>) {
    return $http.post('', params, config)
  },
  // 更新主机
  async updateTableHostsBatch(ctx, { params, config }: {
    params: {bk_host_id: string} & Record<string, any>
    config?: any
  }) {
    return $http.put('table/hosts/batch', params, config).then(() => true)
      .catch(() => false)
  },
  // 更新业务集
  async updateTableBizSet(ctx, { params, config }: IParams<any>) {
    return $http.put('', params, config)
  },
  // 更新业务
  async updateTableBiz(ctx, { params, config }: IParams<any>) {
    return $http.put('', params, config)
  },
  // 更新集合
  async updateTableSet(ctx, { params, config }: IParams<any>) {
    return $http.put('', params, config)
  },
  // 更新模块
  async updateTableModule(ctx, { params, config }: IParams<any>) {
    return $http.put('', params, config)
  },
  // 更新通用表格字段
  async updateTableInstance(ctx, { params, config }: IParams<any>) {
    return $http.put('', params, config)
  }
}

export default {
  namespaced: true,
  actions
}

export {
  actions
}
