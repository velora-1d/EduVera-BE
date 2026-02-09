/**
 * Creates a reusable and configurable WhatsApp WYSIWYG editor.
 *
 * @param {string} containerId - The id of the container element.
 * @param {Object} options - Configuration options:
 *   options.toolbarFormats {Array} - Array of format objects.
 *     Each format object can be of two types:
 *       - Wrap: { id, label, type: 'wrap', markers: [start, end] }
 *       - Prefix: { id, label, type: 'prefix', prefix: string }
 *   options.previewRules {Array} - Array of preview rule objects { regex, replacement }.
 *   options.placeholder {string} - Placeholder text for the editor.
 *
 * The function creates:
 *  - A full‑width toolbar whose buttons equally share the width.
 *  - A full‑width, contenteditable editor.
 *  - A live preview area that renders formatted content.
 *
 * @returns {Object} - An object with methods to retrieve content:
 *   - getContent(): Returns the plain text (with markers).
 *   - getFormattedContent(): Returns the HTML (formatted preview).
 *   - editorElement and previewElement for further manipulation.
 */
function createWhatsAppWysiwyg(containerId, options = {}) {
    // Default configuration
    const defaultConfig = {
        toolbarFormats: [
            {
                id: 'bold',
                label: '<i class="fa-solid fa-bold"></i>',
                type: 'wrap',
                markers: ['*', '*'],
            },
            {
                id: 'italic',
                label: '<i class="fa-solid fa-italic"></i>',
                type: 'wrap',
                markers: ['_', '_'],
            },
            {
                id: 'strike',
                label: '<i class="fa-solid fa-strikethrough"></i>',
                type: 'wrap',
                markers: ['~', '~'],
            },
            {
                id: 'inline-code',
                label: '<i class="fa-solid fa-code"></i>',
                type: 'wrap',
                markers: ['`', '`'],
            },
            {
                id: 'multicode',
                label: '<i class="fa-solid fa-terminal"></i>',
                type: 'wrap',
                markers: ['```', '```'],
            },
            {
                id: 'ulist',
                label: '<i class="fa-solid fa-list-ul"></i>',
                type: 'prefix',
                prefix: '- ',
            },
            {
                id: 'olist',
                label: '<i class="fa-solid fa-list-ol"></i>',
                type: 'prefix',
                prefix: '1. ',
            },
            {
                id: 'blockquote',
                label: '<i class="fa-solid fa-quote-right"></i>',
                type: 'prefix',
                prefix: '> ',
            },
        ],
        previewRules: [
            {
                regex: /```([\s\S]+?)```/g,
                replacement:
          '<pre class="bg-gray-200 rounded overflow-auto"><code>$1</code></pre>',
            },
            // ...other rules...
            { regex: /\n/g, replacement: '<br>' },
        ],
        placeholder: 'Type your message...',
    };

    // Merge options with defaults.
    const config = { ...defaultConfig, ...options };

    // Get the container element.
    const container = document.querySelector(`#${containerId}`);
    if (!container) {
        console.error(`Container with id "${containerId}" not found.`);
        return;
    }
    // Clear any existing content.
    container.innerHTML = '';

    // Create toolbar container.
    const toolbar = document.createElement('div');
    toolbar.id = 'toolbar';
    toolbar.className = 'flex w-full mb-2';
    container.append(toolbar);

    // Create the editor (contenteditable) with full width.
    const editor = document.createElement('div');
    editor.id = 'editor';
    editor.className = 'w-full border border-gray-300 p-2 rounded min-h-[100px]';
    editor.setAttribute('contenteditable', 'true');
    editor.setAttribute('placeholder', config.placeholder);
    container.append(editor);

    // Create the preview area.
    const previewContainer = document.createElement('div');
    previewContainer.className = 'mt-4';
    const previewTitle = document.createElement('h3');
    previewTitle.className = 'font-bold mb-1';
    previewTitle.textContent = 'Preview:';
    previewContainer.append(previewTitle);
    const preview = document.createElement('div');
    preview.id = 'preview';
    preview.className = 'w-full p-2 border border-gray-300 rounded min-h-[100px]';
    previewContainer.append(preview);
    container.append(previewContainer);

    // Helper Functions

    // For "wrap" formatting (e.g., bold, italic, code, etc.)
    function applyFormatting(markerStart, markerEnd) {
        editor.focus();
        const selection = globalThis.getSelection();
        if (!selection.rangeCount) return;
        const range = selection.getRangeAt(0);
        if (range.collapsed) {
            handleCollapsedRange(range, markerStart, markerEnd);
        } else {
            handleSelectionRange(range, markerStart, markerEnd);
        }
        updatePreview();
    }

    function handleCollapsedRange(range, start, end) {
        const textNode = document.createTextNode(start + end);
        range.insertNode(textNode);
        range.setStart(textNode, start.length);
        range.setEnd(textNode, start.length);
        globalThis.getSelection().removeAllRanges().addRange(range);
    }

    function handleSelectionRange(range, start, end) {
        const content = start + range.toString() + end;
        document.execCommand('insertText', false, content);
    }

    // For "prefix" formatting (e.g., lists, block quotes)
    function applyLinePrefix(prefix) {
        editor.focus();
        const selection = globalThis.getSelection();
        if (!selection.rangeCount) return;
        const range = selection.getRangeAt(0);
        if (range.collapsed) {
            document.execCommand('insertText', false, prefix);
        } else {
            const selectedText = range.toString();
            const lines = selectedText.split('\n');
            const prefixedLines = lines.map((line) => prefix + line);
            const newText = prefixedLines.join('\n');
            document.execCommand('insertText', false, newText);
        }
        updatePreview();
    }

    // Update preview by applying preview rules to the editor's plain text.
    function updatePreview() {
        const content = editor.textContent;
        let formattedContent = content;
        for (const rule of config.previewRules) {
            formattedContent = formattedContent.replace(rule.regex, rule.replacement);
        }
        preview.innerHTML = formattedContent;
    }

    // Initialize the toolbar with buttons.
    function initializeToolbar() {
        toolbar.innerHTML = config.toolbarFormats
            .map((format) => {
                if (format.type === 'wrap') {
                    return `
            <button
              id="btn-${format.id}"
              class="format-btn"
              data-start="${format.markers[0]}"
              data-end="${format.markers[1]}"
            >
              ${format.label}
            </button>
          `;
                } else if (format.type === 'prefix') {
                    return `
            <button
              id="btn-${format.id}"
              class="format-btn"
              data-prefix="${format.prefix}"
            >
              ${format.label}
            </button>
          `;
                }
                return '';
            })
            .join('');
        for (const format of config.toolbarFormats) {
            const button = toolbar.querySelector(`#btn-${format.id}`);
            if (!button) continue;
            if (format.type === 'wrap') {
                button.addEventListener('click', () => {
                    applyFormatting(...format.markers);
                });
            } else if (format.type === 'prefix') {
                button.addEventListener('click', () => {
                    applyLinePrefix(format.prefix);
                });
            }
        }
    }

    // Set up event listener to update preview on input.
    editor.addEventListener('input', updatePreview);

    // Initialize toolbar and preview.
    initializeToolbar();
    updatePreview();

    // Return an object for external control.
    return {
        getContent: function () {
            return editor.textContent;
        },
        getFormattedContent: function () {
            return preview.innerHTML;
        },
        editorElement: editor,
        previewElement: preview,
    };
}
