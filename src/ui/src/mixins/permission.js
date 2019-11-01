export default {
    methods: {
        $getPermissionList (permission) {
            const getPermissionText = (data, necessaryKey, extraKey, split = '：') => {
                const text = [data[necessaryKey]]
                if (extraKey && data[extraKey]) {
                    text.push(data[extraKey])
                }
                return text.join(split).trim()
            }

            const list = permission.map(datum => {
                const scope = [datum.scope_type_name]
                if (datum.scope_id) {
                    scope.push(datum.scope_name)
                }
                let resource
                if (datum.resource_type_name) {
                    resource = datum.resource_type_name
                } else {
                    resource = datum.resources.map(resource => {
                        const resourceInfo = resource.map(info => getPermissionText(info, 'resource_type_name', 'resource_name'))
                        return [...new Set(resourceInfo)].join('\n')
                    }).join('\n')
                }
                return {
                    scope: getPermissionText(datum, 'scope_type_name', datum.scope_type === 'system' ? null : 'scope_name'),
                    resource: resource,
                    action: datum.action_name
                }
            })

            const uniqueList = []
            list.forEach(item => {
                const exist = uniqueList.some(unique => {
                    return item.resource === unique.resource
                        && item.scope === unique.scope
                        && item.action === unique.action
                })
                if (!exist) {
                    uniqueList.push(item)
                }
            })

            return uniqueList
        },
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
