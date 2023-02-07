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
  <div class="association">
    <div class="options clearfix">
      <div class="fl" v-if="!readonly" v-show="activeView === viewName.list">
        <cmdb-auth class="inline-block-middle mr10"
          v-if="hasAssociation"
          :auth="HOST_AUTH.U_HOST">
          <bk-button slot-scope="{ disabled }"
            theme="primary"
            class="options-button"
            :disabled="disabled"
            @click="showCreate = true">
            {{$t('新增关联')}}
          </bk-button>
        </cmdb-auth>
        <span class="inline-block-middle mr10" v-else v-bk-tooltips="$t('当前模型暂未定义可用关联')">
          <bk-button theme="primary" class="options-button" disabled>
            {{$t('新增关联')}}
          </bk-button>
        </span>
      </div>
      <div class="fr">
        <bk-checkbox v-if="hasInstance"
          v-show="activeView === viewName.list"
          :size="16" class="options-checkbox"
          @change="handleExpandAll">
          <span class="checkbox-label">{{$t('全部展开')}}</span>
        </bk-checkbox>
        <bk-button class="options-button options-button-view"
          :theme="activeView === viewName.list ? 'primary' : 'default'"
          @click="toggleView(viewName.list)">
          {{$t('列表')}}
        </bk-button>
        <bk-button class="options-button options-button-view"
          :theme="activeView === viewName.graphics ? 'primary' : 'default'"
          @click="toggleView(viewName.graphics)">
          {{$t('拓扑')}}
        </bk-button>
        <bk-button class="options-full-screen"
          v-show="activeView === viewName.graphics"
          v-bk-tooltips="$t('全屏')"
          @click="handleFullScreen">
          <i class="icon-cc-resize-full"></i>
        </bk-button>
      </div>
    </div>
    <div class="association-view">
      <component ref="dynamicComponent" :is="activeView" v-bind="componentProps"></component>
    </div>
    <bk-sideslider v-transfer-dom :is-show.sync="showCreate" :width="800" :title="$t('新增关联')">
      <cmdb-host-association-create slot="content" v-if="showCreate"></cmdb-host-association-create>
    </bk-sideslider>
  </div>
</template>

<script>
  import cmdbHostAssociationList from './association-list.vue'
  import cmdbInstanceAssociation from '@/components/instance/association/index.vue'
  import cmdbHostAssociationCreate from './association-create.vue'
  import { mapGetters } from 'vuex'
  import authMixin from '../mixin-auth'
  import { readonlyMixin } from '../mixin-readonly'

  export default {
    name: 'cmdb-host-association',
    components: {
      cmdbHostAssociationList,
      cmdbInstanceAssociation,
      cmdbHostAssociationCreate
    },
    mixins: [authMixin, readonlyMixin],
    data() {
      return {
        viewName: {
          list: cmdbHostAssociationList.name,
          graphics: cmdbInstanceAssociation.name
        },
        activeView: cmdbHostAssociationList.name,
        showCreate: false
      }
    },
    computed: {
      ...mapGetters('hostDetails', [
        'source',
        'target',
        'allInstances'
      ]),
      hasAssociation() {
        return !!(this.source.length || this.target.length)
      },
      hasInstance() {
        return !!this.allInstances.length
      },
      componentProps() {
        if (this.activeView === cmdbInstanceAssociation.name) {
          const { host } = this.info
          return {
            objId: 'host',
            instId: host.bk_host_id,
            instName: host.bk_host_innerip
          }
        }
        return {}
      }
    },
    beforeDestroy() {
      this.$store.commit('hostDetails/toggleExpandAll', false)
    },
    methods: {
      toggleView(view) {
        this.activeView = view
      },
      handleExpandAll(expandAll) {
        this.$store.commit('hostDetails/toggleExpandAll', expandAll)
      },
      handleFullScreen() {
        this.$refs.dynamicComponent.toggleFullScreen(true)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .association {
        height: 100%;
        .association-view {
            height: calc(100% - 62px);
            @include scrollbar-y;
        }
    }
    .options {
        padding: 14px 0;
        font-size: 0;
        .options-button {
            height: 32px;
            line-height: 30px;
            font-size: 14px;
            &.options-button-view {
                margin: 0 0 0 -1px;
                border-radius: 0;
            }
        }
        .options-checkbox {
            margin: 0 15px 0 0;
            line-height: 32px;
            .checkbox-label {
                padding-left: 4px;
            }
        }
        .options-full-screen {
            width: 32px;
            height: 32px;
            padding: 0;
            text-align: center;
            margin-left: 10px;
        }
    }
</style>
