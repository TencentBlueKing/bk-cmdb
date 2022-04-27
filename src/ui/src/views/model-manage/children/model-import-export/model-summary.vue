<template>
  <div class="model-summary">
    <span class="model-icon-container" :class="{ 'is-builtin': isBuiltin }">
      <i
        class="model-icon"
        :class="modelIconClass"
      ></i>
    </span>
    <span class="model-name">{{ modelName }}</span>
    <span class="model-id">({{ modelId }})</span>
  </div>
</template>
<script>
  import { defineComponent, computed } from '@vue/composition-api'

  export default defineComponent({
    name: 'ModelSummary',
    props: {
      /**
       * 模型数据
       * @property {Object} data
       * @property {Boolean} data.ispre 是否为内置模型
       * @property {String} data.bk_obj_icon 模型图标
       * @property {String} data.bk_obj_id 模型 ID
       * @property {String} data.bk_obj_name 模型名称
       */
      data: {
        type: Object,
        required: true,
        default: null
      }
    },
    setup({ data }) {
      const isBuiltin = computed(() => data.ispre)
      const modelIconClass = computed(() => data.bk_obj_icon || '')
      const modelName = computed(() => data.bk_obj_name || '')
      const modelId = computed(() => data.bk_obj_id || '')

      return {
        isBuiltin,
        modelIconClass,
        modelName,
        modelId
      }
    }
  })
</script>

<style lang="scss" scoped>
.model {
  &-summary {
    display: flex;
    font-size: 12px;
    align-items: center;
  }

  &-icon-container {
    width: 26px;
    height: 26px;
    background: #e1ecff;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    @at-root .is-builtin & {
      background-color: #f0f1f5;
    }
  }

  &-icon {
    color: #3a84ff;

    @at-root .is-builtin & {
      color: #798aad;
    }
  }

  &-name {
    margin-left: 12px;
  }

  &-id {
    color: #979ba5;
    margin-left: 5px;
  }
}
</style>
