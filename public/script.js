const input = document.getElementById('address-input');
const list = document.getElementById('autocomplete-list');
const btn = document.getElementById('select-btn');
const resultBox = document.getElementById('result-display');
const resultText = document.getElementById('selected-text');

let timeout = null;

input.addEventListener('input', () => {
    clearTimeout(timeout);
    const val = input.value;
    list.innerHTML = '';

    if (val.length < 3) return;
    list.style.display = 'block';

    timeout = setTimeout(() => {
        fetch('index.php?q=' + encodeURIComponent(val))
            .then(res => res.json())
            .then(data => {
                data.forEach(addr => {
                    const div = document.createElement('div');
                    div.textContent = addr;
                    div.addEventListener('click', () => {
                        input.value = addr;
                        list.innerHTML = '';
                    });
                    list.appendChild(div);
                });
            });
    }, 300); // Задержка 300мс
});

btn.addEventListener('click', () => {

    if (input.value) {
        resultText.textContent = input.value;
        resultBox.style.display = 'block';
    }
});

// Закрыть список при клике вне его
document.addEventListener('click', (e) => {
    list.style.display = 'none';
    if (e.target !== input) list.innerHTML = '';
});