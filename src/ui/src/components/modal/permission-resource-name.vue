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
  <div>
    <div :class="['skeleton', { 'full-width': relations.length > 1 }]" v-if="fetching"></div>
    <div v-else>{{names.join(' / ') || '--'}}</div>
  </div>
</template>

<script>
  import { IAM_VIEWS_INST_NAME } from './permission-resource-name.js'
  export default {
    props: {
      relations: {
        type: Array,
        default() {
          return []
        }
      }
    },
    data() {
      return {
        fetching: false,
        names: []
      }
    },
    watch: {
      relations() {
        this.fetchName()
      }
    },
    created() {
      this.fetchName()
    },
    methods: {
      async fetchName() {
        this.names = []
        const nameReq = []
        this.relations.forEach((relation) => {
          const [type, id] = relation
          nameReq.push(IAM_VIEWS_INST_NAME[type](this, id, this.relations))
        })
        try {
          this.fetching = true
          const result = await Promise.all(nameReq)
          result.forEach((name, index) => {
            this.names.push(`【${this.relations[index][2]}】${name}`)
          })
        } catch (error) {
          console.error(error)
        } finally {
          this.fetching = false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .skeleton {
        width: 50%;
        height: 24px;
        margin: 2px 0;
        background-image: linear-gradient(90deg, #f2f2f2 25%, #e6e6e6 37%, #f2f2f2 63%);
        background-size: 400% 100%;
        background-position: 100% 50%;
        animation: skeleton-loading 1.4s ease infinite;

        &.full-width {
            width: 100%;
        }
    }
    @keyframes skeleton-loading {
        0% {
            background-position: 100% 50%;
        }
        100% {
            background-position: 0 50%;
        }
    }
</style>
