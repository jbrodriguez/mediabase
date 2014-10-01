(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('options', options);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function options() {
    	var searchTerm = '';

        var filterByOptions = [{value: 'title', display: 'Title'}, {value: 'genre', display: 'Genre'}];
        var filterBy = '';

        var sortByOptions = [{value: 'title', display: 'Title'}, {value: 'runtime', display: 'Runtime'}, {value: 'added', display: 'Added'}, {value: 'last_watched', display: 'Watched'}, {value: 'year', display: 'Year'}, {value: 'imdb_rating', display: 'Rating'}];
        var sortBy = '';

        var sortOrderOptions = ['asc', 'desc'];
        var sortOrder = 'desc';

    	var service = {
            searchTerm: searchTerm,
            filterByOptions: filterByOptions,
            filterBy: filterBy,
            sortByOptions: sortByOptions,
            sortBy: sortBy,
            sortOrderOptions: sortOrderOptions,
            sortOrder: sortOrder
    	};

    	return service;
    }

})();