import { $AxiosQueue } from '@/api/axios'

export default {
    methods: {
        $loading () {
            const queue = $AxiosQueue
            const axiosIds = [].slice.call(arguments)
            if (axiosIds.length) {
                return axiosIds.some(id => queue.includes(id))
            }
            return false
        }
    }
}
