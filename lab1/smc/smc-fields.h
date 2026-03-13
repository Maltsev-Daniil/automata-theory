#pragma once

#include "../src/field.h"
#include "smc_sm.h"
#include <string>
#include <vector>

enum IdentRole { Init, Lhs, Rhs, Att, None };

class SmcClass;

class SmcFields {
public:
    SmcFields();
    SmcClass fsm;

    std::string buffer;
    IdentRole role;
    bool success = false;
    bool reconsume = false;
    Type type = Type::none;
    std::string name_;
    std::vector<std::string> attributes;
    std::string join_lhs;
    std::string join_rhs;

    void reset();

    void clearBuffer() { buffer.clear(); }
    void pushBuffer(char c) { buffer.push_back(c); }
    void setBuffer(const std::string& v) { buffer = v; }

    void setRole(IdentRole r) { role = r; }
    void setSuccess(bool v) { success = v; }
    void setReconsume(bool v) { reconsume = v; }

    void setType(Type t) { type = t; }
    void setName(const std::string& n) { name_ = n; }

    void pushAttribute(const std::string& a) { attributes.push_back(a); }
    void clearAttributes() { attributes.clear(); }

    void setJoinLhs(const std::string& s) { join_lhs = s; }
    void setJoinRhs(const std::string& s) { join_rhs = s; }

    [[nodiscard]] Fields buildCreate() const { return Fields{type, name_, attributes}; }
    [[nodiscard]] Fields buildJoin() const { return Fields{type, name_, join_lhs, join_rhs}; }
};