package llvm

import ("bytes" ; "fmt" ; "strings")

type FString struct {
	Name    string
	Content string
	Length  int
}

var fstringCount int
func init() { fstringCount = 1 }

func NewFStringLabel() string {
	label := fmt.Sprintf("fstr%d", fstringCount)
	fstringCount++
	return label
}

func NewFString(content string) *FString {
	content = strings.ReplaceAll(content, "%d", "%ld")
	content = strings.ReplaceAll(content, "\\n", "\\0A")
	content = content + "\\00"
	length := 1

	for i := 0; i < len(content); i++ {
		if content[i] == '\\' {
			if (i + 1 < len(content) && content[i + 1] == '0') { continue }
		} else if content[i] == '0' {
			if i + 1 < len(content) && (content[i + 1] == '0' || content[i + 1] == 'A') { continue }
		} else { length++ }
	}

	return &FString{NewFStringLabel(), content, length}
}

func (f *FString) String() string {
	var out bytes.Buffer
	out.WriteString("@." + f.Name)
	out.WriteString(" ")
	out.WriteString("=")
	out.WriteString(" ")
	out.WriteString("private")
	out.WriteString(" ")
	out.WriteString("unnamed_addr")
	out.WriteString(" ")
	out.WriteString("constant")
	out.WriteString(" ")
	out.WriteString(fmt.Sprintf("[%d x i8]", f.Length))
	out.WriteString("c")
	out.WriteString("\"")
	out.WriteString(f.Content)
	out.WriteString("\"")
	out.WriteString(",")
	out.WriteString("align")
	out.WriteString(" ")
	out.WriteString("1")
	out.WriteString("\n")
	return out.String()
}
