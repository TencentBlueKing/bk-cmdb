const fs = require('fs')
const path = require('path')

const baseDir = path.resolve(__dirname, '../src')

const fileCN = path.join(baseDir, 'i18n/lang/cn.json')
const fileEN = path.join(baseDir, 'i18n/lang/en.json')
const cn = require(fileCN)
const en = require(fileEN)

// 搜索模板中的语法正则
const searchTplRe = [/path="(.+?)"/gm, /\$t[c]?\((.+?)\)/gm, /i18n.t\((.+?)\)/gm]
// 搜索路由配置中的语法正则
const searchRouteCfgRe = [/i18n:([^,\t\n\}]+)/gm]
// 模板中的中文字符
const searchTplCn = /[\u4E00-\u9FA5]+/gm

// 是否将结果保存为文件
const save = process.argv[2]

function searchPathkey () {
    let fileCount = 0
    let wordCount = 0
    const matched = {}
    const matchedCn = {}

    function find (dir) {
        const dirList = fs.readdirSync(dir, { withFileTypes: true })
        dirList.forEach(dirent => {
            if (dirent.isDirectory()) {
                find(path.join(dir, dirent.name))
            } else if (dirent.isFile()) {
                const filepath = path.join(dir, dirent.name)
                const fileext = path.extname(filepath)
                const content = fs.readFileSync(filepath, { encoding: 'utf8' })
                let match
                let fileCounted = false

                const searchRe = dirent.name === 'router.config.js' ? searchRouteCfgRe : searchTplRe
                searchRe.forEach(re => {
                    while ((match = re.exec(content)) !== null) {
                        const pathkey = match[1].trim()
                        if (matched[filepath]) {
                            matched[filepath].push(pathkey)
                        } else {
                            matched[filepath] = [pathkey]
                        }

                        !fileCounted && fileCount++
                        fileCounted = true
                        wordCount++
                    }
                })

                // 模板中的中文字符
                if (['.vue'].includes(fileext)) {
                    while ((match = searchTplCn.exec(content)) !== null) {
                        if (matchedCn[filepath]) {
                            matchedCn[filepath].push(match[0])
                        } else {
                            matchedCn[filepath] = [match[0]]
                        }
                    }
                }
            }
        })
    }

    find(baseDir)

    return { fileCount, wordCount, matched, matchedCn }
}

console.group('-- 总览 --')
console.time('time')

const { fileCount, wordCount, matched, matchedCn } = searchPathkey()

console.log(`语句查找: ${fileCount} 文件中有 ${wordCount} 个结果`)
console.timeEnd('time')
console.groupEnd('-- 总览 --')

const allPathKeyEN = Object.keys(en)
const allPathKeyCN = Object.keys(cn)

const missingCN = {}
const missingEN = {}
let missingCNCount = 0
let missingENCount = 0
Object.keys(matched).forEach(filepath => {
    const pathkeys = matched[filepath]

    pathkeys.forEach(pathkey => {
        // 获取key处理因正则的匹配结果而定，key值与翻译文件中的key定义一致
        // 现正则在匹配$t[c]语法时当存在多个“(),”字符时会不精准，但可以确保不会产生遗漏
        const key = pathkey.split(',', 1)[0].replace(/['"]/g, '')

        if (!allPathKeyCN.includes(key)) {
            if (missingCN[filepath]) {
                missingCN[filepath].push(key)
            } else {
                missingCN[filepath] = [key]
            }

            missingCNCount++
        }

        if (!allPathKeyEN.includes(key)) {
            if (missingEN[filepath]) {
                missingEN[filepath].push(key)
            } else {
                missingEN[filepath] = [key]
            }

            missingENCount++
        }
    })
})

console.group('-- cn missing --')
console.log(`翻译文件: ${fileCN}`)
console.log(`${Object.keys(missingCN).length} 文件中有 ${missingCNCount} 个结果`)
console.log(missingCN)
console.groupEnd('-- cn missing --')

console.group('-- en missing --')
console.log(`翻译文件: ${fileEN}`)
console.log(`${Object.keys(missingEN).length} 文件中有 ${missingENCount} 个结果`)
console.log(missingEN)
console.groupEnd('-- en missing --')

const missingRawCN = {}
let missingRawCNCount = 0
Object.keys(matchedCn).forEach(filepath => {
    const pathkeys = matchedCn[filepath]

    pathkeys.forEach(pathkey => {
        const key = pathkey

        if (!allPathKeyCN.includes(key) || !allPathKeyEN.includes(key)) {
            if (missingRawCN[filepath]) {
                missingRawCN[filepath].push(key)
            } else {
                missingRawCN[filepath] = [key]
            }

            missingRawCNCount++
        }
    })
})

console.group('-- 模板中的无翻译词的中文(beta) --')
console.log(`${Object.keys(missingRawCN).length} 文件中有 ${missingRawCNCount} 个结果`)
console.log(missingRawCN)
console.groupEnd('-- 模板中的无翻译词的中文(beta) --')

if (save) {
    fs.writeFile(path.join(__dirname, 'trans-mathced-all.json'), JSON.stringify(matched, null, 4), (err) => {
        if (err) throw err
        console.log('trans-mathced-all.json saved')
    })

    fs.writeFile(path.join(__dirname, 'trans-missing-cn.json'), JSON.stringify(missingCN, null, 4), (err) => {
        if (err) throw err
        console.log('trans-missing-cn.json saved')
    })

    fs.writeFile(path.join(__dirname, 'trans-missing-en.json'), JSON.stringify(missingEN, null, 4), (err) => {
        if (err) throw err
        console.log('trans-missing-en.json saved')
    })

    fs.writeFile(path.join(__dirname, 'trans-missing-rawcn.json'), JSON.stringify(missingRawCN, null, 4), (err) => {
        if (err) throw err
        console.log('trans-missing-rawcn.json saved')
    })
}
