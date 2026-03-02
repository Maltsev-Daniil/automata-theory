#pragma once
#include <memory>
#include <unordered_map>

#include "field.h"

struct Relation;

// factory also does validation
class Factory {
    public:
    Factory(std::unordered_map<std::string, std::unique_ptr<Relation>>& relations)
        : relations_(relations) {}
    std::unique_ptr<Relation> createRelation(const Fields& fields);
private:
    std::unordered_map<std::string, std::unique_ptr<Relation>>& relations_;

    std::unique_ptr<Relation> createNew(
        const std::string& name,
        const std::vector<std::string>& attr);

    std::unique_ptr<Relation> createJoin(
        const std::string& name,
        const std::string& lhs,
        const std::string& rhs);

    std::vector<std::string> concatenateNUnite(
        const std::string &lhs, const std::string &rhs);
};
