// Проверка состояния при загрузке страницы
document.addEventListener('DOMContentLoaded', restoreState);

// Открытие файлового проводника при клике на drop-zone
dropZone.addEventListener('click', () => fileInput.click());

// Обновление списка файлов при выборе
fileInput.addEventListener('change', updateFileList);

// Drag-and-drop функциональность
dropZone.addEventListener('dragover', handleDragOver);
dropZone.addEventListener('dragleave', () => dropZone.classList.remove('hover'));
dropZone.addEventListener('drop', handleFileDrop);

// Обработка отправки формы
form.addEventListener('submit', handleFormSubmit);

// Обработка клика на кнопку скачивания Excel
downloadButton.addEventListener('click', () => window.location.href = '/download');

// Закрытие модального окна по кнопке
closeModalButton.addEventListener('click', closeProgressModal);