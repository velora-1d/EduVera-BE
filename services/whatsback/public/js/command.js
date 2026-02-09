const btnAddCommand = document.querySelector('#btn-add-command');
const btnCloseCommand = document.querySelector('#btn-close-command');
const btnCancelCommand = document.querySelector('#btn-cancel-command');
const btnConfirmDeleteCommand = document.querySelector(
    '#btn-confirm-delete-command'
);
let currentCommandId = null;

function openModalCommand(command, type) {
    const modal = document.querySelector('#commandModal');
    if (command && type === 'edit') {
        document.querySelector('#modalTitle').textContent = 'Edit Command';
        document.querySelector('#commandId').value = command.id?.trim();
        document.querySelector('#commandName').value = command.name?.trim();
        document.querySelector('#commandResponse').value = command.response?.trim();
        document.querySelector('#form_action').value = type;
    } else {
        document.querySelector('#modalTitle').textContent = 'Add New Command';
        document.querySelector('#commandId').value = '';
        document.querySelector('#commandName').value = '';
        document.querySelector('#commandResponse').value = '';
        document.querySelector('#form_action').value = '';
    }
    modal.classList.remove('hidden');
}

function closeModalCommand() {
    document.querySelector('#form_action').value = 'add';
    document.querySelector('#commandModal').classList.add('hidden');
}

function openDeleteModalCommand(commandId) {
    currentCommandId = commandId;
    document.querySelector('#form_action').value = 'delete';
    document.querySelector('#deleteModal').classList.remove('hidden');
}

function closeDeleteModalCommand() {
    document.querySelector('#deleteModal').classList.add('hidden');
}

function cancelDeleteModalCommand() {
    document.querySelector('#form_action').value = 'add';
    document.querySelector('#deleteModal').classList.add('hidden');
}

async function confirmDelete(event) {
    event.preventDefault();

    const response = await fetch(
        `${SAFE_API_URL}/api/command/${currentCommandId}`,
        {
            method: 'DELETE',
        }
    );

    const result = await response.json();

    if (result.status) {
        createToast('info', 'Success!', result.message);
        setTimeout(() => {
            loadNewCommands();
        }, 1000);
    } else {
        createToast('warning', 'Failed!', result.message);
    }
    closeDeleteModalCommand();
}

async function loadNewCommands() {
    try {
    // Fetch the updated list of commands from the server
        const response = await fetch(`${SAFE_API_URL}/api/command`);
        const result = await response.json();

        const data = result.data || result;

        // Get the table body element by its ID
        const tbody = document.querySelector('#commandsTableBody');
        // Clear existing rows
        tbody.innerHTML = '';

        // Iterate over the commands and create table rows
        data.commands.forEach((cmd, index) => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${cmd.id}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${cmd.command}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">${cmd.response}</td>
          <td class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
            <button class="text-blue-600 mr-2 px-2 py-1 hover:bg-blue-600 hover:text-white rounded-md edit-btn">
              <i class="fas fa-pen-to-square"></i>
            </button>
            <button class="text-red-600 px-2 py-1 hover:bg-red-600 hover:text-white rounded-md delete-btn">
              <i class="fas fa-trash"></i>
            </button>
          </td>
        `;
            tbody.appendChild(tr);
        });

        document.querySelectorAll('.edit-btn').forEach((btn) => {
            btn.addEventListener('click', () => {
                const row = btn.closest('tr');
                const command = {
                    id: row.cells[0].textContent.trim(),
                    name: row.cells[1].textContent.trim(),
                    response: row.cells[2].textContent.trim(),
                };
                openModalCommand(command, 'edit');
            });
        });

        document.querySelectorAll('.delete-btn').forEach((btn) => {
            btn.addEventListener('click', () => {
                const commandId = btn.closest('tr').cells[0].textContent.trim();
                openDeleteModalCommand(commandId);
            });
        });
    } catch (error) {
        console.error('Error loading commands:', error);
    }
}

// Commands Page
btnAddCommand?.addEventListener('click', () =>
    openModalCommand(currentCommandId, 'add')
);
btnCloseCommand?.addEventListener('click', closeModalCommand);
btnCancelCommand?.addEventListener('click', cancelDeleteModalCommand);
btnConfirmDeleteCommand?.addEventListener('click', confirmDelete);

document
    .querySelector('#commandForm')
    .addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(event.target);
        const collected = Object.fromEntries(formData.entries());
        const data = {
            command_id: collected.commandId,
            command: collected.commandName,
            response: collected.commandResponse,
            action: collected.action,
        };

        switch (data.action) {
        case 'add': {
            const response = await fetch(`${SAFE_API_URL}/api/command`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    command: data.command,
                    response: data.response,
                }),
            });

            const result = await response.json();

            if (result.status) {
                createToast('info', 'Success!', result.message);
                setTimeout(() => {
                    loadNewCommands();
                }, 1000);
            } else {
                createToast('warning', 'Failed!', result.message);
            }

            break;
        }
        case 'edit': {
            const response = await fetch(
                `${SAFE_API_URL}/api/command/${data.command_id}`,
                {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        command: data.command,
                        response: data.response,
                    }),
                }
            );

            const result = await response.json();

            if (result.status) {
                createToast('info', 'Success!', result.message);
                setTimeout(() => {
                    loadNewCommands();
                }, 1000);
            } else {
                createToast('warning', 'Failed!', result.message);
            }

            break;
        }
        }

        setTimeout(closeModalCommand, 300);
    });

document.querySelectorAll('.edit-btn').forEach((btn) => {
    btn.addEventListener('click', () => {
        const row = btn.closest('tr');
        const command = {
            id: row.cells[0].textContent,
            name: row.cells[1].textContent,
            response: row.cells[2].textContent,
        };
        openModalCommand(command, 'edit');
    });
});

document.querySelectorAll('.delete-btn')?.forEach((btn) => {
    btn.addEventListener('click', () => {
        const commandId = btn.closest('tr').cells[0].textContent;
        openDeleteModalCommand(commandId);
    });
});
