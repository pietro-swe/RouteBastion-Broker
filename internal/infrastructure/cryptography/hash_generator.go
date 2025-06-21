package cryptography

type HashGenerator interface {
	Generate(secretKey []byte, plain string) (string, error)
}
