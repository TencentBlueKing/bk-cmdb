<script>
  import PropertySelector from './property-selector.vue'
  export default {
    extends: PropertySelector,
    props: {
      objId: {
        type: String
      },
      properties: {
        type: Array,
        default: () => ([])
      },
      propertyGroups: {
        type: Array,
        default: () => ([])
      },
      propertySelected: {
        type: Array,
        default: () => ([])
      },
      handler: {
        type: Function,
        default: () => {}
      }
    },
    data() {
      return {
        selected: [...this.propertySelected],
      }
    },
    computed: {
      propertyMap() {
        const modelPropertyMap = { [this.objId]: this.properties }
        const ignoreProperties = [] // 预留，需要忽略的属性
        // eslint-disable-next-line max-len
        modelPropertyMap[this.objId] = modelPropertyMap[this.objId].filter(property => !ignoreProperties.includes(property.bk_property_id))
        return modelPropertyMap
      }
    },
    methods: {
      async confirm() {
        this.handler(this.selected)
        this.close()
      }
    }
  }
</script>
