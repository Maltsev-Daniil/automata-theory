cmake_minimum_required(VERSION 3.20)

project(lab1)

set(CMAKE_CXX_STANDARD 20)

add_subdirectory(tests)

add_library(domain
        src/factory.cpp
)

target_include_directories(domain PUBLIC src)


add_library(regex_recognizer
        regex/regex.cpp
)

target_include_directories(regex_recognizer PUBLIC
        regex
        src
)

target_link_libraries(regex_recognizer domain)


add_library(lexer_recognizer
        lexer/lexer.cpp
        lexer/lex.yy.cc
)

target_include_directories(lexer_recognizer PUBLIC
        lexer
        src
)

target_link_libraries(lexer_recognizer domain)


add_library(smc_recognizer
        smc/smc.cpp
        smc/smc_sm.cpp
        smc/smc-fields.cpp
)

target_include_directories(smc_recognizer PUBLIC
        smc
        src
)

target_link_libraries(smc_recognizer domain)


add_executable(main_app
        main.cpp
)

target_link_libraries(main_app
        regex_recognizer
        lexer_recognizer
        smc_recognizer
)


add_executable(timing_app
        timing.cpp
)

target_link_libraries(timing_app
        regex_recognizer
        lexer_recognizer
        smc_recognizer
)