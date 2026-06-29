// sudoku.js
#!/usr/bin/env node
'use strict';

const fs = require('fs');
const readline = require('readline');

const COLORS = {
    reset: '\x1b[0m',
    red: '\x1b[91m',
    green: '\x1b[92m',
    yellow: '\x1b[93m',
    blue: '\x1b[94m',
    cyan: '\x1b[96m',
    gray: '\x1b[90m',
    bold: '\x1b[1m'
};

function colorize(text, color) {
    return COLORS[color] + text + COLORS.reset;
}

class Sudoku {
    constructor(board = null) {
        this.board = board ? board.map(row => [...row]) : Array.from({ length: 9 }, () => Array(9).fill(0));
        this.steps = 0;
    }

    render(highlight = null) {
        const lines = [];
        for (let i = 0; i < 9; i++) {
            if (i % 3 === 0 && i > 0) {
                lines.push(colorize('┃━━━┃━━━┃━━━┃', 'gray'));
            }
            const row = [];
            for (let j = 0; j < 9; j++) {
                if (j % 3 === 0 && j > 0) {
                    row.push(colorize('┃', 'gray'));
                }
                const val = this.board[i][j];
                if (val === 0) {
                    row.push(' ');
                } else {
                    let col = 'green';
                    if (highlight && highlight.some(h => h[0] === i && h[1] === j)) {
                        col = 'cyan';
                    }
                    row.push(colorize(String(val), col));
                }
            }
            lines.push(row.join(' '));
        }
        return lines.join('\n');
    }

    parseString(str) {
        if (str.length !== 81) throw new Error('Строка должна содержать ровно 81 символ');
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                const ch = str[i*9 + j];
                if (ch >= '1' && ch <= '9') {
                    this.board[i][j] = parseInt(ch, 10);
                } else if (ch === '0') {
                    this.board[i][j] = 0;
                } else {
                    throw new Error(`Недопустимый символ: ${ch}`);
                }
            }
        }
    }

    exportString() {
        let s = '';
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                s += this.board[i][j] === 0 ? '0' : String(this.board[i][j]);
            }
        }
        return s;
    }

    isValid() {
        // Строки
        for (let i = 0; i < 9; i++) {
            const seen = new Set();
            for (let j = 0; j < 9; j++) {
                const val = this.board[i][j];
                if (val !== 0) {
                    if (val < 1 || val > 9 || seen.has(val)) return false;
                    seen.add(val);
                }
            }
        }
        // Столбцы
        for (let j = 0; j < 9; j++) {
            const seen = new Set();
            for (let i = 0; i < 9; i++) {
                const val = this.board[i][j];
                if (val !== 0) {
                    if (seen.has(val)) return false;
                    seen.add(val);
                }
            }
        }
        // Блоки
        for (let br = 0; br < 3; br++) {
            for (let bc = 0; bc < 3; bc++) {
                const seen = new Set();
                for (let i = 0; i < 3; i++) {
                    for (let j = 0; j < 3; j++) {
                        const val = this.board[br*3+i][bc*3+j];
                        if (val !== 0) {
                            if (seen.has(val)) return false;
                            seen.add(val);
                        }
                    }
                }
            }
        }
        return true;
    }

    getCandidates(row, col) {
        if (this.board[row][col] !== 0) return [];
        const used = new Set();
        for (let j = 0; j < 9; j++) if (this.board[row][j] !== 0) used.add(this.board[row][j]);
        for (let i = 0; i < 9; i++) if (this.board[i][col] !== 0) used.add(this.board[i][col]);
        const br = Math.floor(row/3)*3, bc = Math.floor(col/3)*3;
        for (let i = 0; i < 3; i++) {
            for (let j = 0; j < 3; j++) {
                const val = this.board[br+i][bc+j];
                if (val !== 0) used.add(val);
            }
        }
        const cands = [];
        for (let v = 1; v <= 9; v++) if (!used.has(v)) cands.push(v);
        return cands;
    }

    findBestEmpty() {
        let bestR = -1, bestC = -1, bestCands = null, minCands = 10;
        for (let i = 0; i < 9; i++) {
            for (let j = 0; j < 9; j++) {
                if (this.board[i][j] === 0) {
                    const cands = this.getCandidates(i, j);
                    if (cands.length < minCands) {
                        minCands = cands.length;
                        bestR = i; bestC = j; bestCands = cands;
                        if (minCands === 1) return [bestR, bestC, bestCands];
                    }
                }
            }
        }
        return [bestR, bestC, bestCands];
    }

    solve(animate = false, delay = 0.1) {
        this.steps = 0;
        const start = Date.now();
        const _solve = () => {
            this.steps++;
            const [r, c, cands] = this.findBestEmpty();
            if (r === -1) return true;
            for (const val of cands) {
                this.board[r][c] = val;
                if (animate) this.printAnimated(r, c, delay);
                if (_solve()) return true;
                this.board[r][c] = 0;
            }
            return false;
        };
        const solved = _solve();
        const elapsed = (Date.now() - start) / 1000;
        return [solved, elapsed];
    }

    printAnimated(row, col, delay) {
        console.clear();
        console.log(this.render([[row, col]]));
        console.log(colorize(`Шаг: ${this.steps}`, 'yellow'));
        // Задержка
        const stop = Date.now() + delay * 1000;
        while (Date.now() < stop) {}
    }

    countSolutions(maxCount = 100) {
        let count = 0;
        const _count = () => {
            if (count >= maxCount) return;
            const [r, c, cands] = this.findBestEmpty();
            if (r === -1) { count++; return; }
            for (const val of cands) {
                this.board[r][c] = val;
                _count();
                this.board[r][c] = 0;
                if (count >= maxCount) return;
            }
        };
        _count();
        return count;
    }

    static generate(level = 'easy') {
        const levels = { easy: [30, 35], medium: [40, 45], hard: [50, 55] };
        const [minEmpty, maxEmpty] = levels[level] || levels.easy;
        // Заполняем доску
        const s = new Sudoku();
        s.solveInternal();
        const filled = s.board.map(row => [...row]);
        const target = minEmpty + Math.floor(Math.random() * (maxEmpty - minEmpty + 1));
        let removed = 0, attempts = 0;
        while (removed < target && attempts < 10000) {
            attempts++;
            const i = Math.floor(Math.random() * 9);
            const j = Math.floor(Math.random() * 9);
            if (filled[i][j] !== 0) {
                const backup = filled[i][j];
                filled[i][j] = 0;
                const test = new Sudoku(filled);
                if (test.countSolutions(2) === 1) {
                    removed++;
                } else {
                    filled[i][j] = backup;
                }
            }
        }
        return new Sudoku(filled);
    }

    solveInternal() {
        const _solve = () => {
            const [r, c, cands] = this.findBestEmpty();
            if (r === -1) return true;
            for (const val of cands) {
                this.board[r][c] = val;
                if (_solve()) return true;
                this.board[r][c] = 0;
            }
            return false;
        };
        _solve();
    }
}

function loadBoard(source) {
    let content;
    if (fs.existsSync(source)) {
        content = fs.readFileSync(source, 'utf8');
    } else {
        content = source;
    }
    const digits = content.split('').filter(ch => ch >= '0' && ch <= '9');
    if (digits.length !== 81) throw new Error('Должно быть ровно 81 цифра');
    const s = new Sudoku();
    s.parseString(digits.join(''));
    return s;
}

function main() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.log(colorize('Usage: node sudoku.js <generate|solve|check|count|export> [options]', 'yellow'));
        console.log('  generate -l <easy|medium|hard>');
        console.log('  solve -i <file|string> [-a] [--delay <sec>]');
        console.log('  check -i <file|string>');
        console.log('  count -i <file|string> [--max <N>]');
        console.log('  export -i <file|string>');
        process.exit(1);
    }

    const cmd = args[0];
    const opts = {};
    for (let i = 1; i < args.length; i++) {
        if (args[i] === '-l' && i+1 < args.length) opts.level = args[++i];
        else if (args[i] === '-i' && i+1 < args.length) opts.input = args[++i];
        else if (args[i] === '-o' && i+1 < args.length) opts.output = args[++i];
        else if (args[i] === '-a') opts.animate = true;
        else if (args[i] === '--delay' && i+1 < args.length) opts.delay = parseFloat(args[++i]);
        else if (args[i] === '--max' && i+1 < args.length) opts.max = parseInt(args[++i]);
    }

    try {
        switch (cmd) {
            case 'generate': {
                const level = opts.level || 'easy';
                const s = Sudoku.generate(level);
                if (opts.output) fs.writeFileSync(opts.output, s.exportString());
                console.log(s.render());
                const empty = s.board.flat().filter(v => v === 0).length;
                console.log(colorize(`Уровень: ${level}, пустых ячеек: ${empty}`, 'yellow'));
                break;
            }
            case 'solve': {
                if (!opts.input) throw new Error('Укажите входную доску через -i');
                const s = loadBoard(opts.input);
                if (!s.isValid()) throw new Error('Доска невалидна!');
                const [solved, elapsed] = s.solve(opts.animate, opts.delay || 0.1);
                if (solved) {
                    if (opts.output) fs.writeFileSync(opts.output, s.exportString());
                    console.log(s.render());
                    console.log(colorize(`Решение найдено за ${elapsed.toFixed(3)} сек, шагов: ${s.steps}`, 'green'));
                } else {
                    console.log(colorize('Решения не существует!', 'red'));
                }
                break;
            }
            case 'check': {
                if (!opts.input) throw new Error('Укажите входную доску через -i');
                const s = loadBoard(opts.input);
                console.log(s.isValid() ? colorize('✅ Доска корректна.', 'green') : colorize('❌ Доска невалидна.', 'red'));
                break;
            }
            case 'count': {
                if (!opts.input) throw new Error('Укажите входную доску через -i');
                const s = loadBoard(opts.input);
                const maxCount = opts.max || 100;
                const count = s.countSolutions(maxCount);
                if (count >= maxCount) {
                    console.log(colorize(`Количество решений >= ${maxCount} (ограничено)`, 'yellow'));
                } else {
                    console.log(colorize(`Количество решений: ${count}`, 'green'));
                }
                break;
            }
            case 'export': {
                if (!opts.input) throw new Error('Укажите входную доску через -i');
                const s = loadBoard(opts.input);
                console.log(s.exportString());
                break;
            }
            default:
                console.log(colorize(`Неизвестная команда: ${cmd}`, 'red'));
        }
    } catch (err) {
        console.log(colorize(`Ошибка: ${err.message}`, 'red'));
        process.exit(1);
    }
}

main();
