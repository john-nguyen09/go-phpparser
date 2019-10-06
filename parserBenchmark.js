const fs = require('fs');
const path = require('path');
const lexer = require('php7parser').Lexer;
const parser = require('php7parser').Parser;
const tokenTypeToString = require('php7parser').tokenTypeToString;
const phraseTypeToString = require('php7parser').phraseTypeToString;

const dir = path.join(__dirname, './cases/moodle');
let totalFiles = 0;
let finishedFiles = 0;
let start = process.hrtime();

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
                    totalFiles++;
                    fs.readFile(filePath, (err, data) => {
                        if (err) {
                            console.log('File not found: ' + file);
                            return;
                        }
                        rootNode = parser.parse(data.toString());
                        finishedFiles++;

                        if (totalFiles == finishedFiles) {
                            const diff = process.hrtime(start);
                            const elapsed = diff[0] * 1000 + diff[1] / 1000000;
                            console.log(`${elapsed.toFixed(2)} ms`);
                        }
                    });
                }
            });
        }
    });
}
