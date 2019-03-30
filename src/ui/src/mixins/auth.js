export default {
    methods: {
        $isAuthorized (auth = '', isView = false) {
            const types = Array.isArray(auth) ? auth : [auth]
            const authorized = types.map(auth => {
                const [ type, action ] = auth.split('.')
                return this.$store.getters['auth/isAuthorized'](type, action, isView)
            })
            return !authorized.some(auth => !auth)
        }
    }
}
