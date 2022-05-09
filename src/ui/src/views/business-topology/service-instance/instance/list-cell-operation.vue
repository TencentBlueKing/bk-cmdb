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
  <div class="instance-operation">
    <cmdb-auth tag="span" class="operation-item" v-test-id.businessHostAndService="'addProcess'"
      v-if="!row.service_template_id"
      :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
      @click.native.stop
      @click="handleAddProcess">
      {{$t('添加进程')}}
    </cmdb-auth>
    <cmdb-auth tag="span" class="operation-item" v-test-id.businessHostAndService="'cloneProcess'"
      v-if="!row.service_template_id"
      :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }"
      @click.native.stop
      @click="handleClone">
      {{$t('克隆')}}
    </cmdb-auth>
    <cmdb-auth tag="span" class="operation-item" v-test-id.businessHostAndService="'delProcess'"
      :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }"
      @click.native.stop>
      <bk-popconfirm trigger="click"
        ext-popover-cls="del-confirm"
        :content="$t('确定删除该服务实例')"
        confirm-loading
        @confirm="handleDelete">
        {{$t('删除')}}
      </bk-popconfirm>
    </cmdb-auth>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import Bus from '../common/bus'
  import createProcessMixin from './create-process-mixin'
  export default {
    mixins: [createProcessMixin],
    props: {
      row: Object
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('businessHost', ['selectedNode'])
    },
    methods: {
      handleClone() {
        this.$routerActions.redirect({
          name: 'cloneServiceInstance',
          params: {
            instanceId: this.row.id,
            hostId: this.row.bk_host_id,
            setId: this.selectedNode.parent.data.bk_inst_id,
            moduleId: this.selectedNode.data.bk_inst_id
          },
          query: {
            title: this.row.name,
            node: this.selectedNode.id
          },
          history: true
        })
      },
      async handleDelete() {
        try {
          await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
            config: {
              data: {
                service_instance_ids: [this.row.id],
                bk_biz_id: this.bizId
              }
            }
          })
          this.$success(this.$t('删除成功'))
          Bus.$emit('delete-complete')
          return true
        } catch (e) {
          console.error(e)
          return false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .instance-operation {
        display: flex;
        align-items: center;
    }
    .operation-item {
        display: inline-block;
        line-height: 32px;
        color: $textColor;
        font-size: 12px;
        cursor: pointer;
        color: $primaryColor;
        &:hover {
            opacity: .7;
        }
        &.disabled {
            color: $textDisabledColor;
            opacity: 1;
        }
        & ~ .operation-item {
            margin-left: 10px;
        }
    }
</style>
