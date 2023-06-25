/* ------------------------------------- UTILS ------------------------------------- */

const modalSearchClose2 = document.getElementById('close-search-modal2');
const modalContainer = document.getElementById('modal-container');
const mensagem = document.getElementById("mensagem-modal");
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

function capt(str) {
    return str.charAt(0).toUpperCase() + str.substring(1).toLowerCase();
}

function modalAviso(mostrar = "Servidor Desligado", tempo = 3000, type = "free") {
    mensagem.innerHTML = mostrar;

    modalContainer.classList.remove('out');
    modalContainer.classList.add('one');

    if (type === "free") {
        modalSearchClose2.style.display = "none";
        let clicked = false;
        let xevent = () => {
            modalContainer.classList.add('out');
            modalContainer.removeEventListener('click', xevent);
            setTimeout(function () {
                modalSearchClose2.style.display = "";
            }, 1000);
            clicked = true;
        }

        modalContainer.addEventListener('click', xevent);

        setTimeout(function () {
            if (!clicked) {
                modalContainer.classList.add('out');
                modalContainer.removeEventListener('click', xevent);
            }
        }, tempo);
    } else if (type === "fix") {
        modalSearchClose2.style.display = "";
        modalSearchClose2.addEventListener('click', z => {
            modalContainer.classList.add('out');

            modalSearchClose2.removeEventListener('click', z);
        });
    }
}

/* ----------------------------------- SIDEBAR ----------------------------------- */

const importarDados = document.getElementById('ImportarDados');
const helpBtn = document.getElementById('Ajuda');

importarDados.onclick = () => {
    if (JSON.parse(localStorage.getItem('zip')) === "true" && JSON.parse(localStorage.getItem("ziptype")) === "lzw") {
        modalAviso(mostrar = "csv compactado em lzw, descompacte primeiro!")
    } else {
        localStorage.setItem('encrypted', JSON.stringify(false));
        localStorage.setItem('zip', JSON.stringify(false));
        localStorage.setItem('ziptype', JSON.stringify("none"));
        sessionStorage.setItem('patternMatchMethod', JSON.stringify(0));
        sessionStorage.setItem('actualPage', JSON.stringify(0));
        sessionStorage.setItem('searchMethod', JSON.stringify(1));
        sessionStorage.setItem('duracao', JSON.stringify(null));
        fetch('http://localhost:8080/loadDatabase')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem, 7000);
                retrieveCardsByPage0();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    }
}

helpBtn.onclick = () => {
    modalAviso("Por favor me dê um emprego (╥﹏╥)");
}

/* ------------------------------------ CARD'S ------------------------------------ */

const tempoDeBusca = document.getElementById("tempoDeBusca");
const balaozinho = document.getElementById("balaozinho")
const dragButton = document.getElementById("dragButton");
const miniModal = document.getElementById("miniModal");
const showAll = document.getElementById('All');
let variavelDeControle = false;
let insertDots = true;
let lastClicked = 1;


const indexMethod = {
    0: "Linear",
    1: "Hashing",
    2: "Arvore B",
    3: "ArvoreB+",
    4: "ArvoreB*",
    5: "Indice Inv.",
    6: "KMP",
    7: "RabinKarp",
};

showAll.onclick = () => {
    let encrypted = localStorage.getItem('encrypted') === "true";
    let zip = localStorage.getItem('zip') === "true";

    if (!encrypted && !zip) {
        fetch('http://localhost:8080/getIdList')
            .then(response => response.json())
            .then(data => {
                paginarIds(data);
                lastClicked = 1;
                insertDots = true;
                sessionStorage.setItem('actualPage', JSON.stringify(1));
                recuperarCards(0)
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else if (encrypted) {
        modalAviso(mostrar = "Sua database foi criptografada e sequestrada<br>Envie um pix pra gente pra recuperar!<br><br>... Ou apenas use a chave q nos fornecemos ...", tempo = 5000)
    } else {
        modalAviso(mostrar = "Database comprimida com " + localStorage.getItem('ziptype'))
    }
}

window.onload = function () {
    if (!localStorage.getItem('encrypted')) {
        localStorage.setItem('encrypted', JSON.stringify(false));
    }
    if (!localStorage.getItem('zip')) {
        localStorage.setItem('zip', JSON.stringify(false));
    }
    sessionStorage.setItem('patternMatchMethod', JSON.stringify(0));
    sessionStorage.setItem('actualPage', JSON.stringify(0));
    sessionStorage.setItem('searchMethod', JSON.stringify(1));
    sessionStorage.setItem('duracao', JSON.stringify(null));
    showAll.click()
};

function paginarIds(ids) {
    const pageSize = 60;
    const pages = Math.ceil(ids.length / pageSize);

    const groups = Array.from({ length: pages }, (_, i) =>
        ids.slice(i * pageSize, (i + 1) * pageSize)
    );

    const object = {
        pages,
        groups
    };

    sessionStorage.setItem('idList', JSON.stringify(object));

    return object;
}

function recuperarCards(pos) {
    const idList = JSON.parse(sessionStorage.getItem('idList'));
    let searchMethod = JSON.parse(sessionStorage.getItem('searchMethod'));

    const ids = idList.groups[pos];
    fetch('http://localhost:8080/getList/?method=' + searchMethod, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(ids)
    })
        .then(response => response.json())
        .then(data => {
            adicionarCards(data.pokemons);
            showTime(indexMethod[searchMethod], data.time);
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

function recuperarCardsIds(ids) {
    const searchMethod = JSON.parse(sessionStorage.getItem('searchMethod'));
    let patternMatchMethod = JSON.parse(sessionStorage.getItem('patternMatchMethod'));
    if (patternMatchMethod == null || patternMatchMethod == 0) {
        patternMatchMethod = searchMethod
    } else {
        patternMatchMethod = +(patternMatchMethod) + 5
    }
    fetch('http://localhost:8080/getList/?method=' + searchMethod, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(ids)
    })
        .then(response => response.json())
        .then(data => {
            adicionarCards(data.pokemons);
            showTime(indexMethod[patternMatchMethod], data.time);
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

function showTime(metodo, tempo) {
    let PMethod = indexMethod[(+JSON.parse(sessionStorage.getItem('patternMatchMethod')) + 5).toString()];
    let duracao = JSON.parse(sessionStorage.getItem('duracao'));
    if (duracao == null) {
        duracao = "--";
    } else {
        duracao += " ms";
    }

    tempoDeBusca.innerHTML = `${metodo}: <strong>${tempo} ms</strong><br>${PMethod}: <strong>${duracao}</strong>`;
    miniModal.classList.add("mostrar");

    if (!variavelDeControle) {
        setTimeout(() => {
            miniModal.classList.add("animated");
            setTimeout(() => {
                miniModal.classList.remove("animated");
            }, 2000);
        }, 2000);
    }

    const intervalId = setInterval(() => {
        if (!variavelDeControle) {
            miniModal.classList.add("animated");
            setTimeout(() => {
                miniModal.classList.remove("animated");
            }, 2000);
        } else {
            // Limpe o intervalo quando a variável de controle mudar para true
            clearInterval(intervalId);
        }
    }, 6000);
};

dragElement(miniModal, dragButton);

function dragElement(element, handle) {
    let pos1 = 0;
    let pos3 = 0;

    handle.addEventListener("mousedown", dragMouseDown);

    function dragMouseDown(e) {
        e = e || window.event;
        e.preventDefault();

        pos3 = e.clientX;

        document.addEventListener("mouseup", closeDragElement);
        document.addEventListener("mousemove", elementDrag);
    }

    function elementDrag(e) {
        e = e || window.event;
        e.preventDefault();

        pos1 = pos3 - e.clientX;
        pos3 = e.clientX;

        let newPosition = element.offsetLeft - pos1;
        let value = window.innerWidth - newPosition;
        if (value < 300 && value > 43) {
            element.style.left = (element.offsetLeft - pos1) + "px";
        }
        /* if (newPosition >= -220 && newPosition <= window.innerWidth - element.offsetWidth - 20) { */
    }

    function closeDragElement() {
        document.removeEventListener("mouseup", closeDragElement);
        document.removeEventListener("mousemove", elementDrag);
    }
}

function adicionarCards(data, gerarPag = true) {
    window.scrollTo(0, 0);
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
    if (gerarPag) {
        gerarPaginacao()
        gerarModalPokemon();
    }
}

function gerarPaginacao() {
    const idList = JSON.parse(sessionStorage.getItem('idList'));
    const cardsHtml = document.getElementById('cards');
    const numPaginas = idList.pages;

    const novaDiv = document.createElement('div');
    novaDiv.classList.add('row');
    novaDiv.classList.add('justify-content-center');

    for (let index = 1; index <= numPaginas; index++) {
        if (index == 1 || index == numPaginas || (index >= lastClicked - 2 && index <= lastClicked + 2)) {
            const novoElemento = document.createElement('button');
            novoElemento.type = 'button';
            novoElemento.classList.add('btn');
            novoElemento.classList.add('btn-Psyduck-mostrarMais');
            novoElemento.id = 'mostrarMais' + index;
            novoElemento.innerHTML = index;
            novaDiv.appendChild(novoElemento);

            novoElemento.onclick = () => {
                lastClicked = index;
                sessionStorage.setItem("actualPage", JSON.stringify(index));
                recuperarCards(index - 1);
                insertDots = true;
            };

            if (index == lastClicked) insertDots = true;
        } else if (insertDots) {
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
}

function retrieveCardsByPage() {
    const actualPage = JSON.parse(sessionStorage.getItem('actualPage'));
    fetch('http://localhost:8080/getIdList')
        .then(response => response.json())
        .then(data => {
            recuperarCardsIds(paginarIds(data).groups[actualPage - 1]);
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

function retrieveCardsByPage0() {
    lastClicked = 1;
    insertDots = true;
    sessionStorage.setItem('actualPage', JSON.stringify(1));
    fetch('http://localhost:8080/getIdList')
        .then(response => response.json())
        .then(data => {
            recuperarCardsIds(paginarIds(data).groups[0]);
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

dragButton.addEventListener("mousedown", function () {
    balaozinho.style.opacity = "0";
    variavelDeControle = true;
    setTimeout(() => {
        balaozinho.style.display = "none";
    }, 1000); // Aguarde a animação de 1 segundo antes de definir display para "none"
});

/* ------------------------------------ UPDATE ------------------------------------ */

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
    <input class="modal-title-jap-input" type="text" name="nome-jap" id="nome-jap" value="${data.nomeJap}">
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
        <p class="poke-desc">Generation</p>
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
            miticoMarca = true;
        } else {
            mitico.classList.remove('mitico-y');
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

const Create = document.querySelector('#Create');
Create.addEventListener('click', function () {
    abrirModal(undefined, true, undefined);
});

/* ------------------------------------ SEARCH ------------------------------------ */

const modalSearchClose = document.getElementById('close-search-modal');
const modalContainer2 = document.getElementById('modal-container2');
const searchIndex = document.getElementById('searchIndex');
const cardsFatherDiv = document.getElementById('cardsFatherDiv');
const mensagem2 = document.getElementById("mensagem-modal2");
const search = document.querySelector('#Search');

search.addEventListener('click', function (event) {
    fatherDivPosition = cardsFatherDiv.style.position;
    cardsFatherDiv.style.position = "fixed";
    modalContainer2.classList.remove('out');
    modalContainer2.classList.add('one');
    modalContainer2.style.zIndex = "9999 !important";

    let modalContent = `
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text"> </p>
        <p class="col-2 modal-search-input4">Descrição</p>
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Nome: </p>
        <input class="col-4 modal-search-input" type="text" name="nome" id="nome" placeholder="Pikachu">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">名前: </p>
        <input class="col-4 modal-search-input" type="text" name="jap" id="japName" placeholder="ピカチュウ">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Especie: </p>
        <input class="col-4 modal-search-input" type="text" name="especie" id="especie" placeholder="Mouse, Fox, Weed">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Tipo: </p>
        <input class="col-4 modal-search-input" type="text" name="tipo" id="tipo" placeholder="Electric, Fire, Poison">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text" style="margin-right:10px; margin-left:-5px">Sobre: </p>
        <textarea class="poke-descricao-input2 scrollbar2" type="text" name="descricao" id="descricao" placeholder="It occasionally uses an electric shock to recharge a fellow Pikachu that is in a weakened state."></textarea>
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text"> </p>
        <p class="col-2 modal-search-input3 right10">Min</p>
        <p class="col-2 modal-search-input3">Max</p>
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Id: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="idI" placeholder="0">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="idF" placeholder="1000">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Hp: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="hpI" placeholder="0">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="hpF" placeholder="200">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Atk: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="atkI" placeholder="0">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="atkF" placeholder="200">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Def: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="defI" placeholder="0">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="defF" placeholder="200">
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Peso: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="pesoI" placeholder="0.1">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="pesoF" placeholder="20.0">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Altura: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="alturaI" placeholder="0.3">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="alturaF" placeholder="0.5">
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Geração: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="geracaoI" placeholder="1">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="geracaoF" placeholder="2">
    </div>
    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Data: </p>
        <input class="col-2 modal-search-input2 right10" type="text" name="tipo-pokemon" id="LancamentoI" placeholder="26/02/1995">
        <input class="col-2 modal-search-input2" type="text" name="tipo-pokemon" id="LancamentoF" placeholder="26/02/1997">
    </div>

    <div class="row justify-content-center">
        <p class="col-4 search-pre-text">Raridade: </p>
        <p class="mit-len-choise lendario2-n col-4 pointer" id="lendario2">Lendario</p>
        <p class="mit-len-choise mitico2-n col-4 pointer" id="mitico2">Mitico</p>
    </div>
    `;

    mensagem2.innerHTML = modalContent;

    const lendario = document.querySelector('#lendario2');
    const mitico = document.querySelector('#mitico2');
    let lendarioMarca = false;
    let miticoMarca = false;

    lendario.addEventListener('click', function () {
        if (!lendarioMarca) {
            lendario.classList.remove('lendario2-n');
            lendario.classList.add('lendario2-y');
            lendarioMarca = true;
        } else {
            lendario.classList.remove('lendario2-y');
            lendario.classList.add('lendario2-n');
            lendarioMarca = false;
        }
    });

    mitico.addEventListener('click', function () {
        if (!miticoMarca) {
            mitico.classList.remove('mitico2-n');
            mitico.classList.add('mitico2-y');
            miticoMarca = true;
        } else {
            mitico.classList.remove('mitico2-y');
            mitico.classList.add('mitico2-n');
            miticoMarca = false;
        }
    });

    let a;

    searchIndex.addEventListener('click', async function mergeSearch(event) {
        a = mergeSearch;
        const nome = document.getElementById('nome').value;
        const especie = document.getElementById('especie').value;
        const tipo = document.getElementById('tipo').value;
        const descricao = document.getElementById('descricao').value;
        const idI = document.getElementById('idI').value;
        let idF = document.getElementById('idF').value;
        const hpI = document.getElementById('hpI').value;
        let hpF = document.getElementById('hpF').value;
        const atkI = document.getElementById('atkI').value;
        let atkF = document.getElementById('atkF').value;
        const defI = document.getElementById('defI').value;
        let defF = document.getElementById('defF').value;
        const pesoI = document.getElementById('pesoI').value;
        let pesoF = document.getElementById('pesoF').value;
        const alturaI = document.getElementById('alturaI').value;
        let alturaF = document.getElementById('alturaF').value;
        const geracaoI = document.getElementById('geracaoI').value;
        let geracaoF = document.getElementById('geracaoF').value;
        const LancamentoI = document.getElementById('LancamentoI').value;
        let LancamentoF = document.getElementById('LancamentoF').value;
        const japName = document.getElementById('japName').value;
        const lendario = "" + (+lendarioMarca);
        const mitico = "" + (+miticoMarca);

        idF = (idF.length == 0) ? idI : idF;
        hpF = (hpF.length == 0) ? hpI : hpF;
        atkF = (atkF.length == 0) ? atkI : atkF;
        defF = (defF.length == 0) ? defI : defF;
        pesoF = (pesoF.length == 0) ? pesoI : pesoF;
        alturaF = (alturaF.length == 0) ? alturaI : alturaF;
        geracaoF = (geracaoF.length == 0) ? geracaoI : geracaoF;
        LancamentoF = (LancamentoF.length == 0) ? LancamentoI : LancamentoF;

        let patternMatch = sessionStorage.getItem("patternMatchMethod");
        if (patternMatch == null) {
            patternMatch = "0";
        }

        fetch('http://localhost:8080/mergeSearch/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                nome: nome,
                especie: especie,
                tipo: tipo,
                descricao: descricao,
                japName: japName,
                idI: idI,
                idF: idF,
                hpI: hpI,
                hpF: hpF,
                atkI: atkI,
                atkF: atkF,
                defI: defI,
                defF: defF,
                pesoI: pesoI,
                pesoF: pesoF,
                alturaI: alturaI,
                alturaF: alturaF,
                geracaoI: geracaoI,
                geracaoF: geracaoF,
                LancamentoI: LancamentoI,
                LancamentoF: LancamentoF,
                Lendario: lendario,
                Mitico: mitico,
                patternMatch: patternMatch,
            })
        })
            .then(response => response.json())
            .then(data => {
                cardsFatherDiv.style.position = fatherDivPosition;
                modalContainer2.classList.add('out');
                paginarIds(data.ids);
                sessionStorage.setItem('duracao', JSON.stringify(data.time));
                lastClicked = 1;
                insertDots = true;
                sessionStorage.setItem('actualPage', JSON.stringify(1));
                recuperarCards(0)
                searchIndex.removeEventListener('click', mergeSearch);
            })
            .catch(error => {
                modalAviso("Nenhum Pokemon Encontrado");
                console.log(error)
            });
    });

    modalSearchClose.addEventListener('click', () => {
        cardsFatherDiv.style.position = fatherDivPosition;
        modalContainer2.classList.add('out');
        searchIndex.removeEventListener('click', a);
    });
});

/* ------------------------------ METODOS DE ORDENAÇAO ------------------------------ */

const ordenar = document.querySelector('#Ordenar');
const ordenarDropdown = document.querySelector('#ordenarDropdown');
const ordenarButtons = document.querySelectorAll('.ordenar-buttons');
const ordenar0 = document.querySelector('#Ordenar0');
const ordenar1 = document.querySelector('#Ordenar1');
const ordenar2 = document.querySelector('#Ordenar2');
const ordenarTransition = ordenar.style.transition;
const ordenarVar3 = ordenar.style.paddingTop;
let ordenarAberto = false;

ordenar.addEventListener('click', function (event) {
    if (event.target === ordenar && !ordenarAberto) {
        ordenar.style.transition = "all 0.4s ease-in-out";
        ordenarDropdown.style.transition = "all 0.4s ease-in-out";
        ordenarDropdown.style.height = "225px";
        ordenarDropdown.style.marginBottom = "15px";
        ordenar.style.height = "225px";
        ordenar.style.paddingTop = "15px";
        ordenarAberto = true;
        window.setTimeout(() => {
            ordenarButtons[0].style.pointerEvents = 'auto';
            ordenarButtons[0].style.opacity = "1";
        }, 100);
        window.setTimeout(() => {
            ordenarButtons[1].style.pointerEvents = 'auto';
            ordenarButtons[1].style.opacity = "1";
        }, 200);
        window.setTimeout(() => {
            ordenarButtons[2].style.pointerEvents = 'auto';
            ordenarButtons[2].style.opacity = "1";
        }, 300);
    } else if (event.target === ordenar) {
        setTimeout(() => {
            ordenarDropdown.style.height = 60 + "px";
            ordenarDropdown.style.marginBottom = "0px";
            ordenar.style.height = 45 + "px";

            ordenar.style.paddingTop = ordenarVar3;
            ordenarButtons.forEach(element => {
                element.style.pointerEvents = 'none';
                element.style.opacity = "0";
            });
            ordenarAberto = false;
            setTimeout(() => {
                ordenar.style.transition = ordenarTransition;
            }, 500);
        }, 200);
        window.setTimeout(() => {
            ordenarButtons[2].style.pointerEvents = 'auto';
            ordenarButtons[2].style.opacity = "0";
        }, 0);
        window.setTimeout(() => {
            ordenarButtons[1].style.pointerEvents = 'auto';
            ordenarButtons[1].style.opacity = "0";
        }, 100);
        window.setTimeout(() => {
            ordenarButtons[0].style.pointerEvents = 'auto';
            ordenarButtons[0].style.opacity = "0";
        }, 200);
    }
})

ordenarButtons.forEach(element => {
    element.style.transition = "all 0.3s ease-in-out";
    element.addEventListener('click', function (event) {
        ordenar.click();
    });
});


ordenar0.onclick = () => {
    if (JSON.parse(localStorage.getItem('zip')) === "false") {
        fetch('http://localhost:8080/ordenacao/?metodo=0')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                showAll.onclick();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else {
        modalAviso(mostrar = "Database comprimida com " + localStorage.getItem('ziptype'));
    }
}
ordenar1.onclick = () => {
    if (JSON.parse(localStorage.getItem('zip')) === "false") {
        fetch('http://localhost:8080/ordenacao/?metodo=1')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                showAll.onclick();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else {
        modalAviso(mostrar = "Database comprimida com " + localStorage.getItem('ziptype'));
    }
}
ordenar2.onclick = () => {
    if (JSON.parse(localStorage.getItem('zip')) === "false") {
        fetch('http://localhost:8080/ordenacao/?metodo=2')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                showAll.onclick();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else {
        modalAviso(mostrar = "Database comprimida com " + localStorage.getItem('ziptype'));
    }
}

/* ------------------------------ ESCOLHA DE INDEXACAO ------------------------------ */

const index = document.querySelector('#Index');
const indexDropdown = document.querySelector('#indexDropdown');
const indexButtons = document.querySelectorAll('.index-buttons');
const indexChoice = document.querySelectorAll('#Index0, #Index1, #Index2, #Index3');
const indexTransition = index.style.transition;
const indexVar3 = index.style.paddingTop;
let indexAberto = false;

index.addEventListener('click', function (event) {
    if (event.target === index && !indexAberto) {
        index.style.transition = "all 0.4s ease-in-out";
        indexDropdown.style.transition = "all 0.4s ease-in-out";
        indexDropdown.style.height = "280px";
        indexDropdown.style.marginBottom = "15px";
        index.style.height = "280px";
        index.style.paddingTop = "15px";
        indexAberto = true;
        window.setTimeout(() => {
            indexButtons[0].style.pointerEvents = 'auto';
            indexButtons[0].style.opacity = "1";
        }, 75);
        window.setTimeout(() => {
            indexButtons[1].style.pointerEvents = 'auto';
            indexButtons[1].style.opacity = "1";
        }, 150);
        window.setTimeout(() => {
            indexButtons[2].style.pointerEvents = 'auto';
            indexButtons[2].style.opacity = "1";
        }, 225);
        window.setTimeout(() => {
            indexButtons[3].style.pointerEvents = 'auto';
            indexButtons[3].style.opacity = "1";
        }, 300);
    } else if (event.target === index) {
        setTimeout(() => {
            indexDropdown.style.height = 60 + "px";
            indexDropdown.style.marginBottom = "0px";
            index.style.height = 45 + "px";

            index.style.paddingTop = indexVar3;
            indexButtons.forEach(element => {
                element.style.pointerEvents = 'none';
                element.style.opacity = "0";
            });
            indexAberto = false;
            setTimeout(() => {
                index.style.transition = indexTransition;
            }, 500);
        }, 200);
        window.setTimeout(() => {
            indexButtons[3].style.pointerEvents = 'auto';
            indexButtons[3].style.opacity = "0";
        }, 0);
        window.setTimeout(() => {
            indexButtons[2].style.pointerEvents = 'auto';
            indexButtons[2].style.opacity = "0";
        }, 75);
        window.setTimeout(() => {
            indexButtons[1].style.pointerEvents = 'auto';
            indexButtons[1].style.opacity = "0";
        }, 150);
        window.setTimeout(() => {
            indexButtons[0].style.pointerEvents = 'auto';
            indexButtons[0].style.opacity = "0";
        }, 225);
    }
})

indexButtons.forEach(element => {
    element.style.transition = "all 0.3s ease-in-out";
    element.addEventListener('click', function (event) {
        index.click();
    });
});

indexChoice.forEach(element => {
    element.onclick = () => {
        const choice = element.id;
        const lastDigit = parseInt(choice.slice(-1));
        sessionStorage.setItem('searchMethod', JSON.stringify(lastDigit));
    }
});

/* ------------------------------ ESCOLHA DE COMPRESSAO ------------------------------ */

const zip = document.querySelector('#Zip');
const zipDropdown = document.querySelector('#zipDropdown');
const zipButtons = document.querySelectorAll('.zip-buttons');
const zip0 = document.querySelector('#Zip0');
const zip1 = document.querySelector('#Zip1');
const zipTransition = zip.style.transition;
const zipVar3 = zip.style.paddingTop;
let zipAberto = false;

zip.addEventListener('click', function (event) {
    if (event.target === zip && !zipAberto) {
        zip.style.transition = "all 0.4s ease-in-out";
        zipDropdown.style.transition = "all 0.4s ease-in-out";
        zipDropdown.style.height = "170px";
        zipDropdown.style.marginBottom = "15px";
        zip.style.height = "170px";
        zip.style.paddingTop = "15px";
        zipAberto = true;
        window.setTimeout(() => {
            zipButtons[0].style.pointerEvents = 'auto';
            zipButtons[0].style.opacity = "1";
        }, 150);
        window.setTimeout(() => {
            zipButtons[1].style.pointerEvents = 'auto';
            zipButtons[1].style.opacity = "1";
        }, 300);
    } else if (event.target === zip) {
        setTimeout(() => {
            zipDropdown.style.height = 60 + "px";
            zipDropdown.style.marginBottom = "0px";
            zip.style.height = 45 + "px";

            zip.style.paddingTop = zipVar3;
            zipButtons.forEach(element => {
                element.style.pointerEvents = 'none';
                element.style.opacity = "0";
            });
            zipAberto = false;
            setTimeout(() => {
                zip.style.transition = zipTransition;
            }, 500);
        }, 200);
        window.setTimeout(() => {
            zipButtons[1].style.pointerEvents = 'auto';
            zipButtons[1].style.opacity = "0";
        }, 0);
        window.setTimeout(() => {
            zipButtons[0].style.pointerEvents = 'auto';
            zipButtons[0].style.opacity = "0";
        }, 150);
    }
})

zipButtons.forEach(element => {
    element.style.transition = "all 0.3s ease-in-out";
    element.addEventListener('click', function (event) {
        zip.click();
    });
});

zip0.onclick = () => {
    if (!JSON.parse(localStorage.getItem('zip'))) {
        localStorage.setItem('zip', JSON.stringify(true));
        localStorage.setItem('ziptype', JSON.stringify("huffman"));

        fetch('http://localhost:8080/zip/?metodo=1')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                const cardsHtml = document.getElementById('cards');
                cardsHtml.innerHTML = '';
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });

    } else if (JSON.parse(localStorage.getItem('ziptype')) === "huffman") {
        localStorage.setItem('zip', JSON.stringify(false));
        localStorage.setItem('ziptype', JSON.stringify("none"));

        fetch('http://localhost:8080/unzip/?metodo=1')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                showAll.onclick();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else {
        modalAviso(mostrar = "Metodo correto = LZW");
    }
};

zip1.onclick = () => {
    if (!JSON.parse(localStorage.getItem('zip'))) {
        localStorage.setItem('zip', JSON.stringify(true));
        localStorage.setItem('ziptype', JSON.stringify("lzw"));

        fetch('http://localhost:8080/zip/?metodo=2')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                const cardsHtml = document.getElementById('cards');
                cardsHtml.innerHTML = '';
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });

    } else if (JSON.parse(localStorage.getItem('ziptype')) === "lzw") {
        localStorage.setItem('zip', JSON.stringify(false));
        localStorage.setItem('ziptype', JSON.stringify("none"));

        fetch('http://localhost:8080/unzip/?metodo=2')
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                showAll.onclick();
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    } else {
        modalAviso(mostrar = "Metodo correto = Huffman");
    }
};

/* ------------------------- ESCOLHA DE CASAMENTO DE PADRAO ------------------------- */

const casamento = document.querySelector('#Casamento');
const casamentoDropdown = document.querySelector('#casamentoDropdown');
const casamentoButtons = document.querySelectorAll('.casamento-buttons');
const casamentoChoice = document.querySelectorAll('#Casamento0, #Casamento1, #Casamento2');
const casamentoTransition = casamento.style.transition;
const casamentoVar3 = casamento.style.paddingTop;
let casamentoAberto = false;

casamento.addEventListener('click', function (event) {
    if (event.target === casamento && !casamentoAberto) {
        casamento.style.transition = "all 0.4s ease-in-out";
        casamentoDropdown.style.transition = "all 0.4s ease-in-out";
        casamentoDropdown.style.height = "225px";
        casamentoDropdown.style.marginBottom = "15px";
        casamento.style.height = "225px";
        casamento.style.paddingTop = "15px";
        casamentoAberto = true;
        window.setTimeout(() => {
            casamentoButtons[0].style.pointerEvents = 'auto';
            casamentoButtons[0].style.opacity = "1";
        }, 100);
        window.setTimeout(() => {
            casamentoButtons[1].style.pointerEvents = 'auto';
            casamentoButtons[1].style.opacity = "1";
        }, 200);
        window.setTimeout(() => {
            casamentoButtons[2].style.pointerEvents = 'auto';
            casamentoButtons[2].style.opacity = "1";
        }, 300);
    } else if (event.target === casamento) {
        setTimeout(() => {
            casamentoDropdown.style.height = 60 + "px";
            casamentoDropdown.style.marginBottom = "0px";
            casamento.style.height = 45 + "px";

            casamento.style.paddingTop = casamentoVar3;
            casamentoButtons.forEach(element => {
                element.style.pointerEvents = 'none';
                element.style.opacity = "0";
            });
            casamentoAberto = false;
            setTimeout(() => {
                casamento.style.transition = casamentoTransition;
            }, 500);
        }, 200);
        window.setTimeout(() => {
            casamentoButtons[2].style.pointerEvents = 'auto';
            casamentoButtons[2].style.opacity = "0";
        }, 0);
        window.setTimeout(() => {
            casamentoButtons[1].style.pointerEvents = 'auto';
            casamentoButtons[1].style.opacity = "0";
        }, 100);
        window.setTimeout(() => {
            casamentoButtons[0].style.pointerEvents = 'auto';
            casamentoButtons[0].style.opacity = "0";
        }, 200);
    }
})

casamentoButtons.forEach(element => {
    element.style.transition = "all 0.3s ease-in-out";
    element.addEventListener('click', function (event) {
        casamento.click();
    });
});

casamentoChoice.forEach(element => {
    element.onclick = () => {
        const choice = element.id;
        const lastDigit = parseInt(choice.slice(-1));
        sessionStorage.setItem('patternMatchMethod', JSON.stringify(lastDigit));
    }
});

/* ------------------------------- ESCOLHA DE CRYPTO ------------------------------- */

const crypto = document.querySelector('#Crypto');
const cryptoDropdown = document.querySelector('#cryptoDropdown');
const cryptoButtons = document.querySelectorAll('.crypto-buttons');
const crypto0 = document.querySelector('#Crypto0');
const crypto1 = document.querySelector('#Crypto1');
const crypto2 = document.querySelector('#Crypto2');
const crypto3 = document.querySelector('#Crypto3');
const cryptoTransition = crypto.style.transition;
const cryptoVar3 = crypto.style.paddingTop;
let cryptoAberto = false;

crypto.addEventListener('click', function (event) {
    if (event.target === crypto && !cryptoAberto) {
        crypto.style.transition = "all 0.4s ease-in-out";
        cryptoDropdown.style.transition = "all 0.4s ease-in-out";
        cryptoDropdown.style.height = "280px";
        cryptoDropdown.style.marginBottom = "15px";
        crypto.style.height = "280px";
        crypto.style.paddingTop = "15px";
        cryptoAberto = true;
        window.setTimeout(() => {
            cryptoButtons[0].style.pointerEvents = 'auto';
            cryptoButtons[0].style.opacity = "1";
        }, 75);
        window.setTimeout(() => {
            cryptoButtons[1].style.pointerEvents = 'auto';
            cryptoButtons[1].style.opacity = "1";
        }, 150);
        window.setTimeout(() => {
            cryptoButtons[2].style.pointerEvents = 'auto';
            cryptoButtons[2].style.opacity = "1";
        }, 225);
        window.setTimeout(() => {
            cryptoButtons[3].style.pointerEvents = 'auto';
            cryptoButtons[3].style.opacity = "1";
        }, 300);
    } else if (event.target === crypto) {
        setTimeout(() => {
            cryptoDropdown.style.height = 60 + "px";
            cryptoDropdown.style.marginBottom = "0px";
            crypto.style.height = 45 + "px";

            crypto.style.paddingTop = cryptoVar3;
            cryptoButtons.forEach(element => {
                element.style.pointerEvents = 'none';
                element.style.opacity = "0";
            });
            cryptoAberto = false;
            setTimeout(() => {
                crypto.style.transition = cryptoTransition;
            }, 500);
        }, 200);
        window.setTimeout(() => {
            cryptoButtons[3].style.pointerEvents = 'auto';
            cryptoButtons[3].style.opacity = "0";
        }, 0);
        window.setTimeout(() => {
            cryptoButtons[2].style.pointerEvents = 'auto';
            cryptoButtons[2].style.opacity = "0";
        }, 75);
        window.setTimeout(() => {
            cryptoButtons[1].style.pointerEvents = 'auto';
            cryptoButtons[1].style.opacity = "0";
        }, 150);
        window.setTimeout(() => {
            cryptoButtons[0].style.pointerEvents = 'auto';
            cryptoButtons[0].style.opacity = "0";
        }, 225);
    }
})

cryptoButtons.forEach(element => {
    element.style.transition = "all 0.3s ease-in-out";
    element.addEventListener('click', function (event) {
        crypto.click();
    });
});

crypto0.onclick = () => criptografar(1);
crypto1.onclick = () => criptografar(2);
crypto2.onclick = () => criptografar(3);
crypto3.onclick = () => criptografar(4);

function criptografar(option = 0) {
    if (JSON.parse(localStorage.getItem('zip')) === "false") {
        let e = localStorage.getItem('encrypted') === "true";
        if (e) {
            modalContainer2.classList.remove('out');
            modalContainer2.classList.add('one');
            modalContainer2.style.zIndex = "9999 !important";

            let modalContent = `
        <div class="row justify-content-center">
            <p class="modal-search-input4">Chave:</p>
        </div>
    
        <div class="row justify-content-center">
            <input class="col-8 modal-search-input" type="text" name="chave" id="chave" placeholder="">
        </div>
        `;

            mensagem2.innerHTML = modalContent;
            let tmp = searchIndex.textContent
            searchIndex.textContent = "descriptografar";

            let xis = (event) => {
                const chave = document.getElementById('chave').value;

                fetch(`http://localhost:8080/decrypt/?metodo=${option}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        key: chave,
                    })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.codigo == 10) {
                            modalContainer2.classList.add('out');
                            searchIndex.removeEventListener('click', xis);
                            setTimeout(function () {
                                searchIndex.textContent = tmp;
                            }, 1000);
                            localStorage.setItem('encrypted', JSON.stringify(false));
                            showAll.click();
                        } else {
                            modalAviso(data.mensagem);
                        }
                    })
                    .catch(error => {
                        modalAviso(error);
                        console.log(error);
                        setTimeout(function () {
                            searchIndex.textContent = tmp;
                        }, 1000);
                    });

            };

            searchIndex.addEventListener('click', xis);

            modalSearchClose.addEventListener('click', () => {
                setTimeout(function () {
                    searchIndex.textContent = tmp;
                }, 1000);
                modalContainer2.classList.add('out');
                searchIndex.removeEventListener('click', xis);
            });
        } else {
            localStorage.setItem('encrypted', JSON.stringify(true))
            fetch(`http://localhost:8080/encrypt/?metodo=${option}`)
                .then(response => response.json())
                .then(data => {
                    modalAviso(mostrar = "Chave:<br>" + data, tempo = 30000, type = "fix");
                    adicionarCards([], false);
                })
                .catch(error => {
                    modalAviso();
                    console.log(error)
                });
        }
    } else {
        modalAviso(mostrar = "Database comprimida com " + localStorage.getItem('ziptype'));
    }
}

/* ------------------------------------- MODAIS ------------------------------------- */

const modalClose = document.querySelector('#close');
const modalSave = document.querySelector('#save');
const modal = document.querySelector('#modalPage');
const meuBotao = document.querySelector('#meu-botao');
const meuBotao2 = document.querySelector('#meu-botao2');
const deleteBtn = document.getElementById('delete');

function getOffset(el) {
    var rect = el.getBoundingClientRect();
    var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    var scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;
    return { top: rect.top + scrollTop, left: rect.left + scrollLeft };
}

function abrirModal(pokemon = "pokebola", editar = false, data) {
    const closeButton = document.getElementById('close');
    const saveButton = document.getElementById('save');
    document.getElementById('delete').hidden = false;

    pokemon = pokemon.toLowerCase();
    if (data === undefined) {
        data = {
            "numero": 0,
            "nome": "Nome",
            "nomeJap": "和名",
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

function gerarModalPokemon() {
    const button = document.querySelectorAll('.card');
    button.forEach(element => {
        let click = function () {
            const deleteButton = document.querySelector('#delete');
            deleteButton.hidden = false;

            // Obtém a posição atual da barra de rolagem
            let scrollTop = window.pageYOffset || document.documentElement.scrollTop;

            // Cria uma cópia do elemento original
            let clonedCard = element.cloneNode(true);
            element.classList.add('originalCard');
            element.style.opacity = "0";
            transitionTmp = element.style.transition;
            element.style.transition = "none";
            document.body.appendChild(clonedCard);

            // Define as propriedades de posição e tamanho da cópia
            clonedCard.id = 'clonedCard';
            clonedCard.style.position = "fixed";
            clonedCard.style.top = (getOffset(element).top - scrollTop) + "px";
            clonedCard.style.left = (getOffset(element).left - 10) + "px";
            clonedCard.style.width = element.offsetWidth + "px";
            clonedCard.style.height = element.offsetHeight + "px";

            // Define uma transição para a cópia
            clonedCard.style.transition = "all 0.5s ease-in-out";

            // Redimensiona a cópia para ocupar a tela inteira
            clonedCard.style.top = "0";
            clonedCard.style.left = "0";
            clonedCard.style.width = "100%";
            clonedCard.style.height = "100%";

            // Remove a classe "card" da cópia e adiciona a classe "card-to-fullscreen"
            clonedCard.classList.remove("card");
            clonedCard.classList.add("card-to-fullscreen");
            clonedCard.classList.add("disabled");
            const cardTitle = clonedCard.querySelector('.card-title');
            cardTitle.remove();
            const cardId = clonedCard.querySelector('.poke-id');
            cardId.remove();
            const image = clonedCard.querySelector('.card-img-top');
            image.style.top = '50%';

            // Obtém todas as classes da variável
            let classes = clonedCard.className.split(' ');
            let bgdClass = classes.find(cls => /^bgd-.+-shadow$/.test(cls));

            clonedCard.classList.remove(bgdClass);

            let randomValue = Math.random() * 400 - 150;
            modal.style.setProperty('--random', randomValue + '%');

            carregarDados(element.id);

            modalClose.addEventListener('click', function destruirClone3() {
                clonedCard.style.transition = "all 1s ease-in-out";
                modal.classList.add("slide-out-right");
                // Adiciona um event listener para a transição
                modal.addEventListener('transitionend', function onModalTransitionEnd() {
                    setTimeout(function () {
                        meuBotao.click();
                        modal.classList.remove("slide-out-right");
                    }, 500);
                    modal.removeEventListener('transitionend', onModalTransitionEnd);
                });

                clonedCard.style.top = (getOffset(element).top - scrollTop - 5) + "px";
                clonedCard.style.left = getOffset(element).left - 15 + "px";
                clonedCard.style.width = element.offsetWidth + "px";
                clonedCard.style.height = element.offsetHeight + "px";
                clonedCard.classList.remove('card-to-fullscreen');
                clonedCard.classList.add('card');

                const div = document.querySelector('.slide-from-left');

                clonedCard.addEventListener("transitionend", () => {
                    element.classList.remove('originalCard')
                    element.style.opacity = "1";
                    setTimeout(function () {
                        element.style.transition = transitionTmp;
                    }, 500);
                    clonedCard.remove();
                });
            });

            modalSave.addEventListener('click', function destruirClone4() {
                clonedCard.style.transition = "all 1s ease-in-out";
                modal.classList.add("slide-out-right");
                // Adiciona um event listener para a transição
                modal.addEventListener('transitionend', function onModalTransitionEnd() {
                    setTimeout(function () {
                        meuBotao.click();
                        modal.classList.remove("slide-out-right");
                    }, 500);
                    modal.removeEventListener('transitionend', onModalTransitionEnd);
                });

                clonedCard.style.top = (getOffset(element).top - scrollTop - 5) + "px";
                clonedCard.style.left = getOffset(element).left - 15 + "px";
                clonedCard.style.width = element.offsetWidth + "px";
                clonedCard.style.height = element.offsetHeight + "px";
                clonedCard.classList.remove('card-to-fullscreen');
                clonedCard.classList.add('card');

                const div = document.querySelector('.slide-from-left');

                clonedCard.addEventListener("transitionend", () => {
                    element.classList.remove('originalCard')
                    element.style.opacity = "1";
                    setTimeout(function () {
                        element.style.transition = transitionTmp;
                    }, 500);
                    clonedCard.remove();
                });
            });
        }


        element.addEventListener('click', click);
    });
}

function adicionarDadosModal(data) {
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
        <p class="modal-title-jap">${data.nomeJap}</p>
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
        <p class="poke-desc">Generation</p>
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
    close = document.getElementById('close');
    deleteBtn.onclick = () => {
        close.click();
        fetch('http://localhost:8080/delete/?id=' + data.numero)
            .then(response => response.json())
            .then(data => {
                modalAviso(data.mensagem);
                retrieveCardsByPage()
            })
            .catch(error => {
                modalAviso();
                console.log(error)
            });
    };
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
    pokemon.nomeJap = await japName.json();
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

    return pokemon;
}

function carregarDados(id) {
    fetch('http://localhost:8080/get/?id=' + id)
        .then(response => response.json())
        .then(data => {
            adicionarDadosModal(data)
        })
        .catch(error => {
            modalAviso(mostrar = error);
            console.log(error);
        });
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
            modalAviso(data.hasOwnProperty('mensagem') ? data.mensagem : "Pokemon registrado com o id: " + data.id);
            retrieveCardsByPage()
        })
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

/* ----------------------------------- SCROLLBAR ----------------------------------- */

const scrollbar = document.querySelector('#scrollbar');
let isMouseDown = false;
let startY;

window.addEventListener('resize', function () {
    const totalHeight = document.documentElement.scrollHeight;
    const scrollbarHeight = window.innerHeight;
    const thumbHeight = Math.max(scrollbarHeight * (window.innerHeight / totalHeight), 20);
    const thumbPosition = (scrollbarHeight - thumbHeight) * (window.scrollY / (totalHeight - window.innerHeight));

    scrollbar.style.height = `${thumbHeight}px`;
    scrollbar.style.top = `${thumbPosition}px`;
});

document.addEventListener("scroll", () => {
    const totalHeight = document.documentElement.scrollHeight;
    const scrollbarHeight = window.innerHeight;
    const thumbHeight = Math.max(scrollbarHeight * (window.innerHeight / totalHeight), 20);
    const thumbPosition = (scrollbarHeight - thumbHeight) * (window.scrollY / (totalHeight - window.innerHeight));

    scrollbar.style.height = `${thumbHeight}px`;
    scrollbar.style.top = `${thumbPosition}px`;
});

let distanceFromTop = 0;
let thumbHeight = 0

scrollbar.addEventListener("mousedown", (e) => {
    distanceFromTop = e.clientY - scrollbar.getBoundingClientRect().y;
    const totalHeight = document.documentElement.scrollHeight;
    const scrollbarHeight = window.innerHeight;
    thumbHeight = Math.max(scrollbarHeight * (window.innerHeight / totalHeight), 20);
});

function handleMouseMove(e) {
    if (!isMouseDown) return;

    const scrollPercentage = ((e.clientY - (distanceFromTop)) / (window.innerHeight - thumbHeight));
    const scrollPosition = (scrollPercentage * (document.documentElement.scrollHeight - window.innerHeight));


    window.scrollTo({
        top: scrollPosition,
        behavior: 'instant'
    });
}

scrollbar.addEventListener("mousedown", (e) => {
    e.preventDefault();
    isMouseDown = true;
    scrollbar.classList.add("dragging");
});

document.addEventListener("mousemove", handleMouseMove);

document.addEventListener("mouseup", () => {
    if (!isMouseDown) return;

    isMouseDown = false;
    scrollbar.classList.remove("dragging");
});