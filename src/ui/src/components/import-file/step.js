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

import { ref } from '@vue/composition-api'
const current = ref(1)
const next = () => {
  current.value = Math.min(current.value + 1, 2)
}
const previous = () => {
  current.value = Math.max(current.value - 1, 1)
}

const reset = () => {
  current.value = 1
}

export default function () {
  return [current, { next, previous, reset }]
}
