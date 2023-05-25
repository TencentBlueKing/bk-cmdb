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
  import { ref } from 'vue'
  import { t } from '@/i18n'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE_CREATE,
    MENU_MODEL_FIELD_TEMPLATE_EDIT,
    MENU_MODEL_FIELD_TEMPLATE_BIND
  } from '@/dictionary/menu-symbol'
  // import fieldTemplateService from '@/service/field-template'

  const searchOptions = [
    {
      id: 'template_name',
      name: t('模板名称')
    },
    {
      id: 'model_name',
      name: t('模型名称')
    },
    {
      id: 'modifier',
      name: t('更新人')
    }
  ]
  const searchQuery = ref([])

  // const getTemplateList = () => {
  //   fieldTemplateService.find()
  // }

  const handleCreate = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_CREATE,
      history: true
    })
  }
  const handleEdit = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_EDIT,
      params: {
        id: 1
      },
      history: true
    })
  }
  const handleBind = () => {
    routerActions.redirect({
      name: MENU_MODEL_FIELD_TEMPLATE_BIND,
      params: {
        id: 1
      },
      history: true
    })
  }
</script>

<template>
  <div class="field-template">
    <cmdb-tips class="mb10" tips-key="fieldTemplateTips">
      {{$t('字段组合模板功能提示')}}
    </cmdb-tips>
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleCreate">
            {{$t('新建')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleBind">
            {{$t('立即绑定')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.U_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleEdit">
            {{$t('编辑')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <bk-search-select
        class="search-select"
        clearable
        :placeholder="$t('请输入模板名称/模型/更新人')"
        :show-popover-tag-change="true"
        :data="searchOptions"
        v-model="searchQuery">
      </bk-search-select>
    </div>
  </div>
</template>
<style lang="scss" scoped>
  .field-template {
    padding: 15px 20px 0;

    .toolbar {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .search-select {
        width: 480px;
      }
    }
  }
</style>
