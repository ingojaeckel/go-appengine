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
	
	$.ajax({url: "/rest/move/" + uuid + "/" + circle.getX() + "/" + circle.getY()})
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

function pollInTheBackground(layer) {
	// console.log("pollInTheBackground");
	
	setTimeout(function() {
		poll(layer);
		pollInTheBackground(layer);
	}, 200);
}

function poll(layer) {
	// console.log("poll");
	$.ajax({
		url: "/rest/poll",
		dataType: "json",
		success: function(response) {
			for (var i=0; i<response.Players.length; i++) {
				if (uuid == response.Players[i].ID) {
					continue; // don't draw yourself
				}
				
				// Create/Update model for this player
				model[response.Players[i].ID] = response.Players[i];
				
				if (view[response.Players[i].ID]) {
					// We alreay know this player. Update view.
					view[response.Players[i].ID].setX(response.Players[i].P.X);
					view[response.Players[i].ID].setY(response.Players[i].P.Y);
					// console.log("updated player " + response.Players[i].ID);
				} else {
					// We don't know this player yet. Create view.
					var newCircle = new Kinetic.Circle({
						x: response.Players[i].P.X,
						y: response.Players[i].P.Y,
						radius: 10,
						fill: 'black',
						stroke: 'black',
						strokeWidth: 2
					});
					view[response.Players[i].ID] = newCircle;
					layer.add(newCircle);
					// console.log("created player " + response.Players[i].ID);						
				}
			}
			layer.draw();
		}
	});
}

function setupChannelApi(token) {
	console.log("setupChannelApi " + token);
	var c = new goog.appengine.Channel(token);
	var s = c.open({
		onopen: function(a) {
			console.log("opened");
		},
		onmessage: function(a) {
			console.log("message");			
		},
		onerror: function(a) {
			console.log("error");
		},
		onclose: function(a) {
			console.log("closed");
		}
	});
}


