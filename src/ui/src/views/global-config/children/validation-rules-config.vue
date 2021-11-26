<template>
  <div class="validation-rules-config" v-bkloading="{ isLoading: globalConfig.loading }">
    <bk-form :label-width="200" ref="ruleFormRef" :model="ruleForm" :rules="ruleFormRules">
      <bk-form-item
        :icon-offset="105"
        v-for="(rule, ruleKey) of ruleForm"
        :key="ruleKey"
        :property="ruleKey" required :label="translateRuleKey(ruleKey)">
        <bk-input style="width: 300px;" :placeholder="$t('请填写验证规则')" v-model="rule.value"></bk-input>
        <bk-input style="width: 400px;" :placeholder="$t('请填写校验提示')" v-model="rule.message"></bk-input>
        <RecoveryButton
          style="margin-left: 8px;"
          :loading="rule.recovering"
          @confirm="resetRule(rule, ruleKey)">
        </RecoveryButton>
      </bk-form-item>
      <bk-form-item>
        <SaveButton @save="save" :loading="globalConfig.updating"></SaveButton>
        <ResetButton @reset="resetAllRules" :loading="globalConfig.resetting"></ResetButton>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import { ref, reactive, computed, defineComponent, onMounted } from '@vue/composition-api'
  import ResetButton from './reset-button.vue'
  import SaveButton from './save-button.vue'
  import RecoveryButton from './recovery-button.vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import { bkMessage } from 'bk-magic-vue'
  import { VALIDATION_RULE_KEYS } from '@/dictionary/validation-rule-keys'
  import cloneDeep from 'lodash/cloneDeep'
  import { updateValidator } from '@/setup/validate'
  import EventBus from '@/utils/bus'

  export default defineComponent({
    name: 'ValidationRulesConfig',
    components: {
      ResetButton,
      SaveButton,
      RecoveryButton,
    },
    setup() {
      const globalConfig = computed(() => store.state.globalConfig)
      const defaultRuleForm = {}
      Object.keys(VALIDATION_RULE_KEYS).forEach((key) => {
        defaultRuleForm[key] = {}
      })
      const ruleForm = reactive(cloneDeep(defaultRuleForm))
      const ruleFormRef = ref(null)
      const ruleFormRules = { }

      const isLegalRegExp = (re) => {
        try {
          new RegExp(re)
          return true
        } catch (err) {
          console.log(err)
          return false
        }
      }

      const translateRuleKey = key => t(VALIDATION_RULE_KEYS[key]) || key

      const initForm = () => {
        const validationRules = cloneDeep(globalConfig.value.config.validationRules)
        const newRules = {}

        Object.keys(validationRules).forEach((key) => {
          newRules[key] = validationRules[key]
          ruleFormRules[key] = [
            {
              required: true,
              validator: val => val.value && isLegalRegExp(val.value) && val.message,
              message: (val) => {
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
              trigger: 'blur'
            }
          ]
        })

        Object.assign(ruleForm, cloneDeep(defaultRuleForm), {
          ...newRules
        })

        ruleFormRef.value.clearError()
      }

      onMounted(() => {
        initForm()
        EventBus.$on('globalConfig/fetched', initForm)
      })

      const save = () => {
        ruleFormRef.value.validate().then(() => {
          store.dispatch('globalConfig/updateConfig', {
            validationRules: ruleForm
          })
            .then(() => {
              bkMessage({
                message: t('保存成功')
              })
            })
        })
          .catch((err) => {
            console.log(err)
          })
      }

      const resetAllRules = () => {
        store.dispatch('globalConfig/resetConfig', 'allValidationRules')
          .then(() => {
            initForm()
            updateValidator()
            bkMessage({
              message: t('重置成功')
            })
          })
      }

      const resetRule = (rule, ruleKey) => {
        rule.recovering = true
        store.dispatch('globalConfig/resetConfig', ruleKey)
          .then(() => {
            initForm()
            updateValidator()
            bkMessage({
              message: t('重置成功')
            })
          })
          .finally(() => {
            rule.recovering = false
          })
      }

      return {
        translateRuleKey,
        globalConfig,
        ruleForm,
        ruleFormRules,
        ruleFormRef,
        save,
        resetRule,
        resetAllRules,
      }
    }
  })
</script>

<style lang="scss" scoped>
.validation-rules-config{
  width: 1000px;
}
</style>
