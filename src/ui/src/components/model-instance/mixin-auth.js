import { mapGetters } from 'vuex'
import {
  BUILTIN_MODELS,
  BUILTIN_MODEL_ROUTEPARAMS_KEYS
} from '@/dictionary/model-constants.js'

export default {
  computed: {
    ...mapGetters('objectModelClassify', ['models', 'getModelById']),
    INST_AUTH() {
      const { params } = this.$route
      const { bizId } = params
      const { instId } = params
      return {
        U_BUSINESS: {
          type: this.$OPERATION.U_BUSINESS,
          relation: [parseInt(bizId, 10)]
        },
        U_BUSINESS_SET: {
          type: this.$OPERATION.U_BUSINESS_SET,
          relation: [parseInt(params[BUILTIN_MODEL_ROUTEPARAMS_KEYS[BUILTIN_MODELS.BUSINESS_SET]], 10)]
        },
        U_INST: {
          type: this.$OPERATION.U_INST,
          relation: [(this.getModelById(params.objId) || {}).id, parseInt(instId, 10)]
        }
      }
    }
  }
}
