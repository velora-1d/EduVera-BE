// ================================================================
// DOM Elements
// ================================================================

const scheduleModal = document.querySelector('#schedule-modal');
const scheduleModalForm = document.querySelector('#schedule-modal form');
const btnAddSchedule = document.querySelector('#btn-add-schedule');
const btnCloseSchedule = document.querySelector('#close-schedule-modal');
const cronDetails = document.querySelector('#cron-details');
const btnOpenAddressBookModal = document.querySelector(
    '#open-address-book-modal'
);
const addressBookModal = document.querySelector('#address-book-modal');
const addressBookModalTitle = document.querySelector(
    '#address-book-modal-title'
);
const addressSearchInput = document.querySelector('#address-search');
const btnCloseAddressBookModal = document.querySelector(
    '#close-address-book-modal'
);
const jobTrigger = document.querySelector('#job-trigger');
const cronHelpModal = document.querySelector('#cron-help-modal');
const btnCloseCronHelp = document.querySelectorAll('.close-cron-help');
const btnShowCronHelp = document.querySelector('#show-cron-help');
const jobCronField = document.querySelector('#job-cron');
const deleteModal = document.querySelector('#delete-schedule-modal');
const btnConfirmDelete = document.querySelector('#confirm-delete-schedule');
const btnCloseDeleteModal = document.querySelector('#close-delete-modal');
const tableContainer = document.querySelector('#job-table-container');
const tableSearchInput = document.querySelector('#search-jobs');
const addressContainer = document.querySelector('#address-container');
const addressLoadingTextElement = document.querySelector(
    '#address-loading-text'
);
const addressList = document.querySelector('#address-list');

// ================================================================
// State Variables
// ================================================================

let currentTriggerType = 'send_message';
let currentDeleteId = null;
let page = 1;
const perPage = 10;
let isLoading = false;
let hasMoreContacts = true;
let searchTerm = '';
let debounceTimeout;

// ================================================================
// Table Pagination
// ================================================================
let tablePage = 1;
const tablePerPage = 10;
let isTableLoading = false;
let hasMoreJobs = true;
let tableSearchTerm = '';
let tableDebounce;

// ================================================================
// Modals
// ================================================================

const openScheduleModal = (action = 'create', job = undefined) => {
    scheduleModal.classList.remove('hidden');
    document.querySelector('#form-action').value = action;

    if (action === 'edit') {
        document.querySelector('#job-id').value = job.id;
        document.querySelector('#job-name').value = job.job_name;
        document.querySelector('#job-trigger').value = job.job_trigger;

        const [name, number] = job.target_contact_or_group.split('|');
        document.querySelector('#recipient').value =
      job.job_trigger === 'send_message'
          ? `${name ? ` ${name}` : 'Unknown'} (${number})`
          : `${name} (Group)`;
        document.querySelector('#recipient-number-or-groupid').value = number;
        document.querySelector('#recipient-name').value = name;

        document.querySelector('#job-message').value = job.message;
        document.querySelector('#job-cron').value = job.job_cron_expression;
    }
};

const closeScheduleModal = () => {
    scheduleModal.classList.add('hidden');
    cronDetails.classList.add('hidden');
    document.querySelector('#schedule-form').reset();
};

const showCronHelp = () => cronHelpModal.classList.remove('hidden');
const closeCronHelp = () => cronHelpModal.classList.add('hidden');

const openAddressBookModal = (triggerType) => {
    currentTriggerType = triggerType;
    page = 1;
    hasMoreContacts = true;
    addressList.innerHTML = '';

    addressBookModal.classList.remove('hidden');
    addressBookModalTitle.textContent =
    triggerType === 'send_message' ? 'Select Contact' : 'Select Group';
    populateAddressList();
};

const closeAddressModal = () => addressBookModal.classList.add('hidden');

function selectAddress(item) {
    const recipientField = document.querySelector('#recipient');
    const recipientNumberOrId = document.querySelector(
        '#recipient-number-or-groupid'
    );
    const recipientName = document.querySelector('#recipient-name');
    recipientField.value = item.isGroup
        ? `${item.name} (Group)`
        : `${item.name || 'Unknown'} (${item.number})`;
    recipientNumberOrId.value = item.number;
    recipientName.value = item.name;
    closeAddressModal();
}

const populateAddressList = async (filter = '') => {
    if (isLoading || !hasMoreContacts) return;
    isLoading = true;
    addressLoadingTextElement.classList.remove('hidden');

    if (page === 1) {
        addressList.innerHTML = '';
    }

    try {
        const fetchContacts = async () => {
            const response = await fetch(
                `${SAFE_API_URL}/api/contacts?page=${page}&perPage=${perPage}&search=${encodeURIComponent(
                    searchTerm
                )}`
            );
            const result = await response.json();
            const contactsData = result.data || result;
            return contactsData.contacts.map((contact) => ({
                name: contact.name,
                number: contact.number,
                isGroup: false,
            }));
        };

        const fetchGroups = async () => {
            const response = await fetch(
                `${SAFE_API_URL}/api/groups?page=${page}&perPage=${perPage}&search=${encodeURIComponent(
                    searchTerm
                )}`
            );
            const result = await response.json();
            const groupsData = result.data || result;
            return groupsData.map((group) => ({
                name: group.groupName,
                number: group.groupId,
                isGroup: true,
            }));
        };

        const contactData = await fetchContacts();
        const groupData = await fetchGroups();
        const items =
      currentTriggerType === 'send_message' ? contactData : groupData;

        items
            .filter((item) => {
                return (
                    !item.name ||
          item.name.toLowerCase().includes(filter.toLowerCase()) ||
          (item.number && item.number.includes(filter))
                );
            })
            .forEach((item) => {
                const li = document.createElement('li');
                li.className =
          'select-address p-3 hover:bg-gray-100 dark:hover:bg-zinc-700 rounded-md cursor-pointer';
                li.onclick = () => selectAddress(item);
                li.innerHTML = `
          <div class="font-medium dark:text-white">${
    item.name || 'Unknown'
}</div>
          <div class="text-sm text-gray-500 dark:text-gray-400">${
    item.number
}</div>
        `;
                addressList.append(li);
            });

        if (items.length < perPage) {
            hasMoreContacts = false;
        }

        page++;
    } catch (error) {
        console.error('Error loading addresses:', error);
    }

    isLoading = false;
    addressLoadingTextElement.classList.add('hidden');
};

const validationProcess = debounce(async (event) => {
    let expr = event.target.value;
    if (expr.length >= 5) {
        const isValidExpression = validateCronExpression(expr);
        if (isValidExpression.valid) {
            const response = await fetch(`${SAFE_API_URL}/api/cron-next-runs`, {
                method: 'post',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    exp: expr,
                }),
            });
            const result = await response.json();
            cronDetails.classList.remove('hidden');
            cronDetails.innerHTML = `<p class="mb-3">${result.description}</p>
        <ul>
          ${result.nexRuns
        .map((run) => `<li>&mdash; ${formatDate(run)}</li>`)
        .join('')}
        </ul>`;
        }
    }
}, 500);

async function loadTableData(append = false) {
    if (isTableLoading) return;
    isTableLoading = true;

    try {
        const response = await fetch(
            `${SAFE_API_URL}/api/jobs?page=${tablePage}&perPage=${tablePerPage}&search=${encodeURIComponent(
                tableSearchTerm
            )}`
        );
        const result = await response.json();

        const jobs = result.data || [];
        const totalItems = result.total || 0;
        const totalPages = Math.ceil(totalItems / tablePerPage);

        const tbody = document.querySelector('#job-table-body');

        if (!append) {
            tbody.innerHTML = '';
            tablePage = 1;
        }

        if (jobs.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = `
          <td colspan="7" class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 text-center">
            Data not found
          </td>
        `;
            tbody.appendChild(row);
        } else {
            jobs.forEach((job) => {
                const [name, numberOrId] = job.target_contact_or_group.split('|');
                const targetDisplay = name;
                const statusClass =
          job.job_status === 1
              ? 'bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100'
              : 'bg-yellow-100 dark:bg-yellow-800 text-yellow-800 dark:text-yellow-100';
                const triggerClass =
          job.job_trigger === 'send_message'
              ? 'bg-blue-100 dark:bg-blue-800 text-blue-800 dark:text-blue-100 rounded px-2 py-1'
              : 'bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-100 rounded px-2 py-1';

                const row = document.createElement('tr');
                row.innerHTML = `
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${
    job.id
}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${
    job.job_name
}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
            <span class="${triggerClass}">${
    job.job_trigger === 'send_message' ? 'Send Message' : 'Group Message'
}</span>
          </td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${targetDisplay}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${
    job.job_cron_expression
}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
            <span class="px-2 py-1 ${statusClass} rounded">${
    job.job_status === 1 ? 'Active' : 'Inactive'
}</span>
          </td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
            <button data-scheduleid="${
    job.id
}" class="edit-selected-schedule text-blue-600 mr-2 px-2 py-1 hover:bg-blue-600 hover:text-white rounded-md">
              <i class="fas fa-pen-to-square"></i>
            </button>
            <button data-scheduleid="${job.id}" data-schedulename="${
    job.job_name
}" class="delete-selected-schedule text-red-600 px-2 py-1 hover:bg-red-600 hover:text-white rounded-md">
              <i class="fas fa-trash"></i>
            </button>
          </td>
        `;
                tbody.appendChild(row);
            });
        }

        hasMoreJobs = tablePage < totalPages;

        if (hasMoreJobs && append) {
            tablePage++;
        }
    } catch (error) {
        console.error('Error loading jobs:', error);
    } finally {
        isTableLoading = false;
    }
}

// ================================================================
// Document Ready
// ================================================================
document.addEventListener('DOMContentLoaded', () => {
    // Address Book Events
    addressContainer.addEventListener('scroll', () => {
        if (
            addressContainer.scrollTop + addressContainer.clientHeight >=
      addressContainer.scrollHeight - 10
        ) {
            populateAddressList();
        }
    });

    document.querySelector('#address-search').addEventListener('input', () => {
        clearTimeout(debounceTimeout);
        debounceTimeout = setTimeout(() => {
            page = 1;
            hasMoreContacts = true;
            searchTerm = document.querySelector('#address-search').value.trim();
            populateAddressList();
        }, 500);
    });

    // Schedule Modal Events
    btnAddSchedule.addEventListener('click', () => openScheduleModal('create'));
    btnCloseSchedule.addEventListener('click', closeScheduleModal);
    jobTrigger.addEventListener('change', (event) =>
        openAddressBookModal(event.target.value)
    );

    // Cron Help Events
    btnShowCronHelp.addEventListener('click', showCronHelp);
    btnCloseCronHelp.forEach(element => {
        element.addEventListener('click', closeCronHelp);
    });

    // Address Book Modal Events
    btnOpenAddressBookModal.addEventListener('click', () => {
        let trigger = document.querySelector('#job-trigger').value;
        openAddressBookModal(trigger);
    });
    btnCloseAddressBookModal.addEventListener('click', closeAddressModal);

    scheduleModalForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const formData = new FormData(scheduleModalForm);

        const formObject = {};
        formData.forEach((value, key) => {
            formObject[key] = value;
        });

        const data = {
            job_name: formObject.jobName,
            job_trigger: formObject.jobTrigger,
            job_target_contact_or_group: `${formObject.jobTargetName}|${formObject.jobTargetNumberOrId}`,
            job_message: formObject.jobMessage,
            job_cron_expression: formObject.jobCron,
        };

        if (
            !data.job_name ||
      !data.job_trigger ||
      !data.job_target_contact_or_group ||
      !data.job_message ||
      !data.job_cron_expression
        ) {
            createToast('warning', 'Failed!', 'All fields are required');
            return;
        }

        if (formObject.formAction === 'edit') {
            try {
                const response = await fetch(
                    `${SAFE_API_URL}/api/jobs/${formObject.jobId}`,
                    {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            job_name: data.job_name,
                            job_trigger: data.job_trigger,
                            target_contact_or_group: data.job_target_contact_or_group,
                            message: data.job_message,
                            job_cron_expression: data.job_cron_expression,
                        }),
                    }
                );

                const result = await response.json();

                if (result.success) {
                    createToast('info', 'Success!', 'Job updated');
                    closeScheduleModal();
                    setTimeout(loadTableData, 1000);
                } else {
                    createToast('warning', 'Failed!', 'Job update failed');
                }
            } catch (error) {
                createToast('debug', 'Ooops!', 'Something happened!');
                console.error(error);
            }
        } else {
            try {
                const response = await fetch(`${SAFE_API_URL}/api/jobs`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(data),
                });

                const result = await response.json();

                if (result.success) {
                    createToast('info', 'Success!', 'New job created');
                    closeScheduleModal();
                } else {
                    createToast('warning', 'Failed!', 'Job creation failed');
                }
            } catch (error) {
                createToast('debug', 'Ooops!', 'Something happened!');
                console.error(error);
            }
        }
    });

    // Table Events
    tableContainer.addEventListener('scroll', () => {
        if (
            !isTableLoading &&
      hasMoreJobs &&
      tableContainer.scrollTop + tableContainer.clientHeight >=
        tableContainer.scrollHeight - 100
        ) {
            loadTableData(true);
        }
    });

    tableSearchInput.addEventListener('input', (e) => {
        clearTimeout(tableDebounce);
        tableDebounce = setTimeout(() => {
            tableSearchTerm = e.target.value;
            hasMoreJobs = true;
            loadTableData();
        }, 500);
    });

    document
        .querySelector('#job-table-body')
        .addEventListener('click', async (event) => {
            const editBtn = event.target.closest('.edit-selected-schedule');
            const deleteBtn = event.target.closest('.delete-selected-schedule');

            if (editBtn) {
                const jobId = editBtn.dataset.scheduleid;
                try {
                    const response = await fetch(`${SAFE_API_URL}/api/jobs/${jobId}`);
                    const job = await response.json();
                    openScheduleModal('edit', job.data);
                } catch (error) {
                    console.error('Error fetching job:', error);
                    createToast('error', 'Error!', 'Failed to load schedule');
                }
            }

            if (deleteBtn) {
                currentDeleteId = deleteBtn.dataset.scheduleid;
                document.querySelector('#delete-job-name').textContent =
          deleteBtn.dataset.schedulename;
                deleteModal.classList.remove('hidden');
            }
        });

    btnConfirmDelete.addEventListener('click', async () => {
        if (!currentDeleteId) return;

        try {
            const response = await fetch(
                `${SAFE_API_URL}/api/jobs/${currentDeleteId}`,
                {
                    method: 'DELETE',
                }
            );

            const result = await response.json();

            if (result.success) {
                createToast('success', 'Deleted!', 'Schedule deleted successfully');
                deleteModal.classList.add('hidden');
                setTimeout(loadTableData, 1000);
            } else {
                createToast(
                    'error',
                    'Error!',
                    result.message || 'Failed to delete schedule'
                );
            }
        } catch (error) {
            createToast('error', 'Error!', 'Failed to delete schedule');
            console.error('Delete error:', error);
        } finally {
            currentDeleteId = null;
        }
    });

    btnCloseDeleteModal.addEventListener('click', () => {
        deleteModal.classList.add('hidden');
        currentDeleteId = null;
    });

    loadTableData();
});

// ================================================================
// Cron Event Fields
// ================================================================
jobCronField.addEventListener('keyup', validationProcess);
jobCronField.addEventListener('change', validationProcess);

document.querySelector('#cronOptions').addEventListener('change', (event) => {
    jobCronField.value = event.target.value;
    jobCronField.dispatchEvent(new Event('change'));
});
