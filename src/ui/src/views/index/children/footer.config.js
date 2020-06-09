const data = [{
    name: 'QQ咨询(800802001)',
    en_name: 'QQ(800802001)',
    href: 'http://wpa.b.qq.com/cgi/wpa.php?ln=1&key=XzgwMDgwMjAwMV80NDMwOTZfODAwODAyMDAxXzJf'
}, {
    name: '蓝鲸论坛',
    en_name: 'Blueking Forum',
    href: 'https://bk.tencent.com/s-mart/community/'
}, {
    name: '蓝鲸官网',
    en_name: 'BlueKing Official',
    href: 'https://bk.tencent.com/index/'
}]
if (window.CMDB_CONFIG.site.bkDesktop) {
    data.push({
        name: '蓝鲸桌面',
        en_name: 'Blueking Desktop',
        href: window.CMDB_CONFIG.site.bkDesktop
    })
}
export default Object.freeze(data)
