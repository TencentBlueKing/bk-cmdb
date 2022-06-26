<script lang="ts">
  import { computed, defineComponent, ref } from '@vue/composition-api'
  import router from '@/router/index.js'
  import { t } from '@/i18n'
  import { $bkInfo, $success } from '@/magicbox/index.js'
  import routerActions from '@/router/actions'
  import ManagementForm from './children/management-form.vue'
  import store from '@/store'
  import serviceTemplateService from '@/services/service-template'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      ManagementForm
    },
    setup() {
      const managementForm = ref(null)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const loading = ref(true)
      const submitDisabled = ref(false)

      const requestIds = {
        update: Symbol('update')
      }

      const updateTemplate = async (formData) => {
        const data = {
          id: templateId.value,
          bk_biz_id: bizId.value,
          ...formData
        }
        await serviceTemplateService.update(data, { requestId: requestIds.update })

        $success(t('保存成功'))

        routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_DETAILS,
          params: {
            templateId: templateId.value
          }
        })
      }

      const handleSubmit = async () => {
        const valid = await managementForm.value.validate()
        if (!valid) {
          return
        }

        const formData = managementForm.value.getData()

        if (!formData.processes.length) {
          $bkInfo({
            title: t('确认提交'),
            subTitle: t('服务模板创建没进程提示'),
            extCls: 'bk-dialog-sub-header-center',
            confirmFn: () => {
              updateTemplate(formData)
            }
          })
          return
        }

        updateTemplate(formData)
      }

      const handleCancel = () => {
        routerActions.back()
      }

      const handleDataLoaded = () => {
        loading.value = false
      }

      return {
        bizId,
        templateId,
        managementForm,
        loading,
        submitDisabled,
        requestIds,
        handleSubmit,
        handleCancel,
        handleDataLoaded
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="edit-layout" v-bkloading="{ isLoading: loading }">
    <div class="layout-main">
      <management-form
        ref="managementForm"
        :data-id="templateId"
        :submit-disabled.sync="submitDisabled"
        @data-loaded="handleDataLoaded">
      </management-form>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_TEMPLATE, relation: [bizId, templateId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="disabled || submitDisabled"
            :loading="$loading(requestIds.update)"
            @click="handleSubmit">
            {{$t('保存')}}
          </bk-button>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
.edit-layout {
  .layout-main {
    padding: 15px 20px 0 20px;
  }

  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 20px;
    margin-top: 8px;
    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
    }
    .auth-box + .bk-button {
      margin-left: 8px;
    }
    &.is-sticky {
      background-color: #fff;
      border-top: 1px solid $borderColor;
    }
  }
}
</style>
