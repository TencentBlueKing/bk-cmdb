<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" v-if="properties" @click="data.linkTo(data.source)">
      <template v-for="(property, childIndex) in properties">
        <div class="desc-item hl"
          :key="childIndex"
          v-html="`${property.bk_property_name}：${getText(property, data)}`">
        </div>
      </template>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from '@vue/composition-api'
  import { getText, getHighlightValue } from './use-item.js'

  export default defineComponent({
    name: 'item-module',
    props: {
      data: {
        type: Object,
        default: () => ({})
      },
      propertyMap: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props, { root }) {
      const { propertyMap } = toRefs(props)

      const properties = computed(() => {
        const properties = (propertyMap.value.module || []).slice()
        properties.unshift({
          bk_property_id: 'bk_module_id',
          bk_property_name: root.$t('ID'),
        })
        return properties
      })

      return {
        properties,
        getText,
        getHighlightValue
      }
    }
  })
</script>
