import http from '@/api'
export default {
  validate: async (regular) => {
    try {
      if (!regular.trim().length) {
        return { valid: true }
      }
      const data = await http.post(`${window.API_HOST}regular/verify_regular_express`, {
        regular
      }, {
        globalError: false
      })
      return { valid: data.is_valid }
    } catch (error) {
      return { valid: false }
    }
  }
}
