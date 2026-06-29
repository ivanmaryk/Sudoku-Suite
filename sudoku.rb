#!/usr/bin/env ruby
# sudoku.rb
# encoding: UTF-8

require 'json'
require 'fileutils'
require 'timeout'

COLORS = {
  reset: "\e[0m",
  red: "\e[91m",
  green: "\e[92m",
  yellow: "\e[93m",
  blue: "\e[94m",
  cyan: "\e[96m",
  gray: "\e[90m",
  bold: "\e[1m"
}

def colorize(text, color)
  "#{COLORS[color]}#{text}#{COLORS[:reset]}"
end

class Sudoku
  attr_reader :board, :steps

  def initialize(board = nil)
    @board = board ? board.map(&:dup) : Array.new(9) { Array.new(9, 0) }
    @steps = 0
  end

  def render(highlight = nil)
    lines = []
    9.times do |i|
      lines << colorize('┃━━━┃━━━┃━━━┃', :gray) if i % 3 == 0 && i > 0
      row = []
      9.times do |j|
        row << colorize('┃', :gray) if j % 3 == 0 && j > 0
        val = @board[i][j]
        if val == 0
          row << ' '
        else
          col = :green
          if highlight && highlight.include?([i, j])
            col = :cyan
          end
          row << colorize(val.to_s, col)
        end
      end
      lines << row.join(' ')
    end
    lines.join("\n")
  end

  def parse_string(str)
    raise "Строка должна содержать 81 символ" if str.size != 81
    9.times do |i|
      9.times do |j|
        ch = str[i*9 + j]
        if ch.between?('1', '9')
          @board[i][j] = ch.to_i
        elsif ch == '0'
          @board[i][j] = 0
        else
          raise "Недопустимый символ: #{ch}"
        end
      end
    end
  end

  def export_string
    @board.flatten.map { |v| v == 0 ? '0' : v.to_s }.join
  end

  def valid?
    # Строки
    9.times do |i|
      seen = {}
      9.times do |j|
        val = @board[i][j]
        next if val == 0
        return false if val < 1 || val > 9 || seen[val]
        seen[val] = true
      end
    end
    # Столбцы
    9.times do |j|
      seen = {}
      9.times do |i|
        val = @board[i][j]
        next if val == 0
        return false if seen[val]
        seen[val] = true
      end
    end
    # Блоки
    3.times do |br|
      3.times do |bc|
        seen = {}
        3.times do |i|
          3.times do |j|
            val = @board[br*3 + i][bc*3 + j]
            next if val == 0
            return false if seen[val]
            seen[val] = true
          end
        end
      end
    end
    true
  end

  def candidates(row, col)
    return [] if @board[row][col] != 0
    used = {}
    9.times { |j| used[@board[row][j]] = true if @board[row][j] != 0 }
    9.times { |i| used[@board[i][col]] = true if @board[i][col] != 0 }
    br, bc = row / 3 * 3, col / 3 * 3
    3.times do |i|
      3.times do |j|
        val = @board[br+i][bc+j]
        used[val] = true if val != 0
      end
    end
    (1..9).reject { |v| used[v] }
  end

  def find_best_empty
    best_r = best_c = -1
    best_cands = nil
    min_cands = 10
    9.times do |i|
      9.times do |j|
        if @board[i][j] == 0
          cands = candidates(i, j)
          if cands.size < min_cands
            min_cands = cands.size
            best_r, best_c = i, j
            best_cands = cands
            return [best_r, best_c, best_cands] if min_cands == 1
          end
        end
      end
    end
    [best_r, best_c, best_cands]
  end

  def solve_internal
    r, c, cands = find_best_empty
    return true if r == -1
    cands.each do |val|
      @board[r][c] = val
      return true if solve_internal
      @board[r][c] = 0
    end
    false
  end

  def solve(animate = false, delay = 0.1)
    @steps = 0
    start = Time.now
    solve_fn = -> {
      @steps += 1
      r, c, cands = find_best_empty
      return true if r == -1
      cands.each do |val|
        @board[r][c] = val
        if animate
          system('clear') || system('cls')
          puts render([[r, c]])
          puts colorize("Шаг: #{@steps}", :yellow)
          sleep(delay)
        end
        return true if solve_fn.call
        @board[r][c] = 0
      end
      false
    }
    solved = solve_fn.call
    elapsed = Time.now - start
    [solved, elapsed]
  end

  def count_solutions(max_count = 100)
    count = 0
    count_fn = -> {
      return if count >= max_count
      r, c, cands = find_best_empty
      if r == -1
        count += 1
        return
      end
      cands.each do |val|
        @board[r][c] = val
        count_fn.call
        @board[r][c] = 0
        return if count >= max_count
      end
    }
    count_fn.call
    count
  end

  def self.generate(level = 'easy')
    levels = { 'easy' => [30, 35], 'medium' => [40, 45], 'hard' => [50, 55] }
    min_empty, max_empty = levels[level] || [30, 35]

    s = Sudoku.new
    s.solve_internal
    filled = s.board.map(&:dup)
    target = rand(min_empty..max_empty)
    removed = 0
    attempts = 0
    while removed < target && attempts < 10000
      attempts += 1
      i, j = rand(9), rand(9)
      if filled[i][j] != 0
        backup = filled[i][j]
        filled[i][j] = 0
        test = Sudoku.new(filled)
        if test.count_solutions(2) == 1
          removed += 1
        else
          filled[i][j] = backup
        end
      end
    end
    Sudoku.new(filled)
  end
end

def load_board(source)
  content = File.exist?(source) ? File.read(source) : source
  digits = content.chars.select { |ch| ch.between?('0', '9') }
  raise "Должно быть 81 цифра" if digits.size != 81
  s = Sudoku.new
  s.parse_string(digits.join)
  s
end

def main
  if ARGV.empty?
    puts colorize("Usage: ruby sudoku.rb <generate|solve|check|count|export> [options]", :yellow)
    puts "  generate -l <easy|medium|hard>"
    puts "  solve -i <file|string> [-a] [--delay <sec>]"
    puts "  check -i <file|string>"
    puts "  count -i <file|string> [--max <N>]"
    puts "  export -i <file|string>"
    exit 1
  end

  cmd = ARGV[0]
  input = nil
  output = nil
  level = 'easy'
  animate = false
  delay = 0.1
  max_count = 100

  i = 1
  while i < ARGV.size
    case ARGV[i]
    when '-i' then input = ARGV[i+1]; i += 1
    when '-o' then output = ARGV[i+1]; i += 1
    when '-l' then level = ARGV[i+1]; i += 1
    when '-a' then animate = true
    when '--delay' then delay = ARGV[i+1].to_f; i += 1
    when '--max' then max_count = ARGV[i+1].to_i; i += 1
    end
    i += 1
  end

  begin
    case cmd
    when 'generate'
      s = Sudoku.generate(level)
      File.write(output, s.export_string) if output
      puts s.render
      empty = s.board.flatten.count(0)
      puts colorize("Уровень: #{level}, пустых ячеек: #{empty}", :yellow)

    when 'solve'
      raise "Укажите входную доску через -i" unless input
      s = load_board(input)
      raise "Доска невалидна!" unless s.valid?
      solved, elapsed = s.solve(animate, delay)
      if solved
        File.write(output, s.export_string) if output
        puts s.render
        puts colorize("Решение найдено за #{elapsed.round(3)} сек, шагов: #{s.steps}", :green)
      else
        puts colorize("Решения не существует!", :red)
      end

    when 'check'
      raise "Укажите входную доску через -i" unless input
      s = load_board(input)
      puts s.valid? ? colorize("✅ Доска корректна.", :green) : colorize("❌ Доска невалидна.", :red)

    when 'count'
      raise "Укажите входную доску через -i" unless input
      s = load_board(input)
      cnt = s.count_solutions(max_count)
      if cnt >= max_count
        puts colorize("Количество решений >= #{max_count} (ограничено)", :yellow)
      else
        puts colorize("Количество решений: #{cnt}", :green)
      end

    when 'export'
      raise "Укажите входную доску через -i" unless input
      s = load_board(input)
      puts s.export_string

    else
      puts colorize("Неизвестная команда: #{cmd}", :red)
    end
  rescue => e
    puts colorize("Ошибка: #{e.message}", :red)
    exit 1
  end
end

main if __FILE__ == $0
