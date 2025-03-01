document.getElementById('twoFAForm').addEventListener('submit', async function (e) {
    e.preventDefault(); // Предотвращаем стандартную отправку формы

    const code = document.getElementById('twoFACode').value;

    // Получаем user_id из URL
    const urlParams = new URLSearchParams(window.location.search);
    const userId = urlParams.get('user_id');

    if (!userId) {
        document.getElementById('errorMessage').textContent = 'User  ID is missing. Please try logging in again.';
        return;
    }

    try {
        const response = await fetch('/auth/verify-2fa', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ user_id: userId, code }), // Передаем user_id как строку
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

document.getElementById('loginForm').addEventListener('submit', async function (event) {
    event.preventDefault(); // Предотвращаем стандартную отправку формы

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // Вызываем функцию signIn
    await signIn(email, password);
});

async function signIn(email, password) {
    const response = await fetch('/auth/sign-in', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (data.requires_2fa) {
        // Показываем поле для ввода 2FA кода
        const code = prompt('Please enter your 2FA code:');
        if (code) {
            await verify2FACode(data.user_id, code); // Передаем user_id
        }
    } else {
        // Переходим на страницу чата
        window.location.href = '/chat.html';
    }
}

async function verify2FACode(userId, code) {
    const response = await fetch('/auth/verify-2fa', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ EncryptedUserID: userId, code }), // Передаем user_id
    });

    const data = await response.json();

    if (response.ok) {
        // Переходим на страницу чата
        window.location.href = '/chat.html';
    } else {
        alert(data.message || 'Invalid 2FA code');
    }
}
