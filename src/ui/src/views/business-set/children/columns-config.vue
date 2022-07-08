<template>
  <bk-sideslider
    v-transfer-dom
    :is-show.sync="isShow"
    :width="600"
    :title="$t('列表显示属性配置')"
    :before-close="hide">
    <cmdb-columns-config
      slot="content"
      v-if="isShow"
      :properties="properties"
      :selected="columnsConfig.selected"
      :disabled-columns="columnsConfig.disabledColumns"
      @on-apply="handleApplayColumnsConfig"
      @on-cancel="hide"
      @on-reset="handleResetColumnsConfig">
    </cmdb-columns-config>
  </bk-sideslider>
</template>

<script>
  import { defineComponent, computed, reactive, toRefs, watchEffect } from '@vue/composition-api'
  import store from '@/store'
  import { getHeaderProperties, getHeaderPropertyName } from '@/utils/tools.js'
  import cmdbColumnsConfig from '@/components/columns-config/columns-config.vue'

  export default defineComponent({
    components: {
      cmdbColumnsConfig
    },
    props: {
      show: {
        type: Boolean,
        default: false
      },
      properties: {
        type: Array,
        required: true,
        default: () => ([])
      }
    },
    setup(props, { emit }) {
      const { show: isShow, properties } = toRefs(props)
      const columnsConfigKey = 'biz_set_custom_table_columns'
      const globalUsercustom = computed(() => store.getters['userCustom/globalUsercustom'])
      const usercustom = computed(() => store.getters['userCustom/usercustom'])
      const customBusinessSetColumns = computed(() => usercustom.value[columnsConfigKey] || [])
      const globalCustomColumns = computed(() =>  globalUsercustom.value?.biz_set_global_custom_table_columns || [])

      const columnsConfig = reactive({
        selected: [],
        disabledColumns: ['bk_biz_set_id', 'bk_biz_set_name']
      })

      watchEffect(() => {
        // eslint-disable-next-line max-len
        const customColumns = customBusinessSetColumns.value.length ? customBusinessSetColumns.value : globalCustomColumns.value
        const headerProperties = getHeaderProperties(properties.value, customColumns, columnsConfig.disabledColumns)
        const tableHeader = headerProperties.map(property => ({
          id: property.bk_property_id,
          name: getHeaderPropertyName(property),
          property
        }))
        columnsConfig.selected = headerProperties.map(property => property.bk_property_id)
        emit('update-header', tableHeader)
      })

      const handleApplayColumnsConfig = (properties) => {
        store.dispatch('userCustom/saveUsercustom', {
          [columnsConfigKey]: properties.map(property => property.bk_property_id)
        })
        hide()
      }

      const handleResetColumnsConfig = () => {
        store.dispatch('userCustom/saveUsercustom', {
          [columnsConfigKey]: []
        })
        hide()
      }

      const hide = () => {
        emit('update:show', false)
      }

      return {
        isShow,
        columnsConfig,
        hide,
        handleApplayColumnsConfig,
        handleResetColumnsConfig
      }
    }
  })
</script>
