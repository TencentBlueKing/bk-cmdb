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
  import { reactive, ref, watchEffect, set, nextTick } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { v4 as uuidv4 } from 'uuid'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox'
  import GridLayout from '@/components/ui/other/grid-layout.vue'
  import GridItem from '@/components/ui/other/grid-item.vue'
  import FieldGrid from '@/components/model-manage/field-grid.vue'
  import FieldSettingForm from '@/components/model-manage/field-group/field-detail/index.vue'
  import UniqueManage from './unique-manage.vue'

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

  const emit = defineEmits(['update-field', 'update-unique'])

  const fieldLocalList = ref([])
  const uniqueLocalList = ref([])
  const settingFormComp = ref(null)
  const uniqueManageComp = ref(null)

  const uniqueSettings = reactive({
    enabled: false,
    type: 'single',
    option: []
  })

  const slider = reactive({
    title: '',
    isShow: false,
    beforeClose: null,
    view: '',
    isEditField: false
  })

  watchEffect(() => {
    fieldLocalList.value = cloneDeep(props.fieldList || [])
    uniqueLocalList.value = cloneDeep(props.uniqueList || [])
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
    console.log(field)
    slider.title = t('编辑字段')
    slider.isEditField = true
    slider.curField = field
    slider.isShow = true
    slider.beforeClose = handleSliderBeforeClose
    slider.view = 'settingForm'
  }
  const handleImportField = () => {}

  const syncField = () => {
    emit('update-field', fieldLocalList.value)
  }
  const appendField = (data) => {
    fieldLocalList.value.push({
      // 在页面中创建的数据，此id键与后台数据有意保持一致为了简化在更新查找时的逻辑
      id: uuidv4(),
      ...data
    })
  }
  const updateField = (data, id) => {
    const fieldIndex = fieldLocalList.value.findIndex(field => field.id === id)
    if (~fieldIndex) {
      set(fieldLocalList.value, fieldIndex, { id, ...data })
    }
  }
  const handleFieldSave = (data, id) => {
    if (id) {
      updateField(data, id)
    } else {
      appendField(data)
    }

    syncField()

    sliderClose()
  }

  const handleUniqueTypeChange = (type) => {
    if (type === 'union') {
      nextTick(() => {
        uniqueManageComp.value?.$el?.scrollIntoView?.()
      })
    }
  }

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
    <field-grid :field-list="fieldLocalList" @click-field="handleEditField" />

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
          @confirm="handleFieldSave"
          @cancel="handleSliderBeforeClose">
          <template #append>
            <grid-layout class="mt20" mode="form" :gap="24" :font-size="'14px'" :max-columns="1">
              <grid-item :label="$t('设置为唯一校验')">
                <bk-switcher v-model="uniqueSettings.enabled" theme="primary"></bk-switcher>
              </grid-item>
              <grid-item required :label="$t('校验类型')">
                <bk-radio-group v-model="uniqueSettings.type" @change="handleUniqueTypeChange">
                  <bk-radio-button value="single">{{$t('单独唯一')}}</bk-radio-button>
                  <bk-radio-button value="union">{{$t('联合唯一')}}</bk-radio-button>
                </bk-radio-group>
              </grid-item>
            </grid-layout>
            <grid-layout v-if="uniqueSettings.type === 'union'"
              class="unique-option-container"
              mode="form"
              :gap="0"
              :font-size="'14px'"
              :max-columns="1">
              <unique-manage :field="slider.curField" ref="uniqueManageComp"></unique-manage>
            </grid-layout>
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

  .unique-option-container {
    width: 100%;
    margin-bottom: 20px;
    margin-top: 20px;
    padding: 16px;
    background: #F5F7FB;
  }
</style>
