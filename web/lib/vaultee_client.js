
function User(json) {
	var self = this;
	self.name = json.name;
	self.email = json.email;
};

function Asset(json) {
	var self = this;

	self.id = json.id;
	self.name = json.name;
	self.category = json.category;
	self.categoryName = json.categoryName;
	self.created = json.created;
	self.createdShort = moment(json.created, "YYYY-MM-DD HH:mmZ").format("MMM DD, YYYY");
	self.createdLong = moment(json.created, "YYYY-MM-DD HH:mmZ").format("MMM DD, YYYY HH:mm");
	self.modified = json.modified;
	self.modifiedShort = moment(json.modified, "YYYY-MM-DD HH:mmZ").format("MMM DD, YYYY");
	self.modifiedLong = moment(json.modified, "YYYY-MM-DD HH:mmZ").format("MMM DD, YYYY HH:mm");

	self.revisions = ko.observableArray([]);
	self.currentRevision = ko.observable(defaultRevision);

	self.items = ko.observableArray([]);
	// for (var i = 0; i < json.revisions.length; i++) {
	// 	self.revisions[i] = new Revision(json.revisions[i]);
	// }


	self.refreshData = function() {
		// if revision != self.currentRevision
		// self.currentRevision(revision);

		console.log("asset-id = " + self.id);

		if (self.id == -1) {
			return;
		}

		vee.Core.publish("vc.getrevisions", {"aid": self.id}, function(reply) {
			if (reply.status == 'ok') {
				console.log("got back ok status =", reply);

				self.revisions(ko.utils.arrayMap(reply.results, function(item) {
					return new Revision(item);
				}));

				if (self.currentRevision() == defaultRevision) {
					self.currentRevision(self.revisions()[0]);
				}
				
				self.loadItems();
				// $.publish('/revisions/loaded', {"asset": asset, "lastRev": revisions[0]});
				// $.publish('/revisions/loaded', []);
			}
		});			
	}

	self.loadItems = function() {
		vee.Core.publish("vc.getitems", {"aid": self.id, "rev": self.currentRevision().id}, function(reply) {
			if (reply.status === 'ok') {
				var items = [];
				for (var i = 0; i < reply.results.length; i++) {
					items[i] = new Item(reply.results[i]);
				}
				self.items(items);

				console.log('got items');
			}
		});			
	}
};

function Item(json) {
	var self = this;

	self.type = json.type;
	self.name = json.name;
	self.quantity = json.quantity;
	self.price = json.price;
};

function Revision(json) {
	var self = this;
	self.id = json.id;
	self.index = json.index;
	self.created = json.created;
	self.createdLong = ko.computed( function() {
		return moment(self.created, "YYYY-MM-DD HH:mmZ").format("MMM DD, YYYY HH:mm")
	});
};

var defaultRevision = new Revision({id: 0, assetId: 0, index: 0, created: "1970-01-01 05:00:00-05"});
var defaultAsset = new Asset({id: 0, name: "", category: 0, categoryName: "", created: "1970-01-01 05:00:00-05", modified: "1970-01-01 05:00:00-05"});

function View(name, category, section) {
	var self = this;

	self.name = name;
	self.category = category;
	self.section = section;

	self.currentAsset = ko.observable(new Asset(defaultAsset));
};

(function VaulteeViewModel() {
	var self = this;

	self.views = ko.observableArray([
			new View("ALL", 0, 0),
			new View("WORKSTATION", 1, 1),
			new View("SERVER", 2, 1),
			new View("AUDIO/VIDEO", 3, 1),
			new View("OTHER", 4, 1)
		]);

	// self.views = [ {name: "ALL", category: 0, section: 0, asset: new Asset(defaultAsset)}, {name: "WORKSTATION", category: 1, section: 1}, {name: "SERVER", category: 2, section: 1}, {name: "AUDIO/VIDEO", category: 3, section: 1}, {name: "OTHER", category: 4, section: 1}]
	// self.addView = new View("ADD ASSET", -1, 2);

	self.currentUser = ko.observable({});

	self.assets = ko.observableArray([]);

	self.filteredAssets = ko.computed(function() {
		return self.assets().filter(function(asset) {
			return asset.category == self.currentView().category;
		})
	});

	self.currentView = ko.observable({});
	// self.currentCategory = ko.observable(0);
	// self.currentSection = ko.observable(0);

	// self.currentAsset = ko.observable(defaultAsset);

	// self.currentRevision = ko.observable({});
	// self.selectedRevision = ko.observable({});
	// self.items = ko.observableArray([]);

	var onEventBusOpened = function() {
		console.log("inside onEventBusOpened");

		self.loadUserData();
		self.loadAssets();

		self.switchToView(self.views()[0]);

		ko.applyBindings(self);
		// self.selectCategory(categories[0]);
		// // self.getAsset();
	}

	$.subscribe('/eventbus/opened', onEventBusOpened);

	self.loadUserData = function() {
		vee.Core.publish("vc.getuserdata", {}, function(reply) {
			// console.log("got here data=" + reply.name);
			self.currentUser(new User(reply));
			// self.username(reply.name);
			// self.email(reply.email);
			console.log("user.name=" + self.currentUser().name);

		});
	};

	self.loadAssets = function() {
		vee.Core.publish("vc.getassets", {}, function(reply) {
			console.log("these are the days = ", reply.results);
	        if (reply.status === 'ok') {
	          var assetArray = [];
	          for (var i = 0; i < reply.results.length; i++) {
	            assetArray[i] = new Asset(reply.results[i]);
	          }
	          self.assets(assetArray);
	        } else {
	          console.error('Failed to retrieve assets: ' + reply.message);
	        }			
		});
	};

	self.switchToView = function(view) {
		console.log("everybody knows it");
		console.log("view="+view);
		self.currentView(view);
		// self.currentCategory(view.category);
		// self.currentSection(view.section);
	};

	self.switchToAsset = function(asset) {
		console.log("asset is " + asset.name);

		if (self.currentView().category != asset.category) {
			console.log("inside !=");
			var view = ko.utils.arrayFirst(self.views(), function(view) {
				console.log("view.cat = " + view.category);
				console.log("asset.cat = " + asset.category);
				return view.category == asset.category;
			});

			self.switchToView(view);
		};

		self.currentView().currentAsset(asset);

		asset.refreshData();
	};

	self.switchToAddAsset = function() {
		console.log("switched to add asset");
		var asset = new Asset({id: -1, name: "untitled", category: self.currentView().category, categoryName: self.currentView().name, created: "1970-01-01 05:00:00-05", modified: "1970-01-01 05:00:00-05"});
		self.assets.push(asset);
		self.switchToAsset(asset);
	};

	// self.getRevisions = function(asset) {
	// 	console.log("about to get revisions for = ", asset.name);
	// 	vee.Core.publish("vc.getrevisions", {"aid": asset.id}, function(reply) {
	// 		if (reply.status == 'ok') {
	// 			console.log("got back ok status =", reply);
	// 			console.log("asset.id =", asset.id);
	// 			var revisions = [];
	// 			for (var i = 0; i < reply.results.length; i++) {
	// 				revisions[i] = new Revision(reply.results[i]);
	// 			}
	// 			asset.revisions(revisions);

	// 			console.log("about to publish /revisions/loaded");
	// 			$.publish('/revisions/loaded', [asset, revisions[0]]);
	// 			// $.publish('/revisions/loaded', {"asset": asset, "lastRev": revisions[0]});
	// 			// $.publish('/revisions/loaded', []);
	// 		}
	// 	});
	// };

	self.getAsset = function() {
		console.log("calling getAsset");
		vee.Core.publish("vc.getasset", {"aid": 1}, function(reply) {
			if (reply.status === 'ok') {
				console.log('el mero mero');
			}
		});
	};
})();
