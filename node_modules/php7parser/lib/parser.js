'use strict';
Object.defineProperty(exports, "__esModule", { value: true });
const lexer_1 = require("./lexer");
var Parser;
(function (Parser) {
    function precedenceAssociativityTuple(t) {
        switch (t.tokenType) {
            case 113:
                return [48, 2];
            case 135:
                return [47, 2];
            case 129:
                return [47, 2];
            case 86:
                return [47, 2];
            case 152:
                return [47, 2];
            case 153:
                return [47, 2];
            case 150:
                return [47, 2];
            case 155:
                return [47, 2];
            case 151:
                return [47, 2];
            case 148:
                return [47, 2];
            case 149:
                return [47, 2];
            case 94:
                return [47, 2];
            case 43:
                return [46, 0];
            case 89:
                return [45, 2];
            case 101:
                return [44, 1];
            case 91:
                return [44, 1];
            case 92:
                return [44, 1];
            case 111:
                return [43, 1];
            case 143:
                return [43, 1];
            case 126:
                return [43, 1];
            case 106:
                return [42, 1];
            case 108:
                return [42, 1];
            case 99:
                return [41, 0];
            case 100:
                return [41, 0];
            case 141:
                return [41, 0];
            case 137:
                return [41, 0];
            case 136:
                return [40, 0];
            case 138:
                return [40, 0];
            case 139:
                return [40, 0];
            case 140:
                return [40, 0];
            case 142:
                return [40, 0];
            case 103:
                return [39, 1];
            case 125:
                return [38, 1];
            case 123:
                return [37, 1];
            case 102:
                return [36, 1];
            case 124:
                return [35, 1];
            case 122:
                return [34, 2];
            case 96:
                return [33, 1];
            case 85:
                return [32, 2];
            case 127:
                return [32, 2];
            case 112:
                return [32, 2];
            case 144:
                return [32, 2];
            case 146:
                return [32, 2];
            case 130:
                return [32, 2];
            case 145:
                return [32, 2];
            case 114:
                return [32, 2];
            case 104:
                return [32, 2];
            case 110:
                return [32, 2];
            case 105:
                return [32, 2];
            case 107:
                return [32, 2];
            case 109:
                return [32, 2];
            case 48:
                return [31, 1];
            case 50:
                return [30, 1];
            case 49:
                return [29, 1];
            default:
                throwUnexpectedTokenError(t);
        }
    }
    const statementListRecoverSet = [
        66,
        38,
        12,
        35,
        9,
        2,
        31,
        63,
        45,
        116,
        39,
        68,
        16,
        33,
        61,
        5,
        13,
        59,
        36,
        60,
        17,
        65,
        34,
        14,
        64,
        62,
        37,
        88,
        158,
        157,
        81,
        156
    ];
    const classMemberDeclarationListRecoverSet = [
        55,
        56,
        54,
        60,
        2,
        31,
        35,
        67,
        12,
        66
    ];
    const encapsulatedVariableListRecoverSet = [
        80,
        131,
        128
    ];
    function binaryOpToPhraseType(t) {
        switch (t.tokenType) {
            case 96:
                return 40;
            case 126:
            case 111:
            case 143:
                return 1;
            case 123:
            case 103:
            case 125:
                return 14;
            case 101:
            case 91:
            case 92:
                return 117;
            case 113:
                return 71;
            case 106:
            case 108:
                return 154;
            case 102:
            case 124:
            case 48:
            case 49:
            case 50:
                return 109;
            case 138:
            case 140:
            case 136:
            case 139:
                return 59;
            case 99:
            case 141:
            case 100:
            case 137:
            case 142:
                return 143;
            case 122:
                return 37;
            case 85:
                return 155;
            case 112:
            case 144:
            case 146:
            case 114:
            case 130:
            case 127:
            case 145:
            case 104:
            case 110:
            case 105:
            case 107:
            case 109:
                return 38;
            case 43:
                return 100;
            default:
                return 0;
        }
    }
    var tokenBuffer;
    var phraseStack;
    var errorPhrase;
    var recoverSetStack;
    function parse(text) {
        init(text);
        let stmtList = statementList([1]);
        hidden(stmtList);
        return stmtList;
    }
    Parser.parse = parse;
    function init(text, lexerModeStack) {
        lexer_1.Lexer.setInput(text, lexerModeStack);
        phraseStack = [];
        tokenBuffer = [];
        recoverSetStack = [];
        errorPhrase = null;
    }
    function start(phraseType, dontPushHiddenToParent) {
        if (!dontPushHiddenToParent) {
            hidden();
        }
        let p = {
            phraseType: phraseType ? phraseType : 0,
            children: []
        };
        phraseStack.push(p);
        return p;
    }
    function end() {
        return phraseStack.pop();
    }
    function hidden(p) {
        if (!p) {
            p = phraseStack[phraseStack.length - 1];
        }
        let t;
        while (true) {
            t = tokenBuffer.length ? tokenBuffer.shift() : lexer_1.Lexer.lex();
            if (t.tokenType < 159) {
                tokenBuffer.unshift(t);
                break;
            }
            else {
                p.children.push(t);
            }
        }
    }
    function optional(tokenType) {
        if (tokenType === peek().tokenType) {
            errorPhrase = null;
            return next();
        }
        else {
            return null;
        }
    }
    function optionalOneOf(tokenTypes) {
        if (tokenTypes.indexOf(peek().tokenType) >= 0) {
            errorPhrase = null;
            return next();
        }
        else {
            return null;
        }
    }
    function next(doNotPush) {
        let t = tokenBuffer.length ? tokenBuffer.shift() : lexer_1.Lexer.lex();
        if (t.tokenType === 1) {
            return t;
        }
        if (t.tokenType >= 159) {
            phraseStack[phraseStack.length - 1].children.push(t);
            return next(doNotPush);
        }
        else if (!doNotPush) {
            phraseStack[phraseStack.length - 1].children.push(t);
        }
        return t;
    }
    function expect(tokenType) {
        let t = peek();
        if (t.tokenType === tokenType) {
            errorPhrase = null;
            return next();
        }
        else if (tokenType === 88 && t.tokenType === 158) {
            return t;
        }
        else {
            error(tokenType);
            if (peek(1).tokenType === tokenType) {
                let predicate = (x) => { return x.tokenType === tokenType; };
                skip(predicate);
                errorPhrase = null;
                return next();
            }
            return null;
        }
    }
    function expectOneOf(tokenTypes) {
        let t = peek();
        if (tokenTypes.indexOf(t.tokenType) >= 0) {
            errorPhrase = null;
            return next();
        }
        else if (tokenTypes.indexOf(88) >= 0 && t.tokenType === 158) {
            return t;
        }
        else {
            error();
            if (tokenTypes.indexOf(peek(1).tokenType) >= 0) {
                let predicate = (x) => { return tokenTypes.indexOf(x.tokenType) >= 0; };
                skip(predicate);
                errorPhrase = null;
                return next();
            }
            return null;
        }
    }
    function peek(n) {
        let k = n ? n + 1 : 1;
        let bufferPos = -1;
        let t;
        while (true) {
            ++bufferPos;
            if (bufferPos === tokenBuffer.length) {
                tokenBuffer.push(lexer_1.Lexer.lex());
            }
            t = tokenBuffer[bufferPos];
            if (t.tokenType < 159) {
                --k;
            }
            if (t.tokenType === 1 || k === 0) {
                break;
            }
        }
        return t;
    }
    function skip(predicate) {
        let t;
        while (true) {
            t = tokenBuffer.length ? tokenBuffer.shift() : lexer_1.Lexer.lex();
            if (predicate(t) || t.tokenType === 1) {
                tokenBuffer.unshift(t);
                break;
            }
            else {
                errorPhrase.children.push(t);
            }
        }
    }
    function error(expected) {
        if (errorPhrase) {
            return;
        }
        errorPhrase = {
            phraseType: 60,
            children: [],
            unexpected: peek()
        };
        if (expected) {
            errorPhrase.expected = expected;
        }
        phraseStack[phraseStack.length - 1].children.push(errorPhrase);
    }
    function list(phraseType, elementFunction, elementStartPredicate, breakOn, recoverSet) {
        let p = start(phraseType);
        let t;
        let recoveryAttempted = false;
        let listRecoverSet = recoverSet ? recoverSet.slice(0) : [];
        if (breakOn) {
            Array.prototype.push.apply(listRecoverSet, breakOn);
        }
        recoverSetStack.push(listRecoverSet);
        while (true) {
            t = peek();
            if (elementStartPredicate(t)) {
                recoveryAttempted = false;
                p.children.push(elementFunction());
            }
            else if (!breakOn || breakOn.indexOf(t.tokenType) >= 0 || recoveryAttempted) {
                break;
            }
            else {
                error();
                t = peek(1);
                if (elementStartPredicate(t) || breakOn.indexOf(t.tokenType) >= 0) {
                    skip((x) => { return x === t; });
                }
                else {
                    defaultSyncStrategy();
                }
                recoveryAttempted = true;
            }
        }
        recoverSetStack.pop();
        return end();
    }
    function defaultSyncStrategy() {
        let mergedRecoverTokenTypeArray = [];
        for (let n = recoverSetStack.length - 1; n >= 0; --n) {
            Array.prototype.push.apply(mergedRecoverTokenTypeArray, recoverSetStack[n]);
        }
        let mergedRecoverTokenTypeSet = new Set(mergedRecoverTokenTypeArray);
        let predicate = (x) => { return mergedRecoverTokenTypeSet.has(x.tokenType); };
        skip(predicate);
    }
    function statementList(breakOn) {
        return list(157, statement, isStatementStart, breakOn, statementListRecoverSet);
    }
    function constDeclaration() {
        let p = start(42);
        next();
        p.children.push(delimitedList(44, constElement, isConstElementStartToken, 93, [88]));
        expect(88);
        return end();
    }
    function isClassConstElementStartToken(t) {
        return t.tokenType === 83 || isSemiReservedToken(t);
    }
    function isConstElementStartToken(t) {
        return t.tokenType === 83;
    }
    function constElement() {
        let p = start(43);
        expect(83);
        expect(85);
        p.children.push(expression(0));
        return end();
    }
    function expression(minPrecedence) {
        let precedence;
        let associativity;
        let op;
        let lhs = expressionAtom();
        let p;
        let rhs;
        let binaryPhraseType;
        while (true) {
            op = peek();
            binaryPhraseType = binaryOpToPhraseType(op);
            if (binaryPhraseType === 0) {
                break;
            }
            [precedence, associativity] = precedenceAssociativityTuple(op);
            if (precedence < minPrecedence) {
                break;
            }
            if (associativity === 1) {
                ++precedence;
            }
            if (binaryPhraseType === 40) {
                lhs = ternaryExpression(lhs);
                continue;
            }
            p = start(binaryPhraseType, true);
            p.children.push(lhs);
            next();
            if (binaryPhraseType === 100) {
                p.children.push(typeDesignator(101));
            }
            else {
                if (binaryPhraseType === 155 &&
                    peek().tokenType === 103) {
                    next();
                    p.phraseType = 16;
                }
                p.children.push(expression(precedence));
            }
            lhs = end();
        }
        return lhs;
    }
    function ternaryExpression(testExpr) {
        let p = start(40, true);
        p.children.push(testExpr);
        next();
        if (optional(87)) {
            p.children.push(expression(0));
        }
        else {
            p.children.push(expression(0));
            expect(87);
            p.children.push(expression(0));
        }
        return end();
    }
    function variableOrExpression() {
        let part = variableAtom();
        let isVariable = part.phraseType === 156;
        if (isDereferenceOperator(peek())) {
            part = variable(part);
            isVariable = true;
        }
        else {
            switch (part.phraseType) {
                case 141:
                case 84:
                case 144:
                    part = constantAccessExpression(part);
                    break;
                default:
                    break;
            }
        }
        if (!isVariable) {
            return part;
        }
        let t = peek();
        if (t.tokenType === 135) {
            return postfixExpression(132, part);
        }
        else if (t.tokenType === 129) {
            return postfixExpression(131, part);
        }
        else {
            return part;
        }
    }
    function constantAccessExpression(qName) {
        let p = start(41, true);
        p.children.push(qName);
        return end();
    }
    function postfixExpression(phraseType, variableNode) {
        let p = start(phraseType, true);
        p.children.push(variableNode);
        next();
        return end();
    }
    function isDereferenceOperator(t) {
        switch (t.tokenType) {
            case 117:
            case 116:
            case 115:
            case 118:
            case 133:
                return true;
            default:
                return false;
        }
    }
    function expressionAtom() {
        let t = peek();
        switch (t.tokenType) {
            case 60:
                if (peek(1).tokenType === 35) {
                    return anonymousFunctionCreationExpression();
                }
                else {
                    return variableOrExpression();
                }
            case 78:
                if (isDereferenceOperator(peek(1))) {
                    return variableOrExpression();
                }
                else {
                    return next(true);
                }
            case 84:
            case 90:
            case 3:
            case 117:
            case 147:
            case 83:
            case 51:
            case 118:
                return variableOrExpression();
            case 135:
                return unaryExpression(134);
            case 129:
                return unaryExpression(133);
            case 111:
            case 143:
            case 89:
            case 86:
                return unaryExpression(174);
            case 94:
                return unaryExpression(63);
            case 152:
            case 153:
            case 150:
            case 155:
            case 151:
            case 148:
            case 149:
                return unaryExpression(19);
            case 47:
                return listIntrinsic();
            case 11:
                return cloneExpression();
            case 52:
                return objectCreationExpression();
            case 79:
            case 82:
            case 73:
            case 72:
            case 71:
            case 77:
            case 75:
            case 74:
            case 76:
            case 10:
                return next(true);
            case 154:
                return heredocStringLiteral();
            case 97:
                return doubleQuotedStringLiteral();
            case 95:
                return shellCommandExpression();
            case 53:
                return printIntrinsic();
            case 69:
                return yieldExpression();
            case 70:
                return yieldFromExpression();
            case 35:
                return anonymousFunctionCreationExpression();
            case 41:
                return scriptInclusion(97);
            case 42:
                return scriptInclusion(98);
            case 57:
                return scriptInclusion(146);
            case 58:
                return scriptInclusion(147);
            case 28:
                return evalIntrinsic();
            case 20:
                return emptyIntrinsic();
            case 29:
                return exitIntrinsic();
            case 46:
                return issetIntrinsic();
            default:
                start(64);
                error();
                return end();
        }
    }
    function exitIntrinsic() {
        let p = start(70);
        next();
        if (optional(118)) {
            if (isExpressionStart(peek())) {
                p.children.push(expression(0));
            }
            expect(121);
        }
        return end();
    }
    function issetIntrinsic() {
        let p = start(107);
        next();
        expect(118);
        p.children.push(variableList([121]));
        expect(121);
        return end();
    }
    function emptyIntrinsic() {
        let p = start(55);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        return end();
    }
    function evalIntrinsic() {
        let p = start(69);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        return end();
    }
    function scriptInclusion(phraseType) {
        let p = start(phraseType);
        next();
        p.children.push(expression(0));
        return end();
    }
    function printIntrinsic() {
        let p = start(135);
        next();
        p.children.push(expression(0));
        return end();
    }
    function yieldFromExpression() {
        let p = start(181);
        next();
        p.children.push(expression(0));
        return end();
    }
    function yieldExpression() {
        let p = start(180);
        next();
        if (!isExpressionStart(peek())) {
            return end();
        }
        let keyOrValue = expression(0);
        p.children.push(keyOrValue);
        if (optional(132)) {
            p.children.push(expression(0));
        }
        return end();
    }
    function shellCommandExpression() {
        let p = start(153);
        next();
        p.children.push(encapsulatedVariableList(95));
        expect(95);
        return end();
    }
    function doubleQuotedStringLiteral() {
        let p = start(50);
        next();
        p.children.push(encapsulatedVariableList(97));
        expect(97);
        return end();
    }
    function encapsulatedVariableList(breakOn) {
        return list(58, encapsulatedVariable, isEncapsulatedVariableStart, [breakOn], encapsulatedVariableListRecoverSet);
    }
    function isEncapsulatedVariableStart(t) {
        switch (t.tokenType) {
            case 80:
            case 84:
            case 131:
            case 128:
                return true;
            default:
                return false;
        }
    }
    function encapsulatedVariable() {
        switch (peek().tokenType) {
            case 80:
                return next(true);
            case 84:
                let t = peek(1);
                if (t.tokenType === 117) {
                    return encapsulatedDimension();
                }
                else if (t.tokenType === 115) {
                    return encapsulatedProperty();
                }
                else {
                    return simpleVariable();
                }
            case 131:
                return dollarCurlyOpenEncapsulatedVariable();
            case 128:
                return curlyOpenEncapsulatedVariable();
            default:
                throwUnexpectedTokenError(peek());
        }
    }
    function curlyOpenEncapsulatedVariable() {
        let p = start(57);
        next();
        p.children.push(variable(variableAtom()));
        expect(119);
        return end();
    }
    function dollarCurlyOpenEncapsulatedVariable() {
        let p = start(57);
        next();
        let t = peek();
        if (t.tokenType === 84) {
            if (peek(1).tokenType === 117) {
                p.children.push(dollarCurlyEncapsulatedDimension());
            }
            else {
                let sv = start(156);
                next();
                p.children.push(end());
            }
        }
        else if (isExpressionStart(t)) {
            p.children.push(expression(0));
        }
        else {
            error();
        }
        expect(119);
        return end();
    }
    function dollarCurlyEncapsulatedDimension() {
        let p = start(160);
        next();
        next();
        p.children.push(expression(0));
        expect(120);
        return end();
    }
    function encapsulatedDimension() {
        let p = start(160);
        p.children.push(simpleVariable());
        next();
        switch (peek().tokenType) {
            case 83:
            case 82:
                next();
                break;
            case 84:
                p.children.push(simpleVariable());
                break;
            case 143:
                let u = start(174);
                next();
                expect(82);
                p.children.push(end());
                break;
            default:
                error();
                break;
        }
        expect(120);
        return end();
    }
    function encapsulatedProperty() {
        let p = start(136);
        p.children.push(simpleVariable());
        next();
        expect(83);
        return end();
    }
    function heredocStringLiteral() {
        let p = start(94);
        next();
        p.children.push(encapsulatedVariableList(27));
        expect(27);
        return end();
    }
    function anonymousClassDeclaration() {
        let p = start(2);
        p.children.push(anonymousClassDeclarationHeader());
        p.children.push(typeDeclarationBody(29, isClassMemberStart, classMemberDeclarationList));
        return end();
    }
    function anonymousClassDeclarationHeader() {
        let p = start(3);
        next();
        if (optional(118)) {
            if (isArgumentStart(peek())) {
                p.children.push(argumentList());
            }
            expect(121);
        }
        if (peek().tokenType === 30) {
            p.children.push(classBaseClause());
        }
        if (peek().tokenType === 40) {
            p.children.push(classInterfaceClause());
        }
        return end();
    }
    function classInterfaceClause() {
        let p = start(31);
        next();
        p.children.push(qualifiedNameList([116]));
        return end();
    }
    function classMemberDeclarationList() {
        return list(32, classMemberDeclaration, isClassMemberStart, [119], classMemberDeclarationListRecoverSet);
    }
    function isClassMemberStart(t) {
        switch (t.tokenType) {
            case 55:
            case 56:
            case 54:
            case 60:
            case 2:
            case 31:
            case 35:
            case 67:
            case 12:
            case 66:
                return true;
            default:
                return false;
        }
    }
    function classMemberDeclaration() {
        let p = start(61);
        let t = peek();
        switch (t.tokenType) {
            case 55:
            case 56:
            case 54:
            case 60:
            case 2:
            case 31:
                let modifiers = memberModifierList();
                t = peek();
                if (t.tokenType === 84) {
                    p.children.push(modifiers);
                    return propertyDeclaration(p);
                }
                else if (t.tokenType === 35) {
                    return methodDeclaration(p, modifiers);
                }
                else if (t.tokenType === 12) {
                    p.children.push(modifiers);
                    return classConstDeclaration(p);
                }
                else {
                    p.children.push(modifiers);
                    error();
                    return end();
                }
            case 35:
                return methodDeclaration(p, null);
            case 67:
                next();
                return propertyDeclaration(p);
            case 12:
                return classConstDeclaration(p);
            case 66:
                return traitUseClause(p);
            default:
                throwUnexpectedTokenError(t);
        }
    }
    function throwUnexpectedTokenError(t) {
        throw new Error(`Unexpected token: ${t.tokenType}`);
    }
    function traitUseClause(p) {
        p.phraseType = 170;
        next();
        p.children.push(qualifiedNameList([88, 116]));
        p.children.push(traitUseSpecification());
        return end();
    }
    function traitUseSpecification() {
        let p = start(171);
        let t = expectOneOf([88, 116]);
        if (t && t.tokenType === 116) {
            if (isTraitAdaptationStart(peek())) {
                p.children.push(traitAdaptationList());
            }
            expect(119);
        }
        return end();
    }
    function traitAdaptationList() {
        return list(163, traitAdaptation, isTraitAdaptationStart, [119]);
    }
    function isTraitAdaptationStart(t) {
        switch (t.tokenType) {
            case 83:
            case 147:
            case 51:
                return true;
            default:
                return isSemiReservedToken(t);
        }
    }
    function traitAdaptation() {
        let p = start(66);
        let t = peek();
        let t2 = peek(1);
        if (t.tokenType === 51 ||
            t.tokenType === 147 ||
            (t.tokenType === 83 &&
                (t2.tokenType === 133 || t2.tokenType === 147))) {
            p.children.push(methodReference());
            if (peek().tokenType === 44) {
                next();
                return traitPrecedence(p);
            }
        }
        else if (t.tokenType === 83 || isSemiReservedToken(t)) {
            let methodRef = start(116);
            methodRef.children.push(identifier());
            p.children.push(end());
        }
        else {
            error();
            return end();
        }
        return traitAlias(p);
    }
    function traitAlias(p) {
        p.phraseType = 164;
        expect(4);
        let t = peek();
        if (t.tokenType === 83 || isReservedToken(t)) {
            p.children.push(identifier());
        }
        else if (isMemberModifier(t)) {
            next();
            t = peek();
            if (t.tokenType === 83 || isSemiReservedToken(t)) {
                p.children.push(identifier());
            }
        }
        else {
            error();
        }
        expect(88);
        return end();
    }
    function traitPrecedence(p) {
        p.phraseType = 169;
        p.children.push(qualifiedNameList([88]));
        expect(88);
        return end();
    }
    function methodReference() {
        let p = start(116);
        p.children.push(qualifiedName());
        expect(133);
        p.children.push(identifier());
        return end();
    }
    function methodDeclarationHeader(memberModifers) {
        let p = start(115, true);
        if (memberModifers) {
            p.children.push(memberModifers);
        }
        next();
        optional(103);
        p.children.push(identifier());
        expect(118);
        if (isParameterStart(peek())) {
            p.children.push(delimitedList(130, parameterDeclaration, isParameterStart, 93, [121]));
        }
        expect(121);
        if (peek().tokenType === 87) {
            p.children.push(returnType());
        }
        return end();
    }
    function methodDeclaration(p, memberModifers) {
        p.phraseType = 113;
        p.children.push(methodDeclarationHeader(memberModifers));
        p.children.push(methodDeclarationBody());
        return end();
    }
    function methodDeclarationBody() {
        let p = start(114);
        if (peek().tokenType === 88) {
            next();
        }
        else {
            p.children.push(compoundStatement());
        }
        return end();
    }
    function identifier() {
        let p = start(95);
        let t = peek();
        if (t.tokenType === 83 || isSemiReservedToken(t)) {
            next();
        }
        else {
            error();
        }
        return end();
    }
    function interfaceDeclaration() {
        let p = start(103);
        p.children.push(interfaceDeclarationHeader());
        p.children.push(typeDeclarationBody(104, isClassMemberStart, interfaceMemberDeclarations));
        return end();
    }
    function typeDeclarationBody(phraseType, elementStartPredicate, listFunction) {
        let p = start(phraseType);
        expect(116);
        if (elementStartPredicate(peek())) {
            p.children.push(listFunction());
        }
        expect(119);
        return end();
    }
    function interfaceMemberDeclarations() {
        return list(106, classMemberDeclaration, isClassMemberStart, [119], classMemberDeclarationListRecoverSet);
    }
    function interfaceDeclarationHeader() {
        let p = start(105);
        next();
        expect(83);
        if (peek().tokenType === 30) {
            p.children.push(interfaceBaseClause());
        }
        return end();
    }
    function interfaceBaseClause() {
        let p = start(102);
        next();
        p.children.push(qualifiedNameList([116]));
        return end();
    }
    function traitDeclaration() {
        let p = start(165);
        p.children.push(traitDeclarationHeader());
        p.children.push(typeDeclarationBody(166, isClassMemberStart, traitMemberDeclarations));
        return end();
    }
    function traitDeclarationHeader() {
        let p = start(167);
        next();
        expect(83);
        return end();
    }
    function traitMemberDeclarations() {
        return list(168, classMemberDeclaration, isClassMemberStart, [119], classMemberDeclarationListRecoverSet);
    }
    function functionDeclaration() {
        let p = start(86);
        p.children.push(functionDeclarationHeader());
        p.children.push(functionDeclarationBody());
        return end();
    }
    function functionDeclarationBody() {
        let cs = compoundStatement();
        cs.phraseType = 87;
        return cs;
    }
    function functionDeclarationHeader() {
        let p = start(88);
        next();
        optional(103);
        expect(83);
        expect(118);
        if (isParameterStart(peek())) {
            p.children.push(delimitedList(130, parameterDeclaration, isParameterStart, 93, [121]));
        }
        expect(121);
        if (peek().tokenType === 87) {
            p.children.push(returnType());
        }
        return end();
    }
    function isParameterStart(t) {
        switch (t.tokenType) {
            case 103:
            case 134:
            case 84:
                return true;
            default:
                return isTypeDeclarationStart(t);
        }
    }
    function classDeclaration() {
        let p = start(28);
        p.children.push(classDeclarationHeader());
        p.children.push(typeDeclarationBody(29, isClassMemberStart, classMemberDeclarationList));
        return end();
    }
    function classDeclarationHeader() {
        let p = start(30);
        optionalOneOf([2, 31]);
        expect(9);
        expect(83);
        if (peek().tokenType === 30) {
            p.children.push(classBaseClause());
        }
        if (peek().tokenType === 40) {
            p.children.push(classInterfaceClause());
        }
        return end();
    }
    function classBaseClause() {
        let p = start(23);
        next();
        p.children.push(qualifiedName());
        return end();
    }
    function compoundStatement() {
        let p = start(39);
        expect(116);
        if (isStatementStart(peek())) {
            p.children.push(statementList([119]));
        }
        expect(119);
        return end();
    }
    function statement() {
        let t = peek();
        switch (t.tokenType) {
            case 51:
                return namespaceDefinition();
            case 66:
                return namespaceUseDeclaration();
            case 38:
                return haltCompilerStatement();
            case 12:
                return constDeclaration();
            case 35:
                {
                    let p1 = peek(1);
                    if (p1.tokenType === 118 ||
                        (p1.tokenType === 103 && peek(2).tokenType === 118)) {
                        return expressionStatement();
                    }
                    else {
                        return functionDeclaration();
                    }
                }
            case 9:
            case 2:
            case 31:
                return classDeclaration();
            case 63:
                return traitDeclaration();
            case 45:
                return interfaceDeclaration();
            case 116:
                return compoundStatement();
            case 39:
                return ifStatement();
            case 68:
                return whileStatement();
            case 16:
                return doStatement();
            case 33:
                return forStatement();
            case 61:
                return switchStatement();
            case 5:
                return breakStatement();
            case 13:
                return continueStatement();
            case 59:
                return returnStatement();
            case 36:
                return globalDeclaration();
            case 60:
                if (peek(1).tokenType === 84 &&
                    [88, 93,
                        158, 85].indexOf(peek(2).tokenType) >= 0) {
                    return functionStaticDeclaration();
                }
                else {
                    return expressionStatement();
                }
            case 81:
            case 156:
            case 158:
                return inlineText();
            case 34:
                return foreachStatement();
            case 14:
                return declareStatement();
            case 64:
                return tryStatement();
            case 62:
                return throwStatement();
            case 37:
                return gotoStatement();
            case 17:
            case 157:
                return echoIntrinsic();
            case 65:
                return unsetIntrinsic();
            case 88:
                return nullStatement();
            case 83:
                if (peek(1).tokenType === 87) {
                    return namedLabelStatement();
                }
            default:
                return expressionStatement();
        }
    }
    function inlineText() {
        let p = start(99);
        optional(158);
        optional(81);
        optional(156);
        return end();
    }
    function nullStatement() {
        start(127);
        next();
        return end();
    }
    function isCatchClauseStart(t) {
        return t.tokenType === 8;
    }
    function tryStatement() {
        let p = start(172);
        next();
        p.children.push(compoundStatement());
        let t = peek();
        if (t.tokenType === 8) {
            p.children.push(list(21, catchClause, isCatchClauseStart));
        }
        else if (t.tokenType !== 32) {
            error();
        }
        if (peek().tokenType === 32) {
            p.children.push(finallyClause());
        }
        return end();
    }
    function finallyClause() {
        let p = start(74);
        next();
        p.children.push(compoundStatement());
        return end();
    }
    function catchClause() {
        let p = start(20);
        next();
        expect(118);
        p.children.push(delimitedList(22, qualifiedName, isQualifiedNameStart, 123, [84]));
        expect(84);
        expect(121);
        p.children.push(compoundStatement());
        return end();
    }
    function declareDirective() {
        let p = start(46);
        expect(83);
        expect(85);
        expectOneOf([82, 79, 78]);
        return end();
    }
    function declareStatement() {
        let p = start(47);
        next();
        expect(118);
        p.children.push(declareDirective());
        expect(121);
        let t = peek();
        if (t.tokenType === 87) {
            next();
            p.children.push(statementList([21]));
            expect(21);
            expect(88);
        }
        else if (isStatementStart(t)) {
            p.children.push(statement());
        }
        else if (t.tokenType === 88) {
            next();
        }
        else {
            error();
        }
        return end();
    }
    function switchStatement() {
        let p = start(161);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        let t = expectOneOf([87, 116]);
        let tCase = peek();
        if (tCase.tokenType === 7 || tCase.tokenType === 15) {
            p.children.push(caseStatements(t && t.tokenType === 87 ?
                25 : 119));
        }
        if (t && t.tokenType === 87) {
            expect(25);
            expect(88);
        }
        else {
            expect(119);
        }
        return end();
    }
    function caseStatements(breakOn) {
        let p = start(18);
        let t;
        let caseBreakOn = [7, 15];
        caseBreakOn.push(breakOn);
        while (true) {
            t = peek();
            if (t.tokenType === 7) {
                p.children.push(caseStatement(caseBreakOn));
            }
            else if (t.tokenType === 15) {
                p.children.push(defaultStatement(caseBreakOn));
            }
            else if (breakOn === t.tokenType) {
                break;
            }
            else {
                error();
                break;
            }
        }
        return end();
    }
    function caseStatement(breakOn) {
        let p = start(17);
        next();
        p.children.push(expression(0));
        expectOneOf([87, 88]);
        if (isStatementStart(peek())) {
            p.children.push(statementList(breakOn));
        }
        return end();
    }
    function defaultStatement(breakOn) {
        let p = start(48);
        next();
        expectOneOf([87, 88]);
        if (isStatementStart(peek())) {
            p.children.push(statementList(breakOn));
        }
        return end();
    }
    function namedLabelStatement() {
        let p = start(118);
        next();
        next();
        return end();
    }
    function gotoStatement() {
        let p = start(92);
        next();
        expect(83);
        expect(88);
        return end();
    }
    function throwStatement() {
        let p = start(162);
        next();
        p.children.push(expression(0));
        expect(88);
        return end();
    }
    function foreachCollection() {
        let p = start(76);
        p.children.push(expression(0));
        return end();
    }
    function foreachKeyOrValue() {
        let p = start(79);
        p.children.push(expression(0));
        if (peek().tokenType === 132) {
            next();
            p.phraseType = 77;
        }
        return end();
    }
    function foreachValue() {
        let p = start(79);
        optional(103);
        p.children.push(expression(0));
        return end();
    }
    function foreachStatement() {
        let p = start(78);
        next();
        expect(118);
        p.children.push(foreachCollection());
        expect(4);
        let keyOrValue = peek().tokenType === 103 ? foreachValue() : foreachKeyOrValue();
        p.children.push(keyOrValue);
        if (keyOrValue.phraseType === 77) {
            p.children.push(foreachValue());
        }
        expect(121);
        let t = peek();
        if (t.tokenType === 87) {
            next();
            p.children.push(statementList([23]));
            expect(23);
            expect(88);
        }
        else if (isStatementStart(t)) {
            p.children.push(statement());
        }
        else {
            error();
        }
        return end();
    }
    function isVariableStart(t) {
        switch (t.tokenType) {
            case 84:
            case 90:
            case 118:
            case 3:
            case 117:
            case 78:
            case 60:
            case 83:
            case 51:
            case 147:
                return true;
            default:
                return false;
        }
    }
    function variableInitial() {
        return variable(variableAtom());
    }
    function variableList(breakOn) {
        return delimitedList(176, variableInitial, isVariableStart, 93, breakOn);
    }
    function unsetIntrinsic() {
        let p = start(175);
        next();
        expect(118);
        p.children.push(variableList([121]));
        expect(121);
        expect(88);
        return end();
    }
    function expressionInitial() {
        return expression(0);
    }
    function echoIntrinsic() {
        let p = start(51);
        next();
        p.children.push(delimitedList(72, expressionInitial, isExpressionStart, 93));
        expect(88);
        return end();
    }
    function isStaticVariableDclarationStart(t) {
        return t.tokenType === 84;
    }
    function functionStaticDeclaration() {
        let p = start(89);
        next();
        p.children.push(delimitedList(159, staticVariableDeclaration, isStaticVariableDclarationStart, 93, [88]));
        expect(88);
        return end();
    }
    function globalDeclaration() {
        let p = start(91);
        next();
        p.children.push(delimitedList(177, simpleVariable, isSimpleVariableStart, 93, [88]));
        expect(88);
        return end();
    }
    function isSimpleVariableStart(t) {
        switch (t.tokenType) {
            case 84:
            case 90:
                return true;
            default:
                return false;
        }
    }
    function staticVariableDeclaration() {
        let p = start(158);
        expect(84);
        if (peek().tokenType === 85) {
            p.children.push(functionStaticInitialiser());
        }
        return end();
    }
    function functionStaticInitialiser() {
        let p = start(90);
        next();
        p.children.push(expression(0));
        return end();
    }
    function continueStatement() {
        let p = start(45);
        next();
        if (isExpressionStart(peek())) {
            p.children.push(expression(0));
        }
        expect(88);
        return end();
    }
    function breakStatement() {
        let p = start(15);
        next();
        if (isExpressionStart(peek())) {
            p.children.push(expression(0));
        }
        expect(88);
        return end();
    }
    function returnStatement() {
        let p = start(148);
        next();
        if (isExpressionStart(peek())) {
            p.children.push(expression(0));
        }
        expect(88);
        return end();
    }
    function forExpressionGroup(phraseType, breakOn) {
        return delimitedList(phraseType, expressionInitial, isExpressionStart, 93, breakOn);
    }
    function forStatement() {
        let p = start(83);
        next();
        expect(118);
        if (isExpressionStart(peek())) {
            p.children.push(forExpressionGroup(82, [88]));
        }
        expect(88);
        if (isExpressionStart(peek())) {
            p.children.push(forExpressionGroup(75, [88]));
        }
        expect(88);
        if (isExpressionStart(peek())) {
            p.children.push(forExpressionGroup(80, [121]));
        }
        expect(121);
        let t = peek();
        if (t.tokenType === 87) {
            next();
            p.children.push(statementList([22]));
            expect(22);
            expect(88);
        }
        else if (isStatementStart(peek())) {
            p.children.push(statement());
        }
        else {
            error();
        }
        return end();
    }
    function doStatement() {
        let p = start(49);
        next();
        p.children.push(statement());
        expect(68);
        expect(118);
        p.children.push(expression(0));
        expect(121);
        expect(88);
        return end();
    }
    function whileStatement() {
        let p = start(179);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        let t = peek();
        if (t.tokenType === 87) {
            next();
            p.children.push(statementList([26]));
            expect(26);
            expect(88);
        }
        else if (isStatementStart(t)) {
            p.children.push(statement());
        }
        else {
            error();
        }
        return end();
    }
    function elseIfClause1() {
        let p = start(53);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        p.children.push(statement());
        return end();
    }
    function elseIfClause2() {
        let p = start(53);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        expect(87);
        p.children.push(statementList([24, 18, 19]));
        return end();
    }
    function elseClause1() {
        let p = start(52);
        next();
        p.children.push(statement());
        return end();
    }
    function elseClause2() {
        let p = start(52);
        next();
        expect(87);
        p.children.push(statementList([24]));
        return end();
    }
    function isElseIfClauseStart(t) {
        return t.tokenType === 19;
    }
    function ifStatement() {
        let p = start(96);
        next();
        expect(118);
        p.children.push(expression(0));
        expect(121);
        let t = peek();
        let elseIfClauseFunction = elseIfClause1;
        let elseClauseFunction = elseClause1;
        let expectEndIf = false;
        if (t.tokenType === 87) {
            next();
            p.children.push(statementList([19, 18, 24]));
            elseIfClauseFunction = elseIfClause2;
            elseClauseFunction = elseClause2;
            expectEndIf = true;
        }
        else if (isStatementStart(t)) {
            p.children.push(statement());
        }
        else {
            error();
        }
        if (peek().tokenType === 19) {
            p.children.push(list(54, elseIfClauseFunction, isElseIfClauseStart));
        }
        if (peek().tokenType === 18) {
            p.children.push(elseClauseFunction());
        }
        if (expectEndIf) {
            expect(24);
            expect(88);
        }
        return end();
    }
    function expressionStatement() {
        let p = start(73);
        p.children.push(expression(0));
        expect(88);
        return end();
    }
    function returnType() {
        let p = start(149);
        next();
        p.children.push(typeDeclaration());
        return end();
    }
    function typeDeclaration() {
        let p = start(173);
        optional(96);
        switch (peek().tokenType) {
            case 6:
            case 3:
                next();
                break;
            case 83:
            case 51:
            case 147:
                p.children.push(qualifiedName());
                break;
            default:
                error();
                break;
        }
        return end();
    }
    function classConstDeclaration(p) {
        p.phraseType = 25;
        next();
        p.children.push(delimitedList(27, classConstElement, isClassConstElementStartToken, 93, [88]));
        expect(88);
        return end();
    }
    function isExpressionStart(t) {
        switch (t.tokenType) {
            case 84:
            case 90:
            case 3:
            case 117:
            case 78:
            case 147:
            case 83:
            case 51:
            case 118:
            case 60:
            case 135:
            case 129:
            case 111:
            case 143:
            case 89:
            case 86:
            case 94:
            case 152:
            case 153:
            case 150:
            case 155:
            case 151:
            case 148:
            case 149:
            case 47:
            case 11:
            case 52:
            case 79:
            case 82:
            case 73:
            case 72:
            case 71:
            case 77:
            case 75:
            case 74:
            case 76:
            case 10:
            case 154:
            case 97:
            case 95:
            case 53:
            case 69:
            case 70:
            case 35:
            case 41:
            case 42:
            case 57:
            case 58:
            case 28:
            case 20:
            case 46:
            case 29:
                return true;
            default:
                return false;
        }
    }
    function classConstElement() {
        let p = start(26);
        p.children.push(identifier());
        expect(85);
        p.children.push(expression(0));
        return end();
    }
    function isPropertyElementStart(t) {
        return t.tokenType === 84;
    }
    function propertyDeclaration(p) {
        let t;
        p.phraseType = 137;
        p.children.push(delimitedList(139, propertyElement, isPropertyElementStart, 93, [88]));
        expect(88);
        return end();
    }
    function propertyElement() {
        let p = start(138);
        expect(84);
        if (peek().tokenType === 85) {
            p.children.push(propertyInitialiser());
        }
        return end();
    }
    function propertyInitialiser() {
        let p = start(140);
        next();
        p.children.push(expression(0));
        return end();
    }
    function memberModifierList() {
        let p = start(110);
        while (isMemberModifier(peek())) {
            next();
        }
        return end();
    }
    function isMemberModifier(t) {
        switch (t.tokenType) {
            case 55:
            case 56:
            case 54:
            case 60:
            case 2:
            case 31:
                return true;
            default:
                return false;
        }
    }
    function qualifiedNameList(breakOn) {
        return delimitedList(142, qualifiedName, isQualifiedNameStart, 93, breakOn);
    }
    function objectCreationExpression() {
        let p = start(128);
        next();
        if (peek().tokenType === 9) {
            p.children.push(anonymousClassDeclaration());
            return end();
        }
        p.children.push(typeDesignator(34));
        if (optional(118)) {
            if (isArgumentStart(peek())) {
                p.children.push(argumentList());
            }
            expect(121);
        }
        return end();
    }
    function typeDesignator(phraseType) {
        let p = start(phraseType);
        let part = classTypeDesignatorAtom();
        while (true) {
            switch (peek().tokenType) {
                case 117:
                    part = subscriptExpression(part, 120);
                    continue;
                case 116:
                    part = subscriptExpression(part, 119);
                    continue;
                case 115:
                    part = propertyAccessExpression(part);
                    continue;
                case 133:
                    let staticPropNode = start(152);
                    staticPropNode.children.push(part);
                    next();
                    staticPropNode.children.push(restrictedScopedMemberName());
                    part = end();
                    continue;
                default:
                    break;
            }
            break;
        }
        p.children.push(part);
        return end();
    }
    function restrictedScopedMemberName() {
        let p = start(151);
        let t = peek();
        switch (t.tokenType) {
            case 84:
                next();
                break;
            case 90:
                p.children.push(simpleVariable());
                break;
            default:
                error();
                break;
        }
        return end();
    }
    function classTypeDesignatorAtom() {
        let t = peek();
        switch (t.tokenType) {
            case 60:
                return relativeScope();
            case 84:
            case 90:
                return simpleVariable();
            case 83:
            case 51:
            case 147:
                return qualifiedName();
            default:
                start(62);
                error();
                return end();
        }
    }
    function cloneExpression() {
        let p = start(35);
        next();
        p.children.push(expression(0));
        return end();
    }
    function listIntrinsic() {
        let p = start(108);
        next();
        expect(118);
        p.children.push(arrayInitialiserList(121));
        expect(121);
        return end();
    }
    function unaryExpression(phraseType) {
        let p = start(phraseType);
        let op = next();
        switch (phraseType) {
            case 133:
            case 134:
                p.children.push(variable(variableAtom()));
                break;
            default:
                p.children.push(expression(precedenceAssociativityTuple(op)[0]));
                break;
        }
        return end();
    }
    function anonymousFunctionHeader() {
        let p = start(5);
        optional(60);
        next();
        optional(103);
        expect(118);
        if (isParameterStart(peek())) {
            p.children.push(delimitedList(130, parameterDeclaration, isParameterStart, 93, [121]));
        }
        expect(121);
        if (peek().tokenType === 66) {
            p.children.push(anonymousFunctionUseClause());
        }
        if (peek().tokenType === 87) {
            p.children.push(returnType());
        }
        return end();
    }
    function anonymousFunctionCreationExpression() {
        let p = start(4);
        p.children.push(anonymousFunctionHeader());
        p.children.push(functionDeclarationBody());
        return end();
    }
    function isAnonymousFunctionUseVariableStart(t) {
        return t.tokenType === 84 ||
            t.tokenType === 103;
    }
    function anonymousFunctionUseClause() {
        let p = start(6);
        next();
        expect(118);
        p.children.push(delimitedList(36, anonymousFunctionUseVariable, isAnonymousFunctionUseVariableStart, 93, [121]));
        expect(121);
        return end();
    }
    function anonymousFunctionUseVariable() {
        let p = start(7);
        optional(103);
        expect(84);
        return end();
    }
    function isTypeDeclarationStart(t) {
        switch (t.tokenType) {
            case 147:
            case 83:
            case 51:
            case 96:
            case 3:
            case 6:
                return true;
            default:
                return false;
        }
    }
    function parameterDeclaration() {
        let p = start(129);
        if (isTypeDeclarationStart(peek())) {
            p.children.push(typeDeclaration());
        }
        optional(103);
        optional(134);
        expect(84);
        if (peek().tokenType === 85) {
            next();
            p.children.push(expression(0));
        }
        return end();
    }
    function variable(variableAtomNode) {
        let count = 0;
        while (true) {
            ++count;
            switch (peek().tokenType) {
                case 133:
                    variableAtomNode = scopedAccessExpression(variableAtomNode);
                    continue;
                case 115:
                    variableAtomNode = propertyOrMethodAccessExpression(variableAtomNode);
                    continue;
                case 117:
                    variableAtomNode = subscriptExpression(variableAtomNode, 120);
                    continue;
                case 116:
                    variableAtomNode = subscriptExpression(variableAtomNode, 119);
                    continue;
                case 118:
                    variableAtomNode = functionCallExpression(variableAtomNode);
                    continue;
                default:
                    if (count === 1 && variableAtomNode.phraseType !== 156) {
                        let errNode = start(67, true);
                        errNode.children.push(variableAtomNode);
                        error();
                        return end();
                    }
                    break;
            }
            break;
        }
        return variableAtomNode;
    }
    function functionCallExpression(lhs) {
        let p = start(85, true);
        p.children.push(lhs);
        expect(118);
        if (isArgumentStart(peek())) {
            p.children.push(argumentList());
        }
        expect(121);
        return end();
    }
    function scopedAccessExpression(lhs) {
        let p = start(65, true);
        p.children.push(lhs);
        next();
        p.children.push(scopedMemberName(p));
        if (optional(118)) {
            p.phraseType = 150;
            if (isArgumentStart(peek())) {
                p.children.push(argumentList());
            }
            expect(121);
            return end();
        }
        else if (p.phraseType === 150) {
            error();
        }
        return end();
    }
    function scopedMemberName(parent) {
        let p = start(151);
        let t = peek();
        switch (t.tokenType) {
            case 116:
                parent.phraseType = 150;
                p.children.push(encapsulatedExpression(116, 119));
                break;
            case 84:
                parent.phraseType = 152;
                next();
                break;
            case 90:
                p.children.push(simpleVariable());
                parent.phraseType = 152;
                break;
            default:
                if (t.tokenType === 83 || isSemiReservedToken(t)) {
                    p.children.push(identifier());
                    parent.phraseType = 24;
                }
                else {
                    error();
                }
                break;
        }
        return end();
    }
    function propertyAccessExpression(lhs) {
        let p = start(136, true);
        p.children.push(lhs);
        next();
        p.children.push(memberName());
        return end();
    }
    function propertyOrMethodAccessExpression(lhs) {
        let p = start(136, true);
        p.children.push(lhs);
        next();
        p.children.push(memberName());
        if (optional(118)) {
            if (isArgumentStart(peek())) {
                p.children.push(argumentList());
            }
            p.phraseType = 112;
            expect(121);
        }
        return end();
    }
    function memberName() {
        let p = start(111);
        switch (peek().tokenType) {
            case 83:
                next();
                break;
            case 116:
                p.children.push(encapsulatedExpression(116, 119));
                break;
            case 90:
            case 84:
                p.children.push(simpleVariable());
                break;
            default:
                error();
                break;
        }
        return end();
    }
    function subscriptExpression(lhs, closeTokenType) {
        let p = start(160, true);
        p.children.push(lhs);
        next();
        if (isExpressionStart(peek())) {
            p.children.push(expression(0));
        }
        expect(closeTokenType);
        return end();
    }
    function argumentList() {
        return delimitedList(8, argumentExpression, isArgumentStart, 93, [121]);
    }
    function isArgumentStart(t) {
        return t.tokenType === 134 || isExpressionStart(t);
    }
    function variadicUnpacking() {
        let p = start(178);
        next();
        p.children.push(expression(0));
        return end();
    }
    function argumentExpression() {
        return peek().tokenType === 134 ?
            variadicUnpacking() : expression(0);
    }
    function qualifiedName() {
        let p = start(141);
        let t = peek();
        if (t.tokenType === 147) {
            next();
            p.phraseType = 84;
        }
        else if (t.tokenType === 51) {
            p.phraseType = 144;
            next();
            expect(147);
        }
        p.children.push(namespaceName());
        return end();
    }
    function isQualifiedNameStart(t) {
        switch (t.tokenType) {
            case 147:
            case 83:
            case 51:
                return true;
            default:
                return false;
        }
    }
    function shortArrayCreationExpression() {
        let p = start(9);
        next();
        if (isArrayElementStart(peek())) {
            p.children.push(arrayInitialiserList(120));
        }
        expect(120);
        return end();
    }
    function longArrayCreationExpression() {
        let p = start(9);
        next();
        expect(118);
        if (isArrayElementStart(peek())) {
            p.children.push(arrayInitialiserList(121));
        }
        expect(121);
        return end();
    }
    function isArrayElementStart(t) {
        return t.tokenType === 103 || isExpressionStart(t);
    }
    function arrayInitialiserList(breakOn) {
        let p = start(11);
        let t;
        let arrayInitialiserListRecoverSet = [breakOn, 93];
        recoverSetStack.push(arrayInitialiserListRecoverSet);
        while (true) {
            if (isArrayElementStart(peek())) {
                p.children.push(arrayElement());
            }
            t = peek();
            if (t.tokenType === 93) {
                next();
            }
            else if (t.tokenType === breakOn) {
                break;
            }
            else {
                error();
                if (isArrayElementStart(t)) {
                    continue;
                }
                else {
                    defaultSyncStrategy();
                    t = peek();
                    if (t.tokenType === 93 || t.tokenType === breakOn) {
                        continue;
                    }
                }
                break;
            }
        }
        recoverSetStack.pop();
        return end();
    }
    function arrayValue() {
        let p = start(13);
        optional(103);
        p.children.push(expression(0));
        return end();
    }
    function arrayKey() {
        let p = start(12);
        p.children.push(expression(0));
        return end();
    }
    function arrayElement() {
        let p = start(10);
        if (peek().tokenType === 103) {
            p.children.push(arrayValue());
            return end();
        }
        let keyOrValue = arrayKey();
        p.children.push(keyOrValue);
        if (!optional(132)) {
            keyOrValue.phraseType = 13;
            return end();
        }
        p.children.push(arrayValue());
        return end();
    }
    function encapsulatedExpression(openTokenType, closeTokenType) {
        let p = start(56);
        expect(openTokenType);
        p.children.push(expression(0));
        expect(closeTokenType);
        return end();
    }
    function relativeScope() {
        let p = start(145);
        next();
        return end();
    }
    function variableAtom() {
        let t = peek();
        switch (t.tokenType) {
            case 84:
            case 90:
                return simpleVariable();
            case 118:
                return encapsulatedExpression(118, 121);
            case 3:
                return longArrayCreationExpression();
            case 117:
                return shortArrayCreationExpression();
            case 78:
                return next(true);
            case 60:
                return relativeScope();
            case 83:
            case 51:
            case 147:
                return qualifiedName();
            default:
                let p = start(68);
                error();
                return end();
        }
    }
    function simpleVariable() {
        let p = start(156);
        let t = expectOneOf([84, 90]);
        if (t && t.tokenType === 90) {
            t = peek();
            if (t.tokenType === 116) {
                p.children.push(encapsulatedExpression(116, 119));
            }
            else if (t.tokenType === 90 || t.tokenType === 84) {
                p.children.push(simpleVariable());
            }
            else {
                error();
            }
        }
        return end();
    }
    function haltCompilerStatement() {
        let p = start(93);
        next();
        expect(118);
        expect(121);
        expect(88);
        return end();
    }
    function namespaceUseDeclaration() {
        let p = start(124);
        next();
        optionalOneOf([35, 12]);
        optional(147);
        let nsNameNode = namespaceName();
        let t = peek();
        if (t.tokenType === 147 || t.tokenType === 116) {
            p.children.push(nsNameNode);
            expect(147);
            expect(116);
            p.children.push(delimitedList(126, namespaceUseGroupClause, isNamespaceUseGroupClauseStartToken, 93, [119]));
            expect(119);
            expect(88);
            return end();
        }
        p.children.push(delimitedList(123, namespaceUseClauseFunction(nsNameNode), isNamespaceUseClauseStartToken, 93, [88], true));
        expect(88);
        return end();
    }
    function isNamespaceUseClauseStartToken(t) {
        return t.tokenType === 83 || t.tokenType === 147;
    }
    function namespaceUseClauseFunction(nsName) {
        return () => {
            let p = start(122, !!nsName);
            if (nsName) {
                p.children.push(nsName);
                nsName = null;
            }
            else {
                p.children.push(namespaceName());
            }
            if (peek().tokenType === 4) {
                p.children.push(namespaceAliasingClause());
            }
            return end();
        };
    }
    function delimitedList(phraseType, elementFunction, elementStartPredicate, delimiter, breakOn, doNotPushHiddenToParent) {
        let p = start(phraseType, doNotPushHiddenToParent);
        let t;
        let delimitedListRecoverSet = breakOn ? breakOn.slice(0) : [];
        delimitedListRecoverSet.push(delimiter);
        recoverSetStack.push(delimitedListRecoverSet);
        while (true) {
            p.children.push(elementFunction());
            t = peek();
            if (t.tokenType === delimiter) {
                next();
            }
            else if (!breakOn || breakOn.indexOf(t.tokenType) >= 0) {
                break;
            }
            else {
                error();
                if (elementStartPredicate(t)) {
                    continue;
                }
                else if (breakOn) {
                    defaultSyncStrategy();
                    if (peek().tokenType === delimiter) {
                        continue;
                    }
                }
                break;
            }
        }
        recoverSetStack.pop();
        return end();
    }
    function isNamespaceUseGroupClauseStartToken(t) {
        switch (t.tokenType) {
            case 12:
            case 35:
            case 83:
                return true;
            default:
                return false;
        }
    }
    function namespaceUseGroupClause() {
        let p = start(125);
        optionalOneOf([35, 12]);
        p.children.push(namespaceName());
        if (peek().tokenType === 4) {
            p.children.push(namespaceAliasingClause());
        }
        return end();
    }
    function namespaceAliasingClause() {
        let p = start(119);
        next();
        expect(83);
        return end();
    }
    function namespaceDefinition() {
        let p = start(120);
        next();
        if (peek().tokenType === 83) {
            p.children.push(namespaceName());
            let t = expectOneOf([88, 116]);
            if (!t || t.tokenType !== 116) {
                return end();
            }
        }
        else {
            expect(116);
        }
        p.children.push(statementList([119]));
        expect(119);
        return end();
    }
    function namespaceName() {
        start(121);
        expect(83);
        while (true) {
            if (peek().tokenType === 147 &&
                peek(1).tokenType === 83) {
                next();
                next();
            }
            else {
                break;
            }
        }
        return end();
    }
    function isReservedToken(t) {
        switch (t.tokenType) {
            case 41:
            case 42:
            case 28:
            case 57:
            case 58:
            case 49:
            case 50:
            case 48:
            case 43:
            case 52:
            case 11:
            case 29:
            case 39:
            case 19:
            case 18:
            case 24:
            case 17:
            case 16:
            case 68:
            case 26:
            case 33:
            case 22:
            case 34:
            case 23:
            case 14:
            case 21:
            case 4:
            case 64:
            case 8:
            case 32:
            case 62:
            case 66:
            case 44:
            case 36:
            case 67:
            case 65:
            case 46:
            case 20:
            case 13:
            case 37:
            case 35:
            case 12:
            case 59:
            case 53:
            case 69:
            case 47:
            case 61:
            case 25:
            case 7:
            case 15:
            case 5:
            case 3:
            case 6:
            case 30:
            case 40:
            case 51:
            case 63:
            case 45:
            case 9:
            case 10:
            case 77:
            case 74:
            case 75:
            case 73:
            case 72:
            case 71:
            case 76:
                return true;
            default:
                return false;
        }
    }
    function isSemiReservedToken(t) {
        switch (t.tokenType) {
            case 60:
            case 2:
            case 31:
            case 54:
            case 56:
            case 55:
                return true;
            default:
                return isReservedToken(t);
        }
    }
    function isStatementStart(t) {
        switch (t.tokenType) {
            case 51:
            case 66:
            case 38:
            case 12:
            case 35:
            case 9:
            case 2:
            case 31:
            case 63:
            case 45:
            case 116:
            case 39:
            case 68:
            case 16:
            case 33:
            case 61:
            case 5:
            case 13:
            case 59:
            case 36:
            case 60:
            case 17:
            case 65:
            case 34:
            case 14:
            case 64:
            case 62:
            case 37:
            case 83:
            case 88:
            case 158:
            case 81:
            case 156:
            case 157:
                return true;
            default:
                return isExpressionStart(t);
        }
    }
})(Parser = exports.Parser || (exports.Parser = {}));
