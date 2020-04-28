# WrapErr
WrapErr is a Golang library for working with errors.  

## How to get
```
go get github.com/Chekunin/wraperr
```

## Usage
Let's look at a simple example, where you make an error tracing.
```go
import (
    "github.com/Chekunin/wraperr
    "errors"
)

func main() {
	if err := somethingGoWrong(); err != nil {
		err = NewWrapErr(fmt.Errorf("somethingGoWrong"), err)
		fmt.Println(err) // somethingGoWrong: any error
	}
}

func somethingGoWrong() error {
	return errors.New("any error")
}
```

But with WrapErr you can also operate with any previous error, therefore, for instance, you can send corresponding error from your API.
```go
var errDivisionByZero = errors.New("division by zero")
func divide(numerator int, denominator int) (int, error) {
	if denominator == 0 {
		return 0, errDivisionByZero
	}
	res := numerator / denominator
	return res, nil
}

func main() {
    http.HandleFunc("/divide", func(w http.ResponseWriter, r *http.Request) {
        numerator := getNumeratorFromRequest(r) // e.g. numerator = 5
        denominator := getDenominatorFromRequest(r) // e.g. denominator = 0
        res, err := divide(numerator, denominator)
        if err != nil {
            err = NewWrapErr(fmt.Errorf("devide with params %d and %d", denominator, denominator), err)
            if errors.Is(err, errDivisionByZero) {
                http.Error(w, "bad request", 400)
            } else {
                http.Error(w, "internal server error", 500)
            }
            log.Print(err) // devide with params 5 and 0: division by zero
            return
        }
        fmt.Fprintf(w, "the answer is %d", res)
    })
    http.ListenAndServe(":8080", nil)
}
``` 
Also WrapErr supports Golang's 1.13 _errors.As_ function.  
Thus you can not only work with just a textual view of previous errors, but with original objects of those errors.