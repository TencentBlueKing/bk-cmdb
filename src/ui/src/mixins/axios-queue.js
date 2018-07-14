import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters(['axiosQueue'])
    },
    methods: {
        $loading () {
            const queue = this.axiosQueue
            const axiosIds = [].slice.call(arguments)
            if (axiosIds.length) {
                return axiosIds.some(id => queue.includes(id))
            }
            return false
        }
    }
}
