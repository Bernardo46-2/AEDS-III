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

const showAll = document.getElementById('All');

showAll.onclick = () => {
    fetch('http://localhost:8080/getAll/?page=0')
        .then(response => response.json())
        .then(data => console.log(data))
        .catch(error => {
            modalAviso();
            console.log(error)
        });
}

/* function adicionarCards(data) {

    const pokemonCard = `
    <div class="card bgd-magikarp bgd-magikarp-shadow col-sm-6 col-lg-3 col-xxl-2" data-bs-toggle="modal" data-bs-target="#exampleModal">
    <img class="card-img-top" src="imagens/magikarp.png" alt="magikarp">
    <h5 class="card-title text-center">${data.}</h5>
    </div>
    `;
} */