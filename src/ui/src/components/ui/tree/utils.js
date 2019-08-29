export const getNodeId = (data, tree) => {
    const idKey = tree.nodeOptions.idKey
    if (typeof idKey === 'function') {
        return idKey(data)
    }
    return data[idKey]
}

export const getNodeIcon = (data, tree) => {
    const icon = {
        expand: tree.expandIcon,
        collapse: tree.collapseIcon,
        node: tree.nodeIcon
    }
    if (typeof icon.node === 'function') {
        icon.node = icon.node(data)
    }
    return icon
}
