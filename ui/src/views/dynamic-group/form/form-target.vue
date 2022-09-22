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
  <div class="search-target">
    <div class="target-item"
      v-for="item in list"
      :key="item.bk_obj_id"
      :class="{
        disabled: disabled,
        selected: selected === item.bk_obj_id
      }"
      @click="setSelected(item)">
      <span class="item-checkbox"></span>
      <span class="item-info">
        <span class="info-name">{{item.bk_obj_name}}</span>
        <span class="info-desc">{{item.desc}}</span>
      </span>
    </div>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      disabled: Boolean
    },
    data() {
      return {
        list: Object.freeze([{
          bk_obj_id: 'host',
          bk_obj_name: this.$t('主机'),
          desc: this.$t('目标为主机列表')
        }, {
          bk_obj_id: 'set',
          bk_obj_name: this.$t('集群'),
          desc: this.$t('目标为集群列表')
        }])
      }
    },
    computed: {
      selected: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value, this.value)
        }
      }
    },
    methods: {
      setSelected(item) {
        if (this.disabled) return
        this.selected = item.bk_obj_id
      }
    }
  }
</script>

<style lang="scss" scoped>
    .search-target {
        display: flex;
    }
    .target-item {
        display: flex;
        flex: 1;
        height: 56px;
        align-items: center;
        border: 1px solid #C4C6CC;
        border-right-width: 0;
        cursor: pointer;
        &.disabled {
            cursor: not-allowed;
            border-color: #DCDEE5;
            .item-checkbox {
                border-color: #DCDEE5;
                background-color: #FAFBFD;
            }
            .item-info {
                .info-name {
                    color: #979BA5;
                }
                .info-desc {
                    color: #C4C6CC;
                }
            }
        }
        &.selected {
            border-color: $primaryColor;
            border-right-width: 1px;
            & + .target-item {
                border-left-width: 0;
            }
            .item-checkbox {
                border-color: $primaryColor;
                background-color: $primaryColor;
            }
            .item-info {
                .info-name,
                .info-desc {
                    color: $primaryColor;
                }
            }
        }
        &:first-child {
            border-radius: 2px 0 0 2px;
        }
        &:last-child {
            border-radius: 0 2px 2px 0;
            border-right-width: 1px;
        }
        .item-checkbox {
            width: 16px;
            height: 16px;
            padding: 3px;
            margin: 0 16px;
            border: 1px solid #979BA5;
            border-radius: 50%;
            background-clip: content-box;
        }
        .item-info {
            line-height: 16px;
            .info-name {
                display: block;
                font-size: 12px;
                font-weight: 700;
                color: $textColor;
            }
            .info-desc {
                font-size: 12px;
                color: #979ba5;
            }
        }
    }
</style>
