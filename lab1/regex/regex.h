#pragma once

#include <vector>

#include "../src/irecognizer.h"

class RegexRec : public IRecognizer {
    public:
    std::optional<Fields> process(const std::string& input) override;
private:
    std::vector<std::string> tokenizeAtt(const std::string& input);
};