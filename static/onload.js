 function ready() {
		console.log("page load");
		$.getJSON("/get_data?width="+$(window).height()+"&height="+$(window).width(), function(json) {
			console.log(json); // this will show the info it in firebug console
		});
 }

 document.addEventListener("DOMContentLoaded", ready);
