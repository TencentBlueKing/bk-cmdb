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
  <div class="details-layout">
    <div v-bkloading="{ isLoading: loading }" style="height: 100%;">
      <cmdb-host-info
        ref="info"
        @info-toggle="setInfoHeight"
        @change="handleInfoChange">
      </cmdb-host-info>
      <bk-tab class="details-tab" v-if="!loading"
        type="unborder-card"
        :active.sync="active"
        :style="{
          '--infoHeight': infoHeight
        }">
        <bk-tab-panel name="property" :label="$t('主机属性')">
          <cmdb-host-property></cmdb-host-property>
        </bk-tab-panel>
        <bk-tab-panel name="service" :label="$t('服务列表')" v-if="isBusinessHost">
          <cmdb-host-service v-if="active === 'service'"></cmdb-host-service>
        </bk-tab-panel>
        <bk-tab-panel name="association" :label="$t('关联')">
          <cmdb-host-association v-if="active === 'association'"></cmdb-host-association>
        </bk-tab-panel>
        <bk-tab-panel name="history" :label="$t('变更记录')">
          <cmdb-audit-history v-if="active === 'history'"
            resource-type="host"
            obj-id="host"
            :biz-id="business"
            :resource-id="id">
          </cmdb-audit-history>
        </bk-tab-panel>
      </bk-tab>
    </div>
  </div>
</template>

<script>
  import { mapState, mapGetters } from 'vuex'
  import cmdbHostInfo from './children/info.vue'
  import cmdbHostAssociation from './children/association.vue'
  import cmdbHostProperty from './children/property.vue'
  import cmdbAuditHistory from '@/components/model-instance/audit-history'
  import cmdbHostService from './children/service-list.vue'
  import RouterQuery from '@/router/query'
  import { hostInfoProxy } from './service-proxy.js'

  export default {
    components: {
      cmdbHostInfo,
      cmdbHostAssociation,
      cmdbHostProperty,
      cmdbAuditHistory,
      cmdbHostService
    },
    data() {
      return {
        active: RouterQuery.get('tab', 'property'),
        infoHeight: '81px',
        loading: true
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapState('hostDetails', ['info', 'isBusinessHost']),
      ...mapGetters('hostDetails', ['isBusinessHost']),
      id() {
        return parseInt(this.$route.params.id, 10)
      },
      business() {
        const business = parseInt(this.$route.params.bizId || this.$route.params.business, 10)
        if (isNaN(business)) {
          return -1
        }
        return business
      }
    },
    watch: {
      info(info) {
        const hostList = info.host.bk_host_innerip.split(',')
        const host = hostList.length > 1 ? `${hostList[0]}...` : hostList[0]
        this.setBreadcrumbs(host)
      },
      id() {
        this.getData()
      },
      business() {
        this.getData()
      },
      active(active) {
        if (active !== 'association') {
          this.$store.commit('hostDetails/toggleExpandAll', false)
        }
        RouterQuery.set({
          tab: active
        })
      }
    },
    created() {
      this.getData()
    },
    methods: {
      setBreadcrumbs(ip) {
        this.$store.commit('setTitle', `${this.$t('主机详情')}【${ip}】`)
      },
      async getData() {
        try {
          this.loading = true
          await Promise.all([
            this.getProperties(),
            this.getPropertyGroups(),
            this.getHostInfo()
          ])
        } catch (error) {
          console.error(error)
        } finally {
          this.loading = false
        }
      },
      async getHostInfo() {
        try {
          const { info } = await hostInfoProxy(this.getSearchHostParams())

          if (info.length) {
            this.$store.commit('hostDetails/setHostInfo', info[0])
          } else {
            this.$routerActions.redirect({ name: 404 })
          }
        } catch (e) {
          console.error(e)
          this.$store.commit('hostDetails/setHostInfo', null)
        }
      },
      getSearchHostParams() {
        const hostCondition = {
          field: 'bk_host_id',
          operator: '$eq',
          value: this.id
        }
        return {
          bk_biz_id: this.business,
          condition: ['biz', 'set', 'module', 'host'].map(model => ({
            bk_obj_id: model,
            condition: model === 'host' ? [hostCondition] : [],
            fields: []
          })),
          ip: { flag: 'bk_host_innerip', exact: 1, data: [] }
        }
      },
      async getProperties() {
        try {
          const params = {
            bk_supplier_account: this.supplierAccount,
            bk_obj_id: 'host'
          }
          if (this.business > 0) {
            params.bk_biz_id = this.business
          }
          const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params
          })
          this.$store.commit('hostDetails/setHostProperties', properties)
        } catch (e) {
          console.error(e)
          this.$store.commit('hostDetails/setHostProperties', [])
        }
      },
      async getPropertyGroups() {
        try {
          const propertyGroups = await this.$store.dispatch('objectModelFieldGroup/searchGroup', {
            objId: 'host',
            params: this.business > 0 ? { bk_biz_id: this.business } : {}
          })
          this.$store.commit('hostDetails/setHostPropertyGroups', propertyGroups)
        } catch (e) {
          console.error(e)
          this.$store.commit('hostDetails/setHostPropertyGroups', [])
        }
      },
      setInfoHeight() {
        this.infoTimer && clearTimeout(this.infoTimer)
        this.infoTimer = setTimeout(() => {
          this.infoHeight = `${this.$refs.info.$el.offsetHeight}px`
        }, 250)
      },
      handleInfoChange() {
        this.getHostInfo()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-layout {
        overflow: hidden;
        .details-tab {
            height: calc(100% - var(--infoHeight)) !important;
            min-height: 400px;
            /deep/ {
                .bk-tab-header {
                    padding: 0;
                    margin: 0 20px;
                }
                .bk-tab-section {
                    @include scrollbar-y;
                    padding-bottom: 10px;
                }
            }
        }
    }
</style>
