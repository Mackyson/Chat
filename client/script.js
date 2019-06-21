function startChat(pattern,userName){
	ws = new WebSocket("ws://localhost:8080/"+pattern)
	ws.addEventListener("open",function(e){
			console.log("WebSocket connected")
	});
	ws.addEventListener("message",function(e){
		msg = e.data
		console.log(msg)
		var li = document.createElement("li");
		li.textContent=msg;
			document.getElementById("list").appendChild(li);
	});
	document.getElementById("sendBtn").addEventListener("click",function(e){
		ws.send(
			document.getElementById("box").value
		)
	});
}
function enter(){
	var roomName = document.getElementById("roomName").value;
	var userName = document.getElementById("userName").value;
	var entrance = document.getElementById("entrance");
	var chat = document.getElementById("chat");
	if (roomName !== "" && userName !== ""){
		entrance.style.display="none";
		chat.style.display="block";
		ws = new WebSocket("ws://localhost:8080/entry")
		ws.addEventListener("open",function(e){ws.send(roomName)});
		ws.addEventListener("message",function(e){
		 if (e.data === "ok"){
			 ws.close();
			startChat(roomName,userName)
		 }
		});
	}
	else{
		entrance.insertAdjacentHTML("beforeend","<br>your input is null")
	}
}
