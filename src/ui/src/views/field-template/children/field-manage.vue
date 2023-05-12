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

<script setup>
  import { reactive, ref } from 'vue'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox'
  import FieldGrid from '@/components/model-manage/field-grid.vue'
  import FieldSettingForm from '@/components/model-manage/field-group/field-detail/index.vue'

  const props = defineProps({
    fieldList: {
      type: Array,
      default: () => ([])
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    }
  })

  const settingFormComp = ref(null)

  const slider = reactive({
    title: '',
    isShow: false,
    beforeClose: null,
    view: '',
    isEditField: false
  })

  console.log(props)

  const handleAddField = () => {
    slider.title = t('新建字段')
    slider.isEditField = false
    slider.curField = {}
    slider.isShow = true
    slider.beforeClose = handleSliderBeforeClose
    slider.view = 'settingForm'
  }
  const handleEditField = (field) => {
    slider.title = t('编辑字段')
    slider.isEditField = true
    slider.curField = field
    slider.isShow = true
    slider.beforeClose = handleSliderBeforeClose
    slider.view = 'settingForm'
  }
  const handleImportField = () => {}

  const handleFieldSave = () => {}

  const sliderClose = () => {
    slider.isShow = false
    slider.curField = {}
    slider.beforeClose = null
  }
  const handleSliderBeforeClose = () => {
    const hasChanged = Object.keys(settingFormComp.value.changedValues).length
    if (hasChanged) {
      return new Promise((resolve) => {
        $bkInfo({
          title: t('确认退出'),
          subTitle: t('退出会导致未保存信息丢失'),
          extCls: 'bk-dialog-sub-header-center',
          confirmFn: () => {
            sliderClose()
            resolve(true)
          },
          cancelFn: () => {
            resolve(false)
          }
        })
      })
    }
    sliderClose()
    return true
  }
  const handleSliderHidden = () => {
    sliderClose()
  }
</script>

<template>
  <div class="field-manage">
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleAddField">
            {{$t('添加字段')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            :disabled="disabled"
            @click="handleImportField">
            {{$t('从模型导入')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <div class="filter">
        <bk-input
          class="search-input"
          :placeholder="$t('请输入字段名称')"
          :right-icon="'bk-icon icon-search'" />
        <bk-button theme="default" class="unique-button">
          {{$t('唯一校验')}}
          <em class="num">0</em>
        </bk-button>
      </div>
    </div>
    <field-grid :field-list="fieldList" @click-field="handleEditField" />

    <bk-sideslider
      ref="sidesliderComp"
      v-transfer-dom
      :width="640"
      :title="slider.title"
      :is-show.sync="slider.isShow"
      :before-close="slider.beforeClose"
      @hidden="handleSliderHidden">
      <template v-if="slider.isShow" #content>
        <!-- 添加字段 -->
        <field-setting-form v-if="slider.view === 'settingForm'"
          ref="settingFormComp"
          scene="setting"
          :is-main-line-model="false"
          :is-read-only="false"
          :is-edit-field="slider.isEditField"
          :field="slider.curField"
          @save="handleFieldSave"
          @cancel="handleSliderBeforeClose">
          <template #append>
            <div>设置为唯一校验</div>
          </template>
        </field-setting-form>
      </template>
    </bk-sideslider>
  </div>
</template>

<style lang="scss" scoped>
  .field-manage {
    margin: 24px 108px;

    .toolbar {
      display: flex;
      margin-bottom: 30px;
      .auth-box {
        & + .auth-box {
          margin-left: 8px;
        }
      }

      .filter {
        margin-left: auto;
        display: flex;
        gap: 8px;
      }

      .search-input {
        width: 430px;
      }
      .unique-button {
        .num {
          font-style: normal;
          font-size: 12px;
          background: #F0F1F5;
          border-radius: 2px;
          padding: 0 .5em;
          color: #979BA5;
          margin-left: 2px;
        }
      }
    }
  }
</style>
