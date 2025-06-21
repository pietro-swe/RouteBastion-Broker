/*
Package cryptography provides common functionality related to cryptography
*/
package cryptography

type HashComparer interface {
	Compare(secretKey []byte, hash string, plain string) (bool, error)
}
