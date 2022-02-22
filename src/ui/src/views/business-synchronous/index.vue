<template>
  <section class="batch-wrapper" v-bkloading="{ isLoading: processListLoading }">
    <cmdb-tips>{{$t('同步模板功能提示')}}</cmdb-tips>
    <h2 class="title">{{$t('将会同步以下信息')}}：</h2>
    <div class="info-layout cleafix">

      <!-- 进程模板列表 -->
      <ul class="process-list fl">
        <li class="process-item"
          v-for="(process, index) in processList"
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
      <div class="change-details"
        v-bkloading="{ isLoading: currentDiffLoading }"
        v-if="currentDiff"
        :key="`${currentDiff.process_template_id}-${currentDiff.process_template_name}`">
        <cmdb-collapse class="details-info">
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

            <!-- 服务分类变更 -->
            <div class="process-info"
              v-if="currentDiff.type === 'others'">
              <div class="info-item" style="width: auto;">
                {{$t('服务分类')}}：
                <span class="info-item-value">
                  {{currentDiff.changed_service_category}}
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
          @collapse-change="handleInstanceCollapseChange(moduleId, $event)">
          <div class="collapse-title" slot="title">
            {{path}} {{$t('涉及实例')}}
            <span v-if="currentDiff.modules[moduleId].serviceInstanceCount !== ''">
              ({{currentDiff.modules[moduleId].serviceInstanceCount}})
            </span>
          </div>

          <!-- 实例列表 -->
          <ul
            class="instance-list"
            v-bkloading="{ isLoading: currentDiff.modules[moduleId].instancesLoading }">
            <li class="instance-item"
              v-for="(instance, instanceIndex) in currentDiff.modules[moduleId].serviceInstances"
              :key="instanceIndex"
              @click="viewInstanceDiff(instance, moduleId)">
              <span class="instance-name" v-bk-overflow-tips>
                {{instance.name}}
                <label
                  :class="['instance-change-type', instance.type]"
                  v-if="translateChangedType(instance.type)">
                  {{translateChangedType(instance.type)}}
                </label>
              </span>
            </li>
          </ul>
        </cmdb-collapse>
      </div>
    </div>

    <div class="batch-options">
      <bk-button class="mr10" theme="primary"
        :loading="confirming"
        @click="confirmAndSync">
        {{$t('确认并同步')}}
      </bk-button>
      <bk-button @click="goBackModule">{{$t('取消')}}</bk-button>
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
  </section>
</template>

<script>
  import ServiceInstanceDetails from './children/service-instance-details.vue'
  import formatter from '@/filters/formatter'
  import { mapGetters } from 'vuex'
  import isEmpty from 'lodash/isEmpty'
  import to from 'await-to-js'

  export default {
    name: 'BusinessSynchronous',
    components: {
      ServiceInstanceDetails
    },
    data() {
      return {
        processListLoading: false,
        properties: [], // 资源的所有属性，用来翻译
        topoPath: {}, // 进程模板涉及的实例的拓扑路径
        processList: [], // 进程模板列表
        currentDiff: {
          process_template_id: 0, // 当前进程模板 id
          process_template_name: '', // 当前进程模板名称
          process_template: {}, // 当前进程模板内容
          changedProperties: [], // 当前进程模板实例具体更改细节
          modules: {} // 当前进程模板下的各个拓扑模块下的实例和变更
        },
        currentDiffLoading: false,
        instancesLoading: false,
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
        },
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
    },
    async created() {
      this.initCurrentModules()
      this.processListLoading = true
      await to(this.loadProperties())
      await to(this.loadTopoPath())
      await to(this.loadProcessList())
      this.processListLoading = false
    },
    methods: {
      /**
       * 初始化当前模块，用于装载各个模块的实例信息
       */
      initCurrentModules() {
        const modules = {}

        this.modules.forEach((m) => {
          modules[m] = {
            serviceInstanceCount: '',
            serviceInstances: []
          }
        })

        this.currentDiff.modules = modules
      },
      /**
       * 加载进程属性，便于转换成可读中文
       */
      loadProperties() {
        return this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
          params: {
            bk_obj_id: 'process',
            bk_supplier_account: this.supplierAccount,
            bk_biz_id: this.bizId
          }
        }).then((data) => {
          this.properties = data
        })
          .catch(() => {
            this.properties = []
          })
      },
      /**
       * 加载拓扑路径，用于加载涉及实例
       */
      loadTopoPath() {
        return this.$store.dispatch('objectMainLineModule/getTopoPath', {
          bizId: this.bizId,
          params: {
            topo_nodes: this.modules.map(moduleId => ({ bk_obj_id: 'module', bk_inst_id: moduleId }))
          }
        }).then(({ nodes }) => {
          const topoPath = {}

          nodes.forEach((node) => {
            topoPath[node.topo_node.bk_inst_id] = node.topo_path.reverse().map(path => path.bk_inst_name)
              .join(' / ')
          })

          this.topoPath = topoPath
        })
          .catch(() => {
            this.topoPath = []
          })
      },
      /**
       * 加载服务下的全部分类，用来翻译服务分类变更内容
       */
      loadServiceCategories() {
        return this.$store.dispatch('serviceClassification/searchServiceCategory', {
          params: { bk_biz_id: this.bizId }
        }).then(({ info }) => {
          this.serviceCategories = info || []
        })
          .catch(() => {
            this.serviceCategories = []
          })
      },
      /**
       * 加载进程模板列表
       */
      loadProcessList() {
        return this.$store.dispatch('businessSynchronous/getAllProcessTplDiffs', {
          params: {
            bk_module_ids: this.modules,
            bk_biz_id: this.bizId,
            service_template_id: this.templateId
          }
        }).then((difference) => {
          const processList = []
          const operationDiffTypes = ['changed', 'added', 'removed']

          // 模板内容变更
          Object.keys(difference).forEach((type) => {
            const diffItem = difference[type]
            if (operationDiffTypes.includes(type) && diffItem) {
              diffItem.forEach(({ id, name }) => {
                processList.push(this.genDiffItem({
                  diffType: type,
                  processId: id,
                  processName: name,
                }))
              })
            }
          })

          if (difference.changed_attribute) {
            processList.push(this.genDiffItem({
              diffType: 'others'
            }))
          }

          const firstProcess = processList[0]
          firstProcess.confirmed = true
          if (firstProcess.type === 'others') {
            this.loadServiceCategory(firstProcess)
          } else {
            this.loadProcessDiff(firstProcess)
          }

          this.processList = processList
        })
      },
      /**
       * 生成对比项
       * @param {string} diffType 必须，变更类型
       * @param {string} processId 非必须，进程模板的变更 ID
       * @param {string} processName 非必须，进程模板的名称
       */
      genDiffItem({
        diffType,
        processId,
        processName
      }) {
        // 服务分类因为是修改服务模板的属性，所以视图比较特别
        const serviceCategoryDiffItem = {
          type: 'others',
          process_template_id: 'service_category_id',
          process_template_name: this.$t('服务分类变更'),
          modules: [],
          confirmed: false
        }

        if (diffType === 'others') {
          return serviceCategoryDiffItem
        }

        return {
          type: diffType,
          process_template_id: processId,
          process_template_name: processName,
          confirmed: false
        }
      },
      /**
       * 加载进程模板变更内容
       * @param {Object} process 进程信息
       */
      async loadProcessDiff(process) {
        this.initCurrentModules()

        if (process.type === 'others') {
          return this.loadServiceCategory(process)
        }

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
      /**
       * 加载服务分类变更
       * @param {Object} process 进程信息
       */
      async loadServiceCategory(process) {
        this.currentDiffLoading = true

        await to(this.loadServiceCategories())
        const [, { template }] = await to(this.getServiceTemplateDetail())
        const newCategoryId = template.service_category_id

        const category = this.getCategoryById(newCategoryId)
        const parentCategory = this.getCategoryById(category.bk_parent_id)

        this.currentDiff.type = process.type
        this.currentDiff.process_template_id = process.process_template_id
        this.currentDiff.process_template_name = process.process_template_name
        this.currentDiff.changed_service_category = `${parentCategory.name} / ${category?.name || ''}`

        this.currentDiffLoading = false
        process.confirmed = true
      },
      /**
       * 通过分类 ID 获取 分类对象
       * @param {Number} categoryId 分类 ID
       * @return {Object} 分类对象
       */
      getCategoryById(categoryId) {
        return this.serviceCategories.find(item => item.category.id === categoryId)?.category || {}
      },
      /**
       * 获取服务模板详情，用于展示最新的服务模板变更信息
       */
      getServiceTemplateDetail() {
        return this.$store.dispatch('serviceTemplate/findServiceTemplate', {
          id: this.templateId
        })
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
        const theModule = this.currentDiff.modules[moduleId]
        this.$set(theModule, 'instancesLoading', true)
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
            this.$set(theModule, 'instancesLoading', false)
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
        } else if (this.currentDiff.type === 'others') {
          // 服务分类单独加载
          params.service_category = true
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
      confirmAndSync() {
        this.confirming = true
        this.$store.dispatch('businessSynchronous/syncServiceInstanceByTemplate', {
          params: {
            service_template_id: this.templateId,
            bk_module_ids: this.modules,
            bk_biz_id: this.bizId
          }
        }).then(() => {
          this.$success(this.$t('同步成功'))
          this.goBackModule()
        })
          .finally(() => {
            this.confirming = false
          })
      },
      goBackModule() {
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
                    color: #ea3636;
                    background: #ffdddd;
                }
                &.added {
                    color: #20a342;
                    background: #dff9e4;
                }
                &.changed,
                &.others {
                    color: $primaryColor;
                    background: #3a84ff29;
                }
            }
        }
    }

    .batch-options {
        margin-top: 20px;
    }
</style>
