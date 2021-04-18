import fs from 'fs'
import path from 'path'
import nodeResolve from '@rollup/plugin-node-resolve'
import scss from 'rollup-plugin-scss';

const destPath = path.join(__dirname, '../../app/assets/public/')

function isStyleFile(fname) {
    return fname.endsWith('.css') || fname.endsWith('.less') || fname.endsWith('.scss')
}

export default {
    input: path.join(__dirname, 'index.js'),
    output: {
        file: path.join(destPath, 'editor.bundle.js'),
        format: 'iife',
    },
    plugins: [
        nodeResolve(),
        {
            name: 'watch-external-styles',
            async transform(_code, id){
                if (!isStyleFile(id)) {
                    return null;
                }
                // if one css changed, then all should be rebuild, 
                // because its can be refernced in each other
                let folder = path.join(__dirname, 'styles')
                fs.readdir(folder, (_, files) => {
                    files.forEach((file) => {
                        if (isStyleFile(file)) {
                            this.addWatchFile(path.join(folder, file));
                        }
                    });
                });
            }
        },
        scss({
            output: true,
            output: path.join(destPath, 'style.css'),
        }),
    ]
}
