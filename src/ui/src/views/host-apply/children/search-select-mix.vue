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
  <bk-search-select
    ref="bkSearchSelect"
    :data="searchOptions"
    :filter="true"
    :filter-menu-method="filterMenuMethod"
    :filter-children-method="filterChildrenMethod"
    :show-condition="false"
    :show-popover-tag-change="false"
    :strink="false"
    v-model.trim="searchValue"
    :placeholder="$t('关键字/字段值')"
    :remote-method="fetchRemoteOptions"
    @input-change="handleInputChange"
    @change="handleChange"
    @menu-select="handleMenuSelect"
    @key-enter="handleKeyEnter"
    @input-focus="handleFocus"
    @input-click-outside="handleBlur">
    <template slot="nextfix">
      <i class="bk-icon icon-close-circle-shape" v-show="showClear && searchValue.length" @click.stop="handleClear"></i>
    </template>
  </bk-search-select>
</template>
<script>
  import TIMEZONE from '@/components/ui/form/timezone.json'
  import Bus from '@/utils/bus'
  import { mapGetters } from 'vuex'
  import has from 'has'
  import { CONFIG_MODE } from '@/service/service-template/index.js'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import useSearchUser from '@/hooks/use-search-user'

  export default {
    props: {
      mode: String
    },
    data() {
      return {
        isTyeing: false,
        showClear: false,
        searchOptions: [],
        fullOptions: [],
        searchValue: [],
        currentMenu: null,
        eventNames: {
          [CONFIG_MODE.MODULE]: 'host-apply-topology-search',
          [CONFIG_MODE.TEMPLATE]: 'host-apply-template-search'
        }
      }
    },
    computed: {
      ...mapGetters('hostApply', ['configPropertyList']),
      searchEventName() {
        return this.eventNames[this.mode]
      }
    },
    watch: {
      mode(newMode, oldMode) {
        // 切换mode时清空搜索
        if (this.searchValue?.length) {
          this.searchValue = []

          // 同时清空切换时所在mode的搜索，实现切换时恢复原数据
          Bus.$emit(this.eventNames[oldMode], this.getSearchValue())
        }
      },
      configPropertyList() {
        this.initOptions()
      },
      searchValue(searchValue) {
        this.searchOptions.forEach((option) => {
          // eslint-disable-next-line max-len
          const selected = searchValue.some(value => value.id === option.id && value.name === option.name && value.type === option.type)
          option.disabled = selected
        })
        this.handleSearch()
      }
    },
    created() {
      this.initOptions()
      this.userSearch = useSearchUser().search
    },
    mounted() {
      Bus.$on('host-apply-clear-search', (value) => {
        this.searchValue = value
        this.clearInputAll()
      })
    },
    methods: {
      async fetchRemoteOptions(val, menu) {
        // 根据具体类型走具体的远程方法拉取options
        const fetchs = {
          [PROPERTY_TYPES.OBJUSER]: this.fetchMember
        }
        return fetchs[menu.type](val, menu)
      },
      async fetchMember(val, menu) {
        if (!this.isTyeing || val?.length < 2 || val === `${menu.name}：`) {
          return []
        }
        const result = await this.userSearch(val)
        return result
      },
      async initOptions() {
        const availableProperties = this.configPropertyList.filter(property => property.host_apply_enabled)
        this.searchOptions = availableProperties.map((property) => {
          const type = property.bk_property_type
          const data = { id: property.id, name: property.bk_property_name, type, disabled: false }
          if (type === PROPERTY_TYPES.ENUM) {
            // eslint-disable-next-line max-len
            data.children = (property.option || []).map(option => ({ id: option.id, name: option.name, disabled: false }))
            data.multiable = true
          } else if (type === PROPERTY_TYPES.LIST) {
            data.children = (property.option || []).map(option => ({ id: option, name: option, disabled: false }))
            data.multiable = true
          } else if (type === PROPERTY_TYPES.TIMEZONE) {
            data.children = TIMEZONE.map(timezone => ({ id: timezone, name: timezone, disabled: false }))
            data.multiable = true
          } else if (type === PROPERTY_TYPES.BOOL) {
            data.children = [{ id: true, name: 'true' }, { id: false, name: 'false' }]
          } else if (type === PROPERTY_TYPES.OBJUSER) {
            data.remote = true
          } else {
            data.children = []
          }
          return data
        })
        this.fullOptions = this.searchOptions.slice(0)
      },
      handleInputChange() {
        this.isTyeing = true
      },
      handleChange(values) {
        const keywords = values.filter(value => !has(value, 'type') && has(value, 'id'))
        if (keywords.length > 1) {
          keywords.pop()
          this.searchValue = values.filter(value => !keywords.includes(value))
        }
      },
      handleKeyEnter() {
        this.currentMenu = null
      },
      handleFocus() {
        this.isTyeing = true
        this.showClear = true
      },
      handleBlur() {
        this.showClear = false
        this.isTyeing = false
      },
      handleClear() {
        this.searchValue = []
        this.clearInputAll()
        Bus.$emit(this.searchEventName, { query_filter: { rules: [] } })
      },
      clearInputAll() {
        this.$refs.bkSearchSelect.inputOutSideClear()
      },
      handleSearch() {
        Bus.$emit(this.searchEventName, this.getSearchValue())
      },
      getSingleTypeVal(type, value) {
        let val = value.id
        switch (type) {
          case PROPERTY_TYPES.ENUM:
            val = value.name
            break
          case PROPERTY_TYPES.OBJUSER:
            val = value.username
            break
          default:
            break
        }
        return val
      },
      getSearchValue() {
        const params = {
          query_filter: {
            condition: 'AND',
            rules: []
          }
        }
        const { rules } = params.query_filter
        this.searchValue.forEach((item) => {
          if (has(item, 'type')) {
            if (item.values.length === 1) {
              // eslint-disable-next-line prefer-destructuring
              const value = item.values[0]
              const isAny = value.id === '*'
              const rule = { field: String(item.id) }
              if (isAny) {
                rule.operator = 'exist'
              } else {
                // 对枚举类型/用户类型特殊处理
                const val = this.getSingleTypeVal(item.type, value)
                rule.operator = 'contains'
                rule.value = String(val).trim()
              }
              rules.push(rule)
            } else {
              const subRule = {
                condition: 'OR',
                rules: []
              }
              item.values.forEach((value) => {
                subRule.rules.push({
                  field: String(item.id),
                  operator: 'contains',
                  value: String(value.id).trim()
                })
              })
              rules.push(subRule)
            }
          } else {
            rules.push({
              field: 'keyword',
              operator: 'contains',
              value: String(item.id).trim()
            })
          }
        })

        return params
      },
      handleMenuSelect(item) {
        this.currentMenu = item
      },
      filterMenuMethod(list, filter) {
        return list.filter(item => item.name.toLowerCase().indexOf(filter.toLowerCase()) > -1)
      },
      filterChildrenMethod(list, filter) {
        if (this.currentMenu && this.currentMenu.children && this.currentMenu.children.length) {
          return this.currentMenu.children.filter(item => item.name.toLowerCase().indexOf(filter.toLowerCase()) > -1)
        }
        return []
      },
      setEmptyKeyword() {
        this.searchValue = []
      }
    }
  }
</script>
<style lang="scss" scoped>
  .icon-close-circle-shape {
    font-size: 14px;
    margin: 0 6px;
    cursor: pointer;
  }
  .icon-search {
    font-size: 16px;
    margin-right: 10px;
    cursor: pointer;
  }
</style>
<style lang="scss">
  .tippy-tooltip.bk-search-select-theme-theme {
    box-shadow: 0 3px 9px 0 rgb(0 0 0 / 10%);
  }
</style>
