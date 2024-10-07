document.addEventListener('DOMContentLoaded', function() {
    fetch('/api/counterparties')
        .then(response => response.json())
        .then(data => {
            const select = document.getElementById('counterparty');
            data.forEach(counterparty => {
                const option = document.createElement('option');
                option.value = counterparty.id;
                option.textContent = counterparty.name;
                select.appendChild(option);
            });
        })
        .catch(error => {
            console.error('Ошибка загрузки контрагентов:', error);
        });
});