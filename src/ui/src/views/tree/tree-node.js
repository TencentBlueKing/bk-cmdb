export default class TreeNode {
    constructor (data, options) {
        this.options = options
        this.data = data

        this.id = this.getNodeId()
        this.name = data[options.nameKey]
        this.level = options.level
        this.index = null
        this.parent = options.parent

        this.expanded = false
        this.selected = false
    }

    getNodeId () {
        const idKey = this.options.idKey
        if (typeof idKey === 'function') {
            return idKey(this.data)
        }
        return this.data[idKey]
    }

    setValue (key, value) {
        if (this.hasOwnProperty(key)) {
            this[key] = value
        }
    }
}
