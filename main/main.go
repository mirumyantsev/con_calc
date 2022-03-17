package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Processing interface {
	valueInputing()
	valueProcessing()
	valuePrinting()
}

type InputProcessor struct {
	value        string
	expr         string
	lowBound     int
	highBound    int
	closeBracket int
	openBracket  int
}

func DoProcessing(ip *InputProcessor) {
	for {
		ip.valueInputing()
		ip.valueProcessing()
		ip.valuePrinting()
	}
}

func (ip *InputProcessor) valueProcessing() {
	ip.deleteWhitespaces()
	ip.deleteLineFeeds()
	ip.checkUserWantQuit()

	ip.calculateBrackets()
}

func (ip *InputProcessor) calculateBrackets() {
	lenOfValue := len(ip.value)
	var bracketOps int

	for i := 0; i < lenOfValue; i++ {
		char := string(ip.value[i])
		if char == ")" {
			bracketOps++
		}
	}

	for bracketOps > 0 {
		ip.performBracketOp()
		bracketOps--
	}

	ip.expr = ip.value
	ip.calculateExpr()
	ip.value = ip.expr
}

func (ip *InputProcessor) performBracketOp() {
	lenOfValue := len(ip.value)

	for j := 0; j < lenOfValue; j++ {
		char := string(ip.value[j])
		if char == ")" {
			ip.closeBracket = j

			for i := j - 1; i >= 0; i-- {
				char := string(ip.value[i])
				if char == "(" {
					ip.openBracket = i
					break
				}
			}
			break
		}
	}

	ip.expr = ip.value[ip.openBracket+1 : ip.closeBracket]
	ip.calculateExpr()
	ip.value = ip.value[:ip.openBracket] + ip.expr + ip.value[ip.closeBracket+1:]
}

func (ip *InputProcessor) calculateExpr() {
	lenOfExpr := len(ip.expr)
	var highOps, lowOps int

	for i := 0; i < lenOfExpr; i++ {
		char := string(ip.expr[i])
		if (char == "*") || (char == "/") {
			highOps++
		} else if (char == "-") || (char == "+") {
			lowOps++
		}
	}

	for highOps > 0 {
		ip.performHighOp()
		highOps--
	}

	for lowOps > 0 {
		ip.performLowOp()
		lowOps--
	}
}

func (ip *InputProcessor) performHighOp() {
	lenOfExpr := len(ip.expr)

	for j := 0; j < lenOfExpr; j++ {
		char := string(ip.expr[j])
		if (char == "*") || (char == "/") {
			var i int

			for i = j - 1; i >= 0; i-- {
				char := string(ip.expr[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			if i > 0 {
				ip.lowBound = i + 1
			} else {
				ip.lowBound = 0
			}

			for i = j + 1; i < lenOfExpr; i++ {
				char := string(ip.expr[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			ip.highBound = i - 1

			break
		}
	}

	binOp := ip.doBinaryOp(ip.expr[ip.lowBound : ip.highBound+1])
	ip.expr = ip.expr[:ip.lowBound] + binOp + ip.expr[ip.highBound+1:]
}

func (ip *InputProcessor) performLowOp() {
	lenOfExpr := len(ip.expr)

	for j := 0; j < lenOfExpr; j++ {
		char := string(ip.expr[j])
		if (char == "-") || (char == "+") {
			var i int

			for i = j - 1; i >= 0; i-- {
				char := string(ip.expr[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			if i > 0 {
				ip.lowBound = i + 1
			} else {
				ip.lowBound = 0
			}

			for i = j + 1; i < lenOfExpr; i++ {
				char := string(ip.expr[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			ip.highBound = i - 1

			break
		}
	}

	binOp := ip.doBinaryOp(ip.expr[ip.lowBound : ip.highBound+1])
	ip.expr = ip.expr[:ip.lowBound] + binOp + ip.expr[ip.highBound+1:]
}

func (ip *InputProcessor) doBinaryOp(binExpr string) (result string) {
	var operatorChar string
	var operatorPos int
	var resultf float64

	for i := 0; i < len(binExpr); i++ {
		char := string(binExpr[i])
		if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
			operatorChar = char
			operatorPos = i
			break
		}
	}

	operand_1 := binExpr[:operatorPos]
	operand_2 := binExpr[operatorPos+1:]

	operand_1f, _ := strconv.ParseFloat(operand_1, 64)
	operand_2f, _ := strconv.ParseFloat(operand_2, 64)

	switch operatorChar {
	case "+":
		resultf = ip.add(operand_1f, operand_2f)
	case "-":
		resultf = ip.subtract(operand_1f, operand_2f)
	case "*":
		resultf = ip.multiply(operand_1f, operand_2f)
	case "/":
		resultf = ip.divide(operand_1f, operand_2f)
	}

	result = fmt.Sprintf("%.3f", resultf)

	return result
}

func (ip *InputProcessor) divide(addendum_1, addendum_2 float64) (result float64) {
	result = addendum_1 / addendum_2
	return result
}

func (ip *InputProcessor) multiply(addendum_1, addendum_2 float64) (result float64) {
	result = addendum_1 * addendum_2
	return result
}

func (ip *InputProcessor) subtract(addendum_1, addendum_2 float64) (result float64) {
	result = addendum_1 - addendum_2
	return result
}

func (ip *InputProcessor) add(addendum_1, addendum_2 float64) (result float64) {
	result = addendum_1 + addendum_2
	return result
}

func (ip *InputProcessor) checkUserWantQuit() {
	if (ip.value == "q") || (ip.value == "exit") {
		os.Exit(0)
	}
}

func (ip *InputProcessor) deleteLineFeeds() {
	ip.value = ip.value[:len(ip.value)-1]
}

func (ip *InputProcessor) deleteWhitespaces() {
	var tempStr string

	for i := 0; i < len(ip.value); i++ {
		char := string(ip.value[i])
		if char != " " {
			tempStr += char
		}
	}

	ip.value = tempStr
}

func (ip *InputProcessor) valuePrinting() {
	fmt.Println("Результат:")
	fmt.Println(ip.value)
}

func (ip *InputProcessor) valueInputing() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите выражение:")
	ip.value, _ = reader.ReadString('\n')
}

func main() {
	S := InputProcessor{}
	DoProcessing(&S)
}
