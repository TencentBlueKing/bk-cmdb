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
  <div class="template-layout">
    <cmdb-tips class="mb10" tips-key="showSetTips">
      <i18n path="集群模板功能提示" class="tips-text">
        <template #topo>
          <a class="tips-link" href="javascript:void(0)" @click="handleGoBusinessTopo">
            {{$t('业务拓扑')}}
          </a></template>
        <template #template>
          <a class="tips-link" href="javascript:void(0)" @click="handleGoServiceTemplate">
            {{$t('服务模板')}}
          </a>
        </template>
      </i18n>
    </cmdb-tips>
    <div class="options clearfix">
      <div class="fl">
        <cmdb-auth class="fl" :auth="{ type: $OPERATION.C_SET_TEMPLATE, relation: [bizId] }">
          <bk-button slot-scope="{ disabled }" v-test-id="'create'"
            theme="primary"
            :disabled="disabled"
            @click="handleCreate">
            {{$t('新建')}}
          </bk-button>
        </cmdb-auth>
      </div>
      <div class="fr">
        <bk-input :style="{ width: '210px' }"
          :placeholder="$t('请输入xx', { name: $t('模板名称') })"
          clearable
          right-icon="icon-search"
          v-model.trim="searchName"
          @enter="handleFilterTemplate"
          @clear="handleClearFilter">
        </bk-input>
      </div>
    </div>
    <bk-table class="template-table" v-test-id.businessSetTemplate="'templateList'"
      v-bkloading="{ isLoading: $loading(request.getSetTemplates) }"
      :data="list"
      :row-style="{ cursor: 'pointer' }"
      @row-click="handleRowClick"
      @sort-change="handleSortChange">
      <bk-table-column label="ID" prop="id" class-name="is-highlight" sortable="custom">
        <div slot-scope="{ row }" v-bk-overflow-tips
          :class="['template-name', { 'need-sync': row._need_sync_ }]">
          {{row.id}}
        </div>
      </bk-table-column>
      <bk-table-column :label="$t('模板名称')" prop="name" sortable="custom"></bk-table-column>
      <bk-table-column :label="$t('应用数量')" prop="set_instance_count" sortable></bk-table-column>
      <bk-table-column :label="$t('修改人')" prop="modifier" sortable="custom"></bk-table-column>
      <bk-table-column :label="$t('修改时间')" prop="last_time" sortable="custom" show-overflow-tooltip>
        <template slot-scope="{ row }">
          <span>{{$tools.formatTime(row.last_time, 'YYYY-MM-DD HH:mm')}}</span>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" width="180">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_SET_TEMPLATE, relation: [bizId, row.id] }">
            <bk-button slot-scope="{ disabled }"
              text
              :disabled="disabled"
              @click="handleEdit(row)">
              {{$t('编辑')}}
            </bk-button>
          </cmdb-auth>
          <span class="text-primary"
            style="color: #dcdee5 !important; cursor: not-allowed;"
            v-if="row.set_instance_count"
            v-bk-tooltips.top="$t('不可删除')">
            {{$t('删除')}}
          </span>
          <cmdb-auth :auth="{ type: $OPERATION.D_SET_TEMPLATE, relation: [bizId, row.id] }" v-else>
            <bk-button slot-scope="{ disabled }" v-test-id="'delTemplate'"
              text
              :disabled="disabled"
              @click="handleDelete(row)">
              {{$t('删除')}}
            </bk-button>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        :auth="{ type: $OPERATION.C_SET_TEMPLATE, relation: [bizId] }"
        @create="handleCreate"
      ></cmdb-table-empty>
    </bk-table>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE,
    MENU_BUSINESS_SET_TEMPLATE_CREATE,
    MENU_BUSINESS_SET_TEMPLATE_DETAILS,
    MENU_BUSINESS_SET_TEMPLATE_EDIT
  } from '@/dictionary/menu-symbol'

  export default {
    data() {
      return {
        list: [],
        originList: [],
        searchName: '',
        table: {
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('集群模板')
            }
          },
          sort: '-last_time'
        },
        request: {
          getSetTemplates: Symbol('getSetTemplates')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId'])
    },
    watch: {
      originList() {
        this.getSyncStatus()
      }
    },
    async created() {
      await this.getSetTemplates()
    },
    methods: {
      async getSetTemplates() {
        const data = await this.$store.dispatch('setTemplate/getSetTemplates', {
          bizId: this.bizId,
          params: {
            page: {
              sort: this.table.sort
            }
          },
          config: {
            requestId: this.request.getSetTemplates
          }
        })
        const list = (data.info || []).map(item => ({
          set_instance_count: item.set_instance_count,
          ...item.set_template,
          _need_sync_: false
        }))
        this.list = list
        this.originList = list
      },
      async getSyncStatus() {
        try {
          if (this.originList.length) {
            const data = await this.$store.dispatch('setTemplate/getSetTemplateStatus', {
              bizId: this.bizId,
              params: {
                set_template_ids: this.originList.map(item => item.id)
              }
            })
            this.originList.forEach((item) => {
              const syncStatus = data.find(status => status.set_template_id === item.id)
              if (syncStatus) {
                this.$set(item, '_need_sync_', syncStatus.need_sync)
              }
            })
          }
        } catch (e) {
          console.error(e)
        }
      },
      handleCreate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_CREATE,
          history: true
        })
      },
      async handleDelete(row) {
        this.$bkInfo({
          title: this.$t('确认删除'),
          subTitle: this.$t('确认删除xx集群模板', { name: row.name }),
          confirmFn: async () => {
            try {
              await this.$store.dispatch('setTemplate/deleteSetTemplate', {
                bizId: this.$store.getters['objectBiz/bizId'],
                config: {
                  data: {
                    set_template_ids: [row.id]
                  }
                }
              })
              this.getSetTemplates()
            } catch (e) {
              console.error(e)
            }
          }
        })
      },
      handleFilterTemplate() {
        const originList = this.$tools.clone(this.originList)
        this.list = this.searchName
          ? originList.filter(template => template.name.indexOf(this.searchName) !== -1)
          : originList
        this.table.stuff.type = this.searchName ? 'search' : 'default'
      },
      handleClearFilter() {
        this.list = this.originList
        this.table.stuff.type = 'default'
      },
      handleSelectable(row) {
        return !row.set_instance_count
      },
      handleEdit(row) {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_EDIT,
          params: {
            templateId: row.id
          },
          history: true
        })
      },
      handleRowClick(row, event, column) {
        if (!column.property) {
          return false
        }
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SET_TEMPLATE_DETAILS,
          params: {
            templateId: row.id
          },
          history: true
        })
      },
      handleSortChange(sort) {
        if (sort.prop === 'set_instance_count') {
          return
        }
        this.table.sort = this.$tools.getSort(sort, '-last_time')
        this.getSetTemplates()
      },
      handleGoBusinessTopo() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_AND_SERVICE
        })
      },
      handleGoServiceTemplate() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .template-layout {
        padding: 15px 20px 0;
    }
    .tips-link {
        color: #3a84ff;
        &:hover {
            text-decoration: underline;
        }
    }
    .options {
        font-size: 0;
    }
    .template-table {
        margin-top: 16px;
        .template-name {
            position: relative;
            display: inline-block;
            vertical-align: middle;
            line-height: 40px;
            padding-right: 10px;
            max-width: 100%;
            @include ellipsis;
            &.need-sync:after {
                content: "";
                position: absolute;
                top: 10px;
                right: 2px;
                width: 6px;
                height: 6px;
                border-radius: 50%;
                background-color: $dangerColor;
            }
        }
    }
</style>
