#pragma once
#include <string>
#include <vector>

enum class Type {
    join,
    create,
    none
};

struct Fields {
    Fields(Type type, std::string name, std::vector<std::string> attribute_list)
        : type(type), name(name), attributes(attribute_list) {}

    Fields(Type type, std::string name, std::string lhs, std::string rhs)
        : type(type), name(name), join_lhs(lhs), join_rhs(rhs) {}

    Type type;
    std::string name;

    std::vector<std::string> attributes;

    std::string join_lhs;
    std::string join_rhs;
};
