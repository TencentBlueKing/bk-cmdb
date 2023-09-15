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
  <bk-select v-if="display === 'selector'"
    searchable
    :clearable="false"
    :readonly="readonly"
    :disabled="disabled"
    :placeholder="$t('请选择xx', { name: $t('地域') })"
    :loading="$loading(request)"
    v-model="selected">
    <bk-option v-for="region in regions"
      :key="region.bk_region"
      :name="region.bk_region_name"
      :id="region.bk_region">
      <div :class="['region-info', { selected: selected === region.bk_region }]">
        <span class="region-name" v-bk-overflow-tips>{{region.bk_region_name}}</span>
        <span class="region-host-count">
          {{region.bk_host_count}}
        </span>
      </div>
    </bk-option>
  </bk-select>
  <span v-else>{{getRegionInfo()}}</span>
</template>

<script>
  import symbols from '../common/symbol'
  export default {
    name: 'task-region-selector',
    props: {
      account: Number,
      display: {
        type: String,
        default: 'selector'
      },
      readonly: Boolean,
      disabled: Boolean,
      value: {
        type: [String, Number]
      },
      withCount: {
        type: Boolean,
        default: true
      }
    },
    data() {
      return {
        regions: [],
        request: symbols.get(`taskRegionSelection-${this.account}`)
      }
    },
    computed: {
      selected: {
        get() {
          return this.value
        },
        set(value, oldValue) {
          this.$emit('input', value)
          this.$emit('change', value, oldValue)
        }
      }
    },
    created() {
      // 为0时是默认管控区域，无地域信息
      this.account && this.getRegions()
    },
    methods: {
      async getRegions() {
        try {
          const regions = await this.$store.dispatch('cloud/resource/findRegion', {
            params: {
              bk_account_id: this.account,
              with_host_count: this.withCount
            },
            config: {
              requestId: this.request,
              fromCache: true,
              globalError: false
            }
          })
          if (!regions) {
            throw new Error('Request account regions failed.')
          }
          this.regions = regions
          this.selected = regions.length ? regions[0].bk_region : ''
        } catch (e) {
          console.error(e)
          this.regions = []
        }
      },
      getRegionInfo() {
        const region = this.regions.find(region => region.bk_region === this.value)
        return region ? region.bk_region_name : (this.value || '--')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .region-info {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 14px;
        &.selected,
        &:hover {
            .region-host-count {
                background-color: #a2c5fd;
                color: #fff;
            }
        }
        .region-name {
            @include ellipsis;
        }
        .region-host-count {
            display: flex;
            height: 18px;
            padding: 0 5px;
            align-items: center;
            margin-left: 15px;
            border-radius: 2px;
            background-color: #f0f1f5;
            color: #979ba5;
            font-size: 12px;
        }
    }
</style>
