<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<script>
  import { computed, defineComponent, ref, watchEffect } from 'vue'
  import has from 'has'
  import store from '@/store'
  import containerPropertyService from '@/service/container/property.js'
  import containerPropertGroupService from '@/service/container/property-group.js'
  import { getContainerNodeType, getContainerInstanceService, getContainerPropertyObjId } from '@/service/container/common.js'

  export default defineComponent({
    setup() {
      const requestIds = {
        property: Symbol()
      }

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const selectedNode = computed(() => store.getters['businessHost/selectedNode'])

      const objId = computed(() => selectedNode.value.data.bk_obj_id)
      const instId = computed(() => selectedNode.value.data.bk_inst_id)
      const propertyObjId = computed(() => getContainerPropertyObjId(objId.value))
      const isWorkload = computed(() => selectedNode.value.data.is_workload)

      const properties = ref([])
      const propertyGroups = ref([])
      const propertyMap = {}
      const propertyGroupMap = {}

      const instance = ref({})

      watchEffect(async () => {
        // 属性（customResource 使用独立模型 id，其它 workload 仍用 workload）
        if (!has(propertyMap, propertyObjId.value)) {
          const objProperties = await containerPropertyService.getMany({
            objId: propertyObjId.value
          }, {
            requestId: requestIds.property
          })
          properties.value = objProperties
          propertyMap[propertyObjId.value] = objProperties
        } else {
          properties.value = propertyMap[propertyObjId.value]
        }

        // 属性分组
        if (!has(propertyGroupMap, propertyObjId.value)) {
          const objPropertyGroups = await containerPropertGroupService.getMany({ objId: propertyObjId.value })
          propertyGroups.value = objPropertyGroups
          propertyGroupMap[propertyObjId.value] = objPropertyGroups
        } else {
          propertyGroups.value = propertyGroupMap[propertyObjId.value]
        }

        // 实例服务由 getContainerNodeType 解析（含 customResource → workload）；具体 kind 见下方 params
        const instanceService = getContainerInstanceService(getContainerNodeType(objId.value))
        const params = { id: instId.value, bizId: bizId.value }
        if (isWorkload.value) {
          params.kind = objId.value
        }
        instance.value = await instanceService.getOne(params)
      })

      return {
        bizId,
        properties,
        propertyGroups,
        instance
      }
    }
  })
</script>

<template>
  <cmdb-details class="topology-details"
    :properties="properties"
    :property-groups="propertyGroups"
    :inst="instance"
    :show-copy="true"
    :show-options="false">
  </cmdb-details>
</template>

<style lang="scss" scoped>
.topology-details.details-layout {
  padding: 0;
  /deep/ {
    .property-group {
      padding-left: 36px;
    }
    .property-list {
      padding-left: 20px;
    }
    .details-options {
      width: 100%;
      margin: 0;
      padding-left: 56px;
    }
  }
}
</style>
