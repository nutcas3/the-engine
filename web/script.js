function toggleTheme() {
    const html = document.documentElement;
    const sunIcon = document.getElementById('theme-icon-sun');
    const moonIcon = document.getElementById('theme-icon-moon');
    
    if (html.classList.contains('dark')) {
        html.classList.remove('dark');
        sunIcon.classList.add('hidden');
        moonIcon.classList.remove('hidden');
        localStorage.setItem('theme', 'light');
    } else {
        html.classList.add('dark');
        sunIcon.classList.remove('hidden');
        moonIcon.classList.add('hidden');
        localStorage.setItem('theme', 'dark');
    }
}

function loadTheme() {
    const savedTheme = localStorage.getItem('theme');
    const sunIcon = document.getElementById('theme-icon-sun');
    const moonIcon = document.getElementById('theme-icon-moon');
    const html = document.documentElement;
    
    if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        html.classList.add('dark');
        sunIcon.classList.remove('hidden');
        moonIcon.classList.add('hidden');
    } else {
        html.classList.remove('dark');
        sunIcon.classList.add('hidden');
        moonIcon.classList.remove('hidden');
    }
}

function connectSSE() {
    const eventSource = new EventSource('/api/stream');
    const output = document.getElementById('sse-output');
    output.innerHTML = '<div class="text-green-400 mb-2">Connected to SSE stream...</div>';
    
    eventSource.onmessage = function(event) {
        const data = JSON.parse(event.data);
        output.innerHTML += '<div class="bg-gray-100 dark:bg-gray-700 rounded-lg p-4 mb-2 border border-gray-300 dark:border-gray-600 text-sm">' + 
            '<span class="text-purple-500 dark:text-purple-400 font-mono">' + new Date().toLocaleTimeString() + '</span>' +
            '<pre class="mt-2 text-gray-700 dark:text-gray-300">' + JSON.stringify(data, null, 2) + '</pre></div>';
    };
    
    eventSource.onerror = function() {
        output.innerHTML += '<div class="text-red-400 mb-2">SSE connection error</div>';
        eventSource.close();
    };
}

setInterval(() => {
    fetch('/api/health/status')
        .then(r => r.json())
        .then(data => {
            const statusEl = document.getElementById('overall-status');
            const healthDisplay = document.getElementById('health-display');
            
            if (data.status === 'healthy') {
                statusEl.className = 'inline-flex items-center px-3 py-1 rounded-full text-white text-sm font-medium bg-green-500/20 border border-green-500/30';
                statusEl.innerHTML = '<span class="w-2 h-2 bg-green-400 rounded-full mr-2"></span>Healthy';
                healthDisplay.textContent = 'Healthy';
                healthDisplay.className = 'text-3xl font-bold text-green-500 dark:text-green-400';
            } else if (data.status === 'degraded') {
                statusEl.className = 'inline-flex items-center px-3 py-1 rounded-full text-white text-sm font-medium bg-yellow-500/20 border border-yellow-500/30';
                statusEl.innerHTML = '<span class="w-2 h-2 bg-yellow-400 rounded-full mr-2"></span>Degraded';
                healthDisplay.textContent = 'Degraded';
                healthDisplay.className = 'text-3xl font-bold text-yellow-500 dark:text-yellow-400';
            } else {
                statusEl.className = 'inline-flex items-center px-3 py-1 rounded-full text-white text-sm font-medium bg-red-500/20 border border-red-500/30';
                statusEl.innerHTML = '<span class="w-2 h-2 bg-red-400 rounded-full mr-2"></span>Unhealthy';
                healthDisplay.textContent = 'Unhealthy';
                healthDisplay.className = 'text-3xl font-bold text-red-500 dark:text-red-400';
            }
        });
}, 30000);

document.addEventListener('DOMContentLoaded', loadTheme);
