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
  <cmdb-sticky-layout class="form-layout">
    <div class="form-groups" v-if="hasAvaliableGroups" ref="formGroups">
      <template v-for="(group, groupIndex) in $sortedGroups">
        <div class="property-group"
          :key="groupIndex"
          v-if="groupedProperties[groupIndex].length">
          <cmdb-collapse
            :label="group['bk_group_name']"
            :collapse.sync="groupState[group['bk_group_id']]">
            <ul class="property-list">
              <template v-for="(property, propertyIndex) in groupedProperties[groupIndex]">
                <li :class="['property-item', property.bk_property_type]"
                  v-if="!uneditableProperties.includes(property.bk_property_id)"
                  :key="propertyIndex">
                  <cmdb-auth tag="div" class="property-name" :title="property['bk_property_name']" v-bind="authProps">
                    <bk-checkbox class="property-name-checkbox" slot-scope="{ disabled }"
                      :id="`property-name-${property['bk_property_id']}`"
                      :disabled="disabled || $tableTypePropertyIds.includes(property['bk_property_id'])"
                      v-model="editable[property['bk_property_id']]">
                      <span class="property-name-text"
                        v-bk-tooltips.top-start="{
                          content: $t('暂不支持'),
                          disabled: !$tableTypePropertyIds.includes(property['bk_property_id'])
                        }"
                        :for="`property-name-${property['bk_property_id']}`"
                        :class="{ required: property['isrequired'] && editable[property.bk_property_id] }">
                        {{property['bk_property_name']}}
                      </span>
                      <i class="property-name-tooltips icon icon-cc-tips"
                        v-if="property.placeholder && isExclmationProperty(property.bk_property_type)"
                        v-bk-tooltips="{
                          trigger: 'mouseenter',
                          content: property.placeholder
                        }">
                      </i>
                    </bk-checkbox>
                  </cmdb-auth>
                  <div class="property-value">
                    <component class="form-component"
                      v-if="property.bk_property_type !== PROPERTY_TYPES.INNER_TABLE"
                      :is="`cmdb-form-${property['bk_property_type']}`"
                      :class="{ error: errors.has(property['bk_property_id']) }"
                      :unit="property['unit']"
                      :row="2"
                      :disabled="!editable[property['bk_property_id']]"
                      :options="property.option || []"
                      :data-vv-name="property['bk_property_id']"
                      :auto-select="false"
                      :multiple="property.ismultiple"
                      :placeholder="$tools.getPropertyPlaceholder(property)"
                      v-bk-tooltips.top="{
                        disabled: !property.placeholder,
                        theme: 'light',
                        trigger: 'click',
                        content: property.placeholder
                      }"
                      v-bind="$tools.getValidateEvents(property)"
                      v-validate="getValidateRules(property)"
                      v-model.trim="values[property['bk_property_id']]">
                    </component>
                    <cmdb-form-innertable v-else
                      :mode="'update'"
                      :property="property"
                      :immediate="false"
                      :readonly="true"
                      v-model.trim="values[property['bk_property_id']]" />
                    <span class="form-error"
                      :title="errors.first(property['bk_property_id'])">
                      {{errors.first(property['bk_property_id'])}}
                    </span>
                  </div>
                </li>
              </template>
            </ul>
          </cmdb-collapse>
        </div>
      </template>
    </div>
    <div class="form-empty" v-else>
      {{$t('暂无可批量更新的属性')}}
    </div>
    <div slot="footer" slot-scope="{ sticky }"
      :class="['form-options', { sticky }]">
      <slot name="details-options">
        <cmdb-auth class="inline-block-middle" v-bind="authProps">
          <bk-button slot-scope="{ disabled }"
            class="button-save"
            theme="primary"
            :loading="loading || allValidating"
            :disabled="disabled || !hasChange"
            @click="handleSave">
            {{$t('保存')}}
          </bk-button>
        </cmdb-auth>
        <bk-button class="button-cancel" :disabled="loading" @click="handleCancel">{{$t('取消')}}</bk-button>
      </slot>
    </div>
  </cmdb-sticky-layout>
</template>

<script>
  import formMixins from '@/mixins/form'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { BUILTIN_UNEDITABLE_FIELDS } from '@/dictionary/model-constants'
  import useSideslider from '@/hooks/use-sideslider'
  import { isExclmationProperty } from '@/utils/util'

  export default {
    name: 'cmdb-form-multiple',
    mixins: [formMixins],
    props: {
      saveAuth: {
        type: [Object, Array],
        default: null
      },
      loading: {
        type: Boolean,
        default: false,
      }
    },
    data() {
      return {
        PROPERTY_TYPES,
        isMultiple: true,
        allValidating: false,
        values: {},
        refrenceValues: {},
        editable: {},
        groupState: {
          none: true
        }
      }
    },
    computed: {
      changedValues() {
        const changedValues = {}
        Object.keys(this.values).forEach((propertyId) => {
          const property = this.getProperty(propertyId)
          if (
            ['bool'].includes(property.bk_property_type)
            || this.values[propertyId] !== this.refrenceValues[propertyId]
          ) {
            changedValues[propertyId] = this.values[propertyId]
          }
        })
        return changedValues
      },
      hasChange() {
        let hasChange = false
        // eslint-disable-next-line no-restricted-syntax
        for (const propertyId in this.editable) {
          if (this.editable[propertyId]) {
            hasChange = true
            break
          }
        }
        return hasChange
      },
      groupedProperties() {
        return this.$groupedProperties.map(properties => properties.filter((property) => {
          const { editable } = property
          const isapi = property.bk_isapi
          const { isonly } = property
          const isAsst = ['singleasst', 'multiasst'].includes(property.bk_property_type)
          const isUneditable = this.uneditableProperties.includes(property.bk_property_id)
          const isBuiltinUneditable = BUILTIN_UNEDITABLE_FIELDS.includes(property.bk_property_id)
          // eslint-disable-next-line max-len
          return editable && !isapi && !isonly && !isAsst && !isUneditable && !isBuiltinUneditable
        }))
      },
      hasAvaliableGroups() {
        return this.groupedProperties.some(properties => !!properties.length)
      },
      authProps() {
        if (this.saveAuth) {
          return {
            auth: this.saveAuth
          }
        }
        return {
          auth: [],
          ignore: true
        }
      }
    },
    watch: {
      properties() {
        this.initValues()
        this.initEditableStatus()
      }
    },
    created() {
      this.initValues()
      this.initEditableStatus()
      const { beforeClose, setChanged } = useSideslider(this.values)
      this.beforeClose = beforeClose
      this.setChanged = setChanged
    },
    methods: {
      isExclmationProperty(type) {
        return isExclmationProperty(type)
      },
      initValues() {
        this.values = this.$tools.getInstFormValues(this.properties, {}, false)
        this.refrenceValues = this.$tools.clone(this.values)
      },
      initEditableStatus() {
        const editable = {}
        this.groupedProperties.forEach((properties) => {
          properties.forEach((property) => {
            editable[property.bk_property_id] = false
          })
        })
        this.editable = editable
      },
      htmlEncode(placeholder) {
        let temp = document.createElement('div')
        temp.innerHTML = placeholder
        const output = temp.innerText
        temp = null
        return output
      },
      getProperty(id) {
        return this.properties.find(property => property.bk_property_id === id)
      },
      getPlaceholder(property) {
        const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
        return this.$t(placeholderTxt, { name: property.bk_property_name })
      },
      getValidateRules(property) {
        if (!this.editable[property.bk_property_id]) {
          return {}
        }
        return this.$tools.getValidateRules(property)
      },
      getMultipleValues() {
        const multipleValues = {}
        Object.keys(this.values).forEach((propertyId) => {
          if (this.editable[propertyId]) {
            multipleValues[propertyId] = this.values[propertyId]
          }
        })
        return this.$tools.formatValues(multipleValues, this.properties)
      },
      handleSave() {
        this.allValidating = true
        this.$validator.validateAll().then((result) => {
          if (result) {
            this.$emit('on-submit', this.getMultipleValues())
          } else {
            this.uncollapseGroup()
          }
        })
          .finally(() => {
            this.allValidating = false
          })
      },
      uncollapseGroup() {
        this.errors.items.forEach((item) => {
          const property = this.properties.find(property => property.bk_property_id === item.field)
          const group = property.bk_property_group
          this.groupState[group] = false
        })
      },
      handleCancel() {
        this.$emit('on-cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-layout {
        height: 100%;
        @include scrollbar;
    }
    .form-groups {
        padding: 0 0 0 32px;
    }
    .property-group {
        padding: 7px 0 10px 0;
        &:first-child{
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
        .property-item {
            margin: 12px 0 0;
            padding: 0 54px 0 0;
            font-size: 12px;
            flex: 0 0 50%;
            max-width: 50%;
            .property-name {
                display: flex;
                margin: 6px 0 10px;
                color: $cmdbTextColor;
                font-size: 0;
                line-height: 18px;
            }
            .property-name-checkbox {
                margin: 0 6px 0 0;
                max-width: calc(100% - 30px);
                display: flex;

                /deep/ .bk-checkbox-text {
                    width: calc(100% - 30px);
                    flex: 1;
                }
            }
            .property-name-text {
                position: relative;
                display: inline-block;
                max-width: 100%;
                padding: 0 10px 0 0;
                vertical-align: top;
                font-size: 14px;
                @include ellipsis;
                &.required:after{
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
                margin: 0 3px 0 -5px;
            }
            .property-value {
                font-size: 0;
                position: relative;
                /deep/ .control-append-group {
                    .bk-input-text {
                        flex: 1;
                    }
                }
            }

            &.innertable {
                flex: 1 1 100%;
                width: 100%;
                max-width: unset;
            }
        }
    }
    .form-options {
        padding: 10px 32px;
        background-color: #fff;
        &.sticky {
            border-top: 1px solid $cmdbBorderColor;
        }
        .button-save {
            min-width: 76px;
            margin-right: 4px;
        }
        .button-cancel {
            min-width: 76px;
            background-color: #fff;
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
    .form-empty {
        height: 100%;
        text-align: center;
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            height: 100%;
            width: 0;
        }
    }
</style>
