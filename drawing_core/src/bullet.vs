// #version 330 core
// layout(location = 0) in vec3 aPos;
// layout(location = 1) in vec3 aNormal;

// out vec3 vNormal;
// out vec3 vPos;
// uniform mat4 projection,view,model;

// void main(){
//     vec4 worldPos=model*vec4(aPos,1.0);
//     vPos = worldPos.xyz;
//     vNormal=mat3(transpose(inverse(model)))*aNormal;
// }

#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out VS_OUT {
    vec3 color;
} vs_out;

uniform float dis;

void main()
{
    vs_out.color = aColor;
    vec3 pos=aPos + dis*aColor;
    pos+=vec3(0.0,-0.05,0.0)*abs(dis)*10;
    gl_Position=vec4(pos.x,pos.y,pos.z,1.0);
    // gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0); 
}