package resend

import "strconv"

// ByteArrayToString converts a byte array to string
// ie: []byte{44,45,46} becomes "[44,45,46]"
func ByteArrayToStringArray(a []byte) []string {
	res := []string{}
	for _, v := range a {
		v := strconv.Itoa(int(v))
		res = append(res, v)
	}
	return res
}
