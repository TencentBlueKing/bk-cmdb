const scrollListeners = []
const resizeListeners = []
export function addMainScrollListener (fn) {
    if (scrollListeners.includes(fn)) return
    if (typeof fn === 'function') {
        scrollListeners.push(fn)
    }
}

export function removeMainScrollListener (fn) {
    const index = scrollListeners.indexOf(fn)
    if (index !== -1) {
        scrollListeners.splice(index, 1)
    }
}

export function execMainScrollListener (event) {
    scrollListeners.forEach(fn => {
        fn(event)
    })
}

export function addMainResizeListener (fn) {
    if (resizeListeners.includes(fn)) return
    if (typeof fn === 'function') {
        resizeListeners.push(fn)
    }
}

export function removeMainResizeListener (fn) {
    const index = resizeListeners.indexOf(fn)
    if (index !== -1) {
        resizeListeners.splice(index, 1)
    }
}

export function execMainResizeListener (event) {
    resizeListeners.forEach(fn => {
        fn(event)
    })
}
