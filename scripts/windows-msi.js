// Hello, and welcome to view the genetor MSI script file.
// Before you modify this file. I think we need these.
// WIX document: https://wixtoolset.org/documentation/tutorial/
// GUID Generator Tool: http://www.guidgen.com/
// As you can see. Currently grm is using WIX wrapper binary file(Only x64)
// But It's only just getting stared. PR welcome.
// A Note. For build speed. And it's a very simple script We don't
// need too much third party to do some easy thing.

const fs = require('fs')
const fsp = require('fs').promises
const path = require('path')
const child_process = require('child_process')
const os = require('os')
const zlib = require('zlib')

const defaultWd = process.cwd()
const WINDOW_DIRECOTRY = path.join(defaultWd, 'windows')

const exist = (path) =>
    fsp
        .access(path, fs.constants.F_OK)
        .then(() => true)
        .catch(() => false)

const remove = (path) => fsp.rm(path, { recursive: true, force: true })

const outputFile = async (file, data, option) => {
    const dirPath = path.dirname(file)
    if (!(await exist(path))) {
        await fsp.mkdir(dirPath, { recursive: true })
    }
    await fsp.writeFile(file, data, option)
}

const execa = (command, argvs) => {
    const cp = child_process.spawn(command, argvs, { stdio: ['inherit'], shell: true })
    cp.on('exit', () => cp.kill('SIGHUP'))
    return new Promise((resolve, reject) => {
        cp.on('close', (code) => (code === 0 ? resolve() : reject()))
        cp.on('error', reject)
    })
}

const unzip = (file, to) => {
    return new Promise((resolve, reject) => {
        fs.createReadStream(file)
            .pipe(zlib.createUnzip())
            .pipe(fs.createWriteStream(to))
            .on('close', resolve)
            .on('error', reject)
    })
}

const MSI_PLACEMENT = {
    author: 'Modern Magic',
    applicationName: 'Grm Windows',
    upgradeCode: '5f25abbf-75c7-4847-91a7-d8ef0e823e95',
    envId: 'GRM_HOME',
    license: 'assets/LICENSE.rtf',
    dialog: 'assets/UIDialog.bmp',
    banner: 'assets/UIBanner.bmp',
    buildSource: path.join(WINDOW_DIRECOTRY, 'grm.exe'),
}

const TPL_REG = /\{\{((?:.|\r?\n)+?)\}\}/g

const buildImpl = async () => {
    const verionTextPath = path.join(defaultWd, 'version.txt')
    const originalPath = path.join(defaultWd, 'build', 'grm-windows-64.tar.gz')
    const depsExist = await (await Promise.all([verionTextPath, originalPath].map((p) => exist(p)))).every(Boolean)
    if (!depsExist) throw new Error("Can't find grm-windows-64.tar.gz or version.text. Please check it exist.")
    const grmVersion = await fsp.readFile(verionTextPath, 'utf8')
    const msiXMLPath = path.join(defaultWd, 'scripts', 'app.wsx.tmpl')
    const msiXML = await fsp.readFile(msiXMLPath, 'utf8')
    Reflect.set(MSI_PLACEMENT, 'version', grmVersion)
    const injected = msiXML.replace(TPL_REG, (_, s) => {
        if (Reflect.has(MSI_PLACEMENT, s)) return Reflect.get(MSI_PLACEMENT, s)
        return _
    })
    await outputFile(path.join(WINDOW_DIRECOTRY, 'app.wsx'), injected, 'utf8')
    /**
     * the compiler step
     * 1. Unzip the compressed package to windows directory
     * 2. use WIX tools gnerator app.wixobj
     *    First we should use candle
     *    Then use light
     */
    await unzip(originalPath, MSI_PLACEMENT.buildSource)
    // https://github.com/wixtoolset/wix3/blob/develop/src/tools/candle/candle.cs#L178-L186
    // From wix v3 we can't write the result to the memory. So we have to do this.
    await execa('candle.exe', [
        '-o',
        `${WINDOW_DIRECOTRY}/app.wixobj`,
        `${WINDOW_DIRECOTRY}/app.wsx`,
        '-arch',
        'x64',
        '-ext',
        'WixUtilExtension',
    ])
    await execa('light.exe', [
        `${WINDOW_DIRECOTRY}/app.wixobj`,
        '-o',
        `${path.dirname(originalPath)}/grm-installer-64.msi`,
        '-ext',
        'WixUIExtension',
        '-ext',
        'WixUtilExtension',
    ])
    const pdb = path.join(defaultWd, 'build', 'grm-installer-64.wixpdb')
    if (await exist(pdb)) await remove(pdb)
}

const main = () => {
    // ensure current node version greater than 16.17.0
    const [major, minor] = process.versions.node.split('.')
    if (+major < 16 || (+major === 16 && +minor < 17)) {
        console.error('Packing script is depend on version above 16.17.0.')
        return 1
    }
    if (os.platform() !== 'win32') {
        console.error("Can't run on platforms other than Windows.")
        return 1
    }

    buildImpl()
        .catch((err) => {
            console.error(err)
            return 1
        })
        .finally(() => remove(WINDOW_DIRECOTRY))
}
if (require.main === module) {
    main()
}
