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
  <div class="options">
    <div class="left"></div>
    <div class="right">
      <bk-checkbox class="options-expand-all" v-model="expandAll" @change="handleExpandAllChange">
        {{$t('全部展开')}}
      </bk-checkbox>
      <bk-input class="options-search ml10"
        ref="searchSelect"
        v-model.trim="searchValue"
        right-icon="bk-icon icon-search"
        clearable
        :max-width="200"
        :placeholder="$t('请输入进程别名')"
        @enter="handleSearch"
        @clear="handleSearch">
      </bk-input>
      <view-switcher :show-tips="false" class="ml10" active="process"></view-switcher>
    </div>
  </div>
</template>

<script>
  import ViewSwitcher from '@/views/business-topology/service-instance/common/view-switcher'
  import Bus from '@/views/business-topology/service-instance/common/bus'
  import { mapGetters } from 'vuex'
  export default {
    components: {
      ViewSwitcher
    },
    data() {
      return {
        withTemplate: true,
        searchData: [],
        searchValue: '',
        expandAll: false,
        selection: {
          process: null,
          value: [],
          requestId: null
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['selectedNode']),
      serviceTemplateId() {
        return this.selectedNode && this.selectedNode.data.service_template_id
      }
    },
    created() {
      Bus.$on('process-selection-change', this.handleProcessSelectionChange)
      Bus.$on('process-list-change', this.handleProcessListChange)
    },
    beforeDestroy() {
      Bus.$off('process-selection-change', this.handleProcessSelectionChange)
      Bus.$off('process-list-change', this.handleProcessListChange)
    },
    methods: {
      handleProcessSelectionChange(process, selection, requestId) {
        if (selection.length) {
          this.selection.process = process
          this.selection.value = selection
          this.selection.requestId = requestId
        } else if (process === this.selection.process) {
          this.selection.process = null
          this.selection.value = []
          this.selection.requestId = null
        }
      },
      handleSearch() {
        Bus.$emit('filter-list', this.searchValue)
      },
      handleExpandAllChange(expand) {
        Bus.$emit('expand-all-change', expand)
      },
      handleProcessListChange() {
        this.expandAll = false
        this.selection = {
          process: null,
          value: [],
          requestId: null
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .options {
        display: flex;
        justify-content: space-between;
        flex-wrap: wrap;
        .left,
        .right {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
        }
    }
    .options-search {
        width: 300px;
    }
</style>
