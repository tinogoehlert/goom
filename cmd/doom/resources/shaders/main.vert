#version 330

in vec3 vertex;
in vec2 vertTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform mat4 ortho;
uniform int is_billboard;
uniform vec3 billboard_pos;
uniform vec2 billboard_size;


out vec2 fragTexCoord;
out float light;

void main()
{
	// hud code
	fragTexCoord = vertTexCoord;
	if (is_billboard == 2) {
		vec3 v = vertex;
		v.x *= billboard_size.x;
		v.y *= billboard_size.y;
		gl_Position = ortho * vec4(v+billboard_pos, 1.0f);
		return;
	}
	// things code
	if (is_billboard == 1) {
		vec3 particleCenter_wordspace = billboard_pos;
		vec3 CameraRight_worldspace = vec3(view[0][0], view[1][0], view[2][0]);
		vec3 CameraUp_worldspace = vec3(view[0][1], view[1][1], view[2][1]);
		vec3 vertexPosition_worldspace = 
			particleCenter_wordspace
			+ CameraRight_worldspace * vertex.x * billboard_size.x
			+ CameraUp_worldspace * vertex.y * billboard_size.y;
		
		gl_Position = projection*view * vec4(vertexPosition_worldspace, 1.0f);
		return;
	} 
	gl_Position = projection * view  * model * vec4(vertex, 1.0);
}