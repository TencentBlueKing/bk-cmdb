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
  <cmdb-sticky-layout class="field-view-layout">
    <div class="field-view-list" ref="fieldList">
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('唯一标识')}}</span>
        </div>
        <span class="property-value">{{field.bk_property_id}}</span>
      </div>
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('字段名称')}}</span>
        </div>
        <span class="property-value">{{field.bk_property_name}}</span>
      </div>
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('字段类型')}}</span>
        </div>
        <span class="property-value">{{fieldTypeMap[field.bk_property_type]}}</span>
      </div>
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('是否可编辑')}}</span>
        </div>
        <span class="property-value">{{field.editable ? $t('可编辑') : $t('不可编辑')}}</span>
      </div>
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('是否必填')}}</span>
        </div>
        <span class="property-value">{{field.isrequired ? $t('必填') : $t('非必填')}}</span>
      </div>
      <div class="property-item" v-if="hasMultipleType">
        <div class="property-name">
          <span>{{$t('是否可多选')}}</span>
        </div>
        <span class="property-value">{{field.ismultiple ? $t('可多选') : $t('不可多选')}}</span>
      </div>
      <div class="property-item" v-if="stringLikes.includes(field.bk_property_type)">
        <div class="property-name">
          <span>{{$t('正则校验')}}</span>
        </div>
        <span class="property-value">{{option || '--'}}</span>
      </div>
      <template v-else-if="numberLikes.includes(field.bk_property_type)">
        <div class="property-item">
          <div class="property-name">
            <span>{{$t('最小值')}}</span>
          </div>
          <span class="property-value">{{option.min || (option.min === 0 ? 0 : '--')}}</span>
        </div>
        <div class="property-item">
          <div class="property-name">
            <span>{{$t('最大值')}}</span>
          </div>
          <span class="property-value">{{option.max || (option.max === 0 ? 0 : '--')}}</span>
        </div>
      </template>
      <div class="property-item" v-if="numberLikes.includes(field.bk_property_type)">
        <div class="property-name">
          <span>{{$t('单位')}}</span>
        </div>
        <span class="property-value">{{field.unit || '--'}}</span>
      </div>
      <div class="property-item">
        <div class="property-name">
          <span>{{$t('用户提示')}}</span>
        </div>
        <span class="property-value">{{field.placeholder || '--'}}</span>
      </div>
      <div class="property-item full-width" v-if="isTableType">
        <div class="property-name">
          <span>{{$t('表格列设置')}}</span>
        </div>
        <div class="property-value mt10">
          <bk-table
            :data="option.header"
            :outer-border="true"
            :header-border="false">
            <bk-table-column
              :label="$t('字段ID')"
              prop="bk_property_id"
              :show-overflow-tooltip="true">
              <template #default="{ row }">
                <div class="property-name-text" :class="{ required: row.isrequired }">
                  {{row.bk_property_id}}
                </div>
              </template>
            </bk-table-column>
            <bk-table-column
              :label="$t('字段名称')"
              prop="bk_property_name"
              :show-overflow-tooltip="true" />
            <bk-table-column
              :label="$t('字段类型')"
              prop="bk_property_type"
              :show-overflow-tooltip="true">
              <template #default="{ row }">
                {{ fieldTypeMap[row.bk_property_type] }}
              </template>
            </bk-table-column>
          </bk-table>
        </div>
      </div>
      <div class="property-item full-width" v-if="isTableType">
        <div class="property-name">
          <span>{{$t('默认值')}}</span>
        </div>
        <div class="property-value mt10">
          <table-default-settings
            v-if="option.header.length > 0"
            :readonly="true"
            :preview="true"
            :headers="option.header"
            :defaults="option.default" />
          <span v-else>--</span>
        </div>
      </div>
      <div class="property-item enum-list" v-if="listLikes.includes(field.bk_property_type)">
        <div class="property-name">
          <span>{{ ['list'].includes(this.type) ? $t('列表值') : $t('枚举值') }}</span>
        </div>
        <div class="property-value">
          <template v-if="getEnumValue.length">
            <p v-for="(val, index) in getEnumValue" :key="index">{{val}}</p>
          </template>
          <span v-else>--</span>
        </div>
      </div>
      <div class="property-item enum-list" v-else-if="!isTableType">
        <div class="property-name">
          <span>{{$t('默认值')}}</span>
        </div>
        <cmdb-property-value
          class="property-value"
          v-if="['organization','enumquote'].includes(field.bk_property_type)"
          :value="defaultValue"
          :property="field">
        </cmdb-property-value>
        <!-- bool类型字段未设置与默认值false都将显示为false -->
        <span v-else
          class="property-value">{{defaultValue ?? '--'}}</span>
      </div>
    </div>
    <template slot="footer" slot-scope="{ sticky }" v-if="canEdit">
      <div class="btns" :class="{ 'is-sticky': sticky }">
        <bk-button class="mr10" theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
        <bk-button class="delete-btn" v-if="!field.ispre && !isTemplateField" @click="handleDelete">
          {{$t('删除')}}
        </bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<script>
  import { PROPERTY_TYPES, PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'
  import TableDefaultSettings from './table-default-settings.vue'

  export default {
    components: {
      TableDefaultSettings
    },
    props: {
      field: {
        type: Object,
        default: () => {}
      },
      canEdit: Boolean
    },
    data() {
      return {
        fieldTypeMap: PROPERTY_TYPE_NAMES,
        numberLikes: [PROPERTY_TYPES.INT, PROPERTY_TYPES.FLOAT],
        stringLikes: [PROPERTY_TYPES.SINGLECHAR, PROPERTY_TYPES.LONGCHAR],
        listLikes: [PROPERTY_TYPES.ENUM, PROPERTY_TYPES.ENUMMULTI, PROPERTY_TYPES.LIST],
        enumLikes: [PROPERTY_TYPES.ENUM, PROPERTY_TYPES.ENUMMULTI],
        scrollbar: false
      }
    },
    computed: {
      type() {
        return this.field.bk_property_type
      },
      option() {
        if (this.numberLikes.includes(this.type)) {
          return this.field.option || { min: null, max: null }
        }
        if (this.listLikes.includes(this.type)) {
          return this.field.option || []
        }
        if (this.isTableType) {
          return this.field.option || {
            header: [],
            default: []
          }
        }
        return this.field.option
      },
      hasMultipleType() {
        const types = [PROPERTY_TYPES.ORGANIZATION, PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.ENUMMULTI]
        return types.includes(this.type)
      },
      isTableType() {
        return this.type === PROPERTY_TYPES.INNER_TABLE
      },
      isTemplateField() {
        return this.field.bk_template_id > 0
      },
      defaultValue() {
        if ([PROPERTY_TYPES.ENUMQUOTE].includes(this.type)) {
          return  this.field.option.map(item => item.bk_inst_id)
        }
        return  this.$tools.getPropertyDefaultValue(this.field)
      },
      getEnumValue() {
        const value = this.field.option
        const type = this.field.bk_property_type
        if (Array.isArray(value)) {
          if (this.enumLikes.includes(type)) {
            const arr = value.map((item) => {
              if (item.is_default) {
                return `${item.name}(${item.id}, ${this.$t('默认值')})`
              }
              return `${item.name}(${item.id})`
            })
            return arr
          }
          if (this.listLikes.includes(type)) {
            const arr = value.map((item) => {
              if (item === this.field.default) {
                return `${item}(${this.$t('默认值')})`
              }
              return item
            })
            return arr
          }
        }
        return value || []
      }
    },
    methods: {
      handleEdit() {
        this.$emit('on-edit')
      },
      handleDelete() {
        this.$emit('on-delete')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .property-name-text {
      display: inline-block;
      position: relative;
      padding-right: 10px;
    }
    .field-view-layout {
        height: 100%;
        @include scrollbar-y;
    }
    .field-view-list {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
        padding: 0 24px 20px;
        font-size: 14px;
        .property-item {
            min-width: 48%;
            padding-top: 20px;
            &.enum-list,
            &.full-width {
                flex: 100%;
                .property-name,
                .property-value {
                    display: inline-block;
                    vertical-align: top;
                }
            }
            &.full-width {
              .property-value {
                display: block;
              }
              .property-name {
                width: auto;
              }
            }
        }
        .property-name {
            display: inline-block;
            color: #63656e;
            width: 90px;
            span {
              display: inline-block;
              position: relative;
              padding-right: 14px;
              &::after {
                position: absolute;
                right: 0;
                content: "：";
              }
            }
        }
        .property-value {
            color: #313237;
        }
    }
    .btns {
        background: #ffffff;
        padding: 10px 24px;
        font-size: 0;
        .delete-btn:hover {
            color: #ffffff;
            background-color: #ff5656;
            border-color: #ff5656;
        }
        &.is-sticky {
            border-top: 1px solid #dcdee5;
        }
    }
</style>
