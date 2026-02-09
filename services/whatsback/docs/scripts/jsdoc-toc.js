(function($) {
    // TODO: make the node ID configurable
    var treeNode = $('#jsdoc-toc-nav');

    // initialize the tree
    treeNode.tree({
        autoEscape: false,
        closedIcon: '&#x21e2;',
        data: [{"label":"<a href=\"global.html\">Globals</a>","id":"global","children":[]},{"label":"<a href=\"Cron.html\">Cron</a>","id":"Cron","children":[]},{"label":"<a href=\"module-WhatsAppClient.html\">WhatsAppClient</a>","id":"module:WhatsAppClient","children":[]},{"label":"<a href=\"module-state.html\">state</a>","id":"module:state","children":[]}],
        openedIcon: ' &#x21e3;',
        saveState: false,
        useContextMenu: false
    });

    // add event handlers
    // TODO
})(jQuery);
