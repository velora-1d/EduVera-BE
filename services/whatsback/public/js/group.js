const groupContactsModal = document.querySelector('#groupContactsModal');
const groupContactsSearch = document.querySelector('#groupContactSearch');
const groupContactsList = document.querySelector('#groupContactsList');
const btnCloseGroupContact = document.querySelector(
    '#btn-close-group-contacts'
);
const btnClearSelectedGroup = document.querySelector(
    '#btn-clear-selected-group'
);
const btnOpenGroupContact = document.querySelector('#btn-open-group-contacts');
const btnSendMessageGroup = document.querySelector('#btn-send-message-group');

groupContactsSearch.addEventListener('input', () => {
    const query = groupContactsSearch.value.toLowerCase().trim();

    groupContactsList.querySelectorAll('li').forEach((contact) => {
        const name = contact.dataset.name.toLowerCase();
        const number = contact.dataset.number.toLowerCase();

        if (name.includes(query) || number.includes(query)) {
            contact.style.display = '';
        } else {
            contact.style.display = 'none';
        }
    });
});

const waEditor = createWhatsAppWysiwyg('editor-container', {
    placeholder: 'Type your WhatsApp message here...',
});

const groupNameField = document.querySelector('#group_name');
const groupNumberField = document.querySelector('#group_number');

btnSendMessageGroup.addEventListener('click', async (event) => {
    const groupName = groupNameField.value;
    const groupNumber = groupNumberField.value;
    const message = waEditor.getContent();

    if (groupNumber && message) {
        try {
            const response = await fetch(
                `${SAFE_API_URL}/api/message/send-group-message`,
                {
                    method: 'post',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        groupId: groupNumber,
                        message,
                    }),
                }
            );
            const result = await response.json();
            if (result.status) {
                groupNameField.value = '';
                groupNumberField.value = '';
                waEditor.editorElement.innerHTML = '';
                waEditor.previewElement.innerHTML = '';
                waEditor.editorElement.dispatchEvent(new Event('input'));
                createToast(
                    'info',
                    'Message sent to group!',
                    'Message successfully sent'
                );
            } else {
                createToast(
                    'warning',
                    'Sending Failed!',
                    'Message failed to sent in the group'
                );
            }
        } catch (error) {
            createToast('debug', 'Ooops!', 'Something happend!');
            console.error(error);
        }
    } else {
        createToast('error', 'Alert!', 'Please enter both number and message');
    }
});

function openGroupContactsModal() {
    groupContactsModal.classList.remove('hidden');
}

function closeGroupContactsModal() {
    groupContactsModal.classList.add('hidden');
}

function selectGroupContact() {
    document.querySelector('#group_name').value = this.dataset.name;
    document.querySelector('#group_number').value = this.dataset.number;
    closeGroupContactsModal();
}

function clearSelectedContact() {
    document.querySelector('#group_name').value = '';
    document.querySelector('#group_number').value = '';
}

groupContactsModal.addEventListener('click', function (e) {
    if (e.target === this) closeGroupContactsModal();
});

btnOpenGroupContact.addEventListener('click', openGroupContactsModal);
btnCloseGroupContact.addEventListener('click', closeGroupContactsModal);
btnClearSelectedGroup.addEventListener('click', clearSelectedContact);

let page = 1;
const perPage = 10;
let isLoading = false;
let hasMoreContacts = true;
let searchTerm = '';
let debounceTimeout;

const loadingTextElement = document.querySelector('#loadingText');
const contactListElement = document.querySelector('#groupContactList');
const contactContainerElement = document.querySelector(
    '#groupContactContainer'
);

async function loadGroups() {
    if (isLoading || !hasMoreContacts) return;
    isLoading = true;
    loadingTextElement.classList.remove('hidden');

    try {
        const response = await fetch(
            `${SAFE_API_URL}/api/groups?page=${page}&perPage=${perPage}&search=${encodeURIComponent(
                searchTerm
            )}`
        );
        const result = await response.json();

        const contactsData = result.data || result;

        if (contactsData.length === 0 && page === 1) {
            const li = document.createElement('li');
            li.classList.add('cursor-pointer', 'hover:bg-gray-200', 'p-2', 'rounded');
            li.innerHTML = `
        <div class="flex items-center justify-start gap-2">
          <div class="bg-black h-10 w-10 rounded-full"></div>
          <div>
            <strong class="font-bold">No contacts found</strong>
          </div>
        </div>
      `;
            contactListElement.append(li);
        } else {
            for (const contact of contactsData) {
                const li = document.createElement('li');
                li.classList.add(
                    'selected-number',
                    'cursor-pointer',
                    'hover:bg-gray-200',
                    'p-2',
                    'rounded'
                );
                li.setAttribute('data-name', contact.groupName);
                li.setAttribute('data-number', contact.groupId);
                li.innerHTML = `
          <div class="flex flex-col gap-2">
            <strong class="font-bold">${
    contact.groupName === null || contact.groupName === ''
        ? 'Unknown'
        : contact.groupName
}</strong>
              <p class="font-mono text-gray-500">
                ${contact.groupId}
              </p>
          </div>
        `;
                li.addEventListener('click', selectGroupContact);
                contactListElement.append(li);
            }
        }

        if (contactsData.length < perPage) {
            hasMoreContacts = false;
        }

        page++;
    } catch (error) {
        console.error('Error loading groups:', error);
    }

    isLoading = false;
    loadingTextElement.classList.add('hidden');
}
loadGroups();
