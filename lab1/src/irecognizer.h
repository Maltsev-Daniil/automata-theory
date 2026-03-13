#pragma once
#include <optional>
#include <string>

struct Fields;

class IRecognizer {
    public:
    virtual ~IRecognizer() = default;
    virtual std::optional<Fields> process(const std::string& input) = 0;
};
