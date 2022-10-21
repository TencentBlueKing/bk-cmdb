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

<template>
  <loading :loading="loading">{{count}}</loading>
</template>

<script>
  import { defineComponent, computed, ref, watchEffect } from 'vue'
  import loading from '@/components/loading/index.vue'
  import { instanceCounts } from './use-instance-count.js'

  export default defineComponent({
    components: { loading },
    props: {
      objId: {
        type: String,
        required: true
      }
    },
    setup(props) {
      const { objId } = props

      // 累积每一次的结果
      const list = computed(() => {
        const list = []
        instanceCounts.value.forEach((item) => {
          list.push(...item)
        })
        return list
      })

      // 找到当前模型实例数据
      const instance = computed(() => list.value.find(inst => inst.bk_obj_id === objId))

      // 确定loading态
      const loading = ref(true)
      watchEffect(() => {
        loading.value = list.value?.findIndex(item => item.bk_obj_id === objId) === -1
      })

      // 当前模型实例count
      const count = computed(() => {
        if (instance.value?.error) {
          return '--'
        }
        return instance.value?.inst_count
      })

      return {
        loading,
        count
      }
    }
  })
</script>
