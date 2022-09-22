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
  <bk-select
    v-if="type === 'selector'"
    v-model="localValue"
    v-bind="$attrs">
    <bk-option
      v-for="biz in businessList"
      :key="biz.bk_biz_id"
      :id="biz.bk_biz_id"
      :name="`[${biz.bk_biz_id}] ${biz.bk_biz_name}`">
    </bk-option>
  </bk-select>
  <span v-else>{{bizName}}</span>
</template>

<script>
  export default {
    props: {
      value: {
        type: [String, Number]
      },
      type: {
        type: String,
        default: 'selector',
        validator(type) {
          return ['selector', 'info'].includes(type)
        }
      }
    },
    data() {
      return {
        businessList: []
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      bizName() {
        const biz = this.businessList.find(biz => biz.bk_biz_id === this.value)
        return biz ? biz.bk_biz_name : '--'
      }
    },
    created() {
      this.getFullAmountBusiness()
    },
    methods: {
      async getFullAmountBusiness() {
        try {
          const data = await this.$http.get('biz/simplify?sort=bk_biz_id', {
            requestId: 'auditBusinessSelector',
            fromCache: true
          })
          this.businessList = Object.freeze(data.info || [])
        } catch (e) {
          console.error(e)
          this.businessList = []
        }
      }
    }
  }
</script>
