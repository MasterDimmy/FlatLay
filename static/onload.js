 function ready() {
		console.log("page load");
		$.getJSON("/get_data?width="+document.getElementById("data").offsetWidth+"&height="+document.getElementById("data").offsetHeight, function(json) {
			console.log(json);
			if (json.Success == true) {
				console.log("ok");
				
			} else {
				console.log("error: "+json.Message);
				alert(json.Message);
			}
		});
 }

 document.addEventListener("DOMContentLoaded", ready);

 window.onresize = ready;
 