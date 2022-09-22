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
  <div class="cmdb-host-topo-path">
    <template v-if="pending">
      <i class="path-pending"></i>
    </template>
    <template v-else>
      <span class="path" v-bk-overflow-tips>
        {{topologyList[0] && topologyList[0].path}}
      </span>
      <template v-if="!isResourcePool">
        <i class="path-single-link icon-cc-share"
          v-if="isSingle"
          @click="handleLinkToTopology(topologyList[0])">
        </i>
        <span v-else
          class="path-count"
          v-bk-tooltips="{
            interactive: true,
            boundary: 'window',
            onShow: showTips
          }">
          {{`+${topologyList.length - 1}`}}
        </span>
      </template>
    </template>
    <div v-if="!isSingle" style="display: none" ref="tooltipContent">
      <div class="path-tooltip-content">
        <div class="path-tooltip-item"
          v-for="(item, index) in topologyList"
          :key="`${item.id}_${index}`">
          <span class="path-tooltip-text" :title="item.path">{{item.path}}</span>
          <i class="path-tooltip-link icon-cc-share"
            @click="handleLinkToTopology(item)">
          </i>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import proxy from './proxy'
  import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'

  export default {
    name: 'cmdb-host-topo-path',
    props: {
      host: {
        type: Object,
        required: true
      },
      isContainerSearchMode: {
        type: Boolean,
        default: false
      },
      isResourceAssigned: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        pending: true,
        paths: [],
        topologyList: []
      }
    },
    computed: {
      bizId() {
        return this?.host?.biz?.[0]?.bk_biz_id
      },
      isResourcePool() {
        return this?.host?.biz?.[0]?.default === 1
      },
      modules() {
        return this.host?.module?.map(module => module.bk_module_id)
      },
      hostId() {
        return this.host?.host?.bk_host_id
      },
      isSingle() {
        return this.topologyList.length === 1
      }
    },
    watch: {
      host: {
        immediate: true,
        handler() {
          this.searchPath()
        }
      }
    },
    methods: {
      async searchPath() {
        proxy.isContainerSearchMode = this.isContainerSearchMode
        proxy.isResourceAssigned = this.isResourceAssigned
        try {
          this.pending = true
          this.paths = await proxy.search({
            bk_biz_id: this.bizId,
            modules: this.modules,
            hostId: this.hostId
          })
        } catch (error) {
          console.error(error)
          this.paths = {}
        } finally {
          this.generateTopologyList()
          this.pending = false
          this.$emit('path-ready', this.getFullPath())
        }
      },
      generateTopologyList() {
        const { container = [], normal = [] } = this.paths

        const normalTopoPaths = normal.map((item) => {
          const instId = item.topo_node.bk_inst_id
          const paths = item.topo_path?.slice()?.reverse()
          return {
            id: instId,
            path: paths?.map(node => node.bk_inst_name)
              .join(' / ')
          }
        })

        const containerTopoPaths = container.map(item => ({
          id: item.bk_cluster_id,
          path: `${item.biz_name} / ${item.cluster_name}`,
          isContainer: true
        }))

        this.topologyList = [...normalTopoPaths, ...containerTopoPaths || []]
      },
      getFullPath() {
        return this.topologyList.map(topo => topo.path)
      },
      showTips(inst) {
        this.$refs.tooltipContent.style.display = 'block'
        inst.setContent(this.$refs.tooltipContent)
      },
      handleLinkToTopology(topo) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE,
          query: {
            node: this.isContainerHost ? `cluster-${topo.id}` : `module-${topo.id}`
          },
          params: {
            bizId: this.bizId
          },
          history: true
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cmdb-host-topo-path {
        display: flex;
        align-items: center;
        &:hover {
            .path-single-link {
                display: inline-block;
            }
        }
        .path-pending {
            display: inline-block;
            vertical-align: middle;
            width: 16px;
            height: 16px;
            background-color: transparent;
            background-image: url("../../assets/images/icon/loading.svg");
        }
        .path {
            display: block;
            line-height: 24px;
            @include ellipsis;
        }
        .path-single-link {
            display: none;
            flex: 16px 0 0;
            font-size: 12px;
            margin-left: 5px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                opacity: .75;
            }
        }
        .path-count {
            padding: 0 6px;
            margin: 0 0 0 5px;
            height: 16px;
            line-height: 16px;
            font-size: 12px;
            border-radius: 8px;
            white-space: nowrap;
            background-color: #dcdee5;
        }
    }
    .path-tooltip-content {
        width: 300px;
        .path-tooltip-item {
            display: flex;
            align-items: center;
            &:hover {
                .path-tooltip-link {
                    display: inline-block;
                }
            }
            .path-tooltip-text {
                display: block;
                line-height: 24px;
                @include ellipsis;
            }
            .path-tooltip-link {
                display: none;
                flex: 16px 0 0;
                font-size: 12px;
                margin-left: 5px;
                color: $primaryColor;
                cursor: pointer;
                &:hover {
                    opacity: .75;
                }
            }
        }
    }
</style>
