var uuid = 0;

function f(circle1, stage, layer, speed) {
	var initialDirection = "r";
	var direction = initialDirection;
	const b = 20; // size of the border
	var anim1 = new Kinetic.Animation(function(frame) {
		if ("r" == direction) {
			var canContinue = (circle1.getX() + b < stage.getWidth());

			if (canContinue) {
				circle1.setX(circle1.getX() + 1 * speed);
				return;
			}

			if (circle1.getY() + b < stage.getHeight()) {
				direction = "d";
				circle1.setY(circle1.getY() + 1 * speed);
			} else if (circle1.getX() - b > 0) {
				direction = "l";
				circle1.setX(circle1.getX() - 1 * speed);
			} else if (circle1.getY() - b > 0) {
				direction = "u";
				circle1.setY(circle1.getY() - 1 * speed);
			} else {
				console.debug("stuck");
			}
		} else if ("d" == direction) {
			var canContinue = (circle1.getY() + b < stage.getHeight());

			if (canContinue) {
				circle1.setY(circle1.getY() + 1 * speed);
				return;
			}

			if (circle1.getX() - b > 0) {
				direction = "l";
				circle1.setX(circle1.getX() - 1 * speed);
			} else if (circle1.getY() - b > 0) {
				direction = "u";
				circle1.setY(circle1.getY() - 1 * speed);
			} else if (circle1.getX() + b < stage.getWidth()) {
				direction = "r";
				circle1.setX(circle1.getX() + 1 * speed);
			} else {
				console.debug("stuck");
			}
		} else if ("l" == direction) {
			var canContinue = (circle1.getX() - b > 0);

			if (canContinue) {
				circle1.setX(circle1.getX() - 1 * speed);
				return;
			}

			if (circle1.getY() - b > 0) {
				direction = "u";
				circle1.setY(circle1.getY() - 1 * speed);
			} else if (circle1.getX() + b < stage.getWidth()) {
				direction = "r";
				circle1.setX(circle1.getX() + 1 * speed);
			} else if (circle1.getY() - b > 0) {
				direction = "u";
				circle1.setY(circle1.getY() - 1 * speed);
			} else {
				console.debug("stuck");
			}
		} else if ("u" == direction) {
			var canContinue = (circle1.getY() - b > 0);

			if (canContinue) {
				circle1.setY(circle1.getY() - 1 * speed);
				return;
			}

			if (circle1.getX() + b < stage.getWidth()) {
				direction = "r";
				circle1.setX(circle1.getX() + 1 * speed);
			} else if (circle1.getY() + b > 0) {
				direction = "d";
				circle1.setY(circle1.getY() + 1 * speed);
			} else if (circle1.getX() - b > 0) {
				direction = "l";
				circle1.setX(circle1.getX() - 1 * speed);
			} else {
				console.debug("stuck");
			}
		}
	}, layer);

	anim1.start();
}

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
