import http from '@/api'
import Vue from 'vue'
import debounce from 'lodash.debounce'
const observer = new Vue({
  data() {
    return {
      queue: [],
      batchValidate: debounce(this.validate, 100)
    }
  },
  watch: {
    queue() {
      this.batchValidate()
    }
  },
  methods: {
    /**
     * { data: { regular: 'xxx', content: 'xxx' }, resolve: function }
     */
    add(item) {
      this.queue.push(item)
    },
    async validate() {
      if (!this.queue.length) return
      const queue = this.queue.splice(0)
      try {
        const results = await http.post(`${window.API_HOST}regular/verify_regular_content_batch`, {
          items: queue.map(({ data }) => data)
        }, {
          globalError: false
        })
        queue.forEach(({ resolve }, index) => resolve({ valid: results[index] }))
      } catch (error) {
        queue.forEach(({ resolve }) => resolve({ valid: false }))
      }
    }
  }
})
export default {
  validate: async (content, { regular }) => new Promise((resolve) => {
    observer.add({
      data: { content, regular },
      resolve
    })
  })
}
