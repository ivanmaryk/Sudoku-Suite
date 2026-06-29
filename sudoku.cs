// sudoku.cs
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading;

class Sudoku
{
    static string Colorize(string text, string color)
    {
        string col = color switch
        {
            "red" => "\x1b[91m",
            "green" => "\x1b[92m",
            "yellow" => "\x1b[93m",
            "blue" => "\x1b[94m",
            "cyan" => "\x1b[96m",
            "gray" => "\x1b[90m",
            "bold" => "\x1b[1m",
            _ => "\x1b[0m"
        };
        return col + text + "\x1b[0m";
    }

    private int[,] board = new int[9, 9];
    public int Steps { get; private set; }

    public Sudoku() { }

    public Sudoku(int[,] src)
    {
        for (int i = 0; i < 9; i++)
            for (int j = 0; j < 9; j++)
                board[i, j] = src[i, j];
    }

    public string Render(List<(int, int)> highlight = null)
    {
        var sb = new StringBuilder();
        for (int i = 0; i < 9; i++)
        {
            if (i % 3 == 0 && i > 0)
                sb.AppendLine(Colorize("┃━━━┃━━━┃━━━┃", "gray"));
            for (int j = 0; j < 9; j++)
            {
                if (j % 3 == 0 && j > 0)
                    sb.Append(Colorize("┃", "gray"));
                int val = board[i, j];
                if (val == 0)
                    sb.Append(" ");
                else
                {
                    string col = "green";
                    if (highlight != null && highlight.Contains((i, j)))
                        col = "cyan";
                    sb.Append(Colorize(val.ToString(), col));
                }
                if (j < 8) sb.Append(" ");
            }
            sb.AppendLine();
        }
        return sb.ToString();
    }

    public void ParseString(string s)
    {
        if (s.Length != 81) throw new Exception("Строка должна содержать 81 символ");
        for (int i = 0; i < 9; i++)
            for (int j = 0; j < 9; j++)
            {
                char ch = s[i * 9 + j];
                if (ch >= '1' && ch <= '9') board[i, j] = ch - '0';
                else if (ch == '0') board[i, j] = 0;
                else throw new Exception($"Недопустимый символ: {ch}");
            }
    }

    public string ExportString()
    {
        var sb = new StringBuilder();
        for (int i = 0; i < 9; i++)
            for (int j = 0; j < 9; j++)
                sb.Append(board[i, j] == 0 ? '0' : (char)(board[i, j] + '0'));
        return sb.ToString();
    }

    public bool IsValid()
    {
        // Строки
        for (int i = 0; i < 9; i++)
        {
            var seen = new HashSet<int>();
            for (int j = 0; j < 9; j++)
            {
                int val = board[i, j];
                if (val != 0)
                {
                    if (val < 1 || val > 9 || seen.Contains(val)) return false;
                    seen.Add(val);
                }
            }
        }
        // Столбцы
        for (int j = 0; j < 9; j++)
        {
            var seen = new HashSet<int>();
            for (int i = 0; i < 9; i++)
            {
                int val = board[i, j];
                if (val != 0)
                {
                    if (seen.Contains(val)) return false;
                    seen.Add(val);
                }
            }
        }
        // Блоки
        for (int br = 0; br < 3; br++)
            for (int bc = 0; bc < 3; bc++)
            {
                var seen = new HashSet<int>();
                for (int i = 0; i < 3; i++)
                    for (int j = 0; j < 3; j++)
                    {
                        int val = board[br * 3 + i, bc * 3 + j];
                        if (val != 0)
                        {
                            if (seen.Contains(val)) return false;
                            seen.Add(val);
                        }
                    }
            }
        return true;
    }

    public List<int> GetCandidates(int row, int col)
    {
        if (board[row, col] != 0) return new List<int>();
        var used = new bool[10];
        for (int j = 0; j < 9; j++) if (board[row, j] != 0) used[board[row, j]] = true;
        for (int i = 0; i < 9; i++) if (board[i, col] != 0) used[board[i, col]] = true;
        int br = row / 3 * 3, bc = col / 3 * 3;
        for (int i = 0; i < 3; i++)
            for (int j = 0; j < 3; j++)
                if (board[br + i, bc + j] != 0) used[board[br + i, bc + j]] = true;
        var cands = new List<int>();
        for (int v = 1; v <= 9; v++) if (!used[v]) cands.Add(v);
        return cands;
    }

    public bool FindBestEmpty(out int r, out int c, out List<int> cands)
    {
        r = -1; c = -1; cands = null;
        int minCands = 10;
        for (int i = 0; i < 9; i++)
            for (int j = 0; j < 9; j++)
                if (board[i, j] == 0)
                {
                    var cnd = GetCandidates(i, j);
                    if (cnd.Count < minCands)
                    {
                        minCands = cnd.Count;
                        r = i; c = j; cands = cnd;
                        if (minCands == 1) return true;
                    }
                }
        return r != -1;
    }

    private bool SolveInternal()
    {
        if (!FindBestEmpty(out int r, out int c, out var cands)) return true;
        foreach (int val in cands)
        {
            board[r, c] = val;
            if (SolveInternal()) return true;
            board[r, c] = 0;
        }
        return false;
    }

    public bool Solve(bool animate = false, double delay = 0.1)
    {
        Steps = 0;
        var start = DateTime.Now;
        bool SolveFn()
        {
            Steps++;
            if (!FindBestEmpty(out int r, out int c, out var cands)) return true;
            foreach (int val in cands)
            {
                board[r, c] = val;
                if (animate)
                {
                    Console.Clear();
                    Console.WriteLine(Render(new List<(int, int)> { (r, c) }));
                    Console.WriteLine(Colorize($"Шаг: {Steps}", "yellow"));
                    Thread.Sleep((int)(delay * 1000));
                }
                if (SolveFn()) return true;
                board[r, c] = 0;
            }
            return false;
        }
        bool solved = SolveFn();
        if (animate) Console.Clear();
        return solved;
    }

    public int CountSolutions(int maxCount = 100)
    {
        int count = 0;
        void CountFn()
        {
            if (count >= maxCount) return;
            if (!FindBestEmpty(out int r, out int c, out var cands)) { count++; return; }
            foreach (int val in cands)
            {
                board[r, c] = val;
                CountFn();
                board[r, c] = 0;
                if (count >= maxCount) return;
            }
        }
        CountFn();
        return count;
    }

    public static Sudoku Generate(string level = "easy")
    {
        var levels = new Dictionary<string, (int, int)> {
            { "easy", (30, 35) }, { "medium", (40, 45) }, { "hard", (50, 55) }
        };
        var (minEmpty, maxEmpty) = levels.ContainsKey(level) ? levels[level] : (30, 35);

        var s = new Sudoku();
        s.SolveInternal();
        var filled = new int[9, 9];
        for (int i = 0; i < 9; i++)
            for (int j = 0; j < 9; j++)
                filled[i, j] = s.board[i, j];

        var rand = new Random();
        int target = rand.Next(minEmpty, maxEmpty + 1);
        int removed = 0, attempts = 0;
        while (removed < target && attempts < 10000)
        {
            attempts++;
            int i = rand.Next(9), j = rand.Next(9);
            if (filled[i, j] != 0)
            {
                int backup = filled[i, j];
                filled[i, j] = 0;
                var test = new Sudoku(filled);
                if (test.CountSolutions(2) == 1) removed++;
                else filled[i, j] = backup;
            }
        }
        return new Sudoku(filled);
    }
}

class Program
{
    static Sudoku LoadBoard(string source)
    {
        string content;
        if (File.Exists(source))
            content = File.ReadAllText(source);
        else
            content = source;
        string digits = new string(content.Where(ch => ch >= '0' && ch <= '9').ToArray());
        if (digits.Length != 81) throw new Exception("Должно быть ровно 81 цифра");
        var s = new Sudoku();
        s.ParseString(digits);
        return s;
    }

    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine(Colorize("Usage: sudoku <generate|solve|check|count|export> [options]", "yellow"));
            Console.WriteLine("  generate -l <easy|medium|hard>");
            Console.WriteLine("  solve -i <file|string> [-a] [--delay <sec>]");
            Console.WriteLine("  check -i <file|string>");
            Console.WriteLine("  count -i <file|string> [--max <N>]");
            Console.WriteLine("  export -i <file|string>");
            return;
        }

        string cmd = args[0];
        string input = "", output = "", level = "easy";
        bool animate = false;
        double delay = 0.1;
        int maxCount = 100;

        for (int i = 1; i < args.Length; i++)
        {
            if (args[i] == "-i" && i + 1 < args.Length) input = args[++i];
            else if (args[i] == "-o" && i + 1 < args.Length) output = args[++i];
            else if (args[i] == "-l" && i + 1 < args.Length) level = args[++i];
            else if (args[i] == "-a") animate = true;
            else if (args[i] == "--delay" && i + 1 < args.Length) delay = double.Parse(args[++i]);
            else if (args[i] == "--max" && i + 1 < args.Length) maxCount = int.Parse(args[++i]);
        }

        try
        {
            switch (cmd)
            {
                case "generate":
                    var sg = Sudoku.Generate(level);
                    if (!string.IsNullOrEmpty(output)) File.WriteAllText(output, sg.ExportString());
                    Console.WriteLine(sg.Render());
                    int empty = 0;
                    for (int i = 0; i < 9; i++)
                        for (int j = 0; j < 9; j++)
                            if (sg.board[i, j] == 0) empty++;
                    Console.WriteLine(Colorize($"Уровень: {level}, пустых ячеек: {empty}", "yellow"));
                    break;

                case "solve":
                    if (string.IsNullOrEmpty(input)) throw new Exception("Укажите входную доску через -i");
                    var ss = LoadBoard(input);
                    if (!ss.IsValid()) throw new Exception("Доска невалидна!");
                    bool solved = ss.Solve(animate, delay);
                    if (solved)
                    {
                        if (!string.IsNullOrEmpty(output)) File.WriteAllText(output, ss.ExportString());
                        Console.WriteLine(ss.Render());
                        Console.WriteLine(Colorize($"Решение найдено! Шагов: {ss.Steps}", "green"));
                    }
                    else
                        Console.WriteLine(Colorize("Решения не существует!", "red"));
                    break;

                case "check":
                    if (string.IsNullOrEmpty(input)) throw new Exception("Укажите входную доску через -i");
                    var sc = LoadBoard(input);
                    Console.WriteLine(sc.IsValid() ? Colorize("✅ Доска корректна.", "green") : Colorize("❌ Доска невалидна.", "red"));
                    break;

                case "count":
                    if (string.IsNullOrEmpty(input)) throw new Exception("Укажите входную доску через -i");
                    var scnt = LoadBoard(input);
                    int cnt = scnt.CountSolutions(maxCount);
                    if (cnt >= maxCount)
                        Console.WriteLine(Colorize($"Количество решений >= {maxCount} (ограничено)", "yellow"));
                    else
                        Console.WriteLine(Colorize($"Количество решений: {cnt}", "green"));
                    break;

                case "export":
                    if (string.IsNullOrEmpty(input)) throw new Exception("Укажите входную доску через -i");
                    var se = LoadBoard(input);
                    Console.WriteLine(se.ExportString());
                    break;

                default:
                    Console.WriteLine(Colorize($"Неизвестная команда: {cmd}", "red"));
                    break;
            }
        }
        catch (Exception e)
        {
            Console.WriteLine(Colorize($"Ошибка: {e.Message}", "red"));
        }
    }
}
