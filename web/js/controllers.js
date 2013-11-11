'use strict';

/* Controllers */
function vaulteeController($scope, core, model, uuid, defaults) {

	// Bind to the core service, that manages the model for the app
	$scope.core = core;

	// $scope.assets = [];

	// $scope.view = null;

	// $scope.types = [];

	// Event handlers and corresponding subscriptions
	var onEventBusOpened = function() {
		console.log("inside onEventBusOpened");

		$scope.core.initialize().
			then(function() {
				$scope.switchToView($scope.core.views[0]);
				$scope.$apply();
			});
	}

	var onCoreUpdated = function() {
		$scope.$apply();
	}

	$.subscribe('/eventbus/opened', onEventBusOpened);
	$.subscribe('/core/updated', onCoreUpdated);


	// // Internal Functions

	// $scope.loadItemTypes = function() {
	// 	core.publish("vc.getcategories", {}, function(reply) {
	// 		if (reply.status === 'ok') {
	//           var typesArray = [];
	//           for (var i = 0; i < reply.results.length; i++) {
	//             typesArray[i] = new ItemType(reply.results[i]);
	//           };

	//           $scope.types = typesArray;
	//           $scope.$apply();
	//       	} else {
	//       		console.error('Failed to retrieve itemtypes: ' + reply.message);
	//       	}
	// 	});
	// };

	// // Functions called from html land

	$scope.viewFirstAssetIndex = function(view) {
		for (var i = 0; i < $scope.core.assets.length; i++) {
			if (view.category === $scope.core.assets[i].category) {
				return i;
			}
		};

		return -1;
	}

	$scope.switchToView = function(view) {
		$scope.core.view = view;

		if ($scope.core.view.asset === null) {
			var index = $scope.viewFirstAssetIndex($scope.core.view)

			console.log('this is the index '+index);

			if (index > -1) {
				// load from local storage the saved asset id
				// search the id in the array ... 
				// didnt find the id, give index id the focus
				console.log('this is the end of innocence');

				// $scope.view.asset = $scope.assets[index];
				$scope.switchToAsset($scope.core.assets[index]);
			}
		}
	};

	$scope.switchToAsset = function(asset) {
		if ($scope.core.view.category !== asset.category) {
			for (var i = 0; i < $scope.core.views.length; i++) {
				if ($scope.core.views[i].category === asset.category) {
					break;
				}
			};

			// $scope.switchToView($scope.views[i]);
			$scope.core.view = $scope.core.views[i];
		};

		console.log('switchtoasset='+JSON.stringify(asset));
		$scope.core.view.asset = asset;

		$scope.core.refreshAsset($scope.core.view.asset).
			then(function() {
				$scope.$apply();
			});
	};

	$scope.switchToRevision = function(rev) {
		if (!angular.equals($scope.core.view.asset.revision, rev)) {
			$scope.core.view.asset.revision = rev;
			// $scope.core.view.asset.loadItems();
			$scope.core.loadItems($scope.core.view.asset).
				then(function() {
					$scope.$apply();
				});
		}
	}

	$scope.switchToAddAsset = function() {
		var current = $scope.core.view.newAsset;
		if (current == null) {
			current = new model.Asset(angular.extend(defaults.asset, {category: $scope.core.view.category, categoryName: $scope.core.view.name}), defaults);
			$scope.core.view.newAsset = current;
			$scope.core.assets.push(current);
		}

		console.log('addasset='+JSON.stringify(current));
		$scope.switchToAsset(current);
	};

	$scope.startEditingName = function(asset) {
		console.log('asset is'+JSON.stringify(asset));
		if (asset.nameEditing) {
			console.log('truth is the matter')
		} else {
			console.log('no one will WORKSTATION');
		}
		console.log('asset.editing ');
		asset.nameEditing = true;
		console.log("asset.editing "+ asset.nameEditing ? "true":"false");
	}

	$scope.stopEditingName = function(asset) {
		console.log('inside stopEditingName');
		asset.nameEditing = false;
	}

	$scope.addItem = function(item) {
		$scope.core.addItem(item)
	}

	$scope.scrapeItem = function(item) {
		if (!item.reference || item.reference === '') {
			return;
		}

		$scope.core.scrapeItem(item).
			then(function(reply) {
				$scope.$apply();
			})
	}

	$scope.formIsUnchanged = function(item) {
		return angular.equals(item, defaults.item);
	}

	$scope.saveAsset =  function(asset) {
		// var tmp = {"id":0,"hash":"a1e6a475","name":"singularity","nameEditing":false,"category":1,"categoryName":"WORKSTATION","created":1,"modified":1,"revisions":[],"revision":{"id":0,"index":0,"created":1,"items":[{"id":0,"hash":"663eb02b","name":"http://www.amazon.com/dp/B003O8J11E/","itemType":{"id":1,"name":"case"},"quantity":1,"price":"39.99"},{"id":0,"hash":"d0d4a192","name":"http://www.amazon.com/dp/B0056G10WK/","itemType":{"id":2,"name":"motherboard"},"quantity":1,"price":"88.24"},{"id":0,"hash":"34554d50","name":"http://www.amazon.com/dp/B0074J7ITG/","itemType":{"id":3,"name":"cpu"},"quantity":1,"price":"124.99"}]},"newItem":{"id":0,"hash":"ea058aa6","name":"","itemType":{"id":null,"name":""},"quantity":null,"price":null}};
		// var tmp = {"id":0,"hash":"a9754b00","name":"infinity","nameEditing":false,"category":1,"categoryName":"WORKSTATION","created":1,"modified":1,"revisions":[],"revision":{"id":0,"index":0,"created":1,"items":[{"hash":"c6ac7c42","product":{"id":0,"name":"Antec Fusion Remote Black Micro ATX Media Center / HTPC Case","asin":"B001BLWVJU","upc":"FUSION REMOTE BLACK","itemType":{"id":1,"name":"case"}},"reference":"http://www.amazon.com/Antec-Fusion-Remote-Black-Center/dp/B001BLWVJU/","quantity":1,"price":179.99},{"hash":"962c5767","product":{"id":0,"name":"Gigabyte GA-Z68XP-UD3 LGA 1155 Intel Z68 HDMI SATA 6Gb/s USB 3.0 ATX Intel Motherboard","asin":"B0054OWTQU","upc":"GA-Z68XP-UD3","itemType":{"id":2,"name":"motherboard"}},"reference":"http://www.amazon.com/Gigabyte-GA-Z68XP-UD3-1155-Intel-Motherboard/dp/B0054OWTQU/","quantity":1,"price":129.99},{"hash":"061df4db","product":{"id":0,"name":"Antec EarthWatts EA-650 Green 650 Watt 80 PLUS BRONZE Power Supply","asin":"B004NBZAES","upc":"EA-650 Green","itemType":{"id":3,"name":"cpu"}},"reference":"http://www.amazon.com/Antec-EarthWatts-EA-650-Green-BRONZE/dp/B004NBZAES/","quantity":1,"price":79.24}]},"newItem":{"hash":"7c9ce4ca","product":{"id":0,"name":"","asin":"","upc":"","itemType":{"id":0,"name":""}},"reference":"","quantity":null,"price":null}};
		$scope.core.saveAsset(asset).
			then(function() {
				console.log('got back from save: ');
				$scope.$apply();
			})
	}

	// Local filters
	$scope.assetFilter = function(asset) {
		return asset.category === $scope.core.view.category;
	}
};
