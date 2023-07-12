Go Futures
This Go package provides a simple implementation of futures, allowing you to execute functions asynchronously and retrieve their results.

Installation
To use this package, you need to have Go installed on your system. You can install it by following the official Go installation instructions.

Next, you can install the package using the go get command:

go get github.com/your-username/go-futures



Usage
Import the package into your Go code:

import (
	"context"
	"fmt"
	"time"

	"github.com/your-username/go-futures"
)


The package provides a futureInterface and futureStruct types for working with futures.

Here's an example of using the futures package:

func longRunningTask() (int, error) {
	// Simulate a time-consuming operation
	time.Sleep(5 * time.Second)
	return 42, nil
}

func main() {
	f := futures.New(func() (interface{}, error) {
		return longRunningTask()
	})

	// Checking cancel call
	go func() {
		time.Sleep(2 * time.Second)
		f.Cancel()
	}()

	result, err := f.Result()
	fmt.Println(result, err, f.Cancelled())

	// Checking get call
	g := futures.New(func() (interface{}, error) {
		return longRunningTask()
	})
	gResult, gErr := g.Result()

	fmt.Println(g.Done(), gResult, gErr, g.Cancelled())
}



In the example above, we create a future using New and pass a function to execute asynchronously. The longRunningTask function represents a time-consuming operation. We can then retrieve the result using the Result method.

Additionally, the code demonstrates canceling a future using the Cancel method and checking the status using Cancelled and Done methods.

For more details, please refer to the Go Futures documentation.


Contributing
Contributions are welcome! If you find any issues or have suggestions for improvement, please create an issue or submit a pull request.

License
This project is licensed under the MIT License.

