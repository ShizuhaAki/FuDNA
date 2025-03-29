package lab1

import (
	"bufio"
	"fmt"
	"os"
)

const base = 131

type Node struct {
	left, right int  // 1-indexed positions in seq1 (or seq2 for second move)
	dir         byte // '+' or '-' indicating the direction
}

func computePowers(n int) []uint64 {
	powers := make([]uint64, n+1)
	powers[0] = 1
	for i := 1; i <= n; i++ {
		powers[i] = powers[i-1] * base
	}
	return powers
}

func computePrefixHash(s string, powers []uint64) []uint64 {
	n := len(s)
	h := make([]uint64, n+1)
	for i := 1; i <= n; i++ {
		h[i] = h[i-1]*base + uint64(s[i-1])
	}
	return h
}

func substringHash(prefix []uint64, powers []uint64, l, r int) uint64 {
	return prefix[r] - prefix[l-1]*powers[r-l+1]
}

func reverseComplement(s string) string {
	n := len(s)
	rc := make([]byte, n)
	for i := 0; i < n; i++ {
		var comp byte
		switch s[i] {
		case 'A':
			comp = 'T'
		case 'T':
			comp = 'A'
		case 'C':
			comp = 'G'
		case 'G':
			comp = 'C'
		default:
			comp = s[i]
		}
		rc[n-1-i] = comp
	}
	return string(rc)
}

func Main() {
	// Read input.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	seq1 := scanner.Text() // First DNA string
	scanner.Scan()
	seq2 := scanner.Text() // Second DNA string

	len1, len2 := len(seq1), len(seq2)

	// Compute the reverse complement of seq2.
	revComp := reverseComplement(seq2)

	// Compute powers up to the maximum needed length.
	maxLen := len1
	if len2 > maxLen {
		maxLen = len2
	}
	powers := computePowers(maxLen)

	// Compute prefix hash arrays (using 1-indexing).
	hashSeq1 := computePrefixHash(seq1, powers)
	hashSeq2 := computePrefixHash(seq2, powers)
	hashRevComp := computePrefixHash(revComp, powers)

	// Build a map from substring hash (from seq1) to list of ending positions.
	entriesByHash := make(map[uint64][]int)
	for i := 1; i <= len1; i++ {
		for j := i; j <= len1; j++ {
			h := substringHash(hashSeq1, powers, i, j)
			entriesByHash[h] = append(entriesByHash[h], j)
		}
	}

	// Initialize DP and transition tables.
	dp := make([][]bool, len2+1)
	transition := make([][]Node, len2+1)
	for i := 0; i <= len2; i++ {
		dp[i] = make([]bool, len1+1)
		transition[i] = make([]Node, len1+1)
	}
	dp[0][0] = true

	// Fill in DP table.
	// Simple character-by-character matching along the diagonal.
	for i := 1; i <= len2; i++ {
		for j := 1; j <= len1; j++ {
			if seq2[i-1] == seq1[j-1] && dp[i-1][j-1] {
				dp[i][j] = true
			}
		}
		// Try to match segments using substring hashes.
		for j := i; j >= 1; j-- {
			segLen := i - j + 1

			// Case 1: Use the segment from seq2.
			currHash := substringHash(hashSeq2, powers, j, i)
			if posList, exists := entriesByHash[currHash]; exists {
				for _, pos := range posList {
					if pos <= len1 && dp[i-segLen][pos] {
						dp[i][pos] = true
						transition[i][pos] = Node{
							left:  pos - segLen + 1,
							right: pos,
							dir:   '+',
						}
					}
				}
			}

			// Case 2: Use the segment from the reverse complement.
			currHash = substringHash(hashRevComp, powers, len2-i+1, len2-j+1)
			if posList, exists := entriesByHash[currHash]; exists {
				for _, pos := range posList {
					if pos <= len1 && dp[i-segLen][pos] {
						dp[i][pos] = true
						transition[i][pos] = Node{
							left:  pos - segLen + 1,
							right: pos,
							dir:   '-',
						}
					}
				}
			}
		}
	}

	// Reconstruct the sequence of operations.
	row, col := len2, len1
	var moves []Node
	for row != 0 {
		if transition[row][col].left == 0 {
			row--
			col--
			continue
		}
		moves = append(moves, transition[row][col])
		segLen := transition[row][col].right - transition[row][col].left + 1
		moves = append(moves, Node{
			left:  row - segLen + 1,
			right: row,
			dir:   transition[row][col].dir,
		})
		row -= segLen
	}

	writer := bufio.NewWriter(os.Stdout)
	for i, j := len(moves)-1, 0; i >= 0; i -= 2 {
		fmt.Fprintf(writer, "Match %d\n", j)
		j++
		fmt.Fprintf(writer, "query    : l = %d, r = %d, direction = %c\n", moves[i].left, moves[i].right, moves[i].dir)
		fmt.Fprintf(writer, "reference: l = %d, r = %d, direction = %c\n", moves[i-1].left, moves[i-1].right, moves[i-1].dir)
	}
	writer.Flush()
}
