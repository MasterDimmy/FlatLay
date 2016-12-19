 var last = 0;
 
 function cumulativeOffset(element) {
	var top = 0, left = 0;
	do {
		top += element.offsetTop  || 0;
		left += element.offsetLeft || 0;
		element = element.offsetParent;
	} while(element);

	return {
		top: top,
		left: left
	};
};

function imgclick () {
	alert(this.src);
};
					
 function ready() { 
 		while (document.getElementById("data").hasChildNodes()) {
			document.getElementById("data").removeChild(document.getElementById("data").lastChild);
		}
 
		var cur = last;
		last++;
		var total_groups = 0;
		$.getJSON("/get_total_groups", function(json) {
			console.log(json);
			if (json.Success == true) {
				total_groups = json.Items;
				document.getElementById("total_groups").innerHTML = total_groups;
				console.log("total_groups = ",total_groups);				
			} else {
				console.log("error: "+json.Message);
				alert(json.Message);
			}
		});
		
		$.getJSON("/get_field?group="+document.getElementById("group").value+"&width="+(document.getElementById("data").offsetWidth-12)+"&height="+(document.getElementById("data").offsetHeight-12), function(json) {
			console.log(json);
			if (json.Success == true) {
				if (cur!=last-1) {
					console.log("cur=",cur," last=",last);
					return; //был запрос новее
				}
				//json.Items: 	Images, Square, MaxAvailable
				//images: [PosX, Pos, Path]
				
				document.getElementById("square").innerHTML = "Taken square = "+json.Items.Square;
				document.getElementById("max_square").innerHTML = "Max possible = " + json.Items.MaxAvailable;
				json.Items.Images.forEach(function(item, i) {
					console.log("item=",item);	
					var image=document.createElement("img");
					image.src="/get_image?path="+item.Path;
					image.onclick = imgclick;
					
					var offset = cumulativeOffset(document.getElementById("data"));
					var x = item.PosX + offset.left+4;
					var y = item.PosY + offset.top+4;

					image.style = "position:absolute;left:"+x+"px;top:"+y+"px;";
					document.getElementById("data").appendChild(image);
				});

			} else {
				console.log("error: "+json.Message);
				alert(json.Message);
			}
		});
 }

 document.addEventListener("DOMContentLoaded", ready);

 window.onresize = ready;
 
$("#group").change(ready);
 
$('body').css('top', -(document.documentElement.scrollTop) + 'px').addClass('noscroll');
