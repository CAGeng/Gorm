package log

import "fmt"

type Tprinter struct {
	Indentlevel int
}

var (
	Mytprinter = NewTprinter()
)

func NewTprinter()*Tprinter{
	return &Tprinter{
		Indentlevel: 0,
	}
}

func (p *Tprinter)IndentLvUp(){
	p.Indentlevel += 4
}

func (p *Tprinter)IndentLvDown(){
	if p.Indentlevel > 0{
		p.Indentlevel -= 4
	}
}

func (p *Tprinter)Print(s interface{}){
	for i := 0;i < p.Indentlevel;i++{
		fmt.Print(" ")
	}
	fmt.Print(s)
	fmt.Print("\n")
}