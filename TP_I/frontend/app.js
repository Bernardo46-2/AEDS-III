const modalClose = document.querySelector('#close');
const modalClose2 = document.querySelector('#close2');
const modalEdit = document.querySelector('#edit');
const modal = document.querySelector('#modalPage');
const meuBotao = document.querySelector('#meu-botao');
const meuBotao2 = document.querySelector('#meu-botao2');
const scrollbar = document.querySelector('.scrollbar');

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

function getOffset(el) {
    var rect = el.getBoundingClientRect();
    var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    var scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;
    return { top: rect.top + scrollTop, left: rect.left + scrollLeft };
}

function gerarModalPokemon() {
    const button = document.querySelectorAll('.card');
    button.forEach(element => {
        let click = function () {
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
            const image = clonedCard.querySelector('.card-img-top');
            image.style.top = '50%';

            // Obtém todas as classes da variável
            let classes = clonedCard.className.split(' ');
            let bgdClass = classes.find(cls => /^bgd-.+-shadow$/.test(cls));

            clonedCard.classList.remove(bgdClass);

            let randomValue = Math.random() * 400 - 150;
            modal.style.setProperty('--random', randomValue + '%');

            carregarDados(element.id);

            modalClose.addEventListener('click', function destruirClone() {
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
                    element.style.transition = "transitionTmp";
                    clonedCard.remove();
                });
            });
        }


        element.addEventListener('click', click);
    });
}


