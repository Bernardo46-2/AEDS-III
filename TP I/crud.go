package main

import (
    "fmt"
)

func create() {
    pokemon := readPokemon()
    fmt.Printf("%s", pokemon.ToString())
    AlterarNumRegistros(1)
    pokeBytes := pokemon.ToBytes()
    
    if err := AppendPokemon(pokeBytes); err != nil {
        fmt.Println("\n", err)
    } else {
        fmt.Println("Pokemon inserido com sucesso")
    }
}

func read() {
    pokemon, err, _ := readBinToPoke(lerInt("Digite o numero da pokedex a pesquisar:\n"))
    if err != nil {
        fmt.Printf("\n%s\n", err)
    } else {
        fmt.Printf(pokemon.ToString())
    }
}

func update() {
    fmt.Printf("Update\n")
}

func delete() {
    pokemon, err, pos := readBinToPoke(lerInt("Digite o numero da pokedex a deletar:\n"))
    if err != nil {
        fmt.Printf("Erro ao excluir\n%s\n", err)
    } else {
        if err = deletarPokemon(pos); err != nil {
            fmt.Printf("Erro ao excluir\n%s\n", err)
        }
        AlterarNumRegistros(-1)
        fmt.Printf("Pokemon %s excluido com exito", pokemon.Nome)
    }
}