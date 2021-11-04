import autocompleteMixin from './autocomplete'
export default {
  mixins: [autocompleteMixin],
  methods: {
    async fuzzySearchMethod(keyword) {
      try {
        const result = await this.$http.get('biz/simplify?sort=bk_biz_id', {
          requestId: `fuzzy_search_${this.type}`,
          cancelPrevious: true
        })
        const list = result.info || []
        const matchRE = new RegExp(keyword, 'i')
        const matched = []
        list.forEach(({ bk_biz_name: name }) => {
          if (matchRE.test(name)) {
            matched.push({ text: name, value: name })
          }
        })
        return Promise.resolve({
          next: false,
          results: matched
        })
      } catch (error) {
        return Promise.reject(error)
      }
    }
  }
}
