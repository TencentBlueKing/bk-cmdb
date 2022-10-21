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

import { ref } from 'vue'
import { find } from '@/service/association/index.js'

/**
 * 加载所有关联关系类型
 * @returns {boolean} loading 类型加载状态
 * @returns {object} associations 所有关联关系类型
 */
export const useAssociations = () => {
  const loading = ref(false)
  const associations = ref([])

  loading.value = true

  find()
    .then(({ info }) => {
      associations.value = info
    })
    .finally(() => {
      loading.value = false
    })

  return {
    loading,
    associations
  }
}
