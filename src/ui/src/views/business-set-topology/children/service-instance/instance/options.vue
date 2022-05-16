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
    <div class="left" v-test-id="'moreAction'">
      <bk-button :disabled="!selection.length" @click="handleCopy(!selection.length)">
        {{$t('复制IP')}}
      </bk-button>
    </div>
    <div class="right">
      <bk-checkbox class="options-expand-all" v-model="allExpanded" @change="handleExpandAll">
        {{$t('全部展开')}}
      </bk-checkbox>
      <bk-search-select class="options-search ml10"
        ref="searchSelect"
        :show-condition="false"
        :placeholder="$t('请输入实例名称或选择标签')"
        :data="searchMenuList"
        v-model="searchValue"
        @change="handleSearch">
      </bk-search-select>
      <view-switcher :show-tips="false" class="ml10" active="instance"></view-switcher>
    </div>
    <span class="tips-content" ref="tipsContent">
      <i18n path="当前模块主机为空" tag="p" v-if="!hasHost"></i18n>
      <i18n path="当前模块所属的服务模板内进程信息为空" tag="p" v-else-if="isEmptyServiceTemplate"></i18n>
    </span>
  </div>
</template>

<script>
  import ViewSwitcher from '@/views/business-topology/service-instance/common/view-switcher'
  import Bus from '@/views/business-topology/service-instance/common/bus'
  import RouterQuery from '@/router/query'
  import { mapGetters, mapState } from 'vuex'
  import { Validator } from 'vee-validate'
  import { MULTIPLE_IP_REGEXP } from '@/dictionary/regexp.js'
  import { ServiceInstanceService } from '@/service/business-set/service-instance.js'
  import { ProcessTemplateService } from '@/service/business-set/process-template.js'
  export default {
    components: {
      ViewSwitcher
    },
    data() {
      return {
        selection: [],
        hasProcessTemplate: false,
        allExpanded: false,
        historyLabels: {},
        searchValue: [],
        request: {
          label: Symbol('label')
        }
      }
    },
    computed: {
      ...mapGetters('businessHost', ['selectedNode']),
      ...mapState('bizSet', ['bizSetId', 'bizId']),
      withTemplate() {
        return this.selectedNode && this.selectedNode.data.service_template_id
      },
      hasHost() {
        return this.selectedNode && this.selectedNode.data?.host_count > 0
      },
      nameFilterIndex() {
        return this.searchValue.findIndex(data => data.id === 'name')
      },
      searchMenuList() {
        const hasHistoryLables = Object.keys(this.historyLabels).length > 0
        const list = [{
          id: 'name',
          name: this.$t('服务实例名')
        }, {
          id: 'tagValue',
          name: this.$t('标签值'),
          conditions: Object.keys(this.historyLabels).map(key => ({ id: key, name: `${key}:` })),
          disabled: !hasHistoryLables
        }, {
          id: 'tagKey',
          name: this.$t('标签键'),
          disabled: !hasHistoryLables
        }]
        if (this.nameFilterIndex > -1) {
          return list.slice(1)
        }
        return list.slice(0)
      },
      isEmptyServiceTemplate() {
        return this.withTemplate && !this.hasProcessTemplate
      }
    },
    watch: {
      withTemplate: {
        immediate: true,
        handler(withTemplate) {
          if (withTemplate) {
            this.checkProcessTemplate()
          }
        }
      },
      searchMenuList() {
        this.$nextTick(() => {
          const menu = this.$refs.searchSelect && this.$refs.searchSelect.menuInstance
          menu && (menu.list = this.$refs.searchSelect.data)
        })
      },
      selectedNode() {
        this.searchValue = []
        RouterQuery.set({
          node: this.selectedNode.id,
          instanceName: ''
        })
        Bus.$emit('filter-change', this.searchValue)
      }
    },
    created() {
      this.unwatch = RouterQuery.watch(['node', 'page'], () => {
        this.allExpanded = false
      })
      Bus.$on('instance-selection-change', this.handleInstanceSelectionChange)
      this.setFilter()
      this.updateHistoryLabels()
    },
    beforeDestroy() {
      this.unwatch()
      Bus.$off('instance-selection-change', this.handleInstanceSelectionChange)
    },
    methods: {
      updateHistoryLabels() {
        ServiceInstanceService.findAggregationLabels(this.bizSetId, {
          bk_biz_id: this.bizId
        })
          .then((data) => {
            this.historyLabels = data
          })
          .catch(() => {
            this.historyLabels = []
          })
      },
      setFilter() {
        const filterName = RouterQuery.get('instanceName')
        if (!filterName) {
          return false
        }
        this.searchValue.push({
          id: 'name',
          name: this.$t('服务实例名'),
          values: [{
            id: filterName,
            name: filterName
          }]
        })
        this.$nextTick(() => {
          // list.vue注册的监听晚于派发，因此nextTick再触发
          Bus.$emit('filter-change', this.searchValue)
        })
      },
      handleInstanceSelectionChange(selection) {
        this.selection = selection
      },
      async handleCopy(disabled) {
        if (disabled) {
          return false
        }
        try {
          const validator = new Validator()
          const validPromise = []
          this.selection.forEach((row) => {
            // eslint-disable-next-line prefer-destructuring
            const ip = row.name.split('_')[0]
            validPromise.push(new Promise(async (resolve) => {
              const { valid } = await validator.verify(ip, MULTIPLE_IP_REGEXP)
              resolve({ valid, ip })
            }))
          })
          const results = await Promise.all(validPromise)
          const validResult = results.filter(result => result.valid).map(result => result.ip)
          const unique = [...new Set(validResult)]
          if (unique.length) {
            await this.$copyText(unique.join('\n'))
            this.$success(this.$t('复制成功'))
          } else {
            this.$warn(this.$t('暂无可复制的IP'))
          }
        } catch (e) {
          console.error(e)
          this.$error(this.$t('复制失败'))
        }
      },
      checkProcessTemplate() {
        ProcessTemplateService.findAll(this.bizSetId, {
          page: { start: 0, limit: 1 },
          bk_biz_id: this.bizId,
          service_template_id: this.selectedNode.data.service_template_id
        })
          .then((data) => {
            this.hasProcessTemplate = data.length > 0
          })
          .catch(() => {
            this.hasProcessTemplate = false
          })
      },
      handleSearch(filters) {
        const transformedFilters = []
        filters.forEach((data) => {
          if (!data.values) {
            const nameIndex = transformedFilters.findIndex(filter => filter.id === 'name')
            if (nameIndex > -1) {
              transformedFilters[nameIndex].values = [{ ...data }]
            } else {
              transformedFilters.push({
                id: 'name',
                name: this.$t('服务实例名'),
                values: [{ ...data }]
              })
            }
          } else if (data.id === 'tagValue') {
            const [{ name }] = data.values
            if (data.condition && name) {
              transformedFilters.push(data)
            } else {
              this.$warn(this.$t('服务实例标签值搜索提示语'))
            }
          } else {
            transformedFilters.push(data)
          }
        })
        this.$nextTick(() => {
          this.searchValue = transformedFilters
          Bus.$emit('filter-change', this.searchValue)
        })
      },
      handleExpandAll(expanded) {
        Bus.$emit('expand-all-change', expanded)
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
    .menu-list {
        .menu-item {
            line-height: 32px;
            .menu-option {
                display: block;
                padding: 0 10px;
                font-size: 14px;
                cursor: pointer;
                &:hover {
                    color: $primaryColor;
                }
                &.disabled {
                    color: $textDisabledColor;
                    cursor: not-allowed;;
                }
            }
        }
    }
    .options-sync {
        display: inline-block;
        position: relative;
        margin-left: 18px;
        padding: 0;
        &::before {
            content: '';
            position: absolute;
            top: 7px;
            left: -11px;;
            width: 1px;
            height: 20px;
            background-color: #dcdee5;
        }
        .sync-wrapper {
            display: flex;
            align-items: center;
        }
        .icon-refresh {
            font-size: 12px;
            margin-right: 4px;
        }
        .topo-status {
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            background-color: #ea3636;
            border-radius: 50%;
        }
    }
    .options-search {
        width: 200px;
    }

    .tips-content {
      display: none;
      .action-link {
        color: $primaryColor;
        margin-left: 2px;
      }
    }
</style>
