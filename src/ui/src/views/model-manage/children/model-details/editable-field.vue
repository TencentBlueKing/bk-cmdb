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
  <div
    class="editable-field clearfix"
    @click.stop
    v-click-outside="{
      handler: handleClickOutSide,
      middleware: clickOutSideMiddleware,
      isActive: isEditing
    }"
    :class="{ 'is-error': isError, 'is-editing': isEditing, 'is-readonly': !editable }">
    <template v-if="!editable">{{ label || value }}</template>
    <div v-else class="editable-field-container">
      <div class="editable-field-content">
        <span class="editable-field-text" v-show="!isEditing">{{ label || value }}</span>
        <div class="editable-field-control" v-show="isEditing">
          <cmdb-singlechar
            v-if="type === 'singlechar'"
            v-autofocus
            @change="handleInputChange"
            @enter="confirmEdit"
            :disabled="isConfirming"
            v-bind="$attrs"
            v-model="innerValue">
            <slot></slot>
          </cmdb-singlechar>
          <!-- 切换编辑状态会初始化组件，目的是为了触发 show-on-init 的效果 -->
          <cmdb-enum
            v-if="type === 'enum' && isEditing"
            :show-on-init="true"
            @on-selected="confirmEdit"
            v-autofocus
            :disabled="isConfirming"
            v-bind="$attrs"
            v-model="innerValue">
          </cmdb-enum>
          <div class="tips-icon">
            <bk-icon
              v-show="isError"
              v-bk-tooltips="{
                content: errorText
              }"
              class="error-tips-icon"
              type="exclamation-circle-shape">
            </bk-icon>
          </div>
        </div>
      </div>
      <div v-show="!isEditing" class="editable-field-edit-button">
        <cmdb-auth
          tag="i"
          class="icon-cc-edit-shape"
          :auth="auth"
          @click="edit">
        </cmdb-auth>
      </div>
    </div>
    <slot name="append"></slot>
  </div>
</template>

<script>
  import { defineComponent, ref, watch } from '@vue/composition-api'
  import { Validator } from 'vee-validate'
  import { autofocus } from '@/directives/autofocus'
  import CmdbSinglechar from '@/components/ui/form/singlechar.vue'
  import CmdbEnum from '@/components/ui/form/enum.vue'

  export default defineComponent({
    name: 'EditableField',
    directives: {
      autofocus
    },
    components: {
      CmdbSinglechar,
      CmdbEnum
    },
    model: {
      prop: 'value',
      event: 'value-confirm'
    },
    props: {
      // 传入的值，支持 v-model
      value: {
        type: [String, Number, Boolean],
        default: ''
      },
      // 展示的内容，有时候需要展示的不是 value 而是 value 对应的 name，此时可以格式化后用 label 展示
      label: {
        type: String,
        default: ''
      },
      // 控件类型，默认为 singlechar
      type: {
        type: String,
        default: 'singlechar',
        validator: val => ['singlechar', 'enum'].includes(val)
      },
      // 权限数据，仅使用在编辑按钮上
      auth: {
        type: Object,
        default: () => ({})
      },
      // 是否支持编辑
      editable: {
        type: Boolean,
        default: true
      },
      // 是否在编辑状态中，支持 .sync 修饰符，便于从外部获得编辑状态
      editing: {
        type: Boolean,
        default: true
      },
      // vee-validate 验证规则，参考：https://vee-validate.logaretm.com/v2/guide/syntax.html#rules-parameters
      validate: {
        type: String,
        default: ''
      }
    },
    setup(props, { emit }) {
      const isEditing = ref(false)
      const isError = ref(false)
      const innerValue = ref('')
      const errorText = ref('')
      const valueRef = ref(props.value)
      const validator = new Validator()
      const isConfirming = ref(false)

      watch(
        valueRef, (val) => {
          innerValue.value = val
        },
        {
          immediate: true
        }
      )

      const edit = () => {
        isEditing.value = true
        emit('update:editing', true)
      }

      const confirmEdit = async () => {
        const { valid, errors } = await validator.verify(innerValue.value, props.validate)

        const stop = () => {
          isConfirming.value = false
        }

        const confirm = () => {
          isEditing.value = false
          stop()
          emit('value-confirm', innerValue.value)
          emit('update:editing', false)
        }

        if (valid) {
          if (innerValue.value !== props.value) {
            isConfirming.value = true
            emit('confirm', { value: innerValue.value, confirm, stop })
          } else {
            confirm()
          }
        } else {
          [errorText.value] = errors
          isError.value = true
        }
      }

      const handleInputChange = () => {
        errorText.value = ''
        isError.value = false
      }

      // 针对 bk-select 的特殊处理，避免点击下拉区域时触发失焦退出编辑
      const clickOutSideMiddleware = event => !event.path.some(node => node.className === 'bk-select-dropdown-content')

      const handleClickOutSide = () => {
        if (isEditing.value) {
          confirmEdit()
        }
      }

      return {
        isError,
        isEditing,
        isConfirming,
        innerValue,
        errorText,
        edit,
        confirmEdit,
        handleInputChange,
        handleClickOutSide,
        clickOutSideMiddleware
      }
    },
  })
</script>

<style lang="scss" scoped>
$editBtnSize: 12px;
$editBtnMarginLeft: 4px;
$restWidth: 10px;

.editable-field {
  display: inline-block;
  max-width: 100%;
  vertical-align: middle;
  font-weight: normal;

  &-container {
    display: flex;
    align-items: center;
  }

  &.is-error {
    /deep/ .bk-form-input {
      border-color: #EA3636;
    }
  }

  &-content {
    max-width: calc(100% - #{$editBtnSize + $editBtnMarginLeft + $restWidth});
    @include ellipsis;
  }

  &.is-editing &-content {
    width: 100%;
  }

  &-control {
    position: relative;
  }

  .bk-form-input,
  .bk-select {
    width: 100%;
  }


  .tips-icon {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 10px;
    font-size: 16px;
    display: flex;
    align-items: center;
  }

  .error-tips-icon {
    color:#EA3636;
  }

  &-edit-button {
    flex: 0 0 $editBtnSize;
    align-items: center;
    font-size: $editBtnSize;
    margin-left: $editBtnMarginLeft;
    cursor: pointer;

    &:hover {
      color: $primaryColor;
    }
  }
}
</style>
