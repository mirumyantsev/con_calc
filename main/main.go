package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Processing interface {
	valueInputing()
	valueProcessing()
	valuePrinting()
}

type InputProcessor struct {
	availableOps  [6]string
	errors        [5]string
	errNo         int
	value         string
	expr          string
	lowBound      int
	highBound     int
	closeBrackets int
	openBracket   int
}

func DoProcessing(ip *InputProcessor) {
	ip.availableOps = [6]string{"+", "-", "*", "/", "^", "%"}
	ip.errors = [5]string{
		"введен неверный символ.",
		"количество открывающих скобок не равно количеству закрывающих.",
		"ошибка буфера ввода.",
		"ошибка в записи числа.",
		"на ноль делить нельзя.",
	}

	for {
		ip.errNo = 0
		ip.expr = ""
		ip.lowBound = 0
		ip.highBound = 0
		ip.closeBrackets = 0
		ip.openBracket = 0

		ip.valueInputing()
		ip.valueProcessing()
		ip.valuePrinting()
	}
}

func (ip *InputProcessor) valueProcessing() {
	ip.deleteWhitespaces()
	ip.deleteLineFeeds()
	ip.checkUserWantQuit()
	if !(ip.hasWrongChars()) {
		ip.calculateBrackets()
	}
}

func (ip *InputProcessor) hasWrongChars() bool {
	for i := 0; i < len(ip.value); i++ {
		char := string(ip.value[i])
		if !((char == ".") || (char == "(") || (char == ")") || ip.isAvailableDigit(char) || ip.isAvailableOp(char)) {
			ip.errNo = 1
			return true
		}
	}
	return false
}

func (ip *InputProcessor) calculateBrackets() {
	lenOfValue := len(ip.value)
	var openBrackets, closeBrackets int

	for i := 0; i < lenOfValue; i++ {
		char := string(ip.value[i])
		if char == ")" {
			closeBrackets++
		} else if char == "(" {
			openBrackets++
		}
	}

	if openBrackets != closeBrackets {
		ip.errNo = 2
		return
	}

	for closeBrackets > 0 {
		ip.performBracketOp()
		closeBrackets--
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
			ip.closeBrackets = j

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

	ip.expr = ip.value[ip.openBracket+1 : ip.closeBrackets]
	ip.calculateExpr()
	ip.value = ip.value[:ip.openBracket] + ip.expr + ip.value[ip.closeBrackets+1:]
}

func (ip *InputProcessor) calculateExpr() {
	lenOfExpr := len(ip.expr)
	var highOps, lowOps int

	for i := 0; i < lenOfExpr; i++ {
		char := string(ip.expr[i])
		if (char == ip.availableOps[2]) || (char == ip.availableOps[3]) || (char == ip.availableOps[4]) || (char == ip.availableOps[5]) {
			highOps++
		} else if (char == ip.availableOps[0]) || (char == ip.availableOps[1]) {
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
		if (char == ip.availableOps[2]) || (char == ip.availableOps[3]) || (char == ip.availableOps[4]) || (char == ip.availableOps[5]) {
			var i int

			for i = j - 1; i >= 0; i-- {
				char := string(ip.expr[i])
				if ip.isAvailableOp(char) {
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
				if ip.isAvailableOp(char) {
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
		if (char == ip.availableOps[0]) || (char == ip.availableOps[1]) {
			var i int

			for i = j - 1; i >= 0; i-- {
				char := string(ip.expr[i])
				if ip.isAvailableOp(char) {
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
				if ip.isAvailableOp(char) {
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
	var err error

	for i := 0; i < len(binExpr); i++ {
		char := string(binExpr[i])
		if ip.isAvailableOp(char) {
			operatorChar = char
			operatorPos = i
			break
		}
	}

	operand_1 := binExpr[:operatorPos]
	operand_2 := binExpr[operatorPos+1:]

	operand_1f, err := strconv.ParseFloat(operand_1, 64)
	operand_2f, err := strconv.ParseFloat(operand_2, 64)

	if err != nil {
		ip.errNo = 4
	}

	switch operatorChar {
	case ip.availableOps[0]:
		resultf = ip.add(operand_1f, operand_2f)
	case ip.availableOps[1]:
		resultf = ip.subtract(operand_1f, operand_2f)
	case ip.availableOps[2]:
		resultf = ip.multiply(operand_1f, operand_2f)
	case ip.availableOps[3]:
		resultf = ip.divide(operand_1f, operand_2f)
	case ip.availableOps[4]:
		resultf = ip.power(operand_1f, operand_2f)
	case ip.availableOps[5]:
		resultf = ip.modulo(operand_1f, operand_2f)
	default:
		ip.errNo = 1
	}

	return fmt.Sprintf("%.3f", resultf)
}

func (ip *InputProcessor) isAvailableOp(op string) bool {
	for _, availableOp := range ip.availableOps {
		if availableOp == op {
			return true
		}
	}
	return false
}

func (ip *InputProcessor) isAvailableDigit(char string) bool {
	availableDigits := [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	for _, availableDigit := range availableDigits {
		if availableDigit == char {
			return true
		}
	}
	return false
}

func (ip *InputProcessor) modulo(operand_1, operand_2 float64) float64 {
	return math.Mod(operand_1, operand_2)
}

func (ip *InputProcessor) power(operand_1, operand_2 float64) float64 {
	return math.Pow(operand_1, operand_2)
}

func (ip *InputProcessor) divide(operand_1, operand_2 float64) float64 {
	if operand_2 != 0.0 {
		return operand_1 / operand_2
	} else {
		ip.errNo = 5
		return 0.0
	}
}

func (ip *InputProcessor) multiply(operand_1, operand_2 float64) float64 {
	return operand_1 * operand_2
}

func (ip *InputProcessor) subtract(operand_1, operand_2 float64) float64 {
	return operand_1 - operand_2
}

func (ip *InputProcessor) add(operand_1, operand_2 float64) float64 {
	return operand_1 + operand_2
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
	if ip.errNo > 0 {
		fmt.Printf("Ошибка: \n- %s\n", ip.errors[ip.errNo-1])
		return
	}
	fmt.Println("Результат:")
	fmt.Println(ip.value)
}

func (ip *InputProcessor) valueInputing() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите выражение:")
	var err error
	ip.value, err = reader.ReadString('\n')
	if err != nil {
		ip.errNo = 3
	}

}

func main() {
	S := InputProcessor{}
	DoProcessing(&S)
}
