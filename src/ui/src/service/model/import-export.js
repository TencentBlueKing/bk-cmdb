import http from '@/api'
const importService = async (objId, event) => {
  try {
    const [file] = event.target.files
    const form = new FormData()
    form.append('file', file)
    return http.post(`${window.API_HOST}object/object/${objId}/import`, form)
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const exportService = objId => http.download({ url: `${window.API_HOST}object/object/${objId}/export`, data: {} })
export default {
  import: importService,
  export: exportService
}
