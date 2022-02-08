# warnings

## 基础知识

### 1. 变量的作用域

代码环境：
```go
func fun() (int, int) {
	return 5, 5
}
func test() {
	i := 1
	for tms := 0; tms < 1; tms++ {
		cache, i := fun()
		if cache == 5 {
			break
		}
	}
	fmt.Println(i)
}
```
上面的代码将会报错：`i decalred but not used`（这个`i`是循环内的i），所以需要单独进行声明再进行函数调用并赋返回值