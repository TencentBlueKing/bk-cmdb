import activeMixin from './active'
export default {
  mixins: [activeMixin],
  props: {
    value: {
      type: [String, Array]
    }
  },
  computed: {
    multiple() {
      return Array.isArray(this.value)
    },
    localValue: {
      get() {
        if (this.multiple) {
          return this.value
        }
        if (this.value.trim().length) {
          return [this.value]
        }
        return []
      },
      set(value) {
        const emitValue = this.multiple ? [...value] : value.toString()
        this.$emit('input', emitValue)
        this.$emit('change', emitValue)
      }
    }
  },
  methods: {
    exactSearchMethod(names) {
      if (Array.isArray(names)) {
        return Promise.resolve(names.map(name => ({ text: name, value: name })))
      }
      return Promise.resolve({
        text: names,
        value: names
      })
    }
  }
}
