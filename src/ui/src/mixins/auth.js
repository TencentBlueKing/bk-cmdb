export default {
    methods: {
        $isAuthorized (auth = '', option = { isView: false }) {
            const types = Array.isArray(auth) ? auth : [auth]
            const authorized = types.map(auth => {
                const [ type, action ] = auth.split('.')
                return this.$store.getters['auth/isAuthorized'](type, action, option)
            })
            return !authorized.some(auth => !auth)
        }
    }
}
