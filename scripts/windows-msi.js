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

const main = async () => {
  const star = pref_hooks.performance.now()
  const root = process.cwd()
  const verPath = path.join(root, 'version.txt')
  const packedPath = path.join(root, 'build')
  try {
    ;[verPath, packedPath].map((p) => fs.existsSync(p))
    const ver = await fs.readFile(verPath, 'utf-8')
    const winTars = ['build/grm-windows-32.tar.gz', 'build/grm-windows-64.tar.gz', 'build/grm-windows-arm64.tar.gz']
    const args = ['-zxvf']
    const windowsDir = []
    await Promise.all(
      winTars.map(async (p) => {
        const out = 'windows' + '/' + p.split('/')[1].replace('.tar.gz', '')
        windowsDir.push(out)
        await fs.ensureDir(out)
        execa('tar', [...args, p, '-C', out])
      })
    )

    const info = Object.assign({}, msiInfo, { version: ver })

    let msiTmpl = await fs.readFile(path.join(root, 'scripts', 'app.wsx.tmpl'), 'utf8')

    // generator all wsx file.

    Object.keys(info).forEach((c) => {
      const reg = new RegExp(`{{.${c}}}`, 'g')
      msiTmpl = msiTmpl.replace(reg, info[c])
    })

    await Promise.all(
      windowsDir.map((dir) => {
        const tpl = msiTmpl.replace('{{.buildSource}}', `${dir}/grm.exe`)
        const out = path.join(root, dir, 'app.wsx')
        fs.outputFile(out, tpl, 'utf8')
      })
    )

    switch (os.platform()) {
      case 'win32':
        await Promise.all(
          windowsDir.map(async (win) => {
            await execa('candle.exe', ['-o', `${win}/app.wixobj`, `${win}/app.wsx`, '-ext', 'WixUtilExtension'])
            await execa('light.exe', [
              `${win}/app.wixobj`,
              '-o',
              `${win}/app.msi`,
              '-ext',
              'WixUIExtension',
              '-ext',
              'WixUtilExtension',
            ])
          })
        )
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
    const end = pref_hooks.performance.now() - star
    console.log('\x1b[36m%s \x1b[36m\x1b[0m', `✨ genreator all msi installer use ${Math.ceil(end)}ms.\n`)
  }
}

if (require.main === module) {
  main()
}
