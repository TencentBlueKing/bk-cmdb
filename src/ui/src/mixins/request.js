import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters('request', {
            '$requestQueue': 'queue',
            '$requestCache': 'cache'
        })
    },
    methods: {
        $loading (requestIds) {
            if (typeof requestIds === 'undefined') {
                return !!this.$requestQueue.length
            } else if (requestIds instanceof Array) {
                return requestIds.some(requestId => this.$requestQueue.some(request => request.requestId === requestId))
            } else {
                return this.$requestQueue.some(request => request.requestId === requestIds)
            }
        }
    }
}
