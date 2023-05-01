<!--
 * Tencent
  components: { topSteps },is pleased to support the open source community by making 蓝鲸 available.
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
  // import { computed } from 'vue'

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

  console.log(props)

  // const emit = defineEmits(['input'])

  const handleCreate = () => {}
</script>

<template>
  <div class="field-manage">
    <div class="toolbar">
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            theme="primary"
            :disabled="disabled"
            @click="handleCreate">
            {{$t('添加字段')}}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.C_FIELD_TEMPLATE }">
        <template #default="{ disabled }">
          <bk-button
            :disabled="disabled"
            @click="handleCreate">
            {{$t('从模型导入')}}
          </bk-button>
        </template>
      </cmdb-auth>
    </div>
    <div>
      list
    </div>
  </div>
</template>
