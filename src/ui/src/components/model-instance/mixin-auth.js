import { mapGetters } from 'vuex'

export default {
    computed: {
        ...mapGetters('objectModelClassify', ['models', 'getModelById']),
        INST_AUTH () {
            const params = this.$route.params
            const bizId = params.bizId
            const instId = params.instId
            return {
                U_BUSINESS: {
                    type: this.$OPERATION.U_BUSINESS,
                    relation: [parseInt(bizId)]
                },
                U_INST: {
                    type: this.$OPERATION.U_INST,
                    relation: [(this.getModelById(params.objId) || {}).id, parseInt(instId)]
                }
            }
        }
    }
}
