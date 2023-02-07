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
  <div class="business-scope-settings-details">
    <div v-if="isMatchAll" class="condition-list"><div class="condition-item">{{$t('全业务')}}</div></div>
    <template v-else>
      <div class="condition-list" v-if="selectedBusiness.length">
        <dl class="condition-item">
          <dt class="item-title">{{$t('业务')}}：</dt>
          <dd class="item-content">
            <span class="item-value"
              v-for="(business, index) in selectedBusiness" :key="index">{{business.bk_biz_name}}</span>
          </dd>
        </dl>
      </div>
      <div class="condition-list" v-if="condition.length">
        <dl class="condition-item" v-for="item in condition" :key="item.property.bk_property_id">
          <dt class="item-title">{{item.property.bk_property_name}}：</dt>
          <cmdb-property-value tag="dd" class="item-content"
            :value="item.value"
            :multiple="Array.isArray(item.value) && item.value.length > 1"
            :property="item.property">
          </cmdb-property-value>
        </dl>
      </div>
    </template>
  </div>
</template>

<script>
  import { defineComponent, watchEffect, watch, reactive, toRefs, ref, computed } from 'vue'
  import Utils from '@/components/filters/utils'
  import propertyService from '@/service/property/property.js'
  import businessService from '@/service/business/search.js'

  export default defineComponent({
    props: {
      data: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props) {
      const { data: scopeData } = toRefs(props)

      // 是否“全业务”（未选择任何业务及条件）
      const isMatchAll = computed(() => scopeData.value.match_all)

      // 业务属性列表
      const properties = ref([])

      // 所有业务列表
      const allBusiness = ref([])

      // 业务范围配置
      const scopeSettings = reactive({
        selectedBusiness: [],
        condition: []
      })

      // 更新范围配置方法
      const updateScopeSettings = () => {
        const selectedBusiness = []
        const condition = []
        scopeData.value?.filter?.rules?.forEach((rule) => {
          if (rule.field === 'bk_biz_id') {
            rule.value.forEach((bizId) => {
              const business = allBusiness.value.find(business => business.bk_biz_id === bizId)
              if (business) {
                selectedBusiness.push(business)
              }
            })
          } else {
            condition.push({
              field: rule.field,
              property: Utils.findProperty(rule.field, properties.value) || {},
              value: (Array.isArray(rule.value) && rule.value.length > 1) ? rule.value : rule.value[0]
            })
          }
        })
        scopeSettings.selectedBusiness = selectedBusiness
        scopeSettings.condition = condition
      }

      // 初始化业务属性和全量业务列表
      watchEffect(async () => {
        const [businessProperties, businessList] = await Promise.all([
          propertyService.findBiz(),
          businessService.findAll()
        ])

        const allowedPropertyTypes = ['organization', 'enum']
        properties.value = businessProperties.filter(item => allowedPropertyTypes.includes(item.bk_property_type))
        allBusiness.value = businessList

        updateScopeSettings()
      })

      // scopeData属性可能后更新，因此需要监听变化并更新配置数据
      watch(scopeData, () => updateScopeSettings())

      return {
        isMatchAll,
        ...toRefs(scopeSettings)
      }
    }
  })
</script>

<style lang="scss" scoped>
  .business-scope-settings-details {
    .condition-list {
      display: flex;
      flex-wrap: wrap;

      .condition-item {
        display: flex;
        align-items: center;
        color: #63656E;
        background: #F0F1F5;
        margin-right: 24px;
        margin-bottom: 8px;
        border-radius: 2px;
        padding: 2px 4px;
        position: relative;

        &:last-child {
          margin-right: 0;
          &::after {
            display: none;
          }
        }

        &::after {
          color: #c4c6cc;
          content: "&";
          position: absolute;
          right: -17px;
        }

        .item-title {
          white-space: nowrap;
        }
        .item-content {
          word-break: keep-all;

          .item-value {
            &::after {
              content: ", ";
            }

            &:last-child {
              &::after {
                display: none;
              }
            }
          }
        }
      }
    }
  }
</style>
