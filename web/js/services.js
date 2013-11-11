'use strict';

/* Services */
// Demonstrate how to register services
// In this case it is a simple value service.
angular.module('vaultee.services', []).
	value('shortversion', '0.5.0').
	value('longversion', '0.5.0-20130421.146').
	service('uuid', function() {
	    this.get = function() {
	        return 'xxxxxxxx'.replace(/[xy]/g, function (c) {
	            var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
	            return v.toString(16);
	        });
	    }		
	}).
	factory('socket', function() {
		var eb = {};

	    var init = function(protocol, hostname, port) {
			eb = new vertx.EventBus(protocol + '//' + hostname + ':' + port + '/eventbus', {});
			
			eb.onopen = function() {
		        console.log("eventbus is open");
		        $.publish('/eventbus/opened', {});
		    }

		   	eb.onclose = function() {
		        console.log("eventbus is closed");
		        eb = null;	    	
		    }
	    };

	    var send = function(message, params) {
	        console.log("message: " + message);
	        console.log("params: " + JSON.stringify(params));

	        var promise = new RSVP.Promise();

	        eb.send(message, $.extend({sessionID: $.cookie('vaultee_sessionid')}, params), handler);

	        function handler(reply) {
	        	console.log('this is what i got: '+JSON.stringify(reply));
	        	if (reply.status === 'ok') {
	        		promise.resolve(reply)
	        	} else {
	        		promise.reject(this);
	        	}
	        }

	        return promise;
		};

		var receive = function(message, callback) {
            eb.registerHandler(message, callback);
        };

		return {
			init: init,
			send: send,
			receive: receive
		};
	}).
	factory('defaults', function() {
		return {
			revision: {id: 0, assetId: 0, index: 0, created: 1},
			asset: {id: 0, name: "untitled", category: 0, categoryName: "", created: 1, modified: 1},
			item: {hash: "", product: {id: 0, name: "", asin: "", upc: "", itemType: {id: 0, name: ""}}, reference: "", quantity: null, price: null}
		}
	}).
	factory('model', function(defaults, uuid) {
		var User = function(json) {
			var self = this;

			self.name = json.name;
			self.email = json.email;
		}

		var Asset = function(json, defaults) {
			var self = this;

			self.id = json.id;
			self.hash = uuid.get();

			self.name = json.name;

			self.nameEditing = false;

			self.category = json.category;
			self.categoryName = json.categoryName;

			self.created = json.created
			self.modified = json.modified;

			self.revisions = [];
			self.revision = new Revision(defaults.revision);

			self.newItem = new Item(defaults.item);
		}

		var Item = function(json) {
			var self = this;

			self.hash = uuid.get();

			self.product = angular.copy(json.product);
			self.reference = json.reference;
			self.quantity = json.quantity;
			self.price = json.price;
		};

		var Revision = function(json) {
			var self = this;

			self.id = json.id;
			self.index = json.index;
			self.created = json.created;

			self.items = [];
		};

		var View = function(name, category, section) {
			var self = this;

			self.name = name;
			self.category = category;
			self.section = section;

			self.asset = null;
			self.newAsset = null;
		};

		var ItemType = function(json) {
			var self = this;

			self.id = json.id;
			self.name = json.name;
		}

		return {
			User: User,
			Asset: Asset,
			Item: Item,
			Revision: Revision,
			View: View,
			ItemType: ItemType
		}
	}).
	service('core', function($rootScope, model, socket, uuid, defaults) {
		var self = this;

		// Base initialization
		self.views = [
				new model.View("ALL", 0, 0),
				new model.View("WORKSTATION", 1, 1),
				new model.View("SERVER", 2, 1),
				new model.View("AUDIO/VIDEO", 3, 1),
				new model.View("OTHER", 4, 1)
			];		

		self.user = {};

		self.assets = [];
		self.itemTypes = [];

		self.view = {};

		// hashMap to hold objects that are waiting for feedback from the server
		self.cache = {};

		self.initialize = function() {
			var promise = new RSVP.Promise();

			socket.send('load:user', {}).
				then(function(user) {
					self.user = new model.User(user);
					return socket.send('load:itemTypes', {})
				}).
				then(function(itemTypes) {
					angular.forEach(itemTypes.results, function(value, key) {
						this.push(new model.ItemType(value));
					}, self.itemTypes);

					return socket.send('load:assets', {});
				}).
				then(function(assets) {
					angular.forEach(assets.results, function(value, key) {
						this.push(new model.Asset(value, defaults));
					}, self.assets);

					promise.resolve();
					// $.publish('/core/initialized');
					// $rootScope.$apply();
				});			

			return promise;
		}

		self.refreshAsset = function(asset) {
			var promise = new RSVP.Promise();

			if (asset.id === 0) {
				promise.reject();
				return promise;
			}

			socket.send('load:revisions', {aid: asset.id}).
				then(function(json) {
					asset.revisions = $.map(json.results, function(revision) {
						return new model.Revision(revision);
					});

					console.log('revisions = ' + JSON.stringify(asset.revisions));

					if (asset.revision.id === 0) {
						asset.revision = asset.revisions[0];
					}

					return socket.send('load:items', {aid: asset.id, rev: asset.revision.id});
				}).
				then(function(reply) {
					asset.revision.items = $.map(reply.results, function(item) {
						return new model.Item(item);
					});

					// $.publish('/model/updated');
					promise.resolve();
				});

			return promise;
		}

		// var findAssetByHash = function(hash) {
		// 	for (var i = 0; i < self.assets.length; i++) {
		// 		console.log('mother: '+self.assets[i].hash+' - '+hash+' :rana');
		// 		if (self.assets[i].hash === hash) {
		// 			return self.assets[i];
		// 		}
		// 	}

		// 	return null;
		// }

		// var findAssetById = function(id) {
		// 	for (var i = 0; i < self.assets.length; i++) {
		// 		if (self.assets[i].id === id) {
		// 			return self.assets[i];
		// 		}
		// 	}

		// 	return null;
		// }			

		self.saveAsset = function(asset) {
			var promise = new RSVP.Promise();

			self.cache[asset.hash] = asset;

			// socket.send('save:asset', {aid: self.view.asset.id, hash: self.view.asset.hash, name: self.view.asset.name, category: self.view.asset.category, items: self.view.asset.revision.items}).
			socket.send('save:asset', {aid: self.view.asset.id, hash: self.view.asset.hash, name: self.view.asset.name, category: self.view.asset.category, items: asset.revision.items}).
				then(function(reply) {
					console.log('saves successfully');
					console.log('brought back'+JSON.stringify(reply));
					console.log('still live'+JSON.stringify(asset));

					var asset = null;
					// if (reply.asset.hash === "") {
					// 	asset = findAssetById(reply.asset.id);
					// } else {
					// 	asset = findAssetByHash(reply.asset.hash);
					// }

					asset = self.cache[reply.asset.hash];

					if (asset === null) {
						promise.reject();
						return;
					}

					angular.extend(asset, reply.asset);
					return self.refreshAsset(asset);
				}).
				then(function() {
					delete self.cache[asset.hash];
					promise.resolve();
				});

			return promise;
		}

		self.scrapeItem = function(item) {
			var promise = new RSVP.Promise();
			
			self.cache[item.hash] = item;			

			socket.send('scrape:item', item).
				then(function(data) {
					item = self.cache[data.hash];

					console.log('item is '+JSON.stringify(item));
					console.log('data is '+JSON.stringify(data));
					promise.resolve(angular.extend(item.product, data.product));

					delete self.cache[item.hash];
				})

			return promise;
		}

		self.loadItems = function(asset) {
			var promise = new RSVP.Promise();

			socket.send('load:items', {aid: asset.id, rev: asset.revision.id}).
				then(function(json) {
					asset.revision.items = $.map(json.results, function(item) {
						return new model.Item(item);
					});

					// $.publish('/model/updated');
					promise.resolve();				
				})

			return promise;
		}

		self.addItem = function(item) {
			self.view.asset.revision.items.push(item);
			self.view.asset.newItem = new model.Item(defaults.item);

			// socket.send('scrape:item', item).
			// 	then(function(json) {
			// 		console.log('but you know what, why are those characteristics');
			// 	})
		}
	});
