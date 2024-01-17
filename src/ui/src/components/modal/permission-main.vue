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
  <div class="permission-main">
    <div class="permission-content">
      <div class="permission-header">
        <bk-exception type="403" scene="part">
          <h3>{{i18n.permissionTitle}}</h3>
        </bk-exception>
      </div>
      <bk-table ref="table"
        :data="list"
        :max-height="193"
        class="permission-table">
        <bk-table-column :label="$t('系统')" width="150">
          {{ $t(permission.system_name) }}
        </bk-table-column>
        <bk-table-column prop="name" :label="$t('需要申请的权限')" width="200"></bk-table-column>
        <bk-table-column prop="resource" :label="$t('关联的资源实例')">
          <template slot-scope="{ row }">
            <div v-if="row.relations.length" style="overflow: auto;">
              <div class="permission-resource"
                v-for="(relation, index) in row.relations"
                v-bk-overflow-tips
                :key="index">
                <permission-resource-name :id="relation.id" :relations="relation" />
              </div>
            </div>
            <span v-else>--</span>
          </template>
        </bk-table-column>
      </bk-table>
      <div class="relate-apply-permission" v-if="relatedList.length">
        <i18n class="caption" tag="div" path="你可以勾选页面内其他待申请的n个资源，一并申请">
          <template #n><span>{{ relatedList.length }}</span></template>
        </i18n>
        <bk-table ref="relateTable"
          :data="relatedList"
          :max-height="260"
          class="permission-table"
          @selection-change="handleRelateSelectChange">
          <bk-table-column width="46" type="selection"></bk-table-column>
          <bk-table-column prop="name" :label="$t('需要申请的权限')" width="204"
            class-name="related-name-col"></bk-table-column>
          <bk-table-column prop="relation" :label="$t('关联的资源实例')">
            <template slot-scope="{ row }">
              <div class="permission-resource" v-bk-overflow-tips>
                <permission-resource-name :relations="row.relation" />
              </div>
            </template>
          </bk-table-column>
        </bk-table>
      </div>
    </div>
    <div class="permission-footer">
      <template v-if="applied && !selectChanged">
        <bk-button theme="primary" @click="handleRefresh">{{ i18n.applied }}</bk-button>
        <bk-button class="ml10" @click="handleClose">{{ i18n.close }}</bk-button>
      </template>
      <template v-else>
        <bk-button theme="primary"
          :loading="$loading('getSkipUrl')"
          @click="handleApply">
          {{ i18n.apply }}
        </bk-button>
        <bk-button class="ml10" @click="handleClose">{{ i18n.cancel }}</bk-button>
      </template>
    </div>
  </div>
</template>
<script>
  import cloneDeep from 'lodash/cloneDeep'
  import { IAM_ACTIONS, IAM_VIEWS_NAME, IAM_VIEWS } from '@/dictionary/iam-auth'
  import { mergeSameActions } from '@/setup/permission'
  import PermissionResourceName from './permission-resource-name.vue'
  export default {
    components: {
      PermissionResourceName
    },
    props: {
      permission: Object,
      relatedPermission: Object,
      applied: Boolean
    },
    data() {
      return {
        list: [],
        relatedList: [],
        selectedRelatedPermission: {},
        selectChanged: false,
        i18n: {
          permissionTitle: this.$t('没有权限访问或操作此资源'),
          system: this.$t('系统'),
          resource: this.$t('资源'),
          requiredPermissions: this.$t('需要申请的权限'),
          noData: this.$t('无数据'),
          apply: this.$t('去申请'),
          applied: this.$t('我已申请'),
          cancel: this.$t('取消'),
          close: this.$t('关闭')
        }
      }
    },
    watch: {
      permission() {
        this.setList()
      },
      relatedPermission() {
        this.setRelatedList()
      }
    },
    created() {
      this.setList()
      this.setRelatedList()
    },
    methods: {
      generateList(permission) {
        const languageIndex = this.$i18n.locale === 'en' ? 1 : 0
        return permission?.actions?.map((action) => {
          const { id: actionId, related_resource_types: relatedResourceTypes = [] } = action
          const definition = Object.values(IAM_ACTIONS).find((definition) => {
            if (typeof definition.id === 'function') {
              return actionId.indexOf(definition.fixedId) > -1
            }
            return definition.id === actionId
          })
          const allRelationPath = []
          relatedResourceTypes.forEach(({ instances = [] }) => {
            instances.forEach((fullPaths) => {
              // 数据格式[type, id, label]
              const topoPath = []
              fullPaths.forEach((pathData) => {
                const isComobj = pathData.type.indexOf('comobj_') > -1
                if (isComobj) {
                  const [, modelId] = pathData.type.split('_')
                  const instId = pathData.id
                  const modelView = IAM_VIEWS.INSTANCE_MODEL
                  const instView = IAM_VIEWS.INSTANCE
                  const modelPath = [modelView, modelId, IAM_VIEWS_NAME[modelView][languageIndex]]
                  const instPath = instId ? [instView, pathData.id, IAM_VIEWS_NAME[instView][languageIndex]] : null
                  topoPath.push(modelPath)
                  instPath && topoPath.push(instPath)
                } else {
                  topoPath.push([pathData.type, pathData.id, IAM_VIEWS_NAME[pathData.type][languageIndex]])
                }
              })
              allRelationPath.push(topoPath)
            })
          })
          if (!allRelationPath.length && actionId.indexOf('comobj') > -1) {
            // 兼容创建模型实例没有relatedResourceTypes的情况
            const [,, modelId] = actionId.split('_')
            const modelView = IAM_VIEWS.INSTANCE_MODEL
            allRelationPath.push([[modelView, modelId, IAM_VIEWS_NAME[modelView][languageIndex]]])
          }
          return {
            id: actionId,
            name: definition.name[languageIndex],
            relations: allRelationPath
          }
        })
      },
      setList() {
        this.list = this.generateList(this.permission)
      },
      setRelatedList() {
        const list = this.generateList(this.relatedPermission)
        const relatedList = []
        list?.forEach((item, actionIndex) => {
          const { id, name, relations } = item
          relations.forEach((relation) =>  {
            // 单个实例可能是多层的，如果 xx业务 / xx动作，这里把所有id和type拼接作为唯一key
            const key = relation?.reduce((acc, cur) => `${acc}/${cur[0]}_${cur[1]}`, id)
            relatedList.push({ id, name, relation, actionIndex, key })
          })
        })
        this.relatedList = relatedList
      },
      handleRelateSelectChange(selection) {
        this.selectedRelatedPermission.system_id = this.relatedPermission.system_id
        this.selectedRelatedPermission.actions = []
        // 根据选择的权限反向查找到对应的permission数据，并过滤掉未选择的实例
        selection.forEach((item) => {
          const { id, actionIndex, key } = item
          const actionItem = cloneDeep(this.relatedPermission.actions?.[actionIndex])
          actionItem.related_resource_types.forEach((resourceItem) => {
            resourceItem.instances = resourceItem.instances.filter((inst) => {
              const keyMatched = inst?.reduce((acc, cur) => `${acc}/${cur.type}_${cur.id}`, id)
              return key === keyMatched
            })
          })
          this.selectedRelatedPermission.actions.push(actionItem)
        })
        this.selectChanged = true
      },
      handleClose() {
        this.$emit('close')
      },
      handleApply() {
        let permission
        if (this.selectedRelatedPermission?.actions?.length) {
          const actions = this.permission.actions.concat(this.selectedRelatedPermission.actions)
          permission = mergeSameActions(actions)
        }
        this.$emit('apply', permission)
        this.selectChanged = false
      },
      handleRefresh() {
        this.$emit('refresh')
      },
      doTableLayout() {
        this.$refs.table.doLayout()
        this.$refs?.relateTable?.doLayout()
      }
    }
  }
</script>
<style lang="scss" scoped>
.permission-content {
  margin-top: -26px;
  padding: 3px 24px 26px;
  .permission-header {
    padding-top: 16px;
    text-align: center;
    .locked-icon {
      height: 66px;
    }
    h3 {
      margin: 6px 0 30px;
      color: #63656e;
      font-size: 24px;
      font-weight: normal;
    }

    /deep/ {
      .bk-exception-img .exception-image {
          height: 130px;
      }
    }
  }
}
.permission-footer {
  text-align: right;
  padding: 12px 24px;
  background-color: #fafbfd;
  border-top: 1px solid #dcdee5;
  border-radius: 2px;
}
.permission-table {
  .permission-resource {
    line-height: 24px;
  }
  /deep/ {
    .bk-table-row {
      td.is-first {
        vertical-align: top;
        line-height: 42px;
      }
    }
  }

  :deep(.related-name-col) {
    .cell {
      padding: 0;
    }
  }
}
.relate-apply-permission {
  margin-top: 12px;
  .caption {
    font-size: 12px;
    font-weight: 700;
    padding: 8px 0;
  }
}
</style>
