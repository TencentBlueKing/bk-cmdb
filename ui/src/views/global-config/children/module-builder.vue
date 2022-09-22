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
  <div class="module-builder" :style="formItemStyle" :class="[{ 'indent-line': indentLine }, innerState]">
    <div class="module-builder-content">
      <bk-input
        type="text"
        v-model="innerModuleId"
        class="module-id-input"
        :placeholder="moduleIdPlaceholder"
        @change="handleModuleIdChange"
        :disabled="innerState === 'readonly' || moduleIdDisabled"></bk-input>
      <bk-input
        type="text"
        :class="{
          'module-name-input': indentLine,
          'cluster-name-input': !indentLine
        }"
        :disabled="innerState === 'readonly'"
        v-model="innerModuleName"
        @change="handleModuleNameChange"
        :placeholder="moduleNamePlaceholder"
        :maxlength="50"
      >
      </bk-input>
      <div class="operation">
        <span class="icon-button icon-cc-edit-shape" @click="edit">
        </span>
        <span class="icon-button remove-button icon-cc-tips-close" v-show="removeable" @click="remove">
        </span>
        <span class="text-button" @click.stop="confirm">{{$t('确定')}}
        </span>
        <span class="text-button" @click="cancel">{{$t('取消')}}</span>
      </div>
    </div>
  </div>
</template>

<script>
  import cloneDeep from 'lodash/cloneDeep'

  export default {
    name: 'module-builder',
    props: {
      moduleId: {
        type: String,
        default: ''
      },
      moduleIdDisabled: {
        type: Boolean,
        default: false
      },
      moduleName: {
        type: String,
        default: ''
      },
      /**
       * 组件状态 readonly(只读)|editting(可编辑)
       */
      state: {
        type: String,
        default: 'readonly'
      },
      /**
       * 是否需要缩进线
       */
      indentLine: {
        type: Boolean,
        default: false
      },
      /**
       * 节点缩进，默认为 0
       */
      nodeIndent: {
        type: [String, Number],
        default: 0
      },
      removeable: {
        type: Boolean,
        default: true
      },
      moduleIdPlaceholder: {
        type: String,
        default: ''
      },
      moduleNamePlaceholder: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        innerModuleName: '',
        innerModuleId: '',
        innerState: 'readonly'
      }
    },
    computed: {
      formItemStyle() {
        return {
          marginLeft: `${parseInt(this.nodeIndent, 10)}px`
        }
      }
    },
    watch: {
      moduleId: {
        immediate: true,
        handler(value) {
          this.innerModuleId = value
        }
      },
      moduleName: {
        immediate: true,
        handler(value) {
          this.innerModuleName = value
        }
      },
      state: {
        immediate: true,
        handler(newState) {
          if (this.innerState !== newState) this.innerState = newState
        }
      },
    },
    methods: {
      handleModuleNameChange() {
        this.$emit('update:moduleName', this.innerModuleName)
      },
      handleModuleIdChange() {
        this.$emit('update:moduleId', this.innerModuleId)
      },
      edit() {
        this.innerState = 'editting'
        this.oldModuleName = cloneDeep(this.innerModuleName)
        this.oldModuleId = cloneDeep(this.innerModuleId)
      },
      remove() {
        this.$emit('remove')
      },
      confirm() {
        const done = () => {
          this.innerState = 'readonly'
          this.$emit('update:state', this.innerState)
        }
        if ('before-confirm' in this.$listeners) {
          this.$emit('before-confirm', done)
        } else {
          done()
        }

        this.$emit('confirm', {
          moduleId: this.innerModuleId,
          moduleName: this.innerModuleName
        })
      },
      cancel() {
        if (this.innerState !== 'editting') return false

        const done = () => {
          this.innerState = 'readonly'
          this.innerModuleName = cloneDeep(this.oldModuleName)
          this.innerModuleId = cloneDeep(this.oldModuleId)
          this.$emit('update:state', this.innerState)
        }

        if ('before-cancel' in this.$listeners) {
          this.$emit('before-cancel', done)
        } else {
          done()
        }

        this.$emit('cancel', {
          moduleId: this.innerModuleId,
          moduleName: this.innerModuleName
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
.module-builder {
    position: relative;

    &-content{
      position: relative;
      display: flex;
      z-index: 1;
    }

    .module-id-input{
      flex: 0 0 140px;
      margin-right: 4px;
    }

    .module-name-input{
      flex: 0 0 338px;
    }

    .cluster-name-input{
      flex: 0 0 (338px + 40px);
    }

    &.readonly .text-button{
      display: none;
    }

    &.editting .icon-button{
      display: none;
    }

    .remove-button{
      font-size: 14px;
    }

    &.indent-line::after{
      content: "";
      position: absolute;
      bottom: 50%;
      right: 100%;
      display: block;
      width: 26px;
      height: 32px + 20px;
      border-bottom: 1px solid #dcdee5;
      border-left: 1px solid #dcdee5;
      z-index: 0;
    }

    .operation {
      flex-grow: 0;
      flex-shrink: 0;
      margin-top: auto;
      margin-bottom: auto;
      margin-left: 10px;
    }

    .text-button {
        display: inline-block;
        line-height: normal;
        vertical-align: middle;
        font-size: 12px;
        color: $primaryColor;
        cursor: pointer;

        & + .text-button {
            $marginLeft: 6px;
            position: relative;
            margin-left: $marginLeft;

            &::after {
              content: '';
              position: absolute;
              display: block;
              top: 50%;
              transform: translateY(-50%);
              left: -$marginLeft;
              width: 1px;
              height: 14px;
              background-color: #dcdee5;
          }
        }
    }

    .icon-button {
      display: inline-block;
      vertical-align: middle;
      cursor: pointer;
      color: #979BA5;
      margin-top: auto;
      margin-bottom: auto;

      &:hover{
        color: $primaryColor;
      }

      &.remove-button {
        margin-left: 14px;
        font-size: 12px;
      }
    }
}
</style>
