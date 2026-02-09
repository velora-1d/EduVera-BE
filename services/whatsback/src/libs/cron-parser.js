const { CronExpressionParser } = require("cron-parser");

class Cron {
    /**
   * Constructs a new Cron instance with a specified cron expression and options.
   *
   * @param {string} exp - The cron expression to be parsed.
   * @param {Object} [options={}] - Optional settings for the cron instance.
   * @param {string} [options.tz="UTC"] - The timezone for parsing the cron expression.
   * @throws {Error} Throws an error if no cron expression is provided.
   */
    constructor(exp, options = {}) {
        if (!exp) {
            throw new Error("Cron expression is required");
        }
        this.exp = exp;
        this.options = options;
        if (!this.options.tz) {
            this.options.tz = "UTC";
        }
        this.interval = CronExpressionParser.parse(exp, { tz: this.options.tz });
    }

    /**
   * Returns a human-readable description of a cron expression.
   *
   * @returns {string} A human-readable description of the cron expression.
   * @example
   * const cron = new Cron("0 0 * * *");
   * cron.translate(); // "at 12:00 AM every day"
   */
    translate() {
        const expression = this.exp;
        const parts = expression.trim().split(/\s+/);
        const hasSeconds = parts.length === 6;
        if (parts.length !== 5 && parts.length !== 6) {
            throw new Error("Cron expression must have 5 or 6 fields.");
        }

        const second = hasSeconds ? parts[0] : "0";
        const minute = hasSeconds ? parts[1] : parts[0];
        const hour = hasSeconds ? parts[2] : parts[1];
        const dayOfMonth = hasSeconds ? parts[3] : parts[2];
        const month = hasSeconds ? parts[4] : parts[3];
        const dayOfWeek = hasSeconds ? parts[5] : parts[4];

        const monthNames = [
            "January",
            "February",
            "March",
            "April",
            "May",
            "June",
            "July",
            "August",
            "September",
            "October",
            "November",
            "December",
        ];
        const dayNames = [
            "Sunday",
            "Monday",
            "Tuesday",
            "Wednesday",
            "Thursday",
            "Friday",
            "Saturday",
        ];

        /**
     * Describes a cron expression field in a human-readable format.
     * @param {string} field - The field to describe.
     * @param {string} singularLabel - The singular label for the field.
     * @param {string} pluralLabel - The plural label for the field.
     * @param {function} [valuesToName] - A function to map values to a name.
     * @returns {string} - A human-readable description of the field.
     */
        function describeField(
            field,
            singularLabel,
            pluralLabel,
            valuesToName = null
        ) {
            if (field === "*" || field === "?") {
                return `every ${singularLabel}`;
            }
            if (field.includes("*/")) {
                return `every ${field.split("*/")[1]} ${pluralLabel}`;
            }
            if (field.includes(",")) {
                const list = field
                    .split(",")
                    .map((v) => (valuesToName ? valuesToName(v) : v));
                return `at ${singularLabel}s ${list.join(", ")}`;
            }
            if (field.includes("-") && field.includes("/")) {
                const [range, step] = field.split("/");
                const [start, end] = range.split("-");
                return `every ${step} ${pluralLabel} from ${start} to ${end}`;
            }
            if (field.includes("-")) {
                const [start, end] = field.split("-");
                const s = valuesToName ? valuesToName(start) : start;
                const e = valuesToName ? valuesToName(end) : end;
                return `every ${singularLabel} from ${s} to ${e}`;
            }
            return `at ${singularLabel} ${
                valuesToName ? valuesToName(field) : field
            }`;
        }

        const secondDesc =
      second === "*" || second === "?"
          ? "every second"
          : second.includes("*/")
              ? `every ${second.split("*/")[1]} seconds`
              : `at second ${second}`;
        const minuteDesc = describeField(minute, "minute", "minutes");
        const hourDesc = describeField(hour, "hour", "hours");

        /**
     * Returns a human-readable description of a day of the month field in a
     * cron expression.
     *
     * @param {string} field - The field to describe.
     * @returns {string} A human-readable description of the field.
     */
        function describeDayOfMonth(field) {
            if (field === "*" || field === "?") {
                return "every day";
            }
            if (field.includes(",")) {
                const items = field.split(",").map(describeDayOfMonthItem);
                return `on the ${items.join(" and ")} day(s) of the month`;
            }
            return `on the ${describeDayOfMonthItem(field)} day of the month`;
        }

        /**
     * Describes a day of the month field item in a cron expression.
     *
     * - "L" indicates the last day of the month.
     * - A range with a step (e.g., "1-5/3") means "every N days from start to end".
     * - A range (e.g., "1-5") means "days from start through end".
     * - A step (e.g., "'*'/3") indicates "every N days".
     * - A specific day is returned with its ordinal suffix.
     *
     * @param {string} item - The day of the month field item.
     * @returns {string} A description of the item.
     */
        function describeDayOfMonthItem(item) {
            if (item === "L") {
                return "last";
            }
            if (item.includes("/") && item.includes("-")) {
                const [range, step] = item.split("/");
                const [start, end] = range.split("-");
                return `every ${step} days from ${start} to ${end}`;
            }
            if (item.includes("-")) {
                const [start, end] = item.split("-");
                return `days ${start} through ${end}`;
            }
            if (item.includes("/")) {
                const [base, step] = item.split("/");
                if (base === "*") {
                    return `every ${step} days`;
                }
                return `every ${step} days starting at ${base}`;
            }
            return `${item}${getOrdinal(item)}`;
        }

        /**
     * Given a cron expression's month field, returns a human-readable description
     * of that field.
     *
     * @param {string} field - The month field from a cron expression.
     * @returns {string} A human-readable description of the month field.
     *
     * @example
     * describeMonth("*"); // "every month"
     * describeMonth("1,3,5"); // "in January, March, and May"
     * describeMonth("2-4"); // "from February through April"
     * describeMonth("AUG"); // "in August"
     */
        function describeMonth(field) {
            if (field === "*" || field === "?") {
                return "every month";
            }
            if (field.includes(",")) {
                const items = field
                    .split(",")
                    .map((m) => (isNaN(m) ? m : monthNames[parseInt(m, 10) - 1]));
                return `in ${items.join(" and ")}`;
            }
            if (field.includes("-")) {
                const [start, end] = field.split("-");
                const s = isNaN(start) ? start : monthNames[parseInt(start, 10) - 1];
                const e = isNaN(end) ? end : monthNames[parseInt(end, 10) - 1];
                return `from ${s} through ${e}`;
            }
            return `in ${isNaN(field) ? field : monthNames[parseInt(field, 10) - 1]}`;
        }

        /**
     * Translate a cron-style day-of-week field into a human-readable string.
     * @param {string} field - Cron-style day-of-week field.
     * @returns {string} - Human-readable string.
     */
        function describeDayOfWeek(field) {
            if (field === "*" || field === "?") {
                return "every day of the week";
            }
            if (field.includes(",")) {
                const items = field.split(",").map(describeDayOfWeekItem);
                return `on ${items.join(" and ")}`;
            }
            if (field.includes("-")) {
                const [start, end] = field.split("-");
                return `every day from ${describeDayOfWeekItem(
                    start
                )} to ${describeDayOfWeekItem(end)}`;
            }
            return `on ${describeDayOfWeekItem(field)}`;
        }

        /**
     * Describes a day of the week field item. If the item has an
     * "L" suffix, it is treated as "last". If the item has a "#"
     * character, it is treated as "every nth". Otherwise, it is
     * treated as a number 1-7. The description is as follows:
     *
     *  - last: the last <dayName>
     *  - nth: the <nth> <dayName>
     *  - num: <dayName>
     *
     * @param {string} item - the day of the week field item
     * @returns {string} a description of the item
     */
        function describeDayOfWeekItem(item) {
            if (item.endsWith("L")) {
                const num = item.slice(0, -1);
                const dayName = dayNames[parseInt(num, 10) % 7];
                return `the last ${dayName}`;
            }
            if (item.includes("#")) {
                const [d, nth] = item.split("#");
                const dayName = dayNames[parseInt(d, 10) % 7];
                return `the ${nth}${getOrdinal(nth)} ${dayName}`;
            }
            let num = parseInt(item, 10);
            return dayNames[num % 7];
        }

        const dayOfMonthDesc = describeDayOfMonth(dayOfMonth);
        const monthDesc = describeMonth(month);
        const dayOfWeekDesc = describeDayOfWeek(dayOfWeek);

        const partsDesc = [];
        if (hasSeconds) partsDesc.push(secondDesc);
        partsDesc.push(
            minuteDesc,
            hourDesc,
            dayOfMonthDesc,
            monthDesc,
            dayOfWeekDesc
        );

        return partsDesc.join(", ");

        /**
     * Retrieves the ordinal suffix for a given number.
     *
     * @param {number|string} n - The number to retrieve the ordinal suffix for.
     * @returns {string} The ordinal suffix for `n`.
     * @example
     * getOrdinal(1) // "st"
     * getOrdinal(2) // "nd"
     * getOrdinal(3) // "rd"
     * getOrdinal(4) // "th"
     * getOrdinal(10) // "th"
     * getOrdinal(11) // "th"
     * getOrdinal(12) // "th"
     * getOrdinal(13) // "th"
     * getOrdinal(14) // "th"
     */
        function getOrdinal(n) {
            n = parseInt(n, 10);
            const s = ["th", "st", "nd", "rd"];
            const v = n % 100;
            return s[(v - 20) % 10] || s[v] || s[0];
        }
    }

    /**
   * Retrieves the next scheduled run times for the cron expression.
   *
   * @param {number} count - The number of future run dates to retrieve (default is 5).
   * @returns {string[]} An array of strings representing the next `count` scheduled run dates.
   */
    getNextRuns(count = 5) {
        return this.interval.take(count).map((date) => date.toString());
    }

    /**
   * Return an object with a human-readable description of the schedule and an
   * array of the next `count` runs of the schedule.
   *
   * @return {Object}
   * @property {string} description - A human-readable description of the schedule.
   * @property {Array<string>} nextRuns - Array of the next `count` runs of the schedule.
   */
    get schedule() {
        return {
            description: this.translate(),
            nextRuns: this.getNextRuns(),
        };
    }
}

module.exports = Cron;
