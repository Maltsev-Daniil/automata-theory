//
// Created by nitsir on 3/11/26.
//

#include "smc-fields.h"
#include "smc_sm.h"

SmcFields::SmcFields()
    : fsm(*this)
{
    fsm.enterStartState();
}

void SmcFields::reset()
{
    buffer.clear();

    role = IdentRole::Init;

    success = false;
    reconsume = false;

    type = Type::none;

    name_.clear();
    attributes.clear();

    join_lhs.clear();
    join_rhs.clear();

    fsm.setState(SmcMap::Waiting);
    fsm.enterStartState();
}