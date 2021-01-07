import { mapState } from 'vuex'
export default {
    computed: {
        ...mapState('hostDetails', ['info']),
        HOST_AUTH () {
            if (!this.info) {
                return { U_HOST: null, D_SERVICE_INSTANCE: null }
            }
            const { biz, module, host } = this.info
            // 已分配主机
            if (biz[0].default === 0) {
                const bizId = biz[0].bk_biz_id
                return {
                    U_HOST: {
                        type: this.$OPERATION.U_HOST,
                        relation: [bizId, host.bk_host_id]
                    },
                    D_SERVICE_INSTANCE: {
                        type: this.$OPERATION.D_SERVICE_INSTANCE,
                        relation: [bizId]
                    }
                }
            }
            return {
                U_HOST: {
                    type: this.$OPERATION.U_RESOURCE_HOST,
                    relation: [module[0].bk_module_id, host.bk_host_id]
                },
                D_SERVICE_INSTANCE: null
            }
        }
    }
}
