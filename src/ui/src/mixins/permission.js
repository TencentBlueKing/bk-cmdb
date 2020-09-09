export default {
    methods: {
        handleApplyPermission () {
            return this.$store.dispatch('auth/getSkipUrl', {
                params: this.permission,
                config: {
                    requestId: 'getSkipUrl'
                }
            }).then(url => {
                window.open(url)
                return url
            })
        }
    }
}
