document.addEventListener('DOMContentLoaded', function() {
    const uploadForm = document.getElementById('uploadForm');
    const progressModal = document.getElementById('progressModal');
    const progressBar = document.getElementById('progress-bar');
    const statusText = document.getElementById('status-text');
    const closeModalButton = document.getElementById('closeModalButton');

    // Обработка отправки формы
    if (uploadForm) {
        uploadForm.addEventListener('submit', function(event) {
            event.preventDefault();
            const files = document.getElementById('fileInput').files;
            if (files.length === 0) {
                alert('Пожалуйста, выберите файлы для загрузки.');
                return;
            }
            uploadFiles(files);
        });
    }

    // Функция загрузки файлов
    function uploadFiles(files) {
        const formData = new FormData();
        for (const file of files) {
            formData.append('files', file);
        }

        // Открываем модальное окно с прогрессом
        showProgressModal(0, 'Загружено: 0%');

        const xhr = new XMLHttpRequest();
        xhr.open('POST', '/upload', true);

        // Отслеживаем прогресс загрузки
        xhr.upload.onprogress = function(event) {
            if (event.lengthComputable) {
                const percent = (event.loaded / event.total) * 100;
                updateProgress(percent, `Загружено: ${Math.round(percent)}%`);
            }
        };

        // Обрабатываем ответ от сервера
        xhr.onload = function() {
            if (xhr.status === 200) {
                updateProgress(100, 'Файлы успешно загружены!');
            } else {
                updateProgress(100, 'Ошибка загрузки!');
            }

            // Показываем кнопку закрытия модального окна
            closeModalButton.classList.remove('hidden');
        };

        // Отправляем файлы
        xhr.send(formData);
    }

    // Функция обновления прогресса
    function updateProgress(percent, status) {
        progressBar.style.width = `${percent}%`;
        statusText.textContent = status;
    }

    // Функция открытия модального окна с прогрессом
    function showProgressModal(percent, status) {
        progressModal.classList.remove('hidden');
        updateProgress(percent, status);
        closeModalButton.classList.add('hidden'); // Скрываем кнопку закрытия, пока идет загрузка
    }

    // Закрытие модального окна
    closeModalButton.addEventListener('click', function() {
        progressModal.classList.add('hidden');
    });
});