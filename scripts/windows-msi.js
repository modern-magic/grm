/**
 * Thanks: https://github.com/go-flutter-desktop/hover/blob/master/cmd/packaging/windows-msi.go
 * Can view issue here: https://github.com/goreleaser/goreleaser/issues/1295
 * From it we can use wix to build us golang library and wrapper to msi installer.
 * Author: Kanno
 * Time: 10/5/2022
 */
const { access } = require('fs')
const fs = require('fs/promises')
const path = require('path')
const { spawn } = require('child_process')
const os = require('os')

const checkExist = (path) =>
  new Promise((resolve, reject) => {
    access(path, (err) => {
      if (err) reject(err)
      resolve()
    })
  })

const main = async () => {
  const root = process.cwd()
  const verPath = path.join(root, 'version.txt')
  const packedPath = path.join(root, 'build')
  try {
    await Promise.all([verPath, packedPath].map((p) => checkExist(p)))
    const ver = await fs.readFile(verPath, 'utf-8')
    const winTars = [
      path.join(packedPath, 'grm-windows-32.tar.gz'),
      path.join(packedPath, 'grm-windows-64.tar.gz'),
      path.join(packedPath, 'grm-windows-arm64.tar.gz'),
    ]
    const args = ['-xzvf']

    switch (os.platform()) {
      case 'win32':
        /**
         * In windows system. I decide use Wix to wrapper binrary to msi.
         */
        console.log('windows')
        return
      case 'linux':
        /**
         * In linux system. Can use wixl to build msi installer.
         * You can view docs here:http://manpages.ubuntu.com/manpages/bionic/man1/wixl.1.html
         * I'm not sure darwin system can use wixl. If you found darwin system can use wixl too.
         * You can modify this branch. Make it better :)
         */
        console.log('linux')
        return
      default:
        throw new Error("can't capture platform, process will exit.")
    }
  } catch (error) {
    console.log(error)
    process.exit(1)
  }
}

if (require.main === module) {
  main()
}
