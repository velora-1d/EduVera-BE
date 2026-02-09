// DOM Elements
const globalSearch = document.querySelector('#globalSearch');
const contactSearch = document.querySelector('#contactSearch');
const groupSearch = document.querySelector('#groupSearch');
const contactsList = document.querySelector('#contactsList');
const groupsList = document.querySelector('#groupsList');
const sendMessageModal = document.querySelector('#sendMessageModal');
const sendGroupModal = document.querySelector('#sendGroupModal');

// Pagination and Loading States
const contactState = {
    page: 1,
    perPage: 10,
    hasMore: true,
    isLoading: false,
    searchTerm: '',
};

const groupState = {
    page: 1,
    perPage: 10,
    hasMore: true,
    isLoading: false,
    searchTerm: '',
};

// Initialization
document.addEventListener('DOMContentLoaded', () => {
    initializeSearch();
    initializeModals();
    setupInfiniteScroll();
    loadContacts();
    loadGroups();
});

// API Loading Functions
async function loadContacts() {
    if (contactState.isLoading || !contactState.hasMore) return;
    contactState.isLoading = true;

    try {
        const response = await fetch(
            `/api/contacts?page=${contactState.page}&perPage=${
                contactState.perPage
            }&search=${encodeURIComponent(contactState.searchTerm)}`
        );
        const result = await response.json();
        const contactsData = result.data || result;

        if (contactsData.contacts.length === 0) {
            contactState.hasMore = false;
            if (contactState.page === 1) {
                contactsList.innerHTML = '<li class="p-3 text-center text-gray-500">No contacts found</li>';
            }
        } else {
            renderContacts(contactsData.contacts);
            contactState.page++;
        }
    } catch (error) {
        console.error('Error loading contacts:', error);
    }

    contactState.isLoading = false;
}

async function loadGroups() {
    if (groupState.isLoading || !groupState.hasMore) return;
    groupState.isLoading = true;

    try {
        const response = await fetch(
            `/api/groups?page=${groupState.page}&perPage=${
                groupState.perPage
            }&search=${encodeURIComponent(groupState.searchTerm)}`
        );
        const result = await response.json();
        const groupsData = result.data || result;

        if (groupsData.length === 0) {
            groupState.hasMore = false;
            if (groupState.page === 1) {
                groupsList.innerHTML = '<li class="p-3 text-center text-gray-500">No groups found</li>';
            }
        } else {
            renderGroups(groupsData);
            groupState.page++;
        }
    } catch (error) {
        console.error('Error loading groups:', error);
    }

    groupState.isLoading = false;
}

// Rendering Functions
function renderContacts(contacts) {
    const fragment = document.createDocumentFragment();

    contacts.forEach((contact) => {
        const li = document.createElement('li');
        li.className =
      'flex justify-between items-center p-3 bg-white dark:bg-zinc-600 rounded-lg';
        li.innerHTML = `
      <div>
        <div class="font-medium dark:text-white">${contact.name || 'Unknown'}</div>
        <div class="text-sm text-gray-500 dark:text-gray-300">${contact.number}</div>
      </div>
      <button data-action="message" data-name="${contact.name}" data-number="${contact.number}"
       class="px-3 py-1 bg-green-500 hover:bg-green-600 text-white rounded">
        <i class="fas fa-comment"></i>
      </button>
    `;
        fragment.appendChild(li);
    });

    contactsList.appendChild(fragment);
}

function renderGroups(groups) {
    const fragment = document.createDocumentFragment();

    groups.forEach((group) => {
        const li = document.createElement('li');
        li.className =
      'flex justify-between items-center p-3 bg-white dark:bg-zinc-600 rounded-lg';
        li.innerHTML = `
      <div>
        <div class="font-medium dark:text-white">${group.groupName || 'Unknown'}</div>
        <div class="text-sm text-gray-500 dark:text-gray-300">ID: ${group.groupId}</div>
      </div>
      <button data-action="group-message" data-name="${group.groupName}" data-group-id="${group.groupId}"
       class="px-3 py-1 bg-green-500 hover:bg-green-600 text-white rounded">
        <i class="fas fa-comment"></i>
      </button>
    `;
        fragment.appendChild(li);
    });

    groupsList.appendChild(fragment);
}

// Search Handlers
function initializeSearch() {
    const handleSearch = (type, term) => {
        if (type === 'contact') {
            contactState.page = 1;
            contactState.hasMore = true;
            contactState.searchTerm = term;
            contactsList.innerHTML = '';
            loadContacts();
        }
        if (type === 'group') {
            groupState.page = 1;
            groupState.hasMore = true;
            groupState.searchTerm = term;
            groupsList.innerHTML = '';
            loadGroups();
        }
    };

    // Global search
    globalSearch.addEventListener(
        'input',
        debounce((e) => {
            const term = e.target.value.toLowerCase();
            handleSearch('contact', term);
            handleSearch('group', term);
        }, 500)
    );

    // Contact search
    contactSearch.addEventListener(
        'input',
        debounce((e) => {
            handleSearch('contact', e.target.value.toLowerCase());
        }, 500)
    );

    // Group search
    groupSearch.addEventListener(
        'input',
        debounce((e) => {
            handleSearch('group', e.target.value.toLowerCase());
        }, 500)
    );
}

// Infinite Scroll
function setupInfiniteScroll() {
    // Get both list containers
    const contactsContainer = contactsList.parentElement;
    const groupsContainer = groupsList.parentElement;

    // Add scroll listener to contacts container
    contactsContainer.addEventListener('scroll', () => {
        const { scrollTop, scrollHeight, clientHeight } = contactsContainer;
        if (
            scrollHeight - scrollTop - clientHeight < 100 &&
      !contactState.isLoading
        ) {
            loadContacts();
        }
    });

    // Add scroll listener to groups container
    groupsContainer.addEventListener('scroll', () => {
        const { scrollTop, scrollHeight, clientHeight } = groupsContainer;
        if (
            scrollHeight - scrollTop - clientHeight < 100 &&
      !groupState.isLoading
        ) {
            loadGroups();
        }
    });
}

// Modal Functions
function initializeModals() {
    // Modal Elements
    const messageForm = document.querySelector('#messageForm');
    const groupMessageForm = document.querySelector('#groupMessageForm');
    const messageRecipient = document.querySelector('#messageRecipient');
    const groupRecipient = document.querySelector('#groupRecipient');
    const messageContent = messageForm.querySelector('#message');
    const groupMessageContent = groupMessageForm.querySelector('#groupMessage');

    let currentNumber = '';
    let currentGroupId = '';

    // Event Delegation
    document.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-action]');
        if (!btn) return;

        const action = btn.dataset.action;
        if (action === 'message') {
            currentNumber = btn.dataset.number;
            messageRecipient.value = `${btn.dataset.name} (${currentNumber})`;
            sendMessageModal.classList.remove('hidden');
            messageContent.focus();
        }
        if (action === 'group-message') {
            currentGroupId = btn.dataset.groupId;
            groupRecipient.value = btn.dataset.name;
            sendGroupModal.classList.remove('hidden');
            groupMessageContent.focus();
        }
    });

    // Form Submissions
    messageForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const submitBtn = messageForm.querySelector('button[type="submit"]');

        try {
            submitBtn.disabled = true;
            submitBtn.innerHTML = 'Sending...';

            sleep(1500);

            const response = await fetch(`${SAFE_API_URL}/api/message/send-message`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    number: currentNumber,
                    message: messageContent.value,
                }),
            });

            if (!response.ok)
                createToast('warning', 'Sending Failed!', 'Message failed to sent');

            createToast('info', 'Message sent!', 'Message successfully sent');
            messageForm.reset();
            sendMessageModal.classList.add('hidden');
        } catch (error) {
            console.error('Send message error:', error);
            createToast('error', 'Ooops!', 'Something happened!');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = 'Send';
        }
    });

    groupMessageForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const submitBtn = groupMessageForm.querySelector('button[type="submit"]');

        try {
            submitBtn.disabled = true;
            submitBtn.innerHTML = 'Sending...';

            sleep(1500);

            const response = await fetch(
                `${SAFE_API_URL}/api/message/send-group-message`,
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        groupId: currentGroupId,
                        message: groupMessageContent.value,
                    }),
                }
            );

            if (!response.ok)
                createToast('warning', 'Sending Failed!', 'Message failed to sent');

            createToast(
                'info',
                'Message sent to group!',
                'Message successfully sent'
            );
            groupMessageForm.reset();
            sendGroupModal.classList.add('hidden');
        } catch (error) {
            console.error('Send group message error:', error);
            createToast('error', 'Ooops!', 'Something happened!');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = 'Send to Group';
        }
    });

    // Close Buttons
    document.querySelectorAll('[id^="btn-close"]').forEach((btn) => {
        btn.addEventListener('click', () => {
            sendMessageModal.classList.add('hidden');
            sendGroupModal.classList.add('hidden');
        });
    });
}
