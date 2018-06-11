import { $AxiosQueue } from '@/api/axios'

export default {
    computed: {
        $AxiosQueue () {
            return $AxiosQueue
        }
    },
    methods: {
        $loading () {
            const queue = this.$AxiosQueue
            const axiosIds = [].slice.call(arguments)
            if (axiosIds.length) {
                return axiosIds.some(id => queue.includes(id))
            }
            return false
        }
    }
}
