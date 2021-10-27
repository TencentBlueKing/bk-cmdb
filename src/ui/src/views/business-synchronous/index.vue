<template>
  <section class="batch-wrapper" v-bkloading="{ isLoading: processLoading }">
    <cmdb-tips>{{$t('同步模板功能提示')}}</cmdb-tips>
    <h2 class="title">{{$t('将会同步以下信息')}}：</h2>
    <div class="info-layout cleafix">
      <ul class="process-list fl">
        <li class="process-item"
          v-for="(process, index) in processList"
          :key="index"
          :class="{
            'show-tips': !process.confirmed,
            'is-active': currentDiff.process_template_id === process.process_template_id,
            'is-remove': process.type === 'removed'
          }"
          @click="handleChangeActive(process, index)">
          <span class="process-name" :title="process.process_template_name">{{process.process_template_name}}</span>
          <span class="process-service-count"
            v-if="process.type !== 'others' && process.total_num"
          >{{process.total_num}}</span>
        </li>
      </ul>
      <div class="change-details"
        v-bkloading="{ isLoading: detailsLoading }"
        v-if="currentDiff"
        :key="currentDiff.process_template_id">
        <cmdb-collapse class="details-info">
          <div class="collapse-title" slot="title">
            {{$t('变更内容')}}
            <span v-if="currentDiff.type === 'changed'">（{{currentDiff.changedProperties.length}}）</span>
          </div>
          <div class="info-content">
            <div class="process-info"
              v-if="currentDiff.type === 'added'">
              <div class="info-item" style="width: auto;">
                {{$t('模板中新增进程')}}
                <span class="info-item-value">{{currentDiff.process_template_name}}</span>
              </div>
            </div>
            <div class="process-info"
              v-if="currentDiff.type === 'removed'">
              <div class="info-item" style="width: auto;">
                <span class="info-item-value" style="font-weight: 700;">{{currentDiff.process_template_name}}</span>
                {{$t('从模板中删除')}}
              </div>
            </div>
            <div class="process-info clearfix"
              v-else-if="currentDiff.type === 'changed'">
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
                    :value="getChangedValue(changed)"
                    :property="changed.property">
                  </cmdb-property-value>
                </span>
              </div>
            </div>
            <div class="process-info"
              v-else-if="currentDiff.type === 'others'">
              <div class="info-item" style="width: auto;">
                {{$t('服务分类')}}：
                <span class="info-item-value">
                  {{currentDiff.service_category.parent_name}} / {{currentDiff.service_category.name}}
                </span>
              </div>
            </div>
          </div>
        </cmdb-collapse>
        <cmdb-collapse
          class="details-modules"
          v-for="(path, moduleId) of topoPath"
          :key="moduleId"
          :collapse="true"
          @collapse-change.once="handleModulesCollapseChange(moduleId)">
          <div class="collapse-title" slot="title">
            {{path}} {{$t('涉及实例')}}
          </div>
          <ul class="instance-list" v-bkloading="{ isLoading: instancesLoading }">
            <li class="instance-item"
              v-for="(instance, instanceIndex) in currentDiff.modules[moduleId].service_instances"
              :key="instanceIndex"
              @click="handleViewDiff(instance, moduleId)">
              <span class="instance-name" v-bk-overflow-tips>{{instance.service_instance.name}}</span>
              <span class="instance-diff-count"
                v-if="instance.changed_attributes">
                ({{instance.changed_attributes.length}})
              </span>
            </li>
          </ul>
        </cmdb-collapse>
      </div>
    </div>
    <div class="batch-options">
      <bk-button class="mr10" theme="primary"
        :disabled="!allConfirmed"
        :loading="confirming"
        @click="handleConfirm">
        {{$t('确认并同步')}}
      </bk-button>
      <bk-button @click="handleGoBackModule">{{$t('取消')}}</bk-button>
    </div>
    <bk-sideslider
      v-transfer-dom
      :width="676"
      :is-show.sync="slider.show"
      :title="slider.title">
      <template slot="content" v-if="slider.show">
        <instance-details slot="content"
          v-if="slider.show"
          v-bind="slider.props"
          :properties="properties">
        </instance-details>
      </template>
    </bk-sideslider>
  </section>
</template>

<script>
  import InstanceDetails from './children/details.vue'
  import formatter from '@/filters/formatter'
  import { mapGetters } from 'vuex'
  import isEmpty from 'lodash/isEmpty'
  import cloneDeep from 'lodash/cloneDeep'

  export default {
    name: 'SyncTemplate',
    components: {
      InstanceDetails
    },
    data() {
      return {
        processList: [], // 进程模板列表
        processLoading: false,
        properties: [], // 资源的所有属性，用来翻译
        topoPath: {}, // 进程模板涉及的实例的拓扑路径
        slider: {
          show: false,
          title: '',
          props: {
            module: null,
            instance: null,
            type: ''
          }
        },
        currentDiff: {
          total_num: 0, // 当前进程模板变更内容总数
          process_template_id: '', // 当前进程模板 id
          process_template_name: '', // 当前进程模板名称
          process_template: {}, // 当前进程模板内容
          changedProperties: [], // 当前进程模板实例具体更改细节
          modules: {} // 当前进程模板下的各个拓扑模块下的实例和变更
        },
        detailsLoading: false,
        instancesLoading: false,
        instanceLimit: 500,
        confirming: false
      }
    },
    computed: {
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectBiz', ['bizId']),
      templateId() {
        return Number(this.$route.params.template)
      },
      modules() {
        return String(this.$route.params.modules).split(',')
          .map(id => Number(id))
      },
      allConfirmed() {
        return this.processList.every(process => process.confirmed)
      }
    },
    async created() {
      this.initCurrentModules()
      try {
        await this.getProperties()
        this.getTopoPath()
        this.loadAllProcessTplDiffs()
      } catch (e) {
        console.error(e)
      }
    },
    methods: {
      initCurrentModules() {
        const modules = {}

        this.modules.forEach((m) => {
          modules[m] = {
            start: 0,
            service_instance_count: 0,
            service_instances: []
          }
        })

        this.currentDiff.modules = modules
      },
      handleChangeActive(process) {
        this.initCurrentModules()
        this.loadProcessTplDiff(process)
      },
      async getProperties() {
        try {
          this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params: {
              bk_obj_id: 'process',
              bk_supplier_account: this.supplierAccount,
              bk_biz_id: this.bizId
            }
          })
        } catch (e) {
          console.error(e)
        }
      },
      async getTopoPath() {
        try {
          const { nodes } = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
            bizId: this.bizId,
            params: {
              topo_nodes: this.modules.map(moduleId => ({ bk_obj_id: 'module', bk_inst_id: moduleId }))
            }
          })
          const topoPath = {}
          nodes.forEach((node) => {
            topoPath[node.topo_node.bk_inst_id] = node.topo_path.reverse().map(path => path.bk_inst_name)
              .join(' / ')
          })
          this.topoPath = topoPath
        } catch (e) {
          console.error(e)
        }
      },
      /**
       * 加载所有进程模板对比概况
       */
      async loadAllProcessTplDiffs() {
        this.processLoading = true
        this.$store.dispatch('businessSynchronous/getAllProcessTplDiffs', {
          params: {
            bk_module_ids: this.modules,
            bk_biz_id: this.bizId
          }
        }).then(async (difference) => {
          const processList = []
          const operationDiffTypes = ['changed', 'added', 'removed']

          // 模板内容变更
          Object.keys(difference).forEach((type) => {
            const diffItem = difference[type]
            if (operationDiffTypes.includes(type) && diffItem) {
              diffItem.forEach(({ id, name }) => {
                processList.push({
                  type,
                  process_template_id: id,
                  process_template_name: name,
                  confirmed: false
                })
              })
            }
          })

          if (difference.changed_attribute) {
            processList.push({
              type: 'others',
              process_template_id: 'service_category_id',
              process_template_name: this.$t('服务分类变更'),
              modules: [],
              confirmed: false
            })
          }

          processList[0].confirmed = true
          this.processList = processList
          this.loadProcessTplDiff(this.processList[0])
        })
          .finally(() => {
            this.processLoading = false
          })
      },
      /**
       * 加载单个 diff 概况
       */
      loadProcessTplDiff(process) {
        const params = {
          bk_module_ids: this.modules,
          service_template_id: this.templateId,
          bk_biz_id: this.bizId,
        }

        if (process.type === 'others') {
          params.service_category = true
        } else {
          params.process_template_id = process.process_template_id
        }

        this.detailsLoading = true
        this.$store.dispatch('businessSynchronous/getProcessTplDiff', {
          params
        }).then((diff) => {
          // 对接口数据进行转换，组成成可以适应老的 UI 模型的数据
          this.currentDiff.type = process.type
          this.currentDiff.process_template_id = process.process_template_id
          this.currentDiff.process_template_name = process.process_template_name
          this.currentDiff.service_category = diff.changed_category

          if (process.type !== 'others') {
            const changedProperties = []

            Object.keys(diff.process_template.property).forEach((key) => {
              const prop = diff.process_template.property[key]
              const formatedProp = this.properties.find(i => i.bk_property_id === key)

              if (!isEmpty(prop.value)) {
                changedProperties.push({
                  property: formatedProp,
                  template_property_value: prop.value
                })
              }
            })

            this.currentDiff.changedProperties = changedProperties
          }
          process.total_num = diff.total_num
          process.confirmed = true
        })
          .finally(() => {
            this.detailsLoading = false
          })
      },
      getChangedValue(changed) {
        const { property } = changed
        let value = changed.template_property_value
        value = Object.prototype.toString.call(value) === '[object Object]' ? value.value : value
        return formatter(value, property)
      },
      async handleModulesCollapseChange(moduleId) {
        this.loadProcessTplDiffDetails(moduleId)
      },
      /**
       * 加载进程模板变更详情
       */
      loadProcessTplDiffDetails(moduleId) {
        const theModule = this.currentDiff.modules[moduleId]
        const params = {
          bk_biz_id: this.bizId,
          bk_module_id: Number(moduleId),
          service_template_id: this.templateId,
        }

        if (this.currentDiff.type === 'others') {
          params.service_category = true
        } else {
          params.process_template_id =  this.currentDiff.process_template_id
        }

        this.instancesLoading = true
        this.$store.dispatch('businessSynchronous/getProcessTplDiffDetails', {
          params
        }).then((res) => {
          /**
           * 因为接口的数据格式变了，但是 UI 模型的结构因为时间关系没有更换，所以需要做一下数据转换，把新数据转换成可以渲染的老 UI 模型的数据。
           * 如果你要修改这块的代码，看明白这里的逻辑以后，可以优化一下。
           */
          let changedServiceInstances = {}
          // 服务分类展示方式和其他进程模板不一样，所以单独处理
          if (this.currentDiff.type === 'others') {
            changedServiceInstances.service_instances = res.service_category_detail?.service_instance.map(i => ({
              service_instance: i,
              type: 'others',
            })) || []

            changedServiceInstances.service_instances.forEach((instance) => {
              instance.changed_attributes = res.service_category_detail?.module_attribute.map(i => ({
                ...i,
                property_name: '服务分类',
                property_value: `${this.currentDiff.service_category.parent_name} / ${this.currentDiff.service_category.name}`
              }))
            })
          } else {
            changedServiceInstances = res?.service_instances
              ?.find(i => i.service_instances[0].type === this.currentDiff.type)

            // 附加单个实例的变更属性
            if (changedServiceInstances?.service_instances?.length) {
              changedServiceInstances?.service_instances
                .forEach((instance) => {
                  if (!instance?.changed_attributes) {
                    instance.changed_attributes = this.currentDiff.changedProperties.map(i => ({
                      property_name: i.property.bk_property_name,
                      ...i
                    }))
                  }
                })
            }
          }

          theModule.service_instance_count = changedServiceInstances?.service_instance_count || 0
          theModule.service_instances = theModule.service_instances.concat(changedServiceInstances?.service_instances)
        })
          .finally(() => {
            this.instancesLoading = false
          })
      },
      handleViewDiff(instance, moduleId) {
        this.slider.title = instance.service_instance.name
        const instanceDetail = cloneDeep(instance)

        // 为了适应老的 UI 模型做的数据转换
        instanceDetail.changed_attributes = instanceDetail.changed_attributes.map((i) => {
          if (!i.property_id) {
            return {
              ...i,
              property_id: i.property.bk_property_id
            }
          }
          return i
        })

        this.slider.props = {
          module: this.currentDiff.modules[moduleId],
          instance: instanceDetail,
          type: this.currentDiff.type
        }
        this.slider.show = true
      },
      handleConfirm() {
        this.confirming = true
        this.$store.dispatch('businessSynchronous/syncServiceInstanceByTemplate', {
          params: {
            service_template_id: this.templateId,
            bk_module_ids: this.modules,
            bk_biz_id: this.bizId
          }
        }).then(() => {
          this.$success(this.$t('同步成功'))
          this.handleGoBackModule()
        })
          .finally(() => {
            this.confirming = false
          })
      },
      handleGoBackModule() {
        this.$routerActions.back()
      },
    }
  }
</script>

<style lang="scss" scoped>
    .batch-wrapper {
        padding: 10px 20px;
        .title {
            margin-top: 24px;
            font-size: 14px;
            line-height: 20px;
        }
        .collapse-title {
            font-size: 14px;
            color: $textColor;
        }
    }
    .info-layout {
        margin-top: 10px;
        border: 1px solid $borderColor;
        border-bottom: none;
        height: calc(100vh - 350px);
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
            height: 61px;
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
                    color: $primaryColor;
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
                line-height: 60px;
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
        margin-top: 60px;
        & ~ .details-modules {
            margin-top: 20px;
        }
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
            .instance-name {
                padding-right: 28px;
                position: relative;
                @include ellipsis;
            }
            .instance-diff-count {
                color: #C4C6CC;
            }
            .instance-change-type {
                position: absolute;
                right: -2px;
                top: -2px;
                width: 30px;
                height: 18px;
                border-radius: 2px;
                text-align: center;
                font-size: 12px;
                transform: scale(0.833);
                &.del {
                    color: #ea3636;
                    background: #ffdddd;
                }
                &.add {
                    color: #20a342;
                    background: #dff9e4;
                }
            }
        }
    }

    .batch-options {
        margin-top: 20px;
    }
</style>
