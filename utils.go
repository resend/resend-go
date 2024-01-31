package resend

import "os"

// BytesToIntArray converts a byte array to rune array
// ie: []byte(`hello`) becomes []int{104,101,108,108,111}
// which will then be properly marshalled into JSON
// in the way Resend supports
func BytesToIntArray(a []byte) []int {
	res := make([]int, len(a))
	for i, v := range a {
		res[i] = int(v)
	}
	return res
}

func getEnv(key, df string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return df
}
