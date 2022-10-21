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

<script lang="ts">
  import { computed, defineComponent } from 'vue'
  import ProcessDifference from './process-difference.vue'
  import PropertyDifference from './property-difference.vue'

  export default defineComponent({
    components: {
      ProcessDifference,
      PropertyDifference
    },
    props: {
      moduleId: {
        type: Number,
        required: true
      },
      templateId: {
        type: Number,
        required: true
      },
      topoPath: {
        type: String,
        default: '',
        required: true
      },
      propertyDiff: {
        type: Array,
        required: true
      },
      processDiff: {
        type: Array,
        required: true
      },
      modelProperty: {
        type: Object,
        default: () => ({}),
        required: true
      },
      collapseSize: {
        type: String
      }
    },
    setup(props) {
      const hasPropertyDiff = computed(() => props.propertyDiff?.length > 0)
      const hasProcessDiff = computed(() => props.processDiff?.length > 0)

      return {
        hasPropertyDiff,
        hasProcessDiff
      }
    }
  })
</script>

<template>
  <div class="module-instance">
    <cmdb-collapse class="property-container" v-if="hasPropertyDiff"
      :label="$t('属性变更')"
      :size="collapseSize"
      arrow-type="filled">
      <property-difference
        :module-id="moduleId"
        :template-id="templateId"
        :property-diff="propertyDiff">
      </property-difference>
    </cmdb-collapse>

    <cmdb-collapse class="process-container" v-if="hasProcessDiff"
      :label="$t('进程信息变更')"
      :size="collapseSize"
      arrow-type="filled">
      <process-difference
        :module-id="moduleId"
        :template-id="templateId"
        :topo-path="topoPath"
        :process-diff="processDiff"
        :properties="modelProperty.process">
      </process-difference>
    </cmdb-collapse>
  </div>
</template>

<style lang="scss" scoped>
.module-instance {
  .property-container {
    .property-difference {
      margin-top: 12px;
    }
  }

  .process-container {
    .process-difference {
      margin-top: 12px;
    }
  }

  .property-container {
    & + .process-container {
      margin-top: 32px;
    }
  }
}
</style>
