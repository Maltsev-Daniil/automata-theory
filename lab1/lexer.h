#pragma once
#include "irecognizer.h"

class yyFlexLexer;

class LexerRec : public IRecognizer {
    public:
    std::optional<Fields> process(const std::string &input) override;
    private:
    std::optional<Fields> matchAtt(yyFlexLexer& lexer, const std::string& ident);
    std::optional<Fields> matchJoin(yyFlexLexer& lexer, const std::string& ident);
};