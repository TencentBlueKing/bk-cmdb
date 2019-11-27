import Bus from '@/utils/bus.js'
export default {
    props: {
        properties: {
            type: Array,
            required: true
        },
        selection: {
            type: Array,
            required: true
        }
    },
    data () {
        return {
            slider: {
                show: false,
                title: '',
                component: null,
                props: {}
            }
        }
    },
    methods: {
        async handleMultipleEdit () {
            try {
                const [objectUnique, groups] = await Promise.all([
                    this.getObjectUnique(),
                    this.getPropertyGroups()
                ])
                this.slider.props.objectUnique = objectUnique
                this.slider.props.propertyGroups = groups
                this.slider.props.properties = this.properties
                this.slider.title = this.$t('主机属性')
                this.slider.component = 'cmdb-form-multiple'
                this.slider.show = true
            } catch (e) {
                console.error(e)
            }
        },
        async handleMultipleSave (changedValues) {
            try {
                await this.$store.dispatch('hostUpdate/updateHost', {
                    params: {
                        ...changedValues,
                        'bk_host_id': this.selection.map(row => row.host.bk_host_id).join(',')
                    }
                })
                this.slider.show = false
                Bus.$emit('refresh-list', 1)
            } catch (e) {
                console.error(e)
            }
        },
        handleSliderBeforeClose () {
            const $form = this.$refs.multipleForm
            if (Object.keys($form.changedValues).length) {
                return new Promise((resolve, reject) => {
                    this.$bkInfo({
                        title: this.$t('确认退出'),
                        subTitle: this.$t('退出会导致未保存信息丢失'),
                        extCls: 'bk-dialog-sub-header-center',
                        confirmFn: () => {
                            this.slider.show = false
                            resolve(true)
                        },
                        cancelFn: () => {
                            resolve(false)
                        }
                    })
                })
            } else {
                this.slider.show = false
            }
        },
        getObjectUnique () {
            return this.$store.dispatch('objectUnique/searchObjectUniqueConstraints', {
                objId: 'host',
                params: {}
            })
        },
        getPropertyGroups () {
            return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                objId: 'host',
                params: this.$injectMetadata()
            })
        }
    }
}
