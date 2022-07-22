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
  <div>
    <div class="form-label">
      <span class="label-text">{{$t('最小值')}}</span>
      <div class="cmdb-form-item" :class="{ 'is-error': errors.has('min') }">
        <bk-input type="text" class="cmdb-form-input"
          v-model="localValue.min"
          @input="handleInput"
          v-validate="`float`"
          :disabled="isReadOnly"
          :name="'min'">
        </bk-input>
        <p class="form-error">{{errors.first('min')}}</p>
      </div>
    </div>
    <div class="form-label">
      <span class="label-text">{{$t('最大值')}}</span>
      <div class="cmdb-form-item" :class="{ 'is-error': errors.has('max') }">
        <bk-input type="text" class="cmdb-form-input"
          v-model="localValue.max"
          name="max"
          @input="handleInput"
          :disabled="isReadOnly"
          v-validate="`float|isBigger:${localValue.min}`">
        </bk-input>
        <p class="form-error">{{errors.first('max')}}</p>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: [Object, String],
        default: {
          min: '',
          max: ''
        }
      },
      isReadOnly: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        localValue: {
          min: '',
          max: ''
        }
      }
    },
    watch: {
      value: {
        handler() {
          this.initValue()
        },
        deep: true
      }
    },
    created() {
      this.initValue()
    },
    methods: {
      initValue() {
        if (this.value === '' || this.value === null) {
          this.localValue = {
            min: '',
            max: ''
          }
        } else {
          this.localValue = this.value
        }
      },
      async handleInput() {
        const res = await this.$validator.validateAll()
        if (res) {
          this.$emit('input', this.localValue)
        }
      },
      validate() {
        return this.$validator.validateAll()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-label {
        &:last-child {
           margin: 0;
        }
    }
</style>
