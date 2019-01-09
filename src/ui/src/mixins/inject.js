export default {
    methods: {
        $injectMetadata (params = {}, shouldClone = false) {
            let injectedParams
            if (shouldClone) {
                injectedParams = this.$tools.clone(params)
            } else {
                injectedParams = params
            }
            if (!this.$store.getters.isAdminView) {
                Object.assign(injectedParams, {
                    metadata: {
                        label: {
                            bk_biz_id: this.$store.getters['objectBiz/bizId']
                        }
                    }
                })
            }
            return injectedParams
        }
    }
}
