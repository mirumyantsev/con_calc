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
	value   string
	binExpr string
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
	ip.calculationCycles()

	// ip.value = ip.doBinaryOp(ip.value)
}

func (ip *InputProcessor) calculationCycles() {
	lenOfValue := len(ip.value)
	var highOps, lowOps int

	for i := 0; i < lenOfValue; i++ {
		char := string(ip.value[i])
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
	lenOfValue := len(ip.value)
	var lowBound, highBound, i int

	for j := 0; j < lenOfValue; j++ {
		char := string(ip.value[j])
		if (char == "*") || (char == "/") {

			for i = j - 1; i >= 0; i-- {
				char := string(ip.value[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			if i > 0 {
				lowBound = i + 1
			} else {
				lowBound = 0
			}

			for i = j + 1; i < lenOfValue; i++ {
				char := string(ip.value[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			highBound = i - 1
		}
	}

	binExpr := ip.doBinaryOp(ip.value[lowBound : highBound+1])
	ip.value = ip.value[:lowBound] + binExpr + ip.value[highBound+1:]
}

func (ip *InputProcessor) performLowOp() {
	lenOfValue := len(ip.value)
	var lowBound, highBound, i int

	for j := 0; j < lenOfValue; j++ {
		char := string(ip.value[j])
		if (char == "-") || (char == "+") {

			for i = j - 1; i >= 0; i-- {
				char := string(ip.value[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			if i > 0 {
				lowBound = i + 1
			} else {
				lowBound = 0
			}

			for i = j + 1; i < lenOfValue; i++ {
				char := string(ip.value[i])
				if (char == "-") || (char == "+") || (char == "*") || (char == "/") {
					break
				}
			}
			highBound = i - 1
		}
	}

	binExpr := ip.doBinaryOp(ip.value[lowBound : highBound+1])
	ip.value = ip.value[:lowBound] + binExpr + ip.value[highBound+1:]
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
