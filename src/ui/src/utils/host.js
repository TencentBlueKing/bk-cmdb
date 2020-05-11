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
