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

<script>
  import { computed, defineComponent, ref, watchEffect } from 'vue'
  import has from 'has'
  import store from '@/store'
  import containerPropertyService from '@/service/container/property.js'
  import containerPropertGroupService from '@/service/container/property-group.js'
  import { getContainerNodeType, getContainerInstanceService } from '@/service/container/common.js'

  export default defineComponent({
    setup() {
      const requestIds = {
        property: Symbol()
      }

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const selectedNode = computed(() => store.getters['businessHost/selectedNode'])

      const objId = computed(() => selectedNode.value.data.bk_obj_id)
      const instId = computed(() => selectedNode.value.data.bk_inst_id)
      const primaryObjId = computed(() => getContainerNodeType(objId.value))
      const isWorkload = computed(() => selectedNode.value.data.is_workload)

      const properties = ref([])
      const propertyGroups = ref([])
      const propertyMap = {}
      const propertyGroupMap = {}

      const instance = ref({})

      watchEffect(async () => {
        // 属性
        if (!has(propertyMap, primaryObjId.value)) {
          const objProperties = await containerPropertyService.getMany({
            objId: primaryObjId.value
          }, {
            requestId: requestIds.property
          })
          properties.value = objProperties
          propertyMap[primaryObjId.value] = objProperties
        } else {
          properties.value = propertyMap[primaryObjId.value]
        }

        // 属性分组
        if (!has(propertyGroupMap, primaryObjId.value)) {
          const objPropertyGroups = await containerPropertGroupService.getMany({ objId: primaryObjId.value })
          propertyGroups.value = objPropertyGroups
          propertyGroupMap[primaryObjId.value] = objPropertyGroups
        } else {
          propertyGroups.value = propertyGroupMap[primaryObjId.value]
        }

        // 实例
        const instanceService = getContainerInstanceService(primaryObjId.value)
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
