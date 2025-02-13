// Пример простого JavaScript для интерактивности
document.addEventListener('DOMContentLoaded', function() {
    const btn = document.querySelector('.hero .btn');
    btn.addEventListener('click', function(event) {
        event.preventDefault();
        alert('Спасибо за интерес к нашему сайту!');
    });
});