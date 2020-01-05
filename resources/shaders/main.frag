#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
flat in vec2 v_r;
in vec4 v_p;

out vec4 outColor;

uniform float sectorLight;
uniform int draw_phase;

void main()
{
    if (draw_phase == 3) {
      vec2 uv = vec2(v_p.x, v_p.y) / v_p.w * vec2(1, -1);
      uv = vec2(uv.x - 4.0 * v_r.x / 3.14159265358, uv.y + 1.0 + v_r.y);
      outColor = texture(tex, uv);
      return;
    }

    float alpha = texture(tex, fragTexCoord).a;
    if (alpha == 1.0) {
	  outColor = texture(tex, fragTexCoord) * vec4(sectorLight/300,sectorLight/300,sectorLight/300,1);
    } else {
      discard;
    }
}
