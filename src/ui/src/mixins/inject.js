export default {
    methods: {
        $injectMetadata (params, shouldClone = true) {
            let injectedParams
            if (shouldClone) {
                injectedParams = this.$tools.clone(params)
            } else {
                injectedParams = params
            }
            if (!this.$store.getters.isAdminView) {
                Object.assign(injectedParams, {
                    metadata: {
                        label: this.$store.getters['objectBiz/bizId']
                    }
                })
            }
            return injectedParams
        }
    }
}
