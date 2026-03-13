//
// Created by nitsir on 3/7/26.
//

#include "smc.h"
#include "../src/field.h"

std::optional<Fields> SmcRec::process(const std::string& input)
{
    context.reset();

    size_t pos = 0;

    // надо добавить в конец строки пробел
    // для грамотного перехода в дефолтное состояние
    const std::string input_local =  input + "  ";

    while (pos < input_local.size())
    {
        if (context.reconsume)
        {
            pos--;
            context.reconsume = false;
        }

        char c = input_local[pos++];
        context.fsm.curr_char(c);
    }

    context.fsm.eof();

    if (!context.success)
        return std::nullopt;

    return context.type == Type::create
        ? Fields{context.type, context.name_, context.attributes}
        : Fields{context.type, context.name_, context.join_lhs, context.join_rhs};
}
