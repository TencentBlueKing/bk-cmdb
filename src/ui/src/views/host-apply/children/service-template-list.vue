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
  <div class="service-template-list" :style="{ '--height': `${$APP.height - 160 - 44}px` }" ref="$templateList">
    <div :class="['template-item', { selected: item.id === current.id, disabled: isNodeDisabled(item) }]"
      :title="isNodeDisabled(item) ? $t('暂无策略') : ''"
      v-for="item in displayList" :key="item.id"
      @click="handleClickItem(item)">
      <bk-checkbox
        v-if="options.showCheckbox"
        :disabled="isNodeDisabled(item)"
        :value="checked[item.id]">
      </bk-checkbox>
      <i class="node-icon">{{modelIcon}}</i>
      <span class="node-name">{{item.name}}</span>
      <span class="config-icon" v-if="item.host_apply_enabled">
        <i class="bk-cc-icon icon-cc-selected"></i>
      </span>
    </div>
    <bk-exception class="empty" :type="isSearch ? 'search-empty' : 'empty'" scene="part" v-if="!displayList.length">
      <span>{{$t('无数据')}}</span>
    </bk-exception>
  </div>
</template>
<script>
  import { defineComponent, reactive, toRefs, watch, watchEffect, set, del, onActivated, computed, onBeforeUnmount, nextTick, ref, onMounted } from '@vue/composition-api'
  import store from '@/store'
  import Bus from '@/utils/bus'
  import router from '@/router/index.js'
  import { sortTopoTree } from '@/utils/tools.js'
  import serviceTemplateService from '@/services/service-template/index.js'

  export default defineComponent({
    props: {
      options: {
        type: Object,
        default: () => ({})
      },
      action: {
        type: String,
        default: 'batch-edit'
      }
    },
    setup(props, { emit }) {
      const { options, action } = toRefs(props)
      const bizId = store.getters['objectBiz/bizId']
      const getModelById = store.getters['objectModelClassify/getModelById']

      const $templateList = ref(null)

      const state = reactive({
        fullList: [],
        displayList: [],
        modelIcon: getModelById('module').bk_obj_name[0],
        current: {},
        isSearch: false
      })
      const checked = reactive({})

      const targetId = computed(() => Number(router.app.$route.query.id))

      watchEffect(async () => {
        const list = await serviceTemplateService.findAll({ bk_biz_id: bizId })

        sortTopoTree(list ?? [], 'name')

        state.fullList = list
        state.displayList = state.fullList

        const current = targetId.value ? state.fullList.find(item => item.id === targetId.value) : state.fullList[0]
        state.current = current ?? {}

        const counts = await store.dispatch('hostApply/getTemplateRuleCount', {
          params: {
            bk_biz_id: bizId,
            service_template_ids: state.fullList.map(item => item.id)
          }
        })
        state.displayList.forEach((node) => {
          node.host_apply_rule_count = counts.find(item => item.service_template_id === node.id)?.count
        })
      })

      const handleClickItem = (item) => {
        if (isNodeDisabled(item)) {
          return
        }

        if (options.value.checkOnClick) {
          set(checked, item.id, !checked[item.id])
          return
        }

        state.current = item
      }

      // 当前点选项变化时抛出事件
      watch(() => state.current, current => emit('selected', { data: current }))

      watch(checked, (checked) => {
        const checkedIds = []
        Object.keys(checked).forEach((key) => {
          if (checked[key]) checkedIds.push(Number(key))
        })
        emit('checked', checkedIds)
      }, { deep: true })

      const isDel = computed(() => action.value === 'batch-del')

      // 提供给批量编辑取消选择的方法
      const removeChecked = (key) => {
        let keys = Object.keys(checked)
        if (key) {
          keys = !Array.isArray(key) ? [key] : key
        }

        keys.forEach((key) => {
          del(checked, key)
        })
      }

      // 更新节点状态，用于关闭应用后状态更新
      const updateNodeStatus = (id, { isClose, isClear }) => {
        const nodeData = state.displayList.find(item => item.id === Number(id))
        if (isClose) {
          nodeData.host_apply_enabled = false
        }
        if (isClear) {
          nodeData.host_apply_rule_count = 0
        }
      }

      const isNodeDisabled = node => isDel.value && !node.host_apply_rule_count

      const scrollSelectedIntoView = () => {
        $templateList.value?.querySelector('.template-item.selected')?.scrollIntoView()
      }

      onActivated(() => {
        const current = state.displayList.find(item => item.id === targetId.value)
        if (current) {
          state.current = current ?? {}
          emit('selected', { data: current })

          nextTick(scrollSelectedIntoView)
        }
      })

      const handleSearch = async (condition) => {
        state.isSearch = true
        try {
          if (condition.query_filter.rules.length) {
            const keywordRuleIndex = condition.query_filter.rules.findIndex(item => item.field === 'keyword')

            if (keywordRuleIndex === -1) {
              // 不存在关键字则只采用接口搜索
              const data = await store.dispatch('hostApply/searchTemplateNode', { params: { bk_biz_id: bizId, ...condition } })
              filterList({ remote: data || [] })
            } else {
              // 先取出关键字的参数值
              const keyword = condition.query_filter.rules[keywordRuleIndex]?.value

              // 接口搜索需要去掉keyword参数
              condition.query_filter.rules.splice(keywordRuleIndex, 1)

              // 两者都存在，混合搜索
              if (condition.query_filter.rules.length) {
                const data = await store.dispatch('hostApply/searchTemplateNode', { params: { bk_biz_id: bizId, ...condition } })
                filterList({
                  keyword,
                  remote: data || [],
                })
              } else {
                // 仅存在关键字搜索
                filterList({ keyword })
              }
            }
          } else {
            // 清空搜索
            filterList({ clear: true })
          }
        } catch (e) {
          console.error(e)
        }
      }

      const filterList = ({ remote: remoteData, keyword, clear }) => {
        if (clear) {
          state.displayList = state.fullList.slice()
          state.isSearch = false
          return
        }

        if (remoteData && keyword) {
          const ids = remoteData.map(item => item.id)
          const reg = new RegExp(keyword, 'i')
          state.displayList = state.fullList.filter(item => ids.includes(item.id) && reg.test(item.name))
        }

        if (remoteData) {
          const ids = remoteData.map(item => item.id)
          state.displayList = state.fullList.filter(item => ids.includes(item.id))
        }

        if (keyword) {
          const reg = new RegExp(keyword, 'i')
          state.displayList = state.fullList.filter(item => reg.test(item.name))
        }
      }

      Bus.$on('host-apply-template-search', handleSearch)

      onBeforeUnmount(() => {
        Bus.$off('host-apply-template-search', handleSearch)
      })

      onMounted(() => {
        setTimeout(scrollSelectedIntoView, 1000)
      })

      return {
        ...toRefs(state),
        checked,
        isDel,
        isNodeDisabled,
        removeChecked,
        updateNodeStatus,
        handleClickItem,
        $templateList
      }
    }
  })
</script>

<style lang="scss" scoped>
  .service-template-list {
    height: var(--height);
    padding: 10px 0;
    margin-right: 2px;
    @include scrollbar-y;

    &::-webkit-scrollbar {
      background: none;
      &-thumb {
        border-radius: 3px;
      }
    }

    .template-item {
      display: flex;
      align-items: center;
      height: 36px;
      margin: 2px 0;
      padding: 0 10px;
      cursor: pointer;

      &:hover {
        background: #f1f7ff;
      }

      &.selected {
        background: #e1ecff;
      }

      &.disabled {
        cursor: not-allowed;
      }

      .node-name {
        flex: 1;
        font-size: 14px;
        margin-right: 8px;
        @include ellipsis;
      }

      .node-icon {
        width: 22px;
        height: 22px;
        line-height: 21px;
        text-align: center;
        font-style: normal;
        font-size: 12px;
        margin: 0 8px 0 6px;
        border-radius: 50%;
        background-color: #c4c6cc;
        color: #fff;
        background-color: #97aed6;

        &.is-selected {
          background-color: #3a84ff;
        }
      }

      .config-icon {
        padding: 0 5px;
        text-align: center;

        .icon-cc-selected {
          font-size: 26px;
          color: #2dcb56;
        }
      }
    }

    .empty {
      ::v-deep {
        .bk-exception-img {
          width: 160px;
          height: 140px;
        }
        .bk-exception-text {
          font-size: 12px;
        }
      }
    }
  }
</style>
