<script>
  import FilterTagForm from './filter-tag-form.vue'
  import { setSearchQueryByCondition } from './general-model-filter.js'

  export default {
    extends: FilterTagForm,
    props: {
      condition: {
        type: Object,
        default: () => ({})
      }
    },
    computed: {
      operator: {
        get() {
          return this.localOperator || this.condition[this.property.id].operator
        },
        set(operator) {
          this.localOperator = operator
        }
      },
      value: {
        get() {
          if (this.localValue === null) {
            return this.condition[this.property.id].value
          }
          return this.localValue
        },
        set(value) {
          this.localValue = value
        }
      }
    },
    methods: {
      handleConfirm() {
        // 构建单个condition
        const condition = {
          [this.property.id]: {
            operator: this.operator,
            value: this.value
          }
        }
        setSearchQueryByCondition(condition, [this.property])
        this.$emit('confirm')
      }
    }
  }
</script>
