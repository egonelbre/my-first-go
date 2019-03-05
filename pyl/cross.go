package main

import . "fmt"
// Problem no. 15
// --------------------
// compilation
// 8g cross.go && 8l -o cross cross.8

func cmod7(number int) bool{
    t := 0
    for number > 0 {
        m := number % 10
        t += m
        number = (number - m)/ 10
    }
    return t % 7 == 0
}

func summod7s(numbers []int) (result int) {
    result = 0
    for _,v := range(numbers){
        if cmod7(v) {
            result += v
        }
    }
    return
}

func test7(value int, expect bool){
    result := cmod7(value)
    var s string
    if result == expect {
        s = "OK"
    } else {
        s = "ERR"
    }
    Printf("%4d => %5v %s\n", value, result, s)
}

func testsum(values []int, expect int){
    result := summod7s(values)
    var s string
    if result == expect {
        s = "OK"
    } else {
        s = "ERR"
    }
    Printf("%v => %4d %s\n", values, result, s)
}

func main() {
    test7(86, true)
    test7(17, false)
    test7(7, true)
    test7(91, false)
    // [...]int{1,2,3,4} creates array
    // [:] converts it to a slice
    testsum([...]int{52,31,10,86}[:],138)
    testsum([...]int{6,7,8,9,11,16,25}[:],48)
    testsum([...]int{7}[:],7)
}
