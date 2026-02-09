function updateCostIndicators(whatsbackCost, twilioCost, watiCost) {
    const cards = [
        {
            id: 'whatsbackCard',
            cost: whatsbackCost,
        },
        {
            id: 'twilioCard',
            cost: twilioCost,
        },
        {
            id: 'watiCard',
            cost: watiCost,
        },
    ];

    const sorted = [...cards].sort((a, b) => a.cost - b.cost);

    const colorMapping = {};
    if (sorted.length === 3) {
        colorMapping[sorted[0].id] = 'border-green-500';
        colorMapping[sorted[1].id] = 'border-yellow-500';
        colorMapping[sorted[2].id] = 'border-red-500';
    }

    cards.forEach((card) => {
        const el = document.getElementById(card.id);
        el.classList.remove(
            'border-green-500',
            'border-yellow-500',
            'border-red-500'
        );

        el.classList.add('border-4');
        el.classList.add(colorMapping[card.id]);
    });
}

function calculateCosts() {
    const totalMessages =
    parseInt(document.getElementById('messageInput').value, 10) || 0;

    document.getElementById('whatsbackMessages').textContent =
    totalMessages.toLocaleString();
    document.getElementById('whatsbackTotalCost').textContent = '$0.00';
    const whatsbackCost = 0.0;

    document.getElementById('twilioMessages').textContent =
    totalMessages.toLocaleString();
    const twilioCost = totalMessages * 0.005;
    document.getElementById('twilioTotalCost').textContent =
    '$' + twilioCost.toFixed(2);

    document.getElementById('watiMessages').textContent =
    totalMessages.toLocaleString();
    const baseCostWATI = 49;
    let additionalMessages = totalMessages - 1000;
    additionalMessages = additionalMessages > 0 ? additionalMessages : 0;
    const additionalCostWATI = additionalMessages * 0.04;
    const watiCost = baseCostWATI + additionalCostWATI;
    document.getElementById('watiTotalCost').textContent =
    '$' + watiCost.toFixed(2);

    updateCostIndicators(whatsbackCost, twilioCost, watiCost);
}

calculateCosts();

document
    .getElementById('messageInput')
    .addEventListener('input', calculateCosts);
