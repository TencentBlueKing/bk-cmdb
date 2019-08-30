export default class StatusError {
    constructor (options) {
        this.name = options.name || 'error'
        this.query = options.query || {}
    }
}
