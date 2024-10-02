const fileInput = document.getElementById('fileInput');
const fileList = document.getElementById('fileList');
const form = document.getElementById('uploadForm');
const progressModal = document.getElementById('progressModal');
const progressBar = document.getElementById('progress-bar');
const statusText = document.getElementById('status-text');
const closeModalButton = document.getElementById('closeModalButton');
const dropZone = document.getElementById('drop-zone');
const downloadButton = document.getElementById('downloadButton');

// Проверка состояния при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    restoreState();
});

// Открытие файлового проводника при клике на drop-zone
dropZone.addEventListener('click', () => {
    fileInput.click(); // Имитируем нажатие на скрытое поле выбора файла
});

// Обновление списка файлов при выборе
fileInput.addEventListener('change', updateFileList);

// Drag-and-drop функциональность
dropZone.addEventListener('dragover', (event) => {
    event.preventDefault();
    dropZone.classList.add('hover');
});

dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('hover');
});

dropZone.addEventListener('drop', (event) => {
    event.preventDefault();
    dropZone.classList.remove('hover');
    fileInput.files = event.dataTransfer.files;
    updateFileList();
});

// Обработка отправки формы
form.addEventListener('submit', (event) => {
    event.preventDefault();
    if (fileInput.files.length === 0) {
        alert('Пожалуйста, выберите файлы для загрузки.');
        return;
    }
    uploadFiles();
});

// Обработка клика на кнопку скачивания Excel
downloadButton.addEventListener('click', () => {
    window.location.href = '/download';
});

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

// Функция загрузки файлов
function uploadFiles() {
    const formData = new FormData();
    for (const file of fileInput.files) {
        formData.append('files', file);
    }

    // Показываем модальное окно с прогрессом
    showProgressModal(0, 'Загружено: 0%');
    saveState(0, 'Загружено: 0%');

    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/upload', true);

    xhr.upload.onprogress = (event) => {
        if (event.lengthComputable) {
            const percent = (event.loaded / event.total) * 100;
            const statusTextValue = `Загружено: ${Math.round(percent)}%`;
            updateProgress(percent, statusTextValue);
            saveState(percent, statusTextValue);
        }
    };

    xhr.onload = () => {
        if (xhr.status === 200) {
            updateProgress(100, 'Файлы успешно загружены!');
        } else {
            updateProgress(100, 'Ошибка загрузки!');
        }

        // Показываем кнопку "Закрыть" только после завершения загрузки
        closeModalButton.classList.remove('hidden');
        clearState(); // Очищаем состояние
    };

    xhr.send(formData);
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

// Восстановление состояния из LocalStorage
function restoreState() {
    const savedProgress = localStorage.getItem('uploadProgress');
    const savedStatusText = localStorage.getItem('statusText');

    if (savedProgress && savedStatusText) {
        showProgressModal(savedProgress, savedStatusText);

        // Если загрузка завершилась, показываем кнопку "Закрыть"
        if (savedProgress === '100') {
            closeModalButton.classList.remove('hidden');
        }
    }
}

// Сохранение состояния в LocalStorage
function saveState(percent, statusTextValue) {
    localStorage.setItem('uploadProgress', percent.toString());
    localStorage.setItem('statusText', statusTextValue);
}

// Очистка состояния из LocalStorage
function clearState() {
    localStorage.removeItem('uploadProgress');
    localStorage.removeItem('statusText');
}

// Закрытие модального окна по кнопке
closeModalButton.addEventListener('click', () => {
    progressModal.classList.add('hidden');
    closeModalButton.classList.add('hidden'); // Скрываем кнопку "Закрыть" после закрытия
    clearState(); // Очищаем состояние при закрытии
});