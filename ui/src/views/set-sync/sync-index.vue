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
  <cmdb-sticky-layout class="set-sync-layout"
    v-bkloading="{ isLoading: $loading([requestIds.properties, requestIds.topopath]) }">

    <div class="layout-top">
      <p class="title" v-if="isSingleSync">{{$t('请确认单个实例更改信息')}}</p>
      <i18n path="请确认实例更改信息"
        tag="p"
        class="title"
        v-else>
        <template #count>
          <span>{{setIds.length}}</span>
        </template>
      </i18n>
      <div class="type-legend">
        <span class="legend-item">
          <i class="dot changed"></i>
          {{$t('变更')}}
        </span>
        <span class="legend-item">
          <i class="dot added"></i>
          {{$t('新增')}}
        </span>
        <span class="legend-item">
          <i class="dot removed"></i>
          {{$t('删除')}}
        </span>
      </div>
    </div>

    <div class="layout-main">
      <div class="set-instance-container" v-if="isSingleSync">
        <set-instance class="instance-item"
          v-bkloading="{ isLoading: singleSet.loading }"
          :property-diff="singleSet.propertyDiff"
          :module-diff="singleSet.moduleDiff"
          :module-host-count="singleSet.moduleHostCount">
        </set-instance>
      </div>
      <div class="set-instance-group" v-else>
        <cmdb-collapse class="set-instance-container"
          v-for="diff in diffList"
          :label="setGroup[diff.setId].topoPath"
          :collapse="setGroup[diff.setId].collapse"
          arrow-type="filled"
          :key="diff.setId"
          @collapse-change="handleSetCollapseChange(diff.setId, $event)">
          <template #title>
            <div class="collapse-title">
              <span class="topopath">{{setGroup[diff.setId].topoPath}}</span>
              <span class="deny-sync-tips" v-if="diff.denySync">
                <i class="bk-icon icon-exclamation"></i>{{$t('不可同步')}}
              </span>
              <i class="bk-icon icon-close"
                v-bk-tooltips="$t('本次不同步')"
                @click.stop="handleRemove(diff)">
              </i>
            </div>
          </template>
          <set-instance class="set-instance-item"
            v-bkloading="{ isLoading: setGroup[diff.setId].loading }"
            collapse-size="small"
            :property-diff="setGroup[diff.setId].propertyDiff"
            :module-diff="setGroup[diff.setId].moduleDiff"
            :module-host-count="setGroup[diff.setId].moduleHostCount">
          </set-instance>
        </cmdb-collapse>
      </div>
    </div>

    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth
          :auth="{ type: $OPERATION.U_TOPO, relation: [bizId] }"
          v-bk-tooltips="{ content: $t(isSingleSync ? '不可同步' : '请先删除不可同步的实例'), disabled: !denySync }">
          <template slot-scope="{ disabled }">
            <bk-button
              theme="primary"
              :loading="$loading(requestIds.syncTemplateToInstances)"
              :disabled="disabled || denySync"
              @click="handleConfirmSync">
              {{$t('确认同步')}}
            </bk-button>
          </template>
        </cmdb-auth>
        <bk-button @click="handleGoback">{{$t('取消')}}</bk-button>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { MENU_BUSINESS_HOST_AND_SERVICE, MENU_BUSINESS_SET_TEMPLATE_DETAILS } from '@/dictionary/menu-symbol'
  import setInstance from './set-instance'
  import setTemplateService from '@/service/set-template'

  export default {
    components: {
      setInstance
    },
    data() {
      const id = `${this.$store.getters['objectBiz/bizId']}_${this.$route.params.setTemplateId}`
      let { syncIdMap } = this.$store.state.setFeatures
      const sessionSyncIdMap = sessionStorage.getItem('setSyncIdMap')
      if (!Object.keys(syncIdMap).length && sessionSyncIdMap) {
        syncIdMap = JSON.parse(sessionSyncIdMap)
        this.$store.commit('setFeatures/resetSyncIdMap', syncIdMap)
      }
      const setIds = syncIdMap[id] || []

      const setGroup = {}
      setIds.forEach((id) => {
        setGroup[id] = {
          propertyDiff: [], // 属性对比数据
          moduleDiff: {}, // 拓扑模板实例对比数据
          topoPath: '', // 拓扑路径
          collapse: true, // 是否展开
          loaded: false, // 是否加载过
          loading: false, // 是否加载中
          moduleHostCount: {} // 集群下模块实例的主机数
        }
      })

      return {
        setIds,
        setProperties: [],
        setGroup,
        diffList: [],
        requestIds: {
          topopath: Symbol(),
          properties: Symbol(),
          syncTemplateToInstances: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      setTemplateId() {
        return this.$route.params.setTemplateId
      },
      isSingleSync() {
        return this.diffList?.length === 1
      },
      singleSet() {
        const setId = this.diffList?.[0]?.setId
        return this.setGroup[setId] || {}
      },
      denySync() {
        return this.diffList.some(item => item.denySync)
      }
    },
    async created() {
      await this.getSetProperties()
      await this.getTopoPath()
      await this.getRemovedModuleHostStatus()

      if (this.diffList?.length > 1) {
        this.$store.commit('setTitle', this.$t('批量同步集群模板'))
      }

      // 默认展开第1个
      this.getDiffData(this.diffList?.[0]?.setId)
    },
    methods: {
      getSetProperties() {
        return this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
          params: {
            bk_biz_id: this.bizId,
            bk_obj_id: 'set',
            bk_supplier_account: this.$store.getters.supplierAccount,
          },
          config: {
            requestId: this.requestIds.properties,
            fromCache: true
          }
        }).then((data) => {
          this.setProperties = data
        })
          .catch(() => {
            this.setProperties = []
          })
      },
      getTopoPath() {
        return this.$store.dispatch('objectMainLineModule/getTopoPath', {
          bizId: this.bizId,
          params: {
            topo_nodes: this.setIds.map(setId => ({ bk_obj_id: 'set', bk_inst_id: setId }))
          },
          config: { requestId: this.requestIds.topopath }
        }).then(({ nodes }) => {
          nodes.forEach((node) => {
            const setId = node.topo_node.bk_inst_id
            this.setGroup[setId].topoPath = node.topo_path.reverse().map(path => path.bk_inst_name)
              .join(' / ')
          })
        })
      },
      async getRemovedModuleHostStatus() {
        const results = await setTemplateService.getRemovedModuleStatus(this.bizId, this.setTemplateId, {
          bk_set_ids: this.setIds
        })

        this.diffList = results.map(item => ({
          setId: item.id,
          denySync: item.has_host // 移除的模块中存在主机不允许同步
        })).sort((setA, setB) => setB.denySync - setA.denySync)
      },
      async getDiffData(setId) {
        const currentSet = this.setGroup[setId]
        try {
          currentSet.loading = true
          currentSet.collapse = false

          const data = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
            bizId: this.bizId,
            setTemplateId: this.setTemplateId,
            params: {
              bk_set_id: setId
            }
          })

          currentSet.moduleHostCount = data.module_host_count || {}

          const { attributes, ...moduleDiff  } = data.difference || {}

          // 属性变更数据，注入原始属性对象
          if (attributes) {
            currentSet.propertyDiff = attributes.map((attr) => {
              const property = this.setProperties.find(prop => prop.id === attr.id)
              return {
                property,
                ...attr
              }
            })
          }

          currentSet.moduleDiff = moduleDiff
        } finally {
          currentSet.loading = false
          currentSet.loaded = true
        }
      },
      async handleConfirmSync() {
        try {
          await this.$store.dispatch('setSync/syncTemplateToInstances', {
            bizId: this.bizId,
            setTemplateId: this.setTemplateId,
            params: {
              bk_set_ids: this.diffList.map(item => item.setId)
            },
            config: {
              requestId: this.requestIds.syncTemplateToInstances
            }
          })
          this.$success(this.$t('提交同步成功'))
          this.$routerActions.redirect({
            name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
            params: {
              templateId: this.setTemplateId
            },
            query: {
              tab: 'instance'
            }
          })
        } catch (e) {
          console.error(e)
        }
      },
      handleSetCollapseChange(setId, collapse) {
        // 打开并且未加载过或者不在加载中状态
        if (!collapse && !this.setGroup[setId].loaded && !this.setGroup[setId].loading) {
          this.getDiffData(setId)
        }
      },
      handleRemove(diff) {
        const index = this.diffList.indexOf(diff)
        if (index !== -1) {
          this.diffList.splice(index, 1)
        }
      },
      handleGoback() {
        const { moduleId } = this.$route.params
        if (moduleId) {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_AND_SERVICE,
            query: {
              node: `set-${moduleId}`
            }
          })
        } else {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
            params: {
              templateId: this.setTemplateId
            },
            query: {
              tab: 'instance'
            }
          })
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
.set-sync-layout {
  .layout-top {
    display: flex;
    margin: 24px;
    .title {
      font-size: 14px;
    }

    .type-legend {
      display: flex;
      align-items: center;
      font-size: 12px;

      .legend-item {
        margin-right: 30px;
        .dot {
          display: inline-block;
          width: 8px;
          height: 8px;
          border-radius: 50%;
          background-color: #2DCB56;
          margin-right: 2px;
          &.added {
            background-color: #2DCB56;
          }
          &.changed {
            background-color: #FF9C01;
          }
          &.removed {
            background-color: #FF5656;
          }
        }
      }
    }
  }

  .layout-main {
    margin: 24px;
  }

  .set-instance-container {
    background: #fff;
    box-shadow: 0 2px 4px 0 rgba(25, 25, 41, 0.05);
    border-radius: 2px;
    padding: 24px;

    & + .set-instance-container {
      margin-top: 16px;
    }

    .collapse-title {
      display: flex;

      .deny-sync-tips {
        display: flex;
        align-items: center;
        font-size: 12px;
        color: #FF5656;
        margin-left: 12px;
        margin-top: -2px;

        .bk-icon {
          width: 14px;
          height: 14px;
          line-height: 14px;
          text-align: center;
          color: #FFFFFF;
          background-color: #FF5656;
          border-radius: 50%;
          margin-right: 4px;
        }
      }

      .icon-close {
        color: #979BA5;
        font-size: 20px;
        margin-left: auto; // 靠右
        cursor: pointer;

        &:hover {
          color: $primaryColor;
        }
      }
    }
  }

  .set-instance-group {
    .set-instance-item {
      margin: 24px 16px 0 16px;
    }
  }

  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 24px;
    margin-top: 8px;
    .bk-button {
      min-width: 86px;

      & + .bk-button {
        margin-left: 8px;
      }
    }
    .auth-box + .bk-button {
      margin-left: 8px;
    }
    &.is-sticky {
      background-color: #fff;
      border-top: 1px solid $borderColor;
    }
  }
}
</style>
