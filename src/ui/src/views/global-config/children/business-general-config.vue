<template>
  <div v-bkloading="{ isLoading: globalConfig.loading }">
    <bk-form ref="bizGeneralFormRef" :rules="bizGeneralFormRules" :label-width="labelWidth" :model="bizGeneralForm">
      <bk-form-item :label="$t('拓扑最大可建层级')" :icon-offset="180" property="maxBizTopoLevel" required>
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
      </bk-form-item>
      <bk-form-item class="form-action-item">
        <SaveButton @save="save" :loading="globalConfig.updating"></SaveButton>
        <bk-popconfirm
          trigger="click"
          :title="$t('确认重置业务通用选项？')"
          @confirm="reset"
          :content="$t('该操作将会把业务通用选项内容重置为最后一次保存的状态，请谨慎操作！')">
          <bk-button class="action-button">{{$t('重置')}}</bk-button>
        </bk-popconfirm>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import { ref, reactive, computed, defineComponent, onMounted, onBeforeUnmount } from '@vue/composition-api'
  import store from '@/store'
  import SaveButton from './save-button.vue'
  import { bkMessage } from 'bk-magic-vue'
  import { language, t } from '@/i18n'
  import cloneDeep from 'lodash/cloneDeep'
  import EventBus from '@/utils/bus'

  export default defineComponent({
    components: {
      SaveButton
    },
    setup() {
      const globalConfig = computed(() => store.state.globalConfig)
      const defaultForm = {
        maxBizTopoLevel: ''
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

      const bizGeneralFormRules = {
        maxBizTopoLevel: [
          {
            required: true,
            message: t('请输入拓扑最大可建层级'),
            trigger: 'blur'
          }
        ]
      }

      const save = () => {
        bizGeneralFormRef.value.validate().then(() => {
          store.dispatch('globalConfig/updateConfig', {
            backend: cloneDeep(bizGeneralForm)
          })
            .then(() => {
              initForm()
              bkMessage({
                theme: 'success',
                message: t('保存成功')
              })
            })
        })
          .catch((err) => {
            console.log(err)
          })
      }

      const reset = () => {
        initForm()
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
  ::v-deep .max-biz-topo-level-input .bk-input-number{
    width: 114px;
  }
</style>
