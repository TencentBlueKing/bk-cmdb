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
  <div class="set-instance-layout">
    <cmdb-collapse class="property-container" v-if="hasPropertyDiff"
      :label="$t('属性变更')"
      :size="collapseSize"
      arrow-type="filled">
      <property-difference
        :property-diff="propertyDiff">
      </property-difference>
    </cmdb-collapse>

    <cmdb-collapse class="module-container" v-if="hasModuleDiff"
      :label="$t('拓扑结构变更')"
      :size="collapseSize"
      arrow-type="filled">
      <module-difference
        :module-diff="moduleDiff"
        :module-host-count="moduleHostCount">
      </module-difference>
    </cmdb-collapse>
  </div>
</template>

<script>
  import ModuleDifference from './module-difference.vue'
  import PropertyDifference from './property-difference.vue'

  export default {
    components: {
      ModuleDifference,
      PropertyDifference
    },
    props: {
      propertyDiff: {
        type: Array,
        required: true
      },
      moduleDiff: {
        type: Object,
        required: true
      },
      moduleHostCount: {
        type: Object,
        default: () => ({})
      },
      collapseSize: {
        type: String
      }
    },
    data() {
      return {
      }
    },
    computed: {
      hasPropertyDiff() {
        return this.propertyDiff?.length > 0
      },
      hasModuleDiff() {
        return this.moduleDiff.module_diffs?.length > 0
      }
    }
  }
</script>

<style lang="scss" scoped>
.set-instance-layout {
  .property-container {
    .property-difference {
      margin-top: 12px;
    }
  }

  .module-container {
    .module-difference {
      margin-top: 12px;
    }
  }

  .property-container {
    & + .module-container {
      margin-top: 32px;
    }
  }
}
</style>
