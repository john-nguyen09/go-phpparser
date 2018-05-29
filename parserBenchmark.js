const fs = require('fs');
const lexer = require('php7parser').Lexer;
const parser = require('php7parser').Parser;
const tokenTypeToString = require('php7parser').tokenTypeToString;
const phraseTypeToString = require('php7parser').phraseTypeToString;

const dir = './cases/moodle';

parseFolder(dir);

function parseFolder(dir) {
    fs.readdir(dir, (err, files) => {
        if (err) {
            console.log('Folder not found: ' + dir);
            return;
        }

        for (const file of files) {
            const filePath = dir + '/' + file;

            fs.stat(filePath, (err, stats) => {
                if (err) {
                    console.log('Cannot stat: ' + filePath);
                    return;
                }

                if (stats.isDirectory()) {
                    parseFolder(filePath);
                } else if (file.endsWith('.php')) {
                    fs.readFile(filePath, (err, data) => {
                        if (err) {
                            console.log('File not found: ' + file);
                            return;
                        }
                        rootNode = parser.parse(data.toString());
                    });
                }
            });
        }
    });
}
