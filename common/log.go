package common

import "io"

type SimpleLogger struct {
	logger *Logger
}

func NewSimpleLogger(w io.Writer) *SimpleLogger {
	return &SimpleLogger{
		logger: New(w, "", Lshortfile|Ldate|Lmicroseconds)}
}

func (s *SimpleLogger) Debugf(f string, m ...interface{}) (err error) {
	s.logger.Debugf(f, m...)
	return
}

func (s *SimpleLogger) Debug(m ...interface{}) (err error) {
	s.logger.Debug(m...)
	return
}

func (s *SimpleLogger) Err(m ...interface{}) (err error) {
	s.logger.Err(m...)
	return
}

func (s *SimpleLogger) Info(m ...interface{}) (err error) {
	s.logger.Info(m...)
	return
}

func (s *SimpleLogger) Warning(m ...interface{}) (err error) {
	s.logger.Warning(m...)
	return
}
