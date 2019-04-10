package logging

func ExampleNewNull() {
	logger := NewNull()
	defer logger.Close()
	logger.Infof("event-1 occurred")
	logger.Infof("event #%d occurred", 2)
	logger.Errorf("error-1 occurred")
	logger.Errorf("error #%d occurred", 2)

	// Output:
	//
}

func ExampleNewNull_asFacade() {
	logger := NewNull()
	defer logger.Close()

	func(f Facade) {
		f.Infof("event-1 occurred")
		f.Infof("event #%d occurred", 2)
		f.Errorf("error-1 occurred")
		f.Errorf("error #%d occurred", 2)
	}(logger)

	// Output:
	//
}
