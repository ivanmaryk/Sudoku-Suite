// sudoku.java
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.*;

public class sudoku {
    private static final String RESET = "\u001B[0m";
    private static final String RED = "\u001B[91m";
    private static final String GREEN = "\u001B[92m";
    private static final String YELLOW = "\u001B[93m";
    private static final String BLUE = "\u001B[94m";
    private static final String CYAN = "\u001B[96m";
    private static final String GRAY = "\u001B[90m";
    private static final String BOLD = "\u001B[1m";

    private static String colorize(String text, String color) {
        return color + text + RESET;
    }

    static class SudokuBoard {
        int[][] board = new int[9][9];
        int steps;

        public SudokuBoard() {}

        public SudokuBoard(int[][] src) {
            for (int i = 0; i < 9; i++)
                System.arraycopy(src[i], 0, board[i], 0, 9);
        }

        public String render(List<int[]> highlight) {
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < 9; i++) {
                if (i % 3 == 0 && i > 0)
                    sb.append(colorize("┃━━━┃━━━┃━━━┃", GRAY)).append("\n");
                for (int j = 0; j < 9; j++) {
                    if (j % 3 == 0 && j > 0)
                        sb.append(colorize("┃", GRAY));
                    int val = board[i][j];
                    if (val == 0) {
                        sb.append(" ");
                    } else {
                        String col = GREEN;
                        if (highlight != null) {
                            for (int[] h : highlight) {
                                if (h[0] == i && h[1] == j) {
                                    col = CYAN;
                                    break;
                                }
                            }
                        }
                        sb.append(colorize(String.valueOf(val), col));
                    }
                    if (j < 8) sb.append(" ");
                }
                sb.append("\n");
            }
            return sb.toString();
        }

        public void parseString(String s) throws Exception {
            if (s.length() != 81) throw new Exception("Строка должна содержать 81 символ");
            for (int i = 0; i < 9; i++) {
                for (int j = 0; j < 9; j++) {
                    char ch = s.charAt(i * 9 + j);
                    if (ch >= '1' && ch <= '9') board[i][j] = ch - '0';
                    else if (ch == '0') board[i][j] = 0;
                    else throw new Exception("Недопустимый символ: " + ch);
                }
            }
        }

        public String exportString() {
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < 9; i++)
                for (int j = 0; j < 9; j++)
                    sb.append(board[i][j] == 0 ? '0' : (char)(board[i][j] + '0'));
            return sb.toString();
        }

        public boolean isValid() {
            // Строки
            for (int i = 0; i < 9; i++) {
                boolean[] seen = new boolean[10];
                for (int j = 0; j < 9; j++) {
                    int val = board[i][j];
                    if (val != 0) {
                        if (val < 1 || val > 9 || seen[val]) return false;
                        seen[val] = true;
                    }
                }
            }
            // Столбцы
            for (int j = 0; j < 9; j++) {
                boolean[] seen = new boolean[10];
                for (int i = 0; i < 9; i++) {
                    int val = board[i][j];
                    if (val != 0) {
                        if (seen[val]) return false;
                        seen[val] = true;
                    }
                }
            }
            // Блоки
            for (int br = 0; br < 3; br++) {
                for (int bc = 0; bc < 3; bc++) {
                    boolean[] seen = new boolean[10];
                    for (int i = 0; i < 3; i++) {
                        for (int j = 0; j < 3; j++) {
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

        public List<Integer> getCandidates(int row, int col) {
            if (board[row][col] != 0) return new ArrayList<>();
            boolean[] used = new boolean[10];
            for (int j = 0; j < 9; j++) if (board[row][j] != 0) used[board[row][j]] = true;
            for (int i = 0; i < 9; i++) if (board[i][col] != 0) used[board[i][col]] = true;
            int br = row / 3 * 3, bc = col / 3 * 3;
            for (int i = 0; i < 3; i++)
                for (int j = 0; j < 3; j++)
                    if (board[br+i][bc+j] != 0) used[board[br+i][bc+j]] = true;
            List<Integer> cands = new ArrayList<>();
            for (int v = 1; v <= 9; v++) if (!used[v]) cands.add(v);
            return cands;
        }

        public boolean findBestEmpty(int[] pos) {
            int bestR = -1, bestC = -1;
            int minCands = 10;
            List<Integer> bestCands = null;
            for (int i = 0; i < 9; i++) {
                for (int j = 0; j < 9; j++) {
                    if (board[i][j] == 0) {
                        List<Integer> cands = getCandidates(i, j);
                        if (cands.size() < minCands) {
                            minCands = cands.size();
                            bestR = i; bestC = j; bestCands = cands;
                            if (minCands == 1) {
                                pos[0] = bestR; pos[1] = bestC;
                                return true;
                            }
                        }
                    }
                }
            }
            if (bestR != -1) {
                pos[0] = bestR; pos[1] = bestC;
                return true;
            }
            return false;
        }

        private boolean solveInternal() {
            int[] pos = new int[2];
            if (!findBestEmpty(pos)) return true;
            int r = pos[0], c = pos[1];
            for (int val : getCandidates(r, c)) {
                board[r][c] = val;
                if (solveInternal()) return true;
                board[r][c] = 0;
            }
            return false;
        }

        public boolean solve(boolean animate, double delay) throws InterruptedException {
            steps = 0;
            long start = System.currentTimeMillis();
            class Solver {
                boolean solveFn() {
                    steps++;
                    int[] pos = new int[2];
                    if (!findBestEmpty(pos)) return true;
                    int r = pos[0], c = pos[1];
                    for (int val : getCandidates(r, c)) {
                        board[r][c] = val;
                        if (animate) {
                            try { Thread.sleep((long)(delay * 1000)); } catch (InterruptedException e) {}
                            System.out.print("\033[H\033[2J");
                            System.out.println(render(Arrays.asList(new int[]{r, c})));
                            System.out.println(colorize("Шаг: " + steps, YELLOW));
                        }
                        if (solveFn()) return true;
                        board[r][c] = 0;
                    }
                    return false;
                }
            }
            boolean solved = new Solver().solveFn();
            if (animate) System.out.print("\033[H\033[2J");
            return solved;
        }

        public int countSolutions(int maxCount) {
            int[] count = new int[]{0};
            class Counter {
                void countFn() {
                    if (count[0] >= maxCount) return;
                    int[] pos = new int[2];
                    if (!findBestEmpty(pos)) { count[0]++; return; }
                    int r = pos[0], c = pos[1];
                    for (int val : getCandidates(r, c)) {
                        board[r][c] = val;
                        countFn();
                        board[r][c] = 0;
                        if (count[0] >= maxCount) return;
                    }
                }
            }
            new Counter().countFn();
            return count[0];
        }

        public static SudokuBoard generate(String level) {
            int minEmpty, maxEmpty;
            if (level.equals("easy")) { minEmpty = 30; maxEmpty = 35; }
            else if (level.equals("medium")) { minEmpty = 40; maxEmpty = 45; }
            else if (level.equals("hard")) { minEmpty = 50; maxEmpty = 55; }
            else { minEmpty = 30; maxEmpty = 35; }

            SudokuBoard s = new SudokuBoard();
            s.solveInternal();
            int[][] filled = new int[9][9];
            for (int i = 0; i < 9; i++) System.arraycopy(s.board[i], 0, filled[i], 0, 9);

            Random rand = new Random();
            int target = minEmpty + rand.nextInt(maxEmpty - minEmpty + 1);
            int removed = 0, attempts = 0;
            while (removed < target && attempts < 10000) {
                attempts++;
                int i = rand.nextInt(9), j = rand.nextInt(9);
                if (filled[i][j] != 0) {
                    int backup = filled[i][j];
                    filled[i][j] = 0;
                    SudokuBoard test = new SudokuBoard(filled);
                    if (test.countSolutions(2) == 1) removed++;
                    else filled[i][j] = backup;
                }
            }
            return new SudokuBoard(filled);
        }
    }

    private static SudokuBoard loadBoard(String source) throws Exception {
        String content;
        if (Files.exists(Paths.get(source))) {
            content = new String(Files.readAllBytes(Paths.get(source)));
        } else {
            content = source;
        }
        StringBuilder digits = new StringBuilder();
        for (char ch : content.toCharArray()) {
            if (ch >= '0' && ch <= '9') digits.append(ch);
        }
        if (digits.length() != 81) throw new Exception("Должно быть 81 цифра");
        SudokuBoard s = new SudokuBoard();
        s.parseString(digits.toString());
        return s;
    }

    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.out.println(colorize("Usage: java sudoku <generate|solve|check|count|export> [options]", YELLOW));
            System.out.println("  generate -l <easy|medium|hard>");
            System.out.println("  solve -i <file|string> [-a] [--delay <sec>]");
            System.out.println("  check -i <file|string>");
            System.out.println("  count -i <file|string> [--max <N>]");
            System.out.println("  export -i <file|string>");
            return;
        }

        String cmd = args[0];
        String input = "", output = "", level = "easy";
        boolean animate = false;
        double delay = 0.1;
        int maxCount = 100;

        for (int i = 1; i < args.length; i++) {
            if (args[i].equals("-i") && i+1 < args.length) input = args[++i];
            else if (args[i].equals("-o") && i+1 < args.length) output = args[++i];
            else if (args[i].equals("-l") && i+1 < args.length) level = args[++i];
            else if (args[i].equals("-a")) animate = true;
            else if (args[i].equals("--delay") && i+1 < args.length) delay = Double.parseDouble(args[++i]);
            else if (args[i].equals("--max") && i+1 < args.length) maxCount = Integer.parseInt(args[++i]);
        }

        try {
            switch (cmd) {
                case "generate": {
                    SudokuBoard s = SudokuBoard.generate(level);
                    if (!output.isEmpty()) Files.write(Paths.get(output), s.exportString().getBytes());
                    System.out.println(s.render(null));
                    int empty = 0;
                    for (int i = 0; i < 9; i++)
                        for (int j = 0; j < 9; j++)
                            if (s.board[i][j] == 0) empty++;
                    System.out.println(colorize("Уровень: " + level + ", пустых ячеек: " + empty, YELLOW));
                    break;
                }
                case "solve": {
                    if (input.isEmpty()) throw new Exception("Укажите входную доску через -i");
                    SudokuBoard s = loadBoard(input);
                    if (!s.isValid()) throw new Exception("Доска невалидна!");
                    boolean solved = s.solve(animate, delay);
                    if (solved) {
                        if (!output.isEmpty()) Files.write(Paths.get(output), s.exportString().getBytes());
                        System.out.println(s.render(null));
                        System.out.println(colorize("Решение найдено! Шагов: " + s.steps, GREEN));
                    } else {
                        System.out.println(colorize("Решения не существует!", RED));
                    }
                    break;
                }
                case "check": {
                    if (input.isEmpty()) throw new Exception("Укажите входную доску через -i");
                    SudokuBoard s = loadBoard(input);
                    System.out.println(s.isValid() ? colorize("✅ Доска корректна.", GREEN) : colorize("❌ Доска невалидна.", RED));
                    break;
                }
                case "count": {
                    if (input.isEmpty()) throw new Exception("Укажите входную доску через -i");
                    SudokuBoard s = loadBoard(input);
                    int cnt = s.countSolutions(maxCount);
                    if (cnt >= maxCount) {
                        System.out.println(colorize("Количество решений >= " + maxCount + " (ограничено)", YELLOW));
                    } else {
                        System.out.println(colorize("Количество решений: " + cnt, GREEN));
                    }
                    break;
                }
                case "export": {
                    if (input.isEmpty()) throw new Exception("Укажите входную доску через -i");
                    SudokuBoard s = loadBoard(input);
                    System.out.println(s.exportString());
                    break;
                }
                default:
                    System.out.println(colorize("Неизвестная команда: " + cmd, RED));
            }
        } catch (Exception e) {
            System.err.println(colorize("Ошибка: " + e.getMessage(), RED));
            System.exit(1);
        }
    }
}
