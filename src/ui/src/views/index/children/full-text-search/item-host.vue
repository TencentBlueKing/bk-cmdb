<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" v-if="properties" @click="data.linkTo(data.source)">
      <div class="desc-item"
        v-html="`${$t('主机ID')}：${getHighlightValue(data.source.bk_host_id, data, 'bk_host_id')}`">
      </div>
      <template v-for="(property, childIndex) in properties">
        <div class="desc-item"
          v-if="data.source[property.bk_property_id]"
          :key="childIndex"
          v-html="`${property.bk_property_name}：${getText(property, data, property.bk_property_id)}`">
        </div>
      </template>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from '@vue/composition-api'
  import { getText, getHighlightValue } from './use-item.js'

  export default defineComponent({
    name: 'item-host',
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
    setup(props) {
      const { propertyMap } = toRefs(props)

      const properties = computed(() => propertyMap.value.host)

      return {
        properties,
        getText,
        getHighlightValue
      }
    }
  })
</script>
