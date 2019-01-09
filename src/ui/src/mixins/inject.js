export default {
    methods: {
        $injectMetadata (params) {
            const injectedParams = this.$tools.clone(params)
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
