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
const file = ref(null)
const state = ref(null)
const error = ref(null)
const response = ref(null)

const change = (event) => {
  const { files: [userFile] } = event.target
  file.value = userFile
}

const clear = () => {
  file.value = null
  state.value = null
  error.value = null
  response.value = null
}

const setState = (value) => {
  state.value = value
}

const setError = (value) => {
  error.value = value
}

const setResponse = (value) => {
  response.value = value
}

export default function () {
  return [{ file, state, error, response }, { change, clear, setState, setError, setResponse }]
}
