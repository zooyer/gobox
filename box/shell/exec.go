package shell

type ExecUnit interface {
	Token() (begin, end string)
	Exec()
}

var execUints = []ExecUnit{
	quoteUnit("'"),
	quoteUnit(`"`),
	quoteUnit("`"),

	// ()、中括号 []、花括号 {}、双中括号 [[ ]]
	//for i in {1..5  # 花括号未闭合
	//# Shell 会继续等待你输入 }

	// if、for、while 等
	//if [ -f myfile ]; then
	//  echo "File exists"
	//# Shell 会继续等待输入 "fi" 结束
	//fi

	// 换行符转义（反斜杠 \）
	//echo "This is a \
	//multi-line command"
	//# Shell 会等待输入下一行

	// 管道符号（|）
	//echo "Hello" |
	//grep H
	//# Shell 会等待管道后续的命令

	// 命令替换使用 $() 或 `command`
	//echo $(ls -l
	//# Shell 会继续等待 )

	// << 或 <<- 进行多行输入，直到输入结束标识符（如 EOF）
	//cat << EOF
	//This is a
	//multi-line
	//string.
	//EOF
	//# Shell 会等待输入 "EOF"

	// 在命令行以 & 结尾启动后台任务，如果忘记闭合括号，Shell 会继续等待。
	//(
	//echo "Task in background" &
	//# Shell 会继续等待闭合 )

	// 函数定义未完成
	//myfunc() {
	//  echo "Hello"
	//# Shell 会继续等待输入 }

	// 未闭合的 case 语句
	//case $var in
	//  1) echo "One";;
	//  2) echo "Two";;
	//# Shell 会继续等待输入 esac
}

var (
	beginUints = make(map[string]ExecUnit)
	endUints   = make(map[string]ExecUnit)
)

func init() {
	for _, unit := range execUints {
		begin, end := unit.Token()
		beginUints[begin] = unit
		endUints[end] = unit
	}
}

type quoteUnit string

func (q quoteUnit) Token() (begin, end string) {
	return string(q), string(q)
}

func (q quoteUnit) Exec() {}
