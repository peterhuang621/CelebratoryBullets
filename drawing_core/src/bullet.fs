// #version 330 core
// in vec3 fNormal;
// out vec4 FragColor;

// void main(){
//     float brightness = dot(normalize(fNormal), normalize(vec3(0.0, 1.0, 0.0)));
//     brightness = clamp(brightness, 0.2, 1.0);
//     FragColor = vec4(vec3(0.5, 0.8, 1.0) * brightness, 1.0);
// }
#version 330 core
out vec4 FragColor;

in vec3 fColor;
uniform float dis;

void main()
{
    FragColor = vec4(fColor, 1.0);   
}