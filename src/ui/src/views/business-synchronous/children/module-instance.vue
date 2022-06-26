<script lang="ts">
  import { computed, defineComponent } from '@vue/composition-api'
  import ProcessDifference from './process-difference.vue'
  import PropertyDifference from './property-difference.vue'

  export default defineComponent({
    components: {
      ProcessDifference,
      PropertyDifference
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
        required: true
      },
      propertyDiff: {
        type: Array,
        required: true
      },
      processDiff: {
        type: Array,
        required: true
      },
      modelProperty: {
        type: Object,
        default: () => ({}),
        required: true
      }
    },
    setup(props) {
      const hasPropertyDiff = computed(() => props.propertyDiff?.length > 0)
      const hasProcessDiff = computed(() => props.processDiff?.length > 0)

      return {
        hasPropertyDiff,
        hasProcessDiff
      }
    }
  })
</script>

<template>
  <div class="module-instance">
    <cmdb-collapse class="property-container" v-if="hasPropertyDiff"
      :label="$t('属性变更')"
      arrow-type="filled">
      <property-difference
        :module-id="moduleId"
        :template-id="templateId"
        :property-diff="propertyDiff">
      </property-difference>
    </cmdb-collapse>

    <cmdb-collapse class="process-container" v-if="hasProcessDiff"
      :label="$t('进程信息变更')"
      arrow-type="filled">
      <process-difference
        :module-id="moduleId"
        :template-id="templateId"
        :topo-path="topoPath"
        :process-diff="processDiff"
        :properties="modelProperty.process">
      </process-difference>
    </cmdb-collapse>
  </div>
</template>

<style lang="scss" scoped>
  .module-instance {
    .property-container {
      .property-difference {
        margin-top: 12px;
      }
    }

    .process-container {
      margin-top: 16px;
      .process-difference {
        margin-top: 12px;
      }
    }
  }
</style>
