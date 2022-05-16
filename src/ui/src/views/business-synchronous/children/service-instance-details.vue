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
  <div class="instance-details-wrapper" v-bkloading="{ isLoading: $loading() }">
    <bk-table
      :data="diffDetails"
      :max-height="$APP.height - 120">
      <bk-table-column prop="property_name" :label="$t('属性名称')" show-overflow-tooltip></bk-table-column>

      <bk-table-column prop="property_value" :label="$t('变更前')">
        <template slot-scope="{ row }">
          <process-bind-info-value
            v-if="row.property_id === 'bind_info'"
            :value="getPropertyValue(row, 'before')"
            :property="properties.find(property => property.bk_property_id === row.property_id)">
          </process-bind-info-value>

          <cmdb-property-value
            v-else
            v-bk-tooltips="{
              content: getPropertyValue(row, 'before'),
              duration: 500
            }"
            :value="getPropertyValue(row, 'before')"
            :property="properties.find(property => property.bk_property_id === row.property_id)">
          </cmdb-property-value>
        </template>
      </bk-table-column>

      <bk-table-column prop="show_value" :label="$t('变更后')">
        <template slot-scope="{ row }">
          <process-bind-info-value
            v-if="row.property_id === 'bind_info'"
            :value="getPropertyValue(row, 'after')"
            :property="properties.find(property => property.bk_property_id === row.property_id)">
          </process-bind-info-value>

          <cmdb-property-value
            v-else
            v-bk-tooltips="{
              content: getPropertyValue(row, 'after'),
              duration: 500
            }"
            :value="getPropertyValue(row, 'after')"
            :property="properties.find(property => property.bk_property_id === row.property_id)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
    </bk-table>
  </div>
</template>

<script>
  import formatter from '@/filters/formatter'
  import { mapGetters } from 'vuex'
  import ProcessBindInfoValue from '@/components/service/process-bind-info-value'
  import to from 'await-to-js'
  import isEmpty from 'lodash/isEmpty'
  export default {
    name: 'service-instance-details',
    components: {
      ProcessBindInfoValue
    },
    props: {
      /**
       * 实例对比请求的参数
       */
      diffRequestParams: {
        type: Object,
        default: () => ({})
      },
      /**
       * 翻译属性的集合
       */
      properties: {
        type: Array,
        default: () => ([])
      },
      getCategoryById: {
        type: Function,
        default: () => {}
      }
    },
    data() {
      return {
        diffDetails: []
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId'])
    },
    created() {
      this.loadInstanceDiff()
    },
    methods: {
      /**
       * 加载实例对比信息
       */
      loadInstanceDiff() {
        this.$store.dispatch('businessSynchronous/getInstanceDiff', { params: this.diffRequestParams })
          .then(async ({ changed_attributes: changedAttributes, module_attribute: moduleAttribute, process, type }) => {
            if (type === 'removed') {
              this.diffDetails = this.loadRemovedData(process)
              return
            }

            if (type === 'added') {
              [, this.diffDetails] = await to(this.loadAddedData())
              return
            }

            if (type === 'others') {
              this.diffDetails = this.loadServiceCategoryData(moduleAttribute)
              return
            }

            this.diffDetails = changedAttributes
          })
          .catch(() => {
            this.diffDetails = []
          })
      },
      loadServiceCategoryData(moduleAttribute) {
        return moduleAttribute.map((attr) => {
          const oldCategory = this.getCategoryById(attr.property_value)
          const category = this.getCategoryById(attr.template_property_value)
          return {
            property_name: this.$t('服务分类'),
            property_value: oldCategory.name,
            template_property_value: category.name,
          }
        })
      },
      loadAddedData() {
        const { process_template_id: processTemplateId, service_template_id } = this.diffRequestParams
        return this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
          params: {
            bk_biz_id: this.bizId,
            service_template_id
          }
        }).then(({ info }) => {
          const process = info.find(process => process.id === processTemplateId)
          const diffDetails = []

          Object.keys(process.property).forEach((key) => {
            const property = this.properties.find(property => property.bk_property_id === key)

            if (!isEmpty(process?.property[key]?.value)) {
              diffDetails.push({
                property_id: key,
                property_name: property.bk_property_name,
                property_value: null,
                template_property_value: {
                  ...process.property[key]
                }
              })
            }
          })

          return diffDetails
        })
          .catch((err) => {
            console.log(err)
          })
      },
      /**
       * 用 process 信息组合已删除的详情信息
       * @param {Object} process 进程对象
       * @return {Array} 对比详情
       */
      loadRemovedData(process) {
        const diffDetails = []

        Object.keys(process).forEach((key) => {
          const property = this.getProperty(key)
          if (property && !isEmpty(process[key])) {
            diffDetails.push({
              property_id: key,
              property_name: property.bk_property_name,
              property_value: process[key],
              template_property_value: {
                value: this.$t('该进程已删除')
              }
            })
          }
        })

        return diffDetails
      },
      getProperty(id) {
        return this.properties.find(property => property.bk_property_id === id)
      },
      /**
       * 获取实例属性值
       * @param instance 实例对象
       * @param type 变更前后类型 after 为获取变更后端额值，before 为获取变更前的值
       */
      getPropertyValue(instance, type) {
        const propertyId = instance.property_id
        let value = instance.property_value
        const templateValue = instance.template_property_value

        if (type === 'after') {
          value = templateValue?.value || templateValue
        }

        // others 单独处理
        if (this.type !== 'others') {
          return formatter(value, this.getProperty(propertyId))
        }

        return formatter(value, 'singlechar')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .instance-details-wrapper {
        padding: 20px;
    }
</style>
