const button = document.querySelectorAll('.card');
const modalClose = document.querySelector('#close');
const modalClose2 = document.querySelector('#close2');
const modalEdit = document.querySelector('#edit');
const modal = document.querySelector('#exampleModal');
const modal2 = document.querySelector('#modalInput');
const meuBotao = document.querySelector('#meu-botao');
const meuBotao2 = document.querySelector('#meu-botao2');
const meuBotao3 = document.querySelector('#meu-botao3');
const scrollbar = document.querySelector('.scrollbar');
const rangers = document.querySelectorAll('.rangers');

window.addEventListener('resize', function() {
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

button.forEach(element => {
    let click = function() {
        // Obtém a posição atual da barra de rolagem
        let scrollTop = window.pageYOffset || document.documentElement.scrollTop;

        // Cria uma cópia do elemento original
        let clonedCard = element.cloneNode(true);
        document.body.appendChild(clonedCard);

        // Define as propriedades de posição e tamanho da cópia
        element.id = 'originalCard';
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

        modalClose.addEventListener('click', function destruirClone() {
            modal.classList.add("slide-out-right");
            // Adiciona um event listener para a transição
            modal.addEventListener('transitionend', function onModalTransitionEnd() {
                setTimeout(function() {
                    meuBotao.click();
                    modal.classList.remove("slide-out-right");
                }, 500);
                modal.removeEventListener('transitionend', onModalTransitionEnd);
            });

            clonedCard.style.top = (getOffset(element).top - scrollTop - 5) + "px";
            clonedCard.style.left = getOffset(element).left - 18 + "px";
            clonedCard.style.width = element.offsetWidth + "px";
            clonedCard.style.height = element.offsetHeight + "px";

            clonedCard.classList.remove('card-to-fullscreen');
            clonedCard.classList.add('card');

            const div = document.querySelector('.slide-from-left');

            clonedCard.addEventListener("transitionend", () => {
                    clonedCard.remove();
            });
        });
    }


    element.addEventListener('click', click);
});

modalEdit.addEventListener('click', function() {
    console.log("uheuheuhe");
})

modalClose2.addEventListener('click', function() {
    const clonedCard = document.querySelector('#clonedCard');
    const original = document.querySelector('#originalCard');

    let scrollTop = window.pageYOffset || document.documentElement.scrollTop;

    modal2.classList.add("slide-out-right");
    // Adiciona um event listener para a transição
    modal2.addEventListener('transitionend', function onModalTransitionEnd() {
        setTimeout(function() {
            meuBotao2.click();
            modal2.classList.remove("slide-out-right");
        }, 500);
        modal2.removeEventListener('transitionend', onModalTransitionEnd);
    });

    if (original) {
        clonedCard.style.top = (getOffset(original).top - scrollTop - 5) + "px";
        clonedCard.style.left = getOffset(original).left - 18 + "px";
        clonedCard.style.width = original.offsetWidth + "px";
        clonedCard.style.height = original.offsetHeight + "px";
    
        clonedCard.classList.remove('card-to-fullscreen');
        clonedCard.classList.add('card');
    
        clonedCard.addEventListener("transitionend", () => {
            clonedCard.remove();
            original.id = "";
        });
    } else {
        clonedCard.style.top = "0";
        clonedCard.style.left = "-100%";
        clonedCard.addEventListener("transitionend", () => {
            clonedCard.remove();
            original.id = "";
        });
    }
})

const rangeValueDisplay = document.querySelector('.input-print');

rangers.forEach(range => {
    const rangeValueDisplay = document.querySelector("#" + range.id + 2);
    const defaultValue = range.value / 3;
    range.style.background = `linear-gradient(to right, var(${range.id}) 0%, var(${range.id}) ${defaultValue}%, #f3f3f3 ${defaultValue}%, #f3f3f3 100%)`;  
    range.addEventListener('input', () => {
        const value = Math.ceil((range.value - range.min) / (range.max - range.min) * 100);
        range.style.background = `linear-gradient(to right, var(${range.id}) 0%, var(${range.id}) ${value}%, #f3f3f3 ${value}%, #f3f3f3 100%)`;
        rangeValueDisplay.textContent = value;
    });
});

const lendario = document.querySelector('#lendario');
const mitico = document.querySelector('#mitico');
let lendarioMarca = false;
let miticoMarca = false;

lendario.addEventListener('click', function() {
    if (!lendarioMarca) {
        lendario.classList.remove('lendario-n');
        lendario.classList.add('lendario-y');
        lendario.classList.add('shadow-effect');
        lendarioMarca = true;
    } else {
        lendario.classList.remove('lendario-y');
        lendario.classList.remove('shadow-effect');
        lendario.classList.add('lendario-n');
        lendarioMarca = false;
    }
});

mitico.addEventListener('click', function() {
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

const registrar = document.querySelector('#Registrar');
registrar.addEventListener('click', function() {
    modal2.classList.add("slide-from-left");
    meuBotao3.click();
    
    let clonedCard = button[0].cloneNode(true);
    document.body.appendChild(clonedCard);
    
    
    clonedCard.id = 'clonedCard';
    clonedCard.style.position = "fixed";
    clonedCard.style.top = "0";
    clonedCard.style.left = "-100%";
    clonedCard.style.width = "100%";
    clonedCard.style.height = "100%";

    clonedCard.style.transition = "all 0.5s ease-in-out";

    setTimeout(function() {
        clonedCard.style.left = "0%";
    }, 0);

    // Remove a classe "card" da cópia e adiciona a classe "card-to-fullscreen"
    clonedCard.classList.remove("card");
    clonedCard.classList.add("card-to-fullscreen");
    clonedCard.classList.add("disabled");

    modal2.classList.remove("slide-from-left");
})

const pesquisar = document.querySelector('#Pesquisar');
const pesquisarForm = document.querySelector('#pesquisar-form');
let pesquisaAberta = false;
pesquisar.addEventListener('click', function() {
    if (!pesquisaAberta) {
        pesquisar.style.height = 100 + "px";
        pesquisarForm.classList.remove('displayNone');
        pesquisarForm.classList.remove('btn-Charmander');
        pesquisaAberta = true;
    } else {
        pesquisar.style.height = 45 + "px";
        pesquisarForm.classList.add('displayNone');
        pesquisarForm.classList.add('btn-Charmander');
        pesquisaAberta = false;
    }
})

const atualizar = document.querySelector('#Atualizar');
const atualizarForm = document.querySelector('#atualizar-form');
let atualizarAberta = false;
atualizar.addEventListener('click', function() {
    if (!atualizarAberta) {
        atualizar.style.height = 100 + "px";
        atualizarForm.classList.remove('displayNone');
        atualizarForm.classList.remove('btn-Squirtle');
        atualizarAberta = true;
    } else {
        atualizar.style.height = 45 + "px";
        atualizarForm.classList.add('displayNone');
        atualizarForm.classList.add('btn-Squirtle');
        atualizarAberta = false;
    }
})

const deletar = document.querySelector('#Deletar');
const deletarForm = document.querySelector('#deletar-form');
let deletarAberta = false;
deletar.addEventListener('click', function() {
    if (!deletarAberta) {
        deletar.style.height = 100 + "px";
        deletarForm.classList.remove('displayNone');
        deletarForm.classList.remove('btn-Eevee');
        deletarAberta = true;
    } else {
        deletar.style.height = 45 + "px";
        deletarForm.classList.add('displayNone');
        deletarForm.classList.add('btn-Eevee');
        deletarAberta = false;
    }
})