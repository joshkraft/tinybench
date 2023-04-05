# tinybench

tinybench is a tiny (~200 lines), minimal Go tool for benchmarking JavaScript code. tinybench is designed for quickly comparing the performance of small code snippets. This tool is not production-ready - for more comprehensive and accurate benchmarking, please use a more robust benchmarking tool.

## How it works

tinybench uses the Go [exec](https://golang.org/pkg/os/exec/) package to run JavaScript code in a Node.js process. Each benchmark segment is run as many times as possible in a 10 second period, and the minimum, maximum, and median execution times are recorded for comparison.

## Requirements

- Go (1.16 or higher)
- Node.js

## Installation

To install, first clone the repository:

```bash
git clone https://github.com/joshkraft/tinybench.git
cd tinybench
```

Next, build the binary:

```bash
go build
```

Finally, add the binary to your system's path. You can do this by copying the binary to a directory that is already on your path (e.g., /usr/local/bin on macOS or Linux). Here's an example:

```bash
cp tinybench /usr/local/bin/
```

You should now be able to run tinybench from anywhere on your system.

## Usage

Use the `// tinybench start` and `// tinybench stop` delimiters in your JavaScript file to define the code snippets you want to benchmark. You can create as many benchmarks as you want by adding more pairs of delimiters. Here is an example of a file that compares two different for loops:

```javascript
function forLoop(arr) {
  for (let i = 0; i < arr.length; i++) {
    const element = arr[i];
  }
}

function forOfLoop(arr) {
  for (const element of arr) {}
}

const arr = new Array(1000000).fill().map(() => Math.random());

// tinybench start
forLoop(arr);
// tinybench stop

// tinybench start
forOfLoop(arr);
// tinybench stop
```

Run the tinybench command and pass in the filepath:

```bash
tinybench path/to/jsFile.js
```


<img alt="tinybench demo" src="examples/demo.gif" width="700"/>

## Contributing

Feel free to open issues or pull requests with improvements, suggestions, or bug fixes.
