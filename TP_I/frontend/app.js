const button = document.querySelectorAll('.card');
const modalClose = document.querySelector('#close');
const modal = document.querySelector('#exampleModal');
const meuBotao = document.querySelector('#meu-botao');

button.forEach(element => {
    let click = function() {

        // Obtém a posição atual da barra de rolagem
        let scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        let scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;

        // Cria uma cópia do elemento original
        let clonedCard = element.cloneNode(true);
        document.body.appendChild(clonedCard);

        // Define as propriedades de posição e tamanho da cópia
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

        modalClose.addEventListener('click', function() {

            modal.classList.add("slide-out-right");
            // Adiciona um event listener para a transição
            modal.addEventListener('transitionend', function onModalTransitionEnd() {
                setTimeout(function() {
                    meuBotao.click();
                    modal.classList.remove("slide-out-right");
                }, 1000);
                modal.removeEventListener('transitionend', onModalTransitionEnd);
            });

            clonedCard.style.top = (getOffset(element).top - scrollTop - 5) + "px";
            clonedCard.style.left = getOffset(element).left - 10 + "px";
            clonedCard.style.width = element.offsetWidth + "px";
            clonedCard.style.height = element.offsetHeight + "px";

            clonedCard.classList.remove('card-to-fullscreen');
            clonedCard.classList.add('card');

            let sidebar = document.querySelector('#sidebar');
            let classes = sidebar.className.split(' ');
            let bgdClass = classes.find(cls => cls.startsWith('bgd-'));
            sidebar.classList.remove(bgdClass);
            sidebar.classList.add('bgd-dark');
            sidebar.classList.add('sidebar-shadow');


            const div = document.querySelector('.slide-from-left');

            clonedCard.addEventListener("transitionend", () => {
                clonedCard.remove();
            });
        });
    }


    element.addEventListener('click', click);
});

function getOffset(el) {
    var rect = el.getBoundingClientRect();
    var scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    var scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;
    return { top: rect.top + scrollTop, left: rect.left + scrollLeft };
}