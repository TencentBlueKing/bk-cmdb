import http from '@/api'

export const requestIds = {
  findmanyTemplate: Symbol('findmanyTemplate')
}

const find = async (params, config) => {
  try {
    const { count = 0, info: list = [] } = await http.post('findmany/proc/service_template', params, {
      requestId: requestIds.findmanyTemplate,
      ...config
    })
    return { count, list }
  } catch (error) {
    console.error(error)
    return { count: 0, list: [] }
  }
}

const findAll = async (params, config) => {
  try {
    let index = 1
    const size = 1000
    const results = []

    const page = index => ({ start: (index - 1) * size, limit: size })

    const req = index => http.post('findmany/proc/service_template', {
      ...params,
      page: page(index)
    }, config)

    const { count = 0, info: list = [] } = await req(index)
    results.push(...list)

    const max = Math.ceil(count / size)

    const reqs = []
    while (index < max) {
      index += 1
      reqs.push(req(index))
    }

    const rest = await Promise.all(reqs)
    rest.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return results
  } catch (error) {
    console.error(error)
    return []
  }
}

const findAllByIds = async (ids, params, config) => {
  try {
    const size = 1000
    const max = Math.ceil(ids.length / size)

    const req = segment => http.post('findmany/proc/service_template', {
      ...params,
      service_template_ids: segment
    }, config)

    const reqs = []
    for (let index = 1; index <= max; index++) {
      const segment = ids.slice((index - 1) * size, size * index)
      reqs.push(req(segment))
    }

    const results = []
    const res = await Promise.all(reqs)
    res.forEach(({ info: list = [] }) => {
      results.push(...list)
    })

    return results
  } catch (error) {
    console.error(error)
    return []
  }
}

export const CONFIG_MODE = {
  MODULE: 'module',
  TEMPLATE: 'template'
}

export default {
  find,
  findAll,
  findAllByIds
}
