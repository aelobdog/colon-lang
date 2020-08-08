package main

import i "colon/colinterp"

func main() {
	code := `
	v: sumToNumLoop = f(num):
		v: iter = 1
		v: sum = 0
		l(iter <= num):
			v: sum = sum + iter
			v: iter = iter + 1
		:l
		r: sum
	:f

	v: sumToNumRec = f(num):
		i(num == 1):
			1
		:i e:
			num + sumToNumRec(num - 1)
		:e
	:f

	v: num = 5

	v: sumL = sumToNumLoop(num)
	print(sumL)
	v: sumR = sumToNumRec(num)
	print(sumR)
	`
	i.Interpret(code)
}
