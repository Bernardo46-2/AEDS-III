package main

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
    "encoding/binary"
    "math"
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

func runCmd(name string, arg ...string) {
    cmd := exec.Command(name, arg...)
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func clearScreen() {
    switch runtime.GOOS {
    case "darwin":
        runCmd("clear")
    case "linux":
        runCmd("clear")
    case "windows":
        runCmd("cmd", "/c", "cls")
    default:
        runCmd("clear")
    }
}

func pause() {
    var input string
    fmt.Printf("\nPressione Enter para continuar...\n")
    fmt.Scanf("%s\n", &input)
}

func lerInt(prompt string) int {
    clearScreen()
    fmt.Printf("%s\n> ", prompt)

    var tmp string
    var result int
    var err error

    if _, err = fmt.Scanln(&tmp); err != nil {
        fmt.Println("\nErro ao ler opção:", err)
        pause()
        result = lerInt(prompt)
    } else {
        if result, err = strconv.Atoi(tmp); err != nil {
            fmt.Println("\nErro ao ler opção:", err)
            pause()
            result = lerInt(prompt)
        }
    }
    return result
}

func lerInt32(prompt string, backup string, limite ...int) (int32, string) {
    clearScreen()
    fmt.Printf(backup)
    fmt.Printf(prompt + ": ")

    var input string
    var result int
    var err error

    if _, err = fmt.Scanln(&input); err != nil {
        fmt.Println("\nErro ao ler opção:", err)
        pause()
        return lerInt32(prompt, backup, limite...)
    } else {
        if result, err = strconv.Atoi(input); err != nil {
            fmt.Println("\nErro ao ler opção:", err)
            pause()
            return lerInt32(prompt, backup)
        }
        if len(limite) > 0 && result > limite[0] {
            fmt.Printf("\nValor fora do range permitido [%d]", limite[0])
            pause()
            return lerInt32(prompt, backup, limite[0])
        }
    }

    return int32(result), backup + prompt + ": " + input + "\n"
}

func lerFloat32(prompt string, backup string) (float32, string) {
    clearScreen()
    fmt.Printf(backup)
    fmt.Printf(prompt + ": ")

    var input string
    var result float64
    var err error

    if _, err = fmt.Scanln(&input); err != nil {
        fmt.Println("\nErro ao ler opção:", err)
        pause()
        return lerFloat32(prompt, backup)
    } else {
        if result, err = strconv.ParseFloat(input, 32); err != nil {
            fmt.Println("\nErro ao ler opção:", err)
            pause()
            return lerFloat32(prompt, backup)
        }
    }

    return float32(result), backup + prompt + ": " + input + "\n"
}

func parseBool(str string) (bool, error) {
    switch str {
    case "1", "t", "T", "true", "TRUE", "True", "s", "S", "sim", "Sim", "SIM":
        return true, nil
    case "0", "f", "F", "false", "FALSE", "False", "n", "N", "sao", "Nao", "NAO":
        return false, nil
    }
    return false, fmt.Errorf("Erro de sintaxe")
}

// lerBool é uma função que lê um valor booleano da entrada do usuário.
// Ela recebe um parâmetro `prompt` como string, que será exibido para o usuário
// antes da entrada do valor booleano. Também recebe uma string `backup`,
// que é usada para limpar a tela e restaurar uma mensagem anterior após a entrada.
// A função retorna o valor booleano lido e uma string contendo a mensagem de prompt e valor lido.
func lerBool(prompt string, backup string) (bool, string) {
    // Limpa a tela e exibe a mensagem de backup
    clearScreen()
    fmt.Printf(backup)
    fmt.Printf(prompt + ": ")

    var input string
    var result bool
    var err error

    // Lê a entrada do usuário
    if _, err = fmt.Scanln(&input); err != nil {
        // Em caso de erro, exibe mensagem de erro e espera que o usuário tente novamente
        fmt.Println("\nErro ao ler opção:", err)
        pause()
        return lerBool(prompt, backup)
    } else {
        // Converte a entrada para um valor booleano
        if result, err = parseBool(input); err != nil {
            // Em caso de erro de conversão, exibe mensagem de erro e espera que o usuário tente novamente
            fmt.Println("\nErro ao ler opção:", err)
            pause()
            return lerBool(prompt, backup)
        }
    }

    // Retorna o valor booleano lido e uma mensagem com o prompt e valor lido
    return result, backup + prompt + ": " + input + "\n"
}

func lerString(prompt string, backup string) (string, string) {
    clearScreen()
    fmt.Printf(backup)
    fmt.Printf(prompt + ": ")

    reader := bufio.NewReader(os.Stdin)
    var input string
    var err error

    if input, err = reader.ReadString('\n'); err != nil {
        fmt.Println("\nErro ao ler opção:", err)
        pause()
        return lerString(prompt, backup)
    }

    return strings.TrimSpace(input), backup + prompt + ": " + input
}

// lerStringSlice lê um slice de strings do usuário, permitindo escolher entre
// guardar até 1 ou 2 strings. Retorna o slice lido e uma mensagem com o que
// foi lido.
func lerStringSlice(prompt string, backup string, maxLen int) ([]string, string) {
    clearScreen()
    fmt.Printf(backup)
    fmt.Printf(prompt)

    var slice []string
    var s string

    numTipos, backup := lerInt32("Numero de "+prompt, backup, maxLen)

    for i := int32(0); i < numTipos; i++ {
        s, backup = lerString(prompt+"["+strconv.Itoa(len(slice)+1)+"]", backup)
        slice = append(slice, s)
    }

    return slice, backup
}
