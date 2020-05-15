export default {
    computed: {
        HOST_AUTH () {
            const params = this.$route.params
            const auth = {
                resource_id: parseInt(params.id)
            }
            const bizId = params.bizId || params.business
            if (bizId) {
                auth.bk_biz_id = parseInt(bizId)
            }
            return {
                U_HOST: {
                    type: this.$OPERATION.U_HOST,
                    ...auth
                },
                D_SERVICE_INSTANCE: {
                    type: this.$OPERATION.D_SERVICE_INSTANCE,
                    ...auth
                }
            }
        }
    }
}
