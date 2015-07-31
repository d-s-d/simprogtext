package simprogtext

import (
    "fmt"
    "io"
)

// Indent is the string that is used as indentation
var Indent = "    "

// SimProgFile is an interface representing a simple program text file.
type SimProgFile interface {
    AddLine(string, ...interface{})
    Unindent()
    Indent()
    AddLineIndent(string, ...interface{})
    AddLineUnindent(string, ...interface{})
    WriteToFile() error
}

// SourceLine is a structure representing one source code line.
type SourceLine struct {
    indent uint
    line string
}

// BufferedSimProgFile represents a buffer for a simple program text file.
type BufferedSimProgFile struct {
    indentLevel uint
    lines []*SourceLine
    fOut io.Writer
}

// NewBufferedSimProgFile create and returns a buffered simple program text
// file.
func NewBufferedSimProgFile(fOut io.Writer) SimProgFile {
    return &BufferedSimProgFile{0, make([]*SourceLine, 0), fOut}
}

// AddLine appends a line to the program text file.
func (sf *BufferedSimProgFile) AddLine(line string, args ...interface{}) {
    formatedLine := fmt.Sprintf(line, args...)
    if sf.lines == nil {
        sf.lines = make([]*SourceLine, 0)
    }
    sf.lines = append(sf.lines, &SourceLine{sf.indentLevel, formatedLine})
}

// Indent increments the current indentation level by 1.
func (sf *BufferedSimProgFile) Indent() {
    sf.indentLevel++
}

// Unindent decrements the current indentation level by 1.
func (sf *BufferedSimProgFile) Unindent() {
    if sf.indentLevel > 0 {
        sf.indentLevel--
    }
}

// AddLineIndent adds the provieded line line and increments the current
// indentation level by 1.
func (sf *BufferedSimProgFile) AddLineIndent(line string,
args ...interface{}) {
    sf.AddLine(line, args...)
    sf.indentLevel++
}

// AddLineUnindent first decrements the current indentation by 1 and
// subsequently adds the provided line.
func (sf *BufferedSimProgFile) AddLineUnindent(line string,
args ...interface{}) {
    sf.Unindent()
    sf.AddLine(line, args...)
}

// WriteToFile writes the content of the buffer to the io.Writer provided
// when the SimProgFile was created.
func (sf *BufferedSimProgFile) WriteToFile() error {
    var err error
    w := sf.fOut
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

// Var is a general interface for an object representing a variable
type Var interface {
    VarName() string
}

// SSAVar is an interface representing an single static assignment variable
type SSAVar interface {
    Var
    Next() string
}

// DynSSAVar is an interface representing a dynamically typed single static
// assignment variable.
type DynSSAVar interface {
    SSAVar
    GetType() string
    SetType(t string)
    NextType(string) string
}

// SimpleVar ...
type SimpleVar struct {
    name string
}

// VarName returns the variable name of the variable.
func (v *SimpleVar) VarName() string {
    return v.name
}

// DynSSAv is the structure that represents a dynamic static single assignment
// variable.
type DynSSAv struct {
    name string
    version uint
    typeName string
}

// NewSimpleVar ...
func NewSimpleVar(name string) Var {
    return &SimpleVar{name}
}

// NewDynSSAVar ...
func NewDynSSAVar(name string, typeName string) DynSSAVar {
    return &DynSSAv{name, 0, ""}
}

// VarName ...
func (v *DynSSAv) VarName() string {
    return fmt.Sprintf("%s_%d", v.name, v.version)
}

// Next increments the version of the variable
func (v *DynSSAv) Next() string {
    v.version++
    return v.VarName()
}

// GetType returns the type of the variable.
func (v *DynSSAv) GetType() string {
    return v.typeName
}

// SetType sets the type of the variable.
func (v *DynSSAv) SetType(n string) {
    v.typeName = n
}

// NextType sets the type of the variable and increases the version count by 1
func (v *DynSSAv) NextType(n string) string {
    v.typeName = n
    return v.Next()
}
