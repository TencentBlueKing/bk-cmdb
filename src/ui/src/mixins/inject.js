export default {
    methods: {
        $injectMetadata (params = {}, options = {}) {
            const mergedOptions = {
                clone: options.clone || false,
                inject: (options.inject === undefined ? true : options.inject) && !this.$store.getters.isAdminView
            }
            let injectedParams
            if (mergedOptions.clone) {
                injectedParams = this.$tools.clone(params)
            } else {
                injectedParams = params
            }
            const bizId = this.$store.getters['objectBiz/bizId']
            if (mergedOptions.inject && bizId !== null) {
                Object.assign(injectedParams, {
                    metadata: {
                        label: {
                            bk_biz_id: String(bizId)
                        }
                    }
                })
            }
            return injectedParams
        }
    }
}
