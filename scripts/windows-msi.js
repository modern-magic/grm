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
const msiInfo = require('./info.json')
const pref_hooks = require('perf_hooks')

/**
 * We decide only provide x64 msi file for user.
 */

const main = async () => {
  const star = pref_hooks.performance.now()
  const root = process.cwd()
  const verPath = path.join(root, 'version.txt')
  const packedPath = path.join(root, 'build')
  try {
    ;[verPath, packedPath].map((p) => fs.existsSync(p))
    const ver = await fs.readFile(verPath, 'utf-8')
    const wingz = 'build/grm-windows-64.tar.gz'
    const wintar = 'windows/grm-windows-64'
    await fs.ensureDir(wintar)
    const args = ['-zxvf', wingz, '-C', wintar]
    await execa('tar', args)

    const info = Object.assign({}, msiInfo, { version: ver })

    let msiTmpl = await fs.readFile(path.join(root, 'scripts', 'app.wsx.tmpl'), 'utf8')

    // generator wsx file.

    Object.keys(info).forEach((c) => {
      const reg = new RegExp(`{{.${c}}}`, 'g')
      msiTmpl = msiTmpl.replace(reg, info[c])
    })

    const tpl = msiTmpl.replace('{{.buildSource}}', `${wintar}/grm.exe`)
    fs.outputFileSync(path.join(root, wintar, 'app.wsx'), tpl, 'utf8')

    switch (os.platform()) {
      case 'win32':
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
        break
      case 'linux':
        /**
         * In linux system. Can use wixl to build msi installer.
         * You can view docs here:http://manpages.ubuntu.com/manpages/bionic/man1/wixl.1.html
         * I'm not sure darwin system can use wixl. If you found darwin system can use wixl too.
         * You can modify this branch. Make it better :)
         */
        break
      default:
        throw new Error("can't capture platform, process will exit.")
    }
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
