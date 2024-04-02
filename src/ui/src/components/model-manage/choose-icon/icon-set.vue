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
  <ul class="icon-set" v-if="curIconList[0]">
    <li class="icon"
      ref="iconItem"
      :class="{ 'active': icon.value === value }"
      v-for="(icon, index) in curIconList"
      :key="index"
      @click="handleChooseIcon(icon.value)">
      <i :class="icon.value" v-bk-tooltips="{ content: language === 'zh_CN' ? icon.nameZh : icon.nameEn }"></i>
      <span class="checked-status"></span>
    </li>
  </ul>
  <cmdb-data-empty
    v-else
    slot="empty"
    :stuff="dataEmpty"
    @clear="handleClearFilter">
  </cmdb-data-empty>
</template>

<script>
  import { mapGetters } from 'vuex'
  export default {
    props: {
      value: {
        type: String,
        default: 'icon-cc-default'
      },
      iconList: {
        type: Array,
        default: () => []
      },
      filterIcon: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        dataEmpty: {
          type: 'search'
        }

      }
    },
    computed: {
      ...mapGetters([
        'language'
      ]),
      curIconList() {
        if (this.filterIcon) {
          // eslint-disable-next-line max-len
          return this.iconList.filter(icon => icon.nameZh.toLowerCase().indexOf(this.filterIcon.toLowerCase()) > -1 || icon.nameEn.toLowerCase().indexOf(this.filterIcon.toLowerCase()) > -1)
        }
        return this.iconList
      }
    },
    methods: {
      handleClearFilter() {
        this.$emit('clear')
      },
      handleChooseIcon(value) {
        this.$emit('input', value)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .data-empty {
      width: 560px;
    }
    .icon-set {
        width: 560px;
        display: flex;
        flex-wrap: wrap;
        padding-bottom: 10px;
        .icon {
            position: relative;
            display: flex;
            justify-content: center;
            align-items: center;
            flex: 0 0 10%;
            height: 50px;
            font-size: 24px;
            outline: 0;
            cursor: pointer;
            &:hover {
                color: #3a84ff;
                background-color: #ebf4ff;
            }
            &.active {
                color: #3a84ff;
                background-color: #ebf4ff;
                border: 1px dashed #3a84ff;
                .checked-status {
                    display: block;
                }
            }
            .checked-status {
                display: none;
                position: absolute;
                bottom: -6px;
                right: -6px;
                width: 18px;
                height: 18px;
                background-color: #2dcb56;
                border-radius: 50%;
                z-index: 2;
                &::before {
                    content: '';
                    position: absolute;
                    bottom: 5px;
                    right: 0;
                    width: 14px;
                    height: 7px;
                    border-bottom: 3px solid #ffffff;
                    border-left: 3px solid #ffffff;
                    transform: rotate(-45deg) scale(.5);
                }
            }
        }
    }
</style>
