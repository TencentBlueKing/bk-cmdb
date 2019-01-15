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
            if (mergedOptions.inject) {
                Object.assign(injectedParams, {
                    metadata: {
                        label: {
                            bk_biz_id: String(this.$store.getters['objectBiz/bizId'])
                        }
                    }
                })
            }
            return injectedParams
        }
    }
}
