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
  import { computed } from 'vue'
  import Loading from '@/components/loading/index.vue'
  import SyncStatus from '@/components/ui/other/sync-status.vue'
  import { statusList, isSyncing, loadingMap } from './use-model-sync-status.js'

  const props = defineProps({
    model: Object,
    mini: Boolean
  })

  const modelStatus = computed(() => statusList.value.find(item => item.object_id === props.model.id))
  const statusValue = computed(() => {
    if (isSyncing(modelStatus.value?.status)) {
      return 'syncing'
    }
    return modelStatus.value?.status
  })
  const statusTips = computed(() => {
    if (statusValue.value === 'failure') {
      return modelStatus.value?.fail_msg
    }
    return ''
  })

  const isLoading = computed(() => loadingMap.value[props.model.id])
</script>

<template>
  <loading :loading="isLoading" v-if="!model.bk_ispaused">
    <sync-status :status="statusValue" :tips="statusTips" :mini="mini"></sync-status>
  </loading>
  <span v-else class="none-placeholder">--</span>
</template>

<style lang="scss" scoped>
.none-placeholder {
  font-size: 12px;
}
</style>
