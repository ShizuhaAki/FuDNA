package common

// Hash related
const BASE = 131

type HashedString struct {
	s    string
	hash []uint64
}

func Hash(s string) HashedString {
	var ret []uint64
	ret = []uint64{0}

	for i := range len(s) {
		ret = append(ret, ret[len(ret)-1]*BASE+uint64(s[i]))
	}

	return HashedString{
		s:    s,
		hash: ret,
	}
}

func quickPow(a uint64, b uint64) uint64 {
	if b == 0 {
		return 1
	}
	ans := uint64(1)
	for ; b > 0; b >>= 1 {
		if b&1 == 1 {
			ans = ans * a
		}
		a = a * a
	}
	return a
}

func GetRangeHash(h *HashedString, l int, r int) uint64 {
	if l == 0 {
		return h.hash[r + 1]
	}
	return h.hash[r + 1] - h.hash[l]*quickPow(BASE, uint64(r-l+1))
}


func GetReverseComplement(s string) string {
	complement := map[byte]byte{
		'A': 'T', 'T': 'A',
		'C': 'G', 'G': 'C',
	}
	n := len(s)
	rc := make([]byte, n)
	for i := 0; i < n; i++ {
		base := s[n-1-i]
		if c, ok := complement[base]; ok {
			rc[i] = c
		} else {
			rc[i] = 'N' // handle unknown bases
		}
	}
	return string(rc)
}