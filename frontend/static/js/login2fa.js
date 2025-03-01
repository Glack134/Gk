document.getElementById('twoFAForm').addEventListener('submit', async function (e) {
    e.preventDefault(); // Предотвращаем стандартную отправку формы

    const code = document.getElementById('twoFACode').value;

    // Получаем user_id из URL (например, http://localhost:8080/login2fa.html?user_id=123)
    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('user_id');

    if (!userId) {
        document.getElementById('errorMessage').textContent = 'User ID is missing. Please try logging in again.';
        return;
    }

    try {
        const response = await fetch('/auth/verify-2fa', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ user_id: parseInt(userId), code }), // Убедитесь, что user_id является числом
            credentials: 'include', // Для отправки кук
        });

        const result = await response.json();

        if (response.ok) {
            // Успешный вход, перенаправляем на страницу чата
            window.location.href = '/chat.html';
        } else {
            // Ошибка при проверке 2FA кода
            document.getElementById('errorMessage').textContent = result.message || 'Invalid 2FA code';
        }
    } catch (error) {
        document.getElementById('errorMessage').textContent = 'An error occurred. Please try again.';
    }
});