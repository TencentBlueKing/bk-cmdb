export default class Meta {
    constructor (data = {}) {
        let menu = false
        if (data.menu) {
            menu = Object.assign({
                id: null,
                i18n: null,
                path: null,
                order: 1,
                parent: null,
                adminView: true,
                businessView: true
            }, data.menu)
        }

        const auth = Object.assign({
            view: null,
            operation: []
        }, data.auth)

        return {
            title: '',
            i18nTitle: '',
            ...data,
            menu,
            auth
        }
    }
}
