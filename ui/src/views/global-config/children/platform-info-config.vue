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
  <div class="platform-info-config"
    :class="`is-${language}`"
    v-bkloading="{ isLoading: globalConfig.loading }">
    <div>
      <bk-popover placement="bottom">
        <p class="setting-title setting-title-ref">{{$t('网页 Title 设置：')}}</p>
        <template #content>
          <img style="width: 270px;" src="../../../assets/images/doc-title-tooltip.webp">
        </template>
      </bk-popover>
    </div>
    <bk-form
      ref="siteFormRef" :model="siteForm" :rules="siteFormRules" form-type="inline" :label-width="labelWidth">
      <bk-form-item ext-cls="site-name-form-item" :label="$t('平台名称')" required property="name">
        <bk-input style="width:240px;" v-model="siteForm.name"></bk-input>
      </bk-form-item>
      <bk-form-item :label="$t('分隔符')" required property="separator">
        <bk-input style="width:100px;" v-model="siteForm.separator"></bk-input>
      </bk-form-item>
    </bk-form>
    <p class="mt40 setting-title">{{$t('页脚信息设置：')}}</p>
    <bk-form :label-width="labelWidth" ref="footerFormRef">
      <bk-form-item :label="$t('联系方式')">
        <bk-input v-model="footerForm.contact"></bk-input>
      </bk-form-item>
      <bk-form-item :label="$t('版权信息')">
        <bk-input v-model="footerForm.copyright"></bk-input>
        <p class="form-item-tips">
          {{copyrightTips}}
        </p>
      </bk-form-item>
      <bk-form-item class="form-action-item">
        <SaveButton @save="save" :loading="globalConfig.updating"></SaveButton>
        <bk-button class="action-button" @click="preview">{{$t('预览')}}</bk-button>
        <bk-popconfirm
          trigger="click"
          :title="$t('确认重置平台信息选项？')"
          @confirm="reset"
          :content="$t('该操作将会把平台信息选项内容重置为最后一次保存的状态，请谨慎操作！')">
          <bk-button class="action-button">{{$t('重置')}}</bk-button>
        </bk-popconfirm>
      </bk-form-item>
    </bk-form>

    <bk-dialog width="1240px" v-model="previewDialogVisible" :show-footer="false">
      <TheFooter :preview-contact="footerForm.contact" :preview-copyright="footerForm.copyright"></TheFooter>
    </bk-dialog>
  </div>
</template>

<script>
  import { ref, reactive, computed, defineComponent, onMounted } from '@vue/composition-api'
  import SaveButton from './save-button.vue'
  import TheFooter from '@/views/index/children/footer.vue'
  import store from '@/store'
  import { bkMessage } from 'bk-magic-vue'
  import { t, language } from '@/i18n'
  import { changeDocumentTitle } from '@/utils/change-document-title'
  import cloneDeep from 'lodash/cloneDeep'
  import EventBus from '@/utils/bus'

  export default defineComponent({
    components: {
      SaveButton,
      TheFooter
    },
    setup() {
      const globalConfig = computed(() => store.state.globalConfig)
      const labelWidth = computed(() => (language === 'zh_CN' ? 80 : 120))
      const siteFormRef = ref(null)
      const defaultSiteForm = {
        name: '',
        separator: '',
      }
      const siteForm = reactive(cloneDeep(defaultSiteForm))
      const siteFormRules = {
        name: [{
          required: true,
          message: t('请输入平台名称'),
          trigger: 'blur'
        }],
        separator: [
          {
            required: true,
            message: t('请输入分隔符'),
            trigger: 'blur'
          }
        ]
      }
      const copyrightTips = t('支持日期变量，{{current_year}} 为当前年份，{{current_month}} 为当前月份，{{current_day}} 为当前日')
      const defaultFooterForm = {
        contact: '',
        copyright: ''
      }
      const footerForm = reactive(cloneDeep(defaultFooterForm))
      const footerFormRef = ref(null)
      const previewDialogVisible = ref(false)

      const initForm = () => {
        const { site, footer } = globalConfig.value.config
        Object.assign(siteForm, cloneDeep(defaultSiteForm), cloneDeep(site))
        Object.assign(footerForm,  cloneDeep(defaultFooterForm), cloneDeep(footer))
        siteFormRef.value.clearError()
        footerFormRef.value.clearError()
        changeDocumentTitle()
      }

      onMounted(() => {
        initForm()
        EventBus.$on('globalConfig/fetched', initForm)
      })

      const save = () => {
        siteFormRef.value.validate().then(() => {
          const newSetting = cloneDeep({
            site: siteForm,
            footer: footerForm
          })
          store.dispatch('globalConfig/updateConfig', newSetting)
            .then(() => {
              initForm()
              bkMessage({
                theme: 'success',
                message: t('保存成功')
              })
            })
        })
      }

      const preview = () => {
        previewDialogVisible.value = true
      }

      const reset = () => {
        initForm()
      }

      return {
        labelWidth,
        globalConfig,
        siteForm,
        siteFormRef,
        siteFormRules,
        footerForm,
        footerFormRef,
        previewDialogVisible,
        save,
        preview,
        reset,
        copyrightTips,
        language
      }
    }
  })
</script>

<style lang="scss" scoped>
@import url("../style.scss");
.platform-info-config {
  width: 760px;

  .bk-form {
    margin-top: 16px;
  }

  &.is-en {
    // inline 表单项不支持 label-width 自定义，所以只能在样式里强制设置
    ::v-deep .site-name-form-item .bk-label{
      width: 120px !important;
      text-align: left;
    }
  }

  &.is-zh_CN {
    // inline 表单项不支持 label-width 自定义，所以只能在样式里强制设置
    ::v-deep .site-name-form-item .bk-label{
      width: 80px !important;
      text-align: left;
    }
  }
}

::v-deep .footer{
  position: static;
  width: auto;
}

.setting-title{
  font-size: 14px;
  font-weight: 600;

  &-ref{
    border-bottom: 1px dashed #979BA5;
  }
}


</style>
