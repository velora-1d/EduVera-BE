/**
 * Validates an individual cron field token.
 * @param {string} field - The full field string (which may be a comma-separated list)
 * @param {number} min - Minimum allowed numeric value.
 * @param {number} max - Maximum allowed numeric value.
 * @param {object} options - Options for special tokens.
 *                           { allowL: boolean, allowHash: boolean }
 * @returns {boolean} - True if valid.
 */
function isValidCronField(field, min, max, options = {}) {
    const { allowL = false, allowHash = false } = options;

    // Split the field by commas to allow lists
    const parts = field.split(',');
    for (let part of parts) {
        part = part.trim();

        // Allow a single asterisk.
        if (part === '*') continue;

        // Allow "*/step" pattern.
        let match = part.match(/^\*\/(\d+)$/);
        if (match) {
            const step = parseInt(match[1], 10);
            if (step < 1 || step > max) return false;
            continue;
        }

        // Allow a plain number.
        if (/^\d+$/.test(part)) {
            const num = parseInt(part, 10);
            if (num < min || num > max) return false;
            continue;
        }

        // Allow a number with an "L" suffix (e.g. "6L") if allowed.
        if (allowL && /^(\d+)L$/.test(part)) {
            const num = parseInt(part, 10);
            if (num < min || num > max) return false;
            continue;
        }

        // Allow a standalone "L" if allowed.
        if (allowL && part === 'L') continue;

        // Allow ranges, with an optional "/step": e.g. "10-30" or "10-30/2"
        match = part.match(/^(\d+)-(\d+)(\/(\d+))?$/);
        if (match) {
            const start = parseInt(match[1], 10);
            const end = parseInt(match[2], 10);
            if (start > end || start < min || end > max) return false;
            if (match[3]) {
                // if there is a step value
                const step = parseInt(match[4], 10);
                if (step < 1 || step > max) return false;
            }
            continue;
        }

        // Allow hash syntax for day-of-week (e.g. "6#3") if allowed.
        if (allowHash) {
            match = part.match(/^(\d+)#(\d+)$/);
            if (match) {
                const num = parseInt(match[1], 10);
                const nth = parseInt(match[2], 10);
                if (num < min || num > max) return false;
                // Typically, the nth occurrence is between 1 and 5.
                if (nth < 1 || nth > 5) return false;
                continue;
            }
        }

        // If nothing matches, the token is invalid.
        return false;
    }
    return true;
}

/**
 * Validates a cron expression (either 5 or 6 fields).
 * For 5 fields (standard): Minute Hour Day-of-Month Month Day-of-Week
 * For 6 fields (Quartz): Second Minute Hour Day-of-Month Month Day-of-Week
 * Special tokens ("L", "#" etc.) are allowed in day-of-month and day-of-week.
 *
 * @param {string} cronExpr - The cron expression.
 * @returns {object} - An object with { valid: boolean, message: string }
 */
function validateCronExpression(cronExpr) {
    const fields = cronExpr.trim().split(/\s+/);
    if (!(fields.length === 5 || fields.length === 6)) {
        return {
            valid: false,
            message: `❌ Invalid field count (${fields.length}). Expected 5 or 6 fields.`,
        };
    }

    // Define validators for each field depending on the number of fields.
    let validators;
    if (fields.length === 5) {
    // Standard cron: Minute, Hour, Day-of-Month, Month, Day-of-Week
        validators = [
            { name: 'Minute', min: 0, max: 59, options: {} },
            { name: 'Hour', min: 0, max: 23, options: {} },
            { name: 'Day of Month', min: 1, max: 31, options: { allowL: true } },
            { name: 'Month', min: 1, max: 12, options: {} },
            {
                name: 'Day of Week',
                min: 0,
                max: 7,
                options: { allowL: true, allowHash: true },
            },
        ];
    } else {
    // Quartz cron (6 fields): Second, Minute, Hour, Day-of-Month, Month, Day-of-Week
        validators = [
            { name: 'Second', min: 0, max: 59, options: {} },
            { name: 'Minute', min: 0, max: 59, options: {} },
            { name: 'Hour', min: 0, max: 23, options: {} },
            { name: 'Day of Month', min: 1, max: 31, options: { allowL: true } },
            { name: 'Month', min: 1, max: 12, options: {} },
            {
                name: 'Day of Week',
                min: 0,
                max: 7,
                options: { allowL: true, allowHash: true },
            },
        ];
    }

    // Validate each field.
    for (let i = 0; i < fields.length; i++) {
        if (
            !isValidCronField(
                fields[i],
                validators[i].min,
                validators[i].max,
                validators[i].options
            )
        ) {
            return {
                valid: false,
                message: `❌ Invalid ${validators[i].name} field: "${fields[i]}".`,
            };
        }
    }
    return { valid: true, message: '✅ Cron expression is valid.' };
}

// ---------------------
// Example Expressions
// ---------------------
// const expressions = [
//   "* * * * * *", // 6-field, Quartz
//   "0 * * * *", // 5-field standard
//   "*/15 * * * *", // 5-field standard with step
//   "0 9 * * *", // 5-field standard
//   "0 0 * * 0", // 5-field standard
//   "30 18 * * 1-5", // 5-field with range in day-of-week
//   "0 7 1 * *", // 5-field standard
//   "0 12 * * 6#3", // 5-field with nth weekday (day-of-week)
//   "0 0 31 12 *", // 5-field standard
//   "0 * * 1,4-10,L * *", // 6-field: day-of-month allows list with L
//   "10-30/2 2 12 8 0", // 5-field: minute has range with step
//   "0 0 0 * * 4,6L", // 6-field: day-of-week list with L suffix
//   "0 12 */5 6 *", // 5-field: day-of-month with step on wildcard
//   "0 15 */5 5 *", // 5-field standard with step
//   "0 0 6-20/2,L 2 *", // 5-field: day-of-month is a list: range with step and L
//   "0 0 0 * * 1L,5L", // 6-field: day-of-week with L suffixes in a list
//   "10 2 12 8 7", // 5-field: day-of-week can be 7
//   "0 8-18/2 * * 1-5", // 5-field: hour field with range/step
//   "0 9 15 * *", // 5-field standard
//   "0 12 1-7 * 1", // 5-field standard with range in day-of-month
//   "0 0 L * *", // 5-field: day-of-month as "L"
// ];

// expressions.forEach((expr) => {
//   const result = validateCronExpression(expr);
//   console.log(`Expression: "${expr}" ➡ ${result.message}`);
// });
