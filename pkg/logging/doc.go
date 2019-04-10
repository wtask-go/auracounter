/*
Package logging exposes common logging facade and interface to manage and to use loggers within application.

It provides several lightweight and customizable loggers, which output targets are:

	- stdout
	- stderr
	- null
	- bytes buffer
	- local file

All provided loggers support decorators to make format of your log rows highly customizable.
All loggers use system log.Logger as backend, except already planned syslog logger.
*/
package logging
