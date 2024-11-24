package cat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"

	"github.com/zooyer/gobox/box"
)

const version = `cat: (by goland) 1.0
Copyright (C) 2024 Free Software Foundation, Inc.
`

const usage = `Usage: cat [OPTION]... [FILE]...
Concatenate FILE(s) to standard output.

  -h, --help     display this help and exit
  -v, --version  output version information and exit
  -              read from standard input
`

const (
	kb = 1024
	mb = kb * 1024
)

const defaultBufferSize = 1 * mb // 默认缓冲区大小

var sizeMaps = map[int64]int{
	1 * mb:   4 * kb,   // 文件 <= 1MB 时使用 4KB 缓冲区
	10 * mb:  32 * kb,  // 文件 <= 10MB 时使用 32KB 缓冲区
	50 * mb:  128 * kb, // 文件 <= 50MB 时使用 128KB 缓冲区
	100 * mb: 256 * kb, // 文件 <= 100MB 时使用 256KB 缓冲区
	200 * mb: 512 * kb, // 文件 <= 200MB 时使用 512KB 缓冲区
}

var bufferSizes = make([]int64, 0, len(sizeMaps)) // 缓冲区大小的有序切片

func init() {
	for size := range sizeMaps {
		bufferSizes = append(bufferSizes, size)
	}
	slices.Sort(bufferSizes) // 排序保证二分查找可用
}

// 根据文件大小返回合适的缓冲区大小
func getBufferSize(fileSize int64) int {
	if fileSize == 0 || len(bufferSizes) == 0 {
		return defaultBufferSize
	}

	// 二分查找
	var index = sort.Search(len(bufferSizes), func(i int) bool {
		return fileSize <= bufferSizes[i]
	})

	if index < len(bufferSizes) {
		// 权重比例 2:1 判断是否使用更大的缓冲区
		if index+1 < len(bufferSizes) && fileSize > (bufferSizes[index]*2/3) {
			index++
		}
		return sizeMaps[bufferSizes[index]]
	}

	return defaultBufferSize // 如果超过最大值，返回默认缓冲区大小
}

// 使用指定大小的缓冲区拷贝数据
func copyWithBuffer(dst io.Writer, src io.Reader, size int) (err error) {
	if _, err = io.CopyBuffer(dst, src, make([]byte, size)); err != nil {
		return
	}

	return
}

// 读取文件内容并输出到指定的 Writer
func readFile(ctx context.Context, filename string, out io.Writer) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}()

	info, err := file.Stat()
	if err != nil {
		return
	}

	var size = getBufferSize(info.Size())

	return copyWithBuffer(out, withContext(ctx, file), size)
}

// 读取文件到out，如果错误则写入err并返回错误码
func readFileErrno(ctx context.Context, filename string, out, err io.Writer) (errno int) {
	if e := readFile(ctx, filename, out); e != nil {
		errno = 3
		_, _ = fmt.Fprintln(err, fmt.Sprintf("cat: %s: %s", filename, e.Error()))
	}

	return
}

func Cat(ctx context.Context, opt box.Option) (errno int) {
	if len(opt.Args) == 1 {
		opt.Args = []string{"-"}
	}

	var (
		eno int
		err error
		end bool
	)

	for _, arg := range opt.Args[1:] {
		if end {
			if eno = readFileErrno(ctx, arg, opt.Stdout, opt.Stderr); eno != 0 {
				errno = eno
			}
			continue
		}

		switch arg {
		case "--":
			end = true
		case "-h", "--help":
			_, _ = fmt.Fprint(opt.Stdout, usage)
		case "-v", "--version":
			_, _ = fmt.Fprint(opt.Stdout, version)
		case "-", "-i", "--stdin":
			if _, err = io.Copy(opt.Stdout, withContext(ctx, opt.Stdin)); err != nil {
				errno = 2
			}
		default:
			if arg[0] == '-' {
				_, _ = fmt.Fprintln(opt.Stderr, fmt.Sprintf("cat: invalid option '%s'", arg))
				_, _ = fmt.Fprintln(opt.Stderr, "Try 'cat --help' for more information.")
				errno = 1
			} else {
				if eno = readFileErrno(ctx, arg, opt.Stdout, opt.Stderr); eno != 0 {
					errno = eno
				}
			}
		}
	}

	return
}