<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" @click="data.linkTo(data.source)">
      <div class="desc-item hl" v-html="`${$t('实例ID')}：${getHighlightValue(data.source.bk_inst_id, data)}`"> </div>
      <template v-for="(property, childIndex) in properties">
        <div class="desc-item hl"
          :key="childIndex"
          v-html="`${getHighlightValue(property.bk_property_name, data)}：${getText(property, data)}`">
        </div>
      </template>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from '@vue/composition-api'
  import { getText, getHighlightValue } from './use-item.js'

  export default defineComponent({
    name: 'item-instance',
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
      const { data, propertyMap } = toRefs(props)

      const properties = computed(() => propertyMap.value[data.value.source.bk_obj_id])

      return {
        properties,
        getText,
        getHighlightValue
      }
    }
  })
</script>
