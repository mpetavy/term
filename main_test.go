package main

import (
	"github.com/mpetavy/common"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func Test_calculate(t *testing.T) {
	list := []struct {
		term   string
		result float64
	}{
		{"2(3)", 6},
		{"1-(12*3-(11+5))/8", -1.5},
		{"-3*-5", 15},
		{"2*(3+1)/(6-2)-1", 1},
		{"2+3*4-5/1+6", 15},
		{"1*3+2/1+7-6", 6},
		{"3/2*5-4+1*2", 5.5},
		{"6-2+1/3*9-1", 6},
		{"5/2-1+3*2+7", 14.5},
		{"1-3/2+4*6/3", 7.5},
		{"2*3+1/2-5+7/2", 5},
		{"7-4*2/4+3*6", 23},
		{"1+5/2-4*1+8/4", 1.5},
		{"4/2*6-7+1*3", 8},
		{"6+4*1/2-8+5/2", 2.5},
		{"2*4-5/5+1*8", 15},
		{"3/1-2*3+4+2/2", 2},
		{"5-3/2*4+1/2*6", 2},
		{"1/2*6+3-1/2+7*1/2", 9},
		{"3/2+2*3-1/2+5/2", 9.5},
		{"1+2*3-4/2+6/3", 7},
		{"5/2-1+4/2*3+7", 14.5},
		{"4-2/2+6*2/4+7", 13},
		{"8/2*3+1-5/5*2", 11},
		{"(3-1)*4/2-1+7/2", 6.5},
		{"2*3/(1-2)+6*2", 6},
		{"2+4*3/2-(1+1)*4", 0},
		{"5+4*2+(2+1)*2-4", 15},
		{"4/2+3*2+(7-3)/2", 10},
		{"(8-4)/2+3*2+1/2", 8.5},
		{"1+5/2+(4-1)*2+7/2", 13},
		{"2+4*2/4-(1+1)*4", -4},
		{"3*2/(1+2)+6*2", 14},
		{"(1+2)*(3-1)/2+7/2", 6.5},
		{"2*3/(1-2)+7*2", 8},
		{"3+2/2*(2+3)+1/2", 8.5},
		{"(2+1)*4/2+1+7/2", 10.5},
		{"2+4/2+(4-1)*2+7/2", 13.5},
		{"1+2/2*(3+1)-4+6/3", 3},
		{"(5+3)*(2-1)/2-1/2", 3.5},
		{"3+2*3/2+(2-1)*4+7/2", 13.5},
	}

	type test struct {
		name    string
		term    string
		want    float64
		wantErr bool
	}

	tests := []test{}

	for i := 0; i < len(list); i++ {
		//ba, err := common.URLGet(fmt.Sprintf("https://api.mathjs.org/v4?expr=%s", url.QueryEscape(list[i].term)))
		//if common.Error(err) {
		//	t.Fatalf("%v", err)
		//}
		//
		//fmt.Printf("{\"%s\",%s},\n", list[i].term, string(ba))

		tests = append(tests, test{
			name:    strconv.Itoa(i),
			term:    list[i].term,
			want:    list[i].result,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term, err := NewTerm(tt.term)
			require.NoError(t, err)

			got, err := term.Calc()
			require.NoError(t, err)

			require.Equal(t, tt.want, got, "calc(%s) got = %v, want %v", tt.term, got, tt.want)
		})
	}

	sb := strings.Builder{}

	for i, t := range list {
		if i > 0 {
			r := common.Rnd(4)

			switch r {
			case 0:
				sb.WriteString("+")
			case 1:
				sb.WriteString("-")
			case 2:
				sb.WriteString("*")
			case 3:
				sb.WriteString("/")
			}
		}

		sb.WriteString(t.term)
	}

	s := sb.String()

	term, err := NewTerm(s)
	require.NoError(t, err)

	_, err = term.Calc()
	require.NoError(t, err)
}
