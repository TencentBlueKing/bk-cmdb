import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters(['user']),
        ...mapGetters('userPrivilege', ['privilege']),
        ...mapGetters('objectModelClassify', ['authorizedNavigation']),
        $authorized () {
            let fullAuthority = ['search', 'update', 'delete']
            let fullAuthorityClassification = ['bk_host_manage', 'bk_back_config', 'bk_index']
            let authorized = []
            if (this.user.admin === '1') {
                authorized = fullAuthority
            } else {
                let modelAuthority = this.privilege['model_config']
                let model = this.$model // $model in mixins of classify
                if (model) {
                    const classificationId = model['bk_classification_id']
                    if (fullAuthorityClassification.includes(classificationId)) {
                        authorized = fullAuthority
                    } else {
                        authorized = modelAuthority.hasOwnProperty(classificationId) ? modelAuthority[classificationId][model['bk_obj_id']] : []
                    }
                }
            }
            return {
                search: authorized.includes('search'),
                update: authorized.includes('update'),
                delete: authorized.includes('delete')
            }
        }
    },
    methods: {
        $getUnauthorized (option) {
            let options = []
            if (typeof option === 'string') {
                options.push(option)
            } else if (Array.isArray(option)) {
                options = [...option]
            } else {
                options = ['search', 'update', 'delete']
            }
            return options.some(optionType => !this.authorized[optionType])
        }
    }
}
