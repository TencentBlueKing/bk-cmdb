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
  <div class="id-generate" v-bkloading="{ isLoading: globalConfig.loading }">
    <cmdb-tips class="tips">{{$t('ID生成器提示语')}}</cmdb-tips>
    <div>
      <bk-form :model="form" ref="formRef">
        <cmdb-collapse :label="$t('同步设置')" arrow-type="filled" class="form-model-config">
          <bk-form-item
            class="form-sync-value"
            :desc="{
              content: $t('同步设置描述')
            }"
            :label="$t('允许数据同步')" property="enabled" required>
            <div>
              <div class="basic-value" v-if="!isEdit">
                {{t('是否允许同步', {
                  value: form.enabled ? t('允许同步') : t('不允许同步')
                })}}
              </div>
              <bk-radio-group v-model="form.enabled" v-else>
                <bk-radio :value="false">{{t('不允许同步')}}</bk-radio>
                <bk-radio :value="true">{{t('允许同步')}}</bk-radio>
              </bk-radio-group>
            </div>
          </bk-form-item>
        </cmdb-collapse>
        <cmdb-collapse :label="$t('ID步长配置')" arrow-type="filled" class="form-model-config">
          <bk-form-item
            :rules="[{ trigger: 'blur', message: $t('ID自增步长必填'), required: true }]"
            :icon-offset="-20"
            :desc="{
              content: $t('ID自增步长描述')
            }"
            :label="$t('ID自增步长')" property="step" required>
            <div>
              <div class="basic-value" v-if="!isEdit">
                {{ form.step }}
              </div>
              <bk-input
                v-else
                :class="[{
                  'has-change': originForm.step !== form.step
                }]"
                type="number"
                :min="1"
                v-model.number.trim="form.step">
              </bk-input>
            </div>
          </bk-form-item>
        </cmdb-collapse>
        <cmdb-collapse :label="$t('起始ID配置')" arrow-type="filled" class="form-model-config">
          <bk-form-item
            v-for="property in modelFormKey"
            :key="property"
            :label="property"
            :property="`init_id.${property}`"
            :icon-offset="-20"
            :rules="[{ trigger: 'blur', message: $t('ID必填'), required: true }]"
            required>
            <div>
              <div class="basic-value" v-if="!isEdit">
                {{ form.current_id[property] }}
              </div>
              <div v-else>
                <bk-input
                  :class="[{
                    'has-change': form.current_id[property] !== form.init_id[property]
                  }]"
                  type="number"
                  :name="property"
                  :min="form.current_id[property]"
                  v-model.number.trim="form.init_id[property]">
                </bk-input>
                <div class="form-model-tip">{{$t('当前设置值', { value: form.current_id[property] })}}</div>
              </div>
            </div>
          </bk-form-item>
        </cmdb-collapse>
      </bk-form>
    </div>
    <div class="footer">
      <bk-button theme="primary" v-show="!isEdit" @click="isEdit = true">
        {{$t('编辑')}}
      </bk-button>
      <bk-button theme="primary" v-show="isEdit" @click="handleSubmit" :disabled="!hasChange">
        {{$t('提交')}}
      </bk-button>
      <bk-button theme="default" v-show="isEdit" @click="handleCancel">
        {{$t('取消')}}
      </bk-button>
    </div>
  </div>
</template>

<script setup>
  import { computed, reactive, ref, onMounted } from 'vue'
  import { bkInfoBox } from 'bk-magic-vue'
  import { t } from '@/i18n'
  import store from '@/store'
  import cloneDeep from 'lodash/cloneDeep'
  import isEqual from 'lodash/isEqual'
  import EventBus from '@/utils/bus'

  const defaultIdGenerateForm = {
    enabled: false,
    step: 0,
    init_id: {},
    current_id: {}
  }
  const form = reactive(cloneDeep(defaultIdGenerateForm))
  const originForm = reactive(cloneDeep(defaultIdGenerateForm))
  const isEdit = ref(false)
  const formRef = ref(null)

  const globalConfig = computed(() => store.state.globalConfig)
  const modelFormKey = computed(() => Object.keys(form.init_id))
  const hasChange = computed(() => !isEqual(originForm, form))

  const initForm = () => {
    const { idGenerator } = globalConfig.value.config
    Object.assign(form, cloneDeep(defaultIdGenerateForm), cloneDeep(idGenerator))
    Object.assign(originForm, cloneDeep(defaultIdGenerateForm), cloneDeep(idGenerator))
    formRef.value.clearError()
  }
  const handleSubmit = () => {
    const { enabled, step, init_id: initId, current_id: currentId } = form
    const submitForm = {
      enabled,
      step,
      init_id: {}
    }
    // 本次发生change的init_id
    const changeInitId = {}
    let hasChange = false
    Object.keys(currentId).forEach((key) => {
      if (currentId[key] !== initId[key]) {
        changeInitId[key] = initId[key]
        hasChange = true
      }
    })
    if (hasChange) submitForm.init_id = changeInitId

    formRef.value.validate().then(() => {
      bkInfoBox({
        title: `${t('确认提交')}?`,
        subTitle: t('确认提交ID生成器配置描述'),
        okText: t('确认提交'),
        cancelText: t('取消'),
        confirmFn: () => {
          store.dispatch('globalConfig/updateConfig', {
            idGenerator: submitForm
          })
            .then(() => {
              initForm()
            })
        }
      })
    })
  }
  const handleCancel = () => {
    isEdit.value = false
    initForm()
  }

  onMounted(() => {
    initForm()
    EventBus.$on('globalConfig/fetched', initForm)
  })

</script>

<style lang="scss" scoped>
.id-generate {
  width: 100%;
  margin-top: -16px;
}
.tips {
  width: 100%;
  white-space: pre-line;
  padding: 8px 10px;
  margin-bottom: 20px;
}

.form-model-config {
    padding: 15px;
    border: 1px solid #dcdee5;
    margin-bottom: 20px;

    :deep(.collapse-content) {
      @include space-between;
      flex-wrap: wrap;
      justify-content: flex-start;

      .bk-form-item {
        width: 33%;
        margin-top: 22px !important;
      }
      .form-sync-value {
        width: 100%;
      }
      .bk-form-radio {
        margin-right: 15px;
      }
    }

    .form-model-tip {
      color: #c4c6cc;
      font-size: 12px;
      font-weight: 400;
      line-height: 100%;
      margin-top: 4px;
    }
  }
  .has-change {
    :deep(.bk-form-input) {
      background-color: #FFF3E1;
    }
  }
  .basic-value {
    font-size: 14px;
    line-height: 32px;
  }

</style>
