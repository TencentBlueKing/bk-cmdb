import Form from '@/components/service/form/form.js'
import Bus from '../common/bus'
export default {
    methods: {
        handleAddProcess () {
            Form.show({
                type: 'create',
                title: `${this.$t('添加进程')}(${this.row.name})`,
                hostId: this.row.bk_host_id,
                bizId: this.bizId,
                submitHandler: this.createSubmitHandler
            })
        },
        async createSubmitHandler (values) {
            try {
                await this.$store.dispatch('processInstance/createServiceInstanceProcess', {
                    params: {
                        bk_biz_id: this.bizId,
                        service_instance_id: this.row.id,
                        processes: [{
                            process_info: values
                        }]
                    }
                })
                this.$emit('refresh-count', this.row, this.row.process_count + 1)
                Bus.$emit('refresh-process-list', this.row)
            } catch (error) {
                console.error(error)
            }
        }
    }
}
