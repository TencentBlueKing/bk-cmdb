<template>
  <div class="details-layout">
    <bk-tab class="details-tab"
      type="unborder-card"
      :active.sync="active">
      <bk-tab-panel name="property" :label="$t('属性')">
        <cmdb-property
          :resource-type="resourceType"
          :properties="properties"
          :property-groups="propertyGroups"
          :inst="inst">
          <template #append>
            <div class="group business-scope-group">
              <div class="group-title">{{$t('资源定义范围')}}</div>
              <ul class="property-list">
                <li class="property-item full-width">
                  <div class="property-name">
                    <span class="property-name-text">{{$t('业务范围')}}</span>
                  </div>
                  <div class="property-value">
                    <business-scope-settings-details
                      class="form-component"
                      :data="inst.bk_scope" />
                  </div>
                </li>
              </ul>
            </div>
          </template>
        </cmdb-property>
      </bk-tab-panel>
      <bk-tab-panel name="history" :label="$t('变更记录')">
        <cmdb-audit-history v-if="active === 'history'"
          :resource-type="resourceType"
          :obj-id="objId"
          :resource-id="bizSetId">
        </cmdb-audit-history>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<script>
  import { defineComponent, computed, reactive, watchEffect, toRefs } from '@vue/composition-api'
  import store from '@/store'
  import {
    BUILTIN_MODELS,
    BUILTIN_MODEL_PROPERTY_KEYS,
    BUILTIN_MODEL_ROUTEPARAMS_KEYS,
    BUILTIN_MODEL_RESOURCE_TYPES
  } from '@/dictionary/model-constants.js'
  import router from '@/router/index.js'
  import RouterQuery from '@/router/query'
  import cmdbProperty from '@/components/model-instance/property'
  import cmdbAuditHistory from '@/components/model-instance/audit-history'
  import businessScopeSettingsDetails from '@/components/business-scope/settings-details.vue'
  import propertySetService from '@/service/business-set/index.js'
  import propertyService from '@/service/property/property.js'
  import propertyGroupService from '@/service/property/group.js'

  export default defineComponent({
    components: {
      cmdbProperty,
      cmdbAuditHistory,
      businessScopeSettingsDetails
    },
    setup() {
      const MODEL_NAME_KEY = BUILTIN_MODEL_PROPERTY_KEYS[BUILTIN_MODELS.BUSINESS_SET].NAME
      const MODEL_ROUTEPARAMS_KEY = BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.BUSINESS_SET]

      const requestId = Symbol('details')

      const query = computed(() => RouterQuery.getAll())
      const bizSetId = computed(() => parseInt(router.app.$route.params[MODEL_ROUTEPARAMS_KEY], 10))

      const getModelById = store.getters['objectModelClassify/getModelById']
      const model = computed(() => getModelById(BUILTIN_MODELS.BUSINESS_SET) || {})

      const resourceType = BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET]

      const state = reactive({
        inst: {},
        objId: BUILTIN_MODELS.BUSINESS_SET,
        properties: [],
        propertyGroups: [],
        active: query.tab || 'property'
      })

      watchEffect(async () => {
        const [modelInst, modelProperties, modelPropertyGroups] = await Promise.all([
          propertySetService.findById(bizSetId.value, { requestId, cancelPrevious: true }),
          propertyService.findBizSet(true),
          propertyGroupService.findBizSet()
        ])
        state.inst = modelInst
        state.properties = modelProperties
        state.propertyGroups = modelPropertyGroups
        store.commit('setTitle', `${model.value.bk_obj_name}【${modelInst[MODEL_NAME_KEY]}】`)
      })

      return {
        bizSetId,
        resourceType,
        ...toRefs(state)
      }
    }
  })
</script>

<style lang="scss" scoped>
  .details-layout {
    overflow: hidden;
    .details-tab {
      min-height: 400px;
      /deep/ {
        .bk-tab-header {
          padding: 0;
          margin: 0 20px;
        }
        .bk-tab-section {
          @include scrollbar-y;
          padding-bottom: 10px;
        }
      }
    }
  }

  .business-scope-group {
    .group-name {
      &::before {
        display: none;
      }
    }
    .property-list {
      .property-item {
        &.full-width {
          flex: 1;
          max-width: 100%;

          .property-value {
            flex: 1;
            max-width: 100%;
            display: block;
          }
        }
      }
    }
  }
</style>
