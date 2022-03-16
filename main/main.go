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
	value string
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
	ip.doSimpleCalc()
}

func (ip *InputProcessor) doSimpleCalc() {
	var opPos int

	for i := 0; i < len(ip.value); i++ {
		char := string(ip.value[i])
		if (char == "-") || (char == "+") {
			opPos = i
			break
		}
	}

	if opPos != 0 {
		operand1 := ip.value[:opPos]
		operand2 := ip.value[opPos+1:]
		operand1_f, _ := strconv.ParseFloat(operand1, 32)
		operand2_f, _ := strconv.ParseFloat(operand2, 32)
		fmt.Printf("%f\n", sum(operand1_f, operand2_f))

	}

}

func sum(addndm_1, addndm_2 float64) (result float64) {
	result = addndm_1 + addndm_2
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
