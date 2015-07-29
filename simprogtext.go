package simprogtext

import (
    "fmt"
    "bytes"
)

var Indent = "    "

type SourceLine struct {
    indent int
    Line string
}

type SimProgFile struct {
    indentLevel uint
    lines []*SourceLine
}

func (sf *SimProgFile) addLine(line string, args ...interface{}) {
    formatedLine := fmt.Sprintf(line, args...)
    if sf.Lines == nil {
        sf.Lines = make([]*SourceLine, 0)
    }
    sf.Lines = append(sf.Lines, &SourceLine{sf.IndentLevel, formatedLine})
}

func (sf *SimProgFile) Unindent() {
    if sf.IndentLevel > 0 {
        sf.IndentLevel -= 1
    }
}

func (sf *SimProgFile) addLineIndent(line string, args ...interface{}) {
    sf.addSourceLine(line, args...)
    sf.IndentLevel += 1
}

func (sf *SimProgFile) addLineUnindent(line string, args ...interface{}) {
    sf.Unindent()
    sf.addSourceLine(line, args...)
}

func (sf *SimProgFile) output(w io.Writer) {
    for _, line := range sf.Lines {
        // indentation
        for i := uint(0); i < line.Indent; i++ {
            w.Write([]byte(Indent))
        }
        // source line
        w.Write([]byte(line.Line))
        // newline
        w.Write([]byte("\n"))
    }
}


type Var interface {
    VarName() string {
    }
}

type SSAVar struct {
    name string
    version uint
}

func NewSSAVar(name string) {
    return &SSAVar{name, 0}
}

func (v *SSAVar) VarName() {
    return fmt.Sprintf("%s_%d", v.name, v.version)
}
