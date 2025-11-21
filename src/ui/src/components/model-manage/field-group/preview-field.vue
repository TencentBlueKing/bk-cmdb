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
  <div class="form-layout">
    <div class="form-groups" ref="formGroups">
      <template v-for="(group, groupIndex) in $sortedGroups">
        <div class="property-group"
          :key="groupIndex"
          v-if="checkGroupAvailable(groupedProperties[groupIndex])">
          <cmdb-collapse
            :label="group['bk_group_name']"
            :collapse.sync="groupState[group['bk_group_id']]">
            <ul class="property-list">
              <template v-for="(property, propertyIndex) in groupedProperties[groupIndex]">
                <li :class="['property-item', { 'full-width': isTableType(property) }]"
                  v-if="checkEditable(property)"
                  :key="propertyIndex">
                  <div class="property-name">
                    <span class="property-name-text" :class="{ required: property['isrequired'] }" v-bk-overflow-tips>
                      {{property['bk_property_name']}}
                    </span>
                    <i class="property-name-tooltips icon-cc-tips"
                      v-if="property['placeholder']"
                      v-bk-tooltips="{
                        content: htmlEncode(property['placeholder'])
                      }">
                    </i>
                  </div>
                  <div class="property-value clearfix">
                    <slot :name="property.bk_property_id">
                      <table-default-settings
                        v-if="isTableType(property)"
                        :preview="true"
                        :headers="property.option.header"
                        :defaults="property.option.default" />
                      <component class="form-component" v-else
                        :is="`cmdb-form-${property['bk_property_type']}`"
                        :class="{ error: errors.has(property['bk_property_id']) }"
                        :placeholder="$tools.getPropertyPlaceholder(property)"
                        :unit="property['unit']"
                        :row="2"
                        :multiple="property.ismultiple"
                        :disabled="checkDisabled(property)"
                        :options="property.option || []"
                        :data-vv-name="property['bk_property_id']"
                        :data-vv-as="property['bk_property_name']"
                        v-bind="$tools.getValidateEvents(property)"
                        v-validate="getValidateRules(property)"
                        v-model.trim="values[property['bk_property_id']]">
                      </component>
                      <span class="form-error"
                        :title="errors.first(property['bk_property_id'])">
                        {{errors.first(property['bk_property_id'])}}
                      </span>
                    </slot>
                  </div>
                </li>
              </template>
            </ul>
          </cmdb-collapse>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
  import formMixins from '@/mixins/form'
  import { filterXSS } from '@/utils/util'
  import TableDefaultSettings from './table-default-settings.vue'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'

  export default {
    components: {
      TableDefaultSettings
    },
    mixins: [formMixins],
    props: {
      inst: {
        type: Object,
        default() {
          return {}
        }
      },
      objId: {
        type: String,
        default: ''
      },
      type: {
        default: 'create',
        validator(val) {
          return ['create', 'update'].includes(val)
        }
      }
    },
    data() {
      return {
        values: {},
        refrenceValues: {},
        groupState: {
          none: true
        }
      }
    },
    computed: {
      groupedProperties() {
        return this.$groupedProperties.map(properties => properties.filter(property => !['singleasst', 'multiasst', 'foreignkey'].includes(property.bk_property_type)))
      }
    },
    watch: {
      inst() {
        this.initValues()
      },
      properties() {
        this.initValues()
      }
    },
    created() {
      this.initValues()
    },
    methods: {
      initValues() {
        this.values = this.$tools.getInstFormValues(this.properties, this.inst, this.type === 'create')
        this.refrenceValues = this.$tools.clone(this.values)
      },
      checkGroupAvailable(properties) {
        const availabelProperties = properties.filter(property => this.checkEditable(property))
        return !!availabelProperties.length
      },
      checkEditable(property) {
        if (this.type === 'create') {
          return !property.bk_isapi
        }
        return property.editable && !property.bk_isapi
      },
      checkDisabled(property) {
        if (this.type === 'create') {
          return false
        }
        return !property.editable || property.isreadonly
      },
      htmlEncode(placeholder) {
        return filterXSS(placeholder)
      },
      getPlaceholder(property) {
        const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
        return this.$t(placeholderTxt, { name: property.bk_property_name })
      },
      getValidateRules(property) {
        return this.$tools.getValidateRules(property)
      },
      isTableType(property) {
        return property.bk_property_type === PROPERTY_TYPES.INNER_TABLE
      },
      uncollapseGroup() {
        this.errors.items.forEach((item) => {
          const property = this.properties.find(property => property.bk_property_id === item.field)
          const group = property.bk_property_group
          this.groupState[group] = false
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .form-groups {
        padding: 0 24px 20px;
    }
    .property-group {
        padding: 7px 0 10px 0;
        &:first-child {
            padding: 28px 0 10px 0;
        }
    }
    .group-name {
        font-size: 14px;
        line-height: 14px;
        color: #333948;
        overflow: visible;
    }
    .property-list {
        padding: 4px 0;
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        .property-item {
            flex: 0 0 48%;
            margin: 8px 0 16px 0;
            font-size: 12px;
            width: 100%;
            .property-name {
                display: block;
                margin: 6px 0 9px;
                color: $cmdbTextColor;
                line-height: 16px;
                font-size: 0;
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                max-width: calc(100% - 20px);
                padding: 0 10px 0 0;
                vertical-align: middle;
                font-size: 14px;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                word-break: break-all;
                @include ellipsis;
                &.required:after {
                    position: absolute;
                    left: 100%;
                    top: 0;
                    margin: 0 0 0 -10px;
                    content: "*";
                    color: #ff5656;
                }
            }
            .property-name-tooltips {
                display: inline-block;
                vertical-align: middle;
                width: 16px;
                height: 16px;
                font-size: 16px;
                color: #c3cdd7;
            }
            .property-value {
                font-size: 0;
                position: relative;
                width: 303px;
                min-height: 32px;
            }

            &.full-width {
              flex: 0 0 100%;
              width: 100%;
              .property-value {
                width: 100%;
              }
            }
        }
    }
    .form-error {
        position: absolute;
        top: 100%;
        left: 0;
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
        max-width: 100%;
        @include ellipsis;
    }
</style>
