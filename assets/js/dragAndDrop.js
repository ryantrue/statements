// dragAndDrop.js

document.addEventListener('DOMContentLoaded', function() {
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('fileInput');

    if (dropZone && fileInput) {
        // Открытие файлового проводника при клике на drop-zone
        dropZone.addEventListener('click', () => {
            fileInput.click(); // Имитируем нажатие на скрытое поле выбора файла
        });

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
            fileInput.files = event.dataTransfer.files; // Передаем файлы в инпут
            updateFileList(); // Обновляем список файлов
        });
    }
});

// Пример функции обновления списка файлов
function updateFileList() {
    const fileList = document.getElementById('fileList');
    const fileInput = document.getElementById('fileInput');

    if (fileList && fileInput) {
        fileList.innerHTML = '';
        for (const file of fileInput.files) {
            const li = document.createElement('li');
            li.textContent = file.name;
            fileList.appendChild(li);
        }
    }
}