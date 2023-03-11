/*  * * * * * * * * * * * * *
 *
 *  Import CSV
 *
 * * * * * * * * * * * * * * */

const importarCsv = document.getElementById('ImportarCSV');
const modalContainer = document.getElementById('modal-container');

function modalAviso(mostrar = "Servidor Desligado") {
    const mensagem = document.getElementById("mensagem-modal");
    mensagem.innerHTML = mostrar;

    modalContainer.classList.remove('out');
    modalContainer.classList.add('one');

    setTimeout(function () {
        modalContainer.classList.add('out');
    }, 4000);

    modalContainer.addEventListener('click', () => {
        modalContainer.classList.add('out');
    });
}

importarCsv.onclick = () => {
    fetch('http://localhost:8080/loadDatabase')
        .then(response => response.json())
        .then(data => modalAviso(data.mensagem))
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

/*  * * * * * * * * * * * * *
 *
 *  Mostrar Todos
 *
 * * * * * * * * * * * * * * */

const showAll = document.getElementById('All');

showAll.onclick = () => {
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
    'pikachu': 'bgd-pikachu',
    'charmander': 'bgd-charmander',
    'squirtle': 'bgd-squirtle',
    'eevee': 'bgd-eevee',
    'gengar': 'bgd-gengar',
    'jigglypuff': 'bgd-jigglypuff',
    'psyduck': 'bgd-psyduck',
    'magikarp': 'bgd-magikarp',
    'abra': 'bgd-abra',
    'machop': 'bgd-machop',
    'geodude': 'bgd-geodude',
    'jolteon': 'bgd-jolteon',
    'vaporeon': 'bgd-vaporeon',
    'flareon': 'bgd-flareon',
    'dragonair': 'bgd-dragonair',
    'zapdos': 'bgd-zapdos',
    'meowth': 'bgd-meowth',
    'minum': 'bgd-minum',
    'quilava': 'bgd-quilava'
};

function adicionarCards(data) {
    const cardsHtml = document.getElementById('cards');
    let content = cardsHtml.innerHTML;
    cardsHtml.innerHTML = "";
    for (let i = 0; i < data.length; i++) {
        let nome = data[i].nome.toLowerCase();


        let bgd = classes[nome] ?? classes[Object.keys(classes)[Math.floor(Math.random() * 18)]];
        let imagem = (classes[nome] === undefined)? "pokebola":nome;
        let pokemonCard = `
        <div class="card ${bgd} ${bgd}-shadow col-sm-6 col-lg-3 col-xxl-2" data-bs-toggle="modal" data-bs-target="#exampleModal" id="${data[i].numero}">
        <img class="card-img-top" src="imagens/${imagem}.png" alt="${nome}">
        <h5 class="card-title text-center">${nome}</h5>
        </div>
        `;

        cardsHtml.innerHTML += pokemonCard;
    }
    cardsHtml.innerHTML += content;
    gerarModalPokemon();
}

/*  * * * * * * * * * * * * *
 *
 *  Mostrar Dados
 *
 * * * * * * * * * * * * * * */

function carregarDados(id) {
    fetch('http://localhost:8080/get/?id='+id)
        .then(response => response.json())
        .then(data => adicionarDadosModal(data))
        .catch(error => {
            modalAviso(error);
            console.log(error)
        });
}

function adicionarDadosModal(data) {
    const conteudoPokemon = document.getElementById('conteudoPokemon');

    if (data.tipo.length < 2) {
        data.tipo.push('Null');
    }

    const dateObj = new Date(data.lancamento);
    const dia = dateObj.getDate().toString().padStart(2, '0');
    const mes = (dateObj.getMonth() + 1).toString().padStart(2, '0');
    const ano = dateObj.getFullYear().toString();
    const dataFormatada = `${dia}/${mes}/${ano}`;

    let lendario = data.lendario?'lendario-y':'lendario-n';
    let mitico = data.mitico?'mitico-y':'mitico-n';

    let modalContent = `
    <div class="row justify-content-center">
        <p class="modal-title" id="exampleModalLabel">${data.nome}</p>
        <p class="modal-title-jap">${data.nome_jap}</p>
        <p class="poke-type">${data.especie}</p>
    </div>
    <div class="row justify-content-center">
        <p class="tipo-pokemon bgd-${data.tipo[0]} col-4">${data.tipo[0]}</p>
        <p class="tipo-pokemon bgd-${data.tipo[1]} col-4">${data.tipo[1]}</p>
    </div>
    <div class="row justify-content-center">
        <div class="col-2 descricao">
            <p class="poke-text">${data.peso} KG</p>
            <p class="poke-desc">Peso</p>
        </div>
        <div class="col-2 descricao">
            <p class="poke-text">${data.altura} M</p>
            <p class="poke-desc">Altura</p>
        </div>
    </div>
    <div class="row">
        <div class="col-12">
            <p class="base-stats">Base Stats</p>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">HP</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-bulbasaur" role="progressbar" style="width: ${Math.min(data.hp/150*100, 150)}%" aria-valuenow="${data.hp}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.hp}</p>
            </div>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">ATK</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-charmander" role="progressbar" style="width: ${Math.min(data.atk/150*100, 150)}%" aria-valuenow="${data.atk}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.atk}</p>
            </div>
            <div class="row linha-status justify-content-center">
                <p class="col-2 allign-text">DEF</p>
                <div class="col-10 progress poke-bars">
                    <div class="progress-bar bgd-squirtle" role="progressbar" style="width: ${Math.min(data.def/150*100, 150)}%" aria-valuenow="${data.def}" aria-valuemin="0" aria-valuemax="300"></div>
                </div>
                <p class="col-2 allign-text2">${data.def}</p>
            </div>
        </div>
    </div>
    <div class="row justify-content-center">
        <div class="col-2 descricao">
            <p class="poke-text">${data.geracao}ª</p>
            <p class="poke-desc">Geração</p>
        </div>
        <div class="col-2 descricao">
            <p class="poke-text">${dataFormatada}</p>
            <p class="poke-desc">Lançamento</p>
        </div>
    </div>
    <div class="row justify-content-center">
        <p class="poke-rare ${lendario} col-4">Lendario</p>
        <p class="poke-rare ${mitico} col-4">Mitico</p>
    </div>
    `;

    conteudoPokemon.innerHTML = modalContent;
}