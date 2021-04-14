import activeMixin from './active'
export default {
  mixins: [activeMixin],
  props: {
    value: {
      type: [String, Array]
    },
    bizId: {
      type: Number,
      required: true
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
    async fuzzySearchMethod(keyword) {
      try {
        const list = await this.$http.post('findmany/object/instances/names', {
          bk_obj_id: this.type,
          bk_biz_id: this.bizId,
          name: keyword
        }, {
          requestId: `fuzzy_search_${this.type}`,
          cancelPrevious: true
        })
        return Promise.resolve({
          next: false,
          results: list.map(name => ({ text: name, value: name }))
        })
      } catch (error) {
        return Promise.reject(error)
      }
    },
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
