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
            await verify2FACode(data.user_id, code);
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
        body: JSON.stringify({ user_id: userId, code }),
    });

    const data = await response.json();

    if (data.access_token && data.refresh_token) {
        // Переходим на страницу чата
        window.location.href = '/chat.html';
    } else {
        alert('Invalid 2FA code');
    }
}
