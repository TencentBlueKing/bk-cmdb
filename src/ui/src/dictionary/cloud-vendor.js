const vendors = [{
    id: 'tencent_cloud',
    name: '腾讯云'
}, {
    id: 'ali_cloud',
    name: '阿里云'
}, {
    id: 'aws_cloud',
    name: '亚马逊云'
}]

export default vendors

export const formatter = function (id) {
    const vendor = vendors.find(vendor => vendor.id === id)
    return vendor ? vendor.name : id
}
