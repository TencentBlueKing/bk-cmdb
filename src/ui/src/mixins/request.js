import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters('request', {
      $requestQueue: 'queue',
      $requestCache: 'cache'
    })
  },
  methods: {
    $loading(requestIds) {
      if (typeof requestIds === 'undefined') {
        return !!this.$requestQueue.length
      } if (requestIds instanceof Array) {
        return requestIds.some(requestId => this.$requestQueue.some(request => request.requestId === requestId))
      } if (typeof requestIds === 'string' && requestIds.startsWith('^=')) {
        // eslint-disable-next-line prefer-destructuring
        const requestId = requestIds.split('^=')[1]
        const matchIndex = this.$requestQueue.findIndex(request => (typeof request.requestId === 'string') && request.requestId.startsWith(requestId))
        return matchIndex !== -1
      }
      return this.$requestQueue.some(request => request.requestId === requestIds)
    }
  }
}
