import fs from 'fs'
import path from 'path'
import nodeResolve from '@rollup/plugin-node-resolve'
import less from 'rollup-plugin-less';

const destPath = path.join(__dirname, '../../app/assets/public/')

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
                if (!id.endsWith('.css') && !id.endsWith('.less')) {
                    return null;
                }
                // if one css changed, then all should be rebuild, 
                // because its can be refernced in each other
                let folder = path.join(__dirname, 'styles')
                fs.readdir(folder, (_, files) => {
                    files.forEach((file) => {
                        if (file.endsWith('.css') || file.endsWith('.less')) {
                            this.addWatchFile(path.join(folder, file));
                        }
                    });
                });
            }
        },
        less({ 
            output: path.join(destPath, 'style.css'),
            include: [__dirname + '/' + '**/*.less', __dirname + '/' + '**/*.css'],
        }),
    ]
}
