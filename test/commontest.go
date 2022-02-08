package test

//////////////// test two-dimensional array ////////////////
func test2DArray() {
	Log("####### test two-dimensional array")
	// define 2D array
	var twoDarr [2][]string
	if len(twoDarr[0]) != 0 || len(twoDarr[1]) != 0 {
		Linfo("2Darray is not null. length:", len(twoDarr[0]), len(twoDarr[1]))
	} else {
		Linfo("2Darray is null.")
	}
	Logf("len 2darray: %d, child array len: %d, %d, content of twoDarr: %+v", len(twoDarr), len(twoDarr[0]), len(twoDarr[1]), twoDarr)
	twoDarr[0] = append(twoDarr[0], "123")
	Logf("len 2darray: %d, content of twoDarr: %+v", len(twoDarr), twoDarr)
	twoDarr[1] = append(twoDarr[1], "456")
	Logf("len 2darray: %d, content of twoDarr: %+v", len(twoDarr), twoDarr)
}

//////////////// test slice format ////////////////

// 测试切片slice格式
//  a[x:y:z] 切片内容 [x:y] 切片长度: y-x 切片容量:z-x
func testSliceFormat() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	anotherSlice := slice[:]
	Linfo("another slice:", anotherSlice)
	d1 := slice[6:8]
	Linfo(d1, len(d1), cap(d1))
	d2 := slice[:6:8]
	Linfo(d2, len(d2), cap(d2))
}

//////////////// test const value ////////////////

const (
	con_var_0 = iota // be init to 0
	con_var_1 = 100  // be init to 100
	con_var_2 = iota // be init to 2
	con_var_3        // be init to 3
	con_var_4
	con_var_5
	con_var_6
	con_var_7
)

func testConstValue() {
	Log("const value of iota(begin):", con_var_0)
	Log("const value of below list:",
		"\nvar_1 =", con_var_1,
		", var_2 =", con_var_2,
		", var_3 =", con_var_3,
		", var_4 =", con_var_4,
		", var_5 =", con_var_5,
		", var_6 =", con_var_6,
		", var_7 =", con_var_7,
	)
}
