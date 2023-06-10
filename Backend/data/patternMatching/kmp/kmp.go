package kmp

import (
	"strings"

	"github.com/Bernardo46-2/AEDS-III/data/binManager"
)

const (
	PatternSize int = 100
)

func SearchNext(haystack string, needle string) int {
	retSlice := kmp(haystack, needle)
	if len(retSlice) > 0 {
		return retSlice[len(retSlice)-1]
	}

	return -1
}

func SearchString(haystack string, needle string) int {
	retSlice := kmp(strings.ToLower(haystack), strings.ToLower(needle))
	if len(retSlice) > 0 {
		return retSlice[0]
	}

	return -1
}

func kmp(haystack string, needle string) []int {
	next := preKMP(needle)
	i := 0
	j := 0
	m := len(needle)
	n := len(haystack)

	x := []byte(needle)
	y := []byte(haystack)
	var ret []int

	//got zero target or want, just return empty result
	if m == 0 || n == 0 {
		return ret
	}

	//want string bigger than target string
	if n < m {
		return ret
	}

	for j < n {
		for i > -1 && x[i] != y[j] {
			i = next[i]
		}
		i++
		j++

		//fmt.Println(i, j)
		if i >= m {
			ret = append(ret, j-i)
			//fmt.Println("find:", j, i)
			i = next[i]
		}
	}

	return ret
}

func preKMP(x string) [PatternSize]int {
	var i, j int
	length := len(x) - 1
	var kmpNext [PatternSize]int
	i = 0
	j = -1
	kmpNext[0] = -1

	for i < length {
		for j > -1 && x[i] != x[j] {
			j = kmpNext[j]
		}

		i++
		j++

		if x[i] == x[j] {
			kmpNext[i] = kmpNext[j]
		} else {
			kmpNext[i] = j
		}
	}
	return kmpNext
}

func SearchPokemon(search string, field string) []int64 {
	controller, _ := binManager.InicializarControleLeitura(binManager.BIN_FILE)
	target := []int64{}

	for err := controller.ReadNext(); err == nil; err = controller.ReadNext() {
		if !controller.RegistroAtual.IsDead() {
			needle := SearchString(controller.RegistroAtual.Pokemon.GetField(field), search)
			if needle != -1 {
				target = append(target, int64(controller.RegistroAtual.Pokemon.Numero))
			}
		}
	}

	return target
}
