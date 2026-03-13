//
// Created by nitsir on 2/28/26.
//

#include "regex.h"

#include <iostream>

#include "../src/field.h"

#include <regex>

std::optional<Fields> RegexRec::process(const std::string& input) {
    static const std::string identifier = R"([a-zA-Z][a-zA-Z0-9._]*)";

    static const std::string pattern_prefix =
        "^create\\s+(" + identifier + ")\\s*";

    static const std::string pattern_create_postfix =
        R"(\s*\((.*)\)\s*$)";

    static const std::string pattern_join_postfix =
        "as\\s+(" + identifier + ")\\s+join\\s+(" + identifier + ")\\s*$";

    static const std::regex prefix_pattern(pattern_prefix, std::regex::optimize);
    static const std::regex join_pattern(pattern_join_postfix, std::regex::optimize);
    static const std::regex create_pattern(pattern_create_postfix, std::regex::optimize);

    std::smatch match;
    std::string relation_name;
    if (!std::regex_search(input, match, prefix_pattern))
        return std::nullopt;

    relation_name = match[1].str();

    size_t pos = match.length();
    std::string tail = input.substr(pos);

    if (std::regex_match(tail, match, join_pattern)) {
        std::string lhs = match[1];
        std::string rhs = match[2];
        return Fields(Type::join, relation_name, lhs, rhs);
    } else {
        size_t lparen = tail.find('(');
        size_t rparen = tail.rfind(')');

        if (lparen == std::string::npos || rparen == std::string::npos || lparen >= rparen)
            return std::nullopt;

        for (size_t i = rparen + 1; i < tail.size(); ++i) {
            if (!isspace(tail[i]))
                return std::nullopt;
        }

        std::string attributes = tail.substr(lparen + 1, rparen - lparen - 1);

        std::vector<std::string> attribute_list = tokenizeAtt(attributes);

        if (attribute_list.empty())
            return std::nullopt;

        return Fields(Type::create, relation_name, attribute_list);
    }
}


std::vector<std::string> RegexRec::tokenizeAtt(const std::string& input) {
    if (!input.empty() && input.back() == ',')
        return {};

    std::vector<std::string> tokens;
    static const std::regex pattern(R"(^[a-zA-Z][a-zA-Z0-9._]*$)", std::regex::optimize);
    std::stringstream ss(input);
    std::string token;

    while (std::getline(ss, token, ',')) {
        size_t start = token.find_first_not_of(" \t");
        size_t end = token.find_last_not_of(" \t");

        if (start == std::string::npos) {
            return {};
        }

        std::string trimmed = token.substr(start, end - start + 1);

        if (!std::regex_match(trimmed, pattern)) {
            return {};
        }

        tokens.push_back(trimmed);
    }

    return tokens;
}