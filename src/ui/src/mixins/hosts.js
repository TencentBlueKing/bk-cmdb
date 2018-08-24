import {mapGetters, mapActions} from 'vuex'
import cmdbHostsFilter from '@/components/hosts-filter/hosts-filter'
import cmdbColumnsConfig from '@/components/columns-config/columns-config'
export default {
    components: {
        cmdbHostsFilter,
        cmdbColumnsConfig
    },
    data () {
        return {
            properties: {
                biz: [],
                host: [],
                set: [],
                module: []
            },
            propertyGroups: [],
            table: {
                checked: [],
                header: [],
                list: [],
                allList: [],
                pagination: {
                    current: 1,
                    size: 10,
                    count: 0
                },
                defaultSort: 'bk_host_id',
                sort: 'bk_host_id',
                exportUrl: `${window.Site.url}hosts/export`
            },
            filter: {
                params: null,
                paramsResolver: null
            },
            slider: {
                show: false,
                title: ''
            },
            tab: {
                active: 'attribute',
                attribute: {
                    type: 'details',
                    inst: {
                        details: {},
                        edit: {}
                    }
                }
            },
            columnsConfig: {
                show: false,
                selected: []
            }
        }
    },
    computed: {
        ...mapGetters(['supplierAccount']),
        ...mapGetters('userCustom', ['usercustom']),
        customColumns () {
            return this.usercustom[this.table.columnsConfigKey] || []
        },
        clipboardList () {
            return this.table.header.filter(header => header.type !== 'checkbox')
        }
    },
    methods: {
        ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
        ...mapActions('objectModelFieldGroup', ['searchGroup']),
        ...mapActions('hostUpdate', ['updateHost']),
        getProperties () {
            return this.batchSearchObjectAttribute({
                params: {
                    bk_obj_id: {'$in': Object.keys(this.properties)},
                    bk_supplier_account: this.supplierAccount
                },
                config: {
                    requestId: 'hostsAttribute',
                    fromCache: true
                }
            }).then(result => {
                Object.keys(this.properties).forEach(objId => {
                    this.properties[objId] = result[objId]
                })
                return result
            })
        },
        getHostPropertyGroups () {
            return this.searchGroup({
                objId: 'host',
                config: {
                    fromCache: true,
                    requestId: 'hostAttributeGroup'
                }
            }).then(groups => {
                this.propertyGroups = groups
                return groups
            })
        },
        setAllHostList (list) {
            if (this.table.allList.length === this.table.pagination.count) return
            const newList = []
            list.forEach(item => {
                const exist = this.table.allList.some(existItem => existItem['host']['bk_host_id'] === item['host']['bk_host_id'])
                if (!exist) {
                    newList.push(item)
                }
            })
            this.table.allList = [...this.table.allList, ...newList]
        },
        getHostCellText (header, item) {
            const objId = header.objId
            const propertyId = header.id
            const headerProperty = this.$tools.getProperty(this.properties[objId], propertyId)
            const flatternedItem = this.$tools.flatternHostItem(headerProperty, item)
            const text = this.$tools.getHostCellText(flatternedItem, objId, propertyId)
            return text
        },
        handlePageChange (current) {
            this.table.pagination.current = current
            this.getHostList()
        },
        handleSizeChange (size) {
            this.table.pagination.size = size
            this.handlePageChange(1)
        },
        handleSortChange (sort) {
            this.table.sort = sort
            this.handlePageChange(1)
        },
        handleCopy (target) {
            const copyList = this.table.allList.filter(item => {
                return this.table.checked.includes(item['host']['bk_host_id'])
            })
            const copyText = []
            this.$tools.clone(copyList).forEach(item => {
                const cellText = this.getHostCellText(target, item)
                if (cellText !== '--') {
                    copyText.push(cellText)
                }
            })
            if (copyText.length) {
                this.$copyText(copyText.join('\n')).then(() => {
                    this.$success(this.$t('Common["复制成功"]'))
                }, () => {
                    this.$error(this.$t('Common["复制失败"]'))
                })
            } else {
                this.$info(this.$t('Common["该字段无可复制的值"]'))
            }
        },
        handleRefresh (params) {
            this.filter.params = params
            if (this.filter.paramsResolver) {
                this.filter.paramsResolver()
            } else {
                this.getHostList()
            }
        },
        async handleCheckAll (type) {
            let list
            if (type === 'current') {
                list = this.table.list
            } else {
                const data = await this.getAllHostList()
                list = data.info
            }
            this.table.checked = list.map(item => item['host']['bk_host_id'])
        },
        handleRowClick (item) {
            const inst = item['host']
            this.slider.show = true
            this.slider.title = `${this.$t("Common['编辑']")} ${inst['bk_host_innerip']}`
            this.tab.attribute.inst.details = inst
            this.tab.attribute.type = 'details'
        },
        handleSave (values, changedValues, inst, type) {
            this.batchUpdate({
                ...changedValues,
                'bk_host_id': inst['bk_host_id'].toString()
            })
        },
        batchUpdate (params) {
            this.updateHost(params).then(() => {
                this.$success(this.$t('Common[\'保存成功\']'))
                this.getHostList()
                this.slider.show = false
            })
        },
        handleCancel () {
            this.tab.attribute.type = 'details'
        },
        async handleEdit (flatternItem) {
            const list = await this.$http.cache.get('hostSearch')
            const originalItem = list.info.find(item => item['host']['bk_host_id'] === flatternItem['bk_host_id'])
            this.tab.attribute.inst.edit = originalItem['host']
            this.tab.attribute.type = 'update'
        },
        handleMultipleEdit () {
            this.tab.attribute.type = 'multiple'
            this.slider.title = this.$t('HostResourcePool[\'主机属性\']')
            this.slider.show = true
        },
        handleMultipleSave (changedValues) {
            this.batchUpdate({
                ...changedValues,
                'bk_host_id': this.table.checked.join(',')
            })
        },
        handleMultipleCancel () {
            this.slider.show = false
        },
        handleApplyColumnsConfig (properties) {
            this.$store.dispatch('userCustom/saveUsercustom', {
                [this.table.columnsConfigKey]: properties.map(property => property['bk_property_id'])
            })
            this.columnsConfig.show = false
        },
        handleResetColumnsConfig () {
            this.$store.dispatch('userCustom/saveUsercustom', {
                [this.table.columnsConfigKey]: []
            })
            this.columnsConfig.show = false
        }
    }
}
