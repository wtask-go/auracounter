package logging

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func ExampleNewFile_withPrefix() {
	filename := RandomTestdataName(".test.log")
	logger, err := NewFile(
		filename,
		// blow out current time
		WithDefaultDecoration("test", &Timer{func() time.Time { return time.Time{} }, DefaultTimeFormat}),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		logger.Close()
		os.Remove(filename)
	}()
	MakeLog(logger)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))

	// Output:
	// test [0001-01-01 00:00:00.000000] INFO new event {Bar}
	// test [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// test [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// test [0001-01-01 00:00:00.000000] ERR new event {Bar}
	// test [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// test [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}

func ExampleNewFile_withoutPrefix() {
	filename := RandomTestdataName(".test.log")
	logger, err := NewFile(
		filename,
		// blow out current time
		WithDefaultDecoration("", &Timer{func() time.Time { return time.Time{} }, DefaultTimeFormat}),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		logger.Close()
		os.Remove(filename)
	}()
	MakeLog(logger)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))

	// Output:
	// [0001-01-01 00:00:00.000000] INFO new event {Bar}
	// [0001-01-01 00:00:00.000000] INFO event-1 occurred
	// [0001-01-01 00:00:00.000000] INFO event #2 occurred
	// [0001-01-01 00:00:00.000000] ERR new event {Bar}
	// [0001-01-01 00:00:00.000000] ERR error-1 occurred
	// [0001-01-01 00:00:00.000000] ERR error #2 occurred
	//
}

func ExampleNewFile_customTimeForma() {
	filename := RandomTestdataName(".test.log")
	logger, err := NewFile(
		filename,
		// blow out current time
		WithDefaultDecoration(
			"",
			&Timer{
				func() time.Time { return time.Time{}.Add(31 * 24 * time.Hour) }, // 1st february
				"02 Jan 2006 15:04:05",
			},
		),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		logger.Close()
		os.Remove(filename)
	}()
	MakeLog(logger)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))

	// Output:
	// [01 Feb 0001 00:00:00] INFO new event {Bar}
	// [01 Feb 0001 00:00:00] INFO event-1 occurred
	// [01 Feb 0001 00:00:00] INFO event #2 occurred
	// [01 Feb 0001 00:00:00] ERR new event {Bar}
	// [01 Feb 0001 00:00:00] ERR error-1 occurred
	// [01 Feb 0001 00:00:00] ERR error #2 occurred
	//
}

func ExampleNewFile_customDecoration() {
	filename := RandomTestdataName(".test.log")
	logger, err := NewFile(filename, WithDecoration(CustomConstantTimeDecorator()))
	if err != nil {
		panic(err)
	}
	defer func() {
		logger.Close()
		os.Remove(filename)
	}()
	MakeLog(logger)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(buf))

	// Output:
	// : [2006-01-02 15:04:05.000000] new event {Bar}
	// : [2006-01-02 15:04:05.000000] event-1 occurred
	// : [2006-01-02 15:04:05.000000] event #2 occurred
	// ! [2006-01-02 15:04:05.000000] new event {Bar}
	// ! [2006-01-02 15:04:05.000000] error-1 occurred
	// ! [2006-01-02 15:04:05.000000] error #2 occurred
	//
}
