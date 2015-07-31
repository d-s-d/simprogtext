package simprogtext

import (
    "fmt"
    "io"
)

var Indent = "    "

type SimProgFile interface {
    AddLine(string, ...interface{})
    Unindent()
    Indent()
    AddLineIndent(string, ...interface{})
    AddLineUnindent(string, ...interface{})
    WriteToFile() error
}

type SourceLine struct {
    indent uint
    line string
}

type BufferedSimProgFile struct {
    indentLevel uint
    lines []*SourceLine
    f_out io.Writer
}

func NewBufferedSimProgFile(f_out io.Writer) SimProgFile {
    return &BufferedSimProgFile{0, make([]*SourceLine, 0), f_out}
}

func (sf *BufferedSimProgFile) AddLine(line string, args ...interface{}) {
    formatedLine := fmt.Sprintf(line, args...)
    if sf.lines == nil {
        sf.lines = make([]*SourceLine, 0)
    }
    sf.lines = append(sf.lines, &SourceLine{sf.indentLevel, formatedLine})
}

func (sf *BufferedSimProgFile) Indent() {
    sf.indentLevel += 1
}

func (sf *BufferedSimProgFile) Unindent() {
    if sf.indentLevel > 0 {
        sf.indentLevel -= 1
    }
}

func (sf *BufferedSimProgFile) AddLineIndent(line string,
args ...interface{}) {
    sf.AddLine(line, args...)
    sf.indentLevel += 1
}

func (sf *BufferedSimProgFile) AddLineUnindent(line string,
args ...interface{}) {
    sf.Unindent()
    sf.AddLine(line, args...)
}

func (sf *BufferedSimProgFile) WriteToFile() error {
    var err error
    w := sf.f_out
    for _, line := range sf.lines {
        // indentation
        for i := uint(0); i < line.indent; i++ {
            _, err = w.Write([]byte(Indent))
            if err != nil { return err }
        }
        // source line
        _, err = w.Write([]byte(line.line))
        if err != nil { return err }
        // newline
        _, err = w.Write([]byte("\n"))
        if err != nil { return err }
    }
    return nil
}

type Var interface {
    VarName() string
}

type SSAVar interface {
    Var
    Next() string
}

type DynSSAVar interface {
    SSAVar
    GetType() string
    SetType(t string)
    NextType(string) string
}

type SimpleVar struct {
    name string
}

func (v *SimpleVar) VarName() string {
    return v.name
}

type DynSSAv struct {
    name string
    version uint
    typeName string
}

func NewSimpleVar(name string) Var {
    return &SimpleVar{name}
}

func NewDynSSAVar(name string, typeName string) DynSSAVar {
    return &DynSSAv{name, 0, ""}
}


func (v *DynSSAv) VarName() string {
    return fmt.Sprintf("%s_%d", v.name, v.version)
}

func (v *DynSSAv) Next() string {
    v.version += 1
    return v.VarName()
}

func (v *DynSSAv) GetType() string {
    return v.typeName
}

func (v *DynSSAv) SetType(n string) {
    v.typeName = n
}

func (v *DynSSAv) NextType(n string) string {
    v.typeName = n
    return v.Next()
}
