package main

import (
    "time"
    "fmt"
    "strconv"
)

type Pokemon struct {
    Numero     int
    Nome       string
    NomeJap    string
    Geracao    int
    Lancamento time.Time
    Especie    string
    Lendario   bool
    Mitico     bool
    Tipo       string
    Atk        int
    Def        int
    Hp         int
    Altura     float64
    Peso       float64
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

func ParsePokemon(line []string) Pokemon{
    var pokemon Pokemon
    
    pokemon.Numero, _ = strconv.Atoi(line[1])
    pokemon.Nome = line[2]
    pokemon.NomeJap = RemoveAfterSpace(line[4])
    geracao, _ := strconv.Atoi(line[5])
    pokemon.Geracao = geracao
    pokemon.Lancamento, _ = time.Parse("2006/01/02", GenReleaseDates[geracao])
    pokemon.Especie = line[9]
    pokemon.Lendario, _ = strconv.ParseBool(line[7])
    pokemon.Mitico, _ = strconv.ParseBool(line[8])
    pokemon.Tipo = line[11] + line[12]
    pokemon.Atk, _ = strconv.Atoi(line[21])
    pokemon.Def, _ = strconv.Atoi(line[22])
    pokemon.Hp, _ = strconv.Atoi(line[20])
    pokemon.Altura, _ = strconv.ParseFloat(line[13], 64)
    pokemon.Peso, _ = strconv.ParseFloat(line[14], 64)

    return pokemon
}
