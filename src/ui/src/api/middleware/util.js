export const getFullUrl = url => {
    return `${window.API_PREFIX}/${url}`
}

export const isSameRequest = (origin, config) => {
    return origin.url === config.url
        && origin.method.toLowerCase() === config.method
}

export const isRedirectResponse = (redirect, { config }) => {
    return getFullUrl(redirect.url) === config.url
        && redirect.method.toLowerCase() === config.method
}
