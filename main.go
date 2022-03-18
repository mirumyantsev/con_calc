package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/mikerumy/concalc/consts"
)

type Processing interface {
	valueInputing()
	valueProcessing()
	valuePrinting()
}

type InputProcessor struct {
	availOps      [6]string
	errors        [5]string
	errNo         int
	value         string
	expr          string
	lowBound      int
	highBound     int
	closeBrackets int
	openBracket   int
}

func (ip *InputProcessor) DoProcessing() {
	ip.availOps = [6]string{"+", "-", "*", "/", "^", "%"}
	ip.errors = [5]string{
		"введен неверный символ.",
		"количество открывающих скобок не равно количеству закрывающих.",
		"ошибка буфера ввода.",
		"ошибка в записи числа.",
		"на ноль делить нельзя.",
	}

	for {
		ip.errNo = 0

		ip.valueInputing()
		ip.valueProcessing()
		ip.valuePrinting()
	}
}

func (ip *InputProcessor) valueProcessing() {
	ip.deleteWhitespaces()
	ip.deleteLineFeeds()
	ip.checkUserWantQuit()
	ip.constantConverter()
	if !ip.hasWrongChars() {
		ip.calculateBrackets()
	}
}

func (ip *InputProcessor) constantConverter() {
	consts.Init()

	var piCount, eCount int

	for i := 0; i < len(ip.value); i++ {
		char := string(ip.value[i])
		if char == "e" {
			eCount++
		}
	}

	for i := 0; i < len(ip.value)-1; i++ {
		word := string(ip.value[i : i+2])
		if word == "pi" {
			piCount++
		}
	}

	for eCount > 0 {
		for i := 0; i < len(ip.value); i++ {
			char := string(ip.value[i])
			if char == "e" {
				ip.value = ip.value[:i] + consts.Values["e"] + ip.value[i+1:]
				break
			}
		}
		eCount--
	}

	for piCount > 0 {
		for i := 0; i < len(ip.value)-1; i++ {
			word := string(ip.value[i : i+2])
			if word == "pi" {
				ip.value = ip.value[:i] + consts.Values["pi"] + ip.value[i+2:]
				break
			}
		}
		piCount--
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
	var lowestOps, lowOps, highOps, highestOps int

	for i := 0; i < lenOfExpr; i++ {
		char := string(ip.expr[i])
		if char == ip.availOps[4] {
			highestOps++
		} else if (char == ip.availOps[2]) || (char == ip.availOps[3]) {
			highOps++
		} else if (char == ip.availOps[0]) || (char == ip.availOps[1]) {
			lowOps++
		} else if char == ip.availOps[5] {
			lowestOps++
		}
	}

	for highestOps > 0 {
		ip.performOp([]string{ip.availOps[4]})
		highestOps--
	}

	for highOps > 0 {
		ip.performOp([]string{ip.availOps[2], ip.availOps[3]})
		highOps--
	}

	for lowOps > 0 {
		ip.performOp([]string{ip.availOps[0], ip.availOps[1]})
		lowOps--
	}

	for lowestOps > 0 {
		ip.performOp([]string{ip.availOps[5]})
		lowestOps--
	}
}

func (ip *InputProcessor) performOp(availOps []string) {
	lenOfExpr := len(ip.expr)
	var passToDoOp bool

	for j := 0; j < lenOfExpr; j++ {
		char := string(ip.expr[j])
		for _, op := range availOps {
			if char == op {
				passToDoOp = true
			}
		}

		if passToDoOp {
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
	case ip.availOps[0]:
		resultf = operand_1f + operand_2f
	case ip.availOps[1]:
		resultf = operand_1f - operand_2f
	case ip.availOps[2]:
		resultf = operand_1f * operand_2f
	case ip.availOps[3]:
		if operand_2f != 0.0 {
			resultf = operand_1f / operand_2f
		} else {
			ip.errNo = 5
		}
	case ip.availOps[4]:
		resultf = math.Pow(operand_1f, operand_2f)
	case ip.availOps[5]:
		resultf = math.Mod(operand_1f, operand_2f)
	default:
		ip.errNo = 1
	}

	return fmt.Sprintf("%.3f", resultf)
}

func (ip *InputProcessor) isAvailableOp(op string) bool {
	for _, availableOp := range ip.availOps {
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
	var S InputProcessor
	S.DoProcessing()
}