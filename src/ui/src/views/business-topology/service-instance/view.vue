<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="view">
    <component :is="activeComponent"></component>
  </div>
</template>

<script>
  import ViewInstance from './instance/view'
  import ViewProcess from './process/view'
  import RouterQuery from '@/router/query'
  export default {
    components: {
      ViewInstance,
      ViewProcess
    },
    data() {
      return {
        activeComponent: null,
        viewMap: Object.freeze({
          instance: ViewInstance.name,
          process: ViewProcess.name
        })
      }
    },
    created() {
      this.unwatchView = RouterQuery.watch('view', this.handleViewChange, { immediate: true })
      this.unwatchTab = RouterQuery.watch('tab', this.handleTabChange)
    },
    beforeDestroy() {
      this.unwatchView()
      this.unwatchTab()
    },
    methods: {
      handleViewChange(view = 'instance') {
        this.activeComponent = this.viewMap[view]
      },
      handleTabChange(tab) {
        if (tab !== 'serviceInstance') {
          this.activeComponent = null
        } else {
          const view = RouterQuery.get('view', 'instance')
          RouterQuery.set({ view })
          this.activeComponent = this.viewMap[view]
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .view {
        position: relative;
        padding: 15px 0;
    }
</style>
