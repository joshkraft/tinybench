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
