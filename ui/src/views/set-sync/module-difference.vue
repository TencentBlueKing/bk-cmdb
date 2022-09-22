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
  <div class="module-difference">
    <div class="diff-table">
      <div class="table-head">
        <div class="col before-col">{{$t('拓扑同步前')}}</div>
        <div class="col after-col">{{$t('拓扑同步后')}}</div>
      </div>
      <div class="table-body">
        <div class="col before-col">
          <div class="set-tree">
            <div class="node-root">
              <i class="node-icon">{{setName[0]}}</i>
              <div :title="setDetail.bk_set_name" class="root-name">{{setDetail.bk_set_name}}</div>
            </div>
            <ul class="node-children">
              <li class="node-child" v-for="(node, index) in beforeList" :key="index">
                <i class="node-icon">{{moduleName[0]}}</i>
                <div class="module-name" v-bk-overflow-tips>{{node.bk_module_name}}</div>
              </li>
            </ul>
          </div>
        </div>
        <div class="col after-col">
          <div class="set-tree">
            <div class="node-root">
              <i class="node-icon">{{setName[0]}}</i>
              <div :title="setDetail.bk_set_name" class="root-name">{{setDetail.bk_set_name}}</div>
            </div>
            <ul class="node-children">
              <li :class="['node-child', node.diff_type]" v-for="(node, index) in afterList" :key="index">
                <i class="node-icon">{{moduleName[0]}}</i>
                <div class="module-name" v-bk-overflow-tips>{{node.bk_module_name}}</div>
                <div class="tips" v-if="node.diff_type === 'remove' && existHost(node.bk_module_id)">
                  <i class="bk-icon icon-exclamation"></i>
                  <i18n path="存在主机不可同步提示" tag="p">
                    <template #btn>
                      <span class="view-btn"
                        @click="handleViewModule(node.bk_module_id)">
                        {{$t('跳转查看')}}
                      </span>
                    </template>
                  </i18n>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
  export default {
    props: {
      moduleDiff: {
        type: Object,
        required: true
      },
      moduleHostCount: {
        type: Object,
        default: () => ({})
      }
    },
    data() {
      return {
      }
    },
    computed: {
      beforeList() {
        return this.moduleDiff.module_diffs.filter(module => module.diff_type !== 'add')
      },
      afterList() {
        return this.moduleDiff.module_diffs
      },
      setDetail() {
        return this.moduleDiff.set_detail
      },
      setName() {
        const setModel = this.$store.getters['objectModelClassify/getModelById']('set') || {}
        return setModel.bk_obj_name || ''
      },
      moduleName() {
        const moduleModel = this.$store.getters['objectModelClassify/getModelById']('module') || {}
        return moduleModel.bk_obj_name || ''
      },
    },
    methods: {
      existHost(moduleId) {
        return this.moduleHostCount[moduleId] > 0
      },
      handleViewModule(moduleId) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            node: `module-${moduleId}`
          },
          history: true
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
.diff-table {
  display: grid;
  grid-template-rows: 32px auto;

  .table-head {
    display: grid;
    gap: 4px;
    grid-template-columns: 1fr 1fr;
    font-size: 12px;
    font-weight: 700;
    line-height: 32px;

    .col {
      padding-left: 24px;
      overflow: hidden;
    }
    .before-col {
      background: #F0F1F5;
    }
    .after-col {
      background: #DCDEE5;
    }
  }

  .table-body {
    display: grid;
    gap: 4px;
    grid-template-columns: 1fr 1fr;
    padding: 24px 0;
    font-size: 12px;
    background: #FAFBFD;

    .col {
      padding: 0 24px 0 90px;
      overflow: hidden;
    }
  }
}

.set-tree {
  .node-root {
    display: flex;
    line-height: 36px;
    .root-name {
      padding: 0 10px 0 0;
      font-size: 14px;
      color: #63656E;
      @include ellipsis;
    }
  }
  .node-children {
    line-height: 36px;
    margin-left: 9px;
    .node-child {
      display: flex;
      align-items: center;
      padding: 0 10px 0 50px;
      position: relative;

      &::before {
        position: absolute;
        left: 0px;
        top: -18px;
        content: "";
        width: 42px;
        height: 36px;
        border-left: 1px dashed #DCDEE5;
        border-bottom: 1px dashed #DCDEE5;
        z-index: 1;
      }

      .module-name {
        padding: 0 10px 0 0;
        font-size: 14px;
        color: #63656E;
        @include ellipsis;
      }

      .tips {
        display: flex;
        flex: none;
        align-items: center;
        font-size: 12px;
        color: #FF5656;
        .bk-icon {
          width: 16px;
          height: 16px;
          line-height: 16px;
          text-align: center;
          color: #FFFFFF;
          background-color: #FF5656;
          border-radius: 50%;
          margin-right: 4px;
        }
        .view-btn {
            color: #3A84FF;
            cursor: pointer;
        }
      }
      &.remove {
        .module-name {
          color: #FF5656;
          text-decoration: line-through;
        }
        .node-icon {
          background-color: #FF5656;
        }
      }
      &.changed {
        .module-name {
          color: #FF9C01;
        }
        .node-icon {
          background-color: #FF9C01;
        }
      }
      &.add {
        .module-name {
          color: #2DCB56;
        }
        .node-icon {
          background-color: #2DCB56;
        }
      }
    }
  }

  .node-icon {
    flex: none;
    margin: 8px 4px 8px 0px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    line-height: 20px;
    text-align: center;
    font-size: 12px;
    font-style: normal;
    color: #fff;
    background-color: #97AED6;
    z-index: 2;
  }
}
</style>
