export default class Meta {
    constructor (data = {}) {
        this.owner = ''
        Object.keys(data).forEach(key => {
            this[key] = data[key]
        })

        this.menu = Object.assign({
            i18n: '',
            parent: null,
            relative: null
        }, data.menu)

        this.auth = Object.assign({
            authScope: 'global',
            view: null,
            operation: []
        }, data.auth)
    }
}
