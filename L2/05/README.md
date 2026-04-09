# L2.5

Что выведет программа?

Объяснить вывод программы.

```go
package main

type customError struct {
    msg string
}

func (e *customError) Error() string {
    return e.msg
}

func test() *customError {
    // ... do something
    return nil
}

func main() {
    var err error
    err = test()
    if err != nil {
        println("error")
        return
    }
    println("ok")
}
```
---

## Вывод

```
error
```

## Объяснение

Данный пример похож на пример [L2.3](https://github.com/imotkin/L2/03) в том смысле, что переменная err имеет тип интерфейса `error` и хотя из функции `test` возвращается `nil`, но внутри интерфейса имеется тип `*customError` и поэтому проверка на `nil` будет `false`, то есть `err` не равен `nil`, так как при сравнении переменной интерфейса с `nil` сравнивается как тип внутри интерфейса, так и его значение. 