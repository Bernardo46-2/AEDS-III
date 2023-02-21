package main

import (
    "time"
    "fmt"
    "strconv"
    "unsafe"
)

const MAX_NAME_LEN = 40

type Pokemon struct {
    Numero     int32
    Nome       string
    NomeJap    string
    Geracao    int32
    Lancamento time.Time
    Especie    string
    Lendario   bool
    Mitico     bool
    Tipo       string
    Atk        int32
    Def        int32
    Hp         int32
    Altura     float64
    Peso       float64
}

type PokemonSize struct {
    Total      uintptr
    Numero     uintptr
    Nome       uintptr
    NomeJap    uintptr
    Geracao    uintptr
    Lancamento uintptr
    Especie    uintptr
    Lendario   uintptr
    Mitico     uintptr
    Tipo       uintptr
    Atk        uintptr
    Def        uintptr
    Hp         uintptr
    Altura     uintptr
    Peso       uintptr
}

var GenReleaseDates = map[int]string {
    1: "1996/02/27",
    2: "1999/11/21",
    3: "2002/11/21",
    4: "2006/09/28",
    5: "2010/09/18",
    6: "2013/10/12",
    7: "2016/11/18",
    8: "2019/11/15",
    9: "2022/11/18",
}

func Atoi32(s string) (int32, error) {
    i, err := strconv.Atoi(s)
    return int32(i), err
}

func (self* Pokemon) ToString() string {
    str := ""

    str += fmt.Sprintf("\n")
    str += fmt.Sprintf("Numero     = %d\n", self.Numero)
    str += fmt.Sprintf("Nome       = %s\n", self.Nome)
    str += fmt.Sprintf("NomeJap    = %s\n", self.NomeJap)
    str += fmt.Sprintf("Geracao    = %d\n", self.Geracao)
    str += fmt.Sprintf("Lancamento = %s\n", self.Lancamento.Format("02/01/2006"))
    str += fmt.Sprintf("Especie    = %s\n", self.Especie)
    str += fmt.Sprintf("Lendario   = %t\n", self.Lendario)
    str += fmt.Sprintf("Mitico     = %t\n", self.Mitico)
    str += fmt.Sprintf("Tipo       = %s\n", self.Tipo)
    str += fmt.Sprintf("Atk        = %d\n", self.Atk)
    str += fmt.Sprintf("Def        = %d\n", self.Def)
    str += fmt.Sprintf("Hp         = %d\n", self.Hp)
    str += fmt.Sprintf("Altura     = %f\n", self.Altura)
    str += fmt.Sprintf("Peso       = %f\n", self.Peso)

    return str
}

func ParsePokemon(line []string) (Pokemon, PokemonSize) {
    var pokemon Pokemon
    var size PokemonSize
    
    pokemon.Numero, _ = Atoi32(line[1])
    pokemon.Nome = line[2]
    pokemon.NomeJap = RemoveAfterSpace(line[4])
    geracao, _ := Atoi32(line[5])
    pokemon.Geracao = geracao
    pokemon.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[int(geracao)])
    pokemon.Especie = line[9]
    pokemon.Lendario, _ = strconv.ParseBool(line[7])
    pokemon.Mitico, _ = strconv.ParseBool(line[8])
    pokemon.Tipo = line[11] + line[12]
    pokemon.Atk, _ = Atoi32(line[21])
    pokemon.Def, _ = Atoi32(line[22])
    pokemon.Hp, _ = Atoi32(line[20])
    pokemon.Altura, _ = strconv.ParseFloat(line[13], 64)
    pokemon.Peso, _ = strconv.ParseFloat(line[14], 64)

    size.Numero = unsafe.Sizeof(pokemon.Numero)
    size.Nome = MAX_NAME_LEN
    size.NomeJap = (uintptr)(len(pokemon.NomeJap) * 4)
    size.Geracao = unsafe.Sizeof(pokemon.Geracao)

    date_size, err := pokemon.Lancamento.MarshalBinary()
    if err != nil {
        panic("Opora")
    }
    
    size.Lancamento = unsafe.Sizeof(len(date_size))
    size.Especie = (uintptr)(len(pokemon.Especie) * 4)
    size.Lendario = unsafe.Sizeof(pokemon.Lendario)
    size.Mitico = unsafe.Sizeof(pokemon.Mitico)
    size.Tipo = (uintptr)(len(pokemon.Tipo) * 4)
    size.Atk = unsafe.Sizeof(pokemon.Atk)
    size.Def = unsafe.Sizeof(pokemon.Def)
    size.Hp = unsafe.Sizeof(pokemon.Hp)
    size.Altura = unsafe.Sizeof(pokemon.Altura)
    size.Peso = unsafe.Sizeof(pokemon.Peso)

    size.Total = size.Numero + 
                 size.Nome + 
                 size.NomeJap + 
                 size.Geracao + 
                 size.Lancamento + 
                 size.Especie + 
                 size.Lendario + 
                 size.Mitico + 
                 size.Tipo + 
                 size.Atk + 
                 size.Def + 
                 size.Hp + 
                 size.Altura + 
                 size.Peso

    return pokemon, size
}
