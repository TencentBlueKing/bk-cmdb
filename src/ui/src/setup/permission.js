import cursor from '@/directives/cursor'

cursor.setOptions({
    globalCallback: options => {
        const permissionModal = window.permissionModal
        permissionModal && permissionModal.show(options.auth)
    }
})
