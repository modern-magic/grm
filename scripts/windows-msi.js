/**
 * Wix docs here: https://wixtoolset.org/documentation/tutorial/
 * GUID Generator: http://www.guidgen.com/
 * Author: Kanno
 * Time: 10/5/2022
 */
const fs = require('fs-extra')
const path = require('path')
const execa = require('execa')
const os = require('os')
const msiMetaData = require('./info.json')
const pref_hooks = require('perf_hooks')

/**
 * We decide only provide x64 msi file for user.
 */

const main = async () => {
    /**
     * Process exit can't be use here. windows-msi will use in makefile. If
     * we let the process break will affect the result of the makefile.
     */
    if (os.platform() !== 'win32') {
        console.error("Can't run on platforms other than Windows.")
        return 1
    }
    const star = pref_hooks.performance.now()
    const root = process.cwd()
    const verPath = path.join(root, 'version.txt')
    const packedPath = path.join(root, 'build', 'grm-windows-64.tar.gz')
    try {
        const noExist = [verPath, packedPath].map((p) => fs.existsSync(p)).some((p) => !p)
        if (noExist) throw new Error("Can't find grm-windows-64.tar.gz or version.text. Please check it exist.")
        const ver = await fs.readFile(verPath, 'utf-8')
        const wingz = 'build/grm-windows-64.tar.gz'
        const wintar = 'windows/grm-windows-64'
        await fs.ensureDir(wintar)
        const args = ['-zxvf', wingz, '-C', wintar]
        await execa('tar', args)

        const info = Object.assign({}, msiMetaData, { version: ver })

        let msiTmpl = await fs.readFile(path.join(root, 'scripts', 'app.wsx.tmpl'), 'utf8')

        Object.keys(info).forEach((c) => {
            const reg = new RegExp(`{{.${c}}}`, 'g')
            msiTmpl = msiTmpl.replace(reg, info[c])
        })

        const tpl = msiTmpl.replace('{{.buildSource}}', `${wintar}/grm.exe`)
        fs.outputFileSync(path.join(root, wintar, 'app.wsx'), tpl, 'utf8')

        await execa('candle.exe', [
            '-o',
            `${wintar}/app.wixobj`,
            `${wintar}/app.wsx`,
            '-arch',
            'x64',
            '-ext',
            'WixUtilExtension',
        ])
        await execa('light.exe', [
            `${wintar}/app.wixobj`,
            '-o',
            `${wintar}/grm-installer-64.msi`,
            '-ext',
            'WixUIExtension',
            '-ext',
            'WixUtilExtension',
        ])
        await fs.copy(`${wintar}/grm-installer-64.msi`, 'build/grm-installer-64.msi')
    } catch (error) {
        console.log(error)
        process.exit(1)
    } finally {
        await fs.remove('windows')
        const end = pref_hooks.performance.now() - star
        console.log('\x1b[36m%s \x1b[36m\x1b[0m', `✨ genreator all msi installer use ${Math.ceil(end)}ms.\n`)
    }
}

if (require.main === module) {
    main()
}
