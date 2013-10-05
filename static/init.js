var uuid = 0;
var circle = null;
var model = {}; // maps UUIDs to player data
var view = {};  // maps UUIDS to UI elements
var moves = 0;
// Every moveThreshold-th move will be sent to the server.
// Set to 1 to send each move to the server or higher than 1 to skip moves.
var moveThreshold = 2;

$(function() {
	var stage = new Kinetic.Stage({
		container: 'container',
		width: 600,
		height: 400
	});
	var layer = new Kinetic.Layer();

	var rect = new Kinetic.Rect({
		x: 0,
		y: 0,
		width: stage.getWidth(),
		height: stage.getHeight(),
		stroke: 'black',
		strokeWidth: 2
	});

	// add the shape to the layer
	layer.add(rect);
	stage.add(layer);

	$.ajax({
		url: "/rest/join",
		dataType: "json",
		success: function(response) {
			circle = new Kinetic.Circle({
				x: 100,
				y: 100,
				radius: 10,
				fill: 'red',
				stroke: 'black',
				strokeWidth: 2
			});
			
			layer.add(circle);
			layer.draw();
				
			uuid = response.UUID;
				
			registerKeyEvents(move, circle);
			setupChannelApi(response.ChannelToken, layer);
		}
	});
});
