#pragma once
#include <string>
#include <vector>

struct Relation {
    Relation(std::string name, std::vector<std::string> attr)
        : name(name), attributes(attr) {}
    std::string name;
    std::vector<std::string> attributes;
};
