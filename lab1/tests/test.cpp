#define CATCH_CONFIG_MAIN
#include <catch2/catch_all.hpp>

#include "../regex/regex.h"
#include "../smc/smc.h"
#include "../lexer/lexer.h"

#include <memory>
#include <vector>
#include <string>

struct TestCase
{
std::string input;

bool success;

Type type;

std::string name;

std::vector<std::string> attributes;

std::string lhs;

std::string rhs;

};

static std::vector<TestCase> test_cases =
{

// CREATE

{
    "create a(b)",
    true,
    Type::create,
    "a",
    {"b"},
    "",
    ""
},

{
    "create users(id,name,age)",
    true,
    Type::create,
    "users",
    {"id","name","age"},
    "",
    ""
},

{
    "create relation(a, b , c)",
    true,
    Type::create,
    "relation",
    {"a","b","c"},
    "",
    ""
},

// JOIN

{
    "create r as a join b",
    true,
    Type::join,
    "r",
    {},
    "a",
    "b"
},

{
    "create result as table1 join table2",
    true,
    Type::join,
    "result",
    {},
    "table1",
    "table2"
},

// INVALID

{
    "create",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a(",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a(    ,   b  )",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a(  b   ,  )",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a as b",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a(a, ,b )",
    false,
    Type::none,
    "",
    {},
    "",
    ""
},

{
    "create a(a b)",
    false,
    Type::none,
    "",
    {},
    "",
    ""
}

};

struct RecognizerEntry
{
std::string name;
std::function<std::unique_ptr<IRecognizer>()> factory;
};

static std::vector<RecognizerEntry> recognizers =
{
{
"regex",
[] { return std::make_unique<RegexRec>(); }
},

{
    "smc",
    [] { return std::make_unique<SmcRec>(); }
},

{
    "dfa",
    [] { return std::make_unique<LexerRec>(); }
}

};

TEST_CASE("All recognizers parse commands identically")
{
for (const auto& rec_entry : recognizers)
{
auto recognizer = rec_entry.factory();

    for (const auto& tc : test_cases)
    {
        DYNAMIC_SECTION(rec_entry.name << " input: " << tc.input)
        {
            auto result = recognizer->process(tc.input);

            if (!tc.success)
            {
                REQUIRE_FALSE(result.has_value());
                continue;
            }

            REQUIRE(result.has_value());

            const auto& fields = result.value();

            REQUIRE(fields.type == tc.type);
            REQUIRE(fields.name == tc.name);

            if (tc.type == Type::create)
            {
                REQUIRE(fields.attributes == tc.attributes);
            }

            if (tc.type == Type::join)
            {
                REQUIRE(fields.join_lhs == tc.lhs);
                REQUIRE(fields.join_rhs == tc.rhs);
            }
        }
    }
}

}