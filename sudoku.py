# sudoku.py
#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
import random
import time
import json
import argparse
from copy import deepcopy
from pathlib import Path

# ANSI colors
COLORS = {
    'reset': '\033[0m',
    'red': '\033[91m',
    'green': '\033[92m',
    'yellow': '\033[93m',
    'blue': '\033[94m',
    'cyan': '\033[96m',
    'gray': '\033[90m',
    'bold': '\033[1m'
}

def colorize(text, color):
    return f"{COLORS.get(color, '')}{text}{COLORS['reset']}"

class Sudoku:
    def __init__(self, board=None):
        if board is None:
            self.board = [[0]*9 for _ in range(9)]
        else:
            self.board = [row[:] for row in board]
        self.size = 9
        self.box_size = 3
        self.steps = 0

    def __str__(self):
        return self.render()

    def render(self, highlight=None):
        """Выводит доску с цветами и разделителями блоков."""
        result = []
        for i in range(9):
            if i % 3 == 0 and i > 0:
                result.append(colorize('┃' + '━' * 9 + '┃' + '━' * 9 + '┃' + '━' * 9 + '┃', 'gray'))
            row = []
            for j in range(9):
                if j % 3 == 0 and j > 0:
                    row.append(colorize('┃', 'gray'))
                val = self.board[i][j]
                if val == 0:
                    row.append(' ')
                else:
                    color = 'cyan' if highlight and (i, j) in highlight else 'green'
                    row.append(colorize(str(val), color))
            result.append(' '.join(row))
        return '\n'.join(result)

    def parse_string(self, s):
        """Заполняет доску из строки из 81 символа (0 для пустых)."""
        if len(s) != 81:
            raise ValueError("Строка должна содержать ровно 81 символ")
        for i in range(9):
            for j in range(9):
                ch = s[i*9 + j]
                if ch.isdigit():
                    self.board[i][j] = int(ch)
                else:
                    self.board[i][j] = 0

    def export_string(self):
        """Экспортирует доску в строку из 81 символа."""
        return ''.join(str(self.board[i][j]) if self.board[i][j] != 0 else '0'
                       for i in range(9) for j in range(9))

    def is_valid(self):
        """Проверяет, является ли доска корректным Судоку."""
        # Проверка строк
        for i in range(9):
            seen = set()
            for j in range(9):
                val = self.board[i][j]
                if val != 0:
                    if val < 1 or val > 9 or val in seen:
                        return False
                    seen.add(val)
        # Проверка столбцов
        for j in range(9):
            seen = set()
            for i in range(9):
                val = self.board[i][j]
                if val != 0:
                    if val in seen:
                        return False
                    seen.add(val)
        # Проверка блоков 3×3
        for block_i in range(3):
            for block_j in range(3):
                seen = set()
                for i in range(3):
                    for j in range(3):
                        val = self.board[block_i*3 + i][block_j*3 + j]
                        if val != 0:
                            if val in seen:
                                return False
                            seen.add(val)
        return True

    def find_empty(self):
        """Находит первую пустую ячейку (для бэктрекинга)."""
        for i in range(9):
            for j in range(9):
                if self.board[i][j] == 0:
                    return (i, j)
        return None

    def find_best_empty(self):
        """Находит ячейку с минимальным количеством кандидатов (MRV)."""
        min_candidates = 10
        best = None
        for i in range(9):
            for j in range(9):
                if self.board[i][j] == 0:
                    candidates = self.get_candidates(i, j)
                    if len(candidates) < min_candidates:
                        min_candidates = len(candidates)
                        best = (i, j)
                        if min_candidates == 1:
                            return best
        return best

    def get_candidates(self, row, col):
        """Возвращает возможные значения для ячейки."""
        if self.board[row][col] != 0:
            return []
        used = set()
        # Строка
        for j in range(9):
            if self.board[row][j] != 0:
                used.add(self.board[row][j])
        # Столбец
        for i in range(9):
            if self.board[i][col] != 0:
                used.add(self.board[i][col])
        # Блок 3×3
        br, bc = row // 3 * 3, col // 3 * 3
        for i in range(3):
            for j in range(3):
                val = self.board[br+i][bc+j]
                if val != 0:
                    used.add(val)
        return [v for v in range(1, 10) if v not in used]

    def solve(self, animate=False, delay=0.1):
        """Решает Судоку с помощью бэктрекинга с MRV."""
        self.steps = 0
        start_time = time.time()

        def _solve():
            self.steps += 1
            empty = self.find_best_empty()
            if empty is None:
                return True
            row, col = empty
            for val in self.get_candidates(row, col):
                self.board[row][col] = val
                if animate:
                    self.print_animated(row, col)
                    time.sleep(delay)
                if _solve():
                    return True
                self.board[row][col] = 0
            return False

        solved = _solve()
        elapsed = time.time() - start_time
        return solved, elapsed

    def print_animated(self, row, col):
        """Вывод доски с подсветкой текущей ячейки."""
        os.system('clear' if os.name == 'posix' else 'cls')
        print(self.render(highlight=[(row, col)]))
        print(colorize(f"Шаг: {self.steps}", 'yellow'))

    def count_solutions(self, max_count=100):
        """Подсчитывает количество решений (до max_count)."""
        count = 0

        def _count():
            nonlocal count
            if count >= max_count:
                return
            empty = self.find_best_empty()
            if empty is None:
                count += 1
                return
            row, col = empty
            for val in self.get_candidates(row, col):
                self.board[row][col] = val
                _count()
                self.board[row][col] = 0
                if count >= max_count:
                    return

        _count()
        return count

    @staticmethod
    def generate(level='easy'):
        """Генерирует новое Судоку с уникальным решением."""
        # Уровни: easy: 30–35 пустых, medium: 40–45, hard: 50–55
        levels = {'easy': (30, 35), 'medium': (40, 45), 'hard': (50, 55)}
        min_empty, max_empty = levels.get(level, (30, 35))

        # Создаём заполненную доску
        board = [[0]*9 for _ in range(9)]
        sudoku = Sudoku(board)
        sudoku.solve()
        filled = [row[:] for row in sudoku.board]

        # Удаляем ячейки
        target_empty = random.randint(min_empty, max_empty)
        removed = 0
        attempts = 0
        while removed < target_empty and attempts < 10000:
            attempts += 1
            i, j = random.randint(0, 8), random.randint(0, 8)
            if filled[i][j] != 0:
                backup = filled[i][j]
                filled[i][j] = 0
                # Проверяем уникальность решения
                test = Sudoku([row[:] for row in filled])
                if test.count_solutions(2) == 1:
                    removed += 1
                else:
                    filled[i][j] = backup

        return Sudoku([row[:] for row in filled])

def main():
    parser = argparse.ArgumentParser(description="Sudoku Generator & Solver")
    subparsers = parser.add_subparsers(dest='command', help='Команда')

    # generate
    gen_parser = subparsers.add_parser('generate', help='Генерация Судоку')
    gen_parser.add_argument('-l', '--level', choices=['easy', 'medium', 'hard'], default='easy')
    gen_parser.add_argument('-o', '--output', help='Сохранить в файл')

    # solve
    solve_parser = subparsers.add_parser('solve', help='Решение Судоку')
    solve_parser.add_argument('-i', '--input', required=True, help='Файл или строка 81 символа')
    solve_parser.add_argument('-o', '--output', help='Сохранить в файл')
    solve_parser.add_argument('-a', '--animate', action='store_true', help='Анимация решения')
    solve_parser.add_argument('--delay', type=float, default=0.1, help='Задержка анимации (сек)')

    # check
    check_parser = subparsers.add_parser('check', help='Проверка корректности')
    check_parser.add_argument('-i', '--input', required=True, help='Файл или строка 81 символа')

    # count
    count_parser = subparsers.add_parser('count', help='Подсчёт решений')
    count_parser.add_argument('-i', '--input', required=True, help='Файл или строка 81 символа')
    count_parser.add_argument('--max', type=int, default=100, help='Максимальное количество решений')

    # export
    export_parser = subparsers.add_parser('export', help='Экспорт в строку')
    export_parser.add_argument('-i', '--input', required=True, help='Файл или строка 81 символа')

    args = parser.parse_args()

    def load_board(source):
        if os.path.exists(source):
            with open(source, 'r') as f:
                content = f.read().strip()
        else:
            content = source
        # Удаляем пробелы и разделители
        content = ''.join(c for c in content if c.isdigit() or c == '0')
        if len(content) != 81:
            raise ValueError("Должно быть ровно 81 цифра")
        sudoku = Sudoku()
        sudoku.parse_string(content)
        return sudoku

    if args.command == 'generate':
        sudoku = Sudoku.generate(args.level)
        if args.output:
            with open(args.output, 'w') as f:
                f.write(sudoku.export_string())
        print(sudoku.render())
        print(colorize(f"Уровень: {args.level}, пустых ячеек: {sum(1 for row in sudoku.board for v in row if v == 0)}", 'yellow'))

    elif args.command == 'solve':
        sudoku = load_board(args.input)
        if not sudoku.is_valid():
            print(colorize("Доска невалидна!", 'red'))
            sys.exit(1)
        solved, elapsed = sudoku.solve(args.animate, args.delay)
        if solved:
            if args.output:
                with open(args.output, 'w') as f:
                    f.write(sudoku.export_string())
            print(sudoku.render())
            print(colorize(f"Решение найдено за {elapsed:.3f} сек, шагов: {sudoku.steps}", 'green'))
        else:
            print(colorize("Решения не существует!", 'red'))

    elif args.command == 'check':
        sudoku = load_board(args.input)
        if sudoku.is_valid():
            print(colorize("✅ Доска корректна.", 'green'))
        else:
            print(colorize("❌ Доска невалидна.", 'red'))

    elif args.command == 'count':
        sudoku = load_board(args.input)
        count = sudoku.count_solutions(args.max)
        if count >= args.max:
            print(colorize(f"Количество решений >= {args.max} (ограничено)", 'yellow'))
        else:
            print(colorize(f"Количество решений: {count}", 'green'))

    elif args.command == 'export':
        sudoku = load_board(args.input)
        print(sudoku.export_string())

    else:
        parser.print_help()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print(colorize("\nПрервано.", 'yellow'))
        sys.exit(0)
