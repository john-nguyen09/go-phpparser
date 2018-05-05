'use strict';
Object.defineProperty(exports, "__esModule", { value: true });
var Token;
(function (Token) {
    function create(type, offset, length, modeStack) {
        return { tokenType: type, offset: offset, length: length, modeStack: modeStack };
    }
    Token.create = create;
})(Token = exports.Token || (exports.Token = {}));
var Lexer;
(function (Lexer) {
    var state;
    function setInput(text, lexerModeStack, position) {
        state = {
            position: position ? position : 0,
            input: text,
            modeStack: lexerModeStack ? lexerModeStack : [0],
            doubleQuoteScannedLength: -1,
            heredocLabel: null
        };
    }
    Lexer.setInput = setInput;
    function lex() {
        if (state.position >= state.input.length) {
            return {
                tokenType: 1,
                offset: state.position,
                length: 0,
                modeStack: state.modeStack
            };
        }
        let t;
        switch (state.modeStack[state.modeStack.length - 1]) {
            case 0:
                t = initial(state);
                break;
            case 1:
                t = scripting(state);
                break;
            case 2:
                t = lookingForProperty(state);
                break;
            case 3:
                t = doubleQuotes(state);
                break;
            case 4:
                t = nowdoc(state);
                break;
            case 5:
                t = heredoc(state);
                break;
            case 6:
                t = endHeredoc(state);
                break;
            case 7:
                t = backtick(state);
                break;
            case 8:
                t = varOffset(state);
                break;
            case 9:
                t = lookingForVarName(state);
                break;
            default:
                throw new Error('Unknown LexerMode');
        }
        return t ? t : lex();
    }
    Lexer.lex = lex;
    function isLabelStart(c) {
        let cp = c.charCodeAt(0);
        return (cp > 0x40 && cp < 0x5b) || (cp > 0x60 && cp < 0x7b) || cp === 0x5f || cp > 0x7f;
    }
    function isLabelChar(c) {
        let cp = c.charCodeAt(0);
        return (cp > 0x2f && cp < 0x3a) || (cp > 0x40 && cp < 0x5b) || (cp > 0x60 && cp < 0x7b) || cp === 0x5f || cp > 0x7f;
    }
    function isWhitespace(c) {
        return c === ' ' || c === '\n' || c === '\r' || c === '\t';
    }
    function initial(s) {
        let l = s.input.length;
        let c = s.input[s.position];
        let start = s.position;
        if (c === '<' && s.position + 1 < l && s.input[s.position + 1] === '?') {
            let tokenType = 156;
            if (s.input.substr(s.position, 5).toLowerCase() === '<?php' &&
                s.position + 5 < l && isWhitespace(s.input[s.position + 5])) {
                if (s.input[s.position + 5] === '\r' && s.position + 6 < l && s.input[s.position + 6] === '\n') {
                    s.position += 7;
                }
                else {
                    s.position += 6;
                }
            }
            else if (s.position + 2 < l && s.input[s.position + 2] === '=') {
                tokenType = 157;
                s.position += 3;
            }
            else {
                s.position += 2;
            }
            let t = { tokenType: tokenType, offset: start, length: s.position - start, modeStack: s.modeStack };
            s.modeStack = s.modeStack.slice(0, -1);
            s.modeStack.push(1);
            return t;
        }
        while (++s.position < l) {
            c = s.input[s.position];
            if (c === '<' && s.position + 1 < l && s.input[s.position + 1] === '?') {
                break;
            }
        }
        return { tokenType: 81, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scripting(s) {
        let c = s.input[s.position];
        let start = s.position;
        let l = s.input.length;
        let modeStack = s.modeStack;
        switch (c) {
            case ' ':
            case '\t':
            case '\n':
            case '\r':
                while (++s.position < l && isWhitespace(s.input[s.position])) { }
                return { tokenType: 161, offset: start, length: s.position - start, modeStack: modeStack };
            case '-':
                return scriptingMinus(s);
            case ':':
                if (++s.position < l && s.input[s.position] === ':') {
                    ++s.position;
                    return { tokenType: 133, offset: start, length: 2, modeStack: modeStack };
                }
                return { tokenType: 87, offset: start, length: 1, modeStack: modeStack };
            case '.':
                return scriptingDot(s);
            case '=':
                return scriptingEquals(s);
            case '+':
                return scriptingPlus(s);
            case '!':
                return scriptingExclamation(s);
            case '<':
                return scriptingLessThan(s);
            case '>':
                return scriptingGreaterThan(s);
            case '*':
                return scriptingAsterisk(s);
            case '/':
                return scriptingForwardSlash(s);
            case '%':
                if (++s.position < l && s.input[s.position] === '=') {
                    ++s.position;
                    return { tokenType: 145, offset: start, length: 2, modeStack: modeStack };
                }
                return { tokenType: 92, offset: start, length: 1, modeStack: modeStack };
            case '&':
                return scriptingAmpersand(s);
            case '|':
                return scriptingBar(s);
            case '^':
                if (++s.position < l && s.input[s.position] === '=') {
                    ++s.position;
                    return { tokenType: 105, offset: start, length: 2, modeStack: modeStack };
                }
                return { tokenType: 125, offset: start, length: 1, modeStack: modeStack };
            case ';':
                ++s.position;
                return { tokenType: 88, offset: start, length: 1, modeStack: modeStack };
            case ',':
                ++s.position;
                return { tokenType: 93, offset: start, length: 1, modeStack: modeStack };
            case '[':
                ++s.position;
                return { tokenType: 117, offset: start, length: 1, modeStack: modeStack };
            case ']':
                ++s.position;
                return { tokenType: 120, offset: start, length: 1, modeStack: modeStack };
            case '(':
                return scriptingOpenParenthesis(s);
            case ')':
                ++s.position;
                return { tokenType: 121, offset: start, length: 1, modeStack: modeStack };
            case '~':
                ++s.position;
                return { tokenType: 86, offset: start, length: 1, modeStack: modeStack };
            case '?':
                return scriptingQuestion(s);
            case '@':
                ++s.position;
                return { tokenType: 94, offset: start, length: 1, modeStack: modeStack };
            case '$':
                return scriptingDollar(s);
            case '#':
                ++s.position;
                return scriptingComment(s, start);
            case '0':
            case '1':
            case '2':
            case '3':
            case '4':
            case '5':
            case '6':
            case '7':
            case '8':
            case '9':
                return scriptingNumeric(s);
            case '{':
                ++s.position;
                s.modeStack = modeStack.slice(0);
                s.modeStack.push(1);
                return { tokenType: 116, offset: start, length: 1, modeStack: modeStack };
            case '}':
                ++s.position;
                if (s.modeStack.length > 1) {
                    s.modeStack = s.modeStack.slice(0, -1);
                }
                return { tokenType: 119, offset: start, length: 1, modeStack: modeStack };
            case '`':
                ++s.position;
                s.modeStack = s.modeStack.slice(0, -1);
                s.modeStack.push(7);
                return { tokenType: 95, offset: start, length: 1, modeStack: modeStack };
            case '\\':
                return scriptingBackslash(s);
            case '\'':
                return scriptingSingleQuote(s, start);
            case '"':
                return scriptingDoubleQuote(s, start);
            default:
                if (isLabelStart(c)) {
                    return scriptingLabelStart(s);
                }
                else {
                    ++s.position;
                    return { tokenType: 0, offset: start, length: 1, modeStack: s.modeStack };
                }
        }
    }
    function lookingForProperty(s) {
        let start = s.position;
        let c = s.input[s.position];
        let l = s.input.length;
        let modeStack = s.modeStack;
        switch (c) {
            case ' ':
            case '\t':
            case '\n':
            case '\r':
                while (++s.position < l && isWhitespace(s.input[s.position])) { }
                return { tokenType: 161, offset: start, length: s.position - start, modeStack: modeStack };
            default:
                if (isLabelStart(c)) {
                    while (++s.position < l && isLabelChar(s.input[s.position])) { }
                    s.modeStack = s.modeStack.slice(0, -1);
                    return { tokenType: 83, offset: start, length: s.position - start, modeStack: modeStack };
                }
                if (c === '-' && s.position + 1 < l && s.input[s.position + 1] === '>') {
                    s.position += 2;
                    return { tokenType: 115, offset: start, length: 2, modeStack: modeStack };
                }
                s.modeStack = s.modeStack.slice(0, -1);
                return null;
        }
    }
    function doubleQuotes(s) {
        let l = s.input.length;
        let c = s.input[s.position];
        let start = s.position;
        let modeStack = s.modeStack;
        let t;
        switch (c) {
            case '$':
                if ((t = encapsulatedDollar(s))) {
                    return t;
                }
                break;
            case '{':
                if (s.position + 1 < l && s.input[s.position + 1] === '$') {
                    s.modeStack = s.modeStack.slice(0);
                    s.modeStack.push(1);
                    ++s.position;
                    return { tokenType: 128, offset: start, length: 1, modeStack: modeStack };
                }
                break;
            case '"':
                s.modeStack = s.modeStack.slice(0, -1);
                s.modeStack.push(1);
                ++s.position;
                return { tokenType: 97, offset: start, length: 1, modeStack: modeStack };
            default:
                break;
        }
        return doubleQuotesAny(s);
    }
    function nowdoc(s) {
        let start = s.position;
        let n = start;
        let l = s.input.length;
        let c;
        let modeStack = s.modeStack;
        while (n < l) {
            c = s.input[n++];
            switch (c) {
                case '\r':
                    if (n < l && s.input[n] === '\n') {
                        ++n;
                    }
                case '\n':
                    if (n < l && s.heredocLabel === s.input.substr(n, s.heredocLabel.length)) {
                        let k = n + s.heredocLabel.length;
                        if (k < l && s.input[k] === ';') {
                            ++k;
                        }
                        if (k < l && (s.input[k] === '\n' || s.input[k] === '\r')) {
                            let nl = s.input.slice(n - 2, n);
                            if (nl === '\r\n') {
                                n -= 2;
                            }
                            else {
                                --n;
                            }
                            s.modeStack = s.modeStack.slice(0, -1);
                            s.modeStack.push(6);
                            break;
                        }
                    }
                default:
                    continue;
            }
            break;
        }
        s.position = n;
        return { tokenType: 80, offset: start, length: s.position - start, modeStack: modeStack };
    }
    function heredoc(s) {
        let l = s.input.length;
        let c = s.input[s.position];
        let start = s.position;
        let modeStack = s.modeStack;
        let t;
        switch (c) {
            case '$':
                if ((t = encapsulatedDollar(s))) {
                    return t;
                }
                break;
            case '{':
                if (s.position + 1 < l && s.input[s.position + 1] === '$') {
                    s.modeStack = s.modeStack.slice(0);
                    s.modeStack.push(1);
                    ++s.position;
                    return { tokenType: 128, offset: start, length: 1, modeStack: modeStack };
                }
                break;
            default:
                break;
        }
        return heredocAny(s);
    }
    function backtick(s) {
        let l = s.input.length;
        let c = s.input[s.position];
        let start = s.position;
        let modeStack = s.modeStack;
        let t;
        switch (c) {
            case '$':
                if ((t = encapsulatedDollar(s))) {
                    return t;
                }
                break;
            case '{':
                if (s.position + 1 < l && s.input[s.position + 1] === '$') {
                    s.modeStack = s.modeStack.slice(0);
                    s.modeStack.push(1);
                    ++s.position;
                    return { tokenType: 128, offset: start, length: 1, modeStack: modeStack };
                }
                break;
            case '`':
                s.modeStack = s.modeStack.slice(0, -1);
                s.modeStack.push(1);
                ++s.position;
                return { tokenType: 95, offset: start, length: 1, modeStack: modeStack };
            default:
                break;
        }
        return backtickAny(s);
    }
    function varOffset(s) {
        let start = s.position;
        let c = s.input[s.position];
        let l = s.input.length;
        let modeStack = s.modeStack;
        switch (s.input[s.position]) {
            case '$':
                if (s.position + 1 < l && isLabelStart(s.input[s.position + 1])) {
                    ++s.position;
                    while (++s.position < l && isLabelChar(s.input[s.position])) { }
                    return { tokenType: 84, offset: start, length: s.position - start, modeStack: s.modeStack };
                }
                break;
            case '[':
                ++s.position;
                return { tokenType: 117, offset: start, length: 1, modeStack: s.modeStack };
            case ']':
                s.modeStack = s.modeStack.slice(0, -1);
                ++s.position;
                return { tokenType: 120, offset: start, length: 1, modeStack: s.modeStack };
            case '-':
                ++s.position;
                return { tokenType: 143, offset: start, length: 1, modeStack: s.modeStack };
            default:
                if (c >= '0' && c <= '9') {
                    return varOffsetNumeric(s);
                }
                else if (isLabelStart(c)) {
                    while (++s.position < l && isLabelChar(s.input[s.position])) { }
                    return { tokenType: 83, offset: start, length: s.position - start, modeStack: s.modeStack };
                }
                break;
        }
        s.modeStack = s.modeStack.slice(0, -1);
        ++s.position;
        return { tokenType: 0, offset: start, length: 1, modeStack: modeStack };
    }
    function lookingForVarName(s) {
        let start = s.position;
        let l = s.input.length;
        let modeStack = s.modeStack;
        if (isLabelStart(s.input[s.position])) {
            let k = s.position + 1;
            while (++k < l && isLabelChar(s.input[k])) { }
            if (k < l && (s.input[k] === '[' || s.input[k] === '}')) {
                s.modeStack = s.modeStack.slice(0, -1);
                s.modeStack.push(1);
                s.position = k;
                return { tokenType: 84, offset: start, length: s.position - start, modeStack: modeStack };
            }
        }
        s.modeStack = s.modeStack.slice(0, -1);
        s.modeStack.push(1);
        return null;
    }
    function varOffsetNumeric(s) {
        let start = s.position;
        let c = s.input[s.position];
        let l = s.input.length;
        if (c === '0') {
            let k = s.position + 1;
            if (k < l && s.input[k] === 'b' && ++k < l && (s.input[k] === '1' || s.input[k] === '0')) {
                while (++k < l && (s.input[k] === '1' || s.input[k] === '0')) { }
                s.position = k;
                return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
            if (k < l && s.input[k] === 'x' && ++k < l && isHexDigit(s.input[k])) {
                while (++k < l && isHexDigit(s.input[k])) { }
                s.position = k;
                return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        while (++s.position < l && s.input[s.position] >= '0' && s.input[s.position] <= '9') { }
        return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function backtickAny(s) {
        let n = s.position;
        let c;
        let start = n;
        let l = s.input.length;
        if (s.input[n] === '\\' && n < l) {
            ++n;
        }
        while (n < l) {
            c = s.input[n++];
            switch (c) {
                case '`':
                    break;
                case '$':
                    if (n < l && (isLabelStart(s.input[n]) || s.input[n] === '{')) {
                        break;
                    }
                    continue;
                case '{':
                    if (n < l && s.input[n] === '$') {
                        break;
                    }
                    continue;
                case '\\':
                    if (n < l) {
                        ++n;
                    }
                default:
                    continue;
            }
            --n;
            break;
        }
        s.position = n;
        return { tokenType: 80, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function heredocAny(s) {
        let start = s.position;
        let n = start;
        let c;
        let l = s.input.length;
        let modeStack = s.modeStack;
        while (n < l) {
            c = s.input[n++];
            switch (c) {
                case '\r':
                    if (n < l && s.input[n] === '\n') {
                        ++n;
                    }
                case '\n':
                    if (n < l && s.input.slice(n, n + s.heredocLabel.length) === s.heredocLabel) {
                        let k = n + s.heredocLabel.length;
                        if (k < l && s.input[k] === ';') {
                            ++k;
                        }
                        if (k < l && (s.input[k] === '\n' || s.input[k] === '\r')) {
                            let nl = s.input.slice(n - 2, n);
                            if (nl === '\r\n') {
                                n -= 2;
                            }
                            else {
                                --n;
                            }
                            s.position = n;
                            s.modeStack = s.modeStack.slice(0, -1);
                            s.modeStack.push(6);
                            return { tokenType: 80, offset: start, length: s.position - start, modeStack: modeStack };
                        }
                    }
                    continue;
                case '$':
                    if (n < l && (isLabelStart(s.input[n]) || s.input[n] === '{')) {
                        break;
                    }
                    continue;
                case '{':
                    if (n < l && s.input[n] === '$') {
                        break;
                    }
                    continue;
                case '\\':
                    if (n < l && s.input[n] !== '\n' && s.input[n] !== '\r') {
                        ++n;
                    }
                default:
                    continue;
            }
            --n;
            break;
        }
        s.position = n;
        return { tokenType: 80, offset: start, length: s.position - start, modeStack: modeStack };
    }
    function endHeredoc(s) {
        let start = s.position;
        while (++s.position < s.input.length && (s.input[s.position] === '\r' || s.input[s.position] === '\n')) { }
        s.position += s.heredocLabel.length;
        s.heredocLabel = null;
        let t = { tokenType: 27, offset: start, length: s.position - start, modeStack: s.modeStack };
        s.modeStack = s.modeStack.slice(0, -1);
        s.modeStack.push(1);
        return t;
    }
    function doubleQuotesAny(s) {
        let start = s.position;
        if (s.doubleQuoteScannedLength > 0) {
            s.position = s.doubleQuoteScannedLength;
            s.doubleQuoteScannedLength = -1;
        }
        else {
            let n = s.position;
            let l = s.input.length;
            ++n;
            if (s.input[s.position] === '\\' && n + 1 < l) {
                ++n;
            }
            let c;
            while (n < l) {
                c = s.input[n++];
                switch (c) {
                    case '"':
                        break;
                    case '$':
                        if (n < l && (isLabelStart(s.input[n]) || s.input[n] == '{')) {
                            break;
                        }
                        continue;
                    case '{':
                        if (n < l && s.input[n] === '$') {
                            break;
                        }
                        continue;
                    case '\\':
                        if (n < l) {
                            ++n;
                        }
                    default:
                        continue;
                }
                --n;
                break;
            }
            s.position = n;
        }
        return { tokenType: 80, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function encapsulatedDollar(s) {
        let start = s.position;
        let l = s.input.length;
        let k = s.position + 1;
        let modeStack = s.modeStack;
        if (k >= l) {
            return null;
        }
        if (s.input[k] === '{') {
            s.position += 2;
            s.modeStack = s.modeStack.slice(0);
            s.modeStack.push(9);
            return { tokenType: 131, offset: start, length: 2, modeStack: modeStack };
        }
        if (!isLabelStart(s.input[k])) {
            return null;
        }
        while (++k < l && isLabelChar(s.input[k])) { }
        if (k < l && s.input[k] === '[') {
            s.modeStack = s.modeStack.slice(0);
            s.modeStack.push(8);
            s.position = k;
            return { tokenType: 84, offset: start, length: s.position - start, modeStack: modeStack };
        }
        if (k < l && s.input[k] === '-') {
            let n = k + 1;
            if (n < l && s.input[n] === '>' && ++n < l && isLabelStart(s.input[n])) {
                s.modeStack = s.modeStack.slice(0);
                s.modeStack.push(2);
                s.position = k;
                return { tokenType: 84, offset: start, length: s.position - start, modeStack: modeStack };
            }
        }
        s.position = k;
        return { tokenType: 84, offset: start, length: s.position - start, modeStack: modeStack };
    }
    function scriptingDoubleQuote(s, start) {
        ++s.position;
        let n = s.position;
        let c;
        let l = s.input.length;
        while (n < l) {
            c = s.input[n++];
            switch (c) {
                case '"':
                    s.position = n;
                    return { tokenType: 78, offset: start, length: s.position - start, modeStack: s.modeStack };
                case '$':
                    if (n < l && (isLabelStart(s.input[n]) || s.input[n] === '{')) {
                        break;
                    }
                    continue;
                case '{':
                    if (n < l && s.input[n] === '$') {
                        break;
                    }
                    continue;
                case '\\':
                    if (n < l) {
                        ++n;
                    }
                default:
                    continue;
            }
            --n;
            break;
        }
        s.doubleQuoteScannedLength = n;
        let modeStack = s.modeStack;
        s.modeStack = s.modeStack.slice(0, -1);
        s.modeStack.push(3);
        return { tokenType: 97, offset: start, length: s.position - start, modeStack: modeStack };
    }
    function scriptingSingleQuote(s, start) {
        let l = s.input.length;
        ++s.position;
        while (true) {
            if (s.position < l) {
                if (s.input[s.position] === '\'') {
                    ++s.position;
                    break;
                }
                else if (s.input[s.position++] === '\\' && s.position < l) {
                    ++s.position;
                }
            }
            else {
                return { tokenType: 80, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        return { tokenType: 78, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingBackslash(s) {
        let start = s.position;
        ++s.position;
        let t;
        if (s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '\'':
                    return scriptingSingleQuote(s, start);
                case '"':
                    return scriptingDoubleQuote(s, start);
                case '<':
                    t = scriptingHeredoc(s, start);
                    if (t) {
                        return t;
                    }
                default:
                    break;
            }
        }
        return { tokenType: 147, offset: start, length: 1, modeStack: s.modeStack };
    }
    const endHeredocLabelRegExp = /^;?(?:\r\n|\n|\r)/;
    function scriptingHeredoc(s, start) {
        let l = s.input.length;
        let k = s.position;
        let labelStart;
        let labelEnd;
        for (let kPlus3 = k + 3; k < kPlus3; ++k) {
            if (k >= l || s.input[k] !== '<') {
                return null;
            }
        }
        while (k < l && (s.input[k] === ' ' || s.input[k] === '\t')) {
            ++k;
        }
        let quote;
        if (k < l && (s.input[k] === '\'' || s.input[k] === '"')) {
            quote = s.input[k];
            ++k;
        }
        labelStart = k;
        if (k < l && isLabelStart(s.input[k])) {
            while (++k < l && isLabelChar(s.input[k])) { }
        }
        else {
            return null;
        }
        labelEnd = k;
        if (quote) {
            if (k < l && s.input[k] === quote) {
                ++k;
            }
            else {
                return null;
            }
        }
        if (k < l) {
            if (s.input[k] === '\r') {
                ++k;
                if (s.input[k] === '\n') {
                    ++k;
                }
            }
            else if (s.input[k] === '\n') {
                ++k;
            }
            else {
                return null;
            }
        }
        s.position = k;
        s.heredocLabel = s.input.slice(labelStart, labelEnd);
        let t = { tokenType: 154, offset: start, length: s.position - start, modeStack: s.modeStack };
        s.modeStack = s.modeStack.slice(0, -1);
        if (quote === '\'') {
            s.modeStack.push(4);
        }
        else {
            s.modeStack.push(5);
        }
        if (s.input.substr(s.position, s.heredocLabel.length) === s.heredocLabel &&
            s.input.substr(s.position + s.heredocLabel.length, 3).search(endHeredocLabelRegExp) >= 0) {
            s.modeStack.pop();
            s.modeStack.push(6);
        }
        return t;
    }
    function scriptingLabelStart(s) {
        let l = s.input.length;
        let start = s.position;
        while (++s.position < l && isLabelChar(s.input[s.position])) { }
        let text = s.input.slice(start, s.position);
        let tokenType = 0;
        if (text[0] === '_') {
            switch (text) {
                case '__CLASS__':
                    tokenType = 10;
                    break;
                case '__TRAIT__':
                    tokenType = 77;
                    break;
                case '__FUNCTION__':
                    tokenType = 74;
                    break;
                case '__METHOD__':
                    tokenType = 75;
                    break;
                case '__LINE__':
                    tokenType = 73;
                    break;
                case '__FILE__':
                    tokenType = 72;
                    break;
                case '__DIR__':
                    tokenType = 71;
                    break;
                case '__NAMESPACE__':
                    tokenType = 76;
                    break;
                default:
                    break;
            }
            if (tokenType > 0) {
                return { tokenType: tokenType, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        text = text.toLowerCase();
        switch (text) {
            case 'exit':
                tokenType = 29;
                break;
            case 'die':
                tokenType = 29;
                break;
            case 'function':
                tokenType = 35;
                break;
            case 'const':
                tokenType = 12;
                break;
            case 'return':
                tokenType = 59;
                break;
            case 'yield':
                return scriptingYield(s, start);
            case 'try':
                tokenType = 64;
                break;
            case 'catch':
                tokenType = 8;
                break;
            case 'finally':
                tokenType = 32;
                break;
            case 'throw':
                tokenType = 62;
                break;
            case 'if':
                tokenType = 39;
                break;
            case 'elseif':
                tokenType = 19;
                break;
            case 'endif':
                tokenType = 24;
                break;
            case 'else':
                tokenType = 18;
                break;
            case 'while':
                tokenType = 68;
                break;
            case 'endwhile':
                tokenType = 26;
                break;
            case 'do':
                tokenType = 16;
                break;
            case 'for':
                tokenType = 33;
                break;
            case 'endfor':
                tokenType = 22;
                break;
            case 'foreach':
                tokenType = 34;
                break;
            case 'endforeach':
                tokenType = 23;
                break;
            case 'declare':
                tokenType = 14;
                break;
            case 'enddeclare':
                tokenType = 21;
                break;
            case 'instanceof':
                tokenType = 43;
                break;
            case 'as':
                tokenType = 4;
                break;
            case 'switch':
                tokenType = 61;
                break;
            case 'endswitch':
                tokenType = 25;
                break;
            case 'case':
                tokenType = 7;
                break;
            case 'default':
                tokenType = 15;
                break;
            case 'break':
                tokenType = 5;
                break;
            case 'continue':
                tokenType = 13;
                break;
            case 'goto':
                tokenType = 37;
                break;
            case 'echo':
                tokenType = 17;
                break;
            case 'print':
                tokenType = 53;
                break;
            case 'class':
                tokenType = 9;
                break;
            case 'interface':
                tokenType = 45;
                break;
            case 'trait':
                tokenType = 63;
                break;
            case 'extends':
                tokenType = 30;
                break;
            case 'implements':
                tokenType = 40;
                break;
            case 'new':
                tokenType = 52;
                break;
            case 'clone':
                tokenType = 11;
                break;
            case 'var':
                tokenType = 67;
                break;
            case 'eval':
                tokenType = 28;
                break;
            case 'include_once':
                tokenType = 42;
                break;
            case 'include':
                tokenType = 41;
                break;
            case 'require_once':
                tokenType = 58;
                break;
            case 'require':
                tokenType = 57;
                break;
            case 'namespace':
                tokenType = 51;
                break;
            case 'use':
                tokenType = 66;
                break;
            case 'insteadof':
                tokenType = 44;
                break;
            case 'global':
                tokenType = 36;
                break;
            case 'isset':
                tokenType = 46;
                break;
            case 'empty':
                tokenType = 20;
                break;
            case '__halt_compiler':
                tokenType = 38;
                break;
            case 'static':
                tokenType = 60;
                break;
            case 'abstract':
                tokenType = 2;
                break;
            case 'final':
                tokenType = 31;
                break;
            case 'private':
                tokenType = 54;
                break;
            case 'protected':
                tokenType = 56;
                break;
            case 'public':
                tokenType = 55;
                break;
            case 'unset':
                tokenType = 65;
                break;
            case 'list':
                tokenType = 47;
                break;
            case 'array':
                tokenType = 3;
                break;
            case 'callable':
                tokenType = 6;
                break;
            case 'or':
                tokenType = 49;
                break;
            case 'and':
                tokenType = 48;
                break;
            case 'xor':
                tokenType = 50;
                break;
            default:
                break;
        }
        if (tokenType > 0) {
            return { tokenType: tokenType, offset: start, length: s.position - start, modeStack: s.modeStack };
        }
        return { tokenType: 83, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingYield(s, start) {
        let l = s.input.length;
        let k = s.position;
        if (k < l && isWhitespace(s.input[k])) {
            while (++k < l && isWhitespace(s.input[k])) { }
            if (s.input.substr(k, 4).toLowerCase() === 'from') {
                s.position = k + 4;
                return { tokenType: 70, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        return { tokenType: 69, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingQuestion(s) {
        let l = s.input.length;
        let start = s.position;
        ++s.position;
        if (s.position < l) {
            if (s.input[s.position] === '?') {
                ++s.position;
                return { tokenType: 122, offset: start, length: 2, modeStack: s.modeStack };
            }
            else if (s.input[s.position] === '>') {
                ++s.position;
                let modeStack = s.modeStack;
                s.modeStack = s.modeStack.slice(0, -1);
                s.modeStack.push(0);
                return { tokenType: 158, offset: start, length: s.position - start, modeStack: modeStack };
            }
        }
        return { tokenType: 96, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingDollar(s) {
        let start = s.position;
        let k = s.position;
        let l = s.input.length;
        ++k;
        if (k < l && isLabelStart(s.input[k])) {
            while (++k < l && isLabelChar(s.input[k])) { }
            s.position = k;
            return { tokenType: 84, offset: start, length: s.position - start, modeStack: s.modeStack };
        }
        ++s.position;
        return { tokenType: 90, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingOpenParenthesis(s) {
        let start = s.position;
        let k = start;
        let l = s.input.length;
        ++k;
        while (k < l && (s.input[k] === ' ' || s.input[k] === '\t')) {
            ++k;
        }
        let keywordStart = k;
        while (k < l && ((s.input[k] >= 'A' && s.input <= 'Z') || (s.input[k] >= 'a' && s.input <= 'z'))) {
            ++k;
        }
        let keywordEnd = k;
        while (k < l && (s.input[k] === ' ' || s.input[k] === '\t')) {
            ++k;
        }
        if (k < l && s.input[k] === ')') {
            let keyword = s.input.slice(keywordStart, keywordEnd).toLowerCase();
            let tokenType = 0;
            switch (keyword) {
                case 'int':
                case 'integer':
                    tokenType = 152;
                    break;
                case 'real':
                case 'float':
                case 'double':
                    tokenType = 153;
                    break;
                case 'string':
                case 'binary':
                    tokenType = 150;
                    break;
                case 'array':
                    tokenType = 155;
                    break;
                case 'object':
                    tokenType = 151;
                    break;
                case 'bool':
                case 'boolean':
                    tokenType = 148;
                    break;
                case 'unset':
                    tokenType = 149;
                    break;
                default:
                    break;
            }
            if (tokenType > 0) {
                s.position = k + 1;
                return { tokenType: tokenType, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        ++s.position;
        return { tokenType: 118, offset: start, length: 1, modeStack: s.modeStack };
    }
    function isHexDigit(c) {
        return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F');
    }
    function scriptingNumeric(s) {
        let start = s.position;
        let l = s.input.length;
        let k = s.position;
        if (s.input[s.position] === '0' && ++k < l) {
            if (s.input[k] === 'b' && ++k < l && (s.input[k] === '0' || s.input[k] === '1')) {
                while (++k < l && (s.input[k] === '0' || s.input[k] === '1')) { }
                s.position = k;
                return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
            k = s.position + 1;
            if (s.input[k] === 'x' && ++k < l && isHexDigit(s.input[k])) {
                while (++k < l && isHexDigit(s.input[k])) { }
                s.position = k;
                return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        while (++s.position < l && s.input[s.position] >= '0' && s.input[s.position] <= '9') { }
        if (s.input[s.position] === '.') {
            ++s.position;
            return scriptingNumericStartingWithDotOrE(s, start, true);
        }
        else if (s.input[s.position] === 'e' || s.input[s.position] === 'E') {
            return scriptingNumericStartingWithDotOrE(s, start, false);
        }
        return { tokenType: 82, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingNumericStartingWithDotOrE(s, start, hasDot) {
        let l = s.input.length;
        while (s.position < l && s.input[s.position] >= '0' && s.input[s.position] <= '9') {
            ++s.position;
        }
        if (s.position < l && (s.input[s.position] === 'e' || s.input[s.position] === 'E')) {
            let k = s.position + 1;
            if (k < l && (s.input[k] === '+' || s.input[k] === '-')) {
                ++k;
            }
            if (k < l && s.input[k] >= '0' && s.input[k] <= '9') {
                while (++k < l && s.input[k] >= '0' && s.input[k] <= '9') { }
                s.position = k;
                return { tokenType: 79, offset: start, length: s.position - start, modeStack: s.modeStack };
            }
        }
        return { tokenType: hasDot ? 79 : 82, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingBar(s) {
        let start = s.position;
        ++s.position;
        if (s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '=':
                    ++s.position;
                    return { tokenType: 110, offset: start, length: 2, modeStack: s.modeStack };
                case '|':
                    ++s.position;
                    return { tokenType: 124, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 123, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingAmpersand(s) {
        let start = s.position;
        ++s.position;
        if (s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '=':
                    ++s.position;
                    return { tokenType: 104, offset: start, length: 2, modeStack: s.modeStack };
                case '&':
                    ++s.position;
                    return { tokenType: 102, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 103, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingInlineCommentOrDocBlock(s) {
        let tokenType = 159;
        let start = s.position - 2;
        let l = s.input.length;
        if (s.position < l && s.input[s.position] === '*' && s.position + 1 < l && s.input[s.position + 1] !== '/') {
            ++s.position;
            tokenType = 160;
        }
        while (s.position < l) {
            if (s.input[s.position] === '*' && s.position + 1 < l && s.input[s.position + 1] === '/') {
                s.position += 2;
                break;
            }
            ++s.position;
        }
        return { tokenType: tokenType, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingForwardSlash(s) {
        let start = s.position;
        ++s.position;
        if (s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '=':
                    ++s.position;
                    return { tokenType: 130, offset: start, length: 2, modeStack: s.modeStack };
                case '*':
                    ++s.position;
                    return scriptingInlineCommentOrDocBlock(s);
                case '/':
                    ++s.position;
                    return scriptingComment(s, start);
                default:
                    break;
            }
        }
        return { tokenType: 91, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingComment(s, start) {
        let l = s.input.length;
        let c;
        while (s.position < l) {
            c = s.input[s.position];
            ++s.position;
            if (c === '\n' ||
                c === '\r' ||
                (c === '?' && s.position < l && s.input[s.position] === '>')) {
                --s.position;
                break;
            }
        }
        return { tokenType: 159, offset: start, length: s.position - start, modeStack: s.modeStack };
    }
    function scriptingAsterisk(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '*':
                    ++s.position;
                    if (s.position < s.input.length && s.input[s.position] === '=') {
                        ++s.position;
                        return { tokenType: 114, offset: start, length: 3, modeStack: s.modeStack };
                    }
                    return { tokenType: 113, offset: start, length: 2, modeStack: s.modeStack };
                case '=':
                    ++s.position;
                    return { tokenType: 146, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 101, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingGreaterThan(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '>':
                    ++s.position;
                    if (s.position < s.input.length && s.input[s.position] === '=') {
                        ++s.position;
                        return { tokenType: 109, offset: start, length: 3, modeStack: s.modeStack };
                    }
                    return { tokenType: 108, offset: start, length: 2, modeStack: s.modeStack };
                case '=':
                    ++s.position;
                    return { tokenType: 137, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 100, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingLessThan(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '>':
                    ++s.position;
                    return { tokenType: 139, offset: start, length: 2, modeStack: s.modeStack };
                case '<':
                    ++s.position;
                    if (s.position < s.input.length) {
                        if (s.input[s.position] === '=') {
                            ++s.position;
                            return { tokenType: 107, offset: start, length: 3, modeStack: s.modeStack };
                        }
                        else if (s.input[s.position] === '<') {
                            s.position -= 2;
                            let heredoc = scriptingHeredoc(s, start);
                            if (heredoc) {
                                return heredoc;
                            }
                            else {
                                s.position += 2;
                            }
                        }
                    }
                    return { tokenType: 106, offset: start, length: 2, modeStack: s.modeStack };
                case '=':
                    ++s.position;
                    if (s.position < s.input.length && s.input[s.position] === '>') {
                        ++s.position;
                        return { tokenType: 142, offset: start, length: 3, modeStack: s.modeStack };
                    }
                    return { tokenType: 141, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 99, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingExclamation(s) {
        let start = s.position;
        if (++s.position < s.input.length && s.input[s.position] === '=') {
            if (++s.position < s.input.length && s.input[s.position] === '=') {
                ++s.position;
                return { tokenType: 140, offset: start, length: 3, modeStack: s.modeStack };
            }
            return { tokenType: 139, offset: start, length: 2, modeStack: s.modeStack };
        }
        return { tokenType: 89, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingPlus(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '=':
                    ++s.position;
                    return { tokenType: 112, offset: start, length: 2, modeStack: s.modeStack };
                case '+':
                    ++s.position;
                    return { tokenType: 135, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 111, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingEquals(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '=':
                    if (++s.position < s.input.length && s.input[s.position] === '=') {
                        ++s.position;
                        return { tokenType: 138, offset: start, length: 3, modeStack: s.modeStack };
                    }
                    return { tokenType: 136, offset: start, length: 2, modeStack: s.modeStack };
                case '>':
                    ++s.position;
                    return { tokenType: 132, offset: start, length: 2, modeStack: s.modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 85, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingDot(s) {
        let start = s.position;
        if (++s.position < s.input.length) {
            let c = s.input[s.position];
            if (c === '=') {
                ++s.position;
                return { tokenType: 127, offset: start, length: 2, modeStack: s.modeStack };
            }
            else if (c === '.' && s.position + 1 < s.input.length && s.input[s.position + 1] === '.') {
                s.position += 2;
                return { tokenType: 134, offset: start, length: 3, modeStack: s.modeStack };
            }
            else if (c >= '0' && c <= '9') {
                return scriptingNumericStartingWithDotOrE(s, start, true);
            }
        }
        return { tokenType: 126, offset: start, length: 1, modeStack: s.modeStack };
    }
    function scriptingMinus(s) {
        let start = s.position;
        let modeStack = s.modeStack;
        if (++s.position < s.input.length) {
            switch (s.input[s.position]) {
                case '>':
                    ++s.position;
                    s.modeStack = s.modeStack.slice(0);
                    s.modeStack.push(2);
                    return { tokenType: 115, offset: start, length: 2, modeStack: modeStack };
                case '-':
                    ++s.position;
                    return { tokenType: 129, offset: start, length: 2, modeStack: modeStack };
                case '=':
                    ++s.position;
                    return { tokenType: 144, offset: start, length: 2, modeStack: modeStack };
                default:
                    break;
            }
        }
        return { tokenType: 143, offset: start, length: 1, modeStack: s.modeStack };
    }
})(Lexer = exports.Lexer || (exports.Lexer = {}));
function tokenTypeToString(type) {
    switch (type) {
        case 0:
            return 'Unknown';
        case 1:
            return 'EndOfFile';
        case 2:
            return 'Abstract';
        case 3:
            return 'Array';
        case 4:
            return 'As';
        case 5:
            return 'Break';
        case 6:
            return 'Callable';
        case 7:
            return 'Case';
        case 8:
            return 'Catch';
        case 9:
            return 'Class';
        case 10:
            return 'ClassConstant';
        case 11:
            return 'Clone';
        case 12:
            return 'Const';
        case 13:
            return 'Continue';
        case 14:
            return 'Declare';
        case 15:
            return 'Default';
        case 16:
            return 'Do';
        case 17:
            return 'Echo';
        case 18:
            return 'Else';
        case 19:
            return 'ElseIf';
        case 20:
            return 'Empty';
        case 21:
            return 'EndDeclare';
        case 22:
            return 'EndFor';
        case 23:
            return 'EndForeach';
        case 24:
            return 'EndIf';
        case 25:
            return 'EndSwitch';
        case 26:
            return 'EndWhile';
        case 27:
            return 'EndHeredoc';
        case 28:
            return 'Eval';
        case 29:
            return 'Exit';
        case 30:
            return 'Extends';
        case 31:
            return 'Final';
        case 32:
            return 'Finally';
        case 33:
            return 'For';
        case 34:
            return 'ForEach';
        case 35:
            return 'Function';
        case 36:
            return 'Global';
        case 37:
            return 'Goto';
        case 38:
            return 'HaltCompiler';
        case 39:
            return 'If';
        case 40:
            return 'Implements';
        case 41:
            return 'Include';
        case 42:
            return 'IncludeOnce';
        case 43:
            return 'InstanceOf';
        case 44:
            return 'InsteadOf';
        case 45:
            return 'Interface';
        case 46:
            return 'Isset';
        case 47:
            return 'List';
        case 48:
            return 'And';
        case 49:
            return 'Or';
        case 50:
            return 'Xor';
        case 51:
            return 'Namespace';
        case 52:
            return 'New';
        case 53:
            return 'Print';
        case 54:
            return 'Private';
        case 55:
            return 'Public';
        case 56:
            return 'Protected';
        case 57:
            return 'Require';
        case 58:
            return 'RequireOnce';
        case 59:
            return 'Return';
        case 60:
            return 'Static';
        case 61:
            return 'Switch';
        case 62:
            return 'Throw';
        case 63:
            return 'Trait';
        case 64:
            return 'Try';
        case 65:
            return 'Unset';
        case 66:
            return 'Use';
        case 67:
            return 'Var';
        case 68:
            return 'While';
        case 69:
            return 'Yield';
        case 70:
            return 'YieldFrom';
        case 71:
            return 'DirectoryConstant';
        case 72:
            return 'FileConstant';
        case 73:
            return 'LineConstant';
        case 74:
            return 'FunctionConstant';
        case 75:
            return 'MethodConstant';
        case 76:
            return 'NamespaceConstant';
        case 77:
            return 'TraitConstant';
        case 78:
            return 'StringLiteral';
        case 79:
            return 'FloatingLiteral';
        case 80:
            return 'EncapsulatedAndWhitespace';
        case 81:
            return 'Text';
        case 82:
            return 'IntegerLiteral';
        case 83:
            return 'Name';
        case 84:
            return 'VariableName';
        case 85:
            return 'Equals';
        case 86:
            return 'Tilde';
        case 87:
            return 'Colon';
        case 88:
            return 'Semicolon';
        case 89:
            return 'Exclamation';
        case 90:
            return 'Dollar';
        case 91:
            return 'ForwardSlash';
        case 92:
            return 'Percent';
        case 93:
            return 'Comma';
        case 94:
            return 'AtSymbol';
        case 95:
            return 'Backtick';
        case 96:
            return 'Question';
        case 97:
            return 'DoubleQuote';
        case 98:
            return 'SingleQuote';
        case 99:
            return 'LessThan';
        case 100:
            return 'GreaterThan';
        case 101:
            return 'Asterisk';
        case 102:
            return 'AmpersandAmpersand';
        case 103:
            return 'Ampersand';
        case 104:
            return 'AmpersandEquals';
        case 105:
            return 'CaretEquals';
        case 106:
            return 'LessThanLessThan';
        case 107:
            return 'LessThanLessThanEquals';
        case 108:
            return 'GreaterThanGreaterThan';
        case 109:
            return 'GreaterThanGreaterThanEquals';
        case 110:
            return 'BarEquals';
        case 111:
            return 'Plus';
        case 112:
            return 'PlusEquals';
        case 113:
            return 'AsteriskAsterisk';
        case 114:
            return 'AsteriskAsteriskEquals';
        case 115:
            return 'Arrow';
        case 116:
            return 'OpenBrace';
        case 117:
            return 'OpenBracket';
        case 118:
            return 'OpenParenthesis';
        case 119:
            return 'CloseBrace';
        case 120:
            return 'CloseBracket';
        case 121:
            return 'CloseParenthesis';
        case 122:
            return 'QuestionQuestion';
        case 123:
            return 'Bar';
        case 124:
            return 'BarBar';
        case 125:
            return 'Caret';
        case 126:
            return 'Dot';
        case 127:
            return 'DotEquals';
        case 128:
            return 'CurlyOpen';
        case 129:
            return 'MinusMinus';
        case 130:
            return 'ForwardslashEquals';
        case 131:
            return 'DollarCurlyOpen';
        case 132:
            return 'FatArrow';
        case 133:
            return 'ColonColon';
        case 134:
            return 'Ellipsis';
        case 135:
            return 'PlusPlus';
        case 136:
            return 'EqualsEquals';
        case 137:
            return 'GreaterThanEquals';
        case 138:
            return 'EqualsEqualsEquals';
        case 139:
            return 'ExclamationEquals';
        case 140:
            return 'ExclamationEqualsEquals';
        case 141:
            return 'LessThanEquals';
        case 142:
            return 'Spaceship';
        case 143:
            return 'Minus';
        case 144:
            return 'MinusEquals';
        case 145:
            return 'PercentEquals';
        case 146:
            return 'AsteriskEquals';
        case 147:
            return 'Backslash';
        case 148:
            return 'BooleanCast';
        case 149:
            return 'UnsetCast';
        case 150:
            return 'StringCast';
        case 151:
            return 'ObjectCast';
        case 152:
            return 'IntegerCast';
        case 153:
            return 'FloatCast';
        case 154:
            return 'StartHeredoc';
        case 155:
            return 'ArrayCast';
        case 156:
            return 'OpenTag';
        case 157:
            return 'OpenTagEcho';
        case 158:
            return 'CloseTag';
        case 159:
            return 'Comment';
        case 160:
            return 'DocumentComment';
        case 161:
            return 'Whitespace';
    }
}
exports.tokenTypeToString = tokenTypeToString;
