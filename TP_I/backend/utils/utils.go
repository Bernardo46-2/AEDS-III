package utils

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Atoi32(s string) (int32, error) {
	i, err := strconv.Atoi(s)
	return int32(i), err
}

// Funcao para converter int32 para []byte
func IntToBytes(n int32) []byte {
	var buf []byte
	return binary.LittleEndian.AppendUint32(buf, uint32(n))
}

// Funcao para converter float32 para []byte
func FloatToBytes(f float32) []byte {
	b := make([]byte, 4)
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(b, bits)

	return b
}

func RemoveAfterSpace(str string) string {
	parts := strings.Split(str, " ")
	return parts[0]
}

func RunCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearScreen() {
	switch runtime.GOOS {
	case "darwin":
		RunCmd("clear")
	case "linux":
		RunCmd("clear")
	case "windows":
		RunCmd("cmd", "/c", "cls")
	default:
		RunCmd("clear")
	}
}

func Pause() {
	var input string
	fmt.Printf("\nPressione Enter para continuar...\n")
	fmt.Scanf("%s\n", &input)
}

func LerInt(prompt string) int {
	ClearScreen()
	fmt.Printf("%s\n> ", prompt)

	var tmp string
	var result int
	var err error

	if _, err = fmt.Scanln(&tmp); err != nil {
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		result = LerInt(prompt)
	} else {
		if result, err = strconv.Atoi(tmp); err != nil {
			fmt.Println("\nErro ao Ler opção:", err)
			Pause()
			result = LerInt(prompt)
		}
	}
	return result
}

func LerInt32(prompt string, backup string, limite ...int) (int32, string) {
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s: ", prompt)

	var input string
	var result int
	var err error

	if _, err = fmt.Scanln(&input); err != nil {
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		return LerInt32(prompt, backup, limite...)
	} else {
		if result, err = strconv.Atoi(input); err != nil {
			fmt.Println("\nErro ao Ler opção:", err)
			Pause()
			return LerInt32(prompt, backup)
		}
		if len(limite) > 0 && result > limite[0] {
			fmt.Printf("\nValor fora do range permitido [%d]", limite[0])
			Pause()
			return LerInt32(prompt, backup, limite[0])
		}
	}

	return int32(result), backup + prompt + ": " + input + "\n"
}

func LerFloat32(prompt string, backup string) (float32, string) {
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s: ", prompt)

	var input string
	var result float64
	var err error

	if _, err = fmt.Scanln(&input); err != nil {
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		return LerFloat32(prompt, backup)
	} else {
		if result, err = strconv.ParseFloat(input, 32); err != nil {
			fmt.Println("\nErro ao Ler opção:", err)
			Pause()
			return LerFloat32(prompt, backup)
		}
	}

	return float32(result), backup + prompt + ": " + input + "\n"
}

func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "s", "S", "sim", "Sim", "SIM", "y", "Y", "yes", "Yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "n", "N", "nao", "Nao", "NAO", "não", "Não", "NÃO", "no", "No", "NO":
		return false, nil
	}
	return false, fmt.Errorf("erro de sintaxe")
}

// LerBool é uma função que lê um valor booleano da entrada do usuário.
// Ela recebe um parâmetro `prompt` como string, que será exibido para o usuário
// antes da entrada do valor booleano. Também recebe uma string `backup`,
// que é usada para limpar a tela e restaurar uma mensagem anterior após a entrada.
// A função retorna o valor booleano lido e uma string contendo a mensagem de prompt e valor lido.
func LerBool(prompt string, backup string) (bool, string) {
	// Limpa a tela e exibe a mensagem de backup
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s: ", prompt)

	var input string
	var result bool
	var err error

	// Lê a entrada do usuário
	if _, err = fmt.Scanln(&input); err != nil {
		// Em caso de erro, exibe mensagem de erro e espera que o usuário tente novamente
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		return LerBool(prompt, backup)
	} else {
		// Converte a entrada para um valor booleano
		if result, err = ParseBool(input); err != nil {
			// Em caso de erro de conversão, exibe mensagem de erro e espera que o usuário tente novamente
			fmt.Println("\nErro ao Ler opção:", err)
			Pause()
			return LerBool(prompt, backup)
		}
	}

	// Retorna o valor booleano lido e uma mensagem com o prompt e valor lido
	return result, backup + prompt + ": " + input + "\n"
}

func LerTime(prompt string, backup string) (time.Time, string) {
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s (dd/mm/aaaa): ", prompt)

	var input string
	var t time.Time
	var err error

	if _, err = fmt.Scanln(&input); err != nil {
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		return LerTime(prompt, backup)
	} else {
		if t, err = time.Parse("02/01/2006", input); err != nil {
			fmt.Println("\nErro ao Ler opção:", err)
			Pause()
			return LerTime(prompt, backup)
		}
	}

	return t, backup + prompt + ": " + input + "\n"
}

func LerString(prompt string, backup string) (string, string) {
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	var input string
	var err error

	if input, err = reader.ReadString('\n'); err != nil {
		fmt.Println("\nErro ao Ler opção:", err)
		Pause()
		return LerString(prompt, backup)
	}

	return strings.TrimSpace(input), backup + prompt + ": " + input
}

// LerStringSlice lê um slice de strings do usuário, permitindo escolher entre
// guardar até 1 ou 2 strings. Retorna o slice lido e uma mensagem com o que
// foi lido.
func LerStringSlice(prompt string, backup string, maxLen int) ([]string, string) {
	ClearScreen()
	fmt.Printf("%s", backup)
	fmt.Printf("%s", prompt)

	var slice []string
	var s string

	numTipos, backup := LerInt32("Numero de "+prompt, backup, maxLen)

	for i := int32(0); i < numTipos; i++ {
		s, backup = LerString(prompt+"["+strconv.Itoa(len(slice)+1)+"]", backup)
		slice = append(slice, s)
	}

	return slice, backup
}

func BytesToVarSize(registro []byte, ptr int) (int, int) {
	return int(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

func BytesToInt32(registro []byte, ptr int) (int32, int) {
	return int32(binary.LittleEndian.Uint32(registro[ptr : ptr+4])), ptr + 4
}

func BytesToString(registro []byte, ptr int) (string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	nomeBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), nomeBytes)
	return strings.TrimSpace(string(nomeBytes)), ptr + size
}

func BytesToFixedSizeString(registro []byte, ptr int, maxSize int) (string, int) {
	nome := make([]byte, maxSize)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+maxSize]), nome)
	return strings.TrimRight(string(nome), "\x00"), ptr + maxSize
}

func BytesToArrayString(registro []byte, ptr int) ([]string, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	stringBytes := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), stringBytes)
	s := strings.TrimRight(string(stringBytes), ",")
	return strings.Split(s, ","), ptr + size
}

func BytesToJapName(registro []byte, ptr int) (string, int) {
	size, ptr := BytesToVarSize(registro, ptr)

	japNameRunes := make([]rune, size/4)
	for i := 0; i < size/4; i++ {
		// Converte os 4 bytes em um uint32 correspondente ao rune.
		runeUint := binary.LittleEndian.Uint32(registro[ptr : ptr+4])
		// Converte o uint32 em um rune e adiciona à slice de runes.
		japNameRunes[i] = rune(runeUint)
		ptr += 4
	}

	return string(japNameRunes), ptr
}

func BytesToTime(registro []byte, ptr int) (time.Time, int) {
	size, ptr := BytesToVarSize(registro, ptr)
	b := make([]byte, size)
	io.ReadFull(bytes.NewReader(registro[ptr:ptr+size]), b)
	var t time.Time
	t.UnmarshalBinary(b)
	return t, ptr + size
}

func BytesToBool(registro []byte, ptr int) (bool, int) {
	if registro[ptr] != 0 {
		return true, ptr + 1
	} else {
		return false, ptr + 1
	}
}

func BytesToFloat32(registro []byte, ptr int) (float32, int) {
	size := 4
	bits := binary.LittleEndian.Uint32(registro[ptr : ptr+size])
	float := math.Float32frombits(bits)
	return float, ptr + size
}
