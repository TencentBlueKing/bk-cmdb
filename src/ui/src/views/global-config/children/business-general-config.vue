<template>
  <div v-bkloading="{ isLoading: globalConfig.loading }">
    <bk-form ref="bizGeneralFormRef" :rules="bizGeneralFormRules" :label-width="labelWidth" :model="bizGeneralForm">
      <bk-form-item :label="$t('业务拓扑层级')" :icon-offset="180" property="maxBizTopoLevel" required>
        <bk-input
          class="max-biz-topo-level-input"
          type="number"
          :min="0"
          :precision="0"
          v-model.number.trim="bizGeneralForm.maxBizTopoLevel">
          <template #append>
            <div class="group-text">{{$t('层')}}</div>
          </template>
        </bk-input>
        <p class="form-item-tips">{{$t('所有业务可创建的最高层级')}}</p>
      </bk-form-item>
      <bk-form-item v-if="bizGeneralForm.itsm" :label="$t('外部主机导入流程')" property="itsm" required>
        <div class="itsm-wrapper">
          <div class="itsm-host-wrapper">
            <bk-input v-model.trim="bizGeneralForm.itsm.itsmHost" :placeholder="$t('请输入流程所在域名')"></bk-input>
            <p class="form-item-tips">{{$t('外 BG 导入主机要执行的流程')}}</p>
          </div>
          <bk-input
            type="number"
            class="itsm-id"
            :show-controls="false"
            :precision="0"
            v-model.number.trim="bizGeneralForm.itsm.itsmId"
            :placeholder="$t('流程 ID')">
          </bk-input>
        </div>
      </bk-form-item>
      <bk-form-item>
        <SaveButton @save="save" :loading="globalConfig.updating"></SaveButton>
        <ResetButton @reset="reset" :loading="globalConfig.resetting"></ResetButton>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import { ref, reactive, computed, defineComponent, onMounted, onBeforeUnmount } from '@vue/composition-api'
  import store from '@/store'
  import ResetButton from './reset-button.vue'
  import SaveButton from './save-button.vue'
  import { bkMessage } from 'bk-magic-vue'
  import { language, t } from '@/i18n'
  import cloneDeep from 'lodash/cloneDeep'
  import EventBus from '@/utils/bus'

  export default defineComponent({
    components: {
      ResetButton,
      SaveButton
    },
    setup() {
      const globalConfig = computed(() => store.state.globalConfig)
      const defaultForm = {
        maxBizTopoLevel: '',
        itsm: null
      }
      const bizGeneralForm = reactive(cloneDeep(defaultForm))
      const bizGeneralFormRef = ref(null)
      const labelWidth = computed(() => (language === 'zh_CN' ? 150 : 180))

      const initForm = () => {
        const { backend } = globalConfig.value.config
        Object.assign(bizGeneralForm, cloneDeep(defaultForm), cloneDeep(backend))
        bizGeneralFormRef.value.clearError()
      }

      onMounted(() => {
        initForm()
        EventBus.$on('globalConfig/fetched', initForm)
      })

      onBeforeUnmount(() => {
        EventBus.$off('globalConfig/fetched', initForm)
      })


      const validateUrl = url => /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)/.test(url)
      const isEmpty = value => value === '' || value === undefined || value === null
      const bizGeneralFormRules = {
        maxBizTopoLevel: [
          {
            required: true,
            message: t('请输入业务拓扑层级'),
            trigger: 'blur'
          }
        ],
        itsm: [
          {
            required: true,
            validator: async () => {
              if (isEmpty(bizGeneralForm.itsm.itsmHost)) return false

              if (!validateUrl(bizGeneralForm.itsm.itsmHost)) return false

              if (isEmpty(bizGeneralForm.itsm.itsmId)) return false

              return true
            },
            message: () => {
              if (isEmpty(bizGeneralForm.itsm.itsmHost)) {
                return t('请输入流程所在域名')
              }

              if (isEmpty(bizGeneralForm.itsm.itsmId)) {
                return t('请输入流程 ID')
              }

              if (!validateUrl(bizGeneralForm.itsm.itsmHost)) {
                return t('流程域名必须为域名地址')
              }

              return ''
            },
            trigger: 'blur'
          }
        ],
      }

      const save = () => {
        bizGeneralFormRef.value.validate().then(() => {
          store.dispatch('globalConfig/updateConfig', {
            backend: cloneDeep(bizGeneralForm)
          })
            .then(() => {
              initForm()
              bkMessage({
                message: t('保存成功')
              })
            })
        })
          .catch((err) => {
            console.log(err)
          })
      }

      const reset = () => {
        store.dispatch('globalConfig/resetConfig', 'businessCommonSetting')
          .then(() => {
            initForm()
            bkMessage({
              message: t('重置成功')
            })
          })
      }

      return {
        bizGeneralForm,
        bizGeneralFormRef,
        bizGeneralFormRules,
        globalConfig,
        save,
        reset,
        labelWidth
      }
    }
  })
</script>

<style lang="scss" scoped>
@import url("../style.scss");
.max-biz-topo-level-input{
  width: 146px;
}

.itsm-wrapper {
  display: flex;
}

.itsm-host-wrapper {
  width: 240px;
}

.itsm-id {
  width: 100px;
  margin-left: 10px;
}
</style>
