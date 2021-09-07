#version 400
layout(location = 0) in vec4 vertex_position;
layout(location = 1) in vec4 vertex_color;

uniform mat4 camera;
out vec4 theColor;

void main() {
	gl_Position = camera * vec4(vertex_position.x, vertex_position.y, vertex_position.z, 1.0);
	theColor = vertex_color;
}