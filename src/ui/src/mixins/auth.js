export default {
    methods: {
        $isAuthorized (auth = '') {
            const types = Array.isArray(auth) ? auth : [auth]
            const authorized = types.map(auth => {
                const [ type, action ] = auth.split('.')
                return this.$store.getters['auth/isAuthorized'](type, action)
            })
            return !authorized.some(auth => !auth)
        }
    }
}
