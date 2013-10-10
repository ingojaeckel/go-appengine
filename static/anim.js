function move(circle, direction) {
	const speed = 5.0;
	
	if ("r" == direction) {
		circle.setX(circle.getX() + 1 * speed);
	} else if ("l" == direction) {
		circle.setX(circle.getX() - 1 * speed);
	} else if ("u" == direction) {
		circle.setY(circle.getY() - 1 * speed);
	} else if ("d" == direction) {
		circle.setY(circle.getY() + 1 * speed);
	}
	
	circle.getLayer().draw();
	
	moves++;
	
	if (moves % moveThreshold == 0) {
		// Don't send every move to the server.
		$.get("/rest/move/" + uuid + "/" + circle.getX() + "/" + circle.getY() + "/");
	} else if (moves % notifyThreshold == 0) {
		$.ajax({
			url: "/rest/notify",
			type: "POST",
			data: JSON.stringify({
				ID: uuid,
				X: circle.getX(),
				Y: circle.getY(),
				Recipients: ["player1", "player2", "player3"]
			})			
		});
	}
}

function registerKeyEvents(moveFn, circle) {
	$(document).keydown(function(e) {
		switch (e.which) {
		case 37:
			moveFn(circle, "l")
			break;
		case 38:
			moveFn(circle, "u")
			break;
		case 39:
			moveFn(circle, "r")			
			break;
		case 40:
			moveFn(circle, "d")
			break;
		}
	});
}

function addPlayers(players, layer) {
	for (var i=0; i<players.length; i++) {
		updatePlayer(players[i], layer);
	}
}

function updatePlayer(player, layer) {
	if (uuid == player.ID) {
		return; // don't update yourself.
	}
	
	model[player.ID] = player;
	
	if (view[player.ID]) {
		// We alreay know this player. Update view.
		view[player.ID].setX(player.P.X);
		view[player.ID].setY(player.P.Y);
	} else {
		// We don't know this player yet. Create view.
		var newCircle = new Kinetic.Circle({
			x: player.P.X,
			y: player.P.Y,
			radius: 10,
			fill: 'black',
			stroke: 'black',
			strokeWidth: 2
		});
		view[player.ID] = newCircle;
		layer.add(newCircle);
		updatePlayerList();
	}
	
	layer.draw();
}

function updatePlayerList() {
	$("#players").empty();
	
	var html = "";
	
	for (var key in model) {
		html += "<li>" + model[key].ID + " " + model[key].Name + "</li>";
	}
	
	$("#players").append("<ol>" + html + "</ol>");
}

function setupChannelApi(token, layer) {
	var channel = new goog.appengine.Channel(token);
	channel.open({
		'onopen': function() {
			console.log("opened");
		},
		'onmessage': function(message) {
			console.log("onmessage");
			
			// Unmarshal channel message
			var data = JSON.parse(message.data);
			
			if (0 == data[0]) {
				// add player
				var player = { ID: data[1], P: { X: 100, Y: 100 } };
				updatePlayer(player, layer)				
			} else if (1 == data[0]) {
				// remove player
				delete model[data[1]];
				if (view[data[1]] != null) {
					view[data[1]].remove();
					layer.draw();
					updatePlayerList();
				}
			} else {
				var player = { ID: data[0], P: { X: data[1], Y: data[2] } };
				updatePlayer(player, layer)
			}			
		},
		'onerror': function(error) {
			console.log("error");
		},
		'onclose': function() {
			console.log("closed");
		}
	});
}


