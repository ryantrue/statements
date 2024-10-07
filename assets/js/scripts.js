// scripts.js

// Пример функции для показа алертов или сообщений
function showAlert(message, type = 'info') {
    // Простая логика для показа сообщений пользователю
    const alertBox = document.createElement('div');
    alertBox.className = `alert alert-${type}`;
    alertBox.textContent = message;

    document.body.appendChild(alertBox);

    // Убираем сообщение через 3 секунды
    setTimeout(() => {
        alertBox.remove();
    }, 3000);
}

// Пример глобальной функции для валидации
function validateForm(formId) {
    const form = document.getElementById(formId);
    if (!form) return false;

    // Здесь можно добавить логику валидации
    // Пример:
    return form.checkValidity();
}