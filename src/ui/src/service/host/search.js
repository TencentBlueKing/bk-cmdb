import http from '@/api'
const findOne = async ({ bk_host_id: hostId, bk_biz_id: bizId, config }) => {
  try {
    const { info } = await http.post('hosts/search', {
      bk_biz_id: bizId || -1,
      condition: [
        { bk_obj_id: 'biz', condition: [], fields: [] },
        { bk_obj_id: 'set', condition: [], fields: [] },
        { bk_obj_id: 'module', condition: [], fields: [] },
        { bk_obj_id: 'host', condition: [{
          field: 'bk_host_id',
          operator: '$eq',
          value: hostId
        }], fields: [] }
      ],
      id: { flag: 'bk_host_innerip', exact: 1, data: [] }
    }, config)
    const [instance] = info
    return instance ? instance.host : null
  } catch (error) {
    console.error(error)
    return null
  }
}

export default {
  findOne
}
