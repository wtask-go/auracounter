package logging

func ExampleNewNull() {
	logger := NewNull()
	defer logger.Close()

	MakeLog(logger)

	// Output:
	//
}
