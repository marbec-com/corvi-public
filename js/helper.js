var ClearObject = function(obj) {
	for (var key in obj){
		if (obj.hasOwnProperty(key)){
			delete obj[key];
		}
	}
};

var StartsWith = function(str, prefix) {
	return (str.indexOf(prefix) === 0)
};