function startChat(pattern,userName){
	ws = new WebSocket("ws://localhost:8080/"+pattern)
	//入室時の発言
	ws.addEventListener("open",function(e){
		console.log("WebSocket connected")
		var data = {}
		data["name"] = "System"
		data["payload"] = userName+" Joined!"
		data["time"] = "" //クライアントの送信した時間ではなくサーバに届いた時間とする．
		ws.send(
			JSON.stringify(data)
		)
	});
	//発言を受信したらlistに表示
	ws.addEventListener("message",function(e){
		json = e.data
		msg = JSON.parse(json)
		console.log(msg)
		var user=msg["name"],payload=msg["payload"],time=msg["time"]
		var li = document.createElement("li");
		li.textContent=user+" : "+payload+" ("+time+")"
			document.getElementById("list").appendChild(li);
	});
	//boxに書いた内容を発言
	document.getElementById("sendBtn").addEventListener("click",function(e){
		var box = document.getElementById("box")
		var data = {}
		data["name"] = userName
		data["payload"] = box.value
		data["time"] = ""
		ws.send(
			JSON.stringify(data)
		)
		box.value = ""
	});
}
function enter(){
	var roomName = document.getElementById("roomName").value;
	var userName = document.getElementById("userName").value;
	var entrance = document.getElementById("entrance");
	var chat     = document.getElementById("chat");
	if (roomName !== "" && userName !== ""){
		//displayの値で画面を切り替え
		entrance.style.display="none";
		chat.style.display="block";
		//エントリー用のハンドラに接続
		ws = new WebSocket("ws://localhost:8080/entry")
		ws.addEventListener("open",function(e){
			var data = {}
			data["pattern"] = roomName
			data["name"] = userName
			ws.send(
				JSON.stringify(data)
			)
		});
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
