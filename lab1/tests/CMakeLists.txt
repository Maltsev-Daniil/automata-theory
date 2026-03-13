find_package(Catch2 3 REQUIRED)

set(TEST_SOURCES
        test.cpp
)

add_executable(test ${TEST_SOURCES})

target_link_libraries(test
        PRIVATE
        Catch2::Catch2WithMain
        regex_recognizer
        lexer_recognizer
        smc_recognizer
        domain
)

target_include_directories(test
        PUBLIC
        ${PROJECT_SOURCE_DIR}
)

add_compile_options(test
        -O0
        -g
        --coverage
)

add_link_options(test
        --coverage
)

include(CTest)
include(Catch)
catch_discover_tests(test)