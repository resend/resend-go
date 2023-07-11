package resend

import "strconv"

// ByteArrayToStringArray converts a byte array to string array
// ie: []byte{44,45,46} becomes []string{44,45,46}
// which will then be properly marshalled into JSON
// in the way Resend supports
func ByteArrayToStringArray(a []byte) []string {
	res := []string{}
	for _, v := range a {
		v := strconv.Itoa(int(v))
		res = append(res, v)
	}
	return res
}
