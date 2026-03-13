//
// Created by nitsir on 3/2/26.
//

#include "factory.h"
#include "relation.h"

std::unique_ptr<Relation> Factory::createRelation(const Fields& fields) {
    if (relations_.contains(fields.name))
        return nullptr;
    switch (fields.type) {
        case Type::create:
            return createNew(fields.name, fields.attributes);
        case Type::join:
            return createJoin(
                fields.name,
                fields.join_lhs,
                fields.join_rhs);
    }

    throw std::invalid_argument("createRelation: unknown Type");
}

std::unique_ptr<Relation> Factory::createNew(
    const std::string &name, const std::vector<std::string> &attr) {
    if (name.empty() || attr.empty()) return nullptr;
    return std::make_unique<Relation>(name, attr);
}

std::unique_ptr<Relation> Factory::createJoin(
    const std::string &name, const std::string &lhs, const std::string &rhs) {
    if (name.empty() || lhs.empty() || rhs.empty()) return nullptr;
    if (!relations_.contains(lhs) || !relations_.contains(rhs)) return nullptr;
    return std::make_unique<Relation>(name, concatenateNUnite(lhs, rhs));
}

std::vector<std::string> Factory::concatenateNUnite(
    const std::string &lhs, const std::string &rhs) {
    std::vector<std::string> result;

    auto lhs_it = relations_.find(lhs);
    auto rhs_it = relations_.find(rhs);

    for (const auto& att : lhs_it->second->attributes) {
        result.push_back(lhs + "." + att);
    }
    for (const auto& att : rhs_it->second->attributes) {
        result.push_back(rhs + "." + att);
    }
    return result;
}