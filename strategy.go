package main

type Strategy interface {
	Init(base *wordList)
	GetFirstMoveTries() ([]string, error)
	GetNextMoveTries(move int, lastTry, gameResponse string) ([]string, error)
}

type stratError struct {
	err     error
	isFatal bool
}

func (se *stratError) Error() string {
	return se.err.Error()
}

func isFatalStratError(err error) bool {
	if err, ok := err.(*stratError); ok {
		return err.isFatal
	} else {
		return false
	}
}

func newStratError(err error, isFatal bool) *stratError {
	return &stratError{err, isFatal}
}
