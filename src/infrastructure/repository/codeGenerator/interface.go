package codeGeneratorRepository

type CodeGeneratorRepository interface {
	GenerateCode(key string) (string, error)
	VerifyCode(key, code string) bool
}
