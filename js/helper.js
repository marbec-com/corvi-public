var ClearObject = function(obj) {
	for (var key in obj){
		if (obj.hasOwnProperty(key)){
			delete obj[key];
		}
	}
};

var ClearObjectArray = function(obj) {
	for (var key in obj){
		if (obj.hasOwnProperty(key)){
			obj[key].length = 0;
		}
	}
};

var StartsWith = function(str, prefix) {
	return (str.indexOf(prefix) === 0)
};