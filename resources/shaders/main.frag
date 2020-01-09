#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
flat in vec2 v_r;
in vec4 v_p;
in float dist;
out vec4 outColor;

uniform float sectorLight;
uniform int draw_phase;


// TODO: use this in combination with the distance.
vec3 saturation(vec3 rgb, float adjustment)
{
    // Algorithm from Chapter 16 of OpenGL Shading Language
    const vec3 W = vec3(0.2125, 0.7154, 0.0721);
    vec3 intensity = vec3(dot(rgb, W));
    return mix(intensity, rgb, adjustment);
}

void main()
{
    if (draw_phase == 3) {
      vec2 uv = vec2(v_p.x, v_p.y) / v_p.w/1.6 * vec2(1, -1);
      uv = vec2(uv.x - 2.0 * v_r.x / 3.14159265358, uv.y + v_r.y);
      uv -= 0.3;
      outColor = texture(tex, uv);
      return;
    }
    float alpha = texture(tex, fragTexCoord).a;
    if (alpha == 1.0) {
      float lighting = 1;
      lighting = sectorLight/400;
      float darken = 0.2 - clamp(dist,0,1000)/1500;
      
      if (draw_phase != 2 && sectorLight < 160) {
        lighting += darken;
      }
      outColor = texture(tex, fragTexCoord) * vec4(vec3(lighting),1);
      if (draw_phase != 2 && sectorLight < 160) {
        outColor.rgb = saturation(outColor.rgb,1.0 - clamp(dist,0,1000)/1000).rgb;
      }
    } else {
      discard;
    }
}
