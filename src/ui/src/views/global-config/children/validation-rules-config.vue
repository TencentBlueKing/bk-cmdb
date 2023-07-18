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
    v-bkloading="{ isLoading: globalConfig.loading }"
    class="validation-rules-config">
    <bk-form
      ref="ruleFormRef"
      :label-width="200"
      :model="ruleForm"
      :rules="ruleFormRules">
      <bk-form-item
        v-for="(rule, ruleKey) of ruleForm"
        :key="ruleKey"
        :icon-offset="105"
        :property="ruleKey"
        required
        :label="translateRuleKey(ruleKey)">
        <bk-input
          v-model="rule.value"
          style="width: 300px"
          :placeholder="$t('请填写验证规则')"></bk-input>
        <bk-input
          v-model="rule.message"
          style="width: 400px"
          :placeholder="$t('请填写校验提示')"></bk-input>
        <bk-button
          size="small"
          text
          class="recovery-button"
          @click="recoveryRule(rule, ruleKey)">
          {{ $t('恢复初始化') }}
        </bk-button>
      </bk-form-item>
      <bk-form-item class="form-action-item">
        <SaveButton :loading="globalConfig.updating" @save="save"></SaveButton>
        <bk-popconfirm
          trigger="click"
          :title="$t('确认重置所有校验规则？')"
          :content="
            $t(
              '该操作将会把所有的校验规则内容重置为最后一次保存的状态，请谨慎操作！'
            )
          "
          @confirm="resetAllRules">
          <bk-button class="action-button">{{ $t('重置') }}</bk-button>
        </bk-popconfirm>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
import { ref, reactive, computed, defineComponent, onMounted } from 'vue'
import { bkMessage } from 'bk-magic-vue'
import cloneDeep from 'lodash/cloneDeep'

import store from '@/store'
import { t } from '@/i18n'
import EventBus from '@/utils/bus'
import { PARAMETER_TYPES } from '@/dictionary/parameter-types'
import { updateValidator } from '@/setup/validate'

import SaveButton from './save-button.vue'

export default defineComponent({
  name: 'ValidationRulesConfig',
  components: {
    SaveButton,
  },
  setup() {
    const globalConfig = computed(() => store.state.globalConfig)
    const defaultConfig = computed(() => store.state.globalConfig.defaultConfig)
    const defaultRuleForm = {}
    Object.keys(PARAMETER_TYPES).forEach(key => {
      defaultRuleForm[key] = {
        reseting: false,
      }
    })
    const ruleForm = reactive(cloneDeep(defaultRuleForm))
    const ruleFormRef = ref(null)
    const ruleFormRules = {}

    const isLegalRegExp = re => {
      try {
        new RegExp(re)
        return true
      } catch (err) {
        console.log(err)
        return false
      }
    }

    const translateRuleKey = key => t(PARAMETER_TYPES[key]) || key

    const initForm = () => {
      const validationRules = cloneDeep(
        globalConfig.value.config.validationRules
      )
      const newRules = {}

      Object.keys(validationRules).forEach(key => {
        newRules[key] = validationRules[key]
        ruleFormRules[key] = [
          {
            required: true,
            validator: val =>
              val.value && isLegalRegExp(val.value) && val.message,
            message: val => {
              if (val.value?.trim() === '') {
                return t('请填写验证规则')
              }
              if (!isLegalRegExp(val.value)) {
                return t('验证规则必须是正则表达式')
              }
              if (val.message?.trim() === '') {
                return t('请填写校验提示')
              }
              return ''
            },
            trigger: 'blur',
          },
        ]
      })

      Object.assign(ruleForm, cloneDeep(defaultRuleForm), {
        ...newRules,
      })

      ruleFormRef.value.clearError()
    }

    onMounted(() => {
      initForm()
      EventBus.$on('globalConfig/fetched', initForm)
    })

    const save = () => {
      ruleFormRef.value
        .validate()
        .then(() => {
          store
            .dispatch('globalConfig/updateConfig', {
              validationRules: ruleForm,
            })
            .then(() => {
              initForm()
              updateValidator()
              bkMessage({
                theme: 'success',
                message: t('保存成功'),
              })
            })
        })
        .catch(err => {
          console.log(err)
        })
    }

    const resetAllRules = () => {
      initForm()
    }

    const recoveryRule = (rule, ruleKey) => {
      rule.message = defaultConfig.value.validationRules[ruleKey].message
      rule.value = defaultConfig.value.validationRules[ruleKey].value
    }

    return {
      translateRuleKey,
      globalConfig,
      ruleForm,
      ruleFormRules,
      ruleFormRef,
      save,
      recoveryRule,
      resetAllRules,
    }
  },
})
</script>

<style lang="scss" scoped>
.validation-rules-config {
  width: 1000px;
}

.recovery-button {
  margin-left: 8px;
  padding: 0;
}
</style>
