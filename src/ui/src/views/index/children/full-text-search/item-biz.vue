<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
      <i class="tag-disabled" v-if="data.source.bk_data_status === 'disabled'">{{$t('已归档')}}</i>
    </div>
    <div class="result-desc" v-if="properties" @click="data.linkTo(data.source)">
      <template v-for="(property, childIndex) in properties">
        <div class="desc-item"
          :key="childIndex"
          v-if="data.source[property.bk_property_id]"
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
    name: 'item-biz',
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

      const properties = computed(() => propertyMap.value.biz)

      return {
        properties,
        getText,
        getHighlightValue
      }
    }
  })
</script>
