const vendors = [{
    id: '1',
    name: '亚马逊云'
}, {
    id: '2',
    name: '腾讯云'
}]

export default vendors

export const formatter = function (id) {
    const vendor = vendors.find(vendor => vendor.id === id)
    return vendor ? vendor.name : id
}
