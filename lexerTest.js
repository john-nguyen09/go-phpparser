const fs = require('fs');
const lexer = require('php7parser').Lexer;
const tokenTypeToString = require('php7parser').tokenTypeToString;

const dir = './lexer/cases';

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

            lexer.setInput(data.toString());

            let output = '';

            for (
                let token = lexer.lex();
                token.tokenType !== 1;
                token = lexer.lex()
            ) {
                output += tokenToString(token) + '\n';
            }

            const outPath = filePath + '.lexed';

            fs.readFile(outPath, (err, data) => {
                if (err) {
                    console.log('Error reading: ' + outPath);
                    return;
                }

                if (data.toString() != output) {
                    console.log(file + ' fail test');
                    console.log(data.toString());
                    console.log(output);
                }
            });
        });
    }
});

function tokenToString(token) {
    let str = tokenTypeToString(token.tokenType) + ' ' + token.offset + ' ' + token.length;

    for (const mode of token.modeStack) {
        str += ' ' + modeToString(mode);
    }

    return str;
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
