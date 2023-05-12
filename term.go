package main

import (
	"fmt"
	"github.com/mpetavy/common"
	"golang.org/x/exp/slices"
	"math"
	"strconv"
	"strings"
)

const (
	signChs       = "+-"
	exponentChs   = "^"
	punctationChs = "*/%"
	lineChs       = "+-"
	operatorChs   = exponentChs + punctationChs + lineChs
	numberChs     = "0123456789."
)

func isSign(txt string) bool {
	return strings.IndexAny(txt, signChs) != -1
}

func isOperator(txt string) bool {
	return strings.IndexAny(txt, operatorChs) != -1
}

func isNumeric(txt string) bool {
	return strings.IndexAny(txt, numberChs) != -1
}

func isExponent(txt string) bool {
	return strings.IndexAny(txt, exponentChs) != -1
}

func isPunctation(txt string) bool {
	return strings.IndexAny(txt, punctationChs) != -1
}

func isLine(txt string) bool {
	return strings.IndexAny(txt, lineChs) != -1
}

type Term struct {
	Text     string
	IsParsed bool
	Result   float64
	Terms    []Term
}

func (term Term) String() string {
	sb := strings.Builder{}

	if term.Text != "" {
		sb.WriteString(term.Text)
	}

	if len(term.Terms) > 0 {
		sb.WriteString("(")
		for i := 0; i < len(term.Terms); i++ {
			sb.WriteString(term.Terms[i].String())
		}
		sb.WriteString(")")
	}

	return sb.String()
}

func (term *Term) Operator() string {
	ch := term.Text[:1]

	if isOperator(ch) {
		return ch
	}

	return "+"
}

func NewTerm(text string) (Term, error) {
	term, _, err := parse(text, 0)
	if common.Error(err) {
		return term, err
	}

	return term, nil
}

func parse(text string, i int) (Term, int, error) {
	current := Term{}
	v := Term{}

	registerTerm := func() {
		if v.Text != "" || len(v.Terms) > 0 {
			current.Terms = append(current.Terms, v)
			v = Term{}
		}
	}

	lastCh := ""
loop:
	for i < len(text) {
		ch := text[i : i+1]

		switch {
		case ch == "(":
			operator := v.Text

			if isNumeric(operator) {
				registerTerm()

				operator = "*"
			}

			var err error

			v, i, err = parse(text, i+1)
			if common.Error(err) {
				return v, 0, err
			}

			v.Text = operator + v.Text

			registerTerm()
		case ch == ")":
			break loop
		default:
			if isOperator(ch) && !isOperator(lastCh) {
				registerTerm()

				v.Text += ch
			} else {
				v.Text += ch
			}
		}

		lastCh = ch

		i++
	}

	registerTerm()

	return current, i, nil
}

func (term *Term) Calc() (float64, error) {
	if len(term.Terms) == 0 {
		if term.IsParsed {
			return term.Result, nil
		}

		number := term.Text

		if strings.HasPrefix(number, "--") {
			number = fmt.Sprintf("+%s", number[2:])
		}
		if strings.HasPrefix(number, "+-") {
			number = fmt.Sprintf("-%s", number[2:])
		}
		if strings.HasPrefix(number, "-+") {
			number = fmt.Sprintf("+%s", number[2:])
		}
		if strings.HasPrefix(number, "++") {
			number = fmt.Sprintf("+%s", number[2:])
		}

		if !isNumeric(number[:1]) && !isSign(number[:1]) {
			number = number[1:]
		}

		var err error

		term.Result, err = strconv.ParseFloat(number, 64)
		if common.Error(err) {
			return 0, err
		}

		term.IsParsed = true

		return term.Result, nil
	}

	if len(term.Terms) == 1 {
		var err error

		term.IsParsed = true
		term.Result, err = term.Terms[0].Calc()
		if common.Error(err) {
			return 0, err
		}
	}

	for operatorGroup := 0; operatorGroup < 3; operatorGroup++ {
		for i := 0; i < len(term.Terms)-1; {
			c := 2
			left := term.Terms[i]
			right := term.Terms[i+1]

			operator := right.Operator()

			if (operatorGroup == 0 && isExponent(operator)) ||
				(operatorGroup == 1 && isPunctation(operator)) ||
				(operatorGroup == 2 && isLine(operator)) {

				var result float64

				leftTerm, err := left.Calc()
				if common.Error(err) {
					return 0, err
				}

				rightTerm, err := right.Calc()
				if common.Error(err) {
					return 0, err
				}

				neg := isPunctation(operator) && leftTerm < 0
				if neg {
					leftTerm = math.Abs(leftTerm)
				}

				switch operator {
				case "^":
					result = math.Pow(leftTerm, rightTerm)
				case "+":
					result = leftTerm + rightTerm
				case "-":
					result = leftTerm + rightTerm
				case "*":
					result = leftTerm * rightTerm
				case "/":
					result = leftTerm / rightTerm
				case "%":
					result = float64(int(leftTerm) % int(rightTerm))
				}

				if neg {
					result = result * -1
				}

				newTerm := Term{}

				if result >= 0 {
					newTerm.Text = fmt.Sprintf("+%v", result)
				} else {
					newTerm.Text = fmt.Sprintf("%v", result)
				}
				newTerm.IsParsed = true
				newTerm.Result = result

				term.Terms = slices.Delete(term.Terms, i, i+c)
				term.Terms = slices.Insert(term.Terms, i, newTerm)

				continue
			}

			i++
		}
	}

	if term.Text != "" && term.Operator() == "-" {
		term.Terms[0].Result = term.Terms[0].Result * -1
	}

	return term.Terms[0].Result, nil
}
