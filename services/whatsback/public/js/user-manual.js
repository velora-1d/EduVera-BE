document.querySelectorAll('.section-link').forEach((link) => {
    link.addEventListener('click', function (e) {
        e.preventDefault();
        const targetId = this.getAttribute('href');
        document.querySelector(targetId).scrollIntoView({
            behavior: 'smooth',
        });
    });
});

/**
 * Updates the active section link based on the current scroll position.
 * This function iterates through all content sections and checks if a section
 * is within the viewport. The corresponding link for the visible section is
 * highlighted by adding the 'active-section' class, while removing it from others.
 */
function updateActiveSection() {
    const sections = document.querySelectorAll('.content-section');
    const links = document.querySelectorAll('.section-link');

    sections.forEach((section, index) => {
        const rect = section.getBoundingClientRect();
        if (rect.top <= 100 && rect.bottom >= 100) {
            links.forEach((link) => link.classList.remove('active-section'));
            links[index].classList.add('active-section');
        }
    });
}

window.addEventListener('scroll', () => {
    updateActiveSection();
});

updateActiveSection();

const messageRange = document.querySelector('#messageRange');

/**
 * Updates the displayed cost values for Twilio and WATI based on the number of messages.
 * - Calculates Twilio cost as $0.005 per message.
 * - Calculates WATI cost with a base of $49, adding $0.04 per message over 1,000.
 * - Updates the savings display with Twilio's cost.
 */
const updateCosts = () => {
    const messages = parseInt(messageRange.value);

    const twilioCost = messages * 0.005;
    document.querySelector('#twilioCost').textContent = twilioCost.toFixed(2);

    const watiCost = messages > 1000 ? 49 + (messages - 1000) * 0.04 : 49;
    document.querySelector('#watiCost').textContent = watiCost.toFixed(2);

    document.querySelector('#savings').textContent = twilioCost.toFixed(2);
};

updateCosts();

messageRange.addEventListener('input', updateCosts);

document.querySelectorAll('#interface-guide img').forEach((img) => {
    img.style.cursor = 'zoom-in';
    img.addEventListener('click', function () {
        const overlay = document.createElement('div');
        overlay.className = 'image-overlay';

        const imgContainer = document.createElement('div');
        imgContainer.className =
      'relative h-full w-full flex justify-center items-center';

        const clonedImg = this.cloneNode();
        clonedImg.className = 'zoomed-image';

        const closeBtn = document.createElement('button');
        closeBtn.className = 'absolute top-4 right-4 text-white text-2xl z-50';
        closeBtn.innerHTML = '&times;';

        imgContainer.appendChild(clonedImg);
        imgContainer.appendChild(closeBtn);
        overlay.appendChild(imgContainer);
        document.body.appendChild(overlay);

        setTimeout(() => {
            overlay.classList.add('active');
            clonedImg.classList.add('active');
        }, 10);

        /**
     * Closes the image overlay by removing the active class and then removing
     * the overlay element after a short delay.
     */
        const closeModal = () => {
            overlay.classList.remove('active');
            clonedImg.classList.remove('active');
            setTimeout(() => overlay.remove(), 300);
        };

        overlay.addEventListener('click', (e) => {
            if (e.target === overlay || e.target === closeBtn) closeModal();
        });

        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') closeModal();
        });
    });
});
