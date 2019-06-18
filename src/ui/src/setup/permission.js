import cursor from '@/directives/cursor'

cursor.setOptions({
    globalCallback: options => {
        const permissionModal = window.permissionModal
        permissionModal && permissionModal.show(options.auth)
    },
    x: 16,
    y: 8
})
