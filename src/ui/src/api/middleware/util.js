export const getFullUrl = url => {
    return `${window.API_PREFIX}/${url}`
}

export const isSameRequest = (origin, config) => {
    const sameUrl = origin.url === config.url
    const sameMethod = origin.method.toLowerCase() === config.method
    return sameUrl && sameMethod
}

export const isRedirectResponse = (redirect, { config }) => {
    return config.redirectId === redirect.redirectId
}

let redirectId = 0
export const getRedirectId = () => {
    return redirectId++
}
