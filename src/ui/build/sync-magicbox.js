const fs = require('fs')
const path = require('path')
const readPkg = require('read-pkg')
const writePkg = require('write-pkg')
const { execSync } = require('child_process')
const dependencyKeys = [
    'dependencies'
]
try {
    const config = {}
    process.argv.slice(2).forEach(str => {
        const argv = str.split('=')
        config[argv[0]] = argv[1]
    })
    if (!config['GIT_MAGICBOX']) {
        throw new Error('请使用npm run sync-magicbox GIT_MAGICBOX=[仓库地址] 指定magic源码仓库地址')
    }
    execSync('rm -rf magicbox-temp')
    execSync(`git clone ${config['GIT_MAGICBOX']} magicbox-temp  -b staging --depth 1`)
    const source = readPkg.sync()
    const from = readPkg.sync({cwd: path.resolve(__dirname, '../magicbox-temp')})
    const to = JSON.parse(JSON.stringify(source))
    dependencyKeys.forEach(dependencyKey => {
        const sourceDep = from[dependencyKey]
        const fromDep = from[dependencyKey]
        if (fromDep) {
            const merged = to[dependencyKey] || {}
            for (let key in fromDep) {
                if (sourceDep && sourceDep[key]) {
                    merged[key] = sourceDep[key]
                } else {
                    merged[key] = fromDep[key]
                }
            }
            to[dependencyKey] = merged
        }
    })
    writePkg.sync(to)
    execSync('rm -rf src/magicbox/src')
    execSync('cp -R magicbox-temp/src src/magicbox/src')
    execSync('rm src/magicbox/src/bk-magic-ui/package.json')
    execSync('rm -rf magicbox-temp')
} catch (error) {
    console.log(error)
}