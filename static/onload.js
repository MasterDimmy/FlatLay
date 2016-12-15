 function ready() {
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
		
		$.getJSON("/get_field?group="+document.getElementById("group").value+"&width="+document.getElementById("data").offsetWidth+"&height="+document.getElementById("data").offsetHeight, function(json) {
			console.log(json);
			if (json.Success == true) {
				//json.Items: 	Images, Square, MaxAvailable
				//images: [PosX, Pos, Path]
				
				while (document.getElementById("data").hasChildNodes()) {
					document.getElementById("data").removeChild(document.getElementById("data").lastChild);
				}

				document.getElementById("square").innerHTML = "Занятая площадь = "+json.Items.Square;
				document.getElementById("max_square").innerHTML = "Максимально возможная = " + json.Items.MaxAvailable;
				json.Items.Images.forEach(function(item, i) {
					console.log("item=",item);	
					var image=document.createElement("img");
					image.src="/get_image?path="+item.Path;
					var x = item.PosX + document.getElementById("data").getBoundingClientRect().left+4;
					var y = item.PosY + document.getElementById("data").getBoundingClientRect().top+4;
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
 