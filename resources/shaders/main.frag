#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
out vec4 outColor;

uniform float sectorLight;

void main()
{

    float alpha = texture(tex, fragTexCoord).a;
    if (alpha == 1.0) {
	  outColor = texture(tex, fragTexCoord) * vec4(sectorLight/300,sectorLight/300,sectorLight/300,1);
    } else {
      discard;
    }
}