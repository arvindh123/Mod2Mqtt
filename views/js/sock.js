url = 'ws://localhost:8080/ws';
c = new WebSocket(url);

document.getElementById("send").onclick = function() {myFunction()};

function myFunction() {
    send(document.getElementById("msg").value)

}

send = function(data){
    $("#output").append((new Date())+ " ==> "+data+"\n")
    c.send(data)
}

c.onmessage = function(msg){
    $("#output").append((new Date())+ " <== "+msg.data+"\n")
    console.log(msg)
}
