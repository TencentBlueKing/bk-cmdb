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

import { ref, watch } from '@vue/composition-api'
import modelAssociationService from '@/service/model/association'
export default function (modelId) {
  const relations = ref([])
  const pending = ref(false)
  const refresh = async (value) => {
    if (!value) return
    pending.value = true
    relations.value = await modelAssociationService.findAll(value)
    pending.value = false
  }
  watch(modelId, refresh, { immediate: true })
  return [{ relations, pending }, refresh]
}
