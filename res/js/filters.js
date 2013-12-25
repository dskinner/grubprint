angular.module("filters", []).
	filter("momentFromNow", function() {
		return function(input) {
			var s = moment(input).fromNow();
			// strip " ago" from end of result
			return s.slice(0, s.length-4);
		};
	});
