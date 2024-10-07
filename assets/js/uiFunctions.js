// Функция обновления списка файлов
function updateFileList() {
    fileList.innerHTML = '';
    for (const file of fileInput.files) {
        const li = document.createElement('li');

        const preview = document.createElement('img');
        preview.src = URL.createObjectURL(file);
        preview.onload = () => URL.revokeObjectURL(preview.src); // Освобождаем память
        li.appendChild(preview);

        const fileName = document.createTextNode(file.name);
        li.appendChild(fileName);

        fileList.appendChild(li);
    }
}

// Функция для отображения и обновления прогресса
function updateProgress(percent, statusTextValue) {
    progressBar.style.width = `${percent}%`;
    statusText.textContent = statusTextValue;
}

// Показ модального окна
function showProgressModal(percent, statusTextValue) {
    progressModal.classList.remove('hidden');
    updateProgress(percent, statusTextValue);
    closeModalButton.classList.add('hidden'); // Скрываем кнопку закрытия во время загрузки
}

// Закрытие модального окна и очистка состояния
function closeProgressModal() {
    progressModal.classList.add('hidden');
    closeModalButton.classList.add('hidden'); // Скрываем кнопку "Закрыть" после закрытия
    clearState(); // Очищаем состояние при закрытии
}