export default async function (app, to, from) {
    const functions = []
    return Promise.all(functions.map(func => func(app, to, from)))
}
