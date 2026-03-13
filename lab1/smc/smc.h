#pragma once

#include "smc-fields.h"
#include "../src/irecognizer.h"

class SmcRec : public IRecognizer {
    public:
    SmcRec() : context(SmcFields()) {}
    std::optional<Fields> process(const std::string &input_local) override;

    private:
    SmcFields context;
};