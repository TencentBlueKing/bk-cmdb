export default {
    methods: {
        async handleApplyPermission () {
            try {
                const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
                    params: this.permission,
                    config: {
                        requestId: 'getSkipUrl'
                    }
                })
                window.open(skipUrl)
            } catch (e) {
                console.error(e)
            }
        }
    }
}
