cmake_minimum_required(VERSION 3.5)
project(lab1 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)

add_executable(lab1_out
        main.cpp
        regex.cpp
        lexer.cpp
        lex.yy.cc
)