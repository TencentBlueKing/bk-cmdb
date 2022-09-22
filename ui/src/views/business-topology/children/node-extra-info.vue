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
  <div class="template-info mb10" v-if="component">
    <component :is="component" :instance="instance"></component>
  </div>
</template>

<script>
  import ServiceTemplate from './node-extra-info-service-template'
  import SetTemplate from './node-extra-info-set-template'
  export default {
    components: {
      [ServiceTemplate.name]: ServiceTemplate,
      [SetTemplate.name]: SetTemplate
    },
    props: {
      instance: {
        type: Object,
        required: true
      }
    },
    data() {
      return {}
    },
    computed: {
      selectedNode() {
        return this.$store.state.businessHost.selectedNode
      },
      component() {
        if (!this.selectedNode) {
          return null
        } if (this.isTypeOfNode('module')) {
          return ServiceTemplate.name
        } if (this.isTypeOfNode('set')) {
          return SetTemplate.name
        }
        return null
      }
    },
    methods: {
      isTypeOfNode(type) {
        return this.selectedNode && this.selectedNode.data.bk_obj_id === type
      }
    }
  }
</script>

<style lang="scss" scoped>
    .template-info {
        font-size: 14px;
        color: #63656e;
        padding: 20px 0 8px 36px;
        margin: 0 20px;
        border-bottom: 1px solid #F0F1F5;
    }
</style>
