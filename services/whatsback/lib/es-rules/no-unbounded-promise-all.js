/**
 * @fileoverview Disallow unbounded use of Promise.all on array.map
 */

"use strict";

module.exports = {
    meta: {
        type: "problem",
        docs: {
            description:
        "Warn when using Promise.all over an array.map without limiting concurrency",
            category: "Performance",
            recommended: false,
        },
        schema: [], // no options
        messages: {
            unboundedPromiseAll:
        "Using Promise.all on an array.map without limiting concurrency can be expensive when processing many items.",
        },
    },

    create(context) {
        return {
            CallExpression(node) {
                // Check if the call is Promise.all(...)
                if (
                    node.callee &&
          node.callee.type === "MemberExpression" &&
          node.callee.object &&
          node.callee.object.name === "Promise" &&
          node.callee.property &&
          node.callee.property.name === "all"
         && // Ensure there's exactly one argument passed to Promise.all
          node.arguments.length === 1) {
                    const argument = node.arguments[0];
                    // Check if the argument is a CallExpression and its callee is a MemberExpression ending in 'map'
                    if (
                        argument.type === "CallExpression" &&
              argument.callee &&
              argument.callee.type === "MemberExpression" &&
              argument.callee.property &&
              argument.callee.property.name === "map"
                    ) {
                        context.report({
                            node,
                            messageId: "unboundedPromiseAll",
                        });
                    }
                }
            },
        };
    },
};
