package input_processor

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
