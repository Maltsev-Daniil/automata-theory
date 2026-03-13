//
// Created by nitsir on 3/13/26.
//

#include <chrono>
#include <fstream>

#include "lexer/lexer.h"
#include "regex/regex.h"
#include <iostream>
#include <memory>
#include <string>

#include "smc/smc.h"
#include "src/relation.h"
#include "src/factory.h"

int main() {
    RegexRec regex_rec;
    LexerRec lexer_rec;
    SmcRec smc_rec;

    std::ifstream input_file("timing_data.txt");
    std::ofstream output_file("timing_result.txt");

    if (!input_file.is_open()) {
        std::cout << "Error: cannot open timing_data.txt\n";
        return 1;
    }

    if (!output_file.is_open()) {
        std::cout << "Error: cannot open timing_result.txt\n";
        return 1;
    }

    std::string line;
    const int runs = 250;

    // прогрев
    regex_rec.process("create a(b,c)");

    while (std::getline(input_file, line)) {

        long long regex_total = 0;
        long long lexer_total = 0;
        long long smc_total = 0;

        for (int i = 0; i < runs; ++i) {

            {
                auto start = std::chrono::high_resolution_clock::now();
                std::optional<Fields> res = regex_rec.process(line);
                auto end = std::chrono::high_resolution_clock::now();

                regex_total += std::chrono::duration_cast<
                    std::chrono::nanoseconds>(end - start).count();
            }

            {
                auto start = std::chrono::high_resolution_clock::now();
                std::optional<Fields> res = lexer_rec.process(line);
                auto end = std::chrono::high_resolution_clock::now();

                lexer_total += std::chrono::duration_cast<
                    std::chrono::nanoseconds>(end - start).count();
            }

            {
                auto start = std::chrono::high_resolution_clock::now();
                std::optional<Fields> res = smc_rec.process(line);
                auto end = std::chrono::high_resolution_clock::now();

                smc_total += std::chrono::duration_cast<
                    std::chrono::nanoseconds>(end - start).count();
            }
        }

        long long regex_avg = regex_total / runs;
        long long lexer_avg = lexer_total / runs;
        long long smc_avg = smc_total / runs;

        output_file << regex_avg << ";"
                    << lexer_avg << ";"
                    << smc_avg << "\n";
    }

    input_file.close();
    output_file.close();

    return 0;
}