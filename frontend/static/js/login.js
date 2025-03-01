document.getElementById('loginForm').addEventListener('submit', async function (e) {
    e.preventDefault(); // Предотвращаем стандартную отправку формы

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('/auth/sign-in', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
            credentials: 'include', // Для отправки кук
        });

        const result = await response.json();

        if (response.ok) {
            if (result.requires_2fa) {
                // Перенаправляем на страницу для ввода 2FA кода
                window.location.href = `/login2fa.html?user_id=${result.user_id}`;
            } else {
                // 2FA не требуется, перенаправляем на страницу чата
                window.location.href = '/chat.html';
            }
        } else {
            // Ошибка при входе
            document.getElementById('errorMessage').textContent = result.message || 'Login failed';
        }
    } catch (error) {
        document.getElementById('errorMessage').textContent = 'An error occurred. Please try again.';
    }
});