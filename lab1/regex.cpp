//
// Created by nitsir on 2/28/26.
//

#include "regex.h"
#include "field.h"

#include <regex>

std::optional<Fields> RegexRec::process(const std::string& input) {
    static const std::string identifier = R"([a-zA-Z][a-zA-Z0-9._]*)";
    static const std::string identifier_list =
        identifier + "(\\s*,\\s*" + identifier + ")*";

    static const std::string pattern1 =
        "^create\\s+(" + identifier + ")\\s*\\((" + identifier_list + ")\\s*\\)\\s*$";
    static const std::string pattern2 =
        "^create\\s+(" + identifier + ")\\s+as\\s+(" + identifier + ")\\s+join\\s+(" + identifier + ")\\s*$";

    static const std::regex pat1(pattern1);
    static const std::regex pat2(pattern2);

    std::smatch match;
    if (std::regex_match(input, match, pat1)) {
        std::string relation_name = match[1];
        std::vector<std::string> attribute_list = tokenizeAtt(match[2]);
        return Fields(Type::create, relation_name, attribute_list);
    } else if (std::regex_match(input, match, pat2)) {
        std::string relation_name = match[1];
        std::string lhs = match[2];
        std::string rhs = match[3];
        return Fields(Type::join, relation_name, lhs, rhs);
    } else {
        return std::nullopt;
    }
}

std::vector<std::string> RegexRec::tokenizeAtt(const std::string& input) {
    static const std::regex pattern(R"(\s*,\s*)");
    std::sregex_token_iterator it(input.begin(), input.end(), pattern, -1);
    std::vector<std::string> tokens;
    for (; it != std::sregex_token_iterator(); ++it) {
        tokens.push_back(*it);
    }
    return tokens;
}