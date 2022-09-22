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
        <bk-table-column prop="name" :label="$t('需要申请的权限')" width="250"></bk-table-column>
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
    </div>
    <div class="permission-footer">
      <template v-if="applied">
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
  import { IAM_ACTIONS, IAM_VIEWS_NAME, IAM_VIEWS } from '@/dictionary/iam-auth'
  import PermissionResourceName from './permission-resource-name.vue'
  export default {
    components: {
      PermissionResourceName
    },
    props: {
      permission: Object,
      applied: Boolean
    },
    data() {
      return {
        list: [],
        i18n: {
          permissionTitle: this.$t('没有权限访问或操作此资源'),
          system: this.$t('系统'),
          resource: this.$t('资源'),
          requiredPermissions: this.$t('需要申请的权限'),
          noData: this.$t('无数据'),
          apply: this.$t('去申请'),
          applied: this.$t('已完成'),
          cancel: this.$t('取消'),
          close: this.$t('关闭')
        }
      }
    },
    watch: {
      permission() {
        this.setList()
      }
    },
    created() {
      this.setList()
    },
    methods: {
      setList() {
        const languageIndex = this.$i18n.locale === 'en' ? 1 : 0
        this.list = this.permission.actions.map((action) => {
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
      handleClose() {
        this.$emit('close')
      },
      handleApply() {
        this.$emit('apply')
      },
      handleRefresh() {
        this.$emit('refresh')
      },
      doTableLayout() {
        this.$refs.table.doLayout()
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
    }
</style>
