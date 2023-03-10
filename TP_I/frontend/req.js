const importCsvBtn = document.getElementById('Importar CSV');

importCsvBtn.onclick = () => {
    fetch('http://localhost:8080/loadDatabase').then(r => console.log(r));
}
