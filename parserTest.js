const fs = require('fs');
const lexer = require('php7parser').Lexer;
const parser = require('php7parser').Parser;
const tokenTypeToString = require('php7parser').tokenTypeToString;
const phraseTypeToString = require('php7parser').phraseTypeToString;

const dir = './cases';

fs.readdir(dir, (err, files) => {
    if (err) {
        console.log('Folder not found: ' + dir);
        return;
    }

    for (const file of files) {
        if (!file.endsWith('.php')) {
            continue;
        }

        const filePath = dir + '/' + file;

        fs.readFile(filePath, (err, data) => {
            if (err) {
                console.log('File not found: ' + file);
                return;
            }
            let phpContent = data.toString();
            
            let rootNode = parser.parse(phpContent);
            let output = traverse(rootNode);

            const outPath = filePath + '.parsed';

            fs.readFile(outPath, (err, data) => {
                if (err) {
                    console.log('Error reading: ' + outPath);
                    return;
                }

                let testContent = data.toString();

                if (testContent != output) {
                    console.log(file + ' fail test');
                    console.log(phpContent + '\n');
                    console.log(testContent);
                    console.log(output);
                }
            });
        });
    }
});

function traverse(node, depth) {
    let str = '';

    if (depth === undefined) {
        depth = 0;
    }

    let indent = '';
    for (let i = 0; i < depth; i++) {
        indent += '-';
    }

    if (node.phraseType !== undefined) {
        str += indent + phraseToString(node) + '\n';
    } else {
        str += indent + tokenToString(node) + '\n';
    }

    if (node.phraseType !== undefined && node.children) {
        for (let i = 0; i < node.children.length; i++) {
            str += traverse(node.children[i], depth + 1);
        }
    }

    return str;
}

function tokenToString(token) {
    let str = tokenTypeToString(token.tokenType) + ' ' + token.offset + ' ' + token.length;

    for (const mode of token.modeStack) {
        str += ' ' + modeToString(mode);
    }

    return str;
}

function phraseToString(phrase) {
    return phraseTypeToString(phrase.phraseType);
}

function parserErrToString(parseErr) {
    return phraseTypeToString(parseError.phraseType);
}

function modeToString(mode) {
    switch (mode) {
        case 0:
            return 'ModeInitial';
        case 1:
            return 'ModeScripting';
        case 2:
            return 'ModeLookingForProperty';
        case 3:
            return 'ModeDoubleQuotes';
        case 4:
            return 'ModeNowDoc';
        case 5:
            return 'ModeHereDoc';
        case 6:
            return 'ModeEndHereDoc';
        case 7:
            return 'ModeBacktick';
        case 8:
            return 'ModeVarOffset';
        case 9:
            return 'ModeLookingForVarName';
    }
}
