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

import Vue from 'vue'
const store = new Vue({
  data() {
    return {
      hosts: [],
      businessList: []
    }
  },
  computed: {
    isSelected() {
      return !!this.hosts.length
    },

    bizSet() {
      const bizSet = new Set()
      this.hosts.forEach((host) => {
        const [biz] = host.biz
        bizSet.add(biz.bk_biz_id)
      })
      return bizSet
    },

    isSameBiz() {
      return this.bizSet.size === 1
    },

    isAllIdleModule() {
      return this.hosts.every((host) => {
        const [module] = host.module
        return module.default === 1
      })
    },

    isAllIdleSet() {
      return this.hosts.every((host) => {
        const [module] = host.module
        return module.default !== 0
      })
    },

    uniqueBusiness() {
      if (this.isSameBiz) {
        const [bizId] = Array.from(this.bizSet)
        return this.businessList.find(business => business.bk_biz_id === bizId)
      }
      return null
    },

    isAllResourceHost() {
      return this.hosts.every((host) => {
        const [biz] = host.biz
        return biz.default === 1
      })
    }
  },
  methods: {
    clear() {
      this.hosts = []
    },

    setSelected(hosts = []) {
      this.hosts = hosts
    },

    getSelected() {
      return this.hosts
    },

    setBusinessList(businessList) {
      this.businessList = businessList
    }
  }
})

export default store
