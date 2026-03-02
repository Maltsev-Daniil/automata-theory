#include "lexer.h"
#include "regex.h"
#include <iostream>
#include <limits>
#include <memory>
#include <string>
#include "relation.h"
#include "factory.h"

void printRelations(
    const std::unordered_map<std::string,
    std::unique_ptr<Relation>>& relations) {
    for (const auto& [name, relation_ptr] : relations) {
        if (!relation_ptr) continue;
        std::cout << name << " (";
        const auto& attrs = relation_ptr->attributes;
        for (std::size_t i = 0; i < attrs.size(); ++i) {
            std::cout << attrs[i];
            if (i + 1 < attrs.size())
                std::cout << ", ";
        }
        std::cout << ")\n";
    }
}

int main() {
    std::unordered_map<std::string, std::unique_ptr<Relation>> relations;
    Factory factory(relations);

    RegexRec regex_rec;
    LexerRec lexer_rec;

    while (true) {
        std::cout << "\n===== MENU =====\n";
        std::cout << "Input source:\n";
        std::cout << "  1 - console\n";
        std::cout << "  2 - file\n";
        std::cout << "Processor:\n";
        std::cout << "  1 - regex\n";
        std::cout << "  2 - lexer\n";
        std::cout << "  3 - smc\n";
        std::cout << "Enter option (XY): ";

        std::string option;
        getline(std::cin, option);

        if (option[0] == '0' || option == "exit") {
            std::cout << "\nExiting...\n";
            break;
        }

        if (option.size() != 2) {
            std::cout << "Error: option must be 2 digits (XY)\n";
            continue;
        }

        char input_mode = option[0];
        char process_mode = option[1];

        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');

        std::string input;

        if (input_mode == '1') {
            std::cout << "Enter text (empty line to finish):\n";
            std::string line;

            while (std::getline(std::cin, line)) {
                if (line.empty())
                    break;
                input += line;
            }
        }
        else if (input_mode == '2') {
            std::cout << "// file input not implemented yet\n";
            continue;
        }
        else {
            std::cout << "Error: invalid input source\n";
            continue;
        }

        std::optional<Fields> result;

        if (process_mode == '1') {
            result = regex_rec.process(input);
        }
        else if (process_mode == '2') {
            result = lexer_rec.process(input);
        }
        else if (process_mode == '3') {
            std::cout << "// smc not implemented yet\n";
            continue;
        }
        else {
            std::cout << "Error: invalid processor\n";
            continue;
        }

        if (result.has_value()) {
            auto relation = factory.createRelation(result.value());
            if (relation) {
                relations.insert({relation->name, std::move(relation)});
            }
        } else {
            std::cout << "Error: invalid relations\n";
        }

        printRelations(relations);
    }

    return 0;
}