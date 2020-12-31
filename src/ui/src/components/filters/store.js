import api from '@/api'
import Vue from 'vue'
import Utils from './utils'
import store from '@/store'
import i18n from '@/i18n'
import QS from 'qs'
import RouterQuery from '@/router/query'
import Throttle from 'lodash.throttle'

function getStorageHeader (type, key, properties) {
    if (!key) {
        return []
    }
    const data = store.state.userCustom[type][key] || []
    const header = []
    data.forEach(propertyId => {
        const property = properties.find(property => property.bk_property_id === propertyId)
        property && header.push(property)
    })
    return header
}

const FilterStore = new Vue({
    data () {
        return {
            config: {},
            components: {},
            properties: [],
            propertyGroups: [],
            selected: [],
            IP: Utils.getDefaultIP(),
            condition: {},
            header: [],
            collections: [],
            activeCollection: null,
            throttleSearch: Throttle(this.dispatchSearch, 100, { leading: false })
        }
    },
    computed: {
        bizId () {
            return this.config.bk_biz_id || void 0
        },
        userBehaviorKey () {
            return this.bizId ? 'topology_host_common_filter' : 'resource_host_common_filter'
        },
        collectable () {
            return !!this.bizId
        },
        request () {
            return {
                property: Symbol(this.bizId || 'property'),
                propertyGroup: Symbol(this.bizId || 'propertyGroup'),
                collections: Symbol(this.bizId || 'collections'),
                deleteCollection: id => `delete_collection_${id}`,
                updateCollection: id => `update_collection_${id}`
            }
        },
        searchHandler () {
            return this.config.searchHandler || (() => {})
        },
        globalHeader () {
            const key = this.config.header && this.config.header.global
            return getStorageHeader('globalUsercustom', key, this.modelPropertyMap.host)
        },
        customHeader () {
            const key = this.config.header && this.config.header.custom
            return getStorageHeader('usercustom', key, this.modelPropertyMap.host)
        },
        presetHeader () {
            const hostProperties = this.getModelProperties('host')
            const headerProperties = Utils.getInitialProperties(hostProperties)
            const fixed = ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
            const fixedProperties = fixed.map(propertyId => Utils.findPropertyByPropertyId(propertyId, hostProperties))
            if (!this.bizId) {
                const topologyProperty = Utils.findPropertyByPropertyId('__bk_host_topology__', hostProperties)
                fixedProperties.push(topologyProperty)
            } else {
                const moduleNameProperty = Utils.findPropertyByPropertyId('bk_module_name', this.properties, 'module')
                const setNameProperty = Utils.findPropertyByPropertyId('bk_set_name', this.properties, 'set')
                fixedProperties.push(moduleNameProperty, setNameProperty)
            }
            return Utils.getUniqueProperties(fixedProperties, headerProperties).slice(0, 6)
        },
        defaultHeader () {
            return this.customHeader.length
                ? this.customHeader
                : this.globalHeader.length
                    ? this.globalHeader
                    : this.presetHeader
        },
        modelPropertyMap () {
            const map = {}
            this.properties.forEach(property => {
                const modelId = property.bk_obj_id
                const modelProperties = map[modelId] || []
                modelProperties.push(property)
                map[modelId] = modelProperties
            })
            return map
        },
        hasCondition () {
            return Object.keys(this.condition).some(id => {
                const value = this.condition[id].value
                return !!String(value).trim().length
            })
        }
    },
    watch: {
        selected () {
            this.initCondition()
        }
    },
    methods: {
        setupPropertyQuery () {
            const query = QS.parse(RouterQuery.get('filter'))
            const properties = []
            const condition = {}
            try {
                Object.keys(query).forEach(key => {
                    const [id, operator] = key.split('.')
                    const property = Utils.findProperty(id, this.properties)
                    const value = query[key].toString().split(',')
                    if (property && operator && value.length) {
                        properties.push(property)
                        condition[property.id] = {
                            operator: '$' + operator,
                            value: Utils.convertValue(value, '$' + operator, property)
                        }
                    }
                })
            } catch (error) {
                Vue.prototype.$warn(i18n.t('解析查询链接出错提示'))
            }
            this.selected = properties
            this.condition = condition
        },
        setupNormalProperty () {
            const userBehavior = store.state.userCustom.usercustom[this.userBehaviorKey] || []
            const normal = userBehavior.length ? userBehavior : [
                ['bk_set_name', 'set'],
                ['bk_module_name', 'module'],
                ['operator', 'host'],
                ['bk_bak_operator', 'host'],
                ['bk_cloud_id', 'host']
            ]
            this.createOrUpdateCondition(normal.map(([field, model]) => ({ field, model })), { createOnly: true, useDefaultData: true })
        },
        setupIPQuery () {
            const query = QS.parse(RouterQuery.get('ip'))
            const { text = '', exact = 'true', inner = 'true', outer = 'true' } = query
            this.IP = {
                text: text.replace(/,/g, '\n'),
                exact: exact.toString() === 'true',
                inner: inner.toString() === 'true',
                outer: outer.toString() === 'true'
            }
        },
        setupHeader () {
            (this.selected.length && this.hasCondition) ? this.setHeader(this.selected) : this.setHeader(this.defaultHeader)
        },
        initCondition () {
            const newConditon = {}
            this.selected.forEach(property => {
                if (this.condition.hasOwnProperty(property.id)) {
                    newConditon[property.id] = this.condition[property.id]
                } else {
                    newConditon[property.id] = Utils.getDefaultData(property)
                }
            })
            this.condition = newConditon
        },
        setCondition (data = {}) {
            this.condition = data.condition || this.condition
            this.IP = data.IP || this.IP
            this.setHeader(this.selected)
            this.throttleSearch()
        },
        updateCondition (property, operator, value) {
            this.condition[property.id] = {
                operator,
                value
            }
            this.throttleSearch()
        },
        createOrUpdateCondition (data, options = {}) {
            const { createOnly = false, useDefaultData = false } = options
            data.forEach(({ field, model, operator, value }) => {
                const existProperty = this.selected.find(property => property.bk_property_id === field && property.bk_obj_id === model)
                if (!existProperty) {
                    const property = Utils.findPropertyByPropertyId(field, this.getModelProperties(model))
                    if (property) {
                        const defaultData = Utils.getDefaultData(property)
                        this.selected.push(property)
                        this.$set(this.condition, property.id, {
                            operator: useDefaultData ? defaultData.operator : operator,
                            value: useDefaultData ? defaultData.value : value
                        })
                    }
                } else if (!createOnly) {
                    const defaultData = Utils.getDefaultData(existProperty)
                    const condition = this.condition[existProperty.id]
                    condition.operator = useDefaultData ? defaultData.operator : operator
                    condition.value = useDefaultData ? defaultData.value : value
                }
            })
            this.throttleSearch()
        },
        updateIP (data = {}) {
            Object.assign(this.IP, data)
            this.throttleSearch()
        },
        updateSelected (selected) {
            this.selected = selected
        },
        removeSelected (property) {
            const index = this.selected.findIndex(target => target.id.toString() === property.id.toString())
            index > -1 && this.selected.splice(index, 1)
            this.throttleSearch()
        },
        resetValue (property, silent = false) {
            const properties = Array.isArray(property) ? property : [property]
            properties.forEach(target => {
                const operator = this.condition[target.id].operator
                const value = Utils.getOperatorSideEffect(target, operator, [])
                this.condition[target.id].value = value
            })
            !silent && this.throttleSearch()
        },
        resetAll (silent) {
            this.IP = Utils.getDefaultIP()
            this.resetValue(this.selected, silent)
        },
        resetIP () {
            this.IP = Utils.getDefaultIP()
            this.throttleSearch()
        },
        dispatchSearch () {
            this.setQuery()
            this.searchHandler(this.condition)
        },
        setQuery () {
            const query = {}
            Object.keys(this.condition).forEach(id => {
                const { operator, value } = this.condition[id]
                if (String(value).length) {
                    query[`${id}.${operator.replace('$', '')}`] = Array.isArray(value) ? value.join(',') : value
                }
            })
            RouterQuery.set({
                filter: QS.stringify(query, { encode: false }),
                ip: QS.stringify(this.IP.text.trim().length ? this.IP : {}, { encode: false }),
                _t: Date.now()
            })
        },
        setHeader (newHeader) {
            newHeader = newHeader.length ? newHeader : this.defaultHeader
            const presetId = ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
            const presetProperty = presetId.map(propertyId => {
                return this.properties.find(property => property.bk_property_id === propertyId)
            })
            this.header = Utils.getUniqueProperties(presetProperty, newHeader)
        },
        setActiveCollection (collection, silent) {
            if (!collection) {
                this.activeCollection = null
                this.resetAll(silent)
                return
            }
            try {
                const IP = JSON.parse(collection.info)
                if (IP.hasOwnProperty('ip_list')) { // 因老数据的操作符不可兼容，应用收藏条件时直接提示错误并返回
                    this.$error(i18n.t('应用收藏条件失败提示'))
                    return false
                }
                const queryParams = JSON.parse(collection.query_params)
                const condition = {}
                const selected = []
                let hasMissedProperty = false
                queryParams.forEach(query => {
                    const property = this.properties.find(property => {
                        return property.bk_obj_id === query.bk_obj_id && property.bk_property_id === query.field
                    })
                    if (property) {
                        selected.push(property)
                        condition[property.id] = {
                            operator: query.operator,
                            value: query.value
                        }
                    } else {
                        hasMissedProperty = true
                    }
                })
                this.updateSelected(selected)
                this.setCondition({ IP, condition })
                this.activeCollection = collection
                hasMissedProperty && this.$warn(i18n.t('收藏条件未完全解析提示'))
            } catch (error) {
                console.error(error)
                this.$error(i18n.t('应用收藏条件失败提示'))
            }
        },
        getSearchParams () {
            const transformedIP = Utils.transformIP(this.IP.text)
            const flag = []
            this.IP.inner && flag.push('bk_host_innerip')
            this.IP.outer && flag.push('bk_host_outerip')
            const params = {
                bk_biz_id: this.bizId, // undefined会被忽略
                ip: {
                    data: transformedIP.data,
                    exact: this.IP.exact ? 1 : 0,
                    flag: flag.join('|')
                },
                condition: Utils.transformCondition(this.condition, this.selected, this.header)
            }
            if (transformedIP.condition) {
                const { condition } = params.condition.find(condition => condition.bk_obj_id === 'host')
                condition.push(transformedIP.condition)
            }
            return params
        },
        getModelProperties (modelId) {
            return [...this.modelPropertyMap[modelId] || []]
        },
        getProperty (propertyId, modelId) {
            return Utils.findPropertyByPropertyId(propertyId, this.getModelProperties(modelId))
        },
        setComponent (name, component) {
            this.components[name] = component
        },
        getComponent (name) {
            return this.components[name]
        },
        async getProperties () {
            const properties = await api.post('find/objectattr', {
                bk_biz_id: this.bizId,
                bk_obj_id: {
                    $in: ['host', 'module', 'set', 'biz']
                }
            }, {
                requestId: this.request.property,
                fromCache: true
            })
            const hostIdProperty = Utils.defineProperty({
                id: 'bk_host_id',
                bk_obj_id: 'host',
                bk_property_id: 'bk_host_id',
                bk_property_name: 'ID',
                bk_property_index: -Infinity,
                bk_property_type: 'int'
            })
            const serviceTemplateProperty = Utils.defineProperty({
                id: 'service_template_id',
                bk_obj_id: 'module',
                bk_property_id: 'service_template_id',
                bk_property_name: i18n.t('服务模板'),
                bk_property_index: Infinity,
                bk_property_type: 'service-template'
            })
            if (this.bizId) {
                this.properties = [...properties, hostIdProperty, serviceTemplateProperty]
            } else {
                const topologyProperty = Utils.defineProperty({
                    id: Date.now() + 2,
                    bk_obj_id: 'host',
                    bk_property_id: '__bk_host_topology__',
                    bk_property_index: Infinity,
                    bk_property_name: i18n.t('业务拓扑'),
                    bk_property_type: 'topology',
                    required: true
                })
                this.properties = [...properties, hostIdProperty, serviceTemplateProperty, topologyProperty]
            }
            return this.properties
        },
        async getPropertyGroups () {
            const groups = await api.post('find/objectattgroup/object/host', {
                bk_biz_id: this.bizId
            }, {
                requestId: this.request.propertyGroup
            })
            this.propertyGroups = groups
            return groups
        },
        async loadCollections () {
            const { info: collections } = await api.post('hosts/favorites/search', {
                condition: {
                    bk_biz_id: this.bizId
                }
            }, {
                requestId: this.request.collections
            })
            this.collections = collections
            return this.collections
        },
        async createCollection (data) {
            const response = await api.post('hosts/favorites', data, {
                requestId: this.request.createCollection,
                globalError: false,
                transformData: false
            })
            if (response.result) {
                const collection = { id: response.data.id, ...data }
                this.collections.unshift(collection)
                this.activeCollection = collection
                this.dispatchSearch()
                return Promise.resolve()
            }
            return Promise.reject(response)
        },
        async removeCollection ({ id }) {
            await api.delete(`hosts/favorites/${id}`, {
                requestId: this.request.deleteCollection(id)
            })
            const index = this.collections.findIndex(target => target.id === id)
            index > -1 && this.collections.splice(index, 1)
            return Promise.resolve()
        },
        async updateCollection (collection) {
            await api.put(`hosts/favorites/${collection.id}`, collection, {
                requestId: this.request.updateCollection(collection.id)
            })
            const raw = this.collections.find(raw => raw.id === collection.id)
            Object.assign(raw, collection)
            return Promise.resolve()
        },
        async updateUserBehavior (properties) {
            await store.dispatch('userCustom/saveUsercustom', {
                [this.userBehaviorKey]: properties.map(property => [property.bk_property_id, property.bk_obj_id])
            })
            return Promise.resolve()
        }
    }
})

/*
* config.bk_biz_id 业务id，仅业务拓扑需要
* config.searchHandler 触发搜索的方法
* config.header.custom 存储用户自定义的表头字段key
* config.header.global 存储全局默认表头字段key
*/

export async function setupFilterStore (config = {}) {
    FilterStore.config = config
    FilterStore.selected = []
    FilterStore.condition = {}
    FilterStore.components = {}
    FilterStore.activeCollection = null
    await Promise.all([
        FilterStore.getProperties(),
        FilterStore.getPropertyGroups()
    ])
    FilterStore.setupIPQuery()
    FilterStore.setupPropertyQuery()
    FilterStore.setupNormalProperty()
    FilterStore.setupHeader()
    FilterStore.throttleSearch()
    return FilterStore
}

window.FilterStore = FilterStore

export default FilterStore
