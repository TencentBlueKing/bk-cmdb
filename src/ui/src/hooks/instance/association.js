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

import service from '@/service/instance/association'
import { ref, reactive, set } from 'vue'
import hasOwnProperty from 'has'
export default function () {
  const instanceMap = reactive({})
  const pending = ref(true)
  const find = async (options) => {
    const response = await service.findTopology(options)
    if (hasOwnProperty(instanceMap, options.bk_obj_id)) {
      set(instanceMap[options.bk_obj_id], options.bk_inst_id, response)
    } else {
      set(instanceMap, options.bk_obj_id, { [options.bk_inst_id]: response })
    }
    pending.value = false
  }
  return [{ map: instanceMap, pending }, find]
}
