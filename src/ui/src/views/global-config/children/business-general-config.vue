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
  <div v-bkloading="{ isLoading: globalConfig.loading }">
    <bk-form ref="bizGeneralFormRef"
      class="config-form" :rules="bizGeneralFormRules" :label-width="labelWidth" :model="bizGeneralForm">
      <bk-form-item
        :desc="{
          width: 400,
          content: $t('平台配置业务快照名提示')
        }"
        :label="$t('业务快照名称')" :icon-offset="bizNameIconOffsetLeft" property="snapshotBizId" required>
        <business-selector class="option-value snapshot-biz-name"
          searchable
          :clearable="false"
          :placeholder="$t('请选择xx', { name: $t('业务') })"
          v-model="bizGeneralForm.snapshotBizId">
        </business-selector>
        <div class="mb0 f12 color-danger" slot="tip" v-show="showBizNameTip">
          {{ $t('业务快照名称切换提示') }}
        </div>
      </bk-form-item>
      <bk-form-item :label="$t('拓扑最大可建层级')" :icon-offset="topoLevelIconOffsetLeft" property="maxBizTopoLevel" required>
        <bk-input
          class="max-biz-topo-level-input"
          type="number"
          :min="3"
          :max="10"
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
  import { ref, reactive, computed, defineComponent, onMounted, onBeforeUnmount, watch } from 'vue'
  import store from '@/store'
  import SaveButton from './save-button.vue'
  import { bkMessage } from 'bk-magic-vue'
  import { language, t } from '@/i18n'
  import cloneDeep from 'lodash/cloneDeep'
  import EventBus from '@/utils/bus'
  import isEqual from 'lodash/isEqual'
  import BusinessSelector from '@/components/audit-history/audit-business-selector'

  export default defineComponent({
    components: {
      SaveButton,
      BusinessSelector
    },
    setup(props, { emit }) {
      const globalConfig = computed(() => store.state.globalConfig)
      const defaultForm = {
        snapshotBizId: '',
        maxBizTopoLevel: ''
      }
      const originForm = {
        snapshotBizId: ''
      }
      const bizGeneralForm = reactive(cloneDeep(defaultForm))
      const originBizGeneralForm = reactive(cloneDeep(defaultForm))
      const bizGeneralFormRef = ref(null)
      const labelWidth = computed(() => (language === 'zh_CN' ? 150 : 230))
      const bizNameIconOffsetLeft =  computed(() => (language === 'zh_CN' ? 30 : 10))
      const topoLevelIconOffsetLeft =  computed(() => (language === 'zh_CN' ? 0 : -20))
      const showBizNameTip = computed(() => bizGeneralForm.snapshotBizId !== originForm.snapshotBizId)
      const hasChange = computed(() => !isEqual(originBizGeneralForm, bizGeneralForm))

      watch(() => hasChange.value, (val) => {
        emit('has-change', val)
      })

      const initForm = () => {
        const { backend } = globalConfig.value.config
        Object.assign(bizGeneralForm, cloneDeep(defaultForm), cloneDeep(backend))
        Object.assign(originBizGeneralForm, cloneDeep(defaultForm), cloneDeep(backend))
        originForm.snapshotBizId = bizGeneralForm.snapshotBizId
        bizGeneralFormRef.value?.clearError()
      }

      onMounted(() => {
        initForm()
        EventBus.$on('globalConfig/fetched', initForm)
      })

      onBeforeUnmount(() => {
        EventBus.$off('globalConfig/fetched', initForm)
      })

      const bizGeneralFormRules = {
        snapshotBizId: [
          {
            required: true,
            message: t('请选择业务'),
            trigger: 'blur'
          }
        ],
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
        labelWidth,
        bizNameIconOffsetLeft,
        topoLevelIconOffsetLeft,
        showBizNameTip
      }
    }
  })
</script>

<style lang="scss" scoped>
  @import url("../style.scss");
  .config-form {
    width: 800px;
  }
  .snapshot-biz-name {
    width: 100%;
  }
  ::v-deep .max-biz-topo-level-input .bk-input-number {
    width: 100%;
  }
  [bk-language="en"] .snapshot-biz-name {
    width: 100%;
  }
  [bk-language="en"] .config-form {
    width: 900px;
  }
</style>
