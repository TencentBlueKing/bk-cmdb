import autocompleteMixin from './autocomplete'
export default {
  mixins: [autocompleteMixin],
  props: {
    bizId: {
      type: Number,
      required: true
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
    }
  }
}
