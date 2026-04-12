async function shortenUrl() {
    const urlInput = document.getElementById('urlInput');
    const resultDiv = document.getElementById('result');
    const shortUrlLink = document.getElementById('shortUrl');
    
    const originalUrl = urlInput.value.trim();
    
    if (!originalUrl) {
        alert('Пожалуйста, введите URL');
        return;
    }
    
    
    const token = localStorage.getItem('token');
    
    if (!token) {
        
        if (confirm('Для сокращения ссылок нужно зарегистрироваться. Перейти к регистрации?')) {
            window.location.href = '/register';
        }
        return;
    }
    
    try {
        const response = await fetch('/api/urls', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ original_url: originalUrl })
        });
        
        if (response.ok) {
            const data = await response.json();
            shortUrlLink.href = data.short_url;
            shortUrlLink.textContent = data.short_url;
            resultDiv.style.display = 'block';
            urlInput.value = '';
        } else {
            const data = await response.json();
            alert('Ошибка: ' + data.error);
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при создании ссылки');
    }
}


document.addEventListener('DOMContentLoaded', () => {
    const urlInput = document.getElementById('urlInput');
    if (urlInput) {
        urlInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                shortenUrl();
            }
        });
    }
});