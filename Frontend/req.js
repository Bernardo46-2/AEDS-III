/*  * * * * * * * * * * * * *
 *
 *  Import CSV
 *
 * * * * * * * * * * * * * * */

const importarCsv = document.getElementById('ImportarCSV');
const modalContainer = document.getElementById('modal-container');
const deleteBtn = document.getElementById('delete');
const helpBtn = document.getElementById('Ajuda');

helpBtn.onclick = () => {
    modalAviso("Por favor me dê um emprego (╥﹏╥)");
}

function modalAviso(mostrar = "Servidor Desligado") {
    const mensagem = document.getElementById("mensagem-modal");
    mensagem.innerHTML = mostrar;

    modalContainer.classList.remove('out');
    modalContainer.classList.add('one');

    setTimeout(function () {
        modalContainer.classList.add('out');
    }, 3000);

    setTimeout(function () {
        modalContainer.addEventListener('click', () => {
            modalContainer.classList.add('out');
        });
    }, 1200);
}

importarCsv.onclick = () => {
    fetch('http://localhost:8080/loadDatabase')
        .then(response => response.json())
        .then(data => {
            modalAviso(data.mensagem);
            showAll.onclick();
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

const IntercalacaoComum = document.getElementById('IntercalacaoComum');
const IntercalacaoVariavel = document.getElementById('IntercalacaoVariavel');
const SelecaoPorSubstituicao = document.getElementById('SelecaoPorSubstituicao');

IntercalacaoComum.onclick = () => {
    fetch('http://localhost:8080/intercalacaoComum/')
        .then(response => response.json())
        .then(data => {
            modalAviso(data.mensagem);
            showAll.onclick();
        })
        .catch(error => {
            modalAviso();
            console.log(error);
        });
}

IntercalacaoVariavel.onclick = () => {
    fetch('http://localhost:8080/intercalacaoVariavel/')
        .then(response => response.json())
        .then(data => {
            modalAviso(data.mensagem);
            showAll.onclick();
        })
}

SelecaoPorSubstituicao.onclick = () => {
    fetch('http://localhost:8080/selecaoPorSubstituicao/')
        .then(response => response.json())
        .then(data => {
            modalAviso(data.mensagem);
            showAll.onclick();
        })
}

/*  * * * * * * * * * * * * *
 *
 *  Mostrar Todos
 *
 * * * * * * * * * * * * * * */

const showAll = document.getElementById('All');

showAll.onclick = () => {
    window.scrollTo(0, 0);
    fetch('http://localhost:8080/getAll/?page=0')
        .then(response => response.json())
        .then(data => adicionarCards(data))
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

window.onload = function () {
    fetch('http://localhost:8080/getAll/?page=0')
        .then(response => response.json())
        .then(data => adicionarCards(data))
        .catch(error => {
            modalAviso();
            console.log(error)
        });
};

const classes = {
    'Normal': 'bgd-Normal',
    'Fire': 'bgd-Fire',
    'Water': 'bgd-Water',
    'Electric': 'bgd-Electric',
    'Grass': 'bgd-Grass',
    'Ice': 'bgd-Ice',
    'Fighting': 'bgd-Fighting',
    'Poison': 'bgd-Poison',
    'Ground': 'bgd-Ground',
    'Flying': 'bgd-Flying',
    'Psychic': 'bgd-Psychic',
    'Bug': 'bgd-Bug',
    'Rock': 'bgd-Rock',
    'Ghost': 'bgd-Ghost',
    'Dragon': 'bgd-Dragon',
    'Dark': 'bgd-Dark',
    'Steel': 'bgd-Steel',
    'Fairy': 'bgd-Fairy'
};

let lastClicked = 1;
let insertDots = true;

function adicionarCards(data) {
    // console.log(data);
    const cardsHtml = document.getElementById('cards');
    cardsHtml.innerHTML = "";
    for (let i = 0; i < data.length; i++) {
        let nome = data[i].nome.toLowerCase();

        let bgd = classes && classes[data[i].tipo[0]] ? classes[data[i].tipo[0]] : classes['Normal'];
        let pokemonCard = `
        <div class="card ${bgd} ${bgd}-shadow col-sm-6 col-lg-3 col-xxl-2" data-bs-toggle="modal" data-bs-target="#modalPage" id="${data[i].numero}">
        <img class="card-img-top" src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/${data[i].numero}.png" alt="${nome}">
        <h5 class="card-title text-center">${capt(nome.split(' ')[0])}</h5>
        <p class="poke-id">#${data[i].numero}</p>
        </div>
        `;

        cardsHtml.innerHTML += pokemonCard;
    }

    fetch('http://localhost:8080/getPagesNumber/')
        .then(response => response.json())
        .then(data => {
            const numPaginas = +data;
            const novaDiv = document.createElement('div');
            novaDiv.classList.add('row');
            novaDiv.classList.add('justify-content-center');
            /* cardsHtml.innerHTML += `<div class="row justify-content-center">`;*/
            for (let index = 1; index <= numPaginas; index++) {
                if(index == 1 || index == numPaginas || (index >= lastClicked - 2 && index <= lastClicked + 2)) {
                    const novoElemento = document.createElement('button');
                    novoElemento.type = 'button';
                    novoElemento.classList.add('btn');
                    novoElemento.classList.add('btn-Psyduck-mostrarMais');
                    novoElemento.id = 'mostrarMais' + index;
                    novoElemento.innerHTML = index;
                    novaDiv.appendChild(novoElemento);
            
                    novoElemento.onclick = () => {
                        lastClicked = index;
                        fetch(`http://localhost:8080/getAll/?page=${index - 1}`)
                            .then(response => response.json())
                            .then(data => adicionarCards(data))
                            .catch(error => {
                                modalAviso();
                                console.log(error);
                            });

                        window.scrollTo(0, 0);
                        insertDots = true;
                    };

                    if(index == lastClicked) insertDots = true;
                } else if(insertDots) {
                    const dots = document.createElement('button');
                    dots.innerHTML = '...';
                    dots.classList.add('btn');
                    dots.classList.add('btn-Psyduck-mostrarMais');
                    novaDiv.appendChild(dots);
                    insertDots = false;
                }
            }

            cardsHtml.appendChild(novaDiv);
            const btn = document.getElementById(`mostrarMais${lastClicked}`);
            btn.classList.add('active');
            btn.classList.remove('btn-Psyduck-mostrarMais');
        })
        .catch(error => {
            modalAviso();
            console.log(error);
        });

    gerarModalPokemon();
}

/*  * * * * * * * * * * * * *
 *
 *  Mostrar Dados
 *
 * * * * * * * * * * * * * * */

function carregarDados(id) {
    fetch('http://localhost:8080/get/?id=' + id)
        .then(response => response.json())
        .then(data => adicionarDadosModal(data))
        .catch(error => {
            modalAviso(error);
            console.log(error);
        });
}

function adicionarDadosModal(data) {
    console.log(data);
    const editButton = document.querySelector('#edit');
    const saveButton = document.querySelector('#save');
    editButton.hidden = false;
    saveButton.hidden = true;

    const conteudoPokemon = document.getElementById('conteudoPokemon');

    if (data.tipo.length < 2) {
        data.tipo.push('Null');
    }

    const dateObj = new Date(data.lancamento);
    const dia = dateObj.getDate().toString().padStart(2, '0');
    const mes = (dateObj.getMonth() + 1).toString().padStart(2, '0');
    const ano = dateObj.getFullYear().toString();
    const dataFormatada = `${dia}/${mes}/${ano}`;

    let lendario = data.lendario ? 'lendario-y' : 'lendario-n';
    let mitico = data.mitico ? 'mitico-y' : 'mitico-n';

    let modalContent = `
    <p class="poke-id2">#${data.numero}</p>
    <div class="row justify-content-center">
        <p class="modal-title" id="modalPage">${capt(data.nome)}</p>
        <p class="modal-title-jap">${data.nome_jap}</p>
        <p class="poke-type">${data.especie}</p>
    </div>
    <div class="row justify-content-center">
        <p class="tipo-pokemon bgd-${data.tipo[0]} col-4">${data.tipo[0]}</p>
        <p class="tipo-pokemon bgd-${data.tipo[1]} col-4">${data.tipo[1]}</p>
    </div>
    <div class="row justify-content-center">
        <p class="poke-text">${data.peso} KG</p>
        <p class="poke-text">${data.altura} M</p>
        </div>
        <div class="row justify-content-center">
        <p class="poke-desc">Peso</p>
        <p class="poke-desc">Altura</p>
    </div>
    </div>
    <div class="row">
        <div class="col-12">
            <p class="base-stats">Base Stats</p>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">HP</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-bulbasaur" role="progressbar" style="width: ${Math.min(data.hp / 2, 200)}%" aria-valuenow="${data.hp}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.hp}</p>
            </div>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">ATK</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-charmander" role="progressbar" style="width: ${Math.min(data.atk / 2, 200)}%" aria-valuenow="${data.atk}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.atk}</p>
            </div>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">DEF</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-squirtle" role="progressbar" style="width: ${Math.min(data.def / 2, 200)}%" aria-valuenow="${data.def}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.def}</p>
            </div>
        </div>
    </div>
    <div class="row justify-content-center">
        <p class="poke-text">${data.geracao}ª</p>
        <p class="poke-text">${dataFormatada}</p>
        </div>
        <div class="row justify-content-center">
        <p class="poke-desc">Geração</p>
        <p class="poke-desc">Lançamento</p>
    </div>
    <div class="row justify-content-center">
        <p class="poke-rare ${lendario} col-4">Lendario</p>
        <p class="poke-rare ${mitico} col-4">Mitico</p>
    </div>
    <div class="row justify-content-center">
        <p class="poke-descricao-titulo">descrição:</p>
        <p class="poke-descricao">${data.descricao}</p>
    </div>

    `;

    conteudoPokemon.innerHTML = modalContent;
    editButton.onclick = () => editarDadosModal(data, false);
    
    const id = document.querySelector('.poke-id2');
    const deleteForm = document.getElementById('deletar-form');
    deleteForm.value = id.innerHTML.substring(1);
    deleteBtn.onclick = () => {
        document.getElementById('close').click();
        document.getElementById('actual-deletar-form').onsubmit();
    }
}

const collectFormData = async () => {
    const pokemon = {};

    const pokeNumber = document.querySelector('.poke-id2');
    const pokeName = document.getElementById('nome');
    const pokeNameJap = document.getElementById('nome-jap');
    const pokeEspecies = document.getElementById('tipo-pokemon');
    const pokeType1 = document.getElementById('tipo1');
    const pokeType2 = document.getElementById('tipo2');
    const pokeWeight = document.getElementById('peso');
    const pokeHeight = document.getElementById('altura');
    const pokeHP = document.getElementById('--bulbasaur');
    const pokeAtk = document.getElementById('--charmander');
    const pokeDef = document.getElementById('--squirtle');
    const pokeGen = document.getElementById('geracao');
    const pokeRelease = document.getElementById('lancamento');
    const pokeLegendary = document.getElementById('lendario');
    const pokeMitic = document.getElementById('mitico');
    const pokeDescription = document.getElementById('descricao-pokemon');

    const japName = await fetch(`http://localhost:8080/toKatakana/?stringToConvert=${pokeNameJap.value}`);
    let number = pokeNumber.innerText;
    number = number == '' ? 1000 : +number.substring(1);

    pokemon.numero = number;
    pokemon.nome = pokeName.value;
    pokemon.nome_jap = await japName.json();
    pokemon.especie = pokeEspecies.value;
    pokemon.tipo = [pokeType1.value, pokeType2.value];
    pokemon.peso = +pokeWeight.value;
    pokemon.altura = +pokeHeight.value;
    pokemon.hp = +pokeHP.value;
    pokemon.atk = +pokeAtk.value;
    pokemon.def = +pokeDef.value;
    pokemon.geracao = +pokeGen.value;
    const dateFields = pokeRelease.value.split('/');
    pokemon.lancamento = dateFields[2] + '-' + dateFields[1] + '-' + dateFields[0] + 'T00:00:00Z';
    pokemon.lendario = pokeLegendary.classList.contains('lendario-y');
    pokemon.mitico = pokeMitic.classList.contains('mitico-y');
    pokemon.descricao = pokeDescription.value;

    console.log(pokemon)

    return pokemon;
}

document.querySelector('#save').onclick = async () => {
    const pokemon = await collectFormData();
    const method = pokemon.numero == 1000 ? 'post' : 'put';

    fetch(`http://localhost:8080/${method}/`, {
        method: 'POST',
        body: JSON.stringify(pokemon),
        headers: {
            'Content-Type': 'application/json'
        }
    })
        .then(res => res.json())
        .then(data => {
            modalAviso(data.hasOwnProperty('mensagem')?data.mensagem:"Pokemon registrado com o id: " + data.id);
            showAll.onclick();
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

/*  * * * * * * * * * * * * *
 *
 *  Editar Dados
 *
 * * * * * * * * * * * * * * */

function editarDadosModal(data, shouldCreate = false) {
    const editButton = document.querySelector('#edit');
    const saveButton = document.querySelector('#save');
    const deleteButton = document.querySelector('#delete');
    editButton.hidden = true;
    saveButton.hidden = false;
    deleteButton.hidden = true;

    const conteudoPokemon = document.getElementById('conteudoPokemon');

    const dateObj = new Date(data.lancamento);
    const dia = dateObj.getDate().toString().padStart(2, '0');
    const mes = (dateObj.getMonth() + 1).toString().padStart(2, '0');
    const ano = dateObj.getFullYear().toString();
    const dataFormatada = `${dia}/${mes}/${ano}`;

    let modalContent = `
    <p class="poke-id2">${shouldCreate ? "" : "#" + data.numero}</p>
    <div class="row justify-content-center">
    <input class="modal-title-input" type="text" name="nome" id="nome" value="${capt(data.nome)}">
    <input class="modal-title-jap-input" type="text" name="nome-jap" id="nome-jap" value="${data.nome_jap}">
    <input class="poke-type-input" type="text" name="tipo-pokemon" id="tipo-pokemon" value="${data.especie}">
    </div>
    <div class="row justify-content-center">
    <input class="tipo-pokemon-input col-4" type="text" name="tipo1" id="tipo1" value="${data.tipo[0]}">
    <input class="tipo-pokemon-input col-4" type="text" name="tipo2" id="tipo2" value="${data.tipo[1]}">
    </div>
    <div class="row justify-content-center">
        <input class="poke-text-input col-4" type="text" name="tipo2" id="peso" value="${data.peso}">
        <input class="poke-text-input col-4" type="text" name="tipo2" id="altura" value="${data.altura}">
    </div>
    <div class="row justify-content-center">
        <p class="poke-desc">Peso</p>
        <p class="poke-desc">Altura</p>
    </div>
    <div class="row">
    <div class="col-12">
        <p class="base-stats">Base Stats</p>
        <div class="row linha-status justify-content-center">
            <p class="col-2 allign-text">HP</p>
            <div class="col-10 progress poke-bars">
                <input class="progress-bar-input rangers" type="range" id="--bulbasaur" min="0" max="200" value="${(Math.min(data.hp / 2, 200))}">
            </div>
            <p class="col-2 allign-text2" id="--bulbasaur2">${data.hp}</p>
        </div>
        <div class="row linha-status justify-content-center">
            <p class="col-2 allign-text">ATK</p>
            <div class="col-10 progress poke-bars">
                <input class="progress-bar-input rangers" type="range" id="--charmander" min="0" max="200" value="${(Math.min(data.atk / 2, 200))}">
            </div>
            <p class="col-2 allign-text2" id="--charmander2">${data.atk}</p>
        </div>
        <div class="row linha-status justify-content-center">
            <p class="col-2 allign-text">DEF</p>
            <div class="col-10 progress poke-bars">
                <input class="progress-bar-input rangers" type="range" id="--squirtle" min="0" max="200" value="${(Math.min(data.def / 2, 200))}">
            </div>
            <p class="col-2 allign-text2" id="--squirtle2">${data.def}</p>
        </div>
    </div>
    </div>
    <div class="row justify-content-center">
        <input class="poke-text-input col-4" type="text" name="tipo2" id="geracao" value="${data.geracao}">
        <input class="poke-text-input col-4" type="text" name="tipo2" id="lancamento" pattern="\d{2}/\d{2}/\d{4}" inputmode="numeric" value="${dataFormatada}">
    </div>
    <div class="row justify-content-center">
        <p class="poke-desc">Geração</p>
        <p class="poke-desc">Lançamento</p>
    </div>
    </div>
    <div class="row justify-content-center">
        <p class="poke-rare lendario-n col-4 pointer" id="lendario">Lendario</p>
        <p class="poke-rare mitico-n col-4 pointer" id="mitico">Mitico</p>
    </div>
    <div class="row justify-content-center">
        <p class="poke-descricao-titulo">descrição:</p>
        <textarea class="poke-descricao-input scrollbar2" type="text" name="tipo-pokemon" id="descricao-pokemon">${data.descricao}</textarea>
    </div>
    `;

    conteudoPokemon.innerHTML = modalContent;

    const lendario = document.querySelector('#lendario');
    const mitico = document.querySelector('#mitico');
    let lendarioMarca = false;
    let miticoMarca = false;

    lendario.addEventListener('click', function () {
        if (!lendarioMarca) {
            lendario.classList.remove('lendario-n');
            lendario.classList.add('lendario-y');
            lendarioMarca = true;
        } else {
            lendario.classList.remove('lendario-y');
            lendario.classList.add('lendario-n');
            lendarioMarca = false;
        }
    });

    mitico.addEventListener('click', function () {
        if (!miticoMarca) {
            mitico.classList.remove('mitico-n');
            mitico.classList.add('mitico-y');
            mitico.classList.add('shadow-effect');
            miticoMarca = true;
        } else {
            mitico.classList.remove('mitico-y');
            mitico.classList.remove('shadow-effect');
            mitico.classList.add('mitico-n');
            miticoMarca = false;
        }
    });

    const rangers = document.querySelectorAll('.rangers');
    rangers.forEach(range => {
        const rangeValueDisplay = document.querySelector("#" + range.id + 2);
        const defaultValue = range.value;
        range.style.background = `linear-gradient(to right, var(${range.id}) 0%, var(${range.id}) ${defaultValue}%, #f3f3f3 ${defaultValue}%, #f3f3f3 100%)`;
        range.addEventListener('input', () => {
            const value = range.value / 2;
            range.style.background = `linear-gradient(to right, var(${range.id}) 0%, var(${range.id}) ${value}%, #f3f3f3 ${value}%, #f3f3f3 100%)`;
            rangeValueDisplay.textContent = Math.floor(value * 2);
        });
    });
}

const registrar = document.querySelector('#Registrar');
registrar.addEventListener('click', function () {
    abrirModal(undefined, true, undefined);
});

/*  * * * * * * * * * * * * *
 *
 *  Barra Lateral
 *
 * * * * * * * * * * * * * * */
const pesquisar = document.querySelector('#Pesquisar');
const pesquisarForm = document.querySelector('#pesquisar-form');
var pesquisarAberto = false;
const atualizar = document.querySelector('#Atualizar');
const atualizarForm = document.querySelector('#atualizar-form');
let atualizarAberto = false;
const deletar = document.querySelector('#Deletar');
const deletarForm = document.querySelector('#deletar-form');
let deletarAberto = false;

pesquisar.addEventListener('click', function (event) {
    if (event.target === pesquisar && !pesquisarAberto) {
        pesquisar.style.height = 100 + "px";
        pesquisarForm.classList.remove('displayNone');
        pesquisarForm.classList.remove('btn-Charmander');
        pesquisarAberto = true;
    } else if (event.target === pesquisar) {
        pesquisar.style.height = 45 + "px";
        pesquisarForm.classList.add('displayNone');
        pesquisarForm.classList.add('btn-Charmander');
        pesquisarAberto = false;
    }
});

document.getElementById('actual-pesquisar-form').onsubmit = e => {
    e.preventDefault();
    pesquisar.style.height = 45 + "px";
    pesquisarForm.classList.add('displayNone');
    pesquisarForm.classList.add('btn-Charmander');
    pesquisarAberto = false;

    fetch('http://localhost:8080/get/?id=' + pesquisarForm.value)
        .then(response => response.json())
        .then(data => {
            if ('mensagem' in data) {
                modalAviso("Pokemon inexistente");
            } else {
                abrirModal(data.nome, false, data)
            }
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
};

atualizar.addEventListener('click', function (event) {
    if (event.target === atualizar && !atualizarAberto) {
        atualizar.style.height = 100 + "px";
        atualizarForm.classList.remove('displayNone');
        atualizarForm.classList.remove('btn-Charmander');
        atualizarAberto = true;
    } else if (event.target === atualizar) {
        atualizar.style.height = 45 + "px";
        atualizarForm.classList.add('displayNone');
        atualizarForm.classList.add('btn-Charmander');
        atualizarAberto = false;
    }
})

document.getElementById('actual-atualizar-form').addEventListener('submit', e => {
    e.preventDefault();

    atualizar.style.height = 45 + "px";
    atualizarForm.classList.add('displayNone');
    atualizarForm.classList.add('btn-Charmander');
    atualizarAberto = false;

    fetch('http://localhost:8080/get/?id=' + atualizarForm.value)
        .then(response => response.json())
        .then(data => {
            if ('mensagem' in data) {
                modalAviso("Pokemon inexistente");
            } else {
                abrirModal(data.nome, true, data);
            }
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
});

deletar.addEventListener('click', function (event) {
    if (event.target === deletar && !deletarAberto) {
        deletar.style.height = 100 + "px";
        deletarForm.classList.remove('displayNone');
        deletarForm.classList.remove('btn-Charmander');
        deletarAberto = true;
    } else if (event.target === deletar) {
        deletar.style.height = 45 + "px";
        deletarForm.classList.add('displayNone');
        deletarForm.classList.add('btn-Charmander');
        deletarAberto = false;
    }
})

document.getElementById('actual-deletar-form').onsubmit = e => {
    if(e !== undefined) e.preventDefault();

    deletar.style.height = 45 + "px";
    deletarForm.classList.add('displayNone');
    deletarForm.classList.add('btn-Charmander');
    deletarAberto = false;

    fetch('http://localhost:8080/delete/?id=' + deletarForm.value)
        .then(response => response.json())
        .then(data => {
            modalAviso(data.mensagem);
            showAll.click();
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
};

function abrirModal(pokemon = "pokebola", editar = false, data) {
    const closeButton = document.getElementById('close');
    const saveButton = document.getElementById('save');
    document.getElementById('delete').hidden = false;

    pokemon = pokemon.toLowerCase();
    if (data === undefined) {
        data = {
            "numero": 0,
            "nome": "Nome",
            "nome_jap": "和名",
            "geracao": 1,
            "lancamento": "1996-02-27T00:00:00Z",
            "especie": "Especie Pokemon",
            "lendario": false,
            "mitico": false,
            "tipo": [
                "Tipo 1",
                "Tipo 2"
            ],
            "atk": Math.floor(Math.random() * 100),
            "def": Math.floor(Math.random() * 100),
            "hp": Math.floor(Math.random() * 100),
            "altura": 0.0,
            "peso": 0.0
        };
    }

    const cardsHtml = document.getElementById('cards');
    modal.style.setProperty('--random', 0 + '%');
    document.querySelector('#meu-botao3').click();

    const novaDiv = document.createElement('div');

    var imagem
    if (classes[pokemon] === undefined) {
        bgdClass = classes[Object.keys(classes)[Math.floor(Math.random() * 18)]];
        imagem = "pokebola";
    } else {
        bgdClass = classes[pokemon];
        imagem = pokemon;
    }

    let bgd = classes && classes[data.tipo[0]] ? classes[data.tipo[0]] : classes['Normal'];
    let pokemonCard
    if (pokemon === "pokebola") {
        pokemonCard = `
        <div class="card ${bgd} col-sm-6 col-lg-3 col-xxl-2" data-bs-toggle="modal" data-bs-target="#modalPage" id="novaDiv">
        <img class="card-img-top" src="imagens/${imagem}.png" alt="pokeball">
        </div>
        `;
    } else {
        pokemonCard = `
        <div class="card ${bgd} col-sm-6 col-lg-3 col-xxl-2" data-bs-toggle="modal" data-bs-target="#modalPage" id="novaDiv">
        <img class="card-img-top" src="https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/${data.numero}.png" alt="${data.nome}">
        </div>
        `;
    }

    novaDiv.innerHTML = pokemonCard;
    let novoPokemon = novaDiv.querySelector('#novaDiv');
    cardsHtml.appendChild(novoPokemon);

    novoPokemon.style.position = "fixed";
    novoPokemon.style.top = "0";
    novoPokemon.style.left = "-40%";
    novoPokemon.style.width = "100%";
    novoPokemon.style.height = "100%";
    novoPokemon.style.transition = "all 0.5s ease-in-out";
    setTimeout(function () {
        novoPokemon.style.left = "0%";
    }, 0);

    novoPokemon.classList.add("card-to-fullscreen");
    novoPokemon.classList.remove("card");
    novoPokemon.classList.add("disabled");

    if (editar) {
        editarDadosModal(data, data.numero == 0 ? true : false);
    } else {
        adicionarDadosModal(data);
    }

    closeButton.addEventListener('click', function destruirClone1() {
        novoPokemon.style.transition = "all 1s ease-in-out";
        modal.classList.add("slide-out-right");
        // Adiciona um event listener para a transição
        modal.addEventListener('transitionend', function onModalTransitionEnd() {
            setTimeout(function () {
                meuBotao.click();
                modal.classList.remove("slide-out-right");
            }, 500);
            modal.removeEventListener('transitionend', onModalTransitionEnd);
        });

        novoPokemon.style.left = "-50%";
        novoPokemon.style.width = "100%";

        novoPokemon.addEventListener("transitionend", () => {
            novoPokemon.remove();
        });
    });

    saveButton.addEventListener('click', function destruirClone2() {
        novoPokemon.style.transition = "all 1s ease-in-out";
        modal.classList.add("slide-out-right");
        // Adiciona um event listener para a transição
        modal.addEventListener('transitionend', function onModalTransitionEnd() {
            setTimeout(function () {
                meuBotao.click();
                modal.classList.remove("slide-out-right");
            }, 500);
            modal.removeEventListener('transitionend', onModalTransitionEnd);
        });

        novoPokemon.style.left = "-50%";
        novoPokemon.style.width = "100%";

        novoPokemon.addEventListener("transitionend", () => {
            novoPokemon.remove();
        });
    });
}