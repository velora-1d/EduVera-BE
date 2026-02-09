const contactModal = document.querySelector('#contactsModal');
const btnContactModal = document.querySelector('#btn-contact-modal');
const btnCloseContact = document.querySelector('#btn-close-contact');
const btnClearSelected = document.querySelector('#clear-selected-contact');

const searchInput = document.querySelector('#contactSearch');
const contactList = document.querySelector('#contactList');
const contacts = contactList.querySelectorAll('li');

const waEditor = createWhatsAppWysiwyg('editor-container', {
    placeholder: 'Type your WhatsApp message here...',
});

const formSendMessage = document.querySelector('#form-send-message');

formSendMessage.addEventListener('click', async (event) => {
    const nameField = document.querySelector('#recipient_name');
    const numberField = document.querySelector('#recipient_number');

    const number = numberField.value;
    const message = waEditor.getContent();

    if (number && message) {
        try {
            const response = await fetch(`${SAFE_API_URL}/api/message/send-message`, {
                method: 'post',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    number,
                    message,
                }),
            });
            const result = await response.json();
            if (result.status) {
                nameField.value = '';
                numberField.value = '';
                waEditor.editorElement.innerHTML = '';
                waEditor.previewElement.innerHTML = '';
                waEditor.editorElement.dispatchEvent(new Event('input'));
                createToast('info', 'Message sent!', 'Message successfully sent');
            } else {
                createToast('warning', 'Sending Failed!', 'Message failed to sent');
            }
        } catch (error) {
            createToast('debug', 'Ooops!', 'Something happend!');
            console.error(error);
        }
    } else {
        createToast('error', 'Alert!', 'Please enter both number and message');
    }
});

function openContactsModal() {
    contactModal.classList.remove('hidden');
}

function closeContactsModal() {
    contactModal.classList.add('hidden');
}

function selectContact() {
    document.querySelector('#recipient_name').value = this.dataset.name;
    document.querySelector('#recipient_number').value = this.dataset.number;
    closeContactsModal();
}

function clearSelectedContact() {
    document.querySelector('#recipient_name').value = '';
    document.querySelector('#recipient_number').value = '';
}

btnContactModal.addEventListener('click', openContactsModal);

btnContactModal.addEventListener('click', function (e) {
    if (e.target === this) closeContactsModal();
});
btnCloseContact.addEventListener('click', closeContactsModal);
btnClearSelected.addEventListener('click', clearSelectedContact);

let page = 1;
const perPage = 10;
let isLoading = false;
let hasMoreContacts = true;
let searchTerm = '';
let debounceTimeout;

const loadingTextElement = document.querySelector('#loadingText');
const contactListElement = document.querySelector('#contactList');
const contactContainerElement = document.querySelector('#contactContainer');

async function loadContacts() {
    if (isLoading || !hasMoreContacts) return;
    isLoading = true;
    loadingTextElement.classList.remove('hidden');

    try {
        const response = await fetch(
            `${SAFE_API_URL}/api/contacts?page=${page}&perPage=${perPage}&search=${encodeURIComponent(
                searchTerm
            )}`
        );
        const result = await response.json();

        const contactsData = result.data || result;

        if (contactsData.contacts.length === 0 && page === 1) {
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
            for (const contact of contactsData.contacts) {
                const li = document.createElement('li');
                li.classList.add(
                    'selected-number',
                    'cursor-pointer',
                    'hover:bg-gray-200',
                    'p-2',
                    'rounded'
                );
                li.setAttribute('data-name', contact.name);
                li.setAttribute('data-number', contact.number);
                li.innerHTML = `
          <div class="flex items-center justify-start gap-2">
            <div class="h-10 w-10 rounded-full bg-black bg-cover bg-center" style="background-image: url('${
    contact.profilePicture
}');"></div>
            <div>
              <strong class="font-bold">${
    contact.name === null ? 'Unknown' : contact.name
}</strong>
              <p class="font-mono text-gray-500">
                ${formatInternationalPhoneNumber(contact.number)}
              </p>
            </div>
          </div>
        `;
                li.addEventListener('click', selectContact);
                contactListElement.append(li);
            }
        }

        if (contactsData.contacts.length < perPage) {
            hasMoreContacts = false;
        }

        page++;
    } catch (error) {
        console.error('Error loading contacts:', error);
    }

    isLoading = false;
    loadingTextElement.classList.add('hidden');
}

function searchContacts() {
    searchTerm = searchInput.value.trim();
    page = 1;
    hasMoreContacts = true;
    contactListElement.innerHTML = '';
    loadContacts();
}

loadContacts();

contactContainerElement.addEventListener('scroll', () => {
    if (
        contactContainerElement.scrollTop + contactContainerElement.clientHeight >=
    contactContainerElement.scrollHeight - 10
    ) {
        loadContacts();
    }
});

searchInput.addEventListener('input', () => {
    clearTimeout(debounceTimeout);
    debounceTimeout = setTimeout(searchContacts, 500);
});

btnCloseContact.addEventListener('click', () => {
    contactModal.classList.add('hidden');
});
