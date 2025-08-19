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
  <div class="form-enum-layout">
    <div class="toolbar">
      <p class="title">{{$t('枚举值')}}</p>
      <i
        v-bk-tooltips.top-start="$t('通过枚举项的值按照0-9，a-z排序')"
        :class="['sort-icon', `icon-cc-sort-${order > 0 ? 'up' : 'down'}`]"
        @click="handleSort">
      </i>
    </div>
    <vue-draggable
      class="form-enum-wrapper"
      tag="ul"
      v-model="enumList"
      :options="dragOptions"
      @end="handleDragEnd">
      <li class="form-item" v-for="(item, index) in enumList" :key="index">
        <div class="enum-id">
          <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`id${index}`) }">
            <bk-input type="text"
              class="cmdb-form-input"
              :placeholder="$t('请输入ID')"
              v-model.trim="item.id"
              v-validate="`required|enumId|length:128|repeat:${getOtherId(index)}`"
              @input="handleInput"
              :disabled="isReadOnly"
              :name="`id${index}`"
              :ref="`id${index}`">
            </bk-input>
            <p class="form-error" :title="errors.first(`id${index}`)">{{errors.first(`id${index}`)}}</p>
          </div>
        </div>
        <div class="enum-name">
          <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`name${index}`) }">
            <bk-input type="text"
              class="cmdb-form-input"
              :placeholder="$t('请输入值')"
              v-model.trim="item.name"
              v-validate="`required|enumName|length:128|repeat:${getOtherName(index)}`"
              @input="handleInput"
              :disabled="isReadOnly"
              :name="`name${index}`">
            </bk-input>
            <p class="form-error" :title="errors.first(`name${index}`)">{{errors.first(`name${index}`)}}</p>
          </div>
        </div>
        <bk-button text class="enum-btn" @click="deleteEnum(index)" :disabled="enumList.length === 1 || isReadOnly">
          <i class="bk-icon icon-minus-circle-shape"></i>
        </bk-button>
        <bk-button text class="enum-btn" @click="addEnum(index)"
          :disabled="isReadOnly" v-if="index === enumList.length - 1">
          <i class="bk-icon icon-plus-circle-shape"></i>
        </bk-button>
      </li>
    </vue-draggable>
    <div class="default-setting">
      <p class="title mb10">{{$t('默认值设置')}}</p>
      <div class="cmdb-form-item" :class="{ 'is-error': errors.has('defaultValueSelect') }">
        <div class="form-item-row">
          <bk-select style="width: 100%;"
            :key="defaultCompKey"
            :scroll-height="150"
            :searchable="true"
            :clearable="false"
            :disabled="isReadOnly"
            :multiple="isDefaultCompMultiple"
            name="defaultValueSelect"
            data-vv-validate-on="change"
            :popover-options="{
              appendTo: 'parent'
            }"
            v-validate="`maxSelectLength:${ multiple ? -1 : 1 }`"
            v-model="defaultValue"
            @change="handleSettingDefault">
            <bk-option v-for="option in settingList"
              :key="option.id"
              :id="option.id"
              :name="option.name">
            </bk-option>
          </bk-select>
          <bk-checkbox
            v-if="isDefaultCompMultiple"
            class="checkbox"
            v-model="isMultiple"
            :disabled="isReadOnly">
            <span>{{$t('可多选')}}</span>
          </bk-checkbox>
        </div>
        <p class="form-error">{{errors.first('defaultValueSelect')}}</p>
      </div>
    </div>
  </div>
</template>

<script>
  import vueDraggable from 'vuedraggable'
  import isEqual from 'lodash/isEqual'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'

  export default {
    components: {
      vueDraggable
    },
    props: {
      value: {
        type: [Array, String],
        default: ''
      },
      isReadOnly: {
        type: Boolean,
        default: false
      },
      multiple: {
        type: Boolean,
        default: false
      },
      type: String
    },
    data() {
      return {
        enumList: [this.generateEnum()],
        settingList: [],
        defaultValue: this.multiple ? [] : '',
        dragOptions: {
          animation: 300,
          disabled: false,
          filter: '.enum-btn, .enum-id, .enum-name',
          preventOnFilter: false,
          ghostClass: 'ghost'
        },
        order: 1,
        defaultCompKey: null
      }
    },
    computed: {
      isDefaultCompMultiple() {
        // 通过类型指定默认值组件是否可多选，用于与可多选配置区分开
        return this.type === PROPERTY_TYPES.ENUMMULTI
      },
      isMultiple: {
        get() {
          return this.multiple
        },
        set(val) {
          this.$emit('update:multiple', val)
        }
      }
    },
    watch: {
      value() {
        this.initValue()
      },
      enumList: {
        deep: true,
        handler(value) {
          // 解决在id或name全部清空的情况下，重新填写的name在下拉框中显示的为上一次name值
          this.defaultCompKey = Date.now()

          // 重复的选项不允许加入的选择列表
          const enumList = []
          if (value.length) {
            enumList.push(value[0])
            value.forEach((data) => {
              if (!enumList.some(item => item.id === data.id || item.name === data.name)) {
                enumList.push(data)
              }
            })
          }
          this.settingList = enumList.filter(item => item.id && item.name)

          // 无默认值选择第0项，有默认值则需要验证值是否存在（列表中可能将其删除）
          if (!this.defaultValue?.length) {
            if (this.isDefaultCompMultiple) {
              this.defaultValue = this.settingList.length ? [this.settingList[0].id] : []
            } else {
              this.defaultValue = this.settingList.length ? this.settingList[0].id : ''
            }
          } else {
            if (this.isDefaultCompMultiple) {
              this.defaultValue = this.settingList.length
                ? this.settingList.filter(item => this.defaultValue.includes(item.id)).map(item => item.id)
                : []
            } else {
              this.defaultValue = this.settingList.length
                ? this.settingList.find(item => this.defaultValue === item.id)?.id ?? ''
                : ''
            }
          }
        }
      },
      defaultValue(val, old) {
        // 检测选中值变化，需要修正is_default，这里值的类型都是string
        if (val && !isEqual(val, old)) {
          this.enumList.forEach((item) => {
            if (this.isDefaultCompMultiple && Array.isArray(val)) {
              item.is_default = val.includes(item.id)
            } else {
              item.is_default = val === item.id
            }
          })
          this.$emit('input', this.enumList)
        }
      },
      multiple() {
        // 多选变化时校验默认值设置
        this.$nextTick(async () => this.$validator.validate('defaultValueSelect'))
      }
    },
    created() {
      this.initValue()
    },
    methods: {
      getOtherId(index) {
        const idList = []
        this.enumList.forEach((item, enumIndex) => {
          if (index !== enumIndex) {
            idList.push(item.id)
          }
        })
        return idList.join(',')
      },
      getOtherName(index) {
        const nameList = []
        this.enumList.forEach((item, enumIndex) => {
          if (index !== enumIndex) {
            nameList.push(item.name)
          }
        })
        return nameList.join(',')
      },
      initValue() {
        // 枚举多选默认值是空数组
        if (this.value === '' || (Array.isArray(this.value) && !this.value.length)) {
          this.enumList = [this.generateEnum()]
        } else {
          this.enumList = this.value.map(data => (this.generateEnum(data)))
          const defaultValues = this.enumList.filter(item => item.is_default).map(item => item.id)
          this.defaultValue = this.isDefaultCompMultiple ? defaultValues : defaultValues[0]
        }
      },
      handleInput() {
        this.$emit('input', this.enumList)
      },
      addEnum(index) {
        this.enumList.push(this.generateEnum({ is_default: false }))
        this.$nextTick(() => {
          this.$refs[`id${index + 1}`] && this.$refs[`id${index + 1}`][0].focus()
        })
      },
      deleteEnum(index) {
        this.enumList.splice(index, 1)
        this.handleInput()
      },
      generateEnum(settings = {}) {
        const defaults = {
          id: '',
          is_default: true,
          name: '',
          type: 'text'
        }
        return { ...defaults, ...settings }
      },
      validate() {
        return this.$validator.validateAll()
      },
      handleSettingDefault(id) {
        if (this.multiple) {
          this.enumList.forEach((item) => {
            item.is_default = id.includes(item.id)
          })
          this.$emit('input', this.enumList)
        } else {
          const itemIndex = this.enumList.findIndex(item => item.id === id)
          if (itemIndex > -1) {
            this.enumList.forEach((item) => {
              item.is_default = item.id === id
            })

            this.$emit('input', this.enumList)
          }
        }
      },
      handleDragEnd() {
        this.$emit('input', this.enumList)
      },
      handleSort() {
        this.order = this.order * -1
        this.enumList.sort((A, B) => A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' }) * this.order)

        this.$emit('input', this.enumList)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .title {
        font-size: 14px;
    }
    .form-enum-wrapper {
        .form-item {
            display: flex;
            align-items: center;
            position: relative;
            font-size: 0;
            margin-bottom: 16px;
            padding: 2px 2px 2px 28px;
            cursor: move;

            &::before {
                content: '';
                position: absolute;
                top: 12px;
                left: 8px;
                width: 3px;
                height: 3px;
                border-radius: 50%;
                background-color: #979ba5;
                box-shadow: 0 5px 0 0 #979ba5,
                    0 10px 0 0 #979ba5,
                    5px 0 0 0 #979ba5,
                    5px 5px 0 0 #979ba5,
                    5px 10px 0 0 #979ba5;
            }

            .enum-id {
                width: 90px;
                margin-right: 10px;
                input {
                    width: 100%;
                }
            }
            .enum-name {
                width: 180px;
                input {
                    width: 100%;
                }
            }
            .enum-btn {
                font-size: 0;
                color: #c4c6cc;
                margin: -2px 0 0 6px;
                .bk-icon {
                    width: 18px;
                    height: 18px;
                    line-height: 18px;
                    font-size: 18px;
                    text-align: center;
                }

                &.is-disabled {
                    color: #dcdee5;
                }
                &:not(.is-disabled):hover {
                    color: #979ba5;
                }
            }
        }
    }

    .toolbar {
        display: flex;
        margin-bottom: 10px;
        align-items: center;
        line-height: 20px;

        .sort-icon {
            width: 20px;
            height: 20px;
            margin-left: 10px;
            border: 1px solid #c4c6cc;
            background: #fff;
            border-radius: 2px;
            font-size: 16px;
            line-height: 18px;
            text-align: center;
            color: #c4c6cc;
            cursor: pointer;

            &:hover {
                color: #979ba5;
            }
        }
    }

    .ghost {
        border: 1px dashed $cmdbBorderFocusColor;
    }

    .form-item-row {
      display: flex;
      align-items: center;
      gap: 12px;
      .checkbox {
        flex: none;
      }
    }
</style>
