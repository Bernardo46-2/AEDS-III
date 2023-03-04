function postData() {
    const numero = document.getElementById("numero").value;
    const nome = document.getElementById("nome").value;
    const nomeJap = document.getElementById("nome_jap").value;
    const geracao = document.getElementById("geracao").value;
    const lancamento = new Date(document.getElementById("lancamento").value).toISOString();
    const especie = document.getElementById("especie").value;
    const lendario = document.getElementById("lendario").checked;
    const mitico = document.getElementById("mitico").checked;
    const tipo = document.getElementById("tipo").value.split(",");
    const atk = document.getElementById("atk").value;
    const def = document.getElementById("def").value;
    const hp = document.getElementById("hp").value;
    const altura = document.getElementById("altura").value;
    const peso = document.getElementById("peso").value;

    const pokemon = {
        numero: parseInt(numero),
        nome: nome,
        nome_jap: nomeJap,
        geracao: parseInt(geracao),
        lancamento: lancamento,
        especie: especie,
        lendario: lendario,
        mitico: mitico,
        tipo: tipo,
        atk: parseInt(atk),
        def: parseInt(def),
        hp: parseInt(hp),
        altura: parseFloat(altura),
        peso: parseFloat(peso)
    };

    const pikachu = {
        numero: 25,
        nome: "Pika Pika",
        nome_jap: "フシギダネ",
        geracao: 1,
        lancamento: "1996-02-27T08:00:00Z",
        especie: "Seed",
        lendario: false,
        mitico: false,
        tipo: ["Grass", "Poison"],
        atk: 49,
        def: 49,
        hp: 45,
        altura: 0.7,
        peso: 6.9
    };

    const url = "http://localhost:8080/put/";
    const options = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(pikachu)
    };

    console.log(options.body)

    fetch(url, options)
        .then(response => {
            if (response.ok) {
                alert("Pokemon modificado com sucesso!");
                window.location.href = "index.html";
            } else {
                alert("Erro ao adicionar pokemon");
            }
        })
        .catch(error => alert("Erro ao adicionar pokemon: " + error));
}