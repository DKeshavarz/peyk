package codegen

type CodeGenerator interface {
	Generate() (string, error)
}


