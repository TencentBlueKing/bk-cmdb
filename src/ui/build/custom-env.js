'use strict'
process.cmdb = process.cmdb || {}
const argv = process.argv.slice(2)
argv.forEach(args => {
    const argsArr = args.split('=')
    process.cmdb[argsArr[0]] = argsArr[1]
})