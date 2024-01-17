/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { ref, watch } from 'vue'
import uniqueCheckService from '@/service/unique-check'
export default function (models, config) {
  const uniqueChecks = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value.length) return
    pending.value = true
    uniqueChecks.value = await uniqueCheckService.findMany(value, config)
    pending.value = false
  }
  watch(models, refresh, { immediate: true })
  return [{ uniqueChecks, pending }, { refresh }]
}
