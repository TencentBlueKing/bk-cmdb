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
  <bk-dialog
    ext-cls="permission-dialog"
    v-model="isModalShow"
    width="740"
    :z-index="2400"
    :close-icon="false"
    :mask-close="true"
    :show-footer="false"
    :position="dialogPos"
    @cancel="onCloseDialog">
    <permission-main ref="main" :permission="permission" :related-permission="relatedPermission" :applied="applied"
      @close="onCloseDialog"
      @apply="handleApply"
      @refresh="handleRefresh" />
  </bk-dialog>
</template>
<script>
  import { IAM_VIEWS } from '@/dictionary/iam-auth'
  import permissionMixins from '@/mixins/permission'
  import PermissionMain from './permission-main.vue'
  export default {
    name: 'permissionModal',
    components: {
      PermissionMain
    },
    mixins: [permissionMixins],
    props: {},
    data() {
      return {
        applied: false,
        isModalShow: false,
        permission: {
          actions: []
        },
        callbackUrl: '',
        relatedPermission: null,
        dialogPos: {}
      }
    },
    watch: {
      isModalShow(val) {
        if (val) {
          setTimeout(() => {
            this.$refs.main.doTableLayout()
            this.setDialogPos()
          }, 0)
        }
      },
      '$APP.height'() {
        this.setDialogPos()
      }
    },
    methods: {
      show(permission, authResults, callbackUrl = '', relatedPermission = null) {
        this.permission = this.getPermission(permission, authResults)
        this.applied = false
        this.callbackUrl = callbackUrl
        this.relatedPermission = relatedPermission
        this.isModalShow = true
      },
      onCloseDialog() {
        this.isModalShow = false
      },
      async handleApply(permission) {
        const finalPermission = permission || this.permission
        try {
          await this.handleApplyPermission(finalPermission)
          this.applied = true
        } catch (error) {}
      },
      handleRefresh() {
        if (this.callbackUrl) {
          window.location.replace(this.callbackUrl)
        } else {
          window.location.reload()
        }
      },
      getPermission(permission, authResults) {
        if (!authResults) {
          return permission
        }

        // 批量鉴权的场景下从permission中过滤掉有权限实例
        const batchInstTypes = [IAM_VIEWS.INSTANCE, IAM_VIEWS.HOST, IAM_VIEWS.PROJECT] // 通用模型实例和主机
        permission.actions.forEach((action) => {
          const { related_resource_types: relatedResourceTypes = [] } = action
          const newInstances = []
          relatedResourceTypes.forEach(({ instances = [] }) => {
            const insts = instances.filter((fullPaths) => {
              const matched = fullPaths.find(item => batchInstTypes.includes(item.type))
              if (matched) {
                const authed = authResults.find(item => String(item.resource_id) === matched.id)
                return !authed?.is_pass
              }
              return true
            })
            newInstances.push(insts)
          })

          // 替换更新整个instances
          relatedResourceTypes.forEach((item, i) => item.instances = newInstances[i])
        })

        return permission
      },
      setDialogPos() {
        const dialogHeight = document.querySelector('.permission-dialog .bk-dialog .bk-dialog-content')?.getBoundingClientRect()?.height
        if (dialogHeight) {
          this.dialogPos = {
            top: `${Math.floor(Math.max(this.$APP.height * 0.4 - (parseInt(dialogHeight, 10) / 2), 20))}`
          }
        }
      }
    }
  }
</script>
<style lang="scss" scoped>
    /deep/ .permission-dialog {
        .bk-dialog-body {
            padding: 0;
        }
    }
</style>
