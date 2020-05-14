import { mapGetters } from 'vuex'

export default {
    computed: {
        ...mapGetters('objectModelClassify', ['models', 'getModelById']),
        INST_AUTH () {
            const params = this.$route.params
            const bizId = params.bizId || params.business
            const instId = params.instId
            const objId = params.objId
            const auth = {}
            if (objId) {
                const model = this.getModelById(objId) || {}
                auth.parent_layers = [{
                    resource_id: model.id,
                    resource_type: 'model'
                }]
            }

            return {
                U_BUSINESS: {
                    type: this.$OPERATION.U_BUSINESS,
                    resource_id: parseInt(bizId),
                    ...auth
                },
                U_INST: {
                    type: this.$OPERATION.U_INST,
                    resource_id: parseInt(instId),
                    ...auth
                }
            }
        }
    }
}
