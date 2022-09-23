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
  <div class="process-difference cleafix">
    <!-- 进程模板列表 -->
    <ul class="process-list fl">
      <li class="process-item"
        v-for="(process, index) in processDiff"
        :key="index"
        :class="{
          'show-tips': !process.confirmed,
          'is-active': currentDiff.process_template_id === process.process_template_id
            && process.process_template_name === currentDiff.process_template_name,
          'is-remove': process.type === 'removed'
        }"
        @click="loadProcessDiff(process, index)">
        <span class="process-name" :title="process.process_template_name">{{process.process_template_name}}</span>
      </li>
    </ul>

    <!-- 变更详情 -->
    <div class="change-details" ref="changeDetails"
      v-bkloading="{ isLoading: currentDiffLoading }"
      v-if="currentDiff"
      :key="`${currentDiff.process_template_id}-${currentDiff.process_template_name}`">
      <cmdb-collapse class="details-info" arrow-type="filled">
        <div class="collapse-title" slot="title">
          {{$t('变更内容')}}
          <span v-if="currentDiff.type === 'changed'">({{currentDiff.changedProperties.length}})</span>
        </div>
        <div class="info-content">

          <!-- 进程新增 -->
          <div class="process-info"
            v-if="currentDiff.type === 'added'">
            <div class="info-item" style="width: auto;">
              {{$t('模板中新增进程')}}
              <span class="info-item-value">{{currentDiff.process_template_name}}</span>
            </div>
          </div>

          <!-- 进程删除 -->
          <div class="process-info"
            v-if="currentDiff.type === 'removed'">
            <div class="info-item" style="width: auto;">
              <span class="info-item-value" style="font-weight: 700;">{{currentDiff.process_template_name}}</span>
              {{$t('从模板中删除')}}
            </div>
          </div>

          <!-- 进程变更 -->
          <div
            class="process-info clearfix"
            v-if="currentDiff.type === 'changed'">
            <div :class="['info-item fl', { table: changed.property.bk_property_type === 'table' }]"
              v-for="(changed, index) in currentDiff.changedProperties"
              :key="index"
              v-bk-overflow-tips>
              {{changed.property.bk_property_name}}：
              <span class="info-item-value">
                <span
                  v-if="changed.property.bk_property_id === 'bind_info' && !changed.template_property_value.length">
                  {{$t('移除所有进程监听信息')}}
                </span>
                <cmdb-property-value v-else
                  :value="formatChangedValue(changed)"
                  :property="changed.property">
                </cmdb-property-value>
              </span>
            </div>
          </div>
        </div>
      </cmdb-collapse>
      <cmdb-collapse
        class="details-modules"
        arrow-type="filled"
        :key="moduleId"
        :collapse="true"
        @collapse-change="handleInstanceCollapseChange(moduleId, $event)">
        <div class="collapse-title" slot="title">
          {{$t('涉及实例')}}
          <span v-if="currentDiff.module.serviceInstanceCount !== ''">
            ({{currentDiff.module.serviceInstanceCount}})
          </span>
        </div>

        <!-- 实例列表 -->
        <ul
          class="instance-list"
          ref="instanceList"
          v-bkloading="{ isLoading: currentDiff.module.instancesLoading }">
          <li class="instance-item"
            v-for="(instance, instanceIndex) in currentDiff.module.serviceInstances"
            :key="instanceIndex"
            @click="viewInstanceDiff(instance, moduleId)">
            <div class="instance-diff">
              <div class="instance-name" v-bk-overflow-tips>
                {{instance.name}}
              </div>
              <label
                :class="['instance-change-type', instance.type]"
                v-if="translateChangedType(instance.type)">
                {{translateChangedType(instance.type)}}
              </label>
            </div>
          </li>
        </ul>
      </cmdb-collapse>
    </div>

    <!-- 实例对比详情 -->
    <bk-sideslider
      v-transfer-dom
      :width="676"
      :is-show.sync="instanceDiffSlider.show"
      :title="instanceDiffSlider.title">
      <template slot="content" v-if="instanceDiffSlider.show">
        <ServiceInstanceDetails
          slot="content"
          v-if="instanceDiffSlider.show"
          v-bind="instanceDiffSlider.props"
          :properties="properties">
        </ServiceInstanceDetails>
      </template>
    </bk-sideslider>
  </div>
</template>

<script>
  import ServiceInstanceDetails from './service-instance-details.vue'
  import formatter from '@/filters/formatter'
  import { mapGetters } from 'vuex'
  import isEmpty from 'lodash/isEmpty'
  import throttle from 'lodash/throttle'

  export default {
    components: {
      ServiceInstanceDetails
    },
    props: {
      moduleId: {
        type: Number,
        required: true
      },
      templateId: {
        type: Number,
        required: true
      },
      topoPath: {
        type: String,
        default: '',
        required: false
      },
      processDiff: {
        type: Array,
        default: () => ([]),
        required: true
      },
      properties: {
        type: Array,
        default: () => ([]),
        required: true
      }
    },
    data() {
      return {
        currentDiff: {
          process_template_id: 0, // 当前进程模板 id
          process_template_name: '', // 当前进程模板名称
          process_template: {}, // 当前进程模板内容
          changedProperties: [], // 当前进程模板实例具体更改细节
          module: {
            serviceInstanceCount: '',
            serviceInstances: [],
            instancesLoading: false
          }
        },
        currentDiffLoading: false,
        confirming: false,
        serviceCategories: [],
        instanceDiffSlider: {
          show: false,
          title: '',
          props: {
            module: null,
            instance: null,
            process: null
          }
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId'])
    },
    created() {
      if (this.processDiff?.length) {
        this.loadProcessDiff(this.processDiff[0])
      }

      this.throttleHideTooltips = throttle(this.hideTooltips, 100)
    },
    mounted() {
      setTimeout(() => {
        this.$refs.changeDetails.addEventListener('scroll', this.throttleHideTooltips)
      }, 1000)
    },
    beforeDestroy() {
      this.$refs.changeDetails.removeEventListener('scroll', this.throttleHideTooltips)
    },
    methods: {
      /**
       * 加载进程模板变更内容
       * @param {Object} process 进程信息
       */
      async loadProcessDiff(process) {
        if (process.type === 'removed') {
          return this.loadRemovedProcess(process)
        }

        this.loadChangedProcess(process)
      },
      /**
       * 渲染被删除的进程模板的变更内容
       */
      loadRemovedProcess(process) {
        this.currentDiff.type = process.type
        this.currentDiff.process_template_name = process.process_template_name
        this.currentDiff.process_template_id = process.process_template_id
        process.confirmed = true
      },
      /**
       * 加载进程
       * @param {Object} process 进程信息
       */
      loadChangedProcess(process) {
        this.currentDiffLoading = true
        this.$store.dispatch('processTemplate/getProcessTemplate', {
          params: {
            processTemplateId: process.process_template_id
          }
        }).then((res) => {
          this.currentDiff.type = process.type
          this.currentDiff.process_template_id = process.process_template_id
          this.currentDiff.process_template_name = process.process_template_name
          this.currentDiff.changedProperties = this.getChangedProperties(res.property)
          process.confirmed = true
        })
          .finally(() => {
            this.currentDiffLoading = false
          })
      },
      getChangedProperties(property) {
        const changedProperties = []

        if (property) {
          Object.keys(property).forEach((key) => {
            const prop = property[key]
            const formatedProp = this.properties.find(i => i.bk_property_id === key)

            if (!isEmpty(prop.value)) {
              changedProperties.push({
                property: formatedProp,
                template_property_value: prop.value
              })
            }
          })
        }

        return changedProperties
      },
      formatChangedValue(changed) {
        const { property } = changed
        const { template_property_value: value } = changed
        return formatter(value?.value || value, property)
      },
      handleInstanceCollapseChange(moduleId, collapse) {
        if (!collapse) {
          this.loadInstances(moduleId)
        }
      },
      /**
       * 加载服务模板涉及的实例
       */
      loadInstances(moduleId) {
        const theModule = this.currentDiff.module
        theModule.instancesLoading = true
        this.$store.dispatch('businessSynchronous/getDiffInstances', { params: {
          ...this.serializeParams(),
          bk_module_id: Number(moduleId),
        } })
          .then(({ service_instances: serviceInstances, total_count: totalCount, type }) => {
            let instancesDiff = []

            instancesDiff = serviceInstances?.map(instance => ({
              ...instance,
              type: type || this.currentDiff.type
            }))
            theModule.serviceInstanceCount = totalCount
            theModule.serviceInstances = instancesDiff || []
          })
          .catch(() => {
            theModule.serviceInstanceCount = ''
            theModule.serviceInstances = []
          })
          .finally(() => {
            theModule.instancesLoading = false
          })
      },
      serializeParams() {
        const params = {
          bk_biz_id: this.bizId,
          service_template_id: this.templateId,
        }

        if (this.currentDiff.type === 'removed') {
          // 因为删除以后模板 id 会变成 0，所以需要模板名称来加载对应的实例
          params.process_template_name = this.currentDiff.process_template_name
          params.process_template_id = this.currentDiff.process_template_id
        } else {
          params.process_template_id = this.currentDiff.process_template_id
        }

        return params
      },
      translateChangedType(type) {
        const types = new Map()
        types.set('added', '新增')
        types.set('removed', '删除')
        types.set('changed', '变更')
        types.set('others', '变更')

        return types.get(type)
      },
      /**
       * 查看实例对比
       * @param {Object} instance 实例对象
       * @param {String} moduleId 模板 ID
       */
      viewInstanceDiff(instance, moduleId) {
        this.instanceDiffSlider.title = instance.name
        this.instanceDiffSlider.props = {
          diffRequestParams: {
            ...this.serializeParams(),
            bk_module_id: Number(moduleId),
            service_instance_id: instance.id
          },
          properties: this.properties,
          getCategoryById: this.getCategoryById
        }
        this.instanceDiffSlider.show = true
      },
      hideTooltips() {
        const instanceNameDoms = this.$refs?.instanceList?.getElementsByClassName('instance-name')
        // eslint-disable-next-line no-underscore-dangle
        const tippyInsts = Array.from(instanceNameDoms, el => el?._tippy)
        tippyInsts.forEach(inst => inst?.hide())
      }
    }
  }
</script>

<style lang="scss" scoped>
.process-difference {
  margin-top: 10px;
  border: 1px solid $borderColor;
  border-bottom: none;
  height: calc(100vh - 450px);
  overflow: hidden;
  .process-list {
    position: relative;
    margin-right: -1px;
    width: 200px;
    height: 100%;
    z-index: 2;
    @include scrollbar-y;
  }
  .change-details {
    position: relative;
    height: 100%;
    padding: 20px;
    background-color: #FFF;
    border-left: 1px solid $borderColor;
    border-bottom: 1px solid $borderColor;
    z-index: 1;
    @include scrollbar-y;
  }
}
.process-list {
  border-bottom: 1px solid $borderColor;
  .process-item {
    display: flex;
    padding: 0 12px 0 14px;
    height: 42px;
    align-items: center;
    justify-content: space-between;
    background-color: #FAFBFD;
    border-right: 1px solid $borderColor;
    border-bottom: 1px solid $borderColor;
    cursor: pointer;
    &.is-active {
      background-color: #FFF;
      border-right: none;
      .process-name {
        font-weight: bold;
        color: #2DCB56;
      }
      &.is-remove {
        .process-name {
          color: $dangerColor;
        }
      }
    }
    &.is-remove {
      .process-name {
        text-decoration: line-through;
      }
    }
    &.show-tips {
      .process-name:after {
        position: absolute;
        width: 6px;
        height: 6px;
        top: 21px;
        right: 4px;
        border-radius: 50%;
        background-color: #FF5656;
        content: "";
        z-index: 1;
      }
    }
    .process-name {
      font-size: 12px;
      position: relative;
      padding: 0 14px 0 0;
      @include ellipsis;
    }
    .process-service-count {
      padding: 0 8px;
      height: 16px;
      line-height: 16px;
      font-size: 12px;
      font-style: normal;
      text-align: center;
      background-color: #c4c6cc;
      color: #fff;
      border-radius: 8px;
    }
  }
}
.details-info {
  .process-info {
    padding: 0 0 0 22px;
    .info-item {
      width: 200px;
      font-size: 14px;
      margin: 20px 40px 0 0;
      @include ellipsis;
      .info-item-value {
        color: #313238;
      }

      &.table {
        width: 100%;
        /deep/ .table-value {
          width: 800px;
        }
      }
    }
  }
}
.details-modules {
  margin-top: 40px;
}
.instance-list {
  padding: 0 0 0 22px;
  .instance-item {
    display: inline-flex;
    align-items: center;
    justify-content: space-between;
    width: 240px;
    margin: 10px 80px 0 0;
    padding: 0 4px;
    height: 22px;
    background-color: #FAFBFD;
    font-size: 12px;
    cursor: pointer;
    &:hover {
      .instance-name,
      .instance-diff-count {
        color: $primaryColor;
      }
      .instance-diff-count {
        font-weight: bold;
      }
    }
    .instance-diff {
      padding-right: 28px;
      position: relative;
      width: 100%
    }
    .instance-name {
      @include ellipsis;
    }
    .instance-diff-count {
      color: #C4C6CC;
    }
    .instance-change-type {
      position: absolute;
      cursor: pointer;
      right: -2px;
      top: -2px;
      width: 30px;
      height: 18px;
      border-radius: 2px;
      text-align: center;
      font-size: 12px;
      transform: scale(0.833);
      &.removed {
        color: #FF5656;
        background: #ffdddd;
      }
      &.added {
        color: #20a342;
        background: #dff9e4;
      }
      &.changed,
      &.others {
        color: #FE9C00;
        background: #FFF1DB;
      }
    }
  }
}
</style>
