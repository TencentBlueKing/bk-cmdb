import RouterQuery from '@/router/query'
export const getIPPayload = function () {
    const ip = RouterQuery.get('ip', '')
    const data = []
    if (ip.length) {
        data.push(...ip.split(','))
    }
    const flag = []
    const inner = parseInt(RouterQuery.get('inner', 1)) === 1
    const outer = parseInt(RouterQuery.get('outer', 1)) === 1
    if (inner) {
        flag.push('bk_host_innerip')
    }
    if (outer) {
        flag.push('bk_host_outerip')
    }
    const exact = parseInt(RouterQuery.get('exact', '1')) === 1 ? 1 : 0
    return {
        data: data,
        flag: flag.join('|'),
        exact: exact
    }
}
