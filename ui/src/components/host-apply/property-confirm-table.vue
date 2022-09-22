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
  <div class="property-confirm-table">
    <bk-table
      :data="table.list"
      :pagination="table.pagination"
      :max-height="maxHeight || ($APP.height - 220 - 119)"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange">
      <bk-table-column :label="$t('内网IP')" min-width="150" show-overflow-tooltip>
        <template slot-scope="{ row }">
          <bk-button class="ip-value" theme="primary" text @click="handleShowDetails(row)">
            {{row.expect_host.bk_host_innerip}}
          </bk-button>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('主机名称')" min-width="160" prop="expect_host.bk_host_name">
        <template slot-scope="{ row }">{{row.expect_host.bk_host_name || '--'}}</template>
      </bk-table-column>
      <bk-table-column
        :label="$t('当前值')"
        min-width="500"
        class-name="table-cell-change-value"
        show-overflow-tooltip
        :render-header="(h, data) => renderTableHeader(h, data, $t('主机属性当前值'), { placement: 'right' })">
        <template slot-scope="{ row }">
          <div class="cell-change-value"><vnodes :vnode="getCurrentValue(row)"></vnodes></div>
        </template>
      </bk-table-column>
      <bk-table-column
        :label="$t('目标值')"
        min-width="500"
        class-name="table-cell-change-value"
        show-overflow-tooltip
        :render-header="(h, data) => renderTableHeader(h, data, $t('属性自动应用配置值'), { placement: 'right' })">
        <template slot-scope="{ row }">
          <div class="cell-change-value"><vnodes :vnode="getTargetValue(row)"></vnodes></div>
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty">
        <div>{{$t('暂无主机新转入的主机将会自动应用模块的主机属性')}}</div>
      </cmdb-table-empty>
    </bk-table>
    <bk-sideslider
      v-transfer-dom
      :width="800"
      :is-show.sync="slider.isShow"
      :title="slider.title"
      @hidden="handleSliderCancel">
      <template slot="content">
        <cmdb-details
          :show-options="false"
          :inst="details.inst"
          :properties="details.properties"
          :property-groups="details.propertyGroups">
        </cmdb-details>
      </template>
    </bk-sideslider>
  </div>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  export default {
    components: {
      vnodes: {
        functional: true,
        render: (h, ctx) => ctx.props.vnode
      }
    },
    props: {
      list: {
        type: Array,
        default: () => ([])
      },
      total: {
        type: Number
      },
      maxHeight: {
        type: [Number, String],
        default: 0
      }
    },
    data() {
      return {
        table: {
          list: [],
          pagination: {
            current: 1,
            count: 0,
            ...this.$tools.getDefaultPaginationConfig()
          }
        },
        details: {
          show: false,
          title: '',
          inst: {},
          properties: [],
          propertyGroups: []
        },
        slider: {
          width: 514,
          isShow: false,
          title: ''
        }
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', [
        'getModelById'
      ]),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('hostApply', ['configPropertyList']),
      ...mapState('hostApply', ['propertyList'])
    },
    watch: {
      list() {
        this.setTableList()
      },
      total(value) {
        this.table.pagination.count = value
      }
    },
    async created() {
      await this.getHostPropertyList()
      this.setTableList()
    },
    methods: {
      async getHostPropertyList() {
        try {
          const data = await this.$store.dispatch('hostApply/getProperties', {
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: 'getHostPropertyList',
              fromCache: true
            }
          })

          this.$store.commit('hostApply/setPropertyList', data)
        } catch (e) {
          console.error(e)
        }
      },
      setTableList() {
        const { start, limit } = this.$tools.getPageParams(this.table.pagination)
        this.table.list = this.list.slice(start, start + limit)
      },
      getCurrentValue(row) {
        const { conflicts } = row
        const resultConflicts = conflicts.map((item) => {
          const property = this.configPropertyList.find(propertyItem => propertyItem.id === item.bk_attribute_id) || {}
          return (
            <span>
              {property.bk_property_name}：<cmdb-property-value value={item.bk_property_value} property={property} />
            </span>
          )
        })
        return (
          <div>
            { resultConflicts.reduce((acc, x) => (acc === null ? [x] : [acc, '；', x]), null) }
          </div>
        )
      },
      getTargetValue(row) {
        const { update_fields: updateFields } = row
        const resultUpdates = updateFields.map((item) => {
          const property = this.configPropertyList.find(propertyItem => propertyItem.id === item.bk_attribute_id) || {}
          return (
            <span>
              {property.bk_property_name}：<cmdb-property-value value={item.bk_property_value} property={property} />
            </span>
          )
        })
        return (
          <div>
            { resultUpdates.reduce((acc, x) => (acc === null ? [x] : [acc, '；', x]), null) }
          </div>
        )
      },
      getPropertyGroups() {
        return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
          objId: 'host',
          params: { bk_biz_id: this.bizId }
        })
      },
      renderTableHeader(h, data, tips, options = {}) {
        const directive = {
          content: tips,
          placement: options.placement || 'top-end'
        }
        if (options.width) {
          directive.width = options.width
        }
        return <span>{ data.column.label } <i class="bk-cc-icon icon-cc-tips" v-bk-tooltips={ directive }></i></span>
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.setTableList()
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.table.pagination.current = 1
        this.setTableList()
      },
      async handleShowDetails(row) {
        this.slider.title = `${this.$t('属性详情')}【${row.expect_host.bk_host_innerip}】`
        const properties = this.propertyList
        // 云区域数据
        row.cloud_area.bk_inst_name = row.cloud_area.bk_cloud_name
        row.cloud_area.bk_inst_id = row.cloud_area.bk_cloud_id
        try {
          const [inst, propertyGroups] = await Promise.all([
            this.getHostInfo(row),
            this.getPropertyGroups()
          ])
          inst.bk_cloud_id = [row.cloud_area]
          this.details.inst = inst
          this.details.properties = properties
          this.details.propertyGroups = propertyGroups
          this.slider.isShow = true
        } catch (e) {
          console.log(e)
          this.details.inst = {}
          this.details.properties = []
          this.details.propertyGroups = []
          this.slider.isShow = false
        }
      },
      async getHostInfo(row) {
        try {
          const { info } = this.$store.dispatch('hostSearch/searchHost', {
            params: {
              bk_biz_id: this.bizId,
              condition: ['biz', 'set', 'module'].map(model => ({
                bk_obj_id: model,
                condition: [],
                fields: [`bk_${model}_id`]
              })).concat({
                bk_obj_id: 'host',
                condition: [{ field: 'bk_host_id', operator: '$eq', value: row.expect_host.bk_host_id }],
                fields: []
              }),
              ip: { flag: 'bk_host_innerip', exact: 1, data: [] },
              page: { start: 0, limit: 1 }
            }
          })
          const host = info ? info.host : {}
          return { ...host, ...row.expect_host }
        } catch (error) {
          console.error(error)
          return { ...row.expect_host }
        }
      },
      handleSliderCancel() {
        this.slider.isShow = false
      }
    }
  }
</script>

<style lang="scss" scoped>
    .cell-change-value {
        padding: 8px 0;
        line-height: 1.6;
        word-break: break-word;
    }
    .ip-value {
      white-space: nowrap;
    }
</style>
<style lang="scss">
  .table-cell-change-value {
    .cell {
      -webkit-line-clamp: unset !important;
      display: block !important;

      .icon-cc-tips {
        margin-top: -2px;
      }
    }
    }
</style>
