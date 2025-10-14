package library

var x int32 = 1e9 + 7
var y int32 = 151

func addSymbol(hash *[2]int64, symbol int64) {
	(*hash)[0] *= int64(y)
	(*hash)[0] %= int64(x)
	(*hash)[0] += symbol
	(*hash)[0] %= int64(x)
}
