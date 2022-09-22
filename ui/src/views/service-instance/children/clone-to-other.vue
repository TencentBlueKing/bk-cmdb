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
  <div class="create-layout clearfix" v-bkloading="{ isLoading: $loading() }">
    <label class="create-label fl">{{$t('添加主机')}}</label>
    <div class="create-hosts">
      <bk-button class="select-host-button" theme="default"
        @click="handleAddHost">
        <i class="bk-icon icon-plus"></i>
        {{$t('添加主机')}}
      </bk-button>
      <div class="create-tables" ref="createTables">
        <transition-group name="service-table-list" tag="div">
          <service-instance-table class="service-instance-table"
            v-for="(data, index) in hosts"
            ref="serviceInstanceTable"
            deletable
            :key="data.host.bk_host_id"
            :index="index"
            :id="data.host.bk_host_id"
            :name="getName(data)"
            :source-processes="sourceProcesses"
            :editing="getEditState(data.instance)"
            :instance="data.instance"
            @delete-instance="handleDeleteInstance"
            @edit-name="handleEditName(data.instance)"
            @confirm-edit-name="handleConfirmEditName(data.instance, ...arguments)"
            @cancel-edit-name="handleCancelEditName(data.instance)"
            @change-process="handleChangeProcess">
          </service-instance-table>
        </transition-group>
      </div>
    </div>
    <div class="buttons" :class="{ 'is-sticky': hasScrollbar }">
      <cmdb-auth class="mr5" :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }">
        <div slot-scope="{ disabled }"
          v-bk-tooltips="{
            content: $t('请补充服务实例的进程等相关配置信息'),
            theme: 'light',
            disabled: !hosts.length || hasProcess
          }">
          <bk-button theme="primary"
            :disabled="!hosts.length || disabled || !hasProcess"
            @click="handleConfirm">
            {{$t('确定')}}
          </bk-button>
        </div>
      </cmdb-auth>
      <bk-button @click="handleBackToModule">{{$t('取消')}}</bk-button>
    </div>
    <cmdb-dialog v-model="dialog.show" :width="1110" :height="650" :show-close-icon="false">
      <component
        :is="dialog.component"
        v-bind="dialog.props"
        @confirm="handleDialogConfirm"
        @cancel="handleDialogCancel">
      </component>
    </cmdb-dialog>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import HostSelector from '@/views/business-topology/service-instance/common/host-selector.vue'
  import ServiceInstanceTable from '@/components/service/instance-table.vue'
  import { addResizeListener, removeResizeListener } from '@/utils/resize-events'
  export default {
    name: 'clone-to-other',
    components: {
      HostSelector,
      ServiceInstanceTable
    },
    props: {
      module: {
        type: Object,
        default() {
          return {}
        }
      },
      sourceProcesses: {
        type: Array,
        default() {
          return {}
        }
      }
    },
    data() {
      return {
        dialog: {
          show: false,
          component: null,
          props: {}
        },
        hosts: [],
        hasScrollbar: false,
        hasProcess: false
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      hostId() {
        return parseInt(this.$route.params.hostId, 10)
      },
      moduleId() {
        return parseInt(this.$route.params.moduleId, 10)
      },
      setId() {
        return parseInt(this.$route.params.setId, 10)
      },
      withTemplate() {
        if (this.module.service_template_id) {
          return true
        }
        return false
      }
    },
    mounted() {
      addResizeListener(this.$refs.createTables, this.resizeHandler)
    },
    beforeDestroy() {
      removeResizeListener(this.$refs.createTables, this.resizeHandler)
    },
    methods: {
      handleAddHost() {
        this.dialog.component = HostSelector
        this.dialog.props = {
          exist: this.hosts,
          moduleId: this.moduleId,
          withTemplate: this.withTemplate
        }
        this.dialog.show = true
      },
      handleDialogConfirm(selected) {
        this.hosts = selected.map(item => ({
          ...item,
          instance: {
            name: '',
            editing: { name: false }
          }
        }))
        this.dialog.show = false
      },
      handleDialogCancel() {
        this.dialog.show = false
      },
      handleDeleteInstance(index) {
        this.hosts.splice(index, 1)
      },
      async handleConfirm() {
        try {
          const serviceInstanceTables = this.$refs.serviceInstanceTable

          const params = {
            name: this.module.bk_module_name,
            bk_biz_id: this.bizId,
            bk_module_id: this.moduleId,
            instances: serviceInstanceTables.filter(table => table.processList?.length).map((table) => {
              const { instance } = this.hosts.find(data => data.host.bk_host_id === table.id)
              return {
                bk_host_id: table.id,
                service_instance_name: instance.name || '',
                processes: table.processList.map(item => ({
                  process_info: item
                }))
              }
            })
          }

          await this.$store.dispatch('serviceInstance/createProcServiceInstanceWithRaw', { params })

          this.$success(this.$t('克隆成功'))
          this.handleBackToModule()
        } catch (e) {
          console.error(e)
        }
      },
      getName(data) {
        if (data.instance.name) {
          return data.instance.name
        }
        return data.host.bk_host_innerip || '--'
      },
      getEditState(instance) {
        return instance.editing
      },
      handleEditName(instance) {
        this.hosts.forEach(data => (data.instance.editing.name = false))
        instance.editing.name = true
      },
      handleConfirmEditName(instance, name) {
        instance.name = name
        instance.editing.name = false
      },
      handleCancelEditName(instance) {
        instance.editing.name = false
      },
      handleBackToModule() {
        this.$routerActions.back()
      },
      handleChangeProcess() {
        this.$nextTick(() => {
          const serviceInstanceTables = this.$refs.serviceInstanceTable
          if (serviceInstanceTables) {
            this.hasProcess = serviceInstanceTables.some(instanceTable => instanceTable?.processList?.length)
          }
        })
      },
      resizeHandler() {
        this.$nextTick(() => {
          const scroller = this.$el.parentElement
          this.hasScrollbar = scroller.scrollHeight > scroller.offsetHeight
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .create-layout {
        margin: 35px 0 0 0;
        font-size: 14px;
        color: #63656E;
    }
    .create-label{
        display: block;
        width: 100px;
        position: relative;
        line-height: 32px;
        text-align: right;
        &:after {
            content: "*";
            margin: 0 0 0 4px;
            color: $cmdbDangerColor;
            @include inlineBlock;
        }
    }
    .create-hosts {
        padding-left: 10px;
        padding-right: 20px;
        height: 100%;
        overflow: hidden;
    }
    .select-host-button {
        height: 32px;
        line-height: 30px;
        font-size: 0;
        .bk-icon {
            position: static;
            height: 30px;
            line-height: 30px;
            font-size: 12px;
            font-weight: bold;
            @include inlineBlock(top);
        }
        /deep/ span {
            font-size: 14px;
        }
    }
    .create-tables {
        height: calc(100% - 54px);
        margin: 20px 0 0 0;
        @include scrollbar-y;
        position: relative;
    }
    .buttons {
        position: sticky;
        bottom: 0;
        left: 0;
        padding: 10px 0 10px 110px;

        &.is-sticky {
            background-color: #FFF;
            border-top: 1px solid $borderColor;
            z-index: 100;
        }
    }
    .service-instance-table {
        margin-bottom: 12px;
    }

    .service-table-list-enter-active, .service-table-list-leave-active {
        transition: all .7s ease-in;
    }
    .service-table-list-leave-to {
        opacity: 0;
        transform: translateX(30px);
    }
</style>
