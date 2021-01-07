const types = [{
    id: 'host',
    name: '主机'
}]
export default types

export const formatter = function (id) {
    const type = types.find(type => type.id === id)
    return type ? type.name : id
}
