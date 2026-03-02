//
// Created by nitsir on 3/1/26.
//
#include "lexer.h"
#include <vector>
#include <string>
#include "field.h"
#include "tokens.h"
#include <sstream>

#include <FlexLexer.h>
extern std::string curr_ident;

std::optional<Fields> LexerRec::process(const std::string &input) {
    yyFlexLexer lexer;
    std::istringstream iss(input);
    yy_buffer_state* buff = lexer.yy_create_buffer(iss, input.size());
    lexer.yy_switch_to_buffer(buff);

    if (CREATE != lexer.yylex()) {
        lexer.yy_delete_buffer(buff);
        return std::nullopt;
    }
    if (TOKEN_IDENT != lexer.yylex()) {
        lexer.yy_delete_buffer(buff);
        return std::nullopt;
    }

    const std::string ident = curr_ident;
    int curr_token = lexer.yylex();

    std::optional<Fields> result;
    if (LPAREN == curr_token) {
        result = matchAtt(lexer, ident);
    }
    else if (AS == curr_token) {
        result = matchJoin(lexer, ident);
    }
    else {
        result = std::nullopt;
    }
    lexer.yy_delete_buffer(buff);
    return result;
}

std::optional<Fields> LexerRec::matchJoin(yyFlexLexer& lexer, const std::string& ident) {
    if (TOKEN_IDENT != lexer.yylex()) {
        return std::nullopt;
    }
    const std::string lhs = curr_ident;
    if (JOIN != lexer.yylex()) {
        return std::nullopt;
    }
    if (TOKEN_IDENT != lexer.yylex()) {
        return std::nullopt;
    }
    const std::string rhs = curr_ident;
    return Fields(Type::join, ident, lhs, rhs);
}

std::optional<Fields> LexerRec::matchAtt(yyFlexLexer& lexer, const std::string &ident) {
    std::vector<std::string> attributes;
    while (true) {
        if (TOKEN_IDENT != lexer.yylex()) {
            return std::nullopt;
        }
        attributes.push_back(curr_ident);
        int curr_token = lexer.yylex();
        if (RPAREN == curr_token) {
            break;
        }
        if (COMMA != curr_token) {
            return std::nullopt;
        }
    }
    return Fields(Type::create, ident, attributes);
}