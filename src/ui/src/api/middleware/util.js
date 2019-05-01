export const getFullUrl = url => {
    return `${window.API_PREFIX}/${url}`
}

export const isSameRequest = (origin, config) => {
    const sameUrl = origin.url === config.url
    const sameMethod = origin.method.toLowerCase() === config.method
    const isSame = sameUrl && sameMethod
    if (isSame) {
        config.intercepted = true
    }
    return isSame
}

export const isRedirectResponse = (redirect, { config }) => {
    return config.intercepted
        && getFullUrl(redirect.url) === config.url
        && redirect.method.toLowerCase() === config.method
}
