module.exports = {
    io: undefined,
    /**
   * Set the Socket.IO server instance to use for emitting events to connected clients.
   * @param {Server} ioInstance - The Socket.IO server instance.
   */
    setIo: function (ioInstance) {
        this.io = ioInstance;
    },
    /**
   * Emits a Socket.IO event to all connected clients.
   * @param {string} event - The name of the event to emit.
   * @param {Object} data - Optional data to send with the event.
   */
    socketEmit: function (event, data) {
        if (this.io) {
            this.io.emit(event, data);
        }
    },
};
