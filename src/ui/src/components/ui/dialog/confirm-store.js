const map = {}
export default {
    install (vm) {
        if (vm.id !== undefined) {
            map[vm.id] = vm
        }
    },
    uninstall (vm) {
        delete map[vm.id]
    },
    async popup (id) {
        const vm = map[id]
        if (vm) {
            vm.show()
            return vm.confirmPromise
        }
        return true
    }
}
