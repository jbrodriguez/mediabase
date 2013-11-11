'use strict';

// Declare app level module which depends on filters, and services
angular.module('vaultee', ['vaultee.filters', 'vaultee.services', 'vaultee.directives']).
	run(function(socket) {
		socket.init(window.location.protocol, window.location.hostname, window.location.port );
	});
