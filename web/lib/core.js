var vee = vee || {};

vee.Core = (function(url, options) {

    var that = this;
    var eb = new vertx.EventBus(url, options);

    eb.onopen = function() {
        console.log("eventbus is open");
        $.publish('/eventbus/opened', {});
    }

    eb.onclose = function() {
        console.log("eventbus is closed");
        eb = null;
    };    

    return {
        publish: function(message, params, callback) {
            console.log("message: " + message);
            console.log("params: " + params);
            eb.send(message, $.extend({sessionID: $.cookie('vaultee_sessionid')}, params), callback);        
        },

        subscribe: function(message, callback) {
            eb.registerHandler(message, callback);
            return that;
        }
    }
}(window.location.protocol + '//' + window.location.hostname + ':' + window.location.port + '/eventbus', {}));
