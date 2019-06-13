ws = new WebSocket("ws://localhost:8080/echo")

ws.addEventListener("open",function(e){
		console.log("接続")
});
ws.addEventListener("message",function(e){
	msg = e.data
	console.log(msg)
	var li = document.createElement("li");
	li.textContent=msg;
	document.getElementById("list").appendChild(li);
});
document.addEventListener("DOMContentLoaded",function(e){
	document.getElementById("sendBtn").addEventListener("click",function(e){
		ws.send(
			document.getElementById("box").value
		)
	});
});
