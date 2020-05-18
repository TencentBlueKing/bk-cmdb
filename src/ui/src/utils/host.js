import RouterQuery from '@/router/query'
export const getIPPayload = function (inputPayload = {}) {
    const ip = RouterQuery.get('ip', '')
    const queryPayload = {
        data: ip.length ? ip.split(',') : [],
        inner: parseInt(RouterQuery.get('inner', 1)) === 1,
        outer: parseInt(RouterQuery.get('outer', 1)) === 1,
        exact: parseInt(RouterQuery.get('exact', 0)) === 1
    }

    const payload = { ...queryPayload, ...inputPayload }

    const flag = []
    if (payload.inner) {
        flag.push('bk_host_innerip')
    }
    if (payload.outer) {
        flag.push('bk_host_outerip')
    }

    return {
        data: payload.data,
        flag: flag.join('|'),
        exact: payload.exact ? 1 : 0
    }
}

export function injectFields (params, tableHeaderList = []) {
    const headerFields = {}
    const fillFields = {
        host: [],
        set: ['bk_set_id'],
        biz: ['bk_biz_id'],
        module: ['bk_module_id']
    }

    tableHeaderList.forEach(header => {
        const objId = header.bk_obj_id || header.objId
        const propertyId = header.bk_property_id || header.id
        if (headerFields[objId]) {
            headerFields[objId].push(propertyId)
        } else {
            headerFields[objId] = [propertyId]
        }
    })

    Object.keys(headerFields).forEach(objId => {
        headerFields[objId] = [...headerFields[objId], ...fillFields[objId] || []]
    })

    params.condition.forEach(condition => {
        condition.fields = Array.from(new Set([...condition.fields, ...headerFields[condition.bk_obj_id] || []]))
    })

    return params
}
