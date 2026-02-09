const sidebar = document.querySelector('aside');
const backdrop = document.querySelector('#backdrop');
const toggleSidebarButton = document.querySelector('#toggle-sidebar');
const menu = document.querySelector('#dropdownMenu');
const menuButton = document.querySelector('#menuButton');
const openModalButton = document.querySelector('#open-modal');
const closeModalButton = document.querySelector('#close-modal');
const themeToggleBtn = document.querySelector('#theme-toggle');
const themeIcon = document.querySelector('#theme-icon');
const themeUserManualToggleBtn = document.querySelector('#user-manual-theme-toggle');
const themeUserManualIcon = document.querySelector('#user-manual-theme-icon');

const socket = io();
const SESSION_NAME = 'whatsapp_session';

const storedTheme = localStorage.getItem('theme');

if (storedTheme) {
    if (storedTheme === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }
} else {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark');
        localStorage.setItem('theme', 'dark');
    } else {
        document.documentElement.classList.remove('dark');
        localStorage.setItem('theme', 'light');
    }
}

// Initial setup
const storedSession = localStorage.getItem(SESSION_NAME);

if (storedTheme && themeIcon) {
    // Set initial icon based on theme
    if (storedTheme === 'dark') {
        themeIcon.classList.remove('fa-moon');
        themeIcon.classList.add('fa-sun');
        if (themeUserManualIcon) {
            themeUserManualIcon.classList.remove('fa-moon');
            themeUserManualIcon.classList.add('fa-sun');
        }
    } else {
        themeIcon.classList.remove('fa-sun');
        themeIcon.classList.add('fa-moon');
        if (themeUserManualIcon) {
            themeUserManualIcon.classList.remove('fa-sun');
            themeUserManualIcon.classList.add('fa-moon');
        }
    }
}

function userDetail(storedSession) {
    const sessionName = storedSession.name || 'Whatsback User';
    const sessionPicture =
    storedSession.picture || 'https://robohash.org/WhatsbackUser';
    const sessionPhone = storedSession.phone || 'XX XXX XXX';

    const sessionNameElement = document.querySelector('#session-name');
    const sessionPictureElement = document.querySelector('#session-picture');

    sessionNameElement.textContent = sessionName;
    sessionPictureElement.src = sessionPicture;

    const profileSessionNameElement = document.querySelector(
        '#profile-session-name'
    );
    const profileSessionPictureElement = document.querySelector(
        '#profile-session-picture'
    );
    const profileSessionPhoneElement = document.querySelector(
        '#profile-session-phone'
    );
    const profileSessionUsernameElement = document.querySelector(
        '#profile-session-username'
    );

    if (profileSessionNameElement) {
        profileSessionNameElement.textContent = sessionName;
        profileSessionPictureElement.src = sessionPicture;
        profileSessionPhoneElement.textContent = sessionPhone;
        profileSessionUsernameElement.textContent =
      sessionName.toLowerCase().replace(' ', '_') || 'whatsback';
    }

    const userProfileNameElement = document.querySelector('#user-profile-name');
    const userProfilePictureElement = document.querySelector(
        '#user-profile-picture'
    );
    const userProfilePhoneElement = document.querySelector('#user-profile-phone');

    if (userProfileNameElement) {
        userProfileNameElement.textContent = sessionName;
        userProfilePictureElement.src = sessionPicture;
        userProfilePhoneElement.textContent = sessionPhone;
    }
}

function openModal(command = null) {
    const modal = document.getElementById('donateModal');
    modal.classList.remove('hidden');
}

function closeModal() {
    document.getElementById('donateModal').classList.add('hidden');
}

function toggleSidebar() {
    sidebar.classList.toggle('-translate-x-full');
    backdrop.classList.toggle('opacity-0');
    backdrop.classList.toggle('pointer-events-none');
}

// Close sidebar on resize if screen is large (>= lg)
if (sidebar) {
    window.addEventListener('resize', () => {
        if (window.innerWidth >= 1024) {
            sidebar.classList.remove('-translate-x-full');
            backdrop.classList.add('opacity-0');
            backdrop.classList.add('pointer-events-none');
        } else {
            sidebar.classList.add('-translate-x-full');
        }
    });
}

// Toggle dropdown menu visibility
function toggleMenu() {
    menu.classList.toggle('hidden');
}

// Function to toggle theme
function toggleTheme() {
    if (document.documentElement.classList.contains('dark')) {
        localStorage.setItem('theme', 'light');
        themeIcon.classList.remove('fa-sun');
        themeIcon.classList.add('fa-moon');
        if (themeUserManualIcon) {
            themeUserManualIcon.classList.remove('fa-sun');
            themeUserManualIcon.classList.add('fa-moon');
        }
    } else {
        localStorage.setItem('theme', 'dark');
        themeIcon.classList.remove('fa-moon');
        themeIcon.classList.add('fa-sun');
        if (themeUserManualIcon) {
            themeUserManualIcon.classList.remove('fa-moon');
            themeUserManualIcon.classList.add('fa-sun');
        }
    }
    document.documentElement.classList.toggle(
        'dark',
        localStorage.theme === 'dark' ||
      (!('theme' in localStorage) &&
        window.matchMedia('(prefers-color-scheme: dark)').matches)
    );
}

// Connection Status and Logs Simulation
document.addEventListener('DOMContentLoaded', function () {
    const statusElement = document.getElementById('connection-status');
    const logsList = document.querySelector('.logs');
    const isAuthenticated =
    document.querySelector('#is-authenticated')?.dataset.value;

    const btnWhatsbackLogout = document.querySelectorAll('.btn-whatsback-logout');

    if (btnWhatsbackLogout) {
        btnWhatsbackLogout.forEach((btn) => {
            btn.addEventListener('click', () => {
                socket.emit('logout');
            });
        });
    }

    // connection process
    socket.on('connected', () => {
        updateStatus('Connecting...', 'text-yellow-500');
        addLog('Starting connection to WhatsApp...');
        isConnected = true;
    });
    if (isAuthenticated === 'yes' || storedSession) {
        socket.on('ready', (data) => {
            updateStatus('Connected', 'text-green-500');
            addLog('Successfully connected to WhatsApp!');
            userDetail(data.user_info);
            if (data.user_info) {
                localStorage.setItem(SESSION_NAME, JSON.stringify(data.user_info));
            }
        });
        socket.on('authenticated', () => {
            addLog('Whatsapp is authenticated!');
        });
    }
    socket.on('auth_failure', () => {
        localStorage.removeItem(SESSION_NAME);
        window.location.reload();
    });
    socket.on('disconnected', () => updateStatus('Disconnected', 'gray'));
    socket.on('logs', (message) => addLog(message));
    socket.on('client_logout', (message) => {
        addLog(message);
        sleep(2500);
        localStorage.removeItem(SESSION_NAME);
        window.location.href = '/';
    });

    function updateStatus(text, colorClass) {
        if (!statusElement) return;
        statusElement.textContent = text;
        statusElement.className = colorClass;
    }

    function addLog(message) {
        const logEntry = document.createElement('li');
        if (!logEntry || !logsList) return;
        logEntry.innerHTML = `
            <span class="text-gray-400">[${new Date().toLocaleTimeString()}]</span>
            <span class="text-green-400">•</span>
            <span class="text-gray-200">${message}</span>
          `;
        logEntry.classList.add(
            'text-sm',
            'font-mono',
            'flex',
            'items-center',
            'gap-2'
        );
        logsList.appendChild(logEntry);

        // Auto-scroll to bottom
        const container = document.querySelector('.logs-container');
        container.scrollTop = container.scrollHeight;
    }

    themeToggleBtn?.addEventListener('click', toggleTheme);
    themeUserManualToggleBtn?.addEventListener('click', toggleTheme);
    backdrop?.addEventListener('click', toggleSidebar);
    toggleSidebarButton?.addEventListener('click', toggleSidebar);
    menuButton?.addEventListener('click', toggleMenu);
    openModalButton?.addEventListener('click', openModal);
    closeModalButton?.addEventListener('click', closeModal);
});

function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
            func.apply(this, args);
        }, delay);
    };
}

function formatDate(dateInput) {
    // Accept a date string or a Date object.
    const date = typeof dateInput === 'string' ? new Date(dateInput) : dateInput;

    // Array of day names.
    const days = [
        'Sunday',
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday',
    ];
    const dayName = days[date.getDay()];

    // Extract date components.
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-indexed.
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');

    return `${dayName}, ${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

function formatInternationalPhoneNumber(number) {
    number = number.toString();
    // Ensure the number starts with "62" (Indonesia's country code)
    if (!number.startsWith('62')) {
        return number;
    }

    // Remove the country code and add the plus sign
    let localNumber = number.slice(2);
    switch (localNumber.length) {
    case 10: {
        return `+62 ${localNumber.replace(/(\d{3})(\d{3})(\d{4})/, '$1-$2-$3')}`;
    }
    case 11: {
        return `+62 ${localNumber.replace(/(\d{3})(\d{4})(\d{4})/, '$1-$2-$3')}`;
    }
    case 12: {
        return `+62 ${localNumber.replace(/(\d{4})(\d{4})(\d{4})/, '$1-$2-$3')}`;
    }
    case 13: {
        return `+62 ${localNumber.replace(/(\d{4})(\d{4})(\d{5})/, '$1-$2-$3')}`;
    }
    case 14: {
        return `+62 ${localNumber.replace(/(\d{4})(\d{5})(\d{5})/, '$1-$2-$3')}`;
    }
    case 15: {
        return `+62 ${localNumber.replace(/(\d{4})(\d{5})(\d{6})/, '$1-$2-$3')}`;
    }
    default: {
        return number;
    }
    }
}

function createToast(type, title, message, duration = 3000) {
    // Define styles for each type
    const styles = {
        info: {
            bgColor: 'bg-blue-100',
            textColor: 'text-blue-800',
            icon: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>`,
            progressColor: 'bg-blue-500',
        },
        warning: {
            bgColor: 'bg-yellow-100',
            textColor: 'text-yellow-800',
            icon: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>`,
            progressColor: 'bg-yellow-500',
        },
        error: {
            bgColor: 'bg-red-100',
            textColor: 'text-red-800',
            icon: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>`,
            progressColor: 'bg-red-500',
        },
        debug: {
            bgColor: 'bg-purple-100',
            textColor: 'text-purple-800',
            icon: `<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
                </svg>`,
            progressColor: 'bg-purple-500',
        },
    };

    // Get the style for the current type
    const style = styles[type] || styles.info;

    // Create toast element
    const toast = document.createElement('div');
    toast.className = `animate-fade-in ${style.bgColor} p-4 rounded-lg shadow-lg w-80 relative overflow-hidden`;

    // Add progress bar
    const progress = document.createElement('div');
    progress.className = `absolute top-0 left-0 h-1 ${style.progressColor} w-full transition-all duration-[3000ms]`;
    progress.style.width = '100%';

    // Toast content
    toast.innerHTML = `
      <div class="flex items-start space-x-3">
          <div class="flex-shrink-0 ${style.textColor}">
              ${style.icon}
          </div>
          <div class="flex-1">
              <h3 class="font-semibold ${style.textColor}">${title}</h3>
              <p class="text-sm ${style.textColor} mt-1">${message}</p>
          </div>
          <button onclick="this.parentElement.parentElement.remove()" 
                  class="text-gray-400 hover:text-gray-600">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
          </button>
      </div>
  `;

    // Add progress bar to toast
    toast.prepend(progress);

    // Add toast to container
    const container = document.querySelector('#toast-container');
    if (container) container.append(toast);

    // Auto-remove after duration
    let timeout = setTimeout(() => toast.remove(), duration);

    // Pause on hover
    toast.addEventListener('mouseenter', () => {
        clearTimeout(timeout);
        progress.style.transition = 'none';
        progress.style.width = '100%';
    });

    toast.addEventListener('mouseleave', () => {
        progress.style.transition = 'width 3s linear';
        progress.style.width = '0%';
        timeout = setTimeout(() => toast.remove(), duration);
    });

    // Start progress bar animation shortly after
    setTimeout(() => (progress.style.width = '0%'), 50);
}

function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
}
