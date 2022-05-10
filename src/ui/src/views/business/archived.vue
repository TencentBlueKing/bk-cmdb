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
  <div class="archived-layout">
    <div class="archived-filter">
      <div class="filter-item">
        <bk-input v-model="filter.name"
          clearable
          :placeholder="$t('请输入xx', { name: $t('业务') })"
          right-icon="bk-icon icon-search"
          @enter="handlePageChange(1, $event)">
        </bk-input>
      </div>
    </div>
    <bk-table class="archived-table"
      :pagination="pagination"
      :data="list"
      :max-height="$APP.height - 200"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange">
      <bk-table-column prop="bk_biz_id" label="ID"></bk-table-column>
      <bk-table-column v-for="column in header"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        show-overflow-tooltip>
        <template slot-scope="{ row }">{{row[column.id] | formatter(column.property)}}</template>
      </bk-table-column>
      <bk-table-column prop="last_time" :label="$t('更新时间')" show-overflow-tooltip>
        <template slot-scope="{ row }">{{$tools.formatTime(row.last_time)}}</template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" width="200" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth class="inline-block-middle"
            :auth="{ type: $OPERATION.BUSINESS_ARCHIVE, relation: [row.bk_biz_id] }">
            <template #default="{ disabled }">
              <bk-button
                theme="primary"
                size="small"
                :disabled="disabled"
                @click="handleRecovery(row)">
                {{$t('恢复业务')}}
              </bk-button>
              <bk-popconfirm
                :title="$t('确认彻底删除该业务？')"
                trigger="click"
                :content="$t('一旦被清除，将无法再恢复，请谨慎操作！')"
                @confirm="handleCompletelyDelete(row)">
                <bk-button
                  class="ml10"
                  :disabled="disabled"
                  theme="danger"
                  :loading="compeletelyDeleting"
                  size="small"
                >
                  {{$t('彻底删除')}}
                </bk-button>
              </bk-popconfirm>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty slot="empty" :stuff="table.stuff"></cmdb-table-empty>
    </bk-table>

    <bk-dialog
      class="recovery-dialog"
      :draggable="false"
      :mask-close="false"
      :title="$t('恢复业务')"
      header-position="left"
      v-model="recovery.show">
      <div class="recovery-dialog-content">
        <span class="label-title">
          {{$t('业务名')}}
          <font color="red">*</font>
        </span>
        <div class="cmdb-form-item" :class="{ 'is-error': errors.has('bizName') }">
          <cmdb-form-singlechar
            v-model="recovery.name"
            v-validate="'required|singlechar|length:256'"
            name="bizName"
            :placeholder="$t('请输入xx', { name: $t('业务名') })">
          </cmdb-form-singlechar>
          <p class="form-error">{{errors.first('bizName')}}</p>
        </div>
      </div>
      <div class="revocer-foolter" slot="footer">
        <bk-button class="mr10" theme="primary" @click="recoveryBiz">{{$t('确定')}}</bk-button>
        <bk-button @click="recovery.show = false">{{$t('取消')}}</bk-button>
      </div>
    </bk-dialog>
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  export default {
    data() {
      return {
        properties: [],
        header: [],
        list: [],
        filter: {
          range: [],
          name: ''
        },
        pagination: {
          current: 1,
          count: 0,
          ...this.$tools.getDefaultPaginationConfig()
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        },
        recovery: {
          show: false,
          biz: {},
          name: ''
        },
        columnsConfigKey: 'biz_custom_table_columns',
        compeletelyDeleting: false,
      }
    },
    computed: {
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('userCustom', ['usercustom']),
      customBusinessColumns() {
        return this.usercustom[this.columnsConfigKey] || []
      }
    },
    async created() {
      try {
        this.properties = await this.searchObjectAttribute({
          params: {
            bk_obj_id: 'biz',
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: 'post_searchObjectAttribute_biz'
          }
        })
        // 配合全文检索过滤列表
        this.filter.name = this.$route.params.bizName
        this.setTableHeader()
        this.getTableData()
      } catch (e) {
        // ignore
      }
    },
    methods: {
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      ...mapActions('objectBiz', ['searchBusiness', 'recoveryBusiness', 'compeltelyDeleteBusinesses']),
      setTableHeader() {
        const headerProperties = this.$tools.getHeaderProperties(this.properties, this.customBusinessColumns, ['bk_biz_name'])
        this.header = headerProperties.map(property => ({
          id: property.bk_property_id,
          name: this.$tools.getHeaderPropertyName(property),
          property
        }))
      },
      getTableData(event) {
        this.searchBusiness({
          params: this.getSearchParams(),
          config: {
            globalPermission: false,
            cancelPrevious: true,
            requestId: 'searchArchivedBusiness'
          }
        }).then((business) => {
          if (business.count && !business.info.length) {
            this.pagination.current -= 1
            this.getTableData()
          }
          this.pagination.count = business.count
          this.list = business.info

          if (event) {
            this.table.stuff.type = 'search'
          }
        })
          .catch(({ permission }) => {
            if (permission) {
              this.table.stuff = {
                type: 'permission',
                payload: { permission }
              }
            }
          })
      },
      getSearchParams() {
        const params = {
          condition: {
            bk_data_status: 'disabled'
          },
          fields: [],
          page: {
            start: (this.pagination.current - 1) * this.pagination.limit,
            limit: this.pagination.limit,
            sort: '-bk_biz_id'
          }
        }
        if (this.filter.range.length) {
          params.condition.last_time = {
            $gte: this.filter.range[0],
            $lte: this.filter.range[1]
          }
        }
        if (this.filter.name) {
          params.condition.bk_biz_name = this.filter.name
        }
        return params
      },
      handleRecovery(biz) {
        this.recovery.show = true
        this.recovery.name = biz.bk_biz_name
        this.recovery.bizId = biz.bk_biz_id
      },
      handleCompletelyDelete(biz) {
        this.compeletelyDeleting = true
        this.compeltelyDeleteBusinesses({
          bizIds: [biz.bk_biz_id]
        }).then(() => {
          this.$bkMessage({
            message: `${biz.bk_biz_name}已彻底删除`
          })
          this.getTableData()
        })
          .catch((err) => {
            console.log(err)
          })
          .finally(() => {
            this.compeletelyDeleting = false
          })
      },
      async recoveryBiz() {
        if (!await this.$validator.validateAll()) return
        this.recoveryBusiness({
          bizId: this.recovery.bizId,
          params: {
            bk_biz_name: this.recovery.name
          },
          config: {
            cancelWhenRouteChange: false
          }
        }).then(() => {
          this.recovery.show = false
          this.$http.cancel('post_searchBusiness_$ne_disabled')
          this.$success(this.$t('恢复业务成功'))
          this.getTableData()
        })
      },
      handleSizeChange(size) {
        this.pagination.limit = size
        this.handlePageChange(1)
      },
      handlePageChange(current, event) {
        this.pagination.current = current
        this.getTableData(event)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .archived-layout{
        padding: 15px 20px 0;
    }
    .archived-filter {
        padding: 0 0 15px 0;
        .filter-item {
            width: 220px;
            margin-right: 5px;
            @include inlineBlock;
        }
    }
    .recovery-dialog {
        /deep/ .bk-dialog-header {
            padding-bottom: 14px;
        }
        .label-title {
            display: inline-block;
            padding-bottom: 10px;
        }
        .revocer-foolter {
            font-size: 0;
        }
    }
</style>
