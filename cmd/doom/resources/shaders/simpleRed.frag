#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
out vec4 outColor;

uniform float sectorLight;

void main()
{
  fragTexCoord;
  outColor = vec4(1,0,0,1);
}