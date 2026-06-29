// sudoku.cpp
#include <iostream>
#include <vector>
#include <string>
#include <fstream>
#include <sstream>
#include <random>
#include <chrono>
#include <thread>
#include <algorithm>
#include <functional>

using namespace std;

const string RESET = "\033[0m";
const string RED = "\033[91m";
const string GREEN = "\033[92m";
const string YELLOW = "\033[93m";
const string BLUE = "\033[94m";
const string CYAN = "\033[96m";
const string GRAY = "\033[90m";
const string BOLD = "\033[1m";

string colorize(const string& text, const string& color) {
    return color + text + RESET;
}

class Sudoku {
private:
    int board[9][9];
    int steps;

public:
    Sudoku() : steps(0) {
        for (int i = 0; i < 9; ++i)
            for (int j = 0; j < 9; ++j)
                board[i][j] = 0;
    }

    string render(const vector<pair<int,int>>& highlight = {}) const {
        ostringstream result;
        for (int i = 0; i < 9; ++i) {
            if (i % 3 == 0 && i > 0)
                result << colorize("┃━━━┃━━━┃━━━┃", GRAY) << "\n";
            for (int j = 0; j < 9; ++j) {
                if (j % 3 == 0 && j > 0)
                    result << colorize("┃", GRAY);
                int val = board[i][j];
                if (val == 0) {
                    result << " ";
                } else {
                    string col = GREEN;
                    for (auto& h : highlight) {
                        if (h.first == i && h.second == j) {
                            col = CYAN;
                            break;
                        }
                    }
                    result << colorize(to_string(val), col);
                }
                if (j < 8) result << " ";
            }
            result << "\n";
        }
        return result.str();
    }

    void parseString(const string& s) {
        if (s.size() != 81) throw runtime_error("Строка должна содержать 81 символ");
        for (int i = 0; i < 9; ++i) {
            for (int j = 0; j < 9; ++j) {
                char ch = s[i*9 + j];
                if (ch >= '1' && ch <= '9') board[i][j] = ch - '0';
                else if (ch == '0') board[i][j] = 0;
                else throw runtime_error("Недопустимый символ");
            }
        }
    }

    string exportString() const {
        string s;
        for (int i = 0; i < 9; ++i)
            for (int j = 0; j < 9; ++j)
                s += board[i][j] == 0 ? '0' : char('0' + board[i][j]);
        return s;
    }

    bool isValid() const {
        // Строки
        for (int i = 0; i < 9; ++i) {
            bool seen[10] = {false};
            for (int j = 0; j < 9; ++j) {
                int val = board[i][j];
                if (val != 0) {
                    if (val < 1 || val > 9 || seen[val]) return false;
                    seen[val] = true;
                }
            }
        }
        // Столбцы
        for (int j = 0; j < 9; ++j) {
            bool seen[10] = {false};
            for (int i = 0; i < 9; ++i) {
                int val = board[i][j];
                if (val != 0) {
                    if (seen[val]) return false;
                    seen[val] = true;
                }
            }
        }
        // Блоки
        for (int br = 0; br < 3; ++br) {
            for (int bc = 0; bc < 3; ++bc) {
                bool seen[10] = {false};
                for (int i = 0; i < 3; ++i) {
                    for (int j = 0; j < 3; ++j) {
                        int val = board[br*3 + i][bc*3 + j];
                        if (val != 0) {
                            if (seen[val]) return false;
                            seen[val] = true;
                        }
                    }
                }
            }
        }
        return true;
    }

    vector<int> getCandidates(int row, int col) const {
        if (board[row][col] != 0) return {};
        bool used[10] = {false};
        for (int j = 0; j < 9; ++j) if (board[row][j] != 0) used[board[row][j]] = true;
        for (int i = 0; i < 9; ++i) if (board[i][col] != 0) used[board[i][col]] = true;
        int br = row / 3 * 3, bc = col / 3 * 3;
        for (int i = 0; i < 3; ++i)
            for (int j = 0; j < 3; ++j)
                if (board[br+i][bc+j] != 0) used[board[br+i][bc+j]] = true;
        vector<int> cands;
        for (int v = 1; v <= 9; ++v) if (!used[v]) cands.push_back(v);
        return cands;
    }

    bool findBestEmpty(int& r, int& c, vector<int>& cands) const {
        int bestR = -1, bestC = -1, minCands = 10;
        for (int i = 0; i < 9; ++i) {
            for (int j = 0; j < 9; ++j) {
                if (board[i][j] == 0) {
                    auto cnd = getCandidates(i, j);
                    if ((int)cnd.size() < minCands) {
                        minCands = cnd.size();
                        bestR = i; bestC = j; cands = cnd;
                        if (minCands == 1) { r = bestR; c = bestC; return true; }
                    }
                }
            }
        }
        if (bestR != -1) { r = bestR; c = bestC; return true; }
        return false;
    }

    bool solveInternal() {
        int r, c; vector<int> cands;
        if (!findBestEmpty(r, c, cands)) return true;
        for (int val : cands) {
            board[r][c] = val;
            if (solveInternal()) return true;
            board[r][c] = 0;
        }
        return false;
    }

    bool solve(bool animate = false, double delay = 0.1) {
        steps = 0;
        auto start = chrono::steady_clock::now();
        function<bool()> solveFn = [&]() -> bool {
            steps++;
            int r, c; vector<int> cands;
            if (!findBestEmpty(r, c, cands)) return true;
            for (int val : cands) {
                board[r][c] = val;
                if (animate) {
                    cout << "\033[2J\033[1;1H";
                    cout << render({{r, c}});
                    cout << colorize("Шаг: " + to_string(steps), YELLOW) << endl;
                    this_thread::sleep_for(chrono::milliseconds((int)(delay*1000)));
                }
                if (solveFn()) return true;
                board[r][c] = 0;
            }
            return false;
        };
        bool solved = solveFn();
        auto elapsed = chrono::duration<double>(chrono::steady_clock::now() - start).count();
        if (animate) cout << render();
        return solved;
    }

    int countSolutions(int maxCount = 100) {
        int count = 0;
        function<void()> countFn = [&]() {
            if (count >= maxCount) return;
            int r, c; vector<int> cands;
            if (!findBestEmpty(r, c, cands)) { count++; return; }
            for (int val : cands) {
                board[r][c] = val;
                countFn();
                board[r][c] = 0;
                if (count >= maxCount) return;
            }
        };
        countFn();
        return count;
    }

    void generate(string level = "easy") {
        // Заполняем доску
        solveInternal();
        int filled[9][9];
        for (int i = 0; i < 9; ++i)
            for (int j = 0; j < 9; ++j)
                filled[i][j] = board[i][j];

        int minEmpty, maxEmpty;
        if (level == "easy") { minEmpty = 30; maxEmpty = 35; }
        else if (level == "medium") { minEmpty = 40; maxEmpty = 45; }
        else if (level == "hard") { minEmpty = 50; maxEmpty = 55; }
        else { minEmpty = 30; maxEmpty = 35; }

        random_device rd;
        mt19937 gen(rd());
        uniform_int_distribution<> dist(0, 8);
        int target = minEmpty + (gen() % (maxEmpty - minEmpty + 1));
        int removed = 0, attempts = 0;
        while (removed < target && attempts < 10000) {
            attempts++;
            int i = dist(gen), j = dist(gen);
            if (filled[i][j] != 0) {
                int backup = filled[i][j];
                filled[i][j] = 0;
                Sudoku test;
                for (int x = 0; x < 9; ++x)
                    for (int y = 0; y < 9; ++y)
                        test.board[x][y] = filled[x][y];
                if (test.countSolutions(2) == 1) {
                    removed++;
                } else {
                    filled[i][j] = backup;
                }
            }
        }
        for (int i = 0; i < 9; ++i)
            for (int j = 0; j < 9; ++j)
                board[i][j] = filled[i][j];
    }

    friend Sudoku loadBoard(const string& source);
};

Sudoku loadBoard(const string& source) {
    string content;
    ifstream f(source);
    if (f) {
        stringstream buf;
        buf << f.rdbuf();
        content = buf.str();
    } else {
        content = source;
    }
    string digits;
    for (char ch : content)
        if (ch >= '0' && ch <= '9') digits += ch;
    if (digits.size() != 81) throw runtime_error("Должно быть 81 цифра");
    Sudoku s;
    s.parseString(digits);
    return s;
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        cout << colorize("Usage: sudoku <generate|solve|check|count|export> [options]", YELLOW) << endl;
        cout << "  generate -l <easy|medium|hard>" << endl;
        cout << "  solve -i <file|string> [-a] [--delay <sec>]" << endl;
        cout << "  check -i <file|string>" << endl;
        cout << "  count -i <file|string> [--max <N>]" << endl;
        cout << "  export -i <file|string>" << endl;
        return 1;
    }

    string cmd = argv[1];
    string input, output, level = "easy";
    bool animate = false;
    double delay = 0.1;
    int maxCount = 100;

    for (int i = 2; i < argc; ++i) {
        string arg = argv[i];
        if (arg == "-i" && i+1 < argc) input = argv[++i];
        else if (arg == "-o" && i+1 < argc) output = argv[++i];
        else if (arg == "-l" && i+1 < argc) level = argv[++i];
        else if (arg == "-a") animate = true;
        else if (arg == "--delay" && i+1 < argc) delay = stod(argv[++i]);
        else if (arg == "--max" && i+1 < argc) maxCount = stoi(argv[++i]);
        else if (arg == "-h") { /* help */ }
    }

    try {
        if (cmd == "generate") {
            Sudoku s;
            s.generate(level);
            if (!output.empty()) {
                ofstream f(output);
                f << s.exportString();
            }
            cout << s.render();
            int empty = 0;
            for (int i = 0; i < 9; ++i)
                for (int j = 0; j < 9; ++j)
                    if (s.board[i][j] == 0) empty++;
            cout << colorize("Уровень: " + level + ", пустых ячеек: " + to_string(empty), YELLOW) << endl;
        } else if (cmd == "solve") {
            if (input.empty()) throw runtime_error("Укажите входную доску через -i");
            auto s = loadBoard(input);
            if (!s.isValid()) throw runtime_error("Доска невалидна!");
            bool solved = s.solve(animate, delay);
            if (solved) {
                if (!output.empty()) {
                    ofstream f(output);
                    f << s.exportString();
                }
                cout << s.render();
                cout << colorize("Решение найдено!", GREEN) << endl;
            } else {
                cout << colorize("Решения не существует!", RED) << endl;
            }
        } else if (cmd == "check") {
            if (input.empty()) throw runtime_error("Укажите входную доску через -i");
            auto s = loadBoard(input);
            cout << (s.isValid() ? colorize("✅ Доска корректна.", GREEN) : colorize("❌ Доска невалидна.", RED)) << endl;
        } else if (cmd == "count") {
            if (input.empty()) throw runtime_error("Укажите входную доску через -i");
            auto s = loadBoard(input);
            int count = s.countSolutions(maxCount);
            if (count >= maxCount) {
                cout << colorize("Количество решений >= " + to_string(maxCount) + " (ограничено)", YELLOW) << endl;
            } else {
                cout << colorize("Количество решений: " + to_string(count), GREEN) << endl;
            }
        } else if (cmd == "export") {
            if (input.empty()) throw runtime_error("Укажите входную доску через -i");
            auto s = loadBoard(input);
            cout << s.exportString() << endl;
        } else {
            cout << colorize("Неизвестная команда: " + cmd, RED) << endl;
        }
    } catch (const exception& e) {
        cout << colorize("Ошибка: " + string(e.what()), RED) << endl;
        return 1;
    }
    return 0;
}
