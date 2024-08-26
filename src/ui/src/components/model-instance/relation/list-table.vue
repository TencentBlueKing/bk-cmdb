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
  <div class="table" v-bkloading="{ isLoading: loading }">
    <div class="table-info clearfix" @click="expanded = !expanded">
      <div class="info-title fl">
        <i class="icon bk-icon icon-right-shape"
          :class="{ 'is-open': expanded }">
        </i>
        <span class="title-text">{{title}}</span>
        <span class="title-count">({{associationInstances.length}})</span>
      </div>
      <div class="info-pagination fr" v-show="pagination.count" @click.stop>
        <span class="pagination-info">{{getPaginationInfo()}}</span>
        <span class="pagination-toggle">
          <i class="pagination-icon bk-icon icon-cc-arrow-down left"
            :class="{ disabled: pagination.current === 1 }"
            @click="togglePage(-1)">
          </i>
          <i class="pagination-icon bk-icon icon-cc-arrow-down right"
            :class="{ disabled: pagination.current === totalPage }"
            @click="togglePage(1)">
          </i>
        </span>
      </div>
    </div>
    <bk-table class="association-table"
      v-show="expanded"
      :data="list"
      :max-height="462"
      empty-block-class-name="empty-block"
      :row-style="{ cursor: 'pointer' }"
      @row-click="handleShowDetails">
      <bk-table-column v-for="(column, index) in header"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        :class-name="index === 0 ? 'is-highlight' : ''"
        :show-overflow-tooltip="$tools.isShowOverflowTips(column.property)">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :show-unit="false"
            :value="row[column.id]"
            :property="column.property"
            :instance="row"
            show-on="cell">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" v-if="table.stuff.type !== 'permission'">
        <template slot-scope="{ row }">
          <cmdb-auth :auth="getInstanceAuth(row)" :ignore-passed-auth="true">
            <bk-button slot-scope="{ disabled }"
              text
              theme="primary"
              :disabled="disabled"
              @click.stop="showTips($event, row)">
              {{$t('取消关联')}}
            </bk-button>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
    </bk-table>
    <div class="confirm-tips" ref="confirmTips" v-show="confirm.show">
      <p class="tips-content">{{$t('确认取消')}}</p>
      <div class="tips-option">
        <bk-button class="tips-button" theme="primary" @click.stop="cancelAssociation">{{$t('确认')}}</bk-button>
        <bk-button class="tips-button" theme="default" @click.stop="hideTips">{{$t('取消')}}</bk-button>
      </div>
    </div>
  </div>
</template>

<script>
  import bus from '@/utils/bus.js'
  import { mapGetters } from 'vuex'
  import authMixin from '../mixin-auth'
  import instanceService from '@/service/instance/instance'
  import hostSearchService from '@/service/host/search'
  import businessSetService from '@/service/business-set/index.js'
  import {
    BUILTIN_MODELS,
    BUILTIN_MODEL_PROPERTY_KEYS,
    BUILTIN_MODEL_RESOURCE_TYPES
  } from '@/dictionary/model-constants.js'
  import { translateAuth } from '@/setup/permission'
  import { getHostInfoTitle } from '@/utils/util'

  export default {
    name: 'cmdb-relation-list-table',
    mixins: [authMixin],
    props: {
      type: {
        type: String,
        required: true
      },
      targetObjId: {
        type: String,
        required: true
      },
      associationType: {
        type: Object,
        required: true
      },
      associationInstances: {
        type: Array,
        required: true
      },
      objId: {
        type: String,
        required: true
      }
    },
    data() {
      return {
        properties: [],
        list: [],
        hostOriginalList: [],
        pagination: {
          count: 0,
          current: 1,
          size: 10
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        expanded: false,
        confirm: {
          instance: null,
          item: null,
          target: null,
          id: null,
          show: false
        }
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['models', 'getModelById']),
      model() {
        return this.getModelById(this.targetObjId) || {}
      },
      title() {
        const desc = this.type === 'source' ? this.associationType.src_des : this.associationType.dest_des
        return `${desc}-${this.model.bk_obj_name}`
      },
      propertyRequest() {
        return `get_${this.targetObjId}_association_list_table_properties`
      },
      instanceRequest() {
        return `get_${this.targetObjId}_association_list_table_instances`
      },
      page() {
        // 每次in当前页的id查询，因此page参数每次都是0-size
        return {
          limit: this.pagination.size,
          start: 0
        }
      },
      totalPage() {
        return Math.ceil(this.pagination.count / this.pagination.size)
      },
      isSource() {
        return this.type === 'source'
      },
      // 模型实例id，用于查询模型实例数据
      instanceIds() {
        // eslint-disable-next-line max-len
        return this.associationInstances.map(instance => (this.isSource ? instance.bk_asst_inst_id : instance.bk_inst_id))
      },
      currentPageInstanceIds() {
        const start = (this.pagination.current - 1) * this.pagination.size
        return this.instanceIds.slice(start, start + this.pagination.size)
      },
      header() {
        const fixedPropertyIds = this.targetObjId === BUILTIN_MODELS.HOST ? ['bk_host_innerip', 'bk_host_innerip_v6'] : []
        const headerProperties = this.$tools.getHeaderProperties(this.properties, [], fixedPropertyIds)
        const header = headerProperties.map(property => ({
          id: property.bk_property_id,
          name: this.$tools.getHeaderPropertyName(property),
          property
        }))
        return header
      },
      loading() {
        return this.$loading([
          this.propertyRequest,
          this.instanceRequest
        ])
      },
      resourceType() {
        return this.$parent.resourceType
      },
      authResources() {
        const authTypes = {
          [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS]]: this.INST_AUTH.U_BUSINESS,
          [BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]]: this.INST_AUTH.U_BUSINESS_SET
        }
        return authTypes[this.resourceType] || this.INST_AUTH.U_INST
      }
    },
    watch: {
      associationInstances: {
        handler(associationInstances) {
          associationInstances.length && this.expanded && this.getData()
        },
        immediate: true
      },
      expanded(expanded) {
        expanded && this.getData()
      }
    },
    created() {
      bus.$on('expand-all-relation-table', (expandAll) => {
        this.expanded = expandAll
      })
    },
    beforeDestroy() {
      bus.$off('expand-all-relation-table')
    },
    methods: {
      getData() {
        this.getProperties()
        this.getInstances()
      },
      getModelPermission() {
        const permissions = {
          [BUILTIN_MODELS.BUSINESS]: {
            type: this.$OPERATION.R_BUSINESS,
            relation: []
          },
          [BUILTIN_MODELS.BUSINESS_SET]: {
            type: this.$OPERATION.R_BUSINESS_SET,
            relation: []
          }
        }

        return translateAuth(permissions[this.targetObjId])
      },
      async getProperties() {
        try {
          this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params: {
              bk_obj_id: this.targetObjId
            },
            config: {
              fromCache: true,
              requestId: this.propertyRequest,
              globalPermission: false
            }
          })
        } catch (e) {
          if (e.permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission: e.permission }
            }
          }
        }
      },
      async getInstances() {
        let promise
        const config = {
          requestId: this.instanceRequest,
          cancelPrevious: true,
          globalError: false,
          globalPermission: false
        }
        try {
          switch (this.targetObjId) {
            case BUILTIN_MODELS.HOST:
              promise = this.getHostInstances(config)
              break
            case BUILTIN_MODELS.BUSINESS:
              promise = this.getBusinessInstances(config)
              break
            case BUILTIN_MODELS.BUSINESS_SET:
              promise = this.getBusinessSetInstances(config)
              break
            default:
              promise = this.getModelInstances(config)
          }
          const data = await promise

          const dataListKeys = {
            [BUILTIN_MODELS.BUSINESS_SET]: 'list'
          }
          const dataListKey = dataListKeys[this.targetObjId] || 'info'

          this.list = data[dataListKey]

          // 此处总数需要使用总实例数，而非data.count因使用in当前页的id查询data.count只是当前页的条数
          this.pagination.count = this.associationInstances?.length

          // 向前翻一页
          if (data.count && this.pagination.count && !data[dataListKey].length) {
            this.pagination.current -= 1
            this.getInstances()
          }

          // 业务/业务集有自身的查看权限，当无权限时接口不会返回permission，此处通过存在关联实例但是未查询到数据判定为无权限
          if ([BUILTIN_MODELS.BUSINESS, BUILTIN_MODELS.BUSINESS_SET].includes(this.targetObjId)
            && this.associationInstances?.length > 0
            && data.count === 0) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission: this.getModelPermission() }
            }
          }
        } catch (e) {
          if (e.permission) {
            this.table.stuff = {
              type: 'permission',
              payload: { permission: e.permission }
            }
          }
        }
      },
      getInstanceAuth(row) {
        const auth = [this.authResources]
        switch (this.model.bk_obj_id) {
          case BUILTIN_MODELS.BUSINESS: {
            auth.push({
              type: this.$OPERATION.U_BUSINESS,
              relation: [row.bk_biz_id]
            })
            break
          }
          case BUILTIN_MODELS.BUSINESS_SET: {
            auth.push({
              type: this.$OPERATION.U_BUSINESS_SET,
              relation: [row[BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID]]
            })
            break
          }
          case BUILTIN_MODELS.HOST: {
            const originalData = this.hostOriginalList.find(data => data.host.bk_host_id === row.bk_host_id)
            const [biz] = originalData.biz
            if (biz.default === 0) {
              auth.push({
                type: this.$OPERATION.U_HOST,
                relation: [biz.bk_biz_id, row.bk_host_id]
              })
            } else {
              const [module] = originalData.module
              auth.push({
                type: this.$OPERATION.U_RESOURCE_HOST,
                relation: [module.bk_module_id, row.bk_host_id]
              })
            }
            break
          }
          default: {
            auth.push({
              type: this.$OPERATION.U_INST,
              relation: [this.model.id, row.bk_inst_id]
            })
          }
        }
        return auth
      },
      getHostInstances(config) {
        const models = ['biz', 'set', 'module', 'host']
        const hostCondition = {
          field: 'bk_host_id',
          operator: '$in',
          value: this.currentPageInstanceIds
        }
        const condition = models.map(model => ({
          bk_obj_id: model,
          fields: [],
          condition: model === 'host' ? [hostCondition] : []
        }))
        return hostSearchService.getHosts({
          params: {
            bk_biz_id: -1,
            condition,
            id: {
              data: [],
              exact: 0,
              flag: 'bk_host_innerip|bk_host_outerip'
            },
            page: {
              ...this.page,
              sort: 'bk_host_id'
            }
          },
          config
        }).then((data) => {
          this.hostOriginalList = data.info
          return {
            count: data.count,
            info: data.info.map(item => item.host)
          }
        })
      },
      getBusinessInstances(config) {
        return this.$store.dispatch('objectBiz/searchBusiness', {
          params: {
            condition: {
              bk_biz_id: {
                $in: this.currentPageInstanceIds
              }
            },
            fields: [],
            page: {
              ...this.page,
              sort: 'bk_biz_id'
            }
          },
          config
        })
      },
      getBusinessSetInstances(config) {
        const params = {
          fields: [],
          bk_biz_set_filter: {
            condition: 'AND',
            rules: [{ field: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID, operator: 'in', value: this.currentPageInstanceIds }]
          },
          page: this.page
        }

        return businessSetService.find(params, config)
      },
      getModelInstances(config) {
        return instanceService.find({
          bk_obj_id: this.targetObjId,
          params: {
            fields: [],
            page: {
              ...this.page,
              sort: 'bk_inst_id'
            },
            conditions: {
              condition: 'AND',
              rules: [{
                field: 'bk_inst_id',
                operator: 'in',
                value: this.currentPageInstanceIds
              }]
            }
          },
          config
        })
      },
      async cancelAssociation() {
        try {
          const asstInstance = this.associationInstances.find((instance) => {
            const key = this.isSource ? 'bk_asst_inst_id' : 'bk_inst_id'
            // 关联实例数据中的模型实例id与模型实例id比较
            return instance[key] === this.confirm.id
          })
          await this.$store.dispatch('objectAssociation/deleteInstAssociation', {
            id: asstInstance.id,
            objId: this.objId,
            config: { data: {} }
          })
          this.hideTips()
          this.$success(this.$t('取消关联成功'))
          this.$emit('delete-association', asstInstance.id)
        } catch (e) {
          console.error(e)
        }
      },
      getPaginationInfo() {
        return this.$tc('页码', this.pagination.current, {
          current: this.pagination.current,
          count: this.pagination.count,
          total: this.totalPage
        })
      },
      togglePage(step) {
        const { current } = this.pagination
        const newCurrent = current + step
        if (newCurrent < 1 || newCurrent > this.totalPage) {
          return false
        }
        this.pagination.current = newCurrent
        this.getInstances()
      },
      hideTips() {
        this.confirm.instance && this.confirm.instance.hide()
      },
      getRowInstId(item) {
        const specialModel = ['host', 'biz', 'set', 'module']
        const mapping = {
          [BUILTIN_MODELS.BUSINESS_SET]: BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID
        }
        specialModel.forEach(key => (mapping[key] = `bk_${key}_id`))
        return item[mapping[this.targetObjId] || 'bk_inst_id']
      },
      showTips(event, item) {
        this.confirm.item = item
        this.confirm.id = this.getRowInstId(item)
        this.confirm.instance = this.$bkPopover(event.target, {
          content: this.$refs.confirmTips,
          theme: 'light',
          zIndex: 9999,
          width: 200,
          trigger: 'click',
          boundary: 'window',
          arrow: true,
          interactive: true,
          onHidden: () => {
            this.confirm.instance && this.confirm.instance.destroy()
            this.confirm.instance = null
          }
        })
        this.confirm.show = true
        this.$nextTick(() => {
          this.confirm.instance.show()
        })
      },
      async handleShowDetails(row, event, column) {
        if (this.header.findIndex(prop => prop.id === column.property) !== 0) {
          return
        }
        const showInstanceDetails = await import('@/components/instance/details')
        const nameMapping = {
          host: 'bk_host_innerip',
          biz: 'bk_biz_name',
          [BUILTIN_MODELS.BUSINESS_SET]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME]
        }
        const idMapping = {
          host: 'bk_host_id',
          biz: 'bk_biz_id',
          [BUILTIN_MODELS.BUSINESS_SET]: [BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].ID]
        }
        let title = `${row[nameMapping[this.targetObjId] || 'bk_inst_name']}`
        if (this.targetObjId === 'host') {
          const {
            bk_host_id: hostId,
            bk_cloud_id: cloud,
            bk_host_innerip: ip,
            bk_host_innerip_v6: ipv6
          } = row
          const cloudId = this.$tools.getPropertyCopyValue(cloud, 'foreignkey')
          title = getHostInfoTitle(ip, ipv6, cloudId, hostId)
        }

        showInstanceDetails.default({
          bk_obj_id: this.targetObjId,
          bk_inst_id: row[idMapping[this.targetObjId] || 'bk_inst_id'],
          title: `${this.model.bk_obj_name}-${title}`
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .table {
        margin: 0 0 12px 0;
    }
    .table-info {
        height: 42px;
        padding: 0 20px;
        border-radius: 2px 2px 0 0;
        line-height: 42px;
        background-color: #DCDEE5;
        font-size: 14px;
    }
    .info-title {
        cursor: pointer;
        .icon {
            display: inline-block;
            vertical-align: middle;
            transition: transform .2s linear;
            margin-top: -3px;
            &.is-open {
                transform: rotate(90deg);
            }
        }
        .title-text {
            color: #000;
        }
        .title-count {
            color: #8b8d95;
        }
    }
    .info-pagination {
        color: #8b8d95;
        .pagination-toggle {
            margin-left: 10px;
            .pagination-icon {
                font-size: 14px;
                color: #979BA5;
                cursor: pointer;
                &.disabled {
                    color: #C4C6CC;
                    cursor: not-allowed;
                }
                &.left {
                    transform: rotate(90deg);
                }
                &.right {
                    transform: rotate(-90deg);
                }
            }
        }
    }
    .confirm-tips {
        padding: 9px 0;
        text-align: center;
        .tips-content {
            color: $cmdbTextColor;
            line-height: 20px;
        }
        .tips-option {
            margin: 12px 0 0 0;
            .tips-button {
                height: 26px;
                line-height: 24px;
                padding: 0 16px;
                min-width: 56px;
                font-size: 12px;
            }
        }
    }

    .association-table {
      :deep(.empty-block) {
        width: 100% !important;
      }
    }
</style>
