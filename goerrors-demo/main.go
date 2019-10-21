package main

import (
	"errors"
	"fmt"

	pkgerr "github.com/pingcap/errors"
)
type Err struct {
	Code int
	Msg string
}
func (e *Err) Error() string  {
	return fmt.Sprintf("code : %d ,msg:%s",e.Code,e.Msg)
}
var A_ERR = &Err{-1,"error"}
func a()  error {
	return A_ERR
}

func b()  error {
	err := a()
	return fmt.Errorf("access denied: %w", err) //使用fmt.Errorf wrap 另一个错误
}

func stackfn1() error {
	return  pkgerr.WithStack(A_ERR)
}

func main()  {
	err := b()
	er := errors.Unwrap(err)
	fmt.Println(er ==A_ERR )


	fmt.Println(errors.Is(err,A_ERR))
	var e = &Err{}
	fmt.Println( errors.As(err, &e))
	if errors.As(err, &e) {
		fmt.Printf("code : %d ,msg:%s",e.Code,e.Msg)
	}


	era := pkgerr.WithStack(err)

	fmt.Printf("%+v",era)
}