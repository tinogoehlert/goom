#version 330

in vec3 vertex;
in vec2 vertTexCoord;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform mat4 ortho;
uniform int draw_phase;

uniform vec3 player_pos;

uniform int billboard_flipped;
uniform vec3 billboard_pos;
uniform vec2 billboard_size;

out vec2 fragTexCoord;
out float light;
flat out vec2 v_r;
out vec4 v_p;
out float dist;


vec4 drawSky() {
	mat4 transform = projection * model * view;
    vec4 forward = transform[2];
    v_r = vec2(atan(forward.x, forward.z), 1);
	vec4 projected_pos = transform * vec4(vertex, 1);
    v_p = projected_pos;
	
    return projected_pos;
}

vec4 drawBillboard() {
	vec3 particleCenter_wordspace = billboard_pos;
	vec3 CameraRight_worldspace = vec3(view[0][0], view[1][0], view[2][0]);
	vec3 CameraUp_worldspace = vec3(view[0][1], view[1][1], view[2][1]);
	vec3 vertexPosition_worldspace = 
		particleCenter_wordspace
		+ CameraRight_worldspace * vertex.x * billboard_size.x
		+ CameraUp_worldspace * vertex.y * billboard_size.y;

	return projection*view * vec4(vertexPosition_worldspace, 1.0f);
}

vec4 DrawHUD() {
	vec3 v = vertex;
	v.x *= billboard_size.x;
	v.y *= billboard_size.y;
	return ortho * vec4(v+billboard_pos, 1.0f);
}

void main()
{
	dist = 1000;
	fragTexCoord = vertTexCoord;
	// things code
	if (draw_phase == 1) {
		if (billboard_flipped == 1) {
			fragTexCoord.x = -fragTexCoord.x;
		}
		gl_Position = drawBillboard();	
		dist = abs(distance(player_pos,billboard_pos));
		return;
	}
	// hud code
	if (draw_phase == 2) {
		fragTexCoord.x = -fragTexCoord.x;
		gl_Position = DrawHUD();	
		return;
	} 
	if (draw_phase == 3) {
		gl_Position = drawSky();
		return;
	}
	dist = abs(distance(player_pos,vertex));
	gl_Position = projection * view  * model * vec4(vertex, 1.0);
}
