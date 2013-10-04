var uuid = 0;
var circle = null;
var model = {}; // maps UUIDs to player data
var view = {};  // maps UUIDS to UI elements

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

	var simpleText1 = new Kinetic.Text({
		x: 5,
		y: stage.getHeight() - 45,
		text: '[1] Player1: 10',
		fontSize: 12,
		fontFamily: 'Calibri',
		fill: 'green'
	});

	stage.add(layer);

	$("#joinBtn").button().click(function(event) {
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
				setupChannelApi(response.ChannelToken);
				// pollInTheBackground(layer);
			}
		});
	});
	$("#moveBtn").button().click(function(event) {
		$.ajax({
			url: "/rest/move/" + uuid + "/150/150"
		});
	});
	$("#pollBtn").button().click(function(event) {
		poll(layer);
	});

	// f(circle1, stage, layer, 1);
});
