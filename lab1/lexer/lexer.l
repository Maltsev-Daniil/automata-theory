%option noyywrap c++
%{
#include "field.h"
#include <string>
#include "tokens.h"

std::string curr_ident{};
%}

identifier [a-zA-Z][a-zA-Z0-9._]*

%%

"create"        { return CREATE; }
"as"            { return AS; }
"join"          { return JOIN; }

"("             { return LPAREN; }
")"             { return RPAREN; }
","             { return COMMA; }

{identifier}    {
    curr_ident = std::string(yytext, yyleng);
    return TOKEN_IDENT; }

[ \t\n]+        {}
.               { return UNKNOWN; }

%%