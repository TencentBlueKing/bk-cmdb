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
  <div class="model-export-setting">
    <bk-form ref="settingFormRef" class="setting-form" :model="settingForm" :rules="settingFormRules">
      <bk-form-item :label="t('压缩包名')" required property="fileName">
        <bk-input class="setting-form-input" :maxlength="40" :show-word-limit="true"
          v-model="settingForm.fileName" :placeholder="t('压缩包名仅支持大小写英文字母、数字、- 或 _')" />
      </bk-form-item>
      <bk-form-item :label="t('文件加密')" required>
        <bk-radio-group class="encrypt-radio-group" v-model="isEncrypt">
          <bk-radio :value="true">{{t('是')}}</bk-radio>
          <bk-radio :value="false">{{t('否')}}</bk-radio>
        </bk-radio-group>
      </bk-form-item>
      <template v-if="isEncrypt">
        <bk-form-item :label="t('密码设置')" required property="password">
          <bk-input class="setting-form-input"
            v-model="settingForm.password" type="password" :placeholder="t('长度 6-20 个字符，必须包含英文字母、数字和特殊符号')"></bk-input>
        </bk-form-item>
        <bk-form-item :label="t('二次确认')" required property="confirmedPassword">
          <bk-input class="setting-form-input" v-model="settingForm.confirmedPassword" type="password"
            :placeholder="t('请输入同样的密码，以确认密码准确')"></bk-input>
        </bk-form-item>
        <bk-form-item :label="t('文件有效期')" required property="expirationTime">
          <bk-radio-group v-model="settingForm.expirationTime">
            <bk-radio-button :value="0">{{t('永久')}}</bk-radio-button>
            <bk-radio-button :value="1">{{t('1天')}}</bk-radio-button>
            <bk-radio-button :value="3">{{t('3天')}}</bk-radio-button>
            <bk-radio-button :value="7">{{t('7天')}}</bk-radio-button>
            <bk-radio-button :value="30">{{t('一个月')}}</bk-radio-button>
            <bk-radio-button :value="null">{{t('自定义')}}</bk-radio-button>
          </bk-radio-group>
          <bk-input
            v-show="settingForm.expirationTime === null"
            type="number"
            class="expiration-time-input"
            v-model.number="customExpirationTime" :placeholder="t('请输入整数')">
            <template slot="append">
              <div class="group-text">{{t('天')}}</div>
            </template>
          </bk-input>
        </bk-form-item>
      </template>
    </bk-form>
  </div>
</template>

<script>
  import { defineComponent, ref, reactive, watch } from '@vue/composition-api'
  import { t } from '@/i18n'
  import moment from 'moment'

  export default defineComponent({
    name: 'ModelExportSetting',
    model: {
      prop: 'value',
      event: 'value-change'
    },
    props: {
      /**
       * 模型导出设置，支持 v-model
       * @property {Object} value
       * @property {String} value.fileName 压缩包名
       * @property {String} value.password 解压缩密码
       * @property {String} value.expirationTime 文件有效期
       */
      value: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props, { emit }) {
      const settingFormRef = ref(null)
      const settingForm = reactive({
        fileName: `bk_cmdb_model_export_${moment().format('YYYYMMDDHHMMSS')}`,
        password: '',
        confirmedPassword: '',
        expirationTime: 0,
      })

      const customExpirationTime = ref(0)

      const passwordRules = [
        {
          required: true,
          trigger: 'blur',
          message: t('请输入密码')
        },
        {
          required: true,
          trigger: 'blur',
          validator: value => /^(?=.*[0-9])(?=.*[a-zA-Z])(?=.*[^(0-9a-zA-Z)]).{6,16}$/.test(value),
          message: t('长度 6-20 个字符，必须包含英文字母、数字和特殊符号'),
        },
      ]

      const settingFormRules = {
        fileName: [
          {
            required: true,
            trigger: 'blur',
            message: t('请输入压缩包名')
          },
          {
            required: true,
            trigger: 'blur',
            validator: value => /^[a-zA-Z0-9-_]+$/.test(value),
            message: t('压缩包名仅支持大小写英文字母、数字、- 或 _')
          },
        ],
        password: [
          ...passwordRules
        ],
        confirmedPassword: [
          ...passwordRules,
          {
            required: true,
            validator: value => value === settingForm.password,
            trigger: 'blur',
            message: t('两次密码输入必须一致'),
          }
        ],
      }

      const isEncrypt = ref(false)

      watch(settingForm, () => {
        emit('value-change', {
          fileName: settingForm.fileName,
          password: settingForm.password,
          expirationTime: settingForm.expirationTime || customExpirationTime.value,
        })
      }, {
        immediate: true,
        deep: true
      })

      const validate = () => settingFormRef.value.validate()

      return {
        t,
        isEncrypt,
        customExpirationTime,
        settingForm,
        settingFormRules,
        settingFormRef,
        validate
      }
    },
  })
</script>

<style lang="scss" scoped>
.model-export-setting{
  display: flex;
  justify-content: center;
  height: 100%;
  overflow-y: auto;
}

.setting-form{
  margin-top: 48px;

  &-input{
    width: 427px;
  }
}

.encrypt-radio-group {
  .bk-form-radio + .bk-form-radio {
    margin-left: 28px;
  }
}

.expiration-time-input{
  margin-top: 13px;
  .group-text{
    width: 87px;
    text-align: center;
    height: 100%;
    background-color: #f0f1f5;
  }
}
</style>
