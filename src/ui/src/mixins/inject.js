export default {
    methods: {
        $injectMetadata (params = {}, options = {}) {
            const mergedOptions = {
                clone: options.clone || false,
                inject: (options.inject === undefined ? true : options.inject) && !this.$store.getters.isAdminView,
                injectBizId: options.injectBizId && !this.$store.getters.isAdminView
            }
            let injectedParams
            if (mergedOptions.clone) {
                injectedParams = this.$tools.clone(params)
            } else {
                injectedParams = params
            }
            const bizId = this.$store.getters['objectBiz/bizId']
            if (mergedOptions.inject && bizId !== null && !mergedOptions.injectBizId) {
                Object.assign(injectedParams, {
                    metadata: {
                        label: {
                            bk_biz_id: String(bizId)
                        }
                    }
                })
            } else if (bizId !== null && mergedOptions.injectBizId) {
                Object.assign(injectedParams, {
                    bk_biz_id: Number(bizId)
                })
            }
            return injectedParams
        }
    }
}
