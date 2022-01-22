<script>
  import FilterTag from './filter-tag.vue'
  import FilterTagItem from './general-model-filter-tag-item.vue'
  import { clearSearchQuery } from './general-model-filter.js'

  export default {
    components: {
      FilterTagItem
    },
    extends: FilterTag,
    provide() {
      return {
        condition: () => this.condition
      }
    },
    props: {
      filterSelected: {
        type: Array,
        default: () => ([])
      },
      filterCondition: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      condition() {
        return this.filterCondition
      },
      showIPTag() {
        // 替换继续的值指定为false
        return false
      },
      selected() {
        return this.filterSelected.filter((property) => {
          const { value } = this.condition[property.id]
          return value !== null && value !== undefined && !!value.toString().length
        })
      }
    },
    methods: {
      handleResetAll() {
        clearSearchQuery()
      }
    }
  }
</script>
