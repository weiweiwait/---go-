package stack

const MAXVEX1 = 30

type Stack struct {
	Num [MAXVEX1]int
	Top int
}

func (s *Stack) InSert(num int) {
	s.Top++
	s.Num[s.Top] = num
}

func (s *Stack) Out(num *int) {
	*num = s.Num[s.Top]
	s.Top--
}

func (s *Stack) Gettop(num *int) {
	*num = s.Num[s.Top]
}

func (s *Stack) IsEmpty() bool {
	if s.Top != -1 {
		return false
	} else {
		return true
	}
}
