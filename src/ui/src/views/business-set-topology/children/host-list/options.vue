<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="options-layout clearfix">
    <div class="options options-left fl">
      <cmdb-clipboard-selector class="options-clipboard"
        v-test-id
        label-key="bk_property_name"
        :list="clipboardList"
        :disabled="!hasSelection"
        @on-copy="handleCopy">
      </cmdb-clipboard-selector>
    </div>
    <div class="options options-right">
      <filter-fast-search class="option-fast-search" v-test-id></filter-fast-search>
      <filter-collection class="option-collection ml10" v-test-id></filter-collection>
      <icon-button :class="['option-filter', 'ml10', { active: hasCondition }]" v-test-id="'advancedSearch'"
        icon="icon-cc-funnel" v-bk-tooltips.top="$t('高级筛选')"
        @click="handleSetFilters">
      </icon-button>
    </div>
  </div>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import FilterForm from '@/components/filters/filter-form.js'
  import FilterCollection from '@/components/filters/filter-collection'
  import FilterFastSearch from '@/components/filters/filter-fast-search'
  import FilterStore from '@/components/filters/store'
  import FilterUtils from '@/components/filters/utils'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import { IPWithCloudSymbol, IPWithCloudFields, IPv6WithCloudSymbol, IPv46WithCloudSymbol, IPv64WithCloudSymbol } from '@/dictionary/ip-with-cloud-symbol'
  import { isEmptyPropertyValue } from '@/utils/tools'

  export default {
    components: {
      FilterCollection,
      FilterFastSearch,
    },
    computed: {
      ...mapGetters('userCustom', ['usercustom']),
      ...mapState('bizSet', ['bizId']),
      ...mapGetters('businessHost', ['selectedNode']),
      hostProperties() {
        return FilterStore.getModelProperties(BUILTIN_MODELS.HOST)
      },
      count() {
        return this.$parent.table.pagination.count
      },
      selection() {
        return this.$parent.table.selection
      },
      hasSelection() {
        return !!this.selection.length
      },
      clipboardList() {
        const IPWithClouds = Object.keys(IPWithCloudFields).map(key => FilterUtils.defineProperty({
          id: key,
          bk_obj_id: 'host',
          bk_property_id: key,
          bk_property_name: IPWithCloudFields[key],
          bk_property_type: 'singlechar'
        }))
        const clipboardList = this.$parent.tableHeader.slice()
        clipboardList.splice(3, 0, ...IPWithClouds)
        return clipboardList
      },
      tableHeaderPropertyIdList() {
        return this.$parent.tableHeader.map(item => item.bk_property_id)
      },
      hasCondition() {
        return FilterStore.hasCondition
      }
    },
    methods: {
      handleCopy(property) {
        const copyText = this.selection.map((data) => {
          const modelId = property.bk_obj_id
          const modelData = data[modelId]

          const IPWithCloudKeys = Object.keys(IPWithCloudFields)
          if (IPWithCloudKeys.includes(property.id)) {
            const cloud = this.$tools.getPropertyCopyValue(modelData.bk_cloud_id, 'foreignkey')
            const ip = this.$tools.getPropertyCopyValue(modelData.bk_host_innerip, 'singlechar')
            const ipv6 = this.$tools.getPropertyCopyValue(modelData.bk_host_innerip_v6, 'singlechar')
            const isEmptyIPv4Value = isEmptyPropertyValue(modelData.bk_host_innerip)
            const isEmptyIPv6Value = isEmptyPropertyValue(modelData.bk_host_innerip_v6)
            if (property.id === IPWithCloudSymbol) {
              return `${cloud}:${ip}`
            }
            if (property.id === IPv6WithCloudSymbol) {
              return `${cloud}:[${ipv6}]`
            }
            if (property.id === IPv46WithCloudSymbol) {
              if (!isEmptyIPv4Value || isEmptyIPv6Value) {
                return `${cloud}:${ip}`
              }
              return `${cloud}:[${ipv6}]`
            }
            if (property.id === IPv64WithCloudSymbol) {
              if (isEmptyIPv4Value || !isEmptyIPv6Value) {
                return `${cloud}:[${ipv6}]`
              }
              return `${cloud}:${ip}`
            }
          }

          const propertyId = property.bk_property_id
          if (Array.isArray(modelData)) {
            const value = modelData.map(item => this.$tools.getPropertyCopyValue(item[propertyId], property))
            return value.join(',')
          }
          return this.$tools.getPropertyCopyValue(modelData[propertyId], property)
        })
        this.$copyText(copyText.join('\n')).then(() => {
          this.$success(this.$t('复制成功'))
        }, () => {
          this.$error(this.$t('复制失败'))
        })
      },
      handleSetFilters() {
        FilterForm.show()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .options-layout {
        margin-top: 12px;
    }
    .options {
        display: flex;
        align-items: center;
        &.options-right {
            overflow: hidden;
            justify-content: flex-end;
        }
        .option {
            display: inline-block;
            vertical-align: middle;
        }
        .option-fast-search {
            flex: 1;
            max-width: 300px;
            margin-left: 10px;
        }
        .option-collection,
        .option-filter {
            flex: 32px 0 0;
            &:hover,
            .active {
                color: $primaryColor;
            }
        }
        .dropdown-icon {
            margin: 0 -4px;
            display: inline-block;
            vertical-align: middle;
            height: auto;
            top: 0px;
            font-size: 20px;
            &.open {
                top: -1px;
                transform: rotate(180deg);
            }
        }
    }
    .bk-dropdown-list {
        font-size: 14px;
        color: $textColor;
        .bk-dropdown-item {
            position: relative;
            display: block;
            padding: 0 20px;
            margin: 0;
            line-height: 32px;
            cursor: pointer;
            @include ellipsis;
            &:not(.disabled):not(.with-auth):hover {
                background-color: #EAF3FF;
                color: $primaryColor;
            }
            &.disabled {
                color: $textDisabledColor;
                cursor: not-allowed;
            }
            &.with-auth {
                padding: 0;
                span {
                    display: block;
                    padding: 0 20px;
                    &:not(.disabled):hover {
                        background-color: #EAF3FF;
                        color: $primaryColor;
                    }
                    &.disabled {
                        color: $textDisabledColor;
                        cursor: not-allowed;
                    }
                }
            }
        }
    }
    /deep/ {
        .collection-item {
            width: 100%;
            display: flex;
            justify-content: space-between;
            align-items: center;
            &:hover {
                .icon-close {
                    display: block;
                }
            }
            .collection-name {
                @include ellipsis;
            }
            .icon-close {
                display: none;
                color: #979BA5;
                font-size: 20px;
                margin-right: -4px;
            }
        }
    }
</style>
