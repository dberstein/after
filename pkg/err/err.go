package err

import (
	"fmt"
	"os"
)

type Err struct {
	code    int
	message string
}

func New(code int, message string) *Err {
	return &Err{code: code, message: message}
}

func Convert(code int, error error) *Err {
	return New(code, "ERROR: "+error.Error())
}

func (e *Err) Error() string {
	return e.message
}

func (e *Err) Code() int {
	return e.code
}

func (e *Err) Print() *Err {
	_, err := fmt.Fprintf(os.Stderr, "%s\n", e.message)
	if err != nil {
		panic(err)
	}
	return e
}

func (e *Err) Exit() {
	os.Exit(e.code)
}
