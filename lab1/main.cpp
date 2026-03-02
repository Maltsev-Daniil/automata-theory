#include "field.h"
#include "lexer.h"
#include "regex.h"
#include <iostream>
#include <string>

void printFields(const Fields& fields) {
    std::cout << "=== Fields ===\n"
              << "Name: " << fields.name << "\n"
              << "Attributes: ";

    if (fields.attributes.empty()) {
        std::cout << "(empty)";
    } else {
        for (size_t i = 0; i < fields.attributes.size(); ++i) {
            if (i > 0) std::cout << ", ";
            std::cout << fields.attributes[i];
        }
    }

    std::cout << "\nJoin LHS: " << (fields.join_lhs.empty() ? "(none)" : fields.join_lhs)
              << "\nJoin RHS: " << (fields.join_rhs.empty() ? "(none)" : fields.join_rhs)
              << std::endl;
}

int main() {
    RegexRec regex;
    LexerRec lexer;
    std::string line, input;
    while (std::getline(std::cin, line)) {
        if (line.empty()) break;
        input.append(line);
    }
    std::cout << "input done" << std::endl;

    auto res_r = regex.process(input);
    auto res_l = lexer.process(input);

    std::cout << "regex" << std::endl;
    if (res_r != std::nullopt)
        printFields(res_r.value());
    std::cout << "lex" << std::endl;
    if (res_l != std::nullopt)
        printFields(res_l.value());

    return 0;
}
