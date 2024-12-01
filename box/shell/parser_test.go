package shell

import (
	"reflect"
	"testing"
)

func TestNewParser1(t *testing.T) {
	var tests = []struct {
		input    string
		expected []Command
		err      bool
	}{
		{
			input: `echo "hello world"`,
			expected: []Command{
				{
					Path: "echo",
					Args: []string{"hello world"},
				},
			},
			err: false,
		},
		{
			input: "cat < input.txt",
			expected: []Command{
				{
					Path:  "cat",
					Input: "input.txt",
				},
			},
			err: false,
		},
		{
			input: "ls -l > output.txt",
			expected: []Command{
				{
					Path:   "ls",
					Args:   []string{"-l"},
					Output: "output.txt",
				},
			},
			err: false,
		},
		{
			input: `echo "log entry" >> logs.txt`,
			expected: []Command{
				{
					Path:   "echo",
					Args:   []string{"log entry"},
					Append: "logs.txt",
				},
			},
			err: false,
		},
		{
			input: "cat << EOF\nThis is a test\nEOF\n",
			expected: []Command{
				{
					Path:    "cat",
					Heredoc: "This is a test",
				},
			},
			err: false,
		},
		{
			input: `ls | grep "test"`,
			expected: []Command{
				{
					Path: "ls",
					Pipe: &Command{
						Path: "grep",
						Args: []string{"test"},
					},
				},
			},
			err: false,
		},
		{
			input: "mkdir new_dir && cd new_dir",
			expected: []Command{
				{
					Path: "mkdir",
					Args: []string{"new_dir"},
					And: &Command{
						Path: "cd",
						Args: []string{"new_dir"},
					},
				},
			},
			err: false,
		},
		{
			input: `false || echo "Command failed"`,
			expected: []Command{
				{
					Path: "false",
					Or: &Command{
						Path: "echo",
						Args: []string{"Command failed"},
					},
				},
			},
			err: false,
		},
		{
			input: "sleep 10 &\n",
			expected: []Command{
				{
					Path:       "sleep",
					Args:       []string{"10"},
					Background: true,
				},
			},
			err: false,
		},
		{
			input: `echo "hello"; ls; pwd`,
			expected: []Command{
				{
					Path: "echo",
					Args: []string{"hello"},
				},
				{
					Path: "ls",
				},
				{
					Path: "pwd",
				},
			},
			err: false,
		},
		{
			input: `cat file.txt | grep "error" && echo "Error found" || echo "No error found"`,
			expected: []Command{
				{
					Path: "cat",
					Args: []string{"file.txt"},
					Pipe: &Command{
						Path: "grep",
						Args: []string{"error"},
						And: &Command{
							Path: "echo",
							Args: []string{"Error found"},
						},
						Or: &Command{
							Path: "echo",
							Args: []string{"No error found"},
						},
					},
				},
			},
			err: false,
		},
		{
			input: `ls | grep "test" && echo "Found" || echo "Not Found"`,
			expected: []Command{
				{
					Path: "ls",
					Pipe: &Command{
						Path: "grep",
						Args: []string{"test"},
						And: &Command{
							Path: "echo",
							Args: []string{"Found"},
						},
						Or: &Command{
							Path: "echo",
							Args: []string{"Not Found"},
						},
					},
				},
			},
			err: false,
		},

		{
			input: `ls | grep "test" > results.txt`,
			expected: []Command{
				{
					Path: "ls",
					Pipe: &Command{
						Path:   "grep",
						Args:   []string{"test"},
						Output: "results.txt",
					},
				},
			},
			err: false,
		},
		{
			input: `mkdir test_dir && cd test_dir || echo "Failed to create directory"`,
			expected: []Command{
				{
					Path: "mkdir",
					Args: []string{"test_dir"},
					And: &Command{
						Path: "cd",
						Args: []string{"test_dir"},
					},
					Or: &Command{
						Path: "echo",
						Args: []string{"Failed to create directory"},
					},
				},
			},
			err: false,
		},
		{
			input: `cat file.txt | grep "error" || (echo "Error found" && exit 1)`,
			expected: []Command{
				{
					Path: "cat",
					Args: []string{"file.txt"},
					Pipe: &Command{
						Path: "grep",
						Args: []string{"error"},
						Or: &Command{
							Path: "(echo",
							Args: []string{"Error found"},
							And: &Command{
								Path: "exit",
								Args: []string{"1)"},
							},
						},
					},
				},
			},
			err: false,
		},
		{
			input: "cat << EOF > output.txt\nLine 1\nLine 2\nEOF\n",
			expected: []Command{
				{
					Path:    "cat",
					Output:  "output.txt",
					Heredoc: "Line 1\nLine 2",
				},
			},
			err: false,
		},
		{
			input: `(sleep 5 && echo "Background task finished") & echo "Started"`,
			expected: []Command{
				{
					Path: "(sleep",
					Args: []string{"5"},
					And: &Command{
						Path:       "echo",
						Args:       []string{"Background task finished)"},
						Background: true,
					},
				},
				{
					Path: "echo",
					Args: []string{"Started"},
				},
			},
		},
		{
			input: `echo "Step 1"; echo "Step 2"; echo "Step 3"`,
			expected: []Command{
				{
					Path: "echo",
					Args: []string{"Step 1"},
				},
				{
					Path: "echo",
					Args: []string{"Step 2"},
				},
				{
					Path: "echo",
					Args: []string{"Step 3"},
				},
			},
			err: false,
		},
		{
			input: `cat < input.txt && echo "File read successfully" || echo "File read failed"`,
			expected: []Command{
				{
					Path:  "cat",
					Input: "input.txt",
					And: &Command{
						Path: "echo",
						Args: []string{"File read successfully"},
					},
					Or: &Command{
						Path: "echo",
						Args: []string{"File read failed"},
					},
				},
			},
			err: false,
		},
		{
			input: `echo "Test log" >> logs.txt && echo "Log written" || echo "Failed to write log"`,
			expected: []Command{
				{
					Path:   "echo",
					Args:   []string{"Test log"},
					Append: "logs.txt",
					And: &Command{
						Path: "echo",
						Args: []string{"Log written"},
					},
					Or: &Command{
						Path: "echo",
						Args: []string{"Failed to write log"},
					},
				},
			},
			err: false,
		},
		{
			input: `ls -l | grep "file" > output.txt && cat output.txt || echo "No files found"`,
			expected: []Command{
				{
					Path: "ls",
					Args: []string{"-l"},
					Pipe: &Command{
						Path:   "grep",
						Args:   []string{"file"},
						Output: "output.txt",
						And: &Command{
							Path: "cat",
							Args: []string{"output.txt"},
						},
						Or: &Command{
							Path: "echo",
							Args: []string{"No files found"},
						},
					},
				},
			},
			err: false,
		},
		{
			input: "cat << EOF &\nBackground task\nEOF\necho \"HereDoc submitted\"\n",
			expected: []Command{
				{
					Path:       "cat",
					Heredoc:    "Background task",
					Background: true,
				},
				{
					Path: "echo",
					Args: []string{"HereDoc submitted"},
				},
			},
			err: false,
		},

		{
			input: `cat < input.txt > output.txt &`,
			expected: []Command{
				{
					Path:       "cat",
					Input:      "input.txt",
					Output:     "output.txt",
					Background: true,
				},
			},
			err: false,
		},
		{
			input: `cat << EOF | grep "pattern"`,
			expected: []Command{
				{
					Path: "cat",
					Pipe: &Command{
						Path: "grep",
						Args: []string{"pattern"},
					},
				},
			},
			err: true,
		},
		{
			input: `echo "Hello" > hello.txt; cat hello.txt`,
			expected: []Command{
				{
					Path:   "echo",
					Args:   []string{"Hello"},
					Output: "hello.txt",
				},
				{
					Path: "cat",
					Args: []string{"hello.txt"},
				},
			},
			err: false,
		},
		{
			input: `ls && mkdir test || rmdir test; echo "Done"`,
			expected: []Command{
				{
					Path: "ls",
					And: &Command{
						Path: "mkdir",
						Args: []string{"test"},
					},
					Or: &Command{
						Path: "rmdir",
						Args: []string{"test"},
					},
				},
				{
					Path: "echo",
					Args: []string{"Done"},
				},
			},
			err: false,
		},
		{
			input: `find / -name "file" >> results.log && echo "Search completed"`,
			expected: []Command{
				{
					Path:   "find",
					Args:   []string{"/", "-name", "file"},
					Append: "results.log",
					And: &Command{
						Path: "echo",
						Args: []string{"Search completed"},
					},
				},
			},
			err: false,
		},
		{
			input: `sleep 5 & echo "Done"`,
			expected: []Command{
				{
					Path:       "sleep",
					Args:       []string{"5"},
					Background: true,
				},
				{
					Path: "echo",
					Args: []string{"Done"},
				},
			},
			err: false,
		},

		{
			input: `cat file1 | grep "error" > output.log && echo "Processed" >> logs.txt || echo "Failed" > error.log; ls -l /tmp | wc -l > count.txt && echo "done"`,
			expected: []Command{
				{
					Path: "cat",
					Args: []string{"file1"},
					Pipe: &Command{
						Path:   "grep",
						Args:   []string{"error"},
						Output: "output.log",
						And: &Command{
							Path:   "echo",
							Args:   []string{"Processed"},
							Append: "logs.txt",
						},
						Or: &Command{
							Path:   "echo",
							Args:   []string{"Failed"},
							Output: "error.log",
						},
					},
				},
				{
					Path: "ls",
					Args: []string{"-l", "/tmp"},
					Pipe: &Command{
						Path:   "wc",
						Args:   []string{"-l"},
						Output: "count.txt",
						And: &Command{
							Path: "echo",
							Args: []string{"done"},
						},
					},
				},
			},
			err: false,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			commands, err := ParseCommands(test.input)
			if (err != nil) != test.err {
				t.Errorf("expected error: %v, got: %v", test.err, err)
			}

			if !reflect.DeepEqual(commands, test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, commands)
			}
		})
	}
}
